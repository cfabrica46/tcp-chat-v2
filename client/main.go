package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
)

func main() {

	var name string

	if len(os.Args) == 3 {

		serverAddress := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])

		conn, err := net.Dial("tcp", serverAddress)

		if err != nil {
			log.Fatal(err)
		}

		err = reciveMessageName(conn)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Scan(&name)

		conn.Write([]byte(name + "\n"))

		fmt.Println()

		go esperandoMensaje(conn)

		for {

			err = enviarMensaje(conn, name)

			if err != nil {
				log.Fatal(err)
			}
		}
	}

}

func reciveMessageName(conn net.Conn) (err error) {

	reader := bufio.NewReader(conn)

	mensaje, err := reader.ReadString('\n')
	if err != nil {
		return
	}

	fmt.Println(mensaje)

	return
}

func esperandoMensaje(conn net.Conn) (err error) {

	r := bufio.NewReader(conn)

	for {

		mensajeRecivido, err := r.ReadString('\n')

		if err != nil {
			return err
		}

		mensajeRecivido = strings.Replace(mensajeRecivido, "\n", "", -1)

		fmt.Printf("\r%s", mensajeRecivido)

		if strings.Contains(mensajeRecivido, ":") {
			fmt.Println()
			fmt.Println()

		}
	}
}

func enviarMensaje(conn net.Conn, name string) (err error) {

	reader := bufio.NewReader(os.Stdin)

	message, err := reader.ReadString('\n')

	if err != nil {
		return
	}

	fmt.Println()

	_, err = conn.Write([]byte(message))

	if err != nil {
		return
	}

	return
}
