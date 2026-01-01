package main

import (
	"fmt"
	"io"
	"net"
	"strings"
)

func main() {
	fmt.Println("Listening on port :6379")

	// Create a new server
	l, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		panic(err)
	}

	defer l.Close()

	// Listen for connections
	for {
		conn, err := l.Accept()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}

		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	resp := NewResp(conn)
	writer := NewWriter(conn)

	for {
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				// client closed connection — normal
				return
			}
			fmt.Println("read error:", err)
			return
		}

		if value.Type != TypeArray || len(value.Array) == 0 {
			writer.Write(Value{
				Type: TypeError,
				Str:  "ERR invalid request",
			})
			continue
		}

		command := strings.ToUpper(value.Array[0].Bulk)
		args := value.Array[1:]

		fmt.Print("Command:", command)
		for _, arg := range args {
			fmt.Print(" ", arg.Bulk)
		}
		fmt.Println()

		handler, ok := Handlers[command]
		if !ok {
			writer.Write(Value{
				Type: TypeError,
				Str:  "ERR unknown command",
			})
			continue
		}

		result := handler(args)
		writer.Write(result)
	}
}
