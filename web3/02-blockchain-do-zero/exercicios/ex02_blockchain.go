package main

// ============================================================================
// EXERCÍCIOS — Blockchain do Zero
// ============================================================================
//
// Exercício 2.1 — Dificuldade Dinâmica
// Modifique a blockchain para que a dificuldade se ajuste automaticamente:
//   - Se o último bloco foi minerado em < 2 segundos → aumenta dificuldade
//   - Se demorou > 10 segundos → diminui dificuldade
// Mine 10 blocos e observe a dificuldade se ajustando.
//
// Exercício 2.2 — Validação de Saldo
// Antes de aceitar uma transação, verifique se o remetente tem saldo suficiente.
//   - Calcule o saldo percorrendo TODOS os blocos
//   - Se saldo < valor → rejeitar a transação
//   - Transações do "SISTEMA" e "REDE" são especiais (não precisam de saldo)
// Teste: tente fazer Alice enviar mais do que tem.
//
// Exercício 2.3 — Pool de Transações com Prioridade
// Implemente uma pool onde transações com MAIOR TAXA são mineradas primeiro:
//   - Adicione campo "Taxa" à Transação
//   - Ao criar um bloco, ordene por taxa (maior primeiro)
//   - Limite de 5 transações por bloco
//   - A taxa vai para o minerador junto com a recompensa
//
// Exercício 2.4 — Persistência em Arquivo
// Salve a blockchain em um arquivo JSON:
//   - SalvarEmArquivo(filename string)
//   - CarregarDeArquivo(filename string) *Blockchain
// Após minerar alguns blocos, salve, feche o programa, e carregue novamente.
// Verifique que a cadeia carregada é válida.
//
// Exercício 2.5 — Simulação de Rede (2 nós)
// Simule 2 nós minerando ao mesmo tempo:
//   - Nó A e Nó B têm cópias da mesma blockchain
//   - Ambos recebem transações diferentes
//   - Ambos mineram em goroutines simultâneas
//   - Quem terminar primeiro, o outro aceita o bloco
//   - Use channels para comunicação entre nós
// Este exercício é DIFÍCIL — é como consenso funciona de verdade!
//
// ============================================================================

func main() {
	// TODO: Exercício 2.1 — Dificuldade Dinâmica

	// TODO: Exercício 2.2 — Validação de Saldo

	// TODO: Exercício 2.3 — Pool com Prioridade

	// TODO: Exercício 2.4 — Persistência em Arquivo

	// TODO: Exercício 2.5 — Simulação de Rede
}
