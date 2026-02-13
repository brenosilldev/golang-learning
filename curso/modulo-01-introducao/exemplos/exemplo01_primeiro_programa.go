package main

import (
	"fmt"
	"runtime"
	"time"
)

// ============================================================================
// MÓDULO 01 — Primeiro Programa em Go
// ============================================================================
//
// Conceitos demonstrados:
//   - package main → todo executável precisa desse pacote
//   - func main()  → ponto de entrada do programa
//   - import       → importar pacotes (só pode importar o que usa!)
//   - fmt.Println  → imprimir com quebra de linha
//   - fmt.Printf   → imprimir formatado (como printf do C)
//   - :=           → declaração curta de variável (infere o tipo)
//
// Rode com: go run exemplo01_primeiro_programa.go
// ============================================================================

func main() {
	// Declaração curta de variáveis (o Go infere o tipo automaticamente)
	nome := "Gopher"
	linguagem := "Go"

	// fmt.Println → imprime e pula linha
	fmt.Println("=== Meu Primeiro Programa Go ===")
	fmt.Println()

	// fmt.Printf → imprime formatado (precisa de \n manual)
	// %s = string, %d = inteiro, %f = float
	fmt.Printf("👋 Olá! Eu sou %s\n", nome)
	fmt.Printf("🔧 Estou aprendendo %s\n", linguagem)
	fmt.Println()

	// Usando pacotes da stdlib
	fmt.Printf("🐹 Versão do Go: %s\n", runtime.Version())
	fmt.Printf("💻 Sistema: %s/%s\n", runtime.GOOS, runtime.GOARCH)
	fmt.Printf("🧵 CPUs disponíveis: %d\n", runtime.NumCPU())
	fmt.Printf("⏰ Data/Hora: %s\n", time.Now().Format("02/01/2006 15:04:05"))
}
