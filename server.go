package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"strconv"
	"time"
)

type Client struct {
	Conn     net.Conn
	Username string
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
		defer connexion.Close()
		if count <= 10 {
			if err != nil {
				log.Fatal(err)
			}
			count++
			usernumber := strconv.Itoa(count)
			v := Client{
				connexion, "User" + usernumber,
			}
			clients["User"+usernumber] = v
			go HandleClient(connexion, count)
		} else {
			bye := "Maximum connections reached"
			connexion.Write([]byte(bye))
		}
	}
}

// function that handles each client's activity on the server
func HandleClient(con net.Conn, count int) {
	con.Write([]byte(strconv.Itoa(count) + "\n"))
	var message string
	//tab := []byte{}
	/*con.Read(tab)
	con.Write(tab)*/
	Bonjour := bufio.NewScanner(con)
	for {
		Bonjour.Scan()
		message = Bonjour.Text() + "\n"
		if message != "\n" {
			con.Write([]byte(message))
		}
	}
}

// funcion that transmits a client's message to everybody else
func Transmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("[%s][%s]: %s\n", Time(), clientstruct.Username, message)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's arrival in the chat to everybody else
func Logtransmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("[%s]: Yay! %s has joined the chat!", Time(), clientstruct.Username)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's exit of the chat to everybody else
func Delogtransmission(clientstruct Client) {
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			fmtMessage := fmt.Sprintf("[%s]: Unfortunately, %s has left us...", Time(), clientstruct.Username)
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that computes and properly formats the date and time
func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
