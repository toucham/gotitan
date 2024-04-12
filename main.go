package main

import (
	"fmt"
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
		fmt.Printf("Received request for method: %s", req.GetMethod())
		fmt.Printf("With body: %s", req.GetBody())
		return nil
	}
	app.AddRoute(server.HTTP_POST, "/", indexAction)
	app.AddRoute(server.HTTP_GET, "/", indexAction)

	app.Start()
}
