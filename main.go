package main

func main() {
	portStr := IsValidArgPort()
	if portStr == nil {
		return
	}
	listener := ServerCreation(portStr)
	file := CreateLogsFile()
	NewUserConnection(listener, file)
}
