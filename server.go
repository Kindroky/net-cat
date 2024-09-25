package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Conn     net.Conn
	Username string
	Reader   io.Reader
	Writer   io.Writer
	Message  string
}

var clients = make(map[string]Client)

func main() {
	listener, err := net.Listen("tcp", ":2525")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	count := 0
	for {
		connexion, err := listener.Accept()
		if count < 10 {
			if err != nil {
				log.Fatal(err)
			}
			count++
			username := HandleUsername(connexion)
			r := bufio.NewReader(connexion)
			w := bufio.NewWriter(connexion)
			newClient := Client{
				connexion,
				username,
				r,
				w,
				"",
			}
			clients[username] = newClient
			if username != "" {
				go HandleClient(clients[username], &count)
				Logtransmission(clients[username])
			}
		} else {
			connexion.Write([]byte("Maximum connections reached"))
			connexion.Close()
		}
	}
}

func HandleUsername(conn net.Conn) string {
	conn.Write([]byte(`Welcome to TCP-Chat!
         _nnnn_
        dGGGGMMb
       @p~qp~~qMb
       M|@||@) M|
       @,----.JM|
      JS^\__/  qKL
     dZP        qKRb
    dZP          qKKb
   fZP            SMMb
   HZM            MMMM
   FqM            MMMM
 __| ".        |\dS"qML
 |    ` + "`" + `.       | ` + "`" + `' \Zq
_)      \.___.,|     .'
\____   )MMMMMP|   .'
     ` + "`" + `-'       ` + "`" + `--'
[ENTER YOUR NAME]:`))
	buf := bufio.NewScanner(conn)
	for {
		buf.Scan()
		name := buf.Text() + ""
		if name != "" {
			return name
		} else {
			HandleUsername(conn)
		}
	}
}

func HandleNewName(conn net.Conn) string {
	conn.Write([]byte("[ENTER YOUR NAME]:"))
	buf := bufio.NewScanner(conn)
	for {
		buf.Scan()
		name := buf.Text() + ""
		if name != "" {
			return name
		} else {
			HandleNewName(conn)
		}
	}
}

func HandleClient(structure Client, count *int) {
	defer structure.Conn.Close()

	// Send initial message with client count
	structure.Conn.Write([]byte(strconv.Itoa(*count) + "\n"))
	var message string
	bufClient := bufio.NewScanner(structure.Reader)
	for {
		// Send the formatted message every time before reading input
		fmtMessage := fmt.Sprintf("[%s][%s]: ", Time(), structure.Username)
		structure.Conn.Write([]byte(fmtMessage))

		// Read the client's message
		bufClient.Scan()
		message = bufClient.Text() + "\n"

		// Ignore empty messages
		if message != "\n" {
			// Handle commands
			if message[0] == '/' {
				if message == "/exit\n" {
					Delogtransmission(structure)
					*count--
					structure.Conn.Close()
				} else if message == "/rename\n" {
					structure.Username = HandleNewName(structure.Conn)
				} else {
					structure.Conn.Write([]byte("Command not found\n"))

				}
			} else {
				// Regular message handling
				structure.Message = message
				Transmission(structure)
			}
		}
	}
}

// funcion that transmits a client's message to everybody else
func Transmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("[%s][%s]: %s", Time(), clientstruct.Username, clientstruct.Message)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's arrival in the chat to everybody else
func Logtransmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("Yay! %s has joined the chat!\n", clientstruct.Username)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's exit of the chat to everybody else
func Delogtransmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("Unfortunately, %s has left us...\n", clientstruct.Username)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that computes and properly formats the date and time
func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func HandleExit(con net.Conn) {
	con.Write([]byte("Exiting..."))
	con.Close()
}
