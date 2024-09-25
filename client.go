package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"os"
	"strconv"
)

type Client struct {
	Conn     net.Conn
	Username string
	Reader   io.Reader
	Writer   io.Writer
	Message  string
}

var clients = make(map[net.Conn]Client)

func StructAndMap(connexion net.Conn) string {
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
	clients[connexion] = newClient
	return username
}

func HandleUsername(connexion net.Conn) string {
	connexion.Write([]byte("[ENTER YOUR NAME]:"))
	validUsername := true
	username := ""
	buf := bufio.NewScanner(connexion)
	for {
		buf.Scan()
		username = buf.Text()
		if username == "" {
			validUsername = false
		}
		for _, client := range clients {
			if username == client.Username {
				validUsername = false
				break
			}
		}
		if !validUsername {
			connexion.Write([]byte("Invalid Username\n"))
			username = HandleUsername(connexion)
		}
		break
	}
	return username
}

/* func HandleNewName(conn net.Conn) string {
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
} */

func HandleClient(structure Client, count *int, file *os.File) {
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
					DelogTransmission(structure, file)
					*count--
					structure.Conn.Close()
				} else if message == "/rename\n" {
					oldUsername := structure.Username
					structure.Username = StructAndMap(structure.Conn)
					RenameTransmission(structure, oldUsername, file)
				} else {
					structure.Conn.Write([]byte("Command not found\n"))
				}
			} else {
				// Regular message handling
				structure.Message = message
				Transmission(structure, file)
			}
		}
	}
}
