package main

// ============================================================================
// EXERCÍCIOS — DeFi
// ============================================================================
//
// Exercício 6.1 — DEX com Múltiplos Pools
// Expanda o AMM para suportar múltiplos pares:
//   - ETH/USDT, ETH/DAI, USDT/DAI
//   - Roteamento: converter ETH → DAI passando por ETH → USDT → DAI
//   - Encontrar a melhor rota (preço mais baixo) automaticamente
//   - Calcular slippage total da rota
//
// Exercício 6.2 — Lending Protocol
// Implemente um protocolo de empréstimo simplificado:
//   - Depositar(token, valor) → ganha juros (APY simulado)
//   - Emprestar(token, valor, colateral) → exige 150% de colateral
//   - Liquidar() → se colateral cair abaixo de 120%, liquidar automaticamente
//   - Simule mudanças de preço e observe liquidações
//
// Exercício 6.3 — Flash Loan Simulator
// Simule um flash loan:
//   - Pool tem 10.000 ETH
//   - Usuário pega empréstimo de 1.000 ETH
//   - Faz arbitragem entre dois pools (compra barato/vende caro)
//   - Devolve 1.000 ETH + taxa (0.09%)
//   - Se não devolver no mesmo "bloco" → REVERT
//   - Calcule o lucro do arbitragista
//
// Exercício 6.4 — Yield Farm Simulator
// Simule yield farming por 30 dias:
//   - Pool de liquidez com APY de 50%
//   - Compounding diário (reinvestir)
//   - Compare: sem compound vs compound diário
//   - Considere impermanent loss com variação de 20% no preço
//   - Calcule o rendimento real (APY - IL)
//
// ============================================================================

func main() {
	// TODO: Exercício 6.1 — DEX Multi-Pool

	// TODO: Exercício 6.2 — Lending Protocol

	// TODO: Exercício 6.3 — Flash Loan Simulator

	// TODO: Exercício 6.4 — Yield Farm Simulator
}
