package main

import (
	"banking/mongodb"
	"banking/server"
	"os"
)

func main() {
	mongodb.Init()
	port := os.Args[1]
	server.RunServer(port)
}
