package main

import (
	"testing"
)

// ============================================================================
// MÓDULO 13 — Testes em Go
// ============================================================================
//
// Conceitos demonstrados:
//   - Test básico com t.Errorf
//   - Table-driven tests (pattern principal!)
//   - Subtestes com t.Run
//   - t.Fatal vs t.Error
//   - t.Skip (pular teste condicionalmente)
//   - Benchmarks com testing.B
//   - Helpers de teste
//
// Rode com:
//   go test -v ./curso/modulo-13-testes/exemplos/
//   go test -bench=. ./curso/modulo-13-testes/exemplos/
//   go test -cover ./curso/modulo-13-testes/exemplos/
// ============================================================================

// --- Teste simples ---
func TestSoma(t *testing.T) {
	resultado := Soma(2, 3)
	if resultado != 5 {
		t.Errorf("Soma(2, 3) = %d; esperado 5", resultado)
	}
}

// --- TABLE-DRIVEN TESTS (o pattern mais importante do Go!) ---
func TestSomaTableDriven(t *testing.T) {
	tests := []struct {
		name     string
		a, b     int
		esperado int
	}{
		{"positivos", 2, 3, 5},
		{"zeros", 0, 0, 0},
		{"negativos", -1, -2, -3},
		{"misto", -5, 10, 5},
		{"grande", 1000000, 2000000, 3000000},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado := Soma(tt.a, tt.b)
			if resultado != tt.esperado {
				t.Errorf("Soma(%d, %d) = %d; esperado %d",
					tt.a, tt.b, resultado, tt.esperado)
			}
		})
	}
}

// --- Testando funções que retornam erro ---
func TestDividir(t *testing.T) {
	tests := []struct {
		name     string
		a, b     float64
		esperado float64
		wantErr  bool
	}{
		{"normal", 10, 2, 5, false},
		{"decimal", 7, 3, 2.3333333333333335, false},
		{"divisão por zero", 10, 0, 0, true},
		{"zero dividido", 0, 5, 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resultado, err := Dividir(tt.a, tt.b)

			if tt.wantErr {
				if err == nil {
					t.Error("esperava erro, mas não recebeu")
				}
				return
			}

			if err != nil {
				t.Fatalf("erro inesperado: %v", err)
			}

			if resultado != tt.esperado {
				t.Errorf("Dividir(%.1f, %.1f) = %.4f; esperado %.4f",
					tt.a, tt.b, resultado, tt.esperado)
			}
		})
	}
}

// --- Testando strings ---
func TestReverter(t *testing.T) {
	tests := []struct {
		input, esperado string
	}{
		{"hello", "olleh"},
		{"", ""},
		{"a", "a"},
		{"Go é 🔥", "🔥 é oG"},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			resultado := Reverter(tt.input)
			if resultado != tt.esperado {
				t.Errorf("Reverter(%q) = %q; esperado %q",
					tt.input, resultado, tt.esperado)
			}
		})
	}
}

func TestEhPalindromo(t *testing.T) {
	tests := []struct {
		input    string
		esperado bool
	}{
		{"ana", true},
		{"arara", true},
		{"hello", false},
		{"", true},
		{"A", true},
		{"Aba", true}, // case insensitive
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			resultado := EhPalindromo(tt.input)
			if resultado != tt.esperado {
				t.Errorf("EhPalindromo(%q) = %t; esperado %t",
					tt.input, resultado, tt.esperado)
			}
		})
	}
}

// --- Testando struct com métodos ---
func TestContaBancaria(t *testing.T) {
	t.Run("criar conta válida", func(t *testing.T) {
		conta, err := NovaConta("Alice", 1000)
		if err != nil {
			t.Fatalf("erro inesperado: %v", err)
		}
		if conta.Titular != "Alice" {
			t.Errorf("titular = %s; esperado Alice", conta.Titular)
		}
		if conta.Saldo != 1000 {
			t.Errorf("saldo = %.2f; esperado 1000", conta.Saldo)
		}
	})

	t.Run("criar conta sem titular", func(t *testing.T) {
		_, err := NovaConta("", 1000)
		if err == nil {
			t.Error("esperava erro ao criar conta sem titular")
		}
	})

	t.Run("depositar e sacar", func(t *testing.T) {
		conta, _ := NovaConta("Bob", 500)

		if err := conta.Depositar(200); err != nil {
			t.Fatalf("erro no depósito: %v", err)
		}
		if conta.Saldo != 700 {
			t.Errorf("saldo após depósito = %.2f; esperado 700", conta.Saldo)
		}

		if err := conta.Sacar(300); err != nil {
			t.Fatalf("erro no saque: %v", err)
		}
		if conta.Saldo != 400 {
			t.Errorf("saldo após saque = %.2f; esperado 400", conta.Saldo)
		}
	})

	t.Run("saque maior que saldo", func(t *testing.T) {
		conta, _ := NovaConta("Carol", 100)
		err := conta.Sacar(500)
		if err == nil {
			t.Error("esperava erro de saldo insuficiente")
		}
	})
}

// --- Helper de teste ---
func assertError(t *testing.T, got error) {
	t.Helper() // Marca como helper — erro mostra quem chamou
	if got == nil {
		t.Error("esperava erro, mas recebeu nil")
	}
}

func assertNoError(t *testing.T, got error) {
	t.Helper()
	if got != nil {
		t.Errorf("não esperava erro, mas recebeu: %v", got)
	}
}

// --- Benchmark ---
func BenchmarkSoma(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Soma(100, 200)
	}
}

func BenchmarkFatorial(b *testing.B) {
	for i := 0; i < b.N; i++ {
		Fatorial(20)
	}
}

func BenchmarkReverter(b *testing.B) {
	texto := "Esta é uma string para benchmark de reversão"
	for i := 0; i < b.N; i++ {
		Reverter(texto)
	}
}
