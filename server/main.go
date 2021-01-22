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
	if len(os.Args) == 3 {

		serverAddress := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])

		listener, err := net.Listen("tcp", serverAddress)

		if err != nil {
			log.Fatal(err)
		}

		defer listener.Close()

		fmt.Println("Listening on: " + serverAddress)
		fmt.Println()

		c := make(chan []byte)

		for {

			fmt.Println("esperando una conexion...")
			fmt.Println()

			conn, err := listener.Accept()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Printf("Un usuario entro al chat\n")

			go aux(conn, c)

		}
	}
}

func aux(conn net.Conn, c chan []byte) {

	reader := bufio.NewReader(conn)

	for {

		go reciveMensaje(conn, c)

		content, err := reader.ReadString('\n')

		if err != nil {
			if strings.Contains(err.Error(), "host") {
				break
			} else {
				log.Fatal(err)
			}
		}

		c <- []byte(content)

		fmt.Println(string(content))
	}
}

func reciveMensaje(conn net.Conn, c chan []byte) {

	for {
		mensaje := <-c

		_, err := conn.Write(mensaje)

		if err != nil {
			log.Fatal(err)
		}
	}
}
