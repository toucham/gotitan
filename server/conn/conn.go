package conn

import (
	"bufio"
	"errors"
	"net"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

const CHANNEL_BUFFER = 5

// for managing TCP connection to align with HTTP/1.1
type HttpConn struct {
	conn           net.Conn                   // TCP connection
	timeout        int32                      // connection timeout in ms
	queue          chan *router.RouterContext // queue response to return in correct order
	isSafePipeline bool                       // determine whether there is only safe methods from all the received requests
	reqBuilder     msg.RequestBuilder         // builder for creating [Request]
	route          router.Route               // for routing request to correct [Action]
	logger         logger.Logger              // logger for [HttpConn]
}

// create connection manager
func HandleConn(conn net.Conn, r router.Route, timeout int32) *HttpConn {
	queue := make(chan *router.RouterContext, CHANNEL_BUFFER) // set buffer size to not block read
	logger := logger.New("HttpConn")
	req := msg.NewHttpReqBuilder()
	return &HttpConn{conn, timeout, queue, true, req, r, logger}
}

// Read convert raw message into [Request] and passes it to [Router.To()].
func (c *HttpConn) Read() {
	scanner := bufio.NewScanner(c.conn)
	scanner.Split(bufio.ScanBytes) // scan by a sequence of octet
	line := ""                     // HTTP message is separated by CRLF until reaches HTTP body
	var err error

	for scanner.Scan() {
		char := scanner.Text() // convert the byte into string

		switch c.reqBuilder.State() {
		case msg.RequestLineBuildState: // parse request-line
			if char != "\n" {
				line += char
			} else {
				err = c.reqBuilder.AddRequestLine(line) // parse into [HttpRequest] line by line
				line = ""
			}
		case msg.HeadersBuildState:
			if char != "\n" {
				line += char
			} else { // if not an empty line then is a header
				err = c.reqBuilder.AddHeader(line)
				line = ""
			}
		case msg.BodyBuildState: // TODO: add validation for not being request-line
			line += char
			err = c.reqBuilder.AddBody(line)
		default:
			err = errors.New("unrecognized state when building request")
		}

		// stop reading if msg is in wrong format
		if err != nil {
			c.logger.Warn(err.Error())
			break
		}

		// check at every byte scan since after body there is no token that signifies ending of message
		if c.reqBuilder.State() == msg.CompleteBuildState {
			req, err := c.reqBuilder.Build() // build request & reset builder
			line = ""
			if err != nil { // if an error when build, closes connection
				break
			}
			ctx := router.CreateContext(req)
			c.queue <- ctx
			if req.IsSafeMethod() && c.isSafePipeline {
				go c.route.To(req, ctx) // send request to route
			} else {
				c.isSafePipeline = false
				c.route.To(req, ctx) // stop reading on this connection
			}
		}
	}
	close(c.queue) // close channel + signal Write() to close channel
}

// Write send HTTP response back to the client in-order of the request
func (c *HttpConn) Write() {
	writer := bufio.NewWriter(c.conn)
	for result := range c.queue { // will break from loop when channel is closed
		<-result.Ready // block to execute in order

		// send HTTP response
		res := result.Response
		msg := res.String()
		if _, err := writer.WriteString(msg); err != nil { // write the [HttpResponse] to buffer
			c.logger.Warn(err.Error())
			break
		} else {
			c.logger.Debug("Respond to client")
			writer.Flush() // respond to client (write to socket)
		}
	}

	if err := c.conn.Close(); err != nil { // close connection
		c.logger.Warn(err.Error())
	}
	c.logger.Debug("close connection")
}
