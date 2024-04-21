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

// TODO: add validation so that it will reset parsing during pipelining

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
				}
				line = ""
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
					state = msg.REQUESTLINE_BS
				}
				line = ""
			}
		case msg.BODY_BS: // TODO: add validation for not being request-line
			if len(line) < req.Headers.ContentLength {
				line += char
			}
			if len(line) == req.Headers.ContentLength {
				err := req.AddBody(line)
				if err != nil { // if parsed fail then discard
					c.logger.Warn(err.Error())
					state = msg.REQUESTLINE_BS
				} else {
					state = msg.COMPLETE_BS
				}
				line = ""
			}
		default:
			c.logger.Warn("Unrecognized state during building request")
			line = ""
		}

		// check at every byte scan since after body there is no token that signifies ending of message
		if state == msg.COMPLETE_BS {
			result := router.CreateResult(req)
			c.channel <- result
			if req.IsSafeMethod() {
				go c.route.To(req, result) // send request to route
				req = msg.NewRequest()     // refresh new request
				state = msg.REQUESTLINE_BS
			} else {
				// stop reading on this connection
				c.route.To(req, result)
				break
			}
		}
	}
	close(c.channel)
}

// Write send HTTP response back to the client in-order of the request
func (c *HttpConn) Write() {
	writer := bufio.NewWriter(c.conn)
	for result := range c.channel {
		<-result.Ready // block to execute in order
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
	if err := c.conn.Close(); err != nil {
		c.logger.Warn(err.Error())
	}
}
