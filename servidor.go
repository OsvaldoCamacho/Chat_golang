package main

import (
	"encoding/gob"
	"fmt"
	"net"
	"os"
)

var s = make([]*Usuario, 0)
var mostrar = true
var ultimo = 0
var allmensajes []string

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

func server() {
	s, err := net.Listen("tcp", ":9999")
	if err != nil {
		fmt.Println(err)
		return
	}
	for {
		c, err := s.Accept()
		if err != nil {

			fmt.Println(err)
			continue
		}
		go handleCliente(c)

	}
}

func handleCliente(c net.Conn) {
	var mensaje Mensaje
	err := gob.NewDecoder(c).Decode(&mensaje)
	if err != nil {
		fmt.Println(err)
	} else {
		//
		if mensaje.Ingreso {
			aux := Usuario{
				Nickname: mensaje.Id,
				Id:       mensaje.Id,
				Conexion: c,
			}
			s = append(s, &aux)
			fmt.Println(mensaje.Id + " Se conecto al servidor")

		} else {

			if mensaje.Salir {
				var a = make([]*Usuario, 0)
				for _, e := range s {
					if mensaje.Id != e.Id {
						a = append(a, e)
					}
				}
				s = a
				fmt.Println(mensaje.Id + " Se Desconecto")

			} else {
				if !mensaje.Archivo {
					fmt.Println(mensaje.Id + ": " + mensaje.Mensaje)
					msg := mensaje.Id + ": " + mensaje.Mensaje
					for _, e := range s {
						err := gob.NewEncoder(e.Conexion).Encode(&mensaje)

						if err != nil {
							fmt.Println(err)
						}
					}
					allmensajes = append(allmensajes, msg)
				} else {

					fmt.Println(mensaje.Id + " envio el Archivo  " + mensaje.Nombrefile)

					for _, e := range s {
						if mensaje.Id == e.Id {
							aux := Mensaje{
								Id:         mensaje.Id,
								Ingreso:    mensaje.Ingreso,
								Salir:      mensaje.Salir,
								Mensaje:    mensaje.Mensaje,
								Nombrefile: mensaje.Nombrefile,
								Data:       nil,
								Archivo:    mensaje.Archivo,
							}
							err := gob.NewEncoder(e.Conexion).Encode(&aux)
							msg := mensaje.Id + ": " + mensaje.Nombrefile
							allmensajes = append(allmensajes, msg)
							if err != nil {
								fmt.Println(err)
							}

						} else {
							err := gob.NewEncoder(e.Conexion).Encode(&mensaje)
							msg := mensaje.Id + ": " + mensaje.Nombrefile
							allmensajes = append(allmensajes, msg)
							if err != nil {
								fmt.Println(err)
							}

						}

					}

				}
			}

		}

	}

}

func main() {
	salir := false
	go server()
	var opcion int
	for !salir {

		fmt.Println("1.- Mostrar los mensajes/nombre de los archivos enviados")
		fmt.Println("2.- Respaldar en un archivo de texto los mensajes/nombre de los archivos enviados")
		fmt.Println("3.- Terminar Servidor")
		fmt.Scanln(&opcion)
		salir = menu(opcion)
	}
	var input string
	fmt.Scanln(&input)

}

func menu(o int) bool {

	switch o {

	case 1:

		fmt.Println("Mensajes y Archivos")
		fmt.Println("------------")
		for _, e := range allmensajes {
			fmt.Println(e)

		}

		return false
	case 2:

		file, err := os.Create("respaldo.txt")
		if err != nil {
			fmt.Println(err)
		}
		defer file.Close()
		for _, e := range allmensajes {
			bytes, err := file.WriteString(e + "\n")
			if err != nil {
				fmt.Println(err, bytes)
			}
		}

		return false
	case 3:
		for _, e := range s {
			aux := Mensaje{
				Id:         e.Id,
				Ingreso:    false,
				Salir:      true,
				Mensaje:    "bye",
				Nombrefile: "",
				Data:       nil,
				Archivo:    false,
			}
			err := gob.NewEncoder(e.Conexion).Encode(&aux)
			if err != nil {
				fmt.Println(err)
			}
		}
		return true

	default:
		fmt.Println("opcion invalida")
		return false

	}

}

func localiza(nick string) int {
	i := 0
	for _, e := range s {
		if e.Id == nick {

			return i
		} else {
			i++
		}

	}
	return i

}
