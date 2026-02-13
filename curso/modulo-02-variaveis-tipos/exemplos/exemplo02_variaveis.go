package main

import "fmt"

// ============================================================================
// MÓDULO 02 — Variáveis, Tipos, Constantes e Conversão
// ============================================================================
//
// Conceitos:
//   - var vs :=
//   - Tipos primitivos (int, float64, string, bool, byte, rune)
//   - Zero values (todo tipo tem valor padrão)
//   - Constantes e iota
//   - Conversão explícita de tipos
//   - fmt.Printf com verbos de formatação
//
// Rode com: go run exemplo02_variaveis.go
// ============================================================================

func main() {
	fmt.Println("=== VARIÁVEIS ===")

	// --- Declaração com var ---
	var nome string = "Go"
	var idade int = 15 // Go foi lançada em 2009
	var ativa bool = true

	fmt.Printf("Linguagem: %s, Idade: %d, Ativa: %t\n", nome, idade, ativa)

	// --- Declaração curta := (infere o tipo) ---
	cidade := "São Paulo" // string
	temperatura := 28.5   // float64
	populacao := 12000000 // int

	fmt.Printf("Cidade: %s, Temp: %.1f°C, Pop: %d\n", cidade, temperatura, populacao)

	// --- Múltiplas variáveis de uma vez ---
	var (
		x int     = 10
		y float64 = 20.5
		z string  = "hello"
	)
	fmt.Printf("x=%d, y=%.1f, z=%s\n", x, y, z)

	// Ou com :=
	a, b, c := 1, 2.0, "tres"
	fmt.Printf("a=%d, b=%.1f, c=%s\n", a, b, c)

	// ==========================================================
	fmt.Println("\n=== ZERO VALUES ===")

	// Em Go, variáveis NUNCA são "undefined" ou "null"
	// Cada tipo tem um "zero value" padrão
	var intZero int
	var floatZero float64
	var stringZero string
	var boolZero bool

	fmt.Printf("int: %d\n", intZero)         // 0
	fmt.Printf("float64: %f\n", floatZero)   // 0.000000
	fmt.Printf("string: '%s'\n", stringZero) // "" (vazio)
	fmt.Printf("bool: %t\n", boolZero)       // false

	// ==========================================================
	fmt.Println("\n=== TIPOS E TAMANHOS ===")

	// int e uint variam de tamanho conforme a plataforma (32 ou 64 bits)
	var i8 int8 = 127     // -128 a 127
	var i16 int16 = 32767 // -32768 a 32767
	var i32 int32 = 100000
	var i64 int64 = 9999999999

	fmt.Printf("int8: %d, int16: %d, int32: %d, int64: %d\n", i8, i16, i32, i64)

	// byte é alias de uint8 (0 a 255)
	var b1 byte = 'A'
	fmt.Printf("byte: %d (char: %c)\n", b1, b1) // 65 (char: A)

	// rune é alias de int32 — representa um caractere Unicode
	var r1 rune = '世'
	fmt.Printf("rune: %d (char: %c)\n", r1, r1) // 19990 (char: 世)

	// ==========================================================
	fmt.Println("\n=== CONSTANTES E IOTA ===")

	// Constantes não podem ser alteradas depois de definidas
	const Pi = 3.14159
	const AppNome = "MeuApp"
	fmt.Printf("Pi: %f, App: %s\n", Pi, AppNome)

	// iota → auto-incrementa dentro de um bloco const
	const (
		Domingo = iota // 0
		Segunda        // 1
		Terca          // 2
		Quarta         // 3
		Quinta         // 4
		Sexta          // 5
		Sabado         // 6
	)
	fmt.Printf("Domingo=%d, Sexta=%d, Sabado=%d\n", Domingo, Sexta, Sabado)

	// iota com expressão (útil para bit flags)
	const (
		Leitura  = 1 << iota // 1  (001)
		Escrita              // 2  (010)
		Execucao             // 4  (100)
	)
	fmt.Printf("Permissões: Leitura=%d, Escrita=%d, Execução=%d\n", Leitura, Escrita, Execucao)

	// Combinando permissões com OR bitwise
	permissao := Leitura | Escrita // 3 (011)
	fmt.Printf("Leitura+Escrita = %d (%b em binário)\n", permissao, permissao)

	// ==========================================================
	fmt.Println("\n=== CONVERSÃO DE TIPOS ===")

	// Go NÃO faz conversão implícita. Você precisa ser explícito.
	inteiro := 42
	decimal := float64(inteiro) // int → float64
	fmt.Printf("int %d → float64 %.2f\n", inteiro, decimal)

	nota := 9.7
	notaInteira := int(nota) // float64 → int (trunca, NÃO arredonda!)
	fmt.Printf("float64 %.1f → int %d (truncado!)\n", nota, notaInteira)

	// int → string (cuidado! converte o code point, NÃO o número)
	// Para converter número para string, use fmt.Sprintf
	numero := 65
	letraErrada := string(rune(numero))     // "A" (code point 65 = 'A')
	letraCerta := fmt.Sprintf("%d", numero) // "65"
	fmt.Printf("string(65) = '%s' (code point!)\n", letraErrada)
	fmt.Printf("Sprintf = '%s' (correto!)\n", letraCerta)

	// ==========================================================
	fmt.Println("\n=== VERBOS DE FORMATAÇÃO (fmt.Printf) ===")

	// Os verbos mais usados:
	valor := 42
	fmt.Printf("%%d  = %d (decimal)\n", valor)
	fmt.Printf("%%b  = %b (binário)\n", valor)
	fmt.Printf("%%o  = %o (octal)\n", valor)
	fmt.Printf("%%x  = %x (hexadecimal)\n", valor)
	fmt.Printf("%%f  = %f (float padrão)\n", float64(valor))
	fmt.Printf("%%.2f = %.2f (2 casas decimais)\n", float64(valor))
	fmt.Printf("%%s  = %s (string)\n", "hello")
	fmt.Printf("%%q  = %q (string com aspas)\n", "hello")
	fmt.Printf("%%t  = %t (boolean)\n", true)
	fmt.Printf("%%T  = %T (tipo da variável)\n", valor)
	fmt.Printf("%%v  = %v (valor padrão — serve pra tudo)\n", valor)
	fmt.Printf("%%+v = %+v (struct com nomes dos campos)\n", struct{ Nome string }{"Go"})
}
