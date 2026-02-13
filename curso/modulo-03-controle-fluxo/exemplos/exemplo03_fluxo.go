package main

import "fmt"

// ============================================================================
// MÓDULO 03 — Controle de Fluxo: if, switch, for
// ============================================================================
//
// Conceitos:
//   - if / else if / else
//   - if com inicializador (pattern do Go)
//   - switch (com e sem expressão)
//   - for (clássico, while-like, infinito, range)
//   - break, continue, labels
//
// Rode com: go run exemplo03_fluxo.go
// ============================================================================

func main() {
	// ==========================================================
	fmt.Println("=== IF / ELSE ===")

	idade := 20

	if idade < 18 {
		fmt.Println("Menor de idade")
	} else if idade < 65 {
		fmt.Println("Adulto")
	} else {
		fmt.Println("Idoso")
	}

	// if com inicializador — a variável só existe dentro do if
	// Muito usado com funções que retornam (valor, erro) ou (valor, ok)
	numeros := map[string]int{"um": 1, "dois": 2, "tres": 3}

	if valor, existe := numeros["dois"]; existe {
		fmt.Printf("Encontrei: %d\n", valor)
	} else {
		fmt.Println("Não encontrei")
	}
	// 'valor' e 'existe' NÃO existem aqui fora!

	// ==========================================================
	fmt.Println("\n=== SWITCH ===")

	dia := "quarta"

	// Switch com expressão
	switch dia {
	case "segunda", "terca", "quarta", "quinta", "sexta":
		fmt.Println("Dia útil")
	case "sabado", "domingo":
		fmt.Println("Final de semana!")
	default:
		fmt.Println("Dia inválido")
	}

	// Switch SEM expressão (funciona como if/else encadeado)
	nota := 8.5

	switch {
	case nota >= 9:
		fmt.Println("Conceito A")
	case nota >= 7:
		fmt.Println("Conceito B")
	case nota >= 5:
		fmt.Println("Conceito C")
	default:
		fmt.Println("Reprovado")
	}

	// Switch com type assertion (veremos mais em interfaces)
	var valor interface{} = 42
	switch v := valor.(type) {
	case int:
		fmt.Printf("É int: %d\n", v)
	case string:
		fmt.Printf("É string: %s\n", v)
	default:
		fmt.Printf("Tipo desconhecido: %T\n", v)
	}

	// ==========================================================
	fmt.Println("\n=== FOR CLÁSSICO ===")

	// For tradicional (C-style)
	for i := 0; i < 5; i++ {
		fmt.Printf("i = %d\n", i)
	}

	// ==========================================================
	fmt.Println("\n=== FOR COMO WHILE ===")

	// Go NÃO tem while. Usa for sem a parte de init e post.
	contador := 0
	for contador < 3 {
		fmt.Printf("contador = %d\n", contador)
		contador++
	}

	// ==========================================================
	fmt.Println("\n=== FOR INFINITO ===")

	// for {} é um loop infinito — use break para sair
	tentativas := 0
	for {
		tentativas++
		if tentativas >= 3 {
			fmt.Println("Saindo do loop infinito após 3 tentativas")
			break
		}
	}

	// ==========================================================
	fmt.Println("\n=== FOR RANGE ===")

	// range em slice — retorna (índice, valor)
	frutas := []string{"🍎 Maçã", "🍌 Banana", "🍊 Laranja"}
	for i, fruta := range frutas {
		fmt.Printf("[%d] %s\n", i, fruta)
	}

	// Ignorando o índice com _
	fmt.Println("\nSó valores:")
	for _, fruta := range frutas {
		fmt.Println(fruta)
	}

	// range em map — retorna (chave, valor)
	capitais := map[string]string{
		"Brasil":    "Brasília",
		"Argentina": "Buenos Aires",
		"Japão":     "Tóquio",
	}
	for pais, capital := range capitais {
		fmt.Printf("%s → %s\n", pais, capital)
	}

	// range em string — retorna (índice do byte, rune)
	for i, char := range "Go é 🔥" {
		fmt.Printf("byte[%d] = %c (rune %d)\n", i, char, char)
	}

	// ==========================================================
	fmt.Println("\n=== BREAK E CONTINUE ===")

	// continue pula para a próxima iteração
	fmt.Println("Números pares de 0 a 9:")
	for i := 0; i < 10; i++ {
		if i%2 != 0 {
			continue // pula ímpares
		}
		fmt.Printf("%d ", i)
	}
	fmt.Println()

	// ==========================================================
	fmt.Println("\n=== LABELS (para loops aninhados) ===")

	// Labels permitem dar break/continue em loops externos
externo:
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			if i == 1 && j == 1 {
				fmt.Println("Quebrando loop externo!")
				break externo // sai dos DOIS loops
			}
			fmt.Printf("i=%d, j=%d\n", i, j)
		}
	}
}
