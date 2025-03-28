package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"time"
)

func scan_file_with_bufio(filename string) {
	/* In this example, we will read a file line by line using bufio.Scanner.
	 * The bufio package provides buffered I/O, which can be more efficient than
	 * reading directly from the file. The Scanner type provides a convenient
	 * way to read data from a file, allowing us to read it line by line.
	* The scanner is designed to split by a SplitFunction and that defaults to the one that splits by line.
	*/
	fmt.Println("======== Demoing bufio.Scanner ===========")
	startTime := time.Now() // Start the timer
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()                // Any open file should be closed when done . So, defer it.
	scanner := bufio.NewScanner(file) // Create a new scanner for the file
	for scanner.Scan() {
		text := scanner.Text()
		fmt.Println(text) // Print each line
	}
	if scanner.Err() != nil { // We don't need to check io.EOF as if EOF is reached, err will still be nil for this API.
		fmt.Println("Error reading file:", scanner.Err())
	}
	endTime := time.Now() // End the timer
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Time taken to read the file: %v\n", elapsedTime)
}

func read_full_file(filename string) {
	/* In this example, we will read the entire file at once using os.ReadFile.
	 * The os package provides a simple way to read the entire contents of a file
	 * into memory.
	 */
	startTime := time.Now() // Start the timer
	fmt.Println("======== Demoing reading the entire file at once ===========")
	dataBytes, err := os.ReadFile(filename) // Read the entire file at once
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	fmt.Println(string(dataBytes)) // Print the entire file content
	endTime := time.Now()          // End the timer
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Time taken to read the file: %v\n", elapsedTime)
}

func read_file_with_buffered_reader(filename string) {
	fmt.Println("======== Demoing reading the entire file with buffered IO ===========")
	/* In this example, we will read the entire file using a buffered reader which has an internal buffer too.
		This may read more data into internal buffer than the size of the dataBytes slice but will provide callers with upto the size of dataBytes.
	 * The bufio package provides buffered I/O, which can be more efficient than
	 * reading directly from the file as it reduces system calls by reading larger chunks of data at once into internal buffer.
	 Note: The Read function doesn't guarantee that it will read the exact number of bytes requested even if more data is available.
	 It may read fewer bytes than the buffer size, so we need to handle that case.
	 * To ensure reading an exact number of bytes, use io.ReadFull.
	*/
	startTime := time.Now() // Start the timer
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()              // Any open file should be closed when done . So, defer it.
	reader := bufio.NewReader(file) // This creates an internal buffer of size 4096 bytes. The internal buffer size can be controlled with bufio.NewReaderSize(file, size).
	dataBytes := make([]byte, 8)    // Create a buffer of size 8 bytes. The internal buffer is different from the dataBytes of size 8 we are creating here.
	for {

		bytesRead, err := reader.Read(dataBytes) // dataBytes is the bytes slice of 8 bytes and we read upto the size of dataBytes.
		// However, the buffered reader will read more than 8 bytes from the file and store it in its internal buffer.
		// Since the internal buffer default is 4096 bytes, it will read upto 4096 bytes from the file and store it in its internal buffer.
		// This was designed to minimize the number of system calls to read from the file.
		if err != nil {
			if err.Error() == "EOF" {
				break // End of file reached
			}
			fmt.Println("Error reading file:", err)
			return
		}
		msg := string(dataBytes[:bytesRead])
		fmt.Print(msg) // Print the read data
	}
	fmt.Println()         // Print a new line after reading the file
	endTime := time.Now() // End the timer
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Time taken to read the file: %v\n", elapsedTime)
}

func read_file_without_bufio(filename string) {
	fmt.Println("======== Demoing reading the file without bufio ===========")
	/* In this example, we will read the file without bufio. This makes too many system calls to read from the file.
	 */

	startTime := time.Now() // Start the timer

	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()           // Any open file should be closed when done . So, defer it.
	dataBytes := make([]byte, 8) // Create a buffer of size 8 bytes
	for {
		bytesRead, err := file.Read(dataBytes) // Doesn't guarantee reading the exact number of bytes requested.
		// It may read fewer bytes than the buffer size, even when more data is available.
		// To ensure reading an exact number of bytes, use io.ReadFull.
		// This will read upto the size of dataBytes.
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			fmt.Println("Error reading file:", err)
			return
		}
		msg := string(dataBytes[:bytesRead])
		fmt.Print(msg) // Print the read data
	}
	fmt.Println()         // Print a new line after reading the file
	endTime := time.Now() // End the timer
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Time taken to read the file: %v\n", elapsedTime)
}

func read_file_with_exact_sizes_and_bufio(filename string) {
	fmt.Println("======== Demoing reading the file with exact sizes with bufio ===========")
	/* In this example, we will read the file with bufio and io.ReadFull to read an exact number of bytes.
	 * The bufio package provides buffered I/O, which can be more efficient than
	 * reading directly from the file. The io.ReadFull function ensures that we read
	 * the exact number of bytes requested. bufio helps to reduce system calls by reading larger chunks of data at once into internal buffer.
	 * The io.ReadFull function ensures that we read the exact number of bytes requested.
	 */

	startTime := time.Now() // Start the timer
	file, err := os.Open(filename)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer file.Close()              // Any open file should be closed when done . So, defer it.
	reader := bufio.NewReader(file) // This creates an internal buffer of size 4096 bytes. The internal buffer size can be controlled with bufio.NewReaderSize(file, size).
	dataBytes := make([]byte, 8)    // Create a buffer of size 8 bytes
	for {
		bytesRead, err := io.ReadFull(reader, dataBytes) // Read exactly 8 bytes from the file using io.ReadFull
		if err != nil {
			if err == io.EOF {
				break // End of file reached
			}
			if err == io.ErrUnexpectedEOF {
				// This means we reached EOF before reading the exact number of bytes requested.
				msg := string(dataBytes[:bytesRead])
				// Print the read data
				fmt.Print(msg)
				break
			}
			fmt.Println("Error reading file:", err)
			return
		}
		msg := string(dataBytes[:bytesRead])
		fmt.Print(msg) // Print the read data
	}
	fmt.Println()         // Print a new line after reading the file
	endTime := time.Now() // End the timer
	elapsedTime := endTime.Sub(startTime)
	fmt.Printf("Time taken to read the file: %v\n", elapsedTime)
}

func main() {

	scan_file_with_bufio("message.txt")                 // Scan a file line by line
	read_full_file("message.txt")                       // Read the entire file at once
	read_file_with_buffered_reader("message.txt")       // Read the entire file with buffered IO
	read_file_without_bufio("message.txt")              // Read a file without bufio but using os.File.Read
	read_file_with_exact_sizes_and_bufio("message.txt") // Read a file with bufio to reduce system calls and read exact size.
}

/* Summary
When to use which:
- Use bufio.Scanner when you want to read a file line by line or tokenizing input. Its default limit is 64kB per token (like a line) but this can be adjusted with bufio.Scanner.Buffer.
- Use os.ReadFile when you want to read the entire file at once and load it into memory. This is useful for small or medium sized files but large files may cause memory issues.
- Use bufio.NewReader when you want to read a file with buffering for better performance to reduce system calls.
- os.File.Read may be problematic because it only guarantees reading up to the provided buffer size, not the full amount requested.
	- os.File.Read does not guarantee reading the exact number of bytes requested.
	- It may read fewer bytes than the buffer size, even when more data is available.
	- To ensure reading an exact number of bytes, use io.ReadFull.
- Use bufio.NewReaderSize when you want to control the size of the internal buffer for better performance.
- When specific size is needed, Use `io.ReadFull(reader, byteSlice)` when you want to read an exact number of bytes from a reader or a file.
	- Use io.ReadFull when you need an exact number of bytes.
	- If the file is shorter than expected, handle io.ErrUnexpectedEOF.
	- This is almost always used when reading exact sizes of data from a reader.
*/
