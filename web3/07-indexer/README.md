# 07 — Indexer e Explorer

[← DeFi](../06-defi/README.md) | [Projeto Final →](../projeto-final/README.md)

---

## 📖 O que é um Indexer?

A blockchain é lenta para consultas. Ler diretamente dela é como procurar um livro em uma biblioteca sem catálogo.

```
DIRETO NA BLOCKCHAIN:
  "Quais transações a Alice fez?" → precisa varrer TODOS os blocos ❌ lento

COM INDEXER:
  Blockchain → Indexer → PostgreSQL → API
  "Quais transações a Alice fez?" → SELECT * FROM txs WHERE de='Alice' ✅ rápido
```

### Todo projeto Web3 sério precisa de um Indexer

| Produto | O que indexa |
|---------|-------------|
| **Etherscan** | Blocos, TXs, contratos de toda a rede |
| **The Graph** | Eventos de smart contracts específicos |
| **Dune Analytics** | Dados DeFi para dashboards |
| **OpenSea** | Metadata e transferências de NFTs |

---

## 🏗️ Arquitetura de um Indexer

```
┌──────────────┐     ┌──────────────┐     ┌──────────────┐
│  Ethereum     │     │   Indexer    │     │  PostgreSQL   │
│  Node         │────▶│  (Go)        │────▶│  (dados)     │
│               │     │              │     │              │
│  Blocos       │     │  • Escutar   │     │  • blocks    │
│  Transações   │     │  • Decodificar│     │  • txs       │
│  Events       │     │  • Salvar    │     │  • events    │
└──────────────┘     └──────────────┘     └──────────────┘
                                                │
                                          ┌─────┴─────┐
                                          │  API REST  │
                                          │  (Go)      │
                                          │            │
                                          │ GET /blocks│
                                          │ GET /txs   │
                                          │ GET /addr  │
                                          └────────────┘
```

### Dois modos de operação

| Modo | Quando usar |
|------|-------------|
| **Polling** | Verificar novos blocos a cada X segundos |
| **WebSocket** | Receber notificação em tempo real quando bloco é criado |

### Polling
```go
for {
    block, _ := client.BlockByNumber(ctx, nil) // último bloco
    if block.Number() > ultimoProcessado {
        processarBloco(block)
        ultimoProcessado = block.Number()
    }
    time.Sleep(2 * time.Second)
}
```

### WebSocket (tempo real)
```go
// Subscrever a novos blocos via WebSocket
headers := make(chan *types.Header)
sub, _ := client.SubscribeNewHead(ctx, headers)

for header := range headers {
    block, _ := client.BlockByNumber(ctx, header.Number)
    processarBloco(block)
}
```

---

## 📜 Decodificando Eventos (Logs)

Eventos são "mensagens" que smart contracts emitem. São o principal dado que indexers processam.

```go
// Filtrar logs de um contrato específico
query := ethereum.FilterQuery{
    FromBlock: big.NewInt(0),
    Addresses: []common.Address{contratoAddress},
}

logs, _ := client.FilterLogs(ctx, query)

for _, vLog := range logs {
    // vLog.Topics[0] = hash do evento (ex: Transfer, Swap)
    // vLog.Data = dados encodados
    // vLog.BlockNumber = bloco onde ocorreu
}
```

### Topics
```
Transfer(address from, address to, uint256 value)
  Topics[0] = keccak256("Transfer(address,address,uint256)") = 0xddf2...
  Topics[1] = from (indexed)
  Topics[2] = to (indexed)
  Data = value (not indexed)
```

---

## 📂 Arquivos

| Arquivo | O que faz |
|---------|----------|
| `exemplos/indexer.go` | Indexer de blockchain simulado em Go puro |
| `exercicios/ex07_indexer.go` | 🏋️ 4 exercícios |

---

[← DeFi](../06-defi/README.md) | [Projeto Final →](../projeto-final/README.md)
