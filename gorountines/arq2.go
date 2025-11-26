package gorountines

import (
	"fmt"
	"sync"
	"time"
)

func ExemploGoroutine2() {

	var wg sync.WaitGroup // WaitGroup = Aguarda todas as goroutines terminarem
	wg.Add(4)             // Add = Adiciona o numero de goroutines que serao executadas
	go ex1(&wg)           // Cria uma goroutine para executar a funcao ex1
	go ex2(&wg)           // Cria uma goroutine para executar a funcao ex2
	go ex3(&wg)           // Cria uma goroutine para executar a funcao ex3
	wg.Wait()             // Wait = Aguarda todas as goroutines terminarem
}

func ex1(wg *sync.WaitGroup) { // Funcao que executa a goroutine 1
	time.Sleep(5 * time.Second) // Sleep = Espera 5 segundos
	fmt.Println("Exemplo de goroutine 1")
	wg.Done() // Done = Decrementa o numero de goroutines que serao executadas
}

func ex2(wg *sync.WaitGroup) { // Funcao que executa a goroutine 2
	time.Sleep(3 * time.Second)
	fmt.Println("Exemplo de goroutine 2")
	wg.Done() // Done = Decrementa o numero de goroutines que serao executadas
}

func ex3(wg *sync.WaitGroup) { // Funcao que executa a goroutine 3
	time.Sleep(1 * time.Second)
	fmt.Println("Exemplo de goroutine 3")
	wg.Done() // Done = Decrementa o numero de goroutines que serao executadas
}
