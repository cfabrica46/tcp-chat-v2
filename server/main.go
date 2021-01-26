package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
	"time"
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

			c, err := listener.Accept()

			if err != nil {
				log.Fatal(err)
			}

			fmt.Println("alguien ha igresado")

			conexiones = append(conexiones, c)

			go util(c)

		}

	}
}

func util(conn net.Conn) {

	defer conn.Close()

	name := recibirNombre(conn)

	dispacherMessage(conn, name)

	reader := bufio.NewReader(conn)

	for {

		content, err := reader.ReadString('\n')

		if err != nil {

			if err == io.EOF {

				eliminarConexion(conn)

			} else if strings.Contains(err.Error(), "host") {

				eliminarConexion(conn)

			} else {

				log.Fatal(err)

			}

			enviarMensajeDesconexion(name)

			break
		}

		enviarMensaje(conn, content)

	}
}

func recibirNombre(conn net.Conn) (name string) {

	reader := bufio.NewReader(conn)

	name, err := reader.ReadString('\n')

	if err != nil {

		if err == io.EOF {

			fmt.Println("conexion fallida")

		} else if strings.Contains(err.Error(), "host") {

			fmt.Println("conexion fallida")

		} else {

			log.Fatal(err)

		}
	}
	name = strings.Replace(name, "\n", "", -1)

	return
}

func dispacherMessage(conn net.Conn, name string) {

	mensajeIngreso := fmt.Sprintf("%s ha ingresado\n", name)

	for i := range conexiones {

		if conexiones[i] != conn {
			_, err := conexiones[i].Write([]byte(mensajeIngreso))

			if err != nil {
				log.Fatal(err)
			}
		} else {
			_, err := conexiones[i].Write([]byte("Has ingresado\n"))

			if err != nil {
				log.Fatal(err)
			}
		}
	}
}

func eliminarConexion(conn net.Conn) {

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

}

func enviarMensajeDesconexion(name string) {

	for i := range conexiones {
		desconectado := fmt.Sprintf("%s se a desconectado\n", name)

		_, err := conexiones[i].Write([]byte(desconectado))
		if err != nil {
			log.Fatal(err)
		}

	}

}

func enviarMensaje(conn net.Conn,content string) {

	for i := range conexiones {

		if conexiones[i] != conn {
			t := time.Now()

			mensaje := fmt.Sprintf("%v:%v %v", t.Hour(), t.Minute(), content)

			_, err := conexiones[i].Write([]byte(mensaje))

			if err != nil {
				log.Fatal(err)
			}
		}
	}

}
