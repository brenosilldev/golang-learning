package main

import (
	"fmt"
)

// ============================================================================
// ETHEREUM — Conectar e Ler Dados
// ============================================================================
//
// Este exemplo mostra como conectar a um nó Ethereum via Go e ler
// informações sobre a rede (blocos, transações, chain ID).
//
// SETUP:
//   1. go get github.com/ethereum/go-ethereum
//   2. Rode uma blockchain local:
//      npx ganache           (ou)
//      npx hardhat node
//   3. Descomente o código abaixo e rode: go run conectar.go
//
// ============================================================================

// import (
// 	"context"
// 	"log"
// 	"math/big"
//
// 	"github.com/ethereum/go-ethereum/ethclient"
// )

// func main() {
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("=== CONECTAR AO NÓ ETHEREUM ===")
//
// 	// Conectar (use localhost para rede local)
// 	client, err := ethclient.Dial("http://localhost:8545")
// 	if err != nil {
// 		log.Fatal("Erro ao conectar:", err)
// 	}
// 	defer client.Close()
// 	fmt.Println("✅ Conectado!")
//
// 	ctx := context.Background()
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== INFORMAÇÕES DA REDE ===")
//
// 	// Chain ID (1 = Mainnet, 11155111 = Sepolia, 1337 = Ganache)
// 	chainID, err := client.ChainID(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Chain ID: %d\n", chainID)
//
// 	// Número do último bloco
// 	blocoAtual, err := client.BlockNumber(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Bloco atual: %d\n", blocoAtual)
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== LER UM BLOCO ===")
//
// 	// Pegar o bloco mais recente
// 	bloco, err := client.BlockByNumber(ctx, nil) // nil = mais recente
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	fmt.Printf("Bloco #%d\n", bloco.Number().Uint64())
// 	fmt.Printf("  Timestamp: %d\n", bloco.Time())
// 	fmt.Printf("  Hash:      %s\n", bloco.Hash().Hex())
// 	fmt.Printf("  Anterior:  %s\n", bloco.ParentHash().Hex())
// 	fmt.Printf("  Gas Used:  %d\n", bloco.GasUsed())
// 	fmt.Printf("  Gas Limit: %d\n", bloco.GasLimit())
// 	fmt.Printf("  TXs:       %d\n", len(bloco.Transactions()))
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== LER TRANSAÇÕES DO BLOCO ===")
//
// 	for i, tx := range bloco.Transactions() {
// 		fmt.Printf("\n  TX #%d\n", i)
// 		fmt.Printf("    Hash:   %s\n", tx.Hash().Hex())
// 		fmt.Printf("    Valor:  %s Wei\n", tx.Value().String())
// 		fmt.Printf("    Gas:    %d\n", tx.Gas())
// 		if tx.To() != nil {
// 			fmt.Printf("    Para:   %s\n", tx.To().Hex())
// 		}
// 	}
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== LER BLOCO ESPECÍFICO ===")
//
// 	// Ler bloco #0 (genesis)
// 	genesis, err := client.BlockByNumber(ctx, big.NewInt(0))
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Bloco Genesis:\n")
// 	fmt.Printf("  Hash:      %s\n", genesis.Hash().Hex())
// 	fmt.Printf("  Gas Limit: %d\n", genesis.GasLimit())
// }

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  ETHEREUM — Conectar e Ler Dados         ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("Para rodar este exemplo:")
	fmt.Println("  1. go get github.com/ethereum/go-ethereum")
	fmt.Println("  2. Rode: npx ganache  (ou: npx hardhat node)")
	fmt.Println("  3. Descomente o código neste arquivo")
	fmt.Println("  4. go run conectar.go")
	fmt.Println()
	fmt.Println("O código conecta ao nó Ethereum e lê:")
	fmt.Println("  • Chain ID da rede")
	fmt.Println("  • Número do bloco atual")
	fmt.Println("  • Detalhes do bloco (timestamp, hash, gas)")
	fmt.Println("  • Transações dentro do bloco")
}
