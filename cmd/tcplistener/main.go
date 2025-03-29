package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"strings"
	"time"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)

	go func() {
		defer close(linesChan)
		dataBytes := make([]byte, 8)
		remaining := ""
		for {
			bytesRead, err := f.Read(dataBytes)
			if err != nil {
				if err == io.EOF {
					break
				}
				fmt.Fprintf(os.Stderr, "Error reading from connection: %v\n", err)
				return
			}
			textToProcess := string(dataBytes[:bytesRead])
			remaining += textToProcess
			parts := strings.Split(remaining, "\n")
			for i := 0; i < len(parts)-1; i++ {
				line := strings.TrimSpace(parts[i])
				if line != "" {
					linesChan <- line
				}
			}
			remaining = parts[len(parts)-1]
		}
		linesChan <- remaining
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
