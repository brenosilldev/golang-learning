package main

import (
	"fmt"
	"strings"
)

// ============================================================================
// MÓDULO 05 — Funções
// ============================================================================
//
// Conceitos:
//   - Funções com múltiplos retornos
//   - Retornos nomeados
//   - Funções variádicas (...T)
//   - Funções como valores (first-class)
//   - Closures
//   - defer (execução adiada)
//   - init() (inicialização de pacote)
//
// Rode com: go run exemplo05_funcoes.go
// ============================================================================

// --- Função simples ---
func soma(a, b int) int {
	return a + b
}

// --- Múltiplos retornos ---
// Pattern mais importante do Go: retornar (resultado, error)
func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("divisão por zero")
	}
	return a / b, nil
}

// --- Retornos nomeados ---
func calculaRetangulo(largura, altura float64) (area float64, perimetro float64) {
	area = largura * altura
	perimetro = 2 * (largura + altura)
	return // "naked return" — retorna as variáveis nomeadas
}

// --- Função variádica ---
// O último parâmetro com ... aceita 0 ou mais argumentos
func somaTotal(nums ...int) int {
	total := 0
	for _, n := range nums {
		total += n
	}
	return total
}

// --- Função que recebe outra função como parâmetro ---
func aplicar(slice []int, fn func(int) int) []int {
	resultado := make([]int, len(slice))
	for i, v := range slice {
		resultado[i] = fn(v)
	}
	return resultado
}

// --- Closure (função que captura variáveis externas) ---
func criarContador(inicio int) func() int {
	valor := inicio
	return func() int {
		valor++
		return valor
	}
}

// --- Closure prática: criar multiplicador ---
func multiplicador(fator int) func(int) int {
	return func(n int) int {
		return n * fator
	}
}

// --- Função que demonstra defer ---
func exemploDefer() {
	fmt.Println("  Início da função")

	// defer executa ao final da função, na ordem LIFO (último a entrar, primeiro a sair)
	defer fmt.Println("  defer 1 (último a executar)")
	defer fmt.Println("  defer 2")
	defer fmt.Println("  defer 3 (primeiro a executar)")

	fmt.Println("  Meio da função")
	fmt.Println("  Fim da função")
	// Agora os defers executam em ordem reversa
}

// --- defer prático: medir tempo de execução ---
func operacaoLenta() {
	// Esse pattern é MUITO usado em Go para medir performance
	// defer avalia os argumentos na hora, mas executa ao final
	defer fmt.Println("  → operacaoLenta() finalizada")
	fmt.Println("  → operacaoLenta() iniciada")

	// Simular trabalho pesado
	resultado := 0
	for i := 0; i < 1000000; i++ {
		resultado += i
	}
	fmt.Printf("  → Resultado: %d\n", resultado)
}

// --- Tipo de função customizado ---
type Transformador func(string) string

func aplicarTransformacoes(texto string, transformacoes ...Transformador) string {
	resultado := texto
	for _, t := range transformacoes {
		resultado = t(resultado)
	}
	return resultado
}

func main() {
	// ==========================================================
	fmt.Println("=== FUNÇÃO SIMPLES ===")
	fmt.Printf("soma(3, 4) = %d\n", soma(3, 4))

	// ==========================================================
	fmt.Println("\n=== MÚLTIPLOS RETORNOS ===")

	resultado, err := divide(10, 3)
	if err != nil {
		fmt.Println("Erro:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", resultado)
	}

	// Usando _ para ignorar um retorno
	_, err2 := divide(10, 0)
	if err2 != nil {
		fmt.Println("Erro esperado:", err2)
	}

	// ==========================================================
	fmt.Println("\n=== RETORNOS NOMEADOS ===")

	area, perimetro := calculaRetangulo(5, 3)
	fmt.Printf("Retângulo 5x3: Área=%.1f, Perímetro=%.1f\n", area, perimetro)

	// ==========================================================
	fmt.Println("\n=== FUNÇÕES VARIÁDICAS ===")

	fmt.Println("soma()        =", somaTotal())
	fmt.Println("soma(1)       =", somaTotal(1))
	fmt.Println("soma(1,2,3)   =", somaTotal(1, 2, 3))
	fmt.Println("soma(1..10)   =", somaTotal(1, 2, 3, 4, 5, 6, 7, 8, 9, 10))

	// Passando slice para variádica (spread com ...)
	numeros := []int{10, 20, 30}
	fmt.Println("soma(slice...) =", somaTotal(numeros...))

	// ==========================================================
	fmt.Println("\n=== FUNÇÕES COMO VALORES ===")

	nums := []int{1, 2, 3, 4, 5}

	// Passando função anônima como argumento
	dobrados := aplicar(nums, func(n int) int {
		return n * 2
	})
	fmt.Println("Dobrados:", dobrados)

	quadrados := aplicar(nums, func(n int) int {
		return n * n
	})
	fmt.Println("Quadrados:", quadrados)

	// ==========================================================
	fmt.Println("\n=== CLOSURES ===")

	// Cada chamada a criarContador retorna uma closure independente
	contA := criarContador(0)
	contB := criarContador(100)

	fmt.Println("A:", contA(), contA(), contA()) // 1, 2, 3
	fmt.Println("B:", contB(), contB())          // 101, 102
	fmt.Println("A:", contA())                   // 4 (continua de onde parou!)

	// Multiplicador (closure prática)
	dobro := multiplicador(2)
	triplo := multiplicador(3)
	fmt.Println("dobro(5)  =", dobro(5))  // 10
	fmt.Println("triplo(5) =", triplo(5)) // 15

	// ==========================================================
	fmt.Println("\n=== DEFER ===")

	exemploDefer()

	fmt.Println("\nDefer prático:")
	operacaoLenta()

	// ==========================================================
	fmt.Println("\n=== TIPO DE FUNÇÃO CUSTOMIZADO ===")

	resultado2 := aplicarTransformacoes("  Hello, World!  ",
		strings.TrimSpace, // remove espaços
		strings.ToUpper,   // maiúsculas
		func(s string) string { // adicionar prefixo
			return ">>> " + s
		},
	)
	fmt.Println(resultado2) // >>> HELLO, WORLD!

	// ==========================================================
	fmt.Println("\n=== FUNÇÃO ANÔNIMA IMEDIATA (IIFE) ===")

	// Definir e executar na mesma hora
	mensagem := func(nome string) string {
		return fmt.Sprintf("Olá, %s!", nome)
	}("Gopher")

	fmt.Println(mensagem)
}
