package server

import (
	"bufio"
	"fmt"
	"io"
	"net"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/msg"
	"github.com/toucham/gotitan/server/router"
)

type HttpServer struct {
	router.Router              // embedded Router
	ln            net.Listener // socket listener
	port          string
	reqMw         []ReqMiddleware
	logger        *logger.Logger
}

func Init(host string, port string) *HttpServer {
	logger := logger.New("Server")
	addr := fmt.Sprintf("%s:%s", host, port)
	ln, err := net.Listen("tcp", addr)
	if err != nil {
		panic(err)
	}

	logger.Info("Listening on address: %s", addr)

	// An array of map for each method ordered as [get, post, put, delete]
	routes := make([]map[string]router.RouterAction, 4)
	for i := 0; i < 4; i++ {
		routes[i] = make(map[string]router.RouterAction)
	}

	routes[0] = make(map[string]router.RouterAction)
	s := HttpServer{
		router.New(),
		ln,
		port,
		make([]ReqMiddleware, 0),
		logger,
	}
	return &s
}

// start the webserver
func (s *HttpServer) Start() {
	defer s.ln.Close()

	// TODO: add logging for # of middlwares, port, ip address
	for {
		c, err := s.ln.Accept() // accepts a TCP connection on the listener
		if err != nil {
			fmt.Println(err)
		}

		// concurrently handle connections
		go s.handleConn(c) // connection is on another fd (accepting conn opens a new fd)
	}

}

func (s *HttpServer) handleConn(c net.Conn) {
	// TODO: how to manage connection?
	defer c.Close()                                       // defer to close TCP connection
	message, err := bufio.NewReader(c).ReadString(('\n')) // wraps the tcp conn with reader and read msg
	if err != io.EOF && err != nil {
		s.logger.Warn(err.Error())
		return
	}

	req, err := msg.NewRequest(message) // instantiate [HttpRequest] from msg
	if err != nil {
		s.logger.Warn("Unable to instantiate [HttpRequest]: %s", err.Error())
		return
	}

	// process middlware
	processMiddleware(req)

	// routing
	res := s.To(req)

	resString, err := res.String() // convert to string to send to socket
	if err != nil {
		s.logger.Warn(err.Error())
		return
	}

	writer := bufio.NewWriter(c)                            // create buffered writer from net.Conn
	if _, err = writer.WriteString(resString); err != nil { // write the [HttpResponse] to buffer
		s.logger.Warn(err.Error())
		return
	} else {
		s.logger.Debug("Respond to client")
		writer.Flush() // respond to client (write to socket)
	}
}

func processMiddleware(req *msg.HttpRequest) {

}
