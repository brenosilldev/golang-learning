# 🎯 Projeto Final — ChainPulse

[← Indexer e Explorer](../07-indexer/README.md) | [Voltar ao Índice](../README.md)

---

## O que é ChainPulse?

Um **Blockchain Explorer + Indexer em tempo real** — como um mini Etherscan que você constrói do zero.

```
┌─────────────┐     ┌─────────────┐     ┌─────────────┐     ┌─────────────┐
│  Ethereum    │     │  Indexer    │     │  PostgreSQL  │     │  API REST   │
│  Node        │────▶│  Engine     │────▶│  Database    │────▶│  + WebSocket│
│  (Ganache)   │     │  (Go)       │     │              │     │  (Go)       │
└─────────────┘     └─────────────┘     └─────────────┘     └─────────────┘
      │                    │                                       │
      │              ┌─────┴──────┐                          ┌─────┴─────┐
      │              │ Decodificar │                          │ Dashboard  │
      │              │ • Blocos    │                          │ Terminal   │
      │              │ • TXs      │                          │ (TUI)      │
      │              │ • Eventos  │                          └───────────┘
      │              │ • ERC-20   │
      │              └────────────┘
      │
      └──▶ WebSocket (novos blocos em tempo real)
```

---

## 🏗️ Funcionalidades (3 fases)

### Fase 1 — Core Indexer (1-2 semanas)
- [ ] Conectar ao nó Ethereum (Ganache) via ethclient
- [ ] Indexar blocos historicamente (bloco 0 até o atual)
- [ ] Armazenar em PostgreSQL: blocos, transações, endereços
- [ ] API REST básica:
  - `GET /api/blocks` — últimos blocos
  - `GET /api/blocks/:number` — detalhes do bloco
  - `GET /api/txs/:hash` — detalhes da transação
  - `GET /api/address/:addr` — saldo + histórico
  - `GET /api/stats` — estatísticas da rede
- [ ] Polling para novos blocos (a cada 2s)

### Fase 2 — Smart Contract Awareness (1-2 semanas)
- [ ] Detectar interações com smart contracts (tx.To = contrato)
- [ ] Decodificar eventos ERC-20 Transfer
- [ ] Indexar tokens: quem tem, quanto, transferências
- [ ] Detectar e indexar ERC-721 (NFTs)
- [ ] Novos endpoints:
  - `GET /api/tokens` — listar tokens encontrados
  - `GET /api/tokens/:address/holders` — holders do token
  - `GET /api/address/:addr/tokens` — tokens de um endereço

### Fase 3 — Real-time + Dashboard (1-2 semanas)
- [ ] WebSocket para receber blocos em tempo real
- [ ] Notificações WebSocket para clientes (`ws://localhost:8080/ws`)
- [ ] Dashboard TUI no terminal mostrando:
  - Blocos chegando em tempo real
  - Top endereços por saldo
  - Transações mais caras
  - Gas por bloco
- [ ] Métricas Prometheus (`/metrics`)
- [ ] Rate limiting na API
- [ ] Dockerizar o projeto completo (API + PostgreSQL + Prometheus)

---

## 📁 Estrutura sugerida

```
chainpulse/
├── cmd/
│   ├── indexer/main.go       # Indexer standalone
│   ├── api/main.go           # API server
│   └── dashboard/main.go     # TUI dashboard
├── internal/
│   ├── indexer/
│   │   ├── indexer.go        # Core: escanear blocos
│   │   ├── decoder.go        # Decodificar TXs e eventos
│   │   └── tokens.go         # Detectar ERC-20/721
│   ├── storage/
│   │   ├── postgres.go       # Queries PostgreSQL
│   │   ├── models.go         # Structs do DB
│   │   └── migrations/       # SQL migrations
│   ├── api/
│   │   ├── server.go         # HTTP server + routes
│   │   ├── handlers.go       # Handler functions
│   │   ├── websocket.go      # WebSocket real-time
│   │   └── middleware.go     # Logging, CORS, rate limit
│   └── dashboard/
│       └── tui.go            # Terminal dashboard
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── go.mod
└── README.md
```

---

## 🛠️ Stack técnica

| Componente | Tecnologia |
|-----------|-----------|
| Linguagem | Go |
| Blockchain | go-ethereum (ethclient) |
| Banco | PostgreSQL + pgx |
| API | net/http (stdlib) |
| WebSocket | gorilla/websocket |
| Dashboard | bubbletea (TUI) |
| Métricas | prometheus/client_golang |
| Container | Docker + docker-compose |

---

## 📋 Conceitos Go aplicados

| Conceito do curso | Como é usado |
|-------------------|-------------|
| Goroutines | Indexer, API e Dashboard em paralelo |
| Channels | Indexer → Dashboard (novos blocos) |
| Context | Cancelamento e timeouts em queries |
| Interfaces | Storage interface (PostgreSQL, SQLite, memória) |
| Generics | Funções utilitárias de conversão |
| Error handling | Erros tipados para cada camada |
| Testing | Testes com banco em memória |
| Packages | Organização em internal/ |
| Structs/JSON | Models e serialização da API |
| Docker | Deploy containerizado |

---

## 🚀 Como começar

```bash
# 1. Criar o projeto
mkdir chainpulse && cd chainpulse
go mod init chainpulse

# 2. Instalar dependências
go get github.com/ethereum/go-ethereum
go get github.com/jackc/pgx/v5
go get github.com/gorilla/websocket

# 3. Subir Ganache + PostgreSQL
docker compose up -d

# 4. Implementar na ordem:
#    a) internal/storage → models + queries
#    b) internal/indexer → ler blocos do Ganache
#    c) cmd/indexer → main que indexa
#    d) internal/api → endpoints REST
#    e) cmd/api → main do server
#    f) Fase 2 e 3...

# 5. Testar
curl http://localhost:8080/api/stats
curl http://localhost:8080/api/blocks
curl http://localhost:8080/api/address/0x...
```

---

## ✅ Critérios de sucesso

- [ ] Indexa 100+ blocos sem erro
- [ ] API responde em < 50ms para queries indexadas
- [ ] WebSocket entrega blocos em < 1s
- [ ] Detecta e indexa pelo menos 1 token ERC-20
- [ ] Dashboard mostra dados em tempo real
- [ ] Projeto roda com `docker compose up`

---

> 🎓 Completar este projeto demonstra habilidade em: Go, concorrência, banco de dados, APIs, WebSocket, blockchain, Docker — tudo que um Go Web3 developer precisa.

[← Indexer e Explorer](../07-indexer/README.md) | [Voltar ao Índice](../README.md)
