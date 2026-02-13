package main

// ============================================================================
// EXERCÍCIO 9 — Tratamento de Erros
// ============================================================================
//
// Exercício 9.1 — Validador de Cadastro
// Crie uma função `validarCadastro(nome, email, senha string) error` que:
//   - Retorna erro se nome estiver vazio
//   - Retorna erro se email não contém "@"
//   - Retorna erro se senha tem menos de 8 caracteres
// Use custom error type com campo indicando qual validação falhou.
//
// Exercício 9.2 — Calculadora Segura
// Crie funções para as 4 operações (+, -, *, /).
// Cada uma retorna (float64, error).
// Crie sentinel errors: ErrDivisaoPorZero, ErrOverflow.
// Processe uma lista de operações e trate cada erro individualmente.
//
// Exercício 9.3 — Parser de Configuração
// Dada uma string no formato "chave=valor" (uma por elemento do slice):
//   config := []string{"host=localhost", "port=8080", "invalida", "timeout=abc"}
// Crie uma função que parse cada linha e retorne map[string]string.
// Trate erros:
//   - Linha sem "=" → erro de formato
//   - Valor que deveria ser número mas não é → erro de tipo
// Use fmt.Errorf com %w para encadear erros.
//
// Exercício 9.4 — Recover de Panic
// Crie uma função que recebe um slice de funções ([]func()) e executa
// cada uma. Se alguma der panic, capture com recover, log o erro,
// e continue executando as próximas.
//
// ============================================================================

func main() {
	// TODO: Exercício 9.1 — Validador de Cadastro

	// TODO: Exercício 9.2 — Calculadora Segura

	// TODO: Exercício 9.3 — Parser de Configuração

	// TODO: Exercício 9.4 — Recover de Panic
}
