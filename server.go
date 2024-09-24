package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
)

type Client struct {
	conn     net.Conn
	username string
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
		if count <= 10 {
			connexion, err := listener.Accept()
			if err != nil {
				log.Fatal(err)
			}
			count++
			go HandleClient(connexion, &count)
		}
	}
}

func HandleClient(con net.Conn, count *int) {
	con.Write([]byte(strconv.Itoa(*count) + "\n"))
	var message string
	//tab := []byte{}
	/*con.Read(tab)
	con.Write(tab)*/
	scanner := bufio.NewScanner(con)
	connected := true
	for connected {
		scanner.Scan()
		message = scanner.Text() + "\n"
		if message != "\n" {
			if message[0] == '/' {
				if message == "/exit\n" {
					*count--
					connected = false
					con.Close()
				} else {
					con.Write([]byte("Command not found\n"))
				}
			} else {
				con.Write([]byte(message))
			}
		}
	}
}
