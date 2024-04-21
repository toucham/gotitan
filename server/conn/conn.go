package conn

import (
	"bufio"
	"net"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

const CHANNEL_BUFFER = 5

// for managing TCP connection to align with HTTP/1.1
type HttpConn struct {
	conn    net.Conn                  // TCP connection
	timeout int32                     // connection timeout in ms
	channel chan *router.RouterResult // queue response to return in correct order
	route   router.Route
	logger  logger.Logger
}

// create connection manager
func HandleConn(conn net.Conn, r router.Route, timeout int32) *HttpConn {
	queue := make(chan *router.RouterResult, CHANNEL_BUFFER) // set buffer size to not block read
	logger := logger.New("HttpConn")
	return &HttpConn{conn, timeout, queue, r, logger}
}

// Read parse app-layer message into [HttpRequest] and execute [Router.To()].
// Have added support for pipelining; however, most browser doesn't support pipelining.
func (c *HttpConn) Read() {
	scanner := bufio.NewScanner(c.conn)
	scanner.Split(bufio.ScanBytes)
	state := msg.REQUESTLINE_BS
	req := msg.NewRequest()
	line := ""
	for scanner.Scan() { // keep scanning TCP connection fd for persistent connection
		char := scanner.Text()
		switch state {
		case msg.REQUESTLINE_BS:
			if char != "\n" {
				line += char
			} else {
				err := req.AddRequestLine(line) // parse into [HttpRequest] line by line
				if err != nil {                 // if parsed fail then discard
					c.logger.Warn(err.Error())
				} else {
					state = msg.HEADERS_BS
					line = ""
				}
			}
		case msg.HEADERS_BS:
			if char != "\n" {
				line += char
			} else if line == "" { // ending headers
				if req.Headers.ContentLength > 0 { // expect body
					state = msg.BODY_BS
				} else {
					state = msg.COMPLETE_BS
				}
			} else { // if not an empty line then is a header
				err := req.AddHeader(line)
				if err != nil { // if parsed fail then discard
					c.logger.Warn(err.Error())
				}
				line = ""
			}
		case msg.BODY_BS:
			if len(line) < req.Headers.ContentLength {
				line += char
			} else if len(line) == req.Headers.ContentLength {
				req.AddBody(line)
			} else {
				c.logger.Warn("length of body is longer than content-length")
			}
		default:
			c.logger.Warn("Unrecognized state during building request")
		}

		// check at every byte scan since after body there is no token that signifies ending of message
		if state == msg.COMPLETE_BS {
			result := router.CreateResult(req)
			c.channel <- result
			go c.route.To(req, result) // send request to route
			req = msg.NewRequest()     // refresh new request
			state = msg.REQUESTLINE_BS
		}
	}
}

// Write send HTTP response back to the client in-order of the request
func (c *HttpConn) Write() {
	writer := bufio.NewWriter(c.conn)
	for result := range c.channel {
		<-result.Ready // block to execute in order
		// TODO: add logic to do smth from fields in [RouterResult] (ex: close connection)

		// send HTTP response
		res := result.Response
		if res, err := res.String(); err != nil {
			c.logger.Warn(err.Error())
		} else {
			if _, err := writer.WriteString(res); err != nil { // write the [HttpResponse] to buffer
				c.logger.Warn(err.Error())
				return
			} else {
				c.logger.Debug("Respond to client")
				writer.Flush() // respond to client (write to socket)
			}
		}
	}
}
