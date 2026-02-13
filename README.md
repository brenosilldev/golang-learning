# 🚀 Go Learning Hub

> Repositório completo de estudo de **Go (Golang)** — da linguagem base até blockchain e Web3.

---

## 📚 O que tem aqui?

Este repositório contém **dois cursos independentes**, organizados para um aprendizado progressivo:

```
golang-learning/
├── 📘 curso/     → Curso de Go (linguagem)
└── 🌐 web3/      → Curso de Web3 com Go (blockchain)
```

---

## 📘 [Curso de Go — Do Básico ao Avançado](./curso/README.md)

Aprenda a linguagem Go do zero até nível profissional. **17 módulos** com teoria, exemplos executáveis e exercícios.

| Trilha | Módulos | O que cobre |
|--------|---------|-------------|
| 🟢 **Fundamentos** | 01-06 | Variáveis, controle de fluxo, coleções, funções, ponteiros |
| 🟡 **Intermediário** | 07-10 | Structs, interfaces, erros, pacotes |
| 🔴 **Avançado** | 11-13 | Concorrência, generics, testes |
| 🏗️ **Produção** | 14-17 | APIs (net/http, Gin, Fiber), gRPC, banco de dados, Docker |
| 🎯 **Projeto Final** | NexusMQ | Message broker distribuído (protocolo TCP, WAL, replicação) |

```bash
# Rodar qualquer exemplo do curso Go:
go run curso/modulo-01-introducao/exemplos/exemplo01_hello.go
```

**→ [Começar o Curso de Go](./curso/README.md)**

---

## 🌐 [Curso de Web3 — Blockchain com Go](./web3/README.md)

Aprenda blockchain, Ethereum, DeFi e como construir infraestrutura Web3 com Go. **7 módulos** com teoria detalhada, simulações funcionais e exercícios práticos.

| Trilha | Módulos | O que cobre |
|--------|---------|-------------|
| 🔐 **Fundamentos** | 01-02 | Criptografia (SHA-256, ECDSA, Merkle Tree), blockchain do zero |
| ⛓️ **Ethereum** | 03-04 | go-ethereum (ethclient), smart contracts + Solidity |
| 💰 **Aplicações** | 05-06 | Tokens (ERC-20/721), DeFi (AMM, lending, flash loans) |
| 🔍 **Infraestrutura** | 07 | Indexer de blockchain (scanner, API, eventos) |
| 🎯 **Projeto Final** | ChainPulse | Blockchain explorer com indexing em tempo real |

```bash
# Rodar exemplos Web3 (vários funcionam sem dependências externas):
go run web3/01-criptografia/exemplos/hashing.go
go run web3/02-blockchain-do-zero/exemplos/blockchain.go
go run web3/05-tokens-nfts/exemplos/tokens.go
go run web3/06-defi/exemplos/amm.go
```

**→ [Começar o Curso de Web3](./web3/README.md)**

---

## 🗺️ Ordem recomendada

```
                    ┌──────────────────────────┐
                    │  📘 Curso Go (módulos 01-17)  │
                    │  Aprenda a linguagem primeiro  │
                    └─────────────┬────────────┘
                                  │
                    ┌─────────────▼────────────┐
                    │  🌐 Curso Web3 (módulos 01-07)│
                    │  Aplique Go em blockchain     │
                    └─────────────┬────────────┘
                                  │
                ┌─────────────────┼─────────────────┐
                ▼                 ▼                  ▼
        🎯 NexusMQ         🎯 ChainPulse        💼 Mercado
        (Message Broker)   (Block Explorer)      (Web3 ou Backend)
```

1. **Primeiro**: Complete o curso Go (módulos 01-17)
2. **Depois**: Faça o curso Web3 se quiser entrar nessa área
3. **Projetos**: Escolha NexusMQ (backend puro) ou ChainPulse (Web3), ou ambos

---

## ⚡ Quick Start

```bash
# Clonar
git clone https://github.com/brenosilldev/golang-learning.git
cd golang-learning

# Primeiro exemplo
go run curso/modulo-01-introducao/exemplos/exemplo01_hello.go
```

---

## 📊 Números

| Métrica | Curso Go | Curso Web3 | Total |
|---------|----------|-----------|-------|
| Módulos | 17 | 7 | 24 |
| READMEs com teoria | 17 | 8 | 25 |
| Exemplos `.go` | 30+ | 12 | 42+ |
| Exercícios | 14 | 7 | 21 |
| Projetos finais | 1 (NexusMQ) | 1 (ChainPulse) | 2 |

---

> 💡 Tudo em Português 🇧🇷 | Go 1.22+ | Nenhum conhecimento prévio de Go necessário
