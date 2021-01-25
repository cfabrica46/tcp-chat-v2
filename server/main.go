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

var conexiones []net.Conn

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

		for {

			fmt.Println("esperando una conexion...")
			fmt.Println()

			c, err := listener.Accept()

			if err != nil {
				log.Fatal(err)
			}

			conexiones = append(conexiones, c)

			name := recivirNombre(c)

			mensajeIngreso := fmt.Sprintf("%s a ingresado\n", name)

			for i := range conexiones {

				if conexiones[i] != c {
					_, err = conexiones[i].Write([]byte(mensajeIngreso))

					if err != nil {
						log.Fatal(err)
					}
				} else {
					_, err = conexiones[i].Write([]byte("Has ingresado\n"))

					if err != nil {
						log.Fatal(err)
					}
				}
			}

			go util(c, name)

		}

	}
}

func util(conn net.Conn, name string) {

	defer conn.Close()
	reader := bufio.NewReader(conn)

	for {

		content, err := reader.ReadString('\n')

		if err != nil {

			if err == io.EOF {

				for i := range conexiones {
					if conexiones[i] == conn {
						if len(conexiones) >= i+1 {
							conexiones = append(conexiones[:i], conexiones[i+1:]...)
							break
						}
						if len(conexiones) == 1 {
							conexiones = []net.Conn{}
						}
					}

				}

			} else if strings.Contains(err.Error(), "host") {

				for i := range conexiones {
					if conexiones[i] == conn {
						if len(conexiones) >= i+1 {
							conexiones = append(conexiones[:i], conexiones[i+1:]...)
							break
						}
					}

				}

			} else {

				log.Fatal(err)

			}

			for i := range conexiones {

				desconectado := fmt.Sprintf("%s se a desconectado\n", name)

				_, err = conexiones[i].Write([]byte(desconectado))
				if err != nil {
					log.Fatal(err)
				}

			}
			break
		}

		for i := range conexiones {

			if conexiones[i] != conn {
				_, err = conexiones[i].Write([]byte(content))

				if err != nil {
					log.Fatal(err)
				}
			}
		}

		fmt.Println(string(content))
	}
}

func recivirNombre(conn net.Conn) (name string) {

	reader := bufio.NewReader(conn)

	name, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	name = strings.Replace(name, "\n", "", -1)

	return
}
