package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {

	var name string

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

			enviarMensaje(conn, name)

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

func enviarMensaje(conn net.Conn, name string) {

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("> ")

	message, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println()

	content := fmt.Sprintf("%s: %s", name, message)

	_, err = conn.Write([]byte(content))

	if err != nil {
		log.Fatal(err)
	}

}
