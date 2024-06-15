package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net"
	"os"
)

const FILEPATH = ""

// TODO:
//  1. create a layer-4 proxy in main()
//     a. create a server
//     b. connect to a client
//     c. pass a message from server to client
//  2. refactor into packages
func main() {
	ctx := context.Background()

	go setupServer()
	go setupClient()

	<-ctx.Done()
}

func setupServer() {
	addr, err := net.ResolveTCPAddr("tcp", "127.0.0.1:8080")
	if err != nil {
		panic(err)
	}

	ln, err := net.ListenTCP("tcp", addr)
	if err != nil {
		panic(err)
	}
	log.Println("Listening on 127.0.0.1:8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Println(err)
		}

		go func() {
			defer conn.Close() // close connection
			log.Println("connection received!")
		}()
	}
}

func setupClient() {
	// open the connection
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}

	fmt.Fprintf(conn, "GET / HTTP/1.1\r\n\r\n")
	status, err := bufio.NewReader(conn).ReadString('\n') // receives the packet back
	if err != nil {
		log.Println(err)
		os.Exit(1)
	}
	log.Println(status)
}
