package main

import "os"

var server Server

func main() {
	server = NewServer(os.Args[1])
	server.Handle()
}