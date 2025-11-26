package interfaces

import "fmt"

// User é uma interface que define o metodo Login / interface vazia
// Toda interface vazia é implementada por qualquer tipo

func Arq3() {
	fmt.Println("----------- Inicio do programa -----------")

	var lista []interface{}

	lista = append(lista, 10)
	lista = append(lista, "Hello")
	lista = append(lista, true)
	lista = append(lista, 10.5)

	for i, usuario := range lista {
		fmt.Println(i, usuario)
	}

	fmt.Println("----------- Fim do programa -----------")
}
