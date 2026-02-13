package main

import (
	"fmt"
	"math"
	"strings"
)

// ============================================================================
// MÓDULO 08 — Interfaces
// ============================================================================
//
// Conceitos:
//   - Interfaces implícitas (sem "implements")
//   - Polimorfismo
//   - interface{} / any (aceita qualquer tipo)
//   - Type assertion e type switch
//   - Interfaces da stdlib (Stringer, error, io.Reader/Writer)
//   - Composição de interfaces
//
// Rode com: go run exemplo08_interfaces.go
// ============================================================================

// --- Interface: define um contrato de comportamento ---
type Forma interface {
	Area() float64
	Perimetro() float64
}

// --- Implementações da interface Forma ---

type Retangulo struct {
	Largura, Altura float64
}

func (r Retangulo) Area() float64 {
	return r.Largura * r.Altura
}

func (r Retangulo) Perimetro() float64 {
	return 2 * (r.Largura + r.Altura)
}

type Circulo struct {
	Raio float64
}

func (c Circulo) Area() float64 {
	return math.Pi * c.Raio * c.Raio
}

func (c Circulo) Perimetro() float64 {
	return 2 * math.Pi * c.Raio
}

type Triangulo struct {
	Base, Altura, LadoA, LadoB, LadoC float64
}

func (t Triangulo) Area() float64 {
	return t.Base * t.Altura / 2
}

func (t Triangulo) Perimetro() float64 {
	return t.LadoA + t.LadoB + t.LadoC
}

// --- Função polimórfica: aceita qualquer Forma ---
func imprimirForma(f Forma) {
	fmt.Printf("  Área: %.2f | Perímetro: %.2f\n", f.Area(), f.Perimetro())
}

// --- Interface fmt.Stringer ---
// Implementar String() faz seu tipo se "apresentar" bonito com fmt.Println
type Pessoa struct {
	Nome  string
	Idade int
}

func (p Pessoa) String() string {
	return fmt.Sprintf("👤 %s (%d anos)", p.Nome, p.Idade)
}

// --- Interface error ---
// Qualquer tipo com método Error() string é um error
type ErroValidacao struct {
	Campo    string
	Mensagem string
}

func (e ErroValidacao) Error() string {
	return fmt.Sprintf("validação falhou no campo '%s': %s", e.Campo, e.Mensagem)
}

// --- Composição de interfaces ---
type Leitor interface {
	Ler() string
}

type Escritor interface {
	Escrever(dados string)
}

// Interface composta
type LeitorEscritor interface {
	Leitor
	Escritor
}

// Implementação
type Arquivo struct {
	conteudo string
}

func (a *Arquivo) Ler() string {
	return a.conteudo
}

func (a *Arquivo) Escrever(dados string) {
	a.conteudo += dados
}

func main() {
	// ==========================================================
	fmt.Println("=== POLIMORFISMO ===")

	// Todas implementam Forma — podemos tratá-las uniformemente
	formas := []Forma{
		Retangulo{Largura: 5, Altura: 3},
		Circulo{Raio: 4},
		Triangulo{Base: 6, Altura: 4, LadoA: 5, LadoB: 5, LadoC: 6},
	}

	for _, forma := range formas {
		fmt.Printf("%T:\n", forma)
		imprimirForma(forma)
	}

	// ==========================================================
	fmt.Println("\n=== INTERFACE VAZIA (any) ===")

	// any (= interface{}) aceita QUALQUER tipo
	var caixa any

	caixa = 42
	fmt.Printf("int: %v (tipo: %T)\n", caixa, caixa)

	caixa = "hello"
	fmt.Printf("string: %v (tipo: %T)\n", caixa, caixa)

	caixa = true
	fmt.Printf("bool: %v (tipo: %T)\n", caixa, caixa)

	caixa = Retangulo{5, 3}
	fmt.Printf("struct: %v (tipo: %T)\n", caixa, caixa)

	// Slice de any — pode misturar tipos
	mistura := []any{1, "dois", 3.0, true, Pessoa{"Go", 15}}
	fmt.Println("\nSlice misto:", mistura)

	// ==========================================================
	fmt.Println("\n=== TYPE ASSERTION ===")

	var valor any = "Hello, Go!"

	// Assertion UNSAFE — panic se o tipo estiver errado!
	str := valor.(string)
	fmt.Println("Assertion:", str)

	// Assertion SAFE — com ok check (SEMPRE use esta!)
	str2, ok := valor.(string)
	if ok {
		fmt.Println("Safe assertion:", str2)
	}

	num, ok := valor.(int)
	if !ok {
		fmt.Printf("Não é int (valor zero: %d)\n", num)
	}

	// ==========================================================
	fmt.Println("\n=== TYPE SWITCH ===")

	valores := []any{42, "hello", 3.14, true, Pessoa{"Alice", 30}, nil}

	for _, v := range valores {
		switch val := v.(type) {
		case int:
			fmt.Printf("int: %d (dobro: %d)\n", val, val*2)
		case string:
			fmt.Printf("string: '%s' (upper: %s)\n", val, strings.ToUpper(val))
		case float64:
			fmt.Printf("float64: %.2f\n", val)
		case bool:
			fmt.Printf("bool: %t\n", val)
		case Pessoa:
			fmt.Printf("Pessoa: %s\n", val) // usa Stringer!
		case nil:
			fmt.Println("nil!")
		default:
			fmt.Printf("tipo desconhecido: %T\n", val)
		}
	}

	// ==========================================================
	fmt.Println("\n=== fmt.Stringer ===")

	// Quando Pessoa implementa String(), fmt.Println usa automaticamente
	p := Pessoa{Nome: "Breno", Idade: 25}
	fmt.Println(p) // Chama p.String() automaticamente!

	// ==========================================================
	fmt.Println("\n=== CUSTOM error ===")

	err := validarEmail("")
	if err != nil {
		fmt.Println("Erro:", err)
	}

	err = validarEmail("invalido")
	if err != nil {
		fmt.Println("Erro:", err)
	}

	// ==========================================================
	fmt.Println("\n=== COMPOSIÇÃO DE INTERFACES ===")

	arq := &Arquivo{}
	arq.Escrever("Primeira linha\n")
	arq.Escrever("Segunda linha\n")

	fmt.Println("Conteúdo:")
	fmt.Println(arq.Ler())

	// Funciona como LeitorEscritor
	var rw LeitorEscritor = arq
	rw.Escrever("Terceira linha\n")
	fmt.Println("Após mais escrita:")
	fmt.Println(rw.Ler())

	// ==========================================================
	fmt.Println("=== VERIFICAR SE IMPLEMENTA INTERFACE ===")

	// Em tempo de compilação, verificar se um tipo implementa uma interface
	var _ Forma = Retangulo{} // Compila ✅
	var _ Forma = Circulo{}   // Compila ✅
	// var _ Forma = Pessoa{}  // NÃO compilaria ❌

	fmt.Println("Todas as implementações verificadas ✅")
}

func validarEmail(email string) error {
	if email == "" {
		return ErroValidacao{Campo: "email", Mensagem: "não pode ser vazio"}
	}
	if !strings.Contains(email, "@") {
		return ErroValidacao{Campo: "email", Mensagem: "formato inválido"}
	}
	return nil
}
