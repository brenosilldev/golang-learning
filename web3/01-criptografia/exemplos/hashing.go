package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// CRIPTOGRAFIA — Hashing com SHA-256
// ============================================================================
//
// SHA-256 é a fundação de TUDO em blockchain.
// Este exemplo demonstra todas as propriedades e como a blockchain usa hashing.
//
// Rode com: go run hashing.go
// ============================================================================

// calcHash é nossa função utilitária — usaremos em TODO o curso
func calcHash(dados string) string {
	hash := sha256.Sum256([]byte(dados))
	return hex.EncodeToString(hash[:])
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       SHA-256 — HASHING BÁSICO           ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	texto := "Olá, blockchain!"
	hash := calcHash(texto)
	fmt.Printf("\nTexto:   %s\n", texto)
	fmt.Printf("SHA-256: %s\n", hash)
	fmt.Printf("Tamanho: %d caracteres (SEMPRE 64)\n", len(hash))

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ PROPRIEDADE 1: DETERMINÍSTICO ━━━")
	// Mesmo input → mesmo hash, sempre, em qualquer computador do mundo

	for i := 1; i <= 3; i++ {
		h := calcHash("Go é incrível")
		fmt.Printf("  Tentativa %d: %s\n", i, h[:32]+"...")
	}
	fmt.Println("  → Todos iguais! ✅")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ PROPRIEDADE 2: EFEITO AVALANCHE ━━━")
	// Mudar UM caractere = hash completamente diferente

	entradas := []string{
		"Blockchain",
		"blockchain",  // apenas 'B' → 'b'
		"Blockchain!", // adicionou '!'
	}

	for _, e := range entradas {
		fmt.Printf("  %-14q → %s\n", e, calcHash(e)[:24]+"...")
	}
	fmt.Println("  → Completamente diferentes! ✅")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ PROPRIEDADE 3: IRREVERSÍVEL ━━━")

	senhaHash := calcHash("minhaSenhaSecreta")
	fmt.Printf("  Hash armazenado: %s\n", senhaHash[:32]+"...")
	fmt.Println("  ❌ Impossível voltar ao texto original")
	fmt.Println("  ✅ Mas posso VERIFICAR:")

	tentativas := []string{"senhaErrada", "minhaSenhaSecreta", "123456"}
	for _, t := range tentativas {
		match := calcHash(t) == senhaHash
		emoji := "❌"
		if match {
			emoji = "✅"
		}
		fmt.Printf("     %q → %s\n", t, emoji)
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ COMO A BLOCKCHAIN USA HASHING ━━━")

	// Simular encadeamento de blocos
	type BlocoSimples struct {
		Numero       int
		Dados        string
		HashAnterior string
		Hash         string
	}

	// Bloco Genesis (primeiro bloco)
	genesis := BlocoSimples{
		Numero:       0,
		Dados:        "Bloco Genesis",
		HashAnterior: strings.Repeat("0", 64), // não tem anterior
	}
	genesis.Hash = calcHash(fmt.Sprintf("%d|%s|%s", genesis.Numero, genesis.Dados, genesis.HashAnterior))

	// Bloco 1 aponta para o genesis
	bloco1 := BlocoSimples{
		Numero:       1,
		Dados:        "Alice→Bob:5ETH",
		HashAnterior: genesis.Hash, // ← encadeado!
	}
	bloco1.Hash = calcHash(fmt.Sprintf("%d|%s|%s", bloco1.Numero, bloco1.Dados, bloco1.HashAnterior))

	// Bloco 2 aponta para o bloco 1
	bloco2 := BlocoSimples{
		Numero:       2,
		Dados:        "Carol→Dave:3ETH",
		HashAnterior: bloco1.Hash, // ← encadeado!
	}
	bloco2.Hash = calcHash(fmt.Sprintf("%d|%s|%s", bloco2.Numero, bloco2.Dados, bloco2.HashAnterior))

	blocos := []BlocoSimples{genesis, bloco1, bloco2}
	for _, b := range blocos {
		fmt.Printf("\n  📦 Bloco #%d\n", b.Numero)
		fmt.Printf("     Dados:    %s\n", b.Dados)
		fmt.Printf("     Anterior: %s...\n", b.HashAnterior[:16])
		fmt.Printf("     Hash:     %s...\n", b.Hash[:16])
	}

	fmt.Println("\n  → Se alguém mudar Bloco #1, o hash muda,")
	fmt.Println("    e o Bloco #2 aponta pro hash ANTIGO → INVÁLIDO! 🔒")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ PROOF OF WORK (MINING) ━━━")

	dificuldade := 4 // "0000" no início
	prefixo := strings.Repeat("0", dificuldade)
	dadosBloco := "Transações do bloco 42"

	fmt.Printf("  Buscando hash que começa com %q...\n", prefixo)

	nonce := 0
	inicio := time.Now()

	for {
		tentativa := fmt.Sprintf("%s|nonce:%d", dadosBloco, nonce)
		hash := calcHash(tentativa)

		if strings.HasPrefix(hash, prefixo) {
			fmt.Printf("\n  ⛏️  MINERADO!\n")
			fmt.Printf("  Nonce encontrado: %d\n", nonce)
			fmt.Printf("  Hash: %s\n", hash)
			fmt.Printf("  Tempo: %v\n", time.Since(inicio))
			fmt.Printf("  Tentativas: %d\n", nonce+1)
			break
		}
		nonce++
	}

	fmt.Println("\n  💡 Quanto mais zeros exigidos, mais difícil:")
	fmt.Println("     '0'    → ~16 tentativas")
	fmt.Println("     '00'   → ~256 tentativas")
	fmt.Println("     '0000' → ~65.536 tentativas")
	fmt.Println("     Bitcoin usa ~19 zeros → trilhões de tentativas!")
}
