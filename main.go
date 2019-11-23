package main

import (
	"banking/mongodb"
	"banking/server"
	"fmt"
	"go/build"
	"os"
)

func main() {
	mongodb.Init()
	gopath := os.Getenv("GOPATH")
	if gopath == "" {
		fmt.Println("bad")
		gopath = build.Default.GOPATH
	}
	fmt.Println(gopath)
	port := os.Args[1]
	server.RunServer(port)
}
