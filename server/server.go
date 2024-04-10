package server

import (
	"bufio"
	"fmt"
	"net"
)

type HttpServer struct {
	ln    net.Listener // socket listener
	port  string
	reqMw []ReqMiddleware
	resMw []ResMiddleware
}

func Init(port string) *HttpServer {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		panic(err)
	}

	s := HttpServer{
		ln,
		port,
		make([]ReqMiddleware, 2), // expect at least 2
		make([]ResMiddleware, 1),
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
		go handleConn(c)
	}

}

func handleConn(c net.Conn) {
	// wraps the tcp connection with reader
	netData, err := bufio.NewReader(c).ReadString(('\n'))
	if err != nil {
		fmt.Println(err)
		return
	}

	// get HttpRequest
	req, err := extractReq(netData)
	if err != nil {
		panic("oh no") // TODO: replace panic to logging
	}

	// process middlware
	processMiddlware(req)

	// routing
}

// parse raw data to instantiate HttpRequest according to HTTP/1.1
func extractReq(msg string) (*HttpRequest, error) {
	req := HttpRequest{}

	// 1) get request line

	// 2) get fields (headers)

	// 3) get body

	return &req, nil
}

func processMiddlware(req *HttpRequest) {

}
