package conn

import (
	"bufio"
	"errors"
	"fmt"
	"net"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

const CHANNEL_BUFFER = 5

type HttpConn struct {
	conn      net.Conn
	route     router.Route
	log       logger.Logger // logger within httpconn
	timeoutMs int           // default timeout at 2s
}

func NewConn(c net.Conn, r *router.Router) *HttpConn {
	return &HttpConn{
		c,
		r,
		logger.New("HttpConn"),
		2000,
	}
}

// create connection handler that read from conn and write to conn when get response from router
func (c *HttpConn) HandleConn() {
	reqQueue := make(chan *routerContext, CHANNEL_BUFFER) // set buffer size to not block read
	resQueue := make(chan *routerContext, CHANNEL_BUFFER) // set buffer size to not block read

	// pipelining
	go read(c.conn, reqQueue, c.log)             // read message and parse request from fd
	go route(c.route, reqQueue, resQueue, c.log) // gets request from queue then pass to writer
	go write(c.conn, resQueue, c.log)            // convert responses to msg and write to fd
}

// routerContext implements the [context.Context] interface for passing in message info across goroutines
type routerContext struct {
	Request   msg.Request   // request from read()
	Response  msg.Response  // response from [RouterAction]
	CloseConn bool          // should close connection after sending response
	Done      chan struct{} // response is ready to be sent to writer
}

func buildContext(req msg.Request) *routerContext {
	// parse header and create result accordingly
	return &routerContext{
		Request:   req,
		CloseConn: false,
		Done:      make(chan struct{}),
	}
}

func route(r router.Route, source <-chan *routerContext, dest chan<- *routerContext, log logger.Logger) {
	isSafePipeline := true
	// for routing and closing [routerContext.Done] channel
	goRoute := func(rc *routerContext) {
		defer close(rc.Done)
		rc.Response = r.To(rc.Request)
	}

	for rc := range source {
		log.Info(fmt.Sprintf("request: %s %s", rc.Request.GetMethod(), rc.Request.GetPath()))
		dest <- rc
		// if request is safe method then pipeline
		if rc.Request.IsSafeMethod() && isSafePipeline {
			go goRoute(rc) // send request to route
		} else {
			isSafePipeline = false
			goRoute(rc) // stop fan-out, since unsafe method shouldn't be parallelized
		}
	}
}

// read convert raw message into [Request] and passes it to [Router.To()].
func read(conn net.Conn, queue chan<- *routerContext, l logger.Logger) {
	reqBuilder := msg.NewHttpReqBuilder() // for building requests
	scanner := bufio.NewScanner(conn)     // to create a scanner
	scanner.Split(bufio.ScanBytes)        // scan by a sequence of octet
	line := ""                            // HTTP message is separated by CRLF until reaches HTTP body
	defer close(queue)                    // close channel in queue to signify write() to stop sending response
	var err error

	for scanner.Scan() {
		char := scanner.Text() // convert the byte into string

		switch reqBuilder.State() {
		case msg.RequestLineBuildState: // parse request-line
			if char != "\n" {
				line += char
			} else {
				err = reqBuilder.AddRequestLine(line) // parse into [HttpRequest] line by line
				line = ""
			}
		case msg.HeadersBuildState: // parse headers
			if char != "\n" {
				line += char
			} else { // if not an empty line then is a header
				err = reqBuilder.AddHeader(line)
				line = ""
			}
		case msg.BodyBuildState: // TODO: add validation for not being request-line
			line += char
			err = reqBuilder.AddBody(line)
		default:
			err = errors.New("unrecognized state when building request")
		}

		// error from bad request (message in wrong format)
		if err != nil {
			l.Warn(err.Error())

			ctx := buildContext(nil)
			queue <- ctx // send a bad request status response
			break        // stop reading if msg is in wrong format
		}

		// check at every byte scan since after body there is no token that signifies ending of message
		if reqBuilder.State() == msg.CompleteBuildState {
			req, err := reqBuilder.Build() // build request & reset builder
			line = ""
			if err != nil { // if an error when build, closes connection
				break
			}
			ctx := buildContext(req)
			queue <- ctx
		}
	}
	// only terminates read after
}

// write send HTTP response back to the client in-order of the request
func write(conn net.Conn, source <-chan *routerContext, l logger.Logger) {
	writer := bufio.NewWriter(conn)
	for routerCtx := range source { // will break from loop when channel is closed
		<-routerCtx.Done // block to execute in order

		// send HTTP response
		res := routerCtx.Response
		if res == nil {
			res = msg.CreateHttpResponse(msg.StatusServerInternalError)
		}
		msg, err := res.String()
		if err != nil {
			l.Fatal(fmt.Sprintf("problems with stringify response %s", err.Error()))
		} else {
			if _, err := writer.WriteString(msg); err != nil { // write the [HttpResponse] to buffer
				l.Warn(err.Error())
				break
			} else {
				writer.Flush() // respond to client (write to socket)
			}
		}
	}

	if err := conn.Close(); err != nil { // close connection
		l.Warn(err.Error())
	}
	l.Debug("close connection")
}
