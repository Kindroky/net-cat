package main

import (
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	portStr := IsValidArgPort()
	if portStr == nil {
		return
	}
	listener := ServerCreation(portStr)
	file := CreateLogsFile()
	NewUserConnection(listener, file)
}
func IsValidArgPort() *string {
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
func ServerCreation(portStr *string) net.Listener {
	if *portStr == ":" {
		*portStr += "8989"
	}
	listener, err := net.Listen("tcp", *portStr)
	if err != nil {
		log.Fatal(err)
	}
	return listener
}
func NewUserConnection(listener net.Listener, file *os.File) {
	count := 0
	for {
		connexion, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
		}
		if count < 10 {
			count++
			LePingouin(connexion)
			StructAndMap(connexion)
			// Non fonctionnel
			/* logsToTransmit := ReadLogs(file)
			PreviousLogsTransmission(clients[connexion], logsToTransmit) */
			go HandleClient(clients[connexion], &count, file)
			LogTransmission(clients[connexion], file)
		} else {
			connexion.Write([]byte("Maximum connections reached"))
			connexion.Close()
		}
	}
}

func CreateLogsFile() *os.File {
	file, err := os.OpenFile("logs.txt", os.O_RDWR|os.O_CREATE|os.O_APPEND|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

// Non fonctionnel
/* func ReadLogs(file *os.File) *string {
	bufLogsFile := bufio.NewScanner(file)
	bufLogsFile.Scan()
	logs := bufLogsFile.Text()
	return &logs
} */

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

func HandleExit(con net.Conn) {
	con.Write([]byte("Exiting..."))
	con.Close()
}
