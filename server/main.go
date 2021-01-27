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

	name, err := recibirNombre(conn)

	if err != nil {
		log.Fatal(err)
	}

	err = dispacherMessage(conn, name)

	if err != nil {
		log.Fatal(err)
	}

	reader := bufio.NewReader(conn)

	for {

		_, err = conn.Write([]byte(">\n"))

		if err != nil {
			log.Fatal(err)
		}

		content, err := reader.ReadString('\n')

		if err != nil {

			if err == io.EOF {

				eliminarConexion(conn)

			} else if strings.Contains(err.Error(), "host") {

				eliminarConexion(conn)

			} else {

				log.Fatal(err)

			}

			err = enviarMensajeDesconexion(name)

			if err != nil {
				log.Fatal(err)
			}

			break
		}

		err = enviarMensaje(conn, content, name)

		if err != nil {
			log.Fatal(err)
		}

		for i := range conexiones {
			if conexiones[i] != conn {
				conexiones[i].Write([]byte(">\n"))
			}
		}

	}
}

func recibirNombre(conn net.Conn) (name string, err error) {

	reader := bufio.NewReader(conn)

	name, err = reader.ReadString('\n')

	if err != nil {

		if err == io.EOF {

			fmt.Println("conexion fallida")

		} else if strings.Contains(err.Error(), "host") {

			fmt.Println("conexion fallida")

		} else {

			return

		}
	}
	name = strings.Replace(name, "\n", "", -1)

	return
}

func dispacherMessage(conn net.Conn, name string) (err error) {

	mensajeIngreso := fmt.Sprintf("%s ha ingresado\n", name)

	for i := range conexiones {

		if conexiones[i] != conn {
			_, err := conexiones[i].Write([]byte(mensajeIngreso))

			if err != nil {
				return err
			}
		} else {
			_, err := conexiones[i].Write([]byte("Has ingresado\n"))

			if err != nil {
				return err
			}
		}
	}
	return
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

func enviarMensajeDesconexion(name string) (err error) {

	for i := range conexiones {
		desconectado := fmt.Sprintf("%s se a desconectado\n", name)

		_, err := conexiones[i].Write([]byte(desconectado))
		if err != nil {
			return err
		}

	}
	return
}

func enviarMensaje(conn net.Conn, content, name string) (err error) {

	for i := range conexiones {

		if conexiones[i] != conn {
			t := time.Now()

			mensaje := fmt.Sprintf("%v:%v %v: %v", t.Hour(), t.Minute(), name, content)

			_, err := conexiones[i].Write([]byte(mensaje))

			if err != nil {
				return err
			}
		}
	}
	return

}
