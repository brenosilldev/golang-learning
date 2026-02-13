package main

import (
	"fmt"
)

// ============================================================================
// SMART CONTRACTS — Interagir via Go
// ============================================================================
//
// Este exemplo mostra como:
//   1. Fazer deploy de um smart contract
//   2. Chamar funções (leitura e escrita)
//   3. Escutar eventos em tempo real
//
// SETUP:
//   1. Compilar o contrato Solidity:
//      solc --abi --bin Cofre.sol -o build/
//
//   2. Gerar bindings Go:
//      abigen --abi=build/Cofre.abi --bin=build/Cofre.bin \
//             --pkg=cofre --out=cofre.go
//
//   3. Descomente o código e rode:
//      go run interagir.go
// ============================================================================

// import (
// 	"context"
// 	"crypto/ecdsa"
// 	"log"
// 	"math/big"
//
// 	"github.com/ethereum/go-ethereum/accounts/abi/bind"
// 	"github.com/ethereum/go-ethereum/common"
// 	"github.com/ethereum/go-ethereum/crypto"
// 	"github.com/ethereum/go-ethereum/ethclient"
//
//  // Este import é do binding gerado pelo abigen
// 	// cofre "./cofre"
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
// 	// Chave privada do deployer (do Ganache)
// 	privateKey, _ := crypto.HexToECDSA("SUA_CHAVE_DO_GANACHE")
// 	publicKey := privateKey.Public().(*ecdsa.PublicKey)
// 	fromAddress := crypto.PubkeyToAddress(*publicKey)
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("=== DEPLOY DO CONTRATO ===")
//
// 	nonce, _ := client.PendingNonceAt(ctx, fromAddress)
// 	gasPrice, _ := client.SuggestGasPrice(ctx)
// 	chainID, _ := client.ChainID(ctx)
//
// 	auth, _ := bind.NewKeyedTransactorWithChainID(privateKey, chainID)
// 	auth.Nonce = big.NewInt(int64(nonce))
// 	auth.Value = big.NewInt(0)       // ETH a enviar
// 	auth.GasLimit = uint64(3000000)  // gas limit para deploy
// 	auth.GasPrice = gasPrice
//
// 	// Deploy! (cofre.DeployCofre é gerado pelo abigen)
// 	// address, tx, instance, err := cofre.DeployCofre(auth, client)
// 	// if err != nil {
// 	// 	log.Fatal(err)
// 	// }
// 	// fmt.Printf("Contrato deployed em: %s\n", address.Hex())
// 	// fmt.Printf("TX Hash: %s\n", tx.Hash().Hex())
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== LER DO CONTRATO (view — gratuito) ===")
//
// 	// Funções view não gastam gas!
// 	// dono, _ := instance.Dono(&bind.CallOpts{})
// 	// fmt.Printf("Dono: %s\n", dono.Hex())
//
// 	// total, _ := instance.TotalDepositado(&bind.CallOpts{})
// 	// fmt.Printf("Total depositado: %s Wei\n", total.String())
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== ESCREVER NO CONTRATO (gasta gas) ===")
//
// 	// Depositar 1 ETH
// 	// auth.Value = big.NewInt(1e18) // 1 ETH em Wei
// 	// tx, _ = instance.Depositar(auth)
// 	// fmt.Printf("Depósito TX: %s\n", tx.Hash().Hex())
//
// 	// Sacar 0.5 ETH
// 	// auth.Value = big.NewInt(0)
// 	// meioETH := new(big.Int).Div(big.NewInt(1e18), big.NewInt(2))
// 	// tx, _ = instance.Sacar(auth, meioETH)
// 	// fmt.Printf("Saque TX: %s\n", tx.Hash().Hex())
//
// 	// ══════════════════════════════════════════════════════════
// 	fmt.Println("\n=== ESCUTAR EVENTOS ===")
//
// 	// Eventos são como webhooks da blockchain
// 	// watchOpts := &bind.WatchOpts{Context: ctx}
//
// 	// Canal de eventos de depósito
// 	// depositoChan := make(chan *cofre.CofreDeposito)
// 	// sub, _ := instance.WatchDeposito(watchOpts, depositoChan, nil)
//
// 	// go func() {
// 	// 	for {
// 	// 		select {
// 	// 		case evento := <-depositoChan:
// 	// 			fmt.Printf("🔔 Depósito! De: %s, Valor: %s Wei\n",
// 	// 				evento.De.Hex(), evento.Valor.String())
// 	// 		case err := <-sub.Err():
// 	// 			log.Fatal(err)
// 	// 		}
// 	// 	}
// 	// }()
// }

func main() {
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║  SMART CONTRACTS — Interação via Go      ║")
	fmt.Println("╚══════════════════════════════════════════╝")
	fmt.Println()
	fmt.Println("O fluxo para interagir com smart contracts:")
	fmt.Println()
	fmt.Println("  1. Escrever contrato: Cofre.sol")
	fmt.Println("  2. Compilar:  solc --abi --bin Cofre.sol -o build/")
	fmt.Println("  3. Gerar Go:  abigen --abi=... --bin=... --pkg=cofre --out=cofre.go")
	fmt.Println("  4. Deploy:    cofre.DeployCofre(auth, client)")
	fmt.Println("  5. Ler:       instance.Dono(&bind.CallOpts{})")
	fmt.Println("  6. Escrever:  instance.Depositar(auth)")
	fmt.Println("  7. Eventos:   instance.WatchDeposito(opts, channel)")
	fmt.Println()
	fmt.Println("Operações VIEW (leitura) são GRÁTIS!")
	fmt.Println("Operações de ESCRITA custam GAS!")
	fmt.Println()
	fmt.Println("📝 Veja o contrato em: Cofre.sol")
}
