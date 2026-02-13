# 01 — Criptografia: A Base de Tudo

[← Índice](../README.md) | [Próximo: Blockchain →](../02-blockchain-do-zero/README.md)

---

## 📖 Por que criptografia?

Sem criptografia, não existe blockchain. **Tudo** depende de 3 conceitos:

```
1. HASHING     → SHA-256   → "impressão digital" de dados
2. ASSINATURA  → ECDSA     → provar que FOI VOCÊ que autorizou
3. MERKLE TREE → Árvore    → verificar dados eficientemente
```

---

## 🔐 1. Hashing (SHA-256)

### O que é?
Uma função que transforma **qualquer dado** em um código fixo de 64 caracteres:

```
"Olá"                              → "e0aa02189b..."
"Olá!"                             → "7b502c3a..." (completamente diferente!)
O texto inteiro da Bíblia          → "a1b2c3d4..." (mesmo tamanho: 64 chars)
```

### Propriedades (decore isso!)

| Propriedade | Significado | Por que importa? |
|-------------|------------|------------------|
| **Determinístico** | Mesmo input = mesmo hash, sempre | Qualquer nó pode verificar |
| **Rápido** | Calcular leva nanossegundos | Processar milhões de transações |
| **Irreversível** | Impossível ir do hash ao dado | Não expõe dados originais |
| **Avalanche** | 1 bit muda = hash 100% diferente | Detecta qualquer adulteração |
| **Único** | 2 inputs diferentes → 2 hashes diferentes | Não há "colisão" prática |

### Onde é usado na blockchain?
- **Hash do bloco**: identifica cada bloco unicamente
- **Merkle root**: resumo de todas as transações
- **Mining**: encontrar hash que começa com zeros
- **Endereços**: derivados do hash da chave pública

---

## 🔑 2. Assinatura Digital (ECDSA)

### O que é?
Um sistema que permite **provar que você autorizou algo** sem revelar sua senha.

### Analogia
```
Chave Privada = sua assinatura manuscrita (secreta!)
Chave Pública = foto da sua assinatura no banco (todos podem verificar)

1. Você ASSINA o cheque (chave privada)
2. O banco VERIFICA a assinatura (chave pública)
3. Ninguém pode forjar sua assinatura
4. Se adultrarem o valor do cheque, a assinatura fica inválida
```

### Como funciona?
```
Chave Privada (256 bits aleatórios)
        │
        ▼ (operação matemática unidirecional)
Chave Pública (ponto na curva elíptica)
        │
        ▼ (hash)
Endereço (0x742d35Cc...)
```

### O fluxo de uma transação
```
1. Alice quer enviar 5 ETH para Bob
2. Alice cria a transação: {de: Alice, para: Bob, valor: 5}
3. Alice calcula o HASH da transação
4. Alice ASSINA o hash com sua chave PRIVADA → gera (r, s)
5. Alice envia: transação + assinatura (r, s)
6. Qualquer nó VERIFICA: assinatura + chave PÚBLICA de Alice → true/false
7. Se válida → transação aceita. Se não → rejeitada.
```

---

## 🌳 3. Merkle Tree

### O que é?
Uma árvore de hashes que permite verificar se um dado existe em um conjunto, sem precisar de TODOS os dados.

### Visualização
```
          Root Hash (Merkle Root)
         /                       \
    Hash(AB)                   Hash(CD)
    /      \                   /      \
Hash(A)   Hash(B)         Hash(C)   Hash(D)
   |         |               |         |
  Tx A      Tx B            Tx C      Tx D
```

### Por que isso é genial?

**Sem Merkle Tree**: Para provar que Tx B existe, preciso baixar Tx A, B, C, D → **tudo**
**Com Merkle Tree**: Preciso apenas de Hash(A) + Hash(CD) + Root → **3 hashes** em vez de 4 transações

No Bitcoin com 2000 transações por bloco:
- Sem Merkle: baixar 2000 transações
- Com Merkle: baixar ~11 hashes (log₂ 2000)

### Onde é usado?
- **Blocos**: cada bloco contém o Merkle Root das transações
- **Light clients**: celulares verificam transações sem baixar blockchain inteira
- **Ethereum state**: toda a "memória" do Ethereum é uma Merkle Patricia Trie

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/hashing.go` | SHA-256, efeito avalanche, proof of work |
| `exemplos/wallet.go` | Criar wallet, assinar transação, verificar |
| `exemplos/merkle.go` | Merkle Tree completa com verificação |
| `exercicios/ex01_cripto.go` | 🏋️ 5 exercícios |

### Como rodar:
```bash
go run web3/01-criptografia/exemplos/hashing.go
go run web3/01-criptografia/exemplos/wallet.go
go run web3/01-criptografia/exemplos/merkle.go
```

---

[← Índice](../README.md) | [Próximo: Blockchain →](../02-blockchain-do-zero/README.md)
