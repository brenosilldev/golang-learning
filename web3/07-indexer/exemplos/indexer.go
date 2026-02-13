package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"strings"
	"time"
)

// ============================================================================
// INDEXER — Indexador de Blockchain Simulado
// ============================================================================
//
// Este exemplo simula um blockchain indexer completo:
//   1. Uma "blockchain" gera blocos com transações
//   2. O indexer escaneia cada bloco
//   3. Armazena em "DB" (map simulando PostgreSQL)
//   4. API permite consultar por endereço, bloco, TX hash
//
// Na vida real, substituiria:
//   - Blockchain simulada → ethclient conectado ao Ethereum
//   - Maps → PostgreSQL com tabelas
//   - API simples → Chi/Gin com endpoints REST
//
// Rode com: go run indexer.go
// ============================================================================

// --- Modelos da Blockchain ---

type Transacao struct {
	Hash      string  `json:"hash"`
	De        string  `json:"de"`
	Para      string  `json:"para"`
	Valor     float64 `json:"valor"`
	Gas       int     `json:"gas"`
	Timestamp int64   `json:"timestamp"`
	BlocoNum  int     `json:"bloco"`
}

type Bloco struct {
	Numero     int         `json:"numero"`
	Hash       string      `json:"hash"`
	Anterior   string      `json:"anterior"`
	Timestamp  int64       `json:"timestamp"`
	Transacoes []Transacao `json:"transacoes"`
	GasUsado   int         `json:"gas_usado"`
}

type Evento struct {
	Contrato string `json:"contrato"`
	Nome     string `json:"nome"`
	Dados    string `json:"dados"`
	BlocoNum int    `json:"bloco"`
	TXHash   string `json:"tx_hash"`
}

// --- Blockchain Simulada ---

type BlockchainSimulada struct {
	Blocos    []Bloco
	Enderecos []string
}

func NovaBlockchainSimulada() *BlockchainSimulada {
	return &BlockchainSimulada{
		Enderecos: []string{
			"0xAlice", "0xBob", "0xCarol", "0xDave", "0xEve",
			"0xUniswap", "0xAave", "0xOpenSea",
		},
	}
}

func (bc *BlockchainSimulada) GerarBlocos(quantidade int) {
	hashAnterior := strings.Repeat("0", 64)

	for i := 0; i < quantidade; i++ {
		numTxs := rand.Intn(5) + 1
		var txs []Transacao
		gasTotal := 0

		for j := 0; j < numTxs; j++ {
			de := bc.Enderecos[rand.Intn(len(bc.Enderecos))]
			para := bc.Enderecos[rand.Intn(len(bc.Enderecos))]
			for para == de {
				para = bc.Enderecos[rand.Intn(len(bc.Enderecos))]
			}
			gas := 21000 + rand.Intn(100000)
			valor := float64(rand.Intn(100)) / 10.0

			txData := fmt.Sprintf("%d|%d|%s|%s|%.4f", i, j, de, para, valor)
			txHash := hashStr(txData)

			txs = append(txs, Transacao{
				Hash:      "0x" + txHash[:16],
				De:        de,
				Para:      para,
				Valor:     valor,
				Gas:       gas,
				Timestamp: time.Now().Add(time.Duration(i) * time.Minute).Unix(),
				BlocoNum:  i,
			})
			gasTotal += gas
		}

		blocoData := fmt.Sprintf("%d|%s|%d", i, hashAnterior, len(txs))
		blocoHash := hashStr(blocoData)

		bloco := Bloco{
			Numero:     i,
			Hash:       "0x" + blocoHash[:16],
			Anterior:   "0x" + hashAnterior[:16],
			Timestamp:  time.Now().Add(time.Duration(i) * time.Minute).Unix(),
			Transacoes: txs,
			GasUsado:   gasTotal,
		}
		bc.Blocos = append(bc.Blocos, bloco)
		hashAnterior = blocoHash
	}
}

// --- Indexer ---

type Indexer struct {
	// "Banco de dados" (na vida real seria PostgreSQL)
	Blocos        map[int]Bloco          // numero → bloco
	Transacoes    map[string]Transacao   // hash → tx
	PorEndereco   map[string][]Transacao // endereço → transações
	Eventos       []Evento
	UltimoBloco   int
	TotalIndexado int
}

func NovoIndexer() *Indexer {
	return &Indexer{
		Blocos:      make(map[int]Bloco),
		Transacoes:  make(map[string]Transacao),
		PorEndereco: make(map[string][]Transacao),
	}
}

// ProcessarBloco indexa um bloco e todas suas transações
func (idx *Indexer) ProcessarBloco(bloco Bloco) {
	// Salvar bloco
	idx.Blocos[bloco.Numero] = bloco
	idx.UltimoBloco = bloco.Numero

	// Indexar cada transação
	for _, tx := range bloco.Transacoes {
		idx.Transacoes[tx.Hash] = tx
		idx.PorEndereco[tx.De] = append(idx.PorEndereco[tx.De], tx)
		idx.PorEndereco[tx.Para] = append(idx.PorEndereco[tx.Para], tx)
		idx.TotalIndexado++
	}
}

// IndexarBlockchain processa todos os blocos
func (idx *Indexer) IndexarBlockchain(bc *BlockchainSimulada) {
	fmt.Printf("  🔄 Indexando %d blocos...\n", len(bc.Blocos))
	inicio := time.Now()

	for _, bloco := range bc.Blocos {
		idx.ProcessarBloco(bloco)
	}

	fmt.Printf("  ✅ Indexado! %d blocos, %d TXs em %v\n",
		len(idx.Blocos), idx.TotalIndexado, time.Since(inicio))
}

// --- API de consulta ---

func (idx *Indexer) BuscarBloco(numero int) (Bloco, bool) {
	bloco, ok := idx.Blocos[numero]
	return bloco, ok
}

func (idx *Indexer) BuscarTX(hash string) (Transacao, bool) {
	tx, ok := idx.Transacoes[hash]
	return tx, ok
}

func (idx *Indexer) HistoricoEndereco(endereco string) []Transacao {
	return idx.PorEndereco[endereco]
}

func (idx *Indexer) SaldoEndereco(endereco string) float64 {
	saldo := 0.0
	for _, tx := range idx.PorEndereco[endereco] {
		if tx.Para == endereco {
			saldo += tx.Valor
		}
		if tx.De == endereco {
			saldo -= tx.Valor
		}
	}
	return saldo
}

func (idx *Indexer) Estatisticas() {
	fmt.Printf("\n  📊 Estatísticas do Indexer:\n")
	fmt.Printf("     Blocos indexados:  %d\n", len(idx.Blocos))
	fmt.Printf("     TXs indexadas:     %d\n", idx.TotalIndexado)
	fmt.Printf("     Endereços únicos:  %d\n", len(idx.PorEndereco))
	fmt.Printf("     Último bloco:      #%d\n", idx.UltimoBloco)

	// Gas total
	gasTotal := 0
	for _, b := range idx.Blocos {
		gasTotal += b.GasUsado
	}
	fmt.Printf("     Gas total:         %d\n", gasTotal)
}

// --- Helpers ---

func hashStr(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║       BLOCKCHAIN INDEXER                  ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// 1. Gerar blockchain simulada com 20 blocos
	bc := NovaBlockchainSimulada()
	bc.GerarBlocos(20)

	fmt.Println("\n━━━ BLOCKCHAIN GERADA ━━━")
	fmt.Printf("  %d blocos criados\n", len(bc.Blocos))

	// 2. Criar indexer e processar todos os blocos
	idx := NovoIndexer()

	fmt.Println("\n━━━ INDEXAÇÃO ━━━")
	idx.IndexarBlockchain(bc)
	idx.Estatisticas()

	// 3. Consultas (como uma API faria)
	fmt.Println("\n━━━ CONSULTA: BLOCO #5 ━━━")
	if bloco, ok := idx.BuscarBloco(5); ok {
		fmt.Printf("  Bloco #%d\n", bloco.Numero)
		fmt.Printf("  Hash: %s\n", bloco.Hash)
		fmt.Printf("  TXs: %d\n", len(bloco.Transacoes))
		fmt.Printf("  Gas: %d\n", bloco.GasUsado)
	}

	fmt.Println("\n━━━ CONSULTA: HISTÓRICO DE ALICE ━━━")
	aliceTxs := idx.HistoricoEndereco("0xAlice")
	fmt.Printf("  0xAlice tem %d transações\n", len(aliceTxs))
	for i, tx := range aliceTxs {
		tipo := "📤 enviou"
		if tx.Para == "0xAlice" {
			tipo = "📥 recebeu"
		}
		fmt.Printf("  %d. %s %.2f ETH (bloco #%d)\n", i+1, tipo, tx.Valor, tx.BlocoNum)
		if i >= 4 {
			fmt.Printf("  ... e mais %d\n", len(aliceTxs)-5)
			break
		}
	}

	fmt.Println("\n━━━ CONSULTA: SALDOS ━━━")
	enderecos := []string{"0xAlice", "0xBob", "0xCarol", "0xUniswap"}
	for _, e := range enderecos {
		saldo := idx.SaldoEndereco(e)
		numTxs := len(idx.HistoricoEndereco(e))
		fmt.Printf("  %s: %.2f ETH (%d txs)\n", e, saldo, numTxs)
	}

	// 4. Simular "novo bloco" chegando (polling)
	fmt.Println("\n━━━ SIMULAÇÃO: NOVO BLOCO ━━━")
	bc.GerarBlocos(1) // gerar mais 1 bloco
	novoBloco := bc.Blocos[len(bc.Blocos)-1]
	novoBloco.Numero = idx.UltimoBloco + 1
	for i := range novoBloco.Transacoes {
		novoBloco.Transacoes[i].BlocoNum = novoBloco.Numero
	}

	idx.ProcessarBloco(novoBloco)
	fmt.Printf("  📦 Novo bloco #%d indexado! (%d txs)\n",
		novoBloco.Numero, len(novoBloco.Transacoes))

	// 5. Export JSON (como API retornaria)
	fmt.Println("\n━━━ EXPORT: JSON DO BLOCO ━━━")
	if b, ok := idx.BuscarBloco(0); ok {
		jsonBytes, _ := json.MarshalIndent(b, "  ", "  ")
		fmt.Printf("  %s\n", string(jsonBytes))
	}
}
