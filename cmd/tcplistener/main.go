package main

import (
	"fmt"
	"net"
	"os"
	"time"

	"github.com/madhu1992blue/httpfromtcp/internal/request"
)

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		fmt.Println("Connection closed")
	}()
	conn.SetDeadline(time.Now().Add(600 * time.Second)) // Set a deadline for the connection to avoid hanging indefinitely.
	request, err := request.RequestFromReader(conn)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request: %v\n", err)
		return
	}
	fmt.Printf(`Request line:
- Method: %s
- Target: %s
- Version: %s`, request.RequestLine.Method, request.RequestLine.RequestTarget, request.RequestLine.HttpVersion)
	fmt.Println()
}
func main() {

	listener, err := net.Listen("tcp", ":42069")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error starting TCP listener: %v\n", err)
		os.Exit(1)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Fprintf(os.Stderr, "listener closed: %v", err)
			break
		}
		fmt.Println("New connection accepted")
		go handleConnection(conn)

	}
}
