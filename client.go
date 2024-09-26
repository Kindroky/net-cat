package main

import (
	"bufio"
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
	return username
}

func HandleClient(structure Client, count *int, file *os.File) {
	defer func() {
		structure.Conn.Close()
		DelogTransmission(structure, file)
		delete(clients, structure.Conn)
		*count--
	}()

	// Send initial message with client count
	structure.Conn.Write([]byte(strconv.Itoa(*count) + "\n"))
	var message string
	bufClient := bufio.NewScanner(structure.Reader)
	for {
		// Read the client's message
		scan := bufClient.Scan()
		if !scan {
			return
		}
		message = bufClient.Text() + "\n"

		// Ignore empty messages
		if message != "\n" {
			// Handle commands
			if message[0] == '/' {
				if message == "/exit\n" {
					return
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
