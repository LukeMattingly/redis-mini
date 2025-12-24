package main

import (
	"fmt"
	"net"
)

func main() {
	connection, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Connected to a server")
}
