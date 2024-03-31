package main

import (
	"toucham/gotitan/server"
)

func main() {
	PORT := "8080"
	app := server.Init(PORT)
	app.Start()
}
