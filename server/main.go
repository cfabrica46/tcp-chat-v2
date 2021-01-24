package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

var conexiones []conexion

type conexion struct {
	id   int
	conn net.Conn
}

func main() {

	var id int

	if len(os.Args) == 3 {

		serverAddress := fmt.Sprintf("%s:%s", os.Args[1], os.Args[2])

		listener, err := net.Listen("tcp", serverAddress)

		if err != nil {
			log.Fatal(err)
		}

		defer listener.Close()

		fmt.Println("Listening on: " + serverAddress)
		fmt.Println()

		for {

			fmt.Println("esperando una conexion...")
			fmt.Println()

			c, err := listener.Accept()

			if err != nil {
				log.Fatal(err)
			}

			conn := conexion{id: id, conn: c}

			conexiones = append(conexiones, conn)

			fmt.Printf("Un usuario entro al chat\n")

			go util(conn, id)

			id++

		}

	}
}

func util(conn conexion, id int) {

	defer conn.conn.Close()
	reader := bufio.NewReader(conn.conn)

	for {

		content, err := reader.ReadString('\n')
		if err != nil {
			if strings.Contains(err.Error(), "host") {
				break
			} else if err == io.EOF {
				for i := range conexiones {
					if conexiones[i].id == conn.id {
						if len(conexiones) >= i+1 {
							conexiones = append(conexiones[:i], conexiones[i+1:]...)
							break
						}
					}
				}

				break
			} else {
				log.Fatal(err)
			}
		}

		for i := range conexiones {

			if conexiones[i].id != conn.id {
				_, err = conexiones[i].conn.Write([]byte(content))

				if err != nil {
					log.Fatal(err)
				}
			}
		}

		fmt.Println(string(content))
	}
}
