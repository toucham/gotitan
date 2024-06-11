package main

import (
	"context"
	"log"
	"net"
)

// TODO:
//  1. create a layer-4 proxy in main()
//     a. create a server
//     b. connect to a client
//     c. pass a message from server to client
//  2. refactor into packages
func main() {
	ctx := context.Background()

	go setupServer(ctx)
	go setupClient(ctx)
}

func setupServer(parentCtx context.Context) {
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

func setupClient(parentCtx context.Context) {

}
