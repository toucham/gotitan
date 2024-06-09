package main

import (
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
		res := msg.CreateHttpResponse(msg.StatusOk)
		res.SetBody("<div> <h1> Welcome to the index page </h1> </div>", "text/html; charset=utf-8")
		return res
	}
	app.AddRoute(msg.HTTP_GET, "/", func(req msg.Request) msg.Response {
		res := msg.CreateHttpResponse(msg.StatusOk)
		res.SetBody("<div> <h1> Welcome to the first page </h1> </div>", "text/html; charset=utf-8")
		return res
	})
	app.AddRoute(msg.HTTP_GET, "/index", indexAction)

	app.Start()
}
