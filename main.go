package main

import (
	"toucham/gotitan/server"
)

func main() {
	PORT := "8080"
	app := server.Init(PORT)
	// add middlware

	// add routing
	app.Start()
}
