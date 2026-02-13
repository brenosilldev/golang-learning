package main

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"
)

// ============================================================================
// CRIPTOGRAFIA — Wallets e Assinatura Digital (ECDSA)
// ============================================================================
//
// Uma "wallet" NÃO guarda dinheiro. Ela guarda suas CHAVES:
//   - Chave Privada → sua senha (NUNCA compartilhe!)
//   - Chave Pública → seu endereço público
//
// ECDSA = Elliptic Curve Digital Signature Algorithm
//         (Como Bitcoin e Ethereum assinam transações)
//
// Rode com: go run wallet.go
// ============================================================================

// --- Wallet ---

type Wallet struct {
	privada  *ecdsa.PrivateKey
	Publica  *ecdsa.PublicKey
	Endereco string
}

func NovaWallet() (*Wallet, error) {
	// Gerar par de chaves na curva P-256
	privada, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	if err != nil {
		return nil, err
	}

	// Endereço = hash da chave pública (primeiros 20 bytes, como Ethereum)
	pubBytes := append(privada.PublicKey.X.Bytes(), privada.PublicKey.Y.Bytes()...)
	hash := sha256.Sum256(pubBytes)
	endereco := "0x" + hex.EncodeToString(hash[:20])

	return &Wallet{
		privada:  privada,
		Publica:  &privada.PublicKey,
		Endereco: endereco,
	}, nil
}

// --- Transação ---

type Transacao struct {
	De    string
	Para  string
	Valor float64
}

func (tx Transacao) Hash() []byte {
	dados := fmt.Sprintf("%s|%s|%.8f", tx.De, tx.Para, tx.Valor)
	hash := sha256.Sum256([]byte(dados))
	return hash[:]
}

func (tx Transacao) String() string {
	return fmt.Sprintf("%s... → %s... (%.2f ETH)",
		tx.De[:10], tx.Para[:10], tx.Valor)
}

// --- Assinatura ---

type Assinatura struct {
	R, S *big.Int
}

// Assinar — usa a chave PRIVADA (só o dono pode fazer)
func (w *Wallet) Assinar(tx Transacao) (*Assinatura, error) {
	r, s, err := ecdsa.Sign(rand.Reader, w.privada, tx.Hash())
	if err != nil {
		return nil, err
	}
	return &Assinatura{R: r, S: s}, nil
}

// Verificar — usa a chave PÚBLICA (qualquer pessoa pode fazer)
func Verificar(tx Transacao, assinatura *Assinatura, publica *ecdsa.PublicKey) bool {
	return ecdsa.Verify(publica, tx.Hash(), assinatura.R, assinatura.S)
}

func main() {
	// ══════════════════════════════════════════════════════════
	fmt.Println("╔══════════════════════════════════════════╗")
	fmt.Println("║     WALLETS E ASSINATURA DIGITAL         ║")
	fmt.Println("╚══════════════════════════════════════════╝")

	// Criar wallets para Alice e Bob
	alice, _ := NovaWallet()
	bob, _ := NovaWallet()

	fmt.Printf("\n👩 Alice: %s\n", alice.Endereco)
	fmt.Printf("👨 Bob:   %s\n", bob.Endereco)

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ ASSINAR UMA TRANSAÇÃO ━━━")

	tx := Transacao{
		De:    alice.Endereco,
		Para:  bob.Endereco,
		Valor: 5.0,
	}
	fmt.Printf("  Transação: %s\n", tx)

	// Alice assina com SUA chave privada
	assinatura, _ := alice.Assinar(tx)
	fmt.Printf("  Assinatura (r): %s...\n", assinatura.R.Text(16)[:16])
	fmt.Printf("  Assinatura (s): %s...\n", assinatura.S.Text(16)[:16])

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ VERIFICAR A ASSINATURA ━━━")

	// Qualquer nó da rede verifica com a chave PÚBLICA de Alice
	valido := Verificar(tx, assinatura, alice.Publica)
	fmt.Printf("  Verificar com chave de Alice: %t ✅\n", valido)

	// Se verificar com a chave de Bob → falha (Bob não assinou)
	invalido := Verificar(tx, assinatura, bob.Publica)
	fmt.Printf("  Verificar com chave de Bob:   %t ❌ (não foi ele)\n", invalido)

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ TENTATIVA DE FRAUDE ━━━")

	// Alguém tenta alterar o valor de 5 para 5000 ETH
	txFraudada := Transacao{
		De:    alice.Endereco,
		Para:  bob.Endereco,
		Valor: 5000.0, // ADULTERADO!
	}

	fraudeDetectada := Verificar(txFraudada, assinatura, alice.Publica)
	fmt.Printf("  TX adulterada (5000 ETH) válida? %t ❌\n", fraudeDetectada)
	fmt.Println("  → A assinatura NÃO bate com dados alterados!")

	// Alguém tenta mudar o destinatário
	carol, _ := NovaWallet()
	txDesviada := Transacao{
		De:    alice.Endereco,
		Para:  carol.Endereco, // DESVIADA para Carol!
		Valor: 5.0,
	}

	desvioDetectado := Verificar(txDesviada, assinatura, alice.Publica)
	fmt.Printf("  TX desviada (para Carol) válida? %t ❌\n", desvioDetectado)
	fmt.Println("  → Mudar QUALQUER campo invalida a assinatura!")

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ SIMULAÇÃO: VÁRIAS TRANSAÇÕES ━━━")

	dave, _ := NovaWallet()
	wallets := map[string]*Wallet{
		"Alice": alice,
		"Bob":   bob,
		"Carol": carol,
		"Dave":  dave,
	}
	_ = wallets

	// Simular um bloco com 4 transações
	transacoes := []struct {
		remetente    *Wallet
		nomeRem      string
		destinatario string
		valor        float64
	}{
		{alice, "Alice", bob.Endereco, 2.5},
		{bob, "Bob", carol.Endereco, 1.0},
		{carol, "Carol", dave.Endereco, 0.5},
		{dave, "Dave", alice.Endereco, 3.0},
	}

	fmt.Println("  Processando bloco com 4 transações:")
	for i, t := range transacoes {
		tx := Transacao{De: t.remetente.Endereco, Para: t.destinatario, Valor: t.valor}
		assin, _ := t.remetente.Assinar(tx)
		valido := Verificar(tx, assin, t.remetente.Publica)

		status := "✅ válida"
		if !valido {
			status = "❌ inválida"
		}
		fmt.Printf("    TX%d: %s envia %.1f ETH → %s\n", i+1, t.nomeRem, t.valor, status)
	}

	// ══════════════════════════════════════════════════════════
	fmt.Println("\n━━━ RESUMO ━━━")
	fmt.Println("  🔑 Chave Privada = sua senha secreta")
	fmt.Println("  📢 Chave Pública = seu endereço público")
	fmt.Println("  ✍️  Assinar = provar que FOI VOCÊ")
	fmt.Println("  🔍 Verificar = qualquer um confere sem saber sua senha")
	fmt.Println("  🛡️  Adulterar = impossível sem a chave privada")
}
