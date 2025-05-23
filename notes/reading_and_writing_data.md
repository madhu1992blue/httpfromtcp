# Reading files in Go

Here's a **refined version** of your summary with a **"Basic vs Advanced APIs"** section for better clarity:  

---

## **Go File Reading APIs: When to Use Which**  

### **🔹 Basic APIs (Great Starting Points)**
| **Function** | **Use Case** |
|-------------|-------------|
| **`bufio.Scanner`** | Read a file **line by line** or tokenize input. Default **64 KB per token limit**, but can be adjusted with `bufio.Scanner.Buffer`. |
| **`os.ReadFile`** | Read the **entire file into memory** at once. Good for **small to medium files**, but large files may cause memory issues. |
| **`bufio.NewReader`** | Read a file **with buffering** to improve performance and **reduce system calls**. Reads upto the provided capacity of data bytes array to the Read but doesn't guarantee reading the eact number of bytes requested, even when more data is available. Use io.ReadFull to ensure full reads. |
| **`os.File.Read`** | os.File.Read reads up to the provided buffer capacity but does not guarantee reading the exact number of bytes requested, even when more data is available. Use io.ReadFull to ensure full reads. |
| **`bufio.NewReaderSize`** | Similar to `bufio.NewReader`, but allows **customizing buffer size** for performance tuning. Reads upto the provided data bytes array to the Read but doesn't guarantee reading the eact number of bytes requested, even when more data is available. Use io.ReadFull to ensure full reads. |
| ** `io.ReadFull`** | Reads from an `io.Reader` until the buffer is full or EOF. Useful for **ensuring complete reads of a specific size**. We need to handle io.ErrUnexpectedEOF when data size requested is more than data available. Note: This can be problematic while reading from non-regular files like connection sockets as data might not be ready yet and io.ErrUnexpectedEOF will be hit. So, best not to use with data streams where data will arrive asynchronously.|
---

### **🔹 Advanced APIs (For Special Use Cases)**
| **Function** | **Use Case** |
|-------------|-------------|
| **`io.ReadAll`** | Reads from any `io.Reader` and returns the entire content as `[]byte`. Useful when working with **network streams or other readers**. |
| **`io.Copy` & `io.CopyBuffer`** | Efficiently copy data between an `io.Reader` and `io.Writer`, useful for **streaming large files**. |
| **`io.LimitReader`** | Wraps an `io.Reader` to **limit how much data can be read**, preventing excessive memory usage. |
| **`io.SectionReader`** | Allows reading a **specific portion of a file** without modifying the original file descriptor. Great for **large files**. |
| **`os.File.ReadAt` & `os.File.WriteAt`** | Perform **random-access reads/writes** at specific file offsets, useful for **binary or structured data**. |
| **`io.NewSectionReader`** | Creates a `SectionReader` that operates on a **subset of an `io.ReaderAt`**, allowing controlled reads within a file. |

---

### **💡 Summary**
- Start with the **Basic APIs** for **most file reading tasks**.  
- Use the **Advanced APIs** when you need **better control over memory, performance, or streaming**.  

Refer to reading_files/reading_files_and_network.go for practical examples for some basic APIs



# Similarity of Reading Files and Network Connections

In Go, we can read files and network connections using similar APIs. This is because both files and network connections implement the `io.Reader` interface, allowing us to use the same functions and methods for reading data from both sources.
Similarly, we can write to files and network connections using the `io.Writer` interface. This means that we can use the same functions and methods for writing data to both files and network connections.
This makes it easy to switch between reading from files and network connections without having to learn new APIs or functions. For example, we can use the `io.Copy` function to copy data from a file to a network connection or vice versa.

#  Difference between reading files and network connections

While similarility exists, there are some differences between reading files and network connections.
While reading a file, the following are in our control: 
- When to read?
- How much to read?
- When to stop reading?
This is called as pulling data from a file.
This is because, a file on disk has a specific size, data is always available, and you can read it in any order. You can also seek to a specific position in the file and read from there.

However, when reading from a network connection, we don't know when data will arrive, how much data will arrive or when it will stop arriving. This is called as pushing data from a network connection. So, APIs like io.ReadFull will not work as expected when reading from a network connection. They may cause io.ErrUnexpectedEOF error when data is not available yet. This may lead to more system calls and performance issues. The reader needs to be prepared to handle this.




