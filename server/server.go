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
	logger := logger.New()
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
		go s.handleConn(c)
	}

}

func (s *HttpServer) handleConn(c net.Conn) {
	defer c.Close()
	// wraps the tcp connection with reader and read app message
	message, err := bufio.NewReader(c).ReadString(('\n'))
	if err != io.EOF && err != nil {
		s.logger.Warn(err.Error())
		return
	}

	fmt.Printf("New message: %s", message)
	// parse app-layer message to create HttpRequest
	req, err := msg.NewRequest(message)
	if err != nil {
		panic("oh no") // TODO: replace panic to logging
	}

	// process middlware
	processMiddleware(req)

	// routing
	res := s.To(req)

	// convert to string to send to socket
	resString, err := res.String()
	if err != nil {
		s.logger.Warn(err.Error())
		return
	}

	// create buffered writer
	writer := bufio.NewWriter(c)
	if _, err = writer.WriteString(resString); err != nil { // write the [HttpResponse] to buffer
		s.logger.Warn(err.Error())
		return
	} else {
		writer.Flush() // respond to client (write to socket)
		s.logger.Debug("Responded to client")
	}
}

func processMiddleware(req *msg.HttpRequest) {

}
