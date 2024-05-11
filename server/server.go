package server

import (
	"fmt"
	"net"

	"github.com/toucham/gotitan/logger"
	"github.com/toucham/gotitan/server/conn"
	"github.com/toucham/gotitan/server/router"
)

const (
	TIMEOUT = 10000 // ms (default at 10s)
)

type HttpServer struct {
	router.Router              // embedded Router
	ln            net.Listener // socket listener
	port          string
	reqMw         []ReqMiddleware
	logger        logger.Logger
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

// Start start the http server to listen on the designated port
func (s *HttpServer) Start() {
	defer s.ln.Close() // listen on the designated {s.port}

	for {
		// block process until it accepts a TCP connection
		if c, err := s.ln.Accept(); err != nil {
			s.logger.Fatal(err.Error())
			return
		} else {
			// create a connection handler
			go conn.HandleConn(c, &s.Router, logger.New("HandleConn"), TIMEOUT)
		}
	}
}
