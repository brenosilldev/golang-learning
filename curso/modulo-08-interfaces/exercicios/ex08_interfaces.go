package main

// ============================================================================
// EXERCÍCIO 8 — Interfaces
// ============================================================================
//
// Exercício 8.1 — Sistema de Pagamento
// Crie uma interface `Pagamento` com:
//   - Processar(valor float64) error
//   - Nome() string
//
// Implemente 3 tipos:
//   - CartaoCredito (bandeira string, limite float64)
//   - Pix (chave string)
//   - Boleto (codigoBarras string)
//
// Cada um deve ter regras próprias:
//   - CartaoCredito: falha se valor > limite
//   - Pix: sempre sucesso
//   - Boleto: falha se valor < 10 (valor mínimo)
//
// Crie uma função `realizarPagamento(p Pagamento, valor float64)` que
// recebe qualquer Pagamento e tenta processar.
//
// Exercício 8.2 — fmt.Stringer
// Crie uma struct `Produto` com Nome, Preco e Categoria.
// Implemente `String() string` para que fmt.Println imprima algo como:
//   📦 Notebook (Eletrônicos) - R$ 4.599,90
//
// Exercício 8.3 — Type Switch Prático
// Crie uma função `descrever(v any) string` que retorna uma descrição
// diferente para cada tipo: int, string, bool, float64, []int, e default.
// Teste com vários valores diferentes.
//
// Exercício 8.4 — Ordenação Customizada (sort.Interface)
// Crie um tipo `Funcionarios` ([]Funcionario) que implementa sort.Interface
// para ordenar por salário (do menor para o maior).
// Struct: Funcionario{Nome string, Salario float64}
// Dica: implement Len(), Less(i,j int) bool, Swap(i,j int)
//
// ============================================================================

func main() {
	// TODO: Exercício 8.1 — Sistema de Pagamento

	// TODO: Exercício 8.2 — fmt.Stringer

	// TODO: Exercício 8.3 — Type Switch Prático

	// TODO: Exercício 8.4 — Ordenação Customizada
}
