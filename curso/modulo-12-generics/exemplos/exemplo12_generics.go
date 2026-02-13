package main

import (
	"cmp"
	"fmt"
)

// ============================================================================
// MÓDULO 12 — Generics (Go 1.18+)
// ============================================================================
//
// Conceitos:
//   - Type parameters [T any]
//   - Constraints (restrições de tipo)
//   - Funções genéricas
//   - Structs genéricas
//   - Interfaces como constraints
//   - cmp.Ordered (Go 1.21+)
//
// Rode com: go run exemplo12_generics.go
// ============================================================================

// --- Constraint customizada ---
type Numero interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 |
		~float32 | ~float64
}

// ~ (til) permite tipos derivados: type Celsius float64 → aceita Celsius

// --- Função genérica simples ---
func Soma[T Numero](a, b T) T {
	return a + b
}

// --- Genéric com any ---
func Contem[T comparable](slice []T, alvo T) bool {
	for _, v := range slice {
		if v == alvo {
			return true
		}
	}
	return false
}

// --- Map genérico ---
func Map[T any, U any](slice []T, fn func(T) U) []U {
	resultado := make([]U, len(slice))
	for i, v := range slice {
		resultado[i] = fn(v)
	}
	return resultado
}

// --- Filter genérico ---
func Filter[T any](slice []T, fn func(T) bool) []T {
	resultado := make([]T, 0)
	for _, v := range slice {
		if fn(v) {
			resultado = append(resultado, v)
		}
	}
	return resultado
}

// --- Reduce genérico ---
func Reduce[T any, U any](slice []T, inicial U, fn func(U, T) U) U {
	acumulador := inicial
	for _, v := range slice {
		acumulador = fn(acumulador, v)
	}
	return acumulador
}

// --- Min/Max genéricos com cmp.Ordered ---
func Min[T cmp.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T cmp.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// --- Struct genérica: Stack ---
type Stack[T any] struct {
	items []T
}

func (s *Stack[T]) Push(item T) {
	s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	item := s.items[len(s.items)-1]
	s.items = s.items[:len(s.items)-1]
	return item, true
}

func (s *Stack[T]) Peek() (T, bool) {
	if len(s.items) == 0 {
		var zero T
		return zero, false
	}
	return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int {
	return len(s.items)
}

// --- Struct genérica: Pair (tupla) ---
type Pair[T, U any] struct {
	First  T
	Second U
}

func NewPair[T, U any](first T, second U) Pair[T, U] {
	return Pair[T, U]{First: first, Second: second}
}

// --- Struct genérica: Result (como Result do Rust) ---
type Result[T any] struct {
	Value T
	Err   error
}

func Ok[T any](valor T) Result[T] {
	return Result[T]{Value: valor, Err: nil}
}

func Errorf[T any](err error) Result[T] {
	return Result[T]{Err: err}
}

func (r Result[T]) IsOk() bool {
	return r.Err == nil
}

func main() {
	// ==========================================================
	fmt.Println("=== FUNÇÕES GENÉRICAS ===")

	// Soma funciona com qualquer tipo numérico
	fmt.Println("Soma int:", Soma(10, 20))
	fmt.Println("Soma float:", Soma(3.14, 2.86))

	// Contem funciona com qualquer tipo comparable
	nums := []int{1, 2, 3, 4, 5}
	fmt.Println("Contém 3?", Contem(nums, 3))
	fmt.Println("Contém 9?", Contem(nums, 9))

	palavras := []string{"go", "rust", "python"}
	fmt.Println("Contém 'go'?", Contem(palavras, "go"))

	// ==========================================================
	fmt.Println("\n=== MAP / FILTER / REDUCE GENÉRICOS ===")

	numeros := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	// Map: dobrar
	dobrados := Map(numeros, func(n int) int { return n * 2 })
	fmt.Println("Dobrados:", dobrados)

	// Map: int → string
	textos := Map(numeros, func(n int) string {
		return fmt.Sprintf("#%d", n)
	})
	fmt.Println("Textos:", textos)

	// Filter: apenas pares
	pares := Filter(numeros, func(n int) bool { return n%2 == 0 })
	fmt.Println("Pares:", pares)

	// Reduce: soma total
	soma := Reduce(numeros, 0, func(acc, n int) int { return acc + n })
	fmt.Println("Soma:", soma)

	// Encadeando: soma dos quadrados dos pares
	resultado := Reduce(
		Map(
			Filter(numeros, func(n int) bool { return n%2 == 0 }),
			func(n int) int { return n * n },
		),
		0,
		func(acc, n int) int { return acc + n },
	)
	fmt.Println("Soma quadrados dos pares:", resultado)

	// ==========================================================
	fmt.Println("\n=== MIN / MAX ===")

	fmt.Println("Min(3, 7):", Min(3, 7))
	fmt.Println("Max(3, 7):", Max(3, 7))
	fmt.Println("Min strings:", Min("abc", "xyz"))

	// ==========================================================
	fmt.Println("\n=== STACK GENÉRICA ===")

	// Stack de int
	intStack := &Stack[int]{}
	intStack.Push(10)
	intStack.Push(20)
	intStack.Push(30)
	fmt.Printf("Stack size: %d\n", intStack.Len())

	if val, ok := intStack.Pop(); ok {
		fmt.Println("Pop:", val)
	}
	if val, ok := intStack.Peek(); ok {
		fmt.Println("Peek:", val)
	}

	// Stack de string — mesmo tipo!
	strStack := &Stack[string]{}
	strStack.Push("hello")
	strStack.Push("world")
	if val, ok := strStack.Pop(); ok {
		fmt.Println("String Pop:", val)
	}

	// ==========================================================
	fmt.Println("\n=== PAIR (TUPLA) ===")

	par := NewPair("Go", 2009)
	fmt.Printf("Linguagem: %s, Ano: %d\n", par.First, par.Second)

	// Pair de tipos diferentes
	par2 := NewPair(3.14, true)
	fmt.Printf("Pi: %.2f, Legal: %t\n", par2.First, par2.Second)

	// ==========================================================
	fmt.Println("\n=== RESULT TYPE ===")

	r1 := Ok(42)
	r2 := Errorf[int](fmt.Errorf("algo deu errado"))

	fmt.Printf("r1: value=%d, ok=%t\n", r1.Value, r1.IsOk())
	fmt.Printf("r2: err=%v, ok=%t\n", r2.Err, r2.IsOk())
}
