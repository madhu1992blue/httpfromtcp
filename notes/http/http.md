# HTTP

TCP is just a transport protocol but it does not define how data transported by it should be interpreted. 
HTTP is one of the application protocols that uses TCP as its transport protocol. HTTP defines how data should be formatted, transmitted, and interpreted by the client and server.

HTTP Versions that are in use:
- HTTP/1.0
- HTTP/1.1
- HTTP/2
- HTTP/3

HTTP (Hypertext Transfer Protocol) has evolved significantly across its versions—HTTP/1.1, HTTP/2, and HTTP/3—each improving performance, efficiency, and security. Below are the key differences:

### 1️⃣ **HTTP/1.1** (1997)
- **Connection Handling:** Uses a single TCP connection per request-response cycle (supports keep-alive but limited concurrency).
- **Head-of-Line Blocking:** Sequential request processing; slow responses delay others in the queue.
- **Compression:** No built-in header compression.
- **Multiplexing:** Not supported; each request needs its own connection.
- **Security:** Works with TLS but not mandatory.
- **Latency:** Higher due to sequential processing and extra round trips.

### 2️⃣ **HTTP/2** (2015)
- **Connection Handling:** Uses a single TCP connection for multiple streams (multiplexing).
- **Head-of-Line Blocking:** Resolved at the HTTP level but still exists at the TCP level.
    - HTTP/2 can multiplex multiple requests over a single TCP connection, but if one TCP packet is lost, all streams sharing that packet must wait for retransmission.
    - QUIC (HTTP/3) solves this by making streams independent, so packet loss on one stream doesn’t delay others.
- **Compression:** Uses **HPACK** for header compression, reducing redundancy.
- **Multiplexing:** Allows multiple streams over a single connection, improving efficiency.
- **Security:** Requires TLS (though not mandated in the spec, all browsers enforce HTTPS).
- **Latency:** Lower latency due to multiplexing and header compression.

### 3️⃣ **HTTP/3** (2022)
- **Connection Handling:** Uses **QUIC** (instead of TCP) to reduce connection setup time.
- **Head-of-Line Blocking:** Fully eliminated since QUIC operates over UDP.
- **Compression:** Uses **QPACK**, an improved version of HPACK, for better header compression.
- **Multiplexing:** Improved with QUIC, ensuring independent request processing.
- **Security:** TLS 1.3 is **built-in** to QUIC, eliminating separate TLS handshakes.
- **Latency:** Further reduced due to faster handshakes and independent request streams.

### 🔥 **Key Takeaways**
| Feature         | HTTP/1.1  | HTTP/2  | HTTP/3  |
|---------------|---------|--------|--------|
| Protocol      | TCP     | TCP    | QUIC (UDP) |
| Multiplexing  | ❌ No   | ✅ Yes  | ✅ Yes  |
| Head-of-Line Blocking | ✅ Yes | ⚠️ TCP-Level | ❌ No |
| Header Compression | ❌ No | ✅ HPACK | ✅ QPACK |
| Security      | Optional TLS | Mandatory TLS | Built-in TLS 1.3 |
| Latency      | High   | Medium | Low |

**In summary**:  
- **HTTP/2** improved efficiency with multiplexing and compression.  
- **HTTP/3** takes it further by using QUIC over UDP, reducing latency and eliminating head-of-line blocking.

Note: The message formatting for HTTP 1.1, 2, and 3 is not the same. HTTP 1.1 uses text-based messages while HTTP 2 and 3 use binary messages. This means that the way data is sent and received is different in each version.


# HTTP/1.1

HTTP/1.1 is a text-based protocol that uses a request-response model. It is the most widely used version of HTTP and is supported by all modern web browsers and servers.

## Request Format
Note: Each Line is terminated by a CRLF (Carriage Return Line Feed) sequence, which is represented as `\r\n` in the request.
The request consists of:
- **Request Line**: Contains the HTTP method (GET, POST, etc.), the URL, and the HTTP version.
    Ex: `GET /index.html HTTP/1.1`
- **Field Lines/Headers**: Contains headers that provide additional information about the request.
    Ex: `Host: www.example.com`
    Some useful headers:
    - `Host`: Specifies the domain name of the server.
    - `User-Agent`: Identifies the client software making the request.
    - `Accept`: Specifies the media types that the client is willing to accept.
    - `Content-Type`: Indicates the media type of the resource being sent in the body.
- Empty Line: Indicates the end of the headers. 
- **Body**: Optional, used for methods like POST to send data to the server.


## Headers

How do we communicate our content size to the server?

### Content Size for HTTP v1.1
One of the following headers is used to communicate the size of the body:
- **Transfer-Encoding: chunked** header: Specifies the encoding used to transfer the body. Commonly used with `chunked` encoding, where the body is sent in chunks of data.
  Each chunk is preceded by its size in bytes.
    - Ex: `Transfer-Encoding: chunked`
    - The body is sent in chunks, each prefixed with its size.
    - The last chunk is followed by a zero-length chunk to indicate the end of the message.
  - Example of chunked transfer encoding request:
    ```
    POST /upload HTTP/1.1
    Host: www.example.com
    Transfer-Encoding: chunked
    4
    Wiki
    5
    Rules
    0
    ```
  - In this example, the body is sent in two chunks: "Wiki" (4 bytes) and "Rules" (5 bytes). The last chunk is of size 0, indicating the end of the message.
- **Content-Length** header: Specifies the size of the body in bytes. This is used when the body is sent in a single chunk.
    - Ex: `Content-Length: 15`
    - Example of content length request:
        ```
        POST /upload HTTP/1.1
        Host: www.example.com
        Content-Length: 15
        Content-Type: application/json

        {"name":"John"}
        ```
    - In this example, the body is sent in a single chunk with a size of 15 bytes.
- **Closing the connection**: The server can also close the connection to indicate the end of the message. This is not commonly used in modern applications.


### Content Size for HTTP v2 and 3

- We don't need to specify the content size in HTTP v2 and 3. In HTTP/2 and HTTP/3, the body is sent in frames, each with a size field, so the server automatically knows the total body length. This is because HTTP/2 and HTTP/3 use binary framing, which improves efficiency by allowing parallel request processing.
- In HTTP v2 and 3, the body is sent in frames, which are small packets of data that can be sent independently. Each frame has a header that specifies its size and type, allowing the server to determine the size of the body without needing to specify it in the request.
- While transfer-encoding and content-length headers are still supported in HTTP v2 and 3, they are not required. The server can determine the size of the body based on the frames being sent.