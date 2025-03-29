# TCP/IP Model


# 4 Layers of Protocol

# 1. Application Layer Protocol: HTTP, FTP, SMTP, DNS, etc.
# 2. Transport Layer: TCP, UDP
# 3. Internet Layer: IP, ICMP
# 4. Link Layer: Ethernet, Wi-Fi, etc.

TCP is great for sending ordered data over the internet reliably.

In the TCP protocol, the data is broken down into TCP segments.
These TCP segments may arrive out of order on the destination and they are reassembled on the other side.
There are TCP segments that get lost. Detecting that and resending is also part of the protocol.


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


## TCP's Data Transmission Process

- The Connection is established using a **3-way handshake**: 
    1. **SYN**: The client sends a SYN TCP segment to the server to initiate a connection.
    2. **SYN-ACK**: The server responds with a SYN-ACK TCP segment to acknowledge the request.
    3. **ACK**: The client sends an ACK TCP segment back to the server, completing the handshake.
    
    (In go , net.Conn object is created after the handshake)

- **Data Transmission**: 
    - After the handshake, data can be sent in both directions.
    - TCP ensures that data is sent in order and without errors.
    - If a TCP segment is lost, TCP will retransmit it.
    - There are at least 3 windows for TCP on each side (client and server):
      - Send Window(SWND): The sender (client or server) determines how much data can be sent before requiring an acknowledgment.
      - Receive Window(RWND): The receiver(client or server) advertises the amount of buffer space available for incoming data.
      - Congestion Window(CWND): The amount of data that can be sent over the network without overwhelming the network. It is controlled by the sender. 
      
      SWND=min(CWND,RWND)  
      Note: Each side advertises its own receive window size, which is the amount of data it can accept without overflowing its buffer. The send window of client should be less than or equal to the receive window of server. The send window of server should be less than or equal to the receive window of client.

- **Data Retransmission**: 
    - TCP uses a sliding window protocol to manage data flow.
    - The sender can send multiple TCP segments before needing an acknowledgment, improving efficiency (SWND).
    - The receiver sends ACK TCP segments to confirm receipt of data.
    - If a TCP segment is lost or corrupted, TCP will retransmit it.

    - How does TCP know if a TCP segment is lost?
        - **Timeout Based Retransmission**: It relies on a Retransmission Timeout(RTO) on the sender where if a sender doesn't get an ACK within Retransmission timeout(RTO), the TCP segment is resent. TCP uses a **timeout** mechanism. If the sender does not receive an ACK for a TCP segment within a certain time, it assumes the TCP segment is lost and retransmits it.

        - **Fast Retransmission**: When the receiver detects new TCP segments, it sees if there is a gap in sequence numbers and the TCP segment is out of order. If yes, it sends a duplicate Ack for the last received TCP segment that was in order. Upon receiving three duplicate ACKs(a total of 4 acks counting the original) for the same segment, the sender assumes packet loss and retransmits it immediately.
      
        - TCP uses a **cumulative acknowledgment** scheme, where the receiver acknowledges all TCP segments up to a certain point. This allows the sender to know which TCP segments have been received successfully. This means if TCP segment 4 is Acked, then TCP segments 1, 2, and 3 are also considered Acked. This is called **cumulative acknowledgment**. That is why when a new TCP segment that is not in order is received, the receiver sends a duplicate Ack for the last received TCP segment that was in order. This is called **duplicate acknowledgment**.

        - There is now a selective acknowledgment (SACK) option in TCP that allows the receiver to inform the sender about all segments that have been received successfully, allowing the sender to retransmit only the missing segments. This is called **selective acknowledgment**. This is an optional feature and not all TCP implementations support it.

- **Connection Termination**:

  Here's a breakdown of the TCP connection termination (FIN sequence) represented in the table, with explanations for each stage:

  **Stage 1: Finish Initiator (FIN, ACK)**

  * **Actor:** The initiator, which is the side that decides to close the connection first, initiates the termination process.
  * **TCP Segment Sent:** The initiator sends a TCP segment with both the FIN (finish) and ACK (acknowledgment) flags set. This signals that it has finished sending data and wants to close the connection.
  * **Initiator State (Before -> After):** The initiator transitions from the ESTABLISHED state to the FIN-WAIT-1 state.
  * **Responder State (Before -> After):** The responder remains in the ESTABLISHED state.

  **Stage 2: Finish Responder (ACK)**

  * **Actor:** The responder acknowledges the initiator's FIN.
  * **TCP Segment Sent:** The responder sends an ACK segment to acknowledge the received FIN.
  * **Initiator State (Before -> After):** The initiator stays in the FIN-WAIT-1 state until it receives this ACK.
  * **Responder State (Before -> After):** The responder transitions from the ESTABLISHED state to the CLOSE-WAIT state, indicating it has received the FIN and will eventually close its side.

  **Stage 3: Finish Responder (FIN)**

  * **Actor:** The responder sends its own FIN to close its side of the connection.
  * **TCP Segment Sent:** The responder sends a FIN segment.
  * **Initiator State (Before -> After):** The initiator can be in either FIN-WAIT-1 or FIN-WAIT-2, depending on the timing of the ACKs. If the initiator has received the ACK to its FIN before the responder sends its FIN, it will be in FIN-WAIT-2.
  * **Responder State (Before -> After):** The responder transitions from the CLOSE-WAIT state to the LAST-ACK state.

  **Stage 4: Finish Initiator (ACK)**

  * **Actor:** The initiator acknowledges the responder's FIN.
  * **TCP Segment Sent:** The initiator sends an ACK segment to acknowledge the received FIN.
  * **Initiator State (Before -> After):** The initiator transitions from FIN-WAIT-1 or FIN-WAIT-2 to the TIME-WAIT state.
  * **Responder State (Before -> After):** The responder transitions from the LAST-ACK state to the CLOSED state.

  The TIME-WAIT state is a special state where the initiator waits for a certain period (2MSL, or Maximum Segment Lifetime) to ensure that any delayed packets are handled before fully closing the connection. After this period, the initiator transitions to the CLOSED state.


## UDP 

Here’s the improved summary:  

---

- **UDP is a connectionless protocol**, meaning it does not establish a connection before sending data.  
- **UDP does not guarantee** delivery, order, or error correction.  
- It is **faster than TCP** due to minimal overhead.  
- UDP is commonly used in **real-time applications** where speed is more critical than reliability, such as **video streaming, online gaming, and VoIP**.  
- UDP transmits data in **discrete units called datagrams**.  
- If a datagram is lost, **UDP does not retransmit it**. However, some applications implement their **own error recovery mechanisms** (e.g., RTP for video streaming).  
- UDP supports **broadcasting and multicasting**, making it useful for **group communication**.  


## Differences between TCP and UDP

| Feature | TCP | UDP |
|---------|-----|-----|
| Connection | Connection-oriented | Connectionless |
| Reliability | Reliable (guarantees delivery) | Unreliable (no delivery guarantee) |
| Order | Ordered delivery | Unordered delivery |
| Flow Control | Yes | No |
| Error Correction | Yes | No |
| Speed | Slower (due to overhead) | Faster (minimal overhead) |
| Use Cases | Web browsing, file transfer, email | Video streaming, online gaming, VoIP |
| Data Transmission | Segments | Datagrams |
| Broadcasting | No | Yes |
| Multicasting | No | Yes |

