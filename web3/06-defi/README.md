# 06 — DeFi (Finanças Descentralizadas)

[← Tokens e NFTs](../05-tokens-nfts/README.md) | [Próximo: Indexer →](../07-indexer/README.md)

---

## 📖 O que é DeFi?

DeFi = **banco sem banco**. Tudo que um banco faz (empréstimo, câmbio, poupança), smart contracts fazem de forma automática e sem intermediários.

```
BANCO TRADICIONAL:                  DeFi:
  Você → Banco → Investimento        Você → Smart Contract → Rendimento
  Banco cobra taxa de 5%             Contrato cobra 0.3%
  Banco funciona 9h-17h              Contrato funciona 24/7/365
  Banco pode negar                   Contrato aceita qualquer pessoa
  Banco pode falir                   Contrato é código imutável
```

---

## 🏦 Os pilares do DeFi

### 1. DEX — Exchange Descentralizada (Uniswap)

Em vez de um livro de ordens (Binance), usa **Automated Market Maker (AMM)**:

```
Pool de Liquidez:
  ┌─────────────────────────────┐
  │  1000 ETH  ↔  2.000.000 DAI │  ← provedores depositaram
  │                              │
  │  Fórmula: x * y = k          │  ← constante!
  │  1000 * 2000000 = 2 bilhões  │
  └─────────────────────────────┘

Se alguém compra 10 ETH:
  ETH no pool: 1000 - 10 = 990
  DAI necessário: k / 990 = 2.020.202 DAI
  Preço pago: 2.020.202 - 2.000.000 = 20.202 DAI por 10 ETH
  Preço por ETH: ~2020 DAI (mais caro que o preço "normal" de 2000!)

→ Isso é "slippage" — quanto mais compra, mais caro fica
→ Por isso pools grandes = menos slippage
```

### 2. Lending — Empréstimo (Aave, Compound)
```
DEPOSITAR (ganhar juros):
  Você deposita 100 ETH → recebe aETH (token que rende juros)
  APY: 3-5% ao ano

EMPRESTAR (pagar juros):
  Você deposita 100 ETH como colateral
  Pode pegar emprestado até 80 USDT (80% do colateral)
  Se ETH cair e colateral ficar < 80% → LIQUIDAÇÃO automática
```

### 3. Flash Loans — Empréstimo instantâneo
```
O conceito mais LOUCO do DeFi:
  1. Você pega emprestado $1.000.000 (sem colateral!)
  2. Usa para arbitragem
  3. Devolve $1.000.000 + taxa
  4. Tudo na MESMA TRANSAÇÃO

Se não devolver? A transação inteira é REVERTIDA.
Como se nunca tivesse acontecido.
```

### 4. Yield Farming — Ganhar rendimento
```
1. Depositar tokens no pool de liquidez
2. Ganhar taxas de trading (0.3% de cada swap)
3. Ganhar tokens de recompensa (governance tokens)
4. Reinvestir → juros compostos automáticos
```

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/amm.go` | Simulação de AMM (Uniswap simplificado) |
| `exercicios/ex06_defi.go` | 🏋️ 4 exercícios |

---

[← Tokens e NFTs](../05-tokens-nfts/README.md) | [Próximo: Indexer →](../07-indexer/README.md)
