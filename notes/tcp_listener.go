package main

import (
	"fmt"
	"net"
	"os"
	"time"
)

func handleConnection(conn net.Conn) {
	defer conn.Close()                                 // This will ensure that the connection is closed when the function returns.
	conn.SetDeadline(time.Now().Add(60 * time.Second)) // Set a deadline for the connection to avoid hanging indefinitely.
	fmt.Println("Handling new connection")
	outgoingData := []byte("Enter your name: ") // Prepare the data to send to the client.
	_, err := conn.Write(outgoingData)          // Write the data to the connection
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error writing to connection: %v\n", err)
		return
	}
	nameBytes := make([]byte, 1024) // Prepare a buffer to read incoming data.
	_, err = conn.Read(nameBytes)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading from connection: %v\n", err)
		return
	}
	name := string(nameBytes)                       // Convert the byte slice to a string.
	fmt.Println("Received data from client:", name) // Print the received data.
	conn.Write([]byte("Hello " + name))             // Send a response back to the client.
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	fmt.Println("We are listening on port 8080 for new connections")
	if err != nil {
		panic(err)
	}
	defer listener.Close()
	for {
		conn, err := listener.Accept()
		fmt.Println("New connection accepted")
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error accepting connection: %v\n", err)
			continue
		}
		go handleConnection(conn) // Let's handle the connection in a goroutine so that we can accept new connections concurrently.
	}
}
