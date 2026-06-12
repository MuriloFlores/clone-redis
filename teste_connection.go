package main

import (
	"fmt"
	"math/rand"
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

func generateRandomCommand() string {
	method := methods[rand.Intn(len(methods))]

	if method == "GET" {
		mu.Lock()
		length := len(reqs)

		if length > 0 {
			index := rand.Intn(length)
			field := reqs[index][0]
			mu.Unlock()

			return buildRESP(method, field)
		}

		mu.Unlock()

		field := fields[rand.Intn(len(fields))]
		return buildRESP(method, field)
	}

	field := fields[rand.Intn(len(fields))]
	value := keys[rand.Intn(len(keys))]

	mu.Lock()
	reqs = append(reqs, []string{
		field, value,
	})
	mu.Unlock()

	return buildRESP(method, field, value)
}

var mu sync.Mutex
var reqs [][]string

var methods []string = []string{
	"GET", "SET",
}

var fields []string = []string{
	"foo",
	"bar",
	"banana",
	"maca",
	"hello",
}

var keys []string = []string{
	"melancia",
	"world",
	"peixe",
	"macaxeira",
	"vermelho",
}

func main() {
	clientsCount := 100

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

			message := generateRandomCommand()

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
