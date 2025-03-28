package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	linesChan := make(chan string)
	bufferedReader := bufio.NewReader(f)
	go func() {

		dataBytes := make([]byte, 8) // Create a byte slice of size 8 bytes
		var currentLine string
		for {
			bytesRead, err := io.ReadFull(bufferedReader, dataBytes)
			if err != nil && err != io.ErrUnexpectedEOF && err != io.EOF {
				// Handle the error if it's not EOF or unexpected EOF
				fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
				os.Exit(1)
			}
			if err == io.EOF {
				break // End of file reached
			}
			msg := currentLine + string(dataBytes[:bytesRead])
			parts := strings.Split(msg, "\n") // Split the message into lines
			for i := 0; i < len(parts)-1; i++ {
				linesChan <- parts[i] // Print the read data
			}
			currentLine = parts[len(parts)-1] // Keep the last part for the next iteration
			if err == io.ErrUnexpectedEOF {
				break // Break the loop if we reach EOF before reading the exact number of bytes requested
			}
		}
		if currentLine != "" {
			linesChan <- currentLine // Print the remaining line
		}
		close(linesChan)
	}()
	return linesChan
}

func main() {
	msgFile, err := os.Open("message.txt")
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error opening file: %v\n", err)
		os.Exit(1)
	}
	defer msgFile.Close()

	for line := range getLinesChannel(msgFile) {
		fmt.Printf("read: %s\n", line)
	}
}
