package main

import (
	"toucham/gotitan/server"
)

func main() {
	PORT := "8080"
	app := server.Init(PORT)
	// add middlware
	app.AddReqMiddlware(func(req *server.HttpRequest) {

	}, server.MiddlwareOptions{})

	// add routing
	indexAction := func(req *server.HttpRequest) *server.HttpResponse {
		return nil
	}
	app.AddRoute(server.HTTP_POST, "/", indexAction)

	app.Start()
}
