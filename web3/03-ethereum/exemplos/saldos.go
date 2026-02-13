package main

import (
	"fmt"
	"math"
	"math/big"
)

// ============================================================================
// ETHEREUM — Saldos e Conversão de Unidades
// ============================================================================
//
// ETH tem 18 casas decimais! Para trabalhar com valores, usamos:
//   Wei  = menor unidade (1 ETH = 10^18 Wei)
//   Gwei = para gas (1 ETH = 10^9 Gwei)
//   ETH  = o que humanos leem
//
// O ethclient SEMPRE retorna valores em Wei (big.Int).
// Você precisa converter para exibir ao usuário.
//
// SETUP:
//   1. go get github.com/ethereum/go-ethereum
//   2. Rode blockchain local: npx ganache
//   3. Descomente o código e rode: go run saldos.go
// ============================================================================

// --- Funções de conversão (ESSENCIAIS em Web3) ---

// WeiParaETH converte Wei (big.Int) para ETH (float64)
func WeiParaETH(wei *big.Int) float64 {
	// 1 ETH = 10^18 Wei
	fWei := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetFloat64(math.Pow10(18))
	resultado, _ := new(big.Float).Quo(fWei, divisor).Float64()
	return resultado
}

// ETHParaWei converte ETH (float64) para Wei (big.Int)
func ETHParaWei(eth float64) *big.Int {
	// Multiplicar por 10^18
	fETH := new(big.Float).SetFloat64(eth)
	multiplicador := new(big.Float).SetFloat64(math.Pow10(18))
	resultado := new(big.Float).Mul(fETH, multiplicador)

	wei := new(big.Int)
	resultado.Int(wei)
	return wei
}

// WeiParaGwei converte Wei para Gwei
func WeiParaGwei(wei *big.Int) float64 {
	// 1 Gwei = 10^9 Wei
	fWei := new(big.Float).SetInt(wei)
	divisor := new(big.Float).SetFloat64(math.Pow10(9))
	resultado, _ := new(big.Float).Quo(fWei, divisor).Float64()
	return resultado
}

// --- Código com ethclient (descomente após instalar go-ethereum) ---

// import (
// 	"context"
// 	"log"
//
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// func main() {
// 	client, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	defer client.Close()
//
// 	ctx := context.Background()
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("=== CONSULTAR SALDO ===")
//
// 	// Endereço de teste do Ganache (substitua pelo seu)
// 	endereco := common.HexToAddress("0x742d35Cc6634C0532925a3b844Bc9e7595f2bD18")
//
// 	// Saldo retorna em Wei (big.Int)
// 	saldoWei, err := client.BalanceAt(ctx, endereco, nil) // nil = bloco mais recente
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	fmt.Printf("Endereço: %s\n", endereco.Hex())
// 	fmt.Printf("Saldo (Wei):  %s\n", saldoWei.String())
// 	fmt.Printf("Saldo (Gwei): %.2f\n", WeiParaGwei(saldoWei))
// 	fmt.Printf("Saldo (ETH):  %.4f\n", WeiParaETH(saldoWei))
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== SALDO EM BLOCO ESPECÍFICO ===")
//
// 	// Saldo no bloco #0 (genesis)
// 	saldoGenesis, _ := client.BalanceAt(ctx, endereco, big.NewInt(0))
// 	fmt.Printf("Saldo no genesis: %.4f ETH\n", WeiParaETH(saldoGenesis))
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== NONCE (contador de transações) ===")
//
// 	// Nonce = quantas transações este endereço já enviou
// 	nonce, _ := client.PendingNonceAt(ctx, endereco)
// 	fmt.Printf("Nonce: %d (= já enviou %d transações)\n", nonce, nonce)
// }

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  SALDOS E CONVERSÃO DE UNIDADES          ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Demonstrar conversões (funciona sem blockchain)
	fmt.Println("\n━━━ CONVERSÕES ━━━")

	// 1 ETH em Wei
	umETH := ETHParaWei(1.0)
	fmt.Printf("1 ETH = %s Wei\n", umETH.String())
	fmt.Printf("       = %.0f Gwei\n", WeiParaGwei(umETH))

	// 0.5 ETH
	meioETH := ETHParaWei(0.5)
	fmt.Printf("0.5 ETH = %s Wei\n", meioETH.String())

	// Converter de volta
	fmt.Printf("Volta: %s Wei = %.4f ETH\n", umETH.String(), WeiParaETH(umETH))

	// Simular saldos de contas
	fmt.Println("\n━━━ SALDOS SIMULADOS ━━━")
	contas := map[string]float64{
		"0xAlice...": 100.0,
		"0xBob...":   45.5,
		"0xCarol...": 0.001,
	}

	for endereco, eth := range contas {
		wei := ETHParaWei(eth)
		fmt.Printf("  %s: %.3f ETH = %s Wei\n", endereco, eth, wei.String())
	}

	// Gas estimation example
	fmt.Println("\n━━━ CÁLCULO DE GAS (exemplo) ━━━")
	gasLimit := int64(21000) // transferência simples
	gasPriceGwei := 20.0     // 20 Gwei
	custoWei := ETHParaWei(gasPriceGwei / 1e9 * float64(gasLimit))
	fmt.Printf("  Gas Limit: %d\n", gasLimit)
	fmt.Printf("  Gas Price: %.0f Gwei\n", gasPriceGwei)
	fmt.Printf("  Custo TX:  %.6f ETH (%s Wei)\n", WeiParaETH(custoWei), custoWei)

	fmt.Println("\n━━━ PARA RODAR COM ETHEREUM REAL ━━━")
	fmt.Println("  1. go get github.com/ethereum/go-ethereum")
	fmt.Println("  2. npx ganache  (blockchain local)")
	fmt.Println("  3. Descomente o código neste arquivo")
}
