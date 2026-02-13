package main

// ============================================================================
// EXERCÍCIOS — Indexer e Explorer
// ============================================================================
//
// Exercício 7.1 — Indexer com Persistência
// Modifique o indexer do exemplo para salvar em arquivo JSON:
//   - Ao indexar, salve em indexer_db.json
//   - Ao iniciar, carregue de indexer_db.json se existir
//   - Só indexe blocos NOVOS (a partir do último indexado)
//   - Teste: rode 2x e veja que não re-indexa blocos antigos
//
// Exercício 7.2 — API REST para o Indexer
// Crie uma API HTTP com net/http que expõe:
//   GET /blocks           → lista últimos 10 blocos
//   GET /blocks/:number   → detalhes de um bloco
//   GET /txs/:hash        → detalhes de uma transação
//   GET /address/:addr    → saldo + histórico de um endereço
//   GET /stats            → estatísticas gerais
// Use JSON como formato de resposta.
//
// Exercício 7.3 — Event Indexer (Smart Contract)
// Crie um indexer especializado em EVENTOS de smart contracts:
//   - Filtre logs por contrato e tipo de evento
//   - Decodifique eventos Transfer (ERC-20) e armazene
//   - Permita consultar: "todas as transferências do token X"
//   - Permita consultar: "todos os tokens enviados por endereço Y"
//   - Se usar Ganache, deploy o ERC-20 do módulo 04 e indexe eventos reais
//
// Exercício 7.4 — Real-time Dashboard
// Combine tudo em um dashboard que:
//   - Conecta via WebSocket (ou polling) ao nó
//   - Mostra blocos chegando em tempo real no terminal
//   - Mostra transações mais caras (maior valor)
//   - Mostra endereços mais ativos
//   - Atualiza a cada novo bloco
//   - Use goroutines: 1 para escutar, 1 para processar, 1 para exibir
//
// ============================================================================

func main() {
	// TODO: Exercício 7.1 — Indexer com Persistência

	// TODO: Exercício 7.2 — API REST

	// TODO: Exercício 7.3 — Event Indexer

	// TODO: Exercício 7.4 — Real-time Dashboard
}
