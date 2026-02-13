package main

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// ============================================================================
// MÓDULO 11 — Concorrência
// ============================================================================
//
// Este é o SUPERPODER do Go. Concorrência é nativa da linguagem.
//
// Conceitos:
//   - Goroutines (funções executando em paralelo)
//   - Channels (comunicação entre goroutines)
//   - select (multiplexação de channels)
//   - sync.WaitGroup, sync.Mutex, sync.RWMutex
//   - Patterns: fan-out, fan-in, worker pool, pipeline
//   - context.Context (cancelamento e timeouts)
//
// Rode com: go run exemplo11_concorrencia.go
// ============================================================================

func main() {
	// ==========================================================
	fmt.Println("=== GOROUTINES BÁSICAS ===")

	// Goroutine = função rodando concorrentemente
	// Custam ~2KB de stack (threads do OS custam ~1MB)
	// O runtime do Go gerencia o scheduling

	go func() {
		fmt.Println("  Olá de uma goroutine!")
	}()

	// Precisamos esperar, senão main termina antes da goroutine
	time.Sleep(100 * time.Millisecond)

	// ==========================================================
	fmt.Println("\n=== WAITGROUP ===")

	var wg sync.WaitGroup

	for i := 1; i <= 5; i++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
			fmt.Printf("  Worker %d finalizado\n", id)
		}(i)
	}

	wg.Wait() // Espera TODAS terminarem
	fmt.Println("  Todas as goroutines finalizaram!")

	// ==========================================================
	fmt.Println("\n=== CHANNELS ===")

	// Channel UNBUFFERED — bloqueia até alguém receber
	ch := make(chan string)

	go func() {
		ch <- "mensagem da goroutine"
	}()

	msg := <-ch // Recebe (bloqueia até ter valor)
	fmt.Println("  Recebido:", msg)

	// Channel BUFFERED — só bloqueia quando cheio
	buffered := make(chan int, 3)
	buffered <- 1
	buffered <- 2
	buffered <- 3
	// buffered <- 4 // Bloquearia! Buffer cheio

	fmt.Printf("  Buffer: %d, %d, %d\n", <-buffered, <-buffered, <-buffered)

	// ==========================================================
	fmt.Println("\n=== CHANNEL COM RANGE ===")

	numeros := make(chan int)

	go func() {
		for i := 1; i <= 5; i++ {
			numeros <- i
		}
		close(numeros) // IMPORTANTE: fechar quando terminar
	}()

	fmt.Print("  Números: ")
	for n := range numeros { // range em channel — para quando fechar
		fmt.Printf("%d ", n)
	}
	fmt.Println()

	// ==========================================================
	fmt.Println("\n=== CHANNEL DIRECIONAL ===")

	// Channels podem ser restringidos a envio ou recebimento
	pipelineDemo()

	// ==========================================================
	fmt.Println("\n=== SELECT ===")

	ch1 := make(chan string)
	ch2 := make(chan string)

	go func() {
		time.Sleep(100 * time.Millisecond)
		ch1 <- "canal 1"
	}()

	go func() {
		time.Sleep(200 * time.Millisecond)
		ch2 <- "canal 2"
	}()

	// Receber de quem chegar primeiro
	for i := 0; i < 2; i++ {
		select {
		case msg := <-ch1:
			fmt.Println("  Recebido de ch1:", msg)
		case msg := <-ch2:
			fmt.Println("  Recebido de ch2:", msg)
		}
	}

	// Select com timeout
	ch3 := make(chan string)
	select {
	case msg := <-ch3:
		fmt.Println(msg)
	case <-time.After(100 * time.Millisecond):
		fmt.Println("  Timeout! Ninguém enviou em 100ms")
	}

	// ==========================================================
	fmt.Println("\n=== MUTEX ===")

	var (
		contador int
		mu       sync.Mutex
		wg2      sync.WaitGroup
	)

	// Sem mutex, teríamos race condition!
	for i := 0; i < 1000; i++ {
		wg2.Add(1)
		go func() {
			defer wg2.Done()
			mu.Lock()
			contador++
			mu.Unlock()
		}()
	}

	wg2.Wait()
	fmt.Printf("  Contador (com Mutex): %d (esperado: 1000)\n", contador)

	// ==========================================================
	fmt.Println("\n=== WORKER POOL ===")
	workerPoolDemo()

	// ==========================================================
	fmt.Println("\n=== CONTEXT (cancelamento) ===")
	contextDemo()
}

// --- Pipeline: produtor → processador → consumidor ---
func produtor(out chan<- int) { // chan<- = só pode ENVIAR
	for i := 1; i <= 5; i++ {
		out <- i
	}
	close(out)
}

func duplicar(in <-chan int, out chan<- int) { // <-chan = só RECEBER, chan<- = só ENVIAR
	for v := range in {
		out <- v * 2
	}
	close(out)
}

func pipelineDemo() {
	nums := make(chan int)
	dobrados := make(chan int)

	go produtor(nums)
	go duplicar(nums, dobrados)

	fmt.Print("  Pipeline: ")
	for v := range dobrados {
		fmt.Printf("%d ", v)
	}
	fmt.Println()
}

// --- Worker Pool ---
func workerPoolDemo() {
	const numWorkers = 3
	const numJobs = 10

	jobs := make(chan int, numJobs)
	results := make(chan int, numJobs)

	// Lançar workers
	var wg sync.WaitGroup
	for w := 1; w <= numWorkers; w++ {
		wg.Add(1)
		go func(id int) {
			defer wg.Done()
			for job := range jobs {
				// Simular trabalho
				time.Sleep(50 * time.Millisecond)
				result := job * job
				fmt.Printf("  Worker %d: job %d → %d\n", id, job, result)
				results <- result
			}
		}(w)
	}

	// Enviar jobs
	for j := 1; j <= numJobs; j++ {
		jobs <- j
	}
	close(jobs)

	// Esperar workers
	go func() {
		wg.Wait()
		close(results)
	}()

	// Coletar resultados
	total := 0
	for r := range results {
		total += r
	}
	fmt.Printf("  Soma dos resultados: %d\n", total)
}

// --- Context ---
func contextDemo() {
	// Context com timeout
	ctx, cancel := context.WithTimeout(context.Background(), 300*time.Millisecond)
	defer cancel()

	resultCh := make(chan string, 1)

	go func() {
		// Simular operação que pode demorar
		time.Sleep(time.Duration(rand.Intn(500)) * time.Millisecond)
		resultCh <- "operação concluída"
	}()

	select {
	case result := <-resultCh:
		fmt.Println(" ", result)
	case <-ctx.Done():
		fmt.Println("  Context cancelado:", ctx.Err())
	}

	// Context com cancelamento manual
	ctx2, cancel2 := context.WithCancel(context.Background())

	go func() {
		time.Sleep(100 * time.Millisecond)
		cancel2() // Cancela manualmente
	}()

	<-ctx2.Done()
	fmt.Println("  Context2 cancelado:", ctx2.Err())
}
