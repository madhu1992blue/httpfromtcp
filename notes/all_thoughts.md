# Create Go Module

```bash
go mod init github.com/madhu1992blue/httpfromtcp
```

# Open a File in Go

```go
package main
import (
    "fmt"
    "os"
)

func main() {
    file, err := os.Open("test.txt")
    if err != nil {
        fmt.Println(err)
        return
    }
    defer file.Close()
    fmt.Println("File opened successfully")
}
```



