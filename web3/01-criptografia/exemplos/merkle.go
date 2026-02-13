package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

// ============================================================================
// CRIPTOGRAFIA — Merkle Tree
// ============================================================================
//
// Merkle Tree permite verificar se um dado existe em um conjunto
// sem precisar de TODOS os dados. É como um índice super eficiente.
//
// Usada em: blocos (verificar transações), light clients, provas de inclusão
//
// Rode com: go run merkle.go
// ============================================================================

// MerkleNode representa um nó da árvore
type MerkleNode struct {
	Esquerda *MerkleNode
	Direita  *MerkleNode
	Hash     string
	Dado     string // só folhas têm dados
}

// NovaFolha cria um nó folha (com dado real)
func NovaFolha(dado string) *MerkleNode {
	hash := hashear(dado)
	return &MerkleNode{Hash: hash, Dado: dado}
}

// NovoNo cria um nó intermediário (hash dos filhos)
func NovoNo(esquerda, direita *MerkleNode) *MerkleNode {
	hash := hashear(esquerda.Hash + direita.Hash)
	return &MerkleNode{
		Esquerda: esquerda,
		Direita:  direita,
		Hash:     hash,
	}
}

// MerkleTree é a árvore completa
type MerkleTree struct {
	Raiz   *MerkleNode
	Folhas []*MerkleNode
}

// ConstruirMerkleTree constrói a árvore a partir de dados
func ConstruirMerkleTree(dados []string) *MerkleTree {
	if len(dados) == 0 {
		return &MerkleTree{}
	}

	// Criar folhas
	var nos []*MerkleNode
	for _, d := range dados {
		nos = append(nos, NovaFolha(d))
	}

	folhas := make([]*MerkleNode, len(nos))
	copy(folhas, nos)

	// Se ímpar, duplicar o último (padrão blockchain)
	if len(nos)%2 != 0 {
		nos = append(nos, nos[len(nos)-1])
	}

	// Construir de baixo para cima
	for len(nos) > 1 {
		var nivel []*MerkleNode
		for i := 0; i < len(nos); i += 2 {
			pai := NovoNo(nos[i], nos[i+1])
			nivel = append(nivel, pai)
		}
		nos = nivel

		// Se o nível ficou ímpar, duplicar
		if len(nos) > 1 && len(nos)%2 != 0 {
			nos = append(nos, nos[len(nos)-1])
		}
	}

	return &MerkleTree{Raiz: nos[0], Folhas: folhas}
}

// MerkleRoot retorna o hash raiz
func (t *MerkleTree) MerkleRoot() string {
	if t.Raiz == nil {
		return ""
	}
	return t.Raiz.Hash
}

// GerarProva gera a prova de inclusão para um dado
// Retorna os hashes necessários para recalcular até a raiz
func (t *MerkleTree) GerarProva(dado string) ([]ProvaItem, bool) {
	hashDado := hashear(dado)

	// Encontrar o índice da folha
	indice := -1
	for i, f := range t.Folhas {
		if f.Hash == hashDado {
			indice = i
			break
		}
	}
	if indice == -1 {
		return nil, false
	}

	// Coletar os hashes irmãos subindo a árvore
	var prova []ProvaItem
	nos := make([]*MerkleNode, len(t.Folhas))
	copy(nos, t.Folhas)

	if len(nos)%2 != 0 {
		nos = append(nos, nos[len(nos)-1])
	}

	idx := indice
	for len(nos) > 1 {
		var nivel []*MerkleNode
		for i := 0; i < len(nos); i += 2 {
			pai := NovoNo(nos[i], nos[i+1])
			nivel = append(nivel, pai)

			// Se nosso índice está neste par, guardar o irmão
			if i == (idx/2)*2 {
				if idx%2 == 0 {
					// Nosso dado está à esquerda, irmão à direita
					prova = append(prova, ProvaItem{Hash: nos[i+1].Hash, Lado: "D"})
				} else {
					// Nosso dado está à direita, irmão à esquerda
					prova = append(prova, ProvaItem{Hash: nos[i].Hash, Lado: "E"})
				}
			}
		}
		nos = nivel
		idx = idx / 2

		if len(nos) > 1 && len(nos)%2 != 0 {
			nos = append(nos, nos[len(nos)-1])
		}
	}

	return prova, true
}

type ProvaItem struct {
	Hash string
	Lado string // "E" = esquerda, "D" = direita
}

// VerificarProva verifica se um dado está na árvore usando a prova
func VerificarProva(dado string, prova []ProvaItem, raizEsperada string) bool {
	hashAtual := hashear(dado)

	for _, p := range prova {
		if p.Lado == "E" {
			hashAtual = hashear(p.Hash + hashAtual)
		} else {
			hashAtual = hashear(hashAtual + p.Hash)
		}
	}

	return hashAtual == raizEsperada
}

func hashear(dados string) string {
	hash := sha256.Sum256([]byte(dados))
	return hex.EncodeToString(hash[:])
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║           MERKLE TREE                     ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Transações de um bloco
	transacoes := []string{
		"Alice→Bob:5ETH",
		"Carol→Dave:3ETH",
		"Eve→Frank:1ETH",
		"Grace→Heidi:2ETH",
	}

	fmt.Println("\n━━━ CONSTRUIR ÁRVORE ━━━")
	fmt.Println("  Transações no bloco:")
	for i, tx := range transacoes {
		fmt.Printf("    TX%d: %s\n", i, tx)
	}

	tree := ConstruirMerkleTree(transacoes)
	fmt.Printf("\n  Merkle Root: %s\n", tree.MerkleRoot()[:32]+"...")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ VISUALIZAÇÃO DA ÁRVORE ━━━")

	fmt.Println("              Root")
	fmt.Printf("          %s...\n", tree.Raiz.Hash[:12])
	fmt.Println("         /              \\")
	fmt.Printf("    %s...    %s...\n",
		tree.Raiz.Esquerda.Hash[:12],
		tree.Raiz.Direita.Hash[:12])
	fmt.Println("    /       \\           /       \\")
	for i, f := range tree.Folhas {
		fmt.Printf("  %s...  ", f.Hash[:8])
		if i%2 == 1 {
			fmt.Println()
		}
	}
	fmt.Println("    |         |           |         |")
	for i, tx := range transacoes {
		fmt.Printf("  %-14s", tx)
		if i%2 == 1 {
			fmt.Println()
		}
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ PROVA DE INCLUSÃO (Merkle Proof) ━━━")
	fmt.Println("  Provar que 'Carol→Dave:3ETH' está no bloco:")

	prova, encontrado := tree.GerarProva("Carol→Dave:3ETH")
	if !encontrado {
		fmt.Println("  TX não encontrada!")
		return
	}

	fmt.Printf("  Prova necessária (%d hashes):\n", len(prova))
	for i, p := range prova {
		fmt.Printf("    %d. [%s] %s...\n", i+1, p.Lado, p.Hash[:24])
	}

	verificado := VerificarProva("Carol→Dave:3ETH", prova, tree.MerkleRoot())
	fmt.Printf("  Verificação: %t ✅\n", verificado)

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ TENTATIVA DE FRAUDE ━━━")

	fraudeVerificada := VerificarProva("Carol→Dave:999ETH", prova, tree.MerkleRoot())
	fmt.Printf("  TX adulterada válida? %t ❌\n", fraudeVerificada)
	fmt.Println("  → Mudar o valor invalida toda a cadeia de hashes!")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ DETECÇÃO DE ADULTERAÇÃO NO BLOCO ━━━")

	// Se alguém mudar UMA transação, o Merkle Root muda
	txOriginais := []string{"A→B:1", "C→D:2", "E→F:3", "G→H:4"}
	txAdulteradas := []string{"A→B:1", "C→D:999", "E→F:3", "G→H:4"} // mudou TX2

	tree1 := ConstruirMerkleTree(txOriginais)
	tree2 := ConstruirMerkleTree(txAdulteradas)

	fmt.Printf("  Root original:   %s...\n", tree1.MerkleRoot()[:24])
	fmt.Printf("  Root adulterada: %s...\n", tree2.MerkleRoot()[:24])
	fmt.Printf("  Iguais? %t → Fraude detectada! 🚨\n", tree1.MerkleRoot() == tree2.MerkleRoot())

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ EFICIÊNCIA ━━━")
	fmt.Println("  Com 1000 transações:")
	fmt.Println("    Sem Merkle: precisaria de 1000 hashes para verificar")
	fmt.Println("    Com Merkle: precisa de ~10 hashes (log₂ 1000)")
	fmt.Println("  Com 1.000.000 transações:")
	fmt.Println("    Sem Merkle: 1.000.000 hashes")
	fmt.Println("    Com Merkle: ~20 hashes 🚀")
}
