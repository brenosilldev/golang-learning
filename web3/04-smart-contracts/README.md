# 04 — Smart Contracts

[← Ethereum e Go](../03-ethereum/README.md) | [Próximo: Tokens e NFTs →](../05-tokens-nfts/README.md)

---

## 📖 O que é um Smart Contract?

É um **programa que roda na blockchain**. Uma vez publicado, ninguém pode alterá-lo.

### Analogia
```
Contrato normal: escrito em papel, precisa de advogado para garantir
Smart Contract:  escrito em código, a BLOCKCHAIN garante execução automática
```

### Características
- ✅ **Imutável** — uma vez na blockchain, não muda (a não ser que tenha upgrade pattern)
- ✅ **Transparente** — qualquer pessoa pode ler o código
- ✅ **Automático** — executa sem intermediário
- ✅ **Determinístico** — mesma entrada = mesma saída, sempre
- ⚠️ **Custa gas** — cada operação custa ETH

---

## 📝 Solidity — A linguagem dos contratos

Solidity é a linguagem mais usada para escrever smart contracts no Ethereum.

### Contrato básico
```solidity
// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

contract Cofre {
    // Variáveis de estado (armazenadas na blockchain)
    address public dono;
    mapping(address => uint256) public saldos;
    
    // Eventos (logs que Go pode escutar)
    event Deposito(address indexed de, uint256 valor);
    event Saque(address indexed para, uint256 valor);
    
    // Constructor (roda UMA vez, no deploy)
    constructor() {
        dono = msg.sender; // quem fez deploy é o dono
    }
    
    // Função para depositar ETH
    function depositar() public payable {
        saldos[msg.sender] += msg.value;
        emit Deposito(msg.sender, msg.value);
    }
    
    // Função para sacar
    function sacar(uint256 valor) public {
        require(saldos[msg.sender] >= valor, "Saldo insuficiente");
        saldos[msg.sender] -= valor;
        payable(msg.sender).transfer(valor);
        emit Saque(msg.sender, valor);
    }
}
```

### Conceitos-chave

| Conceito | O que é |
|----------|---------|
| **address** | Tipo especial: endereço Ethereum (20 bytes) |
| **mapping** | HashMap da blockchain (chave → valor) |
| **msg.sender** | Quem chamou a função (endereço do remetente) |
| **msg.value** | Quantos Wei foram enviados junto |
| **payable** | Indica que a função aceita ETH |
| **require** | Se falhar, reverte tudo (nada muda) |
| **event** | Log na blockchain (Go pode escutar!) |
| **emit** | Emite um evento |
| **view** | Função que só lê (não gasta gas) |
| **pure** | Função que só calcula (nem lê estado) |

---

## 🔗 Go ↔ Smart Contract (O fluxo)

```
1. Escrever contrato em Solidity (.sol)
2. Compilar com solc → gera ABI + Bytecode
3. Usar abigen para gerar código Go tipado
4. No Go: chamar funções do contrato como funções normais!
```

### ABI — Application Binary Interface
O ABI é o "manual de instruções" do contrato. Diz ao Go quais funções existem, quais parâmetros aceitam, e o que retornam.

```json
[
  {
    "name": "depositar",
    "type": "function",
    "inputs": [],
    "outputs": [],
    "stateMutability": "payable"
  },
  {
    "name": "saldos",
    "type": "function",
    "inputs": [{"name": "", "type": "address"}],
    "outputs": [{"name": "", "type": "uint256"}],
    "stateMutability": "view"
  }
]
```

### Gerar binding Go
```bash
# Compilar o contrato
solc --abi --bin Cofre.sol -o build/

# Gerar código Go
abigen --abi=build/Cofre.abi --bin=build/Cofre.bin \
       --pkg=cofre --out=cofre.go

# Agora você pode fazer:
#   contrato, _ := cofre.NewCofre(address, client)
#   contrato.Depositar(opts)
#   saldo, _ := contrato.Saldos(nil, meuEndereco)
```

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/Cofre.sol` | Smart contract Solidity (cofre com depósito/saque) |
| `exemplos/interagir.go` | Interagir com o contrato via Go |
| `exercicios/ex04_contratos.go` | 🏋️ 4 exercícios |

---

[← Ethereum e Go](../03-ethereum/README.md) | [Próximo: Tokens e NFTs →](../05-tokens-nfts/README.md)
