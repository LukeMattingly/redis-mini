package main

import (
	"fmt"
	"net"
	"os"
	"strings"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage:")
		fmt.Println("  client PING")
		fmt.Println("  client SET key value")
		fmt.Println("  client GET key")
		return
	}

	args := os.Args[1:]

	connection, err := net.Dial("tcp", "127.0.0.1:6379")
	if err != nil {
		panic(err)
	}
	defer connection.Close()

	fmt.Println("Connected to a server")

	err = sendCommand(connection, args)
	if err != nil {
		fmt.Println("Send error", err)
		return
	}

	resp, err := readResponse(connection)
	if err != nil {
		fmt.Println("Read error:", err)
		return
	}

	fmt.Println("Response:", parseResponse(resp))
}

func encodeCommand(args []string) []byte {
	result := fmt.Sprintf("*%d\r\n", len(args))

	for _, arg := range args {
		result += fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg)
	}

	return []byte(result)
}

func sendCommand(connection net.Conn, args []string) error {
	cmd := encodeCommand(args)

	_, err := connection.Write(cmd)
	return err
}

func readResponse(connnection net.Conn) (string, error) {
	buffer := make([]byte, 4096)
	line, err := connnection.Read(buffer)
	if err != nil {
		return "", err
	}
	return string(buffer[:line]), nil
}

func parseResponse(resp string) string {
	switch resp[0] {
	case '+':
		return resp[1 : len(resp)-2]
	case '-':
		return "ERR: " + resp[1:len(resp)-2]
	case '$':
		lines := strings.Split(resp, "\r\n")
		if len(lines) >= 2 {
			return lines[1]
		}
		return ""
	default:
		return resp
	}
}
