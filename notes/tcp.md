# TCP/IP Model


# 4 Layers of Protocol

# 1. Application Layer Protocol: HTTP, FTP, SMTP, DNS, etc.
# 2. Transport Layer: TCP, UDP
# 3. Internet Layer: IP, ICMP
# 4. Link Layer: Ethernet, Wi-Fi, etc.

TCP is great for sending ordered data over the internet reliably.

In the TCP protocol, the data is broken down into packets.
These packets may arrive out of order on the destination and they are reassembled on the other side.
There are packets that get lost. Detecting that and resending is also part of the protocol.


## **Conceptual Model for Listen and Connection**

- **The server creates a listening socket** to accept new connections:  
    ```go
    listener, err := net.Listen(network, address) // listener is a net.Listener
    ```
    - This creates a **passive socket** that waits for incoming connections.
    - The OS maintains a **queue (backlog)** for new connections until they are accepted.
    - When a client initiates a connection, it is **added to this queue**.

- **The server accepts a connection from the queue**:  
    ```go
    conn, err := listener.Accept() // conn is a net.Conn
    ```
    - This removes a pending connection from the queue and creates a **new active socket** for communication.
    - The server can now **send and receive data** using:
      ```go
      conn.Read(buffer) // Receive data
      conn.Write(data)  // Send data
      ```
    - "net.Conn allows bidirectional communication, meaning both client and server can send and receive data simultaneously."
    - **Best Practices:**
      - Set a timeout to prevent **hanging connections**:  
        ```go
        conn.SetDeadline(time.Now().Add(30 * time.Second))
        ```
      - Always close the connection when done:  
        ```go
        conn.Close()
        ```
    - Even when the listener is closed, the connection can still be used until it is closed.  
      This is because the listener is a passive socket and the connection is an active socket.
    - **The listener can accept multiple connections**:  
      Each accepted connection creates a new `net.Conn` object, allowing the server to handle multiple clients concurrently.


- **The server can also close the listener**:
    ```go
    listener.Close()
    ```
    - This stops accepting new connections but does not affect existing connections.
    - Existing connections can still be used until they are closed.
    - If listener.Accept() is waiting for a connection, closing the listener will cause it to return an error.
        Ex: "Closing the listener will cause any blocked Accept() calls to return an error."


