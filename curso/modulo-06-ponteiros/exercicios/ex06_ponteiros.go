package main

// ============================================================================
// EXERCÍCIO 6 — Ponteiros
// ============================================================================
//
// Exercício 6.1 — Trocar Valores
// Crie uma função `trocar(a, b *string)` que troca os valores de duas strings.
// Teste com: nome1 := "Alice", nome2 := "Bob"
// Após trocar: nome1 = "Bob", nome2 = "Alice"
//
// Exercício 6.2 — Modificar Slice via Ponteiro
// Crie uma função `adicionarItem(lista *[]string, item string)`
// que adiciona um item a um slice via ponteiro.
// Por que isso é necessário? Reflita: append pode realocar o slice.
//
// Exercício 6.3 — Conta Bancária
// Crie uma struct ContaBancaria com campos Titular (string) e Saldo (float64).
// Implemente:
//   - depositar(conta *ContaBancaria, valor float64)
//   - sacar(conta *ContaBancaria, valor float64) error  (retorna erro se saldo insuficiente)
//   - transferir(origem, destino *ContaBancaria, valor float64) error
// Teste todas as operações e verifique que os saldos são alterados corretamente.
//
// Exercício 6.4 — Linked List Simples
// Crie uma struct Node com Valor (int) e Proximo (*Node).
// Implemente:
//   - inserirFim(head **Node, valor int) — insere no final
//   - imprimir(head *Node) — imprime toda a lista
// Teste criando uma lista: 1 → 2 → 3 → 4 → nil
//
// ============================================================================

func main() {
	// TODO: Exercício 6.1 — Trocar Valores

	// TODO: Exercício 6.2 — Modificar Slice via Ponteiro

	// TODO: Exercício 6.3 — Conta Bancária

	// TODO: Exercício 6.4 — Linked List Simples
}
