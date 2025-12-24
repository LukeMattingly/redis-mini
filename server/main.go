package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func main() {

	//ready to recieve tcp requests
	listener, err := net.Listen("tcp", ":6379")
	fmt.Println("Listening on port :6379")

	if err != nil {
		fmt.Println(err)
		return
	}

	//accept new connection requests
	connection, err := listener.Accept()
	if err != nil {
		fmt.Println(err)
		return
	}

	defer connection.Close()

	//ready to recieve commands from client
	for {
		buffer := make([]byte, 1024)

		//read message from client
		_, err = connection.Read(buffer)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("error reading from client: ", err.Error())
			os.Exit(1)
		}

		//ignore request and send back a PONG
		connection.Write([]byte("+OK\r\n"))
	}
}
