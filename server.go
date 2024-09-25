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
	conn.Write([]byte("Please enter username : "))
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
func HandleClient(structure Client, count *int) {
	defer structure.Conn.Close()
	structure.Conn.Write([]byte(strconv.Itoa(*count) + "\n"))
	var message string
	bufClient := bufio.NewScanner(structure.Reader)
	for {
		bufClient.Scan()
		message = bufClient.Text() + "\n"
		if message != "\n" {
			if message[0] == '/' {
				if message == "/exit\n" {
					Delogtransmission(structure)
					*count--
					structure.Conn.Close()
				} else {
					structure.Conn.Write([]byte("Command not found\n"))
				}
			} else {
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
			fmtMessage := fmt.Sprintf("[%s][%s]: %s\n", Time(), clientstruct.Username, clientstruct.Message)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's arrival in the chat to everybody else
func Logtransmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("[%s]: Yay! %s has joined the chat!\n", Time(), clientstruct.Username)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's exit of the chat to everybody else
func Delogtransmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("[%s]: Unfortunately, %s has left us...\n", Time(), clientstruct.Username)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that computes and properly formats the date and time
func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
