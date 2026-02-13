package main

import (
	"errors"
	"fmt"
	"strings"
)

// ============================================================================
// MÓDULO 13 — Código para ser testado
// ============================================================================
//
// Este arquivo contém funções que serão testadas no arquivo
// exemplo13_codigo_test.go. Demonstra como organizar código testável.
//
// Para rodar os testes:
//   go test -v ./curso/modulo-13-testes/exemplos/
//   go test -cover ./curso/modulo-13-testes/exemplos/
//   go test -bench=. ./curso/modulo-13-testes/exemplos/
// ============================================================================

// --- Funções simples ---

func Soma(a, b int) int {
	return a + b
}

func Dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, errors.New("divisão por zero")
	}
	return a / b, nil
}

func Fatorial(n int) (int, error) {
	if n < 0 {
		return 0, errors.New("número negativo")
	}
	if n <= 1 {
		return 1, nil
	}
	resultado := 1
	for i := 2; i <= n; i++ {
		resultado *= i
	}
	return resultado, nil
}

// --- Funções de string ---

func Reverter(s string) string {
	runes := []rune(s)
	for i, j := 0, len(runes)-1; i < j; i, j = i+1, j-1 {
		runes[i], runes[j] = runes[j], runes[i]
	}
	return string(runes)
}

func EhPalindromo(s string) bool {
	s = strings.ToLower(strings.TrimSpace(s))
	return s == Reverter(s)
}

func ContarPalavras(s string) int {
	if strings.TrimSpace(s) == "" {
		return 0
	}
	return len(strings.Fields(s))
}

// --- Struct com métodos ---

type ContaBancaria struct {
	Titular string
	Saldo   float64
}

func NovaConta(titular string, saldoInicial float64) (*ContaBancaria, error) {
	if titular == "" {
		return nil, errors.New("titular obrigatório")
	}
	if saldoInicial < 0 {
		return nil, errors.New("saldo inicial negativo")
	}
	return &ContaBancaria{Titular: titular, Saldo: saldoInicial}, nil
}

func (c *ContaBancaria) Depositar(valor float64) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}
	c.Saldo += valor
	return nil
}

func (c *ContaBancaria) Sacar(valor float64) error {
	if valor <= 0 {
		return errors.New("valor deve ser positivo")
	}
	if valor > c.Saldo {
		return fmt.Errorf("saldo insuficiente: tem %.2f, pediu %.2f", c.Saldo, valor)
	}
	c.Saldo -= valor
	return nil
}

func main() {
	fmt.Println("Este arquivo contém código para ser testado.")
	fmt.Println("Rode: go test -v ./curso/modulo-13-testes/exemplos/")
}
