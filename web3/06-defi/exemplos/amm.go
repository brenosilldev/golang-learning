package main

import (
	"fmt"
	"math"
)

// ============================================================================
// DeFi — AMM (Automated Market Maker) Simulação
// ============================================================================
//
// Simula uma DEX tipo Uniswap com a fórmula x * y = k
// Entender AMM é FUNDAMENTAL para trabalhar com DeFi.
//
// Rode com: go run amm.go
// ============================================================================

// LiquidityPool representa um par de tokens em um AMM
type LiquidityPool struct {
	TokenA         string
	TokenB         string
	ReservaA       float64
	ReservaB       float64
	K              float64 // constante x*y=k
	TaxaSwap       float64 // 0.003 = 0.3%
	TaxasColetadas float64
}

// NovoPool cria um pool de liquidez
func NovoPool(tokenA, tokenB string, qtdA, qtdB float64) *LiquidityPool {
	return &LiquidityPool{
		TokenA:   tokenA,
		TokenB:   tokenB,
		ReservaA: qtdA,
		ReservaB: qtdB,
		K:        qtdA * qtdB,
		TaxaSwap: 0.003, // 0.3% como Uniswap v2
	}
}

// Preco retorna o preço de A em termos de B
func (p *LiquidityPool) Preco() float64 {
	return p.ReservaB / p.ReservaA
}

// CotarSwap calcula quanto de tokenB você recebe por um dado valor de tokenA
// SEM executar o swap
func (p *LiquidityPool) CotarSwap(qtdA float64) (qtdB float64, slippage float64) {
	// Aplicar taxa
	qtdAComTaxa := qtdA * (1 - p.TaxaSwap)

	// x * y = k → novaReservaA * novaReservaB = k
	novaReservaA := p.ReservaA + qtdAComTaxa
	novaReservaB := p.K / novaReservaA
	qtdB = p.ReservaB - novaReservaB

	// Calcular slippage (diferença do preço ideal)
	precoIdeal := qtdA * p.Preco()
	slippage = (precoIdeal - qtdB) / precoIdeal * 100

	return qtdB, slippage
}

// Swap executa a troca de tokenA por tokenB
func (p *LiquidityPool) Swap(trader string, qtdA float64) (float64, error) {
	if qtdA <= 0 {
		return 0, fmt.Errorf("quantidade deve ser positiva")
	}

	// Calcular taxa
	taxa := qtdA * p.TaxaSwap
	qtdAComTaxa := qtdA - taxa

	// x * y = k
	novaReservaA := p.ReservaA + qtdAComTaxa
	novaReservaB := p.K / novaReservaA
	qtdB := p.ReservaB - novaReservaB

	if qtdB <= 0 {
		return 0, fmt.Errorf("liquidez insuficiente")
	}

	// Atualizar reservas
	precoAntes := p.Preco()
	p.ReservaA = novaReservaA
	p.ReservaB = novaReservaB
	p.TaxasColetadas += taxa
	precoDepois := p.Preco()

	impacto := (precoDepois - precoAntes) / precoAntes * 100

	fmt.Printf("  🔄 Swap: %s troca %.4f %s → %.4f %s\n",
		trader, qtdA, p.TokenA, qtdB, p.TokenB)
	fmt.Printf("     Preço: %.2f → %.2f (impacto: %.2f%%)\n",
		precoAntes, precoDepois, impacto)
	fmt.Printf("     Taxa: %.4f %s\n", taxa, p.TokenA)

	return qtdB, nil
}

// AdicionarLiquidez deposita tokens proporcionalmente
func (p *LiquidityPool) AdicionarLiquidez(provedor string, qtdA float64) float64 {
	// Deve manter a proporção atual
	proporcao := p.ReservaB / p.ReservaA
	qtdB := qtdA * proporcao

	p.ReservaA += qtdA
	p.ReservaB += qtdB
	p.K = p.ReservaA * p.ReservaB // recalcular k

	fmt.Printf("  💧 Liquidez: %s deposita %.2f %s + %.2f %s\n",
		provedor, qtdA, p.TokenA, qtdB, p.TokenB)

	return qtdB
}

// Status mostra o estado atual do pool
func (p *LiquidityPool) Status() {
	fmt.Printf("  📊 Pool %s/%s:\n", p.TokenA, p.TokenB)
	fmt.Printf("     Reserva %s: %.4f\n", p.TokenA, p.ReservaA)
	fmt.Printf("     Reserva %s: %.4f\n", p.TokenB, p.ReservaB)
	fmt.Printf("     Preço 1 %s = %.4f %s\n", p.TokenA, p.Preco(), p.TokenB)
	fmt.Printf("     K = %.2f\n", p.K)
	fmt.Printf("     Taxas coletadas: %.4f %s\n", p.TaxasColetadas, p.TokenA)
}

// ═══════════════════════════════════════════════
// Simulação de Arbitragem
// ═══════════════════════════════════════════════

func verificarArbitragem(poolA, poolB *LiquidityPool) {
	precoA := poolA.Preco()
	precoB := poolB.Preco()
	diff := math.Abs(precoA-precoB) / math.Min(precoA, precoB) * 100

	fmt.Printf("\n  📈 Arbitragem:\n")
	fmt.Printf("     Pool A: 1 %s = %.4f %s\n", poolA.TokenA, precoA, poolA.TokenB)
	fmt.Printf("     Pool B: 1 %s = %.4f %s\n", poolB.TokenA, precoB, poolB.TokenB)
	fmt.Printf("     Diferença: %.2f%%\n", diff)

	if diff > 1 {
		fmt.Println("     🤑 Oportunidade de arbitragem detectada!")
		if precoA < precoB {
			fmt.Printf("     → Comprar em Pool A (mais barato), vender em Pool B\n")
		} else {
			fmt.Printf("     → Comprar em Pool B (mais barato), vender em Pool A\n")
		}
	} else {
		fmt.Println("     → Sem oportunidade significativa")
	}
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║   DeFi — AMM (Uniswap Simplificado)     ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Criar pool ETH/USDT (como Uniswap)
	pool := NovoPool("ETH", "USDT", 1000, 2_000_000) // 1 ETH = 2000 USDT

	fmt.Println("\n━━━ ESTADO INICIAL ━━━")
	pool.Status()

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ COTAÇÃO (sem executar) ━━━")

	tamanhos := []float64{1, 10, 50, 100}
	for _, qtd := range tamanhos {
		recebe, slippage := pool.CotarSwap(qtd)
		fmt.Printf("  %.0f ETH → %.2f USDT (slippage: %.2f%%)\n", qtd, recebe, slippage)
	}
	fmt.Println("  → Quanto MAIS compra, MAIS caro fica (slippage)")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ SWAPS ━━━")

	pool.Swap("Alice", 5)   // Alice compra ~10000 USDT
	pool.Swap("Bob", 10)    // Bob compra ~20000 USDT
	pool.Swap("Carol", 0.5) // Carol faz um swap pequeno

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ESTADO APÓS SWAPS ━━━")
	pool.Status()

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ADICIONAR LIQUIDEZ ━━━")

	pool.AdicionarLiquidez("Dave", 100) // Dave deposita proporcionalmente
	pool.Status()

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ARBITRAGEM ━━━")

	// Dois pools com preços diferentes
	pool1 := NovoPool("ETH", "USDT", 1000, 2_000_000) // 1 ETH = 2000 USDT
	pool2 := NovoPool("ETH", "USDT", 500, 1_200_000)  // 1 ETH = 2400 USDT

	verificarArbitragem(pool1, pool2)

	// Após muitos swaps em pool1, preço se equilibra
	pool1.Swap("Trader", 50)
	verificarArbitragem(pool1, pool2)

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ IMPERMANENT LOSS (risco para LPs) ━━━")

	// Simulação do risco para provedores de liquidez
	poolIL := NovoPool("ETH", "USDT", 10, 20_000) // 1 ETH = 2000 USDT
	depositoInicial := 10*2000.0 + 20_000.0       // $40.000 total

	fmt.Printf("  Depósito inicial: $%.0f (10 ETH + 20.000 USDT)\n", depositoInicial)

	// Simular ETH subindo de preço (muita compra de ETH no pool)
	poolIL.Swap("Mercado", 200_000) // swap inverso: muitos USDT
	// Recalcular: novo preço de ETH

	novoPrecoETH := poolIL.Preco()
	valorNoPool := poolIL.ReservaA*novoPrecoETH + poolIL.ReservaB
	valorSeNaoDepositar := 10*novoPrecoETH + 20_000.0

	fmt.Printf("  Novo preço ETH: $%.2f\n", novoPrecoETH)
	fmt.Printf("  Valor mantendo no pool: $%.2f\n", valorNoPool)
	fmt.Printf("  Valor se tivesse segurado: $%.2f\n", valorSeNaoDepositar)
	fmt.Printf("  Impermanent Loss: $%.2f (%.2f%%)\n",
		valorSeNaoDepositar-valorNoPool,
		(valorSeNaoDepositar-valorNoPool)/valorSeNaoDepositar*100)
}
