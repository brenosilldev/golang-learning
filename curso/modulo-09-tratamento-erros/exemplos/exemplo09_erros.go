package main

import (
	"errors"
	"fmt"
	"strconv"
)

// ============================================================================
// MÓDULO 09 — Tratamento de Erros
// ============================================================================
//
// Conceitos:
//   - Pattern "if err != nil"
//   - errors.New e fmt.Errorf
//   - Error wrapping com %w
//   - Custom error types
//   - Sentinel errors
//   - errors.Is e errors.As
//   - panic e recover (e quando NÃO usar)
//
// Rode com: go run exemplo09_erros.go
// ============================================================================

// --- Sentinel errors (erros como constantes) ---
var (
	ErrNaoEncontrado  = errors.New("recurso não encontrado")
	ErrSemPermissao   = errors.New("sem permissão")
	ErrDivisaoPorZero = errors.New("divisão por zero")
)

// --- Custom error type ---
type ErroValidacao struct {
	Campo    string
	Mensagem string
}

func (e *ErroValidacao) Error() string {
	return fmt.Sprintf("erro de validação: campo '%s' — %s", e.Campo, e.Mensagem)
}

// --- Custom error com Unwrap (para errors.Is funcionar) ---
type ErroOperacao struct {
	Operacao string
	Err      error // erro original
}

func (e *ErroOperacao) Error() string {
	return fmt.Sprintf("falha na operação '%s': %v", e.Operacao, e.Err)
}

func (e *ErroOperacao) Unwrap() error {
	return e.Err // Permite errors.Is e errors.As encontrar o erro original
}

func main() {
	// ==========================================================
	fmt.Println("=== PATTERN BÁSICO: if err != nil ===")

	resultado, err := dividir(10, 3)
	if err != nil {
		fmt.Println("Erro:", err)
	} else {
		fmt.Printf("10 / 3 = %.2f\n", resultado)
	}

	_, err = dividir(10, 0)
	if err != nil {
		fmt.Println("Erro:", err)
	}

	// ==========================================================
	fmt.Println("\n=== errors.New e fmt.Errorf ===")

	err1 := errors.New("algo deu errado")
	fmt.Println(err1)

	// fmt.Errorf — erro com formatação
	usuario := "admin"
	err2 := fmt.Errorf("usuário '%s' não encontrado", usuario)
	fmt.Println(err2)

	// ==========================================================
	fmt.Println("\n=== ERROR WRAPPING com %w ===")

	// %w "embrulha" um erro dentro de outro — preserva a cadeia
	err3 := buscarUsuario(-1)
	if err3 != nil {
		fmt.Println("Erro:", err3)

		// errors.Is verifica se EM ALGUM NÍVEL da cadeia tem o erro
		if errors.Is(err3, ErrNaoEncontrado) {
			fmt.Println("→ É um erro de não encontrado!")
		}
	}

	// ==========================================================
	fmt.Println("\n=== SENTINEL ERRORS ===")

	testCases := []int{1, -1, 0}
	for _, id := range testCases {
		err := fazerOperacao(id)
		switch {
		case err == nil:
			fmt.Printf("ID %d: sucesso\n", id)
		case errors.Is(err, ErrNaoEncontrado):
			fmt.Printf("ID %d: não encontrado\n", id)
		case errors.Is(err, ErrSemPermissao):
			fmt.Printf("ID %d: sem permissão\n", id)
		default:
			fmt.Printf("ID %d: erro desconhecido: %v\n", id, err)
		}
	}

	// ==========================================================
	fmt.Println("\n=== CUSTOM ERROR + errors.As ===")

	err4 := validarFormulario("", "email-invalido")
	if err4 != nil {
		fmt.Println("Erro:", err4)

		// errors.As extrai um tipo específico de erro da cadeia
		var errVal *ErroValidacao
		if errors.As(err4, &errVal) {
			fmt.Printf("→ Campo com problema: '%s'\n", errVal.Campo)
			fmt.Printf("→ Mensagem: %s\n", errVal.Mensagem)
		}
	}

	// ==========================================================
	fmt.Println("\n=== MÚLTIPLOS ERROS em sequência ===")

	// Pattern: tentar converter lista de strings para int
	valores := []string{"10", "abc", "30", "xyz", "50"}
	fmt.Println("Convertendo:", valores)

	for _, v := range valores {
		num, err := strconv.Atoi(v)
		if err != nil {
			fmt.Printf("  '%s' → ERRO: %v\n", v, err)
			continue // pular para o próximo — NÃO parar
		}
		fmt.Printf("  '%s' → %d ✓\n", v, num)
	}

	// ==========================================================
	fmt.Println("\n=== PANIC e RECOVER ===")

	// panic = erro irrecuperável (NÃO use para lógica de negócio!)
	// recover = captura panic dentro de defer

	fmt.Println("Antes de funcaoPerigosa()")
	funcaoPerigosa()
	fmt.Println("Depois de funcaoPerigosa() — programa continua!")

	// Quando usar panic:
	// ✅ Bug impossível no código (assertion failed)
	// ✅ Inicialização falhou (DB não conectou)
	// ❌ NUNCA para erros esperados (arquivo não encontrado, input inválido)
	// ❌ NUNCA para fluxo de controle

	// ==========================================================
	fmt.Println("\n=== PATTERN: Tratar erro uma vez ===")

	// ❌ ERRADO — loga E retorna (trata duas vezes)
	// if err != nil {
	//     log.Printf("erro: %v", err)  // log AQUI
	//     return err                     // e quem chamou loga de novo!
	// }

	// ✅ CERTO — trate OU retorne, nunca ambos
	// Nível baixo: adicionar contexto e retornar
	// Nível alto: tratar (log, resposta HTTP, etc.)
}

func dividir(a, b float64) (float64, error) {
	if b == 0 {
		return 0, ErrDivisaoPorZero
	}
	return a / b, nil
}

func buscarUsuario(id int) error {
	if id < 0 {
		// %w embrulha o erro — errors.Is consegue encontrar ErrNaoEncontrado
		return fmt.Errorf("buscarUsuario(id=%d): %w", id, ErrNaoEncontrado)
	}
	return nil
}

func fazerOperacao(id int) error {
	if id < 0 {
		return &ErroOperacao{Operacao: "buscar", Err: ErrNaoEncontrado}
	}
	if id == 0 {
		return &ErroOperacao{Operacao: "acessar", Err: ErrSemPermissao}
	}
	return nil
}

func validarFormulario(nome, email string) error {
	if nome == "" {
		return &ErroValidacao{Campo: "nome", Mensagem: "não pode ser vazio"}
	}
	if len(email) < 5 {
		return &ErroValidacao{Campo: "email", Mensagem: "muito curto"}
	}
	return nil
}

func funcaoPerigosa() {
	// defer + recover captura o panic
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("Recover capturou: %v\n", r)
		}
	}()

	fmt.Println("→ Vai dar panic...")
	panic("algo terrível aconteceu!")
	// Código abaixo NUNCA executa
}
