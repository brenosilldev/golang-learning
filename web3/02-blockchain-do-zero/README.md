# 02 — Blockchain do Zero

[← Criptografia](../01-criptografia/README.md) | [Próximo: Ethereum →](../03-ethereum/README.md)

---

## 📖 O que vamos construir?

Uma **blockchain funcional** em Go puro — sem bibliotecas externas. Blocos, encadeamento, Proof of Work, validação, e persistência.

---

## ⛓️ Anatomia de um Bloco

```
┌──────────────────────────────────────────┐
│ Bloco #3                                 │
├──────────────────────────────────────────┤
│ Timestamp:      1705334400               │
│ Hash Anterior:  0000a8f3c2...            │ ← link para Bloco #2
│ Transações:                              │
│   • Alice → Bob: 5.0                     │
│   • Carol → Dave: 2.0                    │
│ Merkle Root:    7e91b4...                │
│ Nonce:          84721                    │ ← encontrado por mining
│ Dificuldade:    4                        │
│ Hash:           0000d7e1...              │ ← começa com "0000" ✅
└──────────────────────────────────────────┘
```

## 🔗 O que torna "chain"?

Cada bloco contém o **hash do anterior**. Se alguém mudar um bloco antigo, TODOS os hashes seguintes ficam inválidos:

```
Bloco 0 ──hash──► Bloco 1 ──hash──► Bloco 2 ──hash──► Bloco 3
  │                  │                  │                  │
  └─ válido          └─ válido          └─ ADULTERADO!     └─ hash anterior
                                           O hash mudou!       não bate! ❌
```

## ⛏️ Proof of Work (Consenso)

Como a rede decide quem adiciona o próximo bloco?

```
1. Mineiro junta transações pendentes
2. Tenta nonces: 0, 1, 2, ... (brute force)
3. Calcula: SHA256(bloco + nonce)
4. Se começa com "0000..." → ACHEI! 🎉
5. Transmite o bloco para a rede
6. Outros nós verificam (instantâneo) e aceitam
```

**Dificuldade vs Tempo:**
```
1 zero  ("0...")    → microsegundos
2 zeros ("00...")   → milissegundos  
4 zeros ("0000...") → segundos
6 zeros             → minutos
Bitcoin (~19 zeros) → ~10 minutos (com farms inteiras)
```

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/blockchain.go` | Blockchain completa: blocos, mining, validação |
| `exercicios/ex02_blockchain.go` | 🏋️ 5 exercícios |

```bash
go run web3/02-blockchain-do-zero/exemplos/blockchain.go
```

---

[← Criptografia](../01-criptografia/README.md) | [Próximo: Ethereum →](../03-ethereum/README.md)
