package avancado

import (
	"os"
)

func ExemploPanic() {

	_, err := os.Open("arquivo.txt")

	if err != nil {
		panic(err) // Panic = Gera um erro fatal e para a execucao do programa
	}

}
