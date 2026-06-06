package main

import (
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
	clientsCount := 10000

	var wg sync.WaitGroup

	for i := 1; i <= clientsCount; i++ {
		wg.Add(1)

		go func(id int) {
			defer wg.Done()

			conn, err := net.Dial("tcp", "127.0.0.1:6379")
			if err != nil {
				fmt.Printf("[CLIENT %d] Falha ao conectar: %v\n", id, err)
				return
			}

			defer conn.Close()

			message := fmt.Sprintf("PING %d\n", id)
			_, err = conn.Write([]byte(message))
			if err != nil {
				fmt.Printf("[CLIENT %d] Falha ao enviar mensagem: %v\n", id, err)
				return
			}

			buffer := make([]byte, 1024)

			err = conn.SetReadDeadline(time.Now().Add(2 * time.Second))
			if err != nil {
				fmt.Printf("[CLIENT %d] erro de timeout: %v\n", id, err)
				return
			}

			c, err := conn.Read(buffer)
			if err != nil {
				fmt.Printf("[CLIENT %d] Falha ao receber resposta: %v\n", id, err)
				return
			}

			fmt.Printf("[CLIENT %d] Mensagem recebida: %s\n", id, buffer[:c])
		}(i)
	}

	wg.Wait()
	fmt.Println("Todos os clientes foram executados")
}
