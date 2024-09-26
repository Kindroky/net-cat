package main

import (
	"fmt"
	"log"
	"os"
	"sync"
	"time"
)

var mapMu sync.Mutex // Pour la map clients
var logMu sync.Mutex // Pour Logs

// function that transmits a client's message to everybody else
func Transmission(clientstruct Client, file *os.File) {
	fmtMessage := fmt.Sprintf("[%s][%s]: %s", Time(), clientstruct.Username, clientstruct.Message)
	logMu.Lock()
	_, err := file.Write([]byte(fmtMessage))
	Logs += fmtMessage
	logMu.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	mapMu.Lock()
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		} else {
			client.Conn.Write([]byte("\r\033[1A\033[K" + fmtMessage))
		}
	}
	mapMu.Unlock()
}

// function that transmits a client's arrival in the chat to everybody else
func LogTransmission(clientstruct Client, file *os.File) {
	fmtMessage := fmt.Sprintf("Yay! %s has joined the chat!\n", clientstruct.Username)
	logMu.Lock()
	_, err := file.Write([]byte(fmtMessage))
	Logs += fmtMessage
	logMu.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	mapMu.Lock()
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		} else {
			client.Conn.Write([]byte("\r\033[1A\033[K" + fmtMessage))
		}
	}
	mapMu.Unlock()
}

// function that transmits a client's exit of the chat to everybody else
func DelogTransmission(clientstruct Client, file *os.File) {
	fmtMessage := fmt.Sprintf("Unfortunately, %s has left us...\n", clientstruct.Username)
	logMu.Lock()
	_, err := file.Write([]byte(fmtMessage))
	Logs += fmtMessage
	logMu.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	mapMu.Lock()
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		} else {
			client.Conn.Write([]byte("\r\033[1A\033[K" + fmtMessage))
		}
	}
	mapMu.Unlock()
}

// function that transmits a client's rename in the chat to everybody
func RenameTransmission(clientstruct Client, oldUsername string, file *os.File) {
	fmtMessage := fmt.Sprintf("Wow, %s is now known as %s !\n", oldUsername, clientstruct.Username)
	logMu.Lock()
	_, err := file.Write([]byte(fmtMessage))
	Logs += fmtMessage
	logMu.Unlock()
	if err != nil {
		log.Fatal(err)
	}
	mapMu.Lock()
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		} else {
			client.Conn.Write([]byte("\r\033[1A\033[K" + fmtMessage))
		}
	}
	mapMu.Unlock()
}

// function that computes and properly formats the date and time
func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
