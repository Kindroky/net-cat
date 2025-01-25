# TCP Chat Application

A multi-client TCP chat server built in Go, replicating the TCP chat functionality of `Netcat`. This application allows multiple clients to connect, send messages, and interact with each other in real-time. It features message broadcasting, logging, and basic command handling.

## Features

- **Netcat-like Functionality:** Provides a chat interface similar to `Netcat`'s TCP communication feature.
- **Real-time Messaging:** Broadcast messages to all connected clients.
- **User Commands:**
  - `/exit`: Disconnect from the server.
  - `/rename`: Change your username.
- **Welcome Message:** Displays an ASCII art logo upon connection.
- **Message Logs:** Logs all messages and user activity to a `logs.txt` file.
- **Concurrency Handling:** Uses `sync.Mutex` to manage shared resources safely.
- **Connection Limit:** Supports up to 10 simultaneous client connections.

## Usage

### Server
1. Start the server with a specific port ; in your terminal, type : 
   ```bash
   go run . 2525
   ```
2. Then open 1 or more terminals (representing people chatting with each other) and type :
    ```bash
    nc localhost 2525
    ```
3. You can also chat amongst different devices, provided you are all connected to the same wifi. In that case one should type the following when connecting to the server :
    ```nc [IP address of the server] [provided port]```

