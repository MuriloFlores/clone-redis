package main

import (
	"fmt"
	"net"
	"strings"
	"sync"
	"time"
)

func buildRESP(args ...string) string {
	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("*%d\r\n", len(args)))
	for _, arg := range args {
		sb.WriteString(fmt.Sprintf("$%d\r\n%s\r\n", len(arg), arg))
	}
	return sb.String()
}

func main() {
	clientsCount := 1

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

			message := buildRESP("SET", "foo", "bar")

			// message := buildRESP("GET", "foo")

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

			fmt.Printf("[CLIENT %d] Resposta do servidor: %q\n", id, buffer[:c])
		}(i)
	}

	wg.Wait()
	fmt.Println("Todos os clientes foram executados")
}
