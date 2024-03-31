package server

import (
	"bufio"
	"fmt"
	"net"
	"strings"
)

type HttpServer struct {
	ln   net.Listener
	port string
}

func Init(port string) *HttpServer {
	ln, err := net.Listen("tcp", ":"+port)

	if err != nil {
		panic(err)
	}

	s := HttpServer{
		ln,
		port,
	}
	return &s
}

func (s *HttpServer) Start() {
	defer s.ln.Close()

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

	fmt.Println(string(netData)) // print the app-layer message
	c.Write([]byte("OK\n"))      // send back an OK message

	if strings.TrimSpace((string(netData))) == "STOP" {
		fmt.Println("Exiting TCP server!")
		return
	}
}
