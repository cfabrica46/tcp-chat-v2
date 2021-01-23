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

		reader := bufio.NewReader(os.Stdin)

		fmt.Printf("%v ingresó al chat\n", name)
		fmt.Println()

		go esperandoMensaje(conn)

		for {

			message, err := reader.ReadString('\n')

			if err != nil {
				log.Fatal(err)
			}

			content := fmt.Sprintf("%s: %s", name, message)

			_, err = conn.Write([]byte(content))

			if err != nil {
				log.Fatal(err)
			}

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

		fmt.Println(mensajeRecivido)
	}
}
