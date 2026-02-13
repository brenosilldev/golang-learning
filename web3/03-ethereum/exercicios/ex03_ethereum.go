package main

// ============================================================================
// EXERCÍCIOS — Ethereum e Go
// ============================================================================
//
// Todos os exercícios devem ser feitos com Ganache ou Hardhat rodando localmente.
//
// Exercício 3.1 — Monitor de Blocos
// Conecte ao nó local e implemente um loop que:
//   - Imprime o número do bloco atual a cada 2 segundos
//   - Se o bloco mudou, mostra detalhes (hash, timestamp, TXs)
//   - Funciona até o usuário pressionar Ctrl+C (use signal.Notify)
//
// Exercício 3.2 — Consultor de Carteira
// Crie um programa que recebe um endereço como argumento e mostra:
//   - Saldo em Wei, Gwei e ETH
//   - Nonce (quantas TXs enviadas)
//   - Últimas 5 transações (itere pelos blocos recentes)
//
// Exercício 3.3 — Transferência em Lote
// Envie ETH para múltiplos endereços de uma vez:
//   - Defina uma lista: [ {endereço, valor}, ... ]
//   - Envie para cada um, incrementando o nonce
//   - Mostre o status de cada transferência
//   - Calcule o custo total em gas
//
// Exercício 3.4 — Conversor de Unidades
// Crie uma CLI que converte entre unidades:
//   ./conversor 1.5 ETH        → mostra em Wei e Gwei
//   ./conversor 1000000000 Wei  → mostra em ETH e Gwei
//   ./conversor 20 Gwei         → mostra em Wei e ETH
//
// Exercício 3.5 — Gas Tracker
// Monitore o preço do gas em tempo real:
//   - A cada bloco, consulte SuggestGasPrice
//   - Calcule: preço atual, média dos últimos 10 blocos, máximo, mínimo
//   - Mostre quanto custaria uma transferência simples (21000 gas)
//
// ============================================================================

func main() {
	// TODO: Exercício 3.1 — Monitor de Blocos

	// TODO: Exercício 3.2 — Consultor de Carteira

	// TODO: Exercício 3.3 — Transferência em Lote

	// TODO: Exercício 3.4 — Conversor de Unidades

	// TODO: Exercício 3.5 — Gas Tracker
}
