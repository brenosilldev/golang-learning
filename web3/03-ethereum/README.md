# 03 — Ethereum e Go

[← Blockchain do Zero](../02-blockchain-do-zero/README.md) | [Próximo: Smart Contracts →](../04-smart-contracts/README.md)

---

## 📖 Por que Ethereum?

Bitcoin é dinheiro digital. Ethereum é um **computador mundial** — você pode rodar programas (smart contracts) nele.

```
Bitcoin  → calculadora   → só transações de valor
Ethereum → computador    → transações + programas (smart contracts)
```

### Conceitos Ethereum

| Conceito | O que é | Analogia |
|----------|---------|---------|
| **EOA** | Externally Owned Account | Sua conta pessoal (controlada por chave privada) |
| **Contract Account** | Smart Contract na blockchain | Programa rodando 24/7, sem servidor |
| **Gas** | Taxa por operação | "Gasolina" que o computador precisa para rodar |
| **Wei** | Menor unidade de ETH | 1 ETH = 10¹⁸ Wei (como centavos) |
| **Gwei** | Unidade de gas price | 1 ETH = 10⁹ Gwei |
| **EVM** | Ethereum Virtual Machine | O "processador" que executa smart contracts |
| **ABI** | Application Binary Interface | "Manual de instruções" de um smart contract |

### Unidades
```
1 ETH = 1.000.000.000 Gwei = 1.000.000.000.000.000.000 Wei

Exemplo de uma transação:
  Gas Limit: 21000 (unidades de gas)
  Gas Price: 20 Gwei
  Custo total: 21000 × 20 Gwei = 420000 Gwei = 0.00042 ETH
```

---

## 🦊 Go-Ethereum (Geth)

Geth é o client OFICIAL do Ethereum, escrito em Go. O pacote `ethclient` permite que seu programa Go interaja com a rede Ethereum.

### Como conectar
```go
// Rede local (Ganache/Hardhat)
client, _ := ethclient.Dial("http://localhost:8545")

// Testnet Sepolia via Alchemy
client, _ := ethclient.Dial("https://eth-sepolia.g.alchemy.com/v2/SUA_API_KEY")

// Mainnet (cuidado, é dinheiro real!)
client, _ := ethclient.Dial("https://eth-mainnet.g.alchemy.com/v2/SUA_API_KEY")
```

### O que podemos fazer?
```
✅ Ler blocos e transações
✅ Consultar saldos
✅ Enviar transações (transferir ETH)
✅ Chamar funções de smart contracts
✅ Escutar eventos em tempo real
✅ Estimar gas
✅ Consultar nonce (contador de transações)
```

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/conectar.go` | Conectar ao node e ler informações básicas |
| `exemplos/saldos.go` | Consultar saldos e converter unidades |
| `exemplos/transacoes.go` | Enviar ETH de uma wallet para outra |
| `exercicios/ex03_ethereum.go` | 🏋️ 5 exercícios |

### Setup para rodar os exemplos:
```bash
# 1. Instalar go-ethereum
go get github.com/ethereum/go-ethereum

# 2. Rodar blockchain local (escolha um):
npx ganache                        # Opção 1: Ganache
npx hardhat node                   # Opção 2: Hardhat
# Ambos criam uma rede Ethereum local em localhost:8545

# 3. Rodar exemplos
go run web3/03-ethereum/exemplos/conectar.go
```

---

[← Blockchain do Zero](../02-blockchain-do-zero/README.md) | [Próximo: Smart Contracts →](../04-smart-contracts/README.md)
