package main

import (
	"fmt"
)

// ============================================================================
// TOKENS E NFTs — Simulação em Go puro
// ============================================================================
//
// Este exemplo implementa ERC-20 e ERC-721 simplificados em Go
// para você entender a mecânica antes de interagir com Ethereum.
//
// Na blockchain real, essa lógica roda em Solidity.
// Aqui simulamos em Go para entender os conceitos.
//
// Rode com: go run tokens.go
// ============================================================================

// ═══════════════════════════════════════════════
// ERC-20 — Token Fungível
// ═══════════════════════════════════════════════

type ERC20 struct {
	Nome       string
	Simbolo    string
	Decimais   int
	Supply     float64
	Saldos     map[string]float64
	Allowances map[string]map[string]float64 // dono → spender → valor
}

func NovoERC20(nome, simbolo string, supplyInicial float64, criador string) *ERC20 {
	token := &ERC20{
		Nome:       nome,
		Simbolo:    simbolo,
		Decimais:   18,
		Supply:     supplyInicial,
		Saldos:     make(map[string]float64),
		Allowances: make(map[string]map[string]float64),
	}
	// Todo supply vai para o criador
	token.Saldos[criador] = supplyInicial
	return token
}

func (t *ERC20) BalanceOf(conta string) float64 {
	return t.Saldos[conta]
}

func (t *ERC20) Transfer(de, para string, valor float64) error {
	if t.Saldos[de] < valor {
		return fmt.Errorf("saldo insuficiente: tem %.2f, precisa %.2f", t.Saldos[de], valor)
	}
	t.Saldos[de] -= valor
	t.Saldos[para] += valor
	fmt.Printf("  📤 Transfer: %s → %s: %.2f %s\n", de, para, valor, t.Simbolo)
	return nil
}

func (t *ERC20) Approve(dono, spender string, valor float64) {
	if t.Allowances[dono] == nil {
		t.Allowances[dono] = make(map[string]float64)
	}
	t.Allowances[dono][spender] = valor
	fmt.Printf("  ✅ Approve: %s autoriza %s a usar %.2f %s\n", dono, spender, valor, t.Simbolo)
}

func (t *ERC20) TransferFrom(spender, de, para string, valor float64) error {
	allowance := t.Allowances[de][spender]
	if allowance < valor {
		return fmt.Errorf("allowance insuficiente: tem %.2f, precisa %.2f", allowance, valor)
	}
	if t.Saldos[de] < valor {
		return fmt.Errorf("saldo insuficiente")
	}
	t.Saldos[de] -= valor
	t.Saldos[para] += valor
	t.Allowances[de][spender] -= valor
	fmt.Printf("  📤 TransferFrom: %s move %.2f %s de %s → %s\n", spender, valor, t.Simbolo, de, para)
	return nil
}

// ═══════════════════════════════════════════════
// ERC-721 — NFT (Non-Fungible Token)
// ═══════════════════════════════════════════════

type NFTMetadata struct {
	Nome      string
	Descricao string
	Imagem    string
	Atributos map[string]string
}

type ERC721 struct {
	Nome      string
	Simbolo   string
	Donos     map[int]string      // tokenID → dono
	Metadata  map[int]NFTMetadata // tokenID → metadata
	ProximoID int
}

func NovoERC721(nome, simbolo string) *ERC721 {
	return &ERC721{
		Nome:      nome,
		Simbolo:   simbolo,
		Donos:     make(map[int]string),
		Metadata:  make(map[int]NFTMetadata),
		ProximoID: 1,
	}
}

func (n *ERC721) Mint(para string, metadata NFTMetadata) int {
	id := n.ProximoID
	n.Donos[id] = para
	n.Metadata[id] = metadata
	n.ProximoID++
	fmt.Printf("  🎨 Mint NFT #%d '%s' → %s\n", id, metadata.Nome, para)
	return id
}

func (n *ERC721) OwnerOf(tokenID int) (string, error) {
	dono, existe := n.Donos[tokenID]
	if !existe {
		return "", fmt.Errorf("token #%d não existe", tokenID)
	}
	return dono, nil
}

func (n *ERC721) TransferFrom(de, para string, tokenID int) error {
	dono, existe := n.Donos[tokenID]
	if !existe {
		return fmt.Errorf("token #%d não existe", tokenID)
	}
	if dono != de {
		return fmt.Errorf("%s não é dono do token #%d", de, tokenID)
	}
	n.Donos[tokenID] = para
	fmt.Printf("  📤 Transfer NFT #%d: %s → %s\n", tokenID, de, para)
	return nil
}

func (n *ERC721) BalanceOf(conta string) int {
	count := 0
	for _, dono := range n.Donos {
		if dono == conta {
			count++
		}
	}
	return count
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       TOKENS E NFTs                      ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ERC-20: TOKEN FUNGÍVEL ━━━")

	// Criar token (como USDT, LINK, etc.)
	token := NovoERC20("GoToken", "GTK", 1_000_000, "Alice")
	fmt.Printf("\n  Token: %s (%s)\n", token.Nome, token.Simbolo)
	fmt.Printf("  Supply: %.0f %s\n", token.Supply, token.Simbolo)
	fmt.Printf("  Alice: %.0f %s\n\n", token.BalanceOf("Alice"), token.Simbolo)

	// Transferências
	token.Transfer("Alice", "Bob", 10_000)
	token.Transfer("Alice", "Carol", 5_000)

	fmt.Printf("\n  Saldos:\n")
	for _, conta := range []string{"Alice", "Bob", "Carol"} {
		fmt.Printf("    %s: %.0f %s\n", conta, token.BalanceOf(conta), token.Simbolo)
	}

	// Approve + TransferFrom (como DeFi funciona)
	fmt.Println("\n  --- Approve + TransferFrom (DeFi) ---")
	token.Approve("Bob", "DEX_Uniswap", 5_000)                       // Bob autoriza
	token.TransferFrom("DEX_Uniswap", "Bob", "LiquidityPool", 3_000) // DEX move

	fmt.Printf("\n  Saldos após DeFi:\n")
	fmt.Printf("    Bob: %.0f %s\n", token.BalanceOf("Bob"), token.Simbolo)
	fmt.Printf("    Pool: %.0f %s\n", token.BalanceOf("LiquidityPool"), token.Simbolo)

	// Tentativa inválida
	fmt.Println("\n  --- Tentativa de fraude ---")
	err := token.TransferFrom("DEX_Uniswap", "Bob", "Hacker", 999_999)
	if err != nil {
		fmt.Printf("  ❌ Bloqueado: %s\n", err)
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ERC-721: NFT ━━━")

	nft := NovoERC721("CryptoGophers", "CGPH")
	fmt.Printf("\n  Coleção: %s (%s)\n\n", nft.Nome, nft.Simbolo)

	// Mint (criar) NFTs
	nft.Mint("Alice", NFTMetadata{
		Nome:      "Gopher Dourado",
		Descricao: "Um gopher extremamente raro",
		Imagem:    "ipfs://Qm.../1.png",
		Atributos: map[string]string{"Raridade": "Lendário", "Cor": "Dourado"},
	})

	nft.Mint("Alice", NFTMetadata{
		Nome:      "Gopher Azul",
		Descricao: "Um gopher comum mas bonito",
		Imagem:    "ipfs://Qm.../2.png",
		Atributos: map[string]string{"Raridade": "Comum", "Cor": "Azul"},
	})

	nft.Mint("Bob", NFTMetadata{
		Nome:      "Gopher Hacker",
		Descricao: "Gopher com laptop e hoodie",
		Imagem:    "ipfs://Qm.../3.png",
		Atributos: map[string]string{"Raridade": "Raro", "Acessório": "Laptop"},
	})

	// Quem é o dono?
	dono, _ := nft.OwnerOf(1)
	fmt.Printf("\n  Dono do NFT #1: %s\n", dono)
	fmt.Printf("  Alice tem %d NFTs\n", nft.BalanceOf("Alice"))
	fmt.Printf("  Bob tem %d NFTs\n", nft.BalanceOf("Bob"))

	// Transferir NFT
	fmt.Println()
	nft.TransferFrom("Alice", "Bob", 2) // Alice vende Gopher Azul para Bob

	fmt.Printf("\n  Após transferência:\n")
	fmt.Printf("  Alice tem %d NFTs\n", nft.BalanceOf("Alice"))
	fmt.Printf("  Bob tem %d NFTs\n", nft.BalanceOf("Bob"))

	// Tentativa inválida
	fmt.Println()
	err = nft.TransferFrom("Alice", "Carol", 3) // Alice NÃO é dona do #3
	if err != nil {
		fmt.Printf("  ❌ Bloqueado: %s\n", err)
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ METADATA DO NFT ━━━")
	meta := nft.Metadata[1]
	fmt.Printf("  NFT #1:\n")
	fmt.Printf("    Nome: %s\n", meta.Nome)
	fmt.Printf("    Descrição: %s\n", meta.Descricao)
	fmt.Printf("    Imagem: %s\n", meta.Imagem)
	fmt.Printf("    Atributos:\n")
	for k, v := range meta.Atributos {
		fmt.Printf("      %s: %s\n", k, v)
	}
}
