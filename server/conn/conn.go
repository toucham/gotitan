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
	route   *router.Router
	logger  *logger.Logger
}

// create connection manager
func HandleConn(conn net.Conn, r *router.Router, timeout int32) *HttpConn {
	queue := make(chan *router.RouterResult, CHANNEL_BUFFER) // set buffer size to not block read
	logger := logger.New("HttpConn")
	return &HttpConn{conn, timeout, queue, r, logger}
}

// Read parse app-layer message into [HttpRequest] and execute [Router.To()]
func (c *HttpConn) Read() {
	scanner := bufio.NewScanner(c.conn)
	scanner.Split(bufio.ScanLines)
	for { // keep scanning TCP connection fd for persistent connection
		req := new(msg.HttpRequest)
		scanner.Scan() // scan one line from the buffer filled by fd
		err := scanner.Err()
		if err != nil {
			c.logger.Warn(err.Error())
		} else {
			err = req.Next(scanner.Text()) // parse into [HttpRequest] line by line
			if err != nil {
				c.logger.Warn(err.Error())
				req = new(msg.HttpRequest)
			}
		}

		if req.IsReady() { // if ready then send to router
			result := router.CreateResult(req)
			c.channel <- result
			go c.route.To(req, result) // send request to route
			req = new(msg.HttpRequest) // refresh new request
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
