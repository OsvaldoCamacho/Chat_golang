package main

import (
	"bufio"
	"encoding/gob"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strings"
)

type Mensaje struct {
	Id         string
	Ingreso    bool
	Salir      bool
	Mensaje    string
	Nombrefile string
	Data       []byte
	Archivo    bool
}

type Usuario struct {
	Nickname string
	Id       string
	Conexion net.Conn
}

var miusuario Usuario
var serverDesconectado = false
var allmensajes []string
var scanner = bufio.NewScanner(os.Stdin)

func cliente(p Mensaje) {

	fmt.Println(miusuario.Conexion.LocalAddr().String())

	go handleCliente()

}
func handleCliente() {

	for {
		var mensaje Mensaje
		var user string
		err := gob.NewDecoder(miusuario.Conexion).Decode(&mensaje)
		if err != nil {
			fmt.Println(err)
		} else {
			if mensaje.Archivo == false {
				if mensaje.Salir {
					println("El Servidor se Desconecto")
					println("Desconectando del servidor....")
					serverDesconectado = true
					break

				}

				if mensaje.Id == miusuario.Nickname {

					user = "Tu"
				} else {

					user = mensaje.Id
				}
				msg := user + ": " + mensaje.Mensaje
				fmt.Println(msg)
				allmensajes = append(allmensajes, msg)
			} else {

				if mensaje.Id == miusuario.Nickname {

					user = "Tu"
				} else {

					user = mensaje.Id
					err = ioutil.WriteFile(mensaje.Nombrefile, mensaje.Data, 0644)
					if err != nil {
						log.Fatal(err)
					}

				}
				msg := user + ":File " + mensaje.Nombrefile
				fmt.Println(msg)
				allmensajes = append(allmensajes, msg)

			}
		}
	}

	return

}

var opcion int

func main() {
	salir := false
	fmt.Print("Nombre de usuario:")
	scanner.Scan()
	nickname := scanner.Text()

	p := Mensaje{
		Id:         nickname,
		Ingreso:    true,
		Salir:      false,
		Mensaje:    "nuevo",
		Nombrefile: "",
		Data:       nil,
		Archivo:    false,
	}

	c, err := net.Dial("tcp", ":9999")
	if err != nil {
		fmt.Println(err)

	}
	err = gob.NewEncoder(c).Encode(p)
	miusuario.Conexion = c
	miusuario.Id = p.Id
	miusuario.Nickname = p.Id

	go handleCliente()

	for !salir {
		if serverDesconectado {
			break
		}
		fmt.Println("1) Enviar Mensaje")
		fmt.Println("2) Enviar Archivo")
		fmt.Println("3) Mostrar Mensajes")
		fmt.Println("4) Salir")
		fmt.Scanln(&opcion)
		salir = menu(opcion)
	}
	var input string
	fmt.Scanln(&input)

}

func menu(o int) bool {

	switch o {

	case 1:

		fmt.Print("Mensaje:")
		scanner.Scan()
		line := scanner.Text()

		msg := Mensaje{
			Id:         miusuario.Nickname,
			Ingreso:    false,
			Salir:      false,
			Mensaje:    line,
			Nombrefile: "",
			Data:       nil,
			Archivo:    false,
		}
		c, err := net.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)

		}

		err = gob.NewEncoder(c).Encode(&msg)
		if err != nil {
			fmt.Println(err)

		}
		c.Close()

		return false
	case 2:
		cargarArchivo()

		return false
	case 3:
		for _, e := range allmensajes {
			fmt.Println(e)

		}
		return false
	case 4:
		msg := Mensaje{
			Id:         miusuario.Nickname,
			Ingreso:    false,
			Salir:      true,
			Mensaje:    miusuario.Nickname + "se Desconecto",
			Nombrefile: "",
			Data:       nil,
			Archivo:    false,
		}
		c, err := net.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)

		}

		err = gob.NewEncoder(c).Encode(&msg)
		if err != nil {
			fmt.Println(err)

		}
		c.Close()

		return true
	default:
		if serverDesconectado {
			return false
		} else {
			fmt.Println("opcion invalida")
			return false
		}

	}

}

func cargarArchivo() {
	var ruta string
	fmt.Println("ingresa la ruta del archivo a enviar:")
	fmt.Scanln(&ruta)
	var aux []string

	aux = strings.Split(ruta, "\\")

	nombre := aux[len(aux)-1]

	data, err := ioutil.ReadFile(ruta)

	if err != nil {
		fmt.Println("Error no existe ese archivo o directorio")
	} else {
		Sdata := []byte(data)
		msg := Mensaje{
			Id:         miusuario.Nickname,
			Ingreso:    false,
			Salir:      false,
			Mensaje:    "Archivo:",
			Nombrefile: nombre,
			Data:       Sdata,
			Archivo:    true,
		}
		c, err := net.Dial("tcp", ":9999")
		if err != nil {
			fmt.Println(err)

		}

		err = gob.NewEncoder(c).Encode(&msg)
		fmt.Println("Se a Enviado", msg.Nombrefile)
		if err != nil {
			fmt.Println(err)

		}
		c.Close()

	}

}
