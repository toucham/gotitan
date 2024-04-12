package server

import (
	"bufio"
	"fmt"
	"net"
)

type HttpServer struct {
	ln     net.Listener // socket listener
	port   string
	reqMw  []ReqMiddleware
	resMw  []ResMiddleware
	routes []map[string]HttpAction
}

func Init(port string) *HttpServer {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		panic(err)
	}

	// HashMap for each method ordered as [get, post, put, delete]
	routes := make([]map[string]HttpAction, 4)
	for i := 0; i < 4; i++ {
		routes[i] = make(map[string]HttpAction)
	}

	routes[0] = make(map[string]HttpAction)
	s := HttpServer{
		ln,
		port,
		make([]ReqMiddleware, 2), // expect at least 2
		make([]ResMiddleware, 1),
		routes,
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

	// parse http message to create HttpRequest
	req, err := ExtractRequest(netData)
	if err != nil {
		panic("oh no") // TODO: replace panic to logging
	}

	// process middlware
	processMiddlware(req)

	// routing
	reqToRoute(req)
}

func processMiddlware(req *HttpRequest) {

}

func reqToRoute(req *HttpRequest) {

}
