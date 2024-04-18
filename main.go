package main

import (
	"fmt"

	"github.com/toucham/gotitan/server"
	"github.com/toucham/gotitan/server/msg"
)

func main() {
	PORT := "8080"
	HOST := "127.0.0.1"

	app := server.Init(HOST, PORT)
	// add middlware
	// app.AddReqMiddlware(func(req *server.HttpRequest) {

	// }, server.MiddlwareOptions{})

	// add routing
	indexAction := func(req *msg.HttpRequest) *msg.HttpResponse {
		fmt.Printf("Received request for method: %s", req.GetMethod())
		fmt.Printf("With body: %s", req.GetBody())
		return new(msg.HttpResponse).SetStatus(400)
	}
	app.AddRoute(msg.HTTP_POST, "/", indexAction)
	app.AddRoute(msg.HTTP_GET, "/", indexAction)

	app.Start()
}
