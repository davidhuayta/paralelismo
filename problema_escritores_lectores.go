package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"sync"
	"time"
)

// funciones para leer en un archivo
func leer() string {

	file_data, err := ioutil.ReadFile("./tex.txt")
	if err != nil {
		fmt.Println("hubo un error")
	}
	fmt.Println(string(file_data))
	return ""
}

//debuelve el contenido del archivo
func tex() string {
	file_data, err := ioutil.ReadFile("./tex.txt")
	if err != nil {
		fmt.Println("hubo un error")
	}
	return string(file_data)
}

// funciones para escribir en un archivo
func Escribir(nuevo string, texto string) string {

	message := []byte(texto + nuevo)
	err := ioutil.WriteFile("./tex.txt", message, 0644)

	if err != nil {
		log.Fatal(err)
	}
	return ""
}

//funcion mutex para el acceso a la base de datos cuando se escribe
type AccesDB struct {
	sync.Mutex
}

//estructura usuario para escritores y lectores
type User struct {
	nombre          string
	mensaje         string   //mensaje q escribira el escritor
	TipoUser, count int      //tipo_user: 0=Lector , 1=Escritor
	BD              *AccesDB //mutex
}

//funcion principal q ejecuta las acciones segun el tipo de usuario
func (p User) Ejecutar(c chan *User) {
	for i := 0; i < 5; i++ { //numerop de accesos permitidos a la base de datos todos los usuarios
		c <- &p
		if p.TipoUser == 0 { //si es lector ,leer
			fmt.Println("El Lector", p.nombre, "esta leyendo")

		}
		//seccIon critica
		if p.TipoUser == 1 { //si es escritor ,escritor

			if p.count < 2 { //limita la edicion para cada editor en 2
				p.BD.Lock()
				fmt.Println("El escritor", p.nombre, "edita por", p.count+1, " vez")
				fmt.Println("El escritor", p.nombre, "escribe ,", p.mensaje, "en la base de datos")
				Escribir(p.nombre+" dice "+p.mensaje+"\n", tex())
				p.count = p.count + 1
				//wg.Done() //disminuir el numero de ediciones totales a la base de datos
				fmt.Println("El escritor", p.nombre, " termino de Escribir")
				p.BD.Unlock()
			}

		}
	}

}

//bloquea el canal
func despejar(c chan *User) {
	for {
		if len(c) == 1 {
			<-c
			time.Sleep(20 * time.Millisecond)
		}
	}
}
func main() {
	//incializacion de variables
	var wait string
	var i int
	//var wg sync.WaitGroup
	//canal c para enviara usuarios, de 1 en 1
	c := make(chan *User, 1)
	//numero de escrituras que se admiten
	//wg.Add(4)
	//inicializacion de mutex para la base de datos
	DB := make([]*AccesDB, 1)
	DB[0] = new(AccesDB)
	//nombre de pila
	nombre := [8]string{"Pepe", "Pablo", "Mario", "Luigi", "Leo", "Toño", "Matias", "Alarak"}
	mensaje := [8]string{"Helo word", "goland", "mensaje 1", "hola mundo", "mensaje 1", "conspiraciones", "diasdas", "213233312"}
	//declaracion e  inicializacion de de usurios(lectores y escritores)
	user := make([]*User, 8)
	for i = 0; i < 8; i++ {
		user[i] = &User{nombre[i], mensaje[i], i % 2, 0, DB[0]}
	}
	// bloqueo de canales
	go despejar(c)

	//ejecuta los usuarios
	for i = 0; i < 8; i++ {
		go user[i].Ejecutar(c)
	}
	//wg.Wait()

	//
	fmt.Scanln(&wait)
}
