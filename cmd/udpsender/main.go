package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
)

func main() {
	udpConn, err := net.Dial("udp", "127.0.0.1:42049") // While the return type calls it net.Conn, UDP is a connectionless protocol.
	if err != nil {
		fmt.Println("Error dialing UDP address:", err)
		return
	}
	defer udpConn.Close() // Close the UDP client when done
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print(">")
		line, err := reader.ReadString('\n') // Read a line from standard input
		if err != nil {
			fmt.Println("Error reading input:", err)
			break
		}
		udpConn.Write([]byte(line)) // Send the message to the UDP server
	}
	if err != nil {
		fmt.Println("Error sending message:", err)
		return
	}
	fmt.Println("Message sent to UDP server")

}
