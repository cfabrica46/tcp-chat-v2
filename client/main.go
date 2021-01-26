package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"sync"
)

func main() {

	var name string
	var m sync.Mutex

	if len(os.Args) == 3 {

		fmt.Println("ingrese su nombre")

		fmt.Scan(&name)

		serverAddress := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])

		conn, err := net.Dial("tcp", serverAddress)

		if err != nil {
			log.Fatal(err)
		}

		conn.Write([]byte(name + "\n"))

		fmt.Println()

		go esperandoMensaje(conn)

		for {
			m.Lock()
			go enviarMensaje(conn, name, &m)

		}
	}

}

func esperandoMensaje(conn net.Conn) {

	r := bufio.NewReader(conn)

	for {
		mensajeRecivido, err := r.ReadString('\n')

		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("\r%s\n", mensajeRecivido)
		fmt.Print(">")
	}
}

func enviarMensaje(conn net.Conn, name string, m *sync.Mutex) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("> ")
	message, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	content := fmt.Sprintf("%s: %s", name, message)

	_, err = conn.Write([]byte(content))

	if err != nil {
		log.Fatal(err)
	}

	m.Unlock()
}
