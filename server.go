package main

import (
	"bufio"
	"log"
	"net"
	"strconv"
	//"time"
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

/*
func Transmission(sender Client, message string, clients map) {
	for _, client := range clients {
		if client != sender {
			client.conn.Write(fmt.Sprintf("[%s][%s]: %s", currentTime(), username, []byte(message + "\n")))
		}
	}
}

func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
*/
