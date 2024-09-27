package main

import (
	"bufio"
	"io"
	"net"
	"os"
)

type Client struct {
	Conn     net.Conn
	Username string
	Reader   io.Reader
	Writer   io.Writer
	Message  string
}

var clients = make(map[net.Conn]Client) // map with connection as key and structure as value

func StructAndMap(connexion net.Conn) string { // add each client to the map
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
	mapMu.Lock()
	clients[connexion] = newClient
	mapMu.Unlock()
	return username
}

func HandleUsername(connexion net.Conn) string { // Ask the client for its username
	connexion.Write([]byte("[ENTER YOUR NAME]: "))
	validUsername := true
	username := ""
	buf := bufio.NewScanner(connexion)
	buf.Scan()
	username = buf.Text()
	if username == "" {
		validUsername = false
	}
	mapMu.Lock()
	for _, client := range clients {
		if username == client.Username {
			validUsername = false
			break
		}
	}
	mapMu.Unlock()
	if !validUsername {
		connexion.Write([]byte("Invalid Username\n"))
		username = HandleUsername(connexion)
	}
	return username
}

func HandleClient(structure Client, count *int, file *os.File) {
	// delog part
	defer func() {
		structure.Conn.Close()
		DelogTransmission(structure, file)
		delete(clients, structure.Conn)
		countMu.Lock()
		*count--
		countMu.Unlock()
	}()

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

// display linux logo on entry
func LePingouin(connexion net.Conn) {
	connexion.Write([]byte(`Welcome to TCP-Chat!
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
`))
}
