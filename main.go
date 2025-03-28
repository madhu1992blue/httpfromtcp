package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)
	scanner := bufio.NewScanner(f)
	go func() {
		for scanner.Scan() {
			linesChan <- scanner.Text() // Send the line to the channel
		}
		close(linesChan)
	}()
	return linesChan
}

func handleConnection(conn net.Conn) {
	defer func() {
		conn.Close()
		fmt.Println("Connection closed")
	}()
	conn.SetDeadline(time.Now().Add(600 * time.Second)) // Set a deadline for the connection to avoid hanging indefinitely.
	linesChan := getLinesChannel(conn)
	for line := range linesChan {
		fmt.Println(line)
	}
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
