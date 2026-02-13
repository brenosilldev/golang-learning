# 05 — Tokens e NFTs

[← Smart Contracts](../04-smart-contracts/README.md) | [Próximo: DeFi →](../06-defi/README.md)

---

## 📖 O que são Tokens?

Tokens são **ativos digitais** que vivem dentro de smart contracts. O contrato define as regras (quem tem, como transferir, quanto existe).

### Os 3 padrões mais importantes

| Padrão | Para que serve | Exemplo real |
|--------|---------------|-------------|
| **ERC-20** | Token fungível (moeda) | USDT, LINK, UNI, DAI |
| **ERC-721** | NFT (único, não-fungível) | Bored Apes, CryptoPunks |
| **ERC-1155** | Multi-token (combina ambos) | Itens de jogos (CS skins) |

### Fungível vs Não-Fungível
```
FUNGÍVEL (ERC-20):
  1 BRL = 1 BRL (são iguais, intercambiáveis)
  1 USDT = 1 USDT

NÃO-FUNGÍVEL (ERC-721):
  Mona Lisa ≠ Starry Night (cada um é único)
  CryptoPunk #7523 ≠ CryptoPunk #3100
```

---

## 💰 ERC-20 — Token Fungível

### Interface obrigatória
```solidity
interface IERC20 {
    function totalSupply() external view returns (uint256);
    function balanceOf(address account) external view returns (uint256);
    function transfer(address to, uint256 amount) external returns (bool);
    function allowance(address owner, address spender) external view returns (uint256);
    function approve(address spender, uint256 amount) external returns (bool);
    function transferFrom(address from, address to, uint256 amount) external returns (bool);
    
    event Transfer(address indexed from, address indexed to, uint256 value);
    event Approval(address indexed owner, address indexed spender, uint256 value);
}
```

### Como funciona na prática
```
1. Deploy do contrato com supply total (ex: 1.000.000 tokens)
2. Todos os tokens vão para o criador
3. criador.transfer(alice, 1000) → Alice recebe 1000 tokens
4. alice.approve(DEX, 500) → DEX pode usar 500 tokens da Alice
5. DEX.transferFrom(alice, bob, 200) → DEX transfere 200 da Alice pro Bob
```

### Por que `approve` + `transferFrom`?
Para que **smart contracts** (DEX, lending, etc.) possam mover seus tokens POR VOCÊ, com sua autorização prévia. Sem isso, DeFi não existiria.

---

## 🎨 ERC-721 — NFT

### Interface obrigatória
```solidity
interface IERC721 {
    function balanceOf(address owner) external view returns (uint256);
    function ownerOf(uint256 tokenId) external view returns (address);
    function transferFrom(address from, address to, uint256 tokenId) external;
    function approve(address to, uint256 tokenId) external;
    
    event Transfer(address indexed from, address indexed to, uint256 indexed tokenId);
}
```

### Diferença chave
```
ERC-20:  balanceOf(alice) → 1000    (quantidade)
ERC-721: ownerOf(tokenId) → alice   (quem é o dono deste ID específico)
```

### Metadata (o que torna cada NFT único)
```json
{
  "name": "CryptoApe #42",
  "description": "Um macaco digital muito raro",
  "image": "ipfs://QmX.../42.png",
  "attributes": [
    {"trait_type": "Pelo", "value": "Dourado"},
    {"trait_type": "Olhos", "value": "Laser"},
    {"trait_type": "Raridade", "value": "Lendário"}
  ]
}
```

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/tokens.go` | Simulação de ERC-20 e ERC-721 em Go puro |
| `exercicios/ex05_tokens.go` | 🏋️ 4 exercícios |

---

[← Smart Contracts](../04-smart-contracts/README.md) | [Próximo: DeFi →](../06-defi/README.md)
