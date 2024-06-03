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
	indexAction := func(req msg.Request) msg.Response {
		fmt.Printf("Received request for method: %s", req.GetMethod())
		fmt.Printf("With body: %s", req.GetBody())
		return new(msg.HttpResponse)
	}
	app.AddRoute(msg.HTTP_POST, "/", func(req msg.Request) msg.Response {
		res := msg.NewHttpResponse()
		res.SetBody("<div> <h1> Welcome to the index page </h1> </div>", "text/html; charset=utf-8")
		return res
	})
	app.AddRoute(msg.HTTP_GET, "/index.html", indexAction)

	app.Start()
}
