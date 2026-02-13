package main

import (
	"encoding/json"
	"fmt"
)

// ============================================================================
// MÓDULO 07 — Structs, Métodos e Composição
// ============================================================================
//
// Conceitos:
//   - Declaração de structs
//   - Métodos (value receiver vs pointer receiver)
//   - Composição (embedding) — a "herança" do Go
//   - Tags JSON
//   - Pattern NewXxx (construtor)
//   - Struct anônima
//
// Rode com: go run exemplo07_structs.go
// ============================================================================

// --- Struct básica ---
type Pessoa struct {
	Nome  string
	Idade int
	Email string
}

// --- Método com VALUE receiver (não modifica) ---
func (p Pessoa) Apresentar() string {
	return fmt.Sprintf("Olá, sou %s e tenho %d anos", p.Nome, p.Idade)
}

// --- Método com POINTER receiver (pode modificar) ---
func (p *Pessoa) Aniversario() {
	p.Idade++
}

func (p *Pessoa) SetEmail(email string) {
	p.Email = email
}

// --- Construtor (pattern Go) ---
// Go não tem construtores nativos. O pattern é criar uma função NewXxx
func NewPessoa(nome string, idade int) *Pessoa {
	return &Pessoa{
		Nome:  nome,
		Idade: idade,
	}
}

// --- Composição (embedding) ---
type Endereco struct {
	Rua    string
	Cidade string
	Estado string
	CEP    string
}

func (e Endereco) Completo() string {
	return fmt.Sprintf("%s, %s - %s, %s", e.Rua, e.Cidade, e.Estado, e.CEP)
}

// Funcionario "herda" Pessoa e Endereco via embedding
type Funcionario struct {
	Pessoa   // embedding — campos e métodos de Pessoa ficam disponíveis
	Endereco // embedding — campos e métodos de Endereco ficam disponíveis
	Cargo    string
	Salario  float64
}

func (f Funcionario) Ficha() string {
	return fmt.Sprintf("%s | Cargo: %s | Salário: R$%.2f | End: %s",
		f.Apresentar(), f.Cargo, f.Salario, f.Completo())
}

// --- JSON Tags ---
type Produto struct {
	ID        int     `json:"id"`
	Nome      string  `json:"nome"`
	Preco     float64 `json:"preco"`
	Descricao string  `json:"descricao,omitempty"` // omitempty = omitir se vazio
	Interno   string  `json:"-"`                   // "-" = nunca serializar
}

func main() {
	// ==========================================================
	fmt.Println("=== CRIANDO STRUCTS ===")

	// Forma 1: campos nomeados (recomendado)
	p1 := Pessoa{
		Nome:  "Alice",
		Idade: 30,
		Email: "alice@email.com",
	}
	fmt.Println(p1)

	// Forma 2: posicional (evitar — frágil)
	p2 := Pessoa{"Bob", 25, "bob@email.com"}
	fmt.Println(p2)

	// Forma 3: usando construtor
	p3 := NewPessoa("Carol", 28)
	fmt.Println(*p3)

	// Forma 4: zero value (todos os campos zerados)
	var p4 Pessoa
	fmt.Printf("Zero value: %+v\n", p4)

	// ==========================================================
	fmt.Println("\n=== MÉTODOS ===")

	pessoa := NewPessoa("Go", 15)
	fmt.Println(pessoa.Apresentar())
	fmt.Printf("Idade antes: %d\n", pessoa.Idade)

	pessoa.Aniversario()
	fmt.Printf("Idade depois: %d\n", pessoa.Idade)

	pessoa.SetEmail("go@golang.org")
	fmt.Printf("Email: %s\n", pessoa.Email)

	// ==========================================================
	fmt.Println("\n=== COMPOSIÇÃO (embedding) ===")

	func1 := Funcionario{
		Pessoa: Pessoa{
			Nome:  "Maria",
			Idade: 35,
			Email: "maria@empresa.com",
		},
		Endereco: Endereco{
			Rua:    "Rua das Flores, 123",
			Cidade: "São Paulo",
			Estado: "SP",
			CEP:    "01001-000",
		},
		Cargo:   "Engenheira de Software",
		Salario: 15000.00,
	}

	// Acessar campos de Pessoa DIRETAMENTE (sem func1.Pessoa.Nome)
	fmt.Println("Nome:", func1.Nome) // Funciona por causa do embedding
	fmt.Println("Cargo:", func1.Cargo)
	fmt.Println()
	fmt.Println("Ficha completa:")
	fmt.Println(func1.Ficha())

	// Também pode acessar explicitamente
	fmt.Println("\nEndereço:", func1.Endereco.Completo())

	// ==========================================================
	fmt.Println("\n=== JSON SERIALIZAÇÃO ===")

	produto := Produto{
		ID:      1,
		Nome:    "Notebook",
		Preco:   4599.90,
		Interno: "dado-secreto",
		// Descricao está vazia — será omitida pelo omitempty
	}

	// Struct → JSON
	jsonBytes, _ := json.MarshalIndent(produto, "", "  ")
	fmt.Println("Struct → JSON:")
	fmt.Println(string(jsonBytes))
	// Note: "interno" não aparece (tag json:"-")
	// Note: "descricao" não aparece (omitempty + vazio)

	// JSON → Struct
	jsonStr := `{"id": 2, "nome": "Mouse", "preco": 89.90, "descricao": "Mouse gamer"}`
	var produto2 Produto
	json.Unmarshal([]byte(jsonStr), &produto2)
	fmt.Printf("\nJSON → Struct: %+v\n", produto2)

	// ==========================================================
	fmt.Println("\n=== STRUCT ANÔNIMA ===")

	// Struct sem nome — útil para uso rápido/temporário
	config := struct {
		Host string
		Port int
	}{
		Host: "localhost",
		Port: 8080,
	}
	fmt.Printf("Config: %+v\n", config)

	// ==========================================================
	fmt.Println("\n=== SLICE DE STRUCTS ===")

	pessoas := []Pessoa{
		{Nome: "Alice", Idade: 30},
		{Nome: "Bob", Idade: 25},
		{Nome: "Carol", Idade: 28},
	}

	fmt.Println("Maiores de 27:")
	for _, p := range pessoas {
		if p.Idade > 27 {
			fmt.Printf("  %s (%d anos)\n", p.Nome, p.Idade)
		}
	}

	// ==========================================================
	fmt.Println("\n=== COMPARAÇÃO DE STRUCTS ===")

	// Structs são comparáveis se todos os campos forem comparáveis
	a := Pessoa{Nome: "Go", Idade: 15}
	b := Pessoa{Nome: "Go", Idade: 15}
	c := Pessoa{Nome: "Rust", Idade: 9}

	fmt.Printf("a == b: %t\n", a == b) // true
	fmt.Printf("a == c: %t\n", a == c) // false
}
