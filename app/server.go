package main

import (
	"fmt"
	"net"
	"os"
)

var MessageBuffer chan []byte =  make(chan []byte)

//Função responsavel por lidar com o envio de mensagens para o client
//Function responsible for handling sending messages to the client
func handleMessage(conn net.Conn) {
	_, err := conn.Read(<-MessageBuffer)
	if err != nil {
		panic(err)
	}
	
	for _, message := range <-MessageBuffer {
		msg := string(message)
		fmt.Println(msg)
	}
	
}

func main() {
	fmt.Println("Logs from your program will appear here!")

	l, err := net.Listen("tcp", "0.0.0.0:6379")
	if err != nil {
		fmt.Println("Failed to bind to port 6379")
		os.Exit(1)
	}
	conn, err := l.Accept()
	if err != nil {
		fmt.Println("Error accepting connection: ", err.Error())
		os.Exit(1)
	}

	//enviando a mensagem 
	//seding message
	handleMessage(conn)
}
