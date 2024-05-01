package main

import (
	"fmt"
	"net"
	"os"
)

var MessageBuffer []byte = make([]byte, 1024)

// Função responsavel por lidar com o envio de mensagens para o client
// Function responsible for handling sending messages to the client
func handleMessage(conn net.Conn, message string) {
	_, err := conn.Read(MessageBuffer)
	if err != nil {
		panic(err)
	}

	_, err = conn.Write([]byte(message))
	if err != nil {
		panic(err)
	}
}

func handleConnection(port string) net.Conn {
	for {

		l, err := net.Listen("tcp", port)
		if err != nil {
			fmt.Println("Failed to bind to port 6379")
			os.Exit(1)
		}

		conn, err := l.Accept()
		if err != nil {
			fmt.Println("Error accepting connection: ", err.Error())
			os.Exit(1)
		}

		handleMessage(conn, "+PONG\r\n")
	}
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	handleConnection("0.0.0.0:6379")

}
