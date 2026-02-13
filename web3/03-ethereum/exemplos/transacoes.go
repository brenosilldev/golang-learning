package main

import (
	"fmt"
)

// ============================================================================
// ETHEREUM — Enviar Transações
// ============================================================================
//
// Para enviar ETH de uma conta para outra, você precisa:
//   1. Chave privada do remetente
//   2. Endereço do destinatário
//   3. Valor em Wei
//   4. Gas limit e gas price
//   5. Nonce (contador de transações do remetente)
//
// SETUP:
//   1. go get github.com/ethereum/go-ethereum
//   2. Rode blockchain local: npx ganache
//   3. Copie uma chave privada do Ganache
//   4. Descomente o código e rode: go run transacoes.go
//
// ⚠️  NUNCA coloque chaves privadas REAIS no código!
//     Use variáveis de ambiente: os.Getenv("PRIVATE_KEY")
// ============================================================================

// import (
// 	"context"
// 	"crypto/ecdsa"
// 	"log"
// 	"math/big"
//
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/core/types"
// 	"github.com/ethereum/go-ethereum/crypto"
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
// 	fmt.Println("=== PREPARAR TRANSAÇÃO ===")
//
// 	// 1. Chave privada (do Ganache — NUNCA use chaves reais assim!)
// 	// Pegue uma das chaves que o Ganache mostra ao iniciar
// 	privateKeyHex := "SUA_CHAVE_PRIVADA_DO_GANACHE_AQUI" // sem 0x
//
// 	privateKey, err := crypto.HexToECDSA(privateKeyHex)
// 	if err != nil {
// 		log.Fatal("Chave privada inválida:", err)
// 	}
//
// 	// 2. Derivar endereço do remetente
// 	publicKey := privateKey.Public().(*ecdsa.PublicKey)
// 	fromAddress := crypto.PubkeyToAddress(*publicKey)
// 	fmt.Printf("De: %s\n", fromAddress.Hex())
//
// 	// 3. Endereço de destino
// 	toAddress := common.HexToAddress("0xDESTINATARIO_AQUI")
// 	fmt.Printf("Para: %s\n", toAddress.Hex())
//
// 	// 4. Obter nonce (próximo número de transação)
// 	nonce, err := client.PendingNonceAt(ctx, fromAddress)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Nonce: %d\n", nonce)
//
// 	// 5. Definir valor (1 ETH = 10^18 Wei)
// 	valor := new(big.Int)
// 	valor.SetString("1000000000000000000", 10) // 1 ETH em Wei
// 	fmt.Printf("Valor: 1 ETH\n")
//
// 	// 6. Gas
// 	gasLimit := uint64(21000) // 21000 = transferência simples de ETH
// 	gasPrice, err := client.SuggestGasPrice(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("Gas Limit: %d\n", gasLimit)
// 	fmt.Printf("Gas Price: %s Wei\n", gasPrice.String())
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== CRIAR E ASSINAR TRANSAÇÃO ===")
//
// 	// Criar transação
// 	tx := types.NewTransaction(nonce, toAddress, valor, gasLimit, gasPrice, nil)
//
// 	// Obter Chain ID para EIP-155 (proteção contra replay)
// 	chainID, err := client.ChainID(ctx)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// Assinar com a chave privada
// 	signedTx, err := types.SignTx(tx, types.NewEIP155Signer(chainID), privateKey)
// 	if err != nil {
// 		log.Fatal(err)
// 	}
// 	fmt.Printf("TX Hash: %s\n", signedTx.Hash().Hex())
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== ENVIAR TRANSAÇÃO ===")
//
// 	err = client.SendTransaction(ctx, signedTx)
// 	if err != nil {
// 		log.Fatal("Erro ao enviar:", err)
// 	}
// 	fmt.Printf("✅ Transação enviada! Hash: %s\n", signedTx.Hash().Hex())
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== VERIFICAR RECIBO ===")
//
// 	// Aguardar confirmação (em rede local é instantâneo)
// 	receipt, err := client.TransactionReceipt(ctx, signedTx.Hash())
// 	if err != nil {
// 		fmt.Println("Aguardando confirmação...")
// 	} else {
// 		fmt.Printf("Status: %d (1 = sucesso, 0 = falha)\n", receipt.Status)
// 		fmt.Printf("Bloco:  %d\n", receipt.BlockNumber.Uint64())
// 		fmt.Printf("Gas:    %d\n", receipt.GasUsed)
// 	}
// }

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  ENVIAR TRANSAÇÕES ETHEREUM              ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("O fluxo para enviar ETH:")
	fmt.Println()
	fmt.Println("  1. Carregar CHAVE PRIVADA do remetente")
	fmt.Println("  2. Derivar ENDEREÇO do remetente")
	fmt.Println("  3. Obter NONCE (contador de TXs)")
	fmt.Println("  4. Definir VALOR em Wei")
	fmt.Println("  5. Estimar GAS (limite e preço)")
	fmt.Println("  6. Criar TX e ASSINAR com chave privada")
	fmt.Println("  7. ENVIAR para a rede")
	fmt.Println("  8. Verificar RECIBO (status, bloco, gas usado)")
	fmt.Println()
	fmt.Println("⚠️  NUNCA coloque chaves privadas no código!")
	fmt.Println("   Use: os.Getenv(\"PRIVATE_KEY\")")
	fmt.Println()
	fmt.Println("Para rodar:")
	fmt.Println("  1. go get github.com/ethereum/go-ethereum")
	fmt.Println("  2. npx ganache  (blockchain local)")
	fmt.Println("  3. Copie uma private key do Ganache")
	fmt.Println("  4. Descomente o código neste arquivo")
}
