package gorountines

import "fmt"

/**

Goroutine = Funcao que roda em paralelo
Go = Cria uma goroutine


*/

func ExemploGoroutine() {
	for i := 0; i < 50; i++ {
		go ex1Simples() // Cria uma goroutine para executar a funcao ex1
		go ex2Simples() // Cria uma goroutine para executar a funcao ex2
		go ex3Simples() // Cria uma goroutine para executar a funcao ex3
	}
}

func ex1Simples() {
	fmt.Println("Exemplo de goroutine 1")
}

func ex2Simples() {
	fmt.Println("Exemplo de goroutine 2")
}

func ex3Simples() {
	fmt.Println("Exemplo de goroutine 3")
}
