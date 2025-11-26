package genericsts

import (
	"fmt"
)

type Tipos interface { // constrstraint - restringe o tipo de dados que a funcao pode receber
	int | string
}

func reverse[T Tipos](slice []T) []T { // Funcao generica que recebe um slice de inteiros ou strings e retorna um slice de inteiros ou strings

	newSlice := make([]T, len(slice)) // Cria um novo slice com o mesmo tamanho do slice original
	//newInts := len(slice) - 1         // Pega o ultimo indice do slice

	for i := 0; i < len(slice); i++ {
		newSlice[i] = slice[len(slice)-1-i] // Adiciona o ultimo elemento do slice original no novo slice
		// newSlice[i] = slice[newInts]
		// newInts--
		// Decrementa o indice do novo slice
	}

	return newSlice // Retorna o novo slice
}

func Arq1() {
	fmt.Println(reverse([]int{1, 2, 3, 4, 5}))
}
