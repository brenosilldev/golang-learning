# 🌐 Curso Web 3.0 com Go — Do Zero ao Deploy

> **Pré-requisito**: Ter completado o curso Go (módulos 01-17) ou ter experiência equivalente  
> **Foco**: Blockchain, Ethereum, Smart Contracts, DeFi — tudo construído e integrado com Go

---

## 🤔 O que é Web 3.0?

```
Web 1.0 → Ler           → Sites estáticos. Você consome.
Web 2.0 → Ler + Escrever → Apps e redes sociais. Empresas controlam seus dados.
Web 3.0 → Ler + Escrever + POSSUIR → Blockchain. VOCÊ controla seus dados e ativos.
```

**Na prática**: Web3 é um ecossistema de tecnologias onde dados, dinheiro e identidade são controlados pelo USUÁRIO, não por empresas. A blockchain é o "banco de dados" público e imutável que torna isso possível.

---

## 🗺️ Trilha do Curso

### 🔐 Fundamentos
| # | Módulo | O que aprende |
|---|--------|--------------|
| 01 | [Criptografia](./01-criptografia/README.md) | SHA-256, ECDSA, Merkle Tree — as 3 bases da blockchain |
| 02 | [Blockchain do Zero](./02-blockchain-do-zero/README.md) | Construir uma blockchain completa em Go puro |

### ⛓️ Ethereum
| # | Módulo | O que aprende |
|---|--------|--------------|
| 03 | [Ethereum e Go](./03-ethereum/README.md) | Conectar, ler blocos, consultar saldos, enviar TX |
| 04 | [Smart Contracts](./04-smart-contracts/README.md) | Solidity básico, ABI, interagir via Go, eventos |

### 💰 Aplicações
| # | Módulo | O que aprende |
|---|--------|--------------|
| 05 | [Tokens e NFTs](./05-tokens-nfts/README.md) | ERC-20, ERC-721, criar e interagir via Go |
| 06 | [DeFi](./06-defi/README.md) | DEX, AMM, liquidity pools, flash loans |

### 🔍 Infraestrutura
| # | Módulo | O que aprende |
|---|--------|--------------|
| 07 | [Indexer e Explorer](./07-indexer/README.md) | Escanear blocos, indexar eventos, WebSocket |

### 🎯 Projeto
| # | Módulo | O que aprende |
|---|--------|--------------|
| 🔥 | [Projeto Final](./projeto-final/README.md) | ChainPulse — Blockchain Explorer com indexing em tempo real |

---

## 🧰 Setup do Ambiente

```bash
# go-ethereum (interagir com Ethereum)
go get github.com/ethereum/go-ethereum

# Solidity compiler (smart contracts)
sudo add-apt-repository ppa:ethereum/ethereum
sudo apt update && sudo apt install solc

# abigen (gerar bindings Go de contratos)
go install github.com/ethereum/go-ethereum/cmd/abigen@latest

# Ganache (blockchain local para testes)
npm install -g ganache
# Ou use Hardhat: npx hardhat node
```

---

## 💡 Dica importante

Todos os exemplos deste curso rodam **localmente** ou contra **testnets gratuitas**. Você **NÃO precisa gastar dinheiro real** para aprender. Use:

| Rede | Para que | Como acessar |
|------|---------|-------------|
| **Ganache** | Blockchain local (instantânea) | `ganache` no terminal |
| **Hardhat** | Blockchain local + debug | `npx hardhat node` |
| **Sepolia** | Testnet pública Ethereum | Registre-se no [Alchemy](https://alchemy.com) ou [Infura](https://infura.io) |

---

> 🚀 Comece pelo módulo 01 e siga em ordem. Cada módulo depende do anterior.
