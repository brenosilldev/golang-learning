package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// ============================================================================
// BLOCKCHAIN COMPLETA DO ZERO — Em Go puro
// ============================================================================
//
// Tudo que uma blockchain real tem (simplificado):
//   - Blocos com transações
//   - Encadeamento via hashes
//   - Proof of Work (mining)
//   - Validação da cadeia
//   - Bloco genesis
//
// Rode com: go run blockchain.go
// ============================================================================

// --- Transação ---

type Transacao struct {
	De    string  `json:"de"`
	Para  string  `json:"para"`
	Valor float64 `json:"valor"`
}

func (tx Transacao) String() string {
	return fmt.Sprintf("%s→%s:%.2f", tx.De, tx.Para, tx.Valor)
}

// --- Bloco ---

type Bloco struct {
	Indice       int         `json:"indice"`
	Timestamp    int64       `json:"timestamp"`
	Transacoes   []Transacao `json:"transacoes"`
	HashAnterior string      `json:"hash_anterior"`
	Hash         string      `json:"hash"`
	Nonce        int         `json:"nonce"`
	Dificuldade  int         `json:"dificuldade"`
}

// calcularHash gera o hash do bloco baseado em TODO seu conteúdo
func (b *Bloco) calcularHash() string {
	// Serializar transações
	txStr := ""
	for _, tx := range b.Transacoes {
		txStr += tx.String() + "|"
	}

	dados := fmt.Sprintf("%d|%d|%s|%s|%d",
		b.Indice, b.Timestamp, txStr, b.HashAnterior, b.Nonce)

	hash := sha256.Sum256([]byte(dados))
	return hex.EncodeToString(hash[:])
}

// minerar encontra o nonce correto (Proof of Work)
func (b *Bloco) minerar() {
	prefixo := strings.Repeat("0", b.Dificuldade)

	for {
		b.Hash = b.calcularHash()
		if strings.HasPrefix(b.Hash, prefixo) {
			return // Encontrou!
		}
		b.Nonce++
	}
}

// --- Blockchain ---

type Blockchain struct {
	Cadeia      []*Bloco
	Dificuldade int
	Pendentes   []Transacao // pool de transações aguardando
}

// NovaBlockchain cria uma blockchain com o bloco genesis
func NovaBlockchain(dificuldade int) *Blockchain {
	bc := &Blockchain{
		Dificuldade: dificuldade,
		Pendentes:   []Transacao{},
	}

	// Bloco Genesis (o primeiro bloco, sem anterior)
	genesis := &Bloco{
		Indice:       0,
		Timestamp:    time.Now().Unix(),
		Transacoes:   []Transacao{{De: "SISTEMA", Para: "Genesis", Valor: 0}},
		HashAnterior: strings.Repeat("0", 64),
		Dificuldade:  dificuldade,
	}
	genesis.minerar()
	bc.Cadeia = append(bc.Cadeia, genesis)

	return bc
}

// AdicionarTransacao coloca TX na pool de pendentes
func (bc *Blockchain) AdicionarTransacao(de, para string, valor float64) {
	bc.Pendentes = append(bc.Pendentes, Transacao{De: de, Para: para, Valor: valor})
}

// MinerarBloco cria um novo bloco com as transações pendentes
func (bc *Blockchain) MinerarBloco(minerador string) *Bloco {
	ultimo := bc.Cadeia[len(bc.Cadeia)-1]

	// Recompensa do minerador (como Bitcoin!)
	recompensa := Transacao{De: "REDE", Para: minerador, Valor: 6.25}
	transacoes := append(bc.Pendentes, recompensa)

	bloco := &Bloco{
		Indice:       len(bc.Cadeia),
		Timestamp:    time.Now().Unix(),
		Transacoes:   transacoes,
		HashAnterior: ultimo.Hash,
		Dificuldade:  bc.Dificuldade,
	}

	fmt.Printf("⛏️  Minerando bloco #%d (%d transações)...\n", bloco.Indice, len(transacoes))
	inicio := time.Now()
	bloco.minerar()
	duracao := time.Since(inicio)

	bc.Cadeia = append(bc.Cadeia, bloco)
	bc.Pendentes = []Transacao{} // limpar pendentes

	fmt.Printf("✅ Minerado! Nonce: %d | Tempo: %v | Hash: %s...\n",
		bloco.Nonce, duracao, bloco.Hash[:16])

	return bloco
}

// Valida verifica TODA a cadeia
func (bc *Blockchain) Valida() bool {
	for i := 1; i < len(bc.Cadeia); i++ {
		atual := bc.Cadeia[i]
		anterior := bc.Cadeia[i-1]

		// O hash armazenado bate com o recalculado?
		if atual.Hash != atual.calcularHash() {
			fmt.Printf("❌ Bloco #%d: hash não confere!\n", i)
			return false
		}

		// O hash anterior bate com o hash do bloco anterior?
		if atual.HashAnterior != anterior.Hash {
			fmt.Printf("❌ Bloco #%d: hash anterior não confere!\n", i)
			return false
		}

		// O hash começa com a dificuldade correta?
		prefixo := strings.Repeat("0", atual.Dificuldade)
		if !strings.HasPrefix(atual.Hash, prefixo) {
			fmt.Printf("❌ Bloco #%d: dificuldade não confere!\n", i)
			return false
		}
	}
	return true
}

// Saldo calcula o saldo de um endereço
func (bc *Blockchain) Saldo(endereco string) float64 {
	saldo := 0.0
	for _, bloco := range bc.Cadeia {
		for _, tx := range bloco.Transacoes {
			if tx.Para == endereco {
				saldo += tx.Valor
			}
			if tx.De == endereco {
				saldo -= tx.Valor
			}
		}
	}
	return saldo
}

// ImprimirCadeia mostra toda a blockchain formatada
func (bc *Blockchain) ImprimirCadeia() {
	for _, bloco := range bc.Cadeia {
		fmt.Printf("\n  📦 Bloco #%d\n", bloco.Indice)
		fmt.Printf("     Hash:     %s...\n", bloco.Hash[:24])
		fmt.Printf("     Anterior: %s...\n", bloco.HashAnterior[:24])
		fmt.Printf("     Nonce:    %d\n", bloco.Nonce)
		fmt.Printf("     TXs:      %d\n", len(bloco.Transacoes))
		for _, tx := range bloco.Transacoes {
			fmt.Printf("       • %s → %s: %.2f\n", tx.De, tx.Para, tx.Valor)
		}
	}
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║    BLOCKCHAIN COMPLETA DO ZERO           ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Criar blockchain com dificuldade 4 (4 zeros)
	bc := NovaBlockchain(4)
	fmt.Printf("Blockchain criada! Dificuldade: %d\n", bc.Dificuldade)

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ADICIONANDO TRANSAÇÕES ━━━")

	// Bloco 1: transferências iniciais
	bc.AdicionarTransacao("Alice", "Bob", 5.0)
	bc.AdicionarTransacao("Bob", "Carol", 2.0)
	bc.MinerarBloco("Minerador1")

	// Bloco 2: mais transferências
	bc.AdicionarTransacao("Carol", "Dave", 1.5)
	bc.AdicionarTransacao("Alice", "Dave", 3.0)
	bc.AdicionarTransacao("Dave", "Bob", 0.5)
	bc.MinerarBloco("Minerador2")

	// Bloco 3
	bc.AdicionarTransacao("Bob", "Alice", 1.0)
	bc.MinerarBloco("Minerador1")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ BLOCKCHAIN COMPLETA ━━━")
	bc.ImprimirCadeia()

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ VALIDAÇÃO ━━━")
	if bc.Valida() {
		fmt.Println("✅ Blockchain VÁLIDA!")
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ SALDOS ━━━")
	enderecos := []string{"Alice", "Bob", "Carol", "Dave", "Minerador1", "Minerador2"}
	for _, e := range enderecos {
		fmt.Printf("  💰 %s: %.2f\n", e, bc.Saldo(e))
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ TESTE DE ADULTERAÇÃO ━━━")

	// Tentativa de fraude: mudar valor de uma transação
	fmt.Println("  Tentando alterar Bloco #1 (Alice→Bob de 5 para 999)...")
	bc.Cadeia[1].Transacoes[0].Valor = 999.0

	if !bc.Valida() {
		fmt.Println("  🚨 FRAUDE DETECTADA! Blockchain inválida!")
	}

	// Reverter a fraude para continuar
	bc.Cadeia[1].Transacoes[0].Valor = 5.0

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ EXPORTAR COMO JSON ━━━")
	jsonBytes, _ := json.MarshalIndent(bc.Cadeia[1], "", "  ")
	fmt.Printf("  Bloco #1 em JSON:\n%s\n", string(jsonBytes))
}
