package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

var Logs string

var mapMu sync.Mutex   // Pour la map clients
var logMu sync.Mutex   // Pour Logs
var countMu sync.Mutex // Pour la variable count
func IsValidArgPort() *string { // check if argument is a valide port
	portStr := ":"
	isValid := true
	args := os.Args[1:]
	if len(args) == 0 {
		return &portStr
	} else if len(args) > 1 || len(args[0]) != 4 {
		isValid = false
	} else {
		for _, b := range args[0] {
			if b < 48 || b > 57 {
				isValid = false
				break
			}
		}
	}
	if !isValid {
		fmt.Println("[USAGE]: ./TCPChat $port")
		return nil
	} else {
		portStr += args[0]
		return &portStr
	}
}
func ServerCreation(portStr *string) net.Listener { // create server
	if *portStr == ":" {
		*portStr += "8989" // default port
	}
	listener, err := net.Listen("tcp", *portStr)
	if err != nil {
		log.Fatal(err)
	}
	return listener
}
func NewUserConnection(listener net.Listener, file *os.File) { //manage new connections to the server
	count := 0
	for {
		connexion, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		if count > 9 {
			connexion.Write([]byte("Maximum connections reached"))
			connexion.Close()
			continue

		}
		countMu.Lock()
		count++
		countMu.Unlock()
		go func(connexion net.Conn) {
			LePingouin(connexion)
			StructAndMap(connexion)
			connexion.Write([]byte(Logs))
			go HandleClient(clients[connexion], &count, file)
			LogTransmission(clients[connexion], file)
		}(connexion)
	}
}

func CreateLogsFile() *os.File { // create logs text file
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}
