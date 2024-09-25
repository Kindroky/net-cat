package main

import (
	"fmt"
	"log"
	"os"
	"time"
)

// funcion that transmits a client's message to everybody else
func Transmission(clientstruct Client, file *os.File) {
	fmtMessage := fmt.Sprintf("[%s][%s]: %s", Time(), clientstruct.Username, clientstruct.Message)
	_, err := file.Write([]byte(fmtMessage))
	if err != nil {
		log.Fatal(err)
	}
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's arrival in the chat to everybody else
func LogTransmission(clientstruct Client, file *os.File) {
	fmtMessage := fmt.Sprintf("Yay! %s has joined the chat!\n", clientstruct.Username)
	_, err := file.Write([]byte(fmtMessage))
	if err != nil {
		log.Fatal(err)
	}
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's exit of the chat to everybody else
func DelogTransmission(clientstruct Client, file *os.File) {
	fmtMessage := fmt.Sprintf("Unfortunately, %s has left us...\n", clientstruct.Username)
	_, err := file.Write([]byte(fmtMessage))
	if err != nil {
		log.Fatal(err)
	}
	for _, client := range clients {
		if client.Username != clientstruct.Username {
			client.Conn.Write([]byte(fmtMessage))
		}
	}
}

// function that transmits a client's rename in the chat to everybody
func RenameTransmission(clientstruct Client, oldUsername string, file *os.File) {
	fmtMessage := fmt.Sprintf("Wow, %s is now known as %s !\n", oldUsername, clientstruct.Username)
	_, err := file.Write([]byte(fmtMessage))
	if err != nil {
		log.Fatal(err)
	}
	for _, client := range clients {
		client.Conn.Write([]byte(fmtMessage))
	}
}

// Non fonctionnel
/* func PreviousLogsTransmission(clientstruct Client, logs *string) {
	clientstruct.Conn.Write([]byte(*logs))
} */

// function that computes and properly formats the date and time
func Time() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
