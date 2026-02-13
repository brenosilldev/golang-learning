package main

import "fmt"

// ============================================================================
// MÓDULO 06 — Ponteiros
// ============================================================================
//
// Conceitos:
//   - & (endereço de) e * (valor apontado por)
//   - Tipo ponteiro (*int, *string, etc.)
//   - Passagem por valor vs por referência
//   - Ponteiros com structs (açúcar sintático)
//   - nil pointer e segurança
//   - new() vs &
//
// Rode com: go run exemplo06_ponteiros.go
// ============================================================================

func main() {
	// ==========================================================
	fmt.Println("=== PONTEIROS BÁSICOS ===")

	x := 42
	p := &x // p é um *int (ponteiro para int) que aponta para x

	fmt.Printf("x  = %d\n", x)
	fmt.Printf("&x = %p (endereço de x)\n", &x)
	fmt.Printf("p  = %p (valor de p = endereço de x)\n", p)
	fmt.Printf("*p = %d (valor apontado por p)\n", *p)

	// Alterar o valor via ponteiro
	*p = 100
	fmt.Printf("\nApós *p = 100:\n")
	fmt.Printf("x  = %d (mudou!)\n", x)
	fmt.Printf("*p = %d\n", *p)

	// ==========================================================
	fmt.Println("\n=== PASSAGEM POR VALOR vs REFERÊNCIA ===")

	valor := 10
	fmt.Printf("Antes: %d\n", valor)

	tentarAlterar(valor)
	fmt.Printf("Após tentarAlterar (por valor): %d (NÃO mudou)\n", valor)

	alterarDeVerdade(&valor)
	fmt.Printf("Após alterarDeVerdade (por ref): %d (MUDOU!)\n", valor)

	// ==========================================================
	fmt.Println("\n=== PONTEIROS COM STRUCTS ===")

	type Pessoa struct {
		Nome  string
		Idade int
	}

	// Criar struct
	pessoa := Pessoa{Nome: "Go", Idade: 15}
	fmt.Println("Antes:", pessoa)

	// Ponteiro para struct
	pp := &pessoa

	// Go tem açúcar sintático: pp.Nome é equivalente a (*pp).Nome
	pp.Nome = "Golang"
	pp.Idade = 16
	fmt.Println("Depois:", pessoa) // Pessoa original foi modificada!

	// Função que modifica struct via ponteiro
	fazerAniversario(&pessoa)
	fmt.Println("Após aniversário:", pessoa)

	// ==========================================================
	fmt.Println("\n=== new() vs & ===")

	// new(T) aloca memória e retorna ponteiro para zero value
	pInt := new(int)    // *int apontando para 0
	pStr := new(string) // *string apontando para ""
	pBool := new(bool)  // *bool apontando para false

	fmt.Printf("new(int):    %d (tipo: %T)\n", *pInt, pInt)
	fmt.Printf("new(string): '%s' (tipo: %T)\n", *pStr, pStr)
	fmt.Printf("new(bool):   %t (tipo: %T)\n", *pBool, pBool)

	// & faz a mesma coisa, mas com valor inicial
	pInt2 := &[]int{1, 2, 3} // ponteiro para slice
	fmt.Printf("&slice: %v\n", *pInt2)

	// ==========================================================
	fmt.Println("\n=== NIL POINTER ===")

	// Ponteiro não inicializado é nil
	var pNil *int
	fmt.Printf("Ponteiro nil: %v\n", pNil)

	// Acessar *pNil causaria PANIC! Sempre verifique:
	if pNil != nil {
		fmt.Println(*pNil)
	} else {
		fmt.Println("Ponteiro é nil — não pode desreferenciar!")
	}

	// ==========================================================
	fmt.Println("\n=== SWAP COM PONTEIROS ===")

	a, b := 10, 20
	fmt.Printf("Antes: a=%d, b=%d\n", a, b)
	swap(&a, &b)
	fmt.Printf("Depois: a=%d, b=%d\n", a, b)

	// ==========================================================
	fmt.Println("\n=== PONTEIRO PARA PONTEIRO ===")

	val := 42
	p1 := &val // *int
	p2 := &p1  // **int

	fmt.Printf("val = %d\n", val)
	fmt.Printf("*p1 = %d\n", *p1)
	fmt.Printf("**p2 = %d\n", **p2)

	**p2 = 999
	fmt.Printf("val após **p2 = 999: %d\n", val)
}

// Recebe por valor — NÃO modifica o original
func tentarAlterar(n int) {
	n = 999 // Modifica apenas a cópia local
}

// Recebe ponteiro — MODIFICA o original
func alterarDeVerdade(n *int) {
	*n = 999
}

// Pattern comum: método que modifica struct via ponteiro
func fazerAniversario(p *struct {
	Nome  string
	Idade int
}) {
	p.Idade++
}

// Swap clássico com ponteiros
func swap(a, b *int) {
	*a, *b = *b, *a
}
