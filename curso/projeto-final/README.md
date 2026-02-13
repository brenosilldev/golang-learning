# 🔥 Projeto Final — NexusMQ: Message Broker Distribuído

[← APIs](../modulo-14-apis/README.md) | [Voltar ao índice](../README.md)

---

## ⚡ O Desafio

Construa um **message broker distribuído** do zero, como um mini Apache Kafka. Nada de bibliotecas de message queue — você implementa o protocolo, o storage engine, o roteamento, a replicação. **Tudo.**

Isso é um projeto de nível **Staff/Principal Engineer**. Vai doer. Vai ser incrível.

---

## 🎯 O que NexusMQ faz?

```
Producer ──► NexusMQ Broker ──► Consumer Group A (3 consumers)
                   │
                   └──► Consumer Group B (2 consumers)
```

1. **Producers** publicam mensagens em **topics**
2. **Topics** são divididos em **partitions** (para paralelismo)
3. **Consumer Groups** leem mensagens com coordenação (cada mensagem vai para UM consumer do grupo)
4. **Mensagens são persistidas** em disco com WAL (Write-Ahead Log)
5. **Replicação** entre nós para tolerância a falhas
6. **Protocolo binário customizado** sobre TCP (não HTTP!)

---

## 📐 Arquitetura

```
┌─────────────────────────────────────────────────────┐
│                    NEXUSMQ BROKER                    │
│                                                     │
│  ┌──────────┐  ┌──────────────┐  ┌───────────────┐ │
│  │ TCP      │  │ Protocol     │  │ Topic         │ │
│  │ Server   ├──► Decoder/     ├──► Router         │ │
│  │ (net)    │  │ Encoder      │  │               │ │
│  └──────────┘  └──────────────┘  └───────┬───────┘ │
│                                          │         │
│                    ┌─────────────────────┤         │
│                    │                     │         │
│  ┌─────────────────▼──┐  ┌──────────────▼───────┐ │
│  │ Partition Manager   │  │ Consumer Group       │ │
│  │                     │  │ Coordinator          │ │
│  │ ┌────┐ ┌────┐      │  │                      │ │
│  │ │ P0 │ │ P1 │ ...  │  │ ┌─────┐ ┌─────┐     │ │
│  │ └──┬─┘ └──┬─┘      │  │ │ CG1 │ │ CG2 │     │ │
│  │    │      │         │  │ └─────┘ └─────┘     │ │
│  └────┼──────┼─────────┘  └─────────────────────┘ │
│       │      │                                     │
│  ┌────▼──────▼─────────┐  ┌─────────────────────┐ │
│  │ Storage Engine       │  │ Replication          │ │
│  │ (WAL + Segment Files)│  │ (Leader/Follower)    │ │
│  └──────────────────────┘  └─────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

---

## 📋 Funcionalidades por Fase

### Fase 1 — O Core (Obrigatório, ~2-3 semanas)

**Protocolo binário customizado sobre TCP:**
- [ ] Definir formato de frames: `[length:4][type:1][payload:N]`
- [ ] Tipos: `PRODUCE`, `CONSUME`, `ACK`, `SUBSCRIBE`, `CREATE_TOPIC`, `ERROR`, `HEARTBEAT`
- [ ] Encoder/Decoder eficiente (sem JSON — use `encoding/binary`)
- [ ] Connection pool no client

**Storage Engine:**
- [ ] Write-Ahead Log (WAL) — append-only file por partition
- [ ] Formato de segmento: `[offset:8][timestamp:8][key_len:4][key:N][val_len:4][val:N][crc32:4]`
- [ ] Segmentação: quando o arquivo atinge X MB, cria novo segmento
- [ ] Index file: mapeamento offset → posição no arquivo (para seeks rápidos)
- [ ] Compactação: remover mensagens antigas por retenção (tempo ou tamanho)

**Topic & Partition:**
- [ ] Topic com N partitions (definido na criação)
- [ ] Partitioning por key hash (mensagens com mesma key vão para mesma partition)
- [ ] Round-robin se não tiver key
- [ ] Cada partition é um WAL independente

**Producer:**
- [ ] Client TCP que conecta ao broker
- [ ] Envia mensagens com `topic`, `key` (opcional) e `value`
- [ ] Batching: agrupa mensagens para enviar de uma vez (configurable)
- [ ] Retry com backoff exponencial

**Consumer & Consumer Groups:**
- [ ] Consumers se registram em um Consumer Group
- [ ] Cada partition é atribuída a UM consumer do grupo
- [ ] Rebalancing quando consumer entra/sai
- [ ] Offset tracking (cada consumer sabe onde parou)
- [ ] Commit manual e auto-commit

---

### Fase 2 — Production-Ready (~1-2 semanas)

**Admin API (HTTP):**
- [ ] `POST /topics` — criar topic
- [ ] `GET /topics` — listar topics
- [ ] `GET /topics/:name` — detalhes (partitions, offsets, consumers)
- [ ] `GET /groups` — listar consumer groups
- [ ] `GET /metrics` — métricas (msg/s, bytes/s, lag por grupo)
- [ ] `DELETE /topics/:name` — deletar topic

**CLI (Cobra):**
- [ ] `nexusmq topic create --name orders --partitions 6`
- [ ] `nexusmq topic list`
- [ ] `nexusmq produce --topic orders --key "user-123" --value "order created"`
- [ ] `nexusmq consume --topic orders --group my-service`
- [ ] `nexusmq status` — dashboard no terminal

**Observabilidade:**
- [ ] Métricas Prometheus: mensagens/s, bytes/s, latência, consumer lag
- [ ] Structured logging (slog ou zerolog)
- [ ] Health check endpoint
- [ ] Graceful shutdown com `context.Context` + OS signals

---

### Fase 3 — Distribuído (Desafio Máximo, ~2+ semanas)

**Replicação Leader/Follower:**
- [ ] Broker cluster com N nós
- [ ] Cada partition tem 1 leader e N-1 followers
- [ ] Producers escrevem apenas no leader
- [ ] Leader replica para followers via protocolo interno
- [ ] ISR (In-Sync Replicas) — track de quem está sincronizado
- [ ] Leader election quando leader cai (via heartbeats)

**Exactly-Once Semantics:**
- [ ] Producer ID + Sequence Number para deduplicação
- [ ] Idempotent producer

**Compressão:**
- [ ] Suportar gzip, snappy, lz4 nas mensagens
- [ ] Compressão configurável por topic

---

## 📁 Estrutura do Projeto

```
nexusmq/
├── cmd/
│   ├── broker/
│   │   └── main.go                 # Ponto de entrada do broker
│   ├── cli/
│   │   └── main.go                 # CLI (Cobra)
│   └── benchmark/
│       └── main.go                 # Benchmark tool
│
├── internal/
│   ├── protocol/
│   │   ├── frame.go                # Frame format [len][type][payload]
│   │   ├── encoder.go              # Binary encoding
│   │   ├── decoder.go              # Binary decoding
│   │   ├── messages.go             # ProduceRequest, ConsumeResponse, etc.
│   │   └── protocol_test.go
│   │
│   ├── storage/
│   │   ├── wal.go                  # Write-Ahead Log
│   │   ├── segment.go              # Segment file management
│   │   ├── index.go                # Offset → position index
│   │   ├── compactor.go            # Retenção e limpeza
│   │   └── storage_test.go
│   │
│   ├── broker/
│   │   ├── broker.go               # Broker principal
│   │   ├── topic.go                # Topic management
│   │   ├── partition.go            # Partition (wraps WAL)
│   │   ├── router.go               # Roteamento key → partition
│   │   └── broker_test.go
│   │
│   ├── consumer/
│   │   ├── group.go                # Consumer Group coordination
│   │   ├── assignment.go           # Partition → Consumer assignment
│   │   ├── offset.go               # Offset tracking & commit
│   │   └── group_test.go
│   │
│   ├── network/
│   │   ├── server.go               # TCP server (accept, handle conns)
│   │   ├── connection.go           # Connection wrapper
│   │   ├── handler.go              # Request dispatcher
│   │   └── server_test.go
│   │
│   ├── replication/                # Fase 3
│   │   ├── leader.go
│   │   ├── follower.go
│   │   ├── election.go
│   │   └── isr.go
│   │
│   ├── admin/
│   │   ├── api.go                  # HTTP admin API
│   │   └── metrics.go              # Prometheus metrics
│   │
│   └── config/
│       └── config.go               # Configuração YAML
│
├── pkg/
│   └── client/
│       ├── producer.go             # Producer client library
│       ├── consumer.go             # Consumer client library
│       └── client_test.go
│
├── configs/
│   └── nexusmq.yaml                # Config padrão
│
├── go.mod
├── go.sum
├── Makefile
├── Dockerfile
└── README.md
```

---

## 🔧 Detalhes Técnicos que Você Vai Precisar Resolver

### 1. Protocolo Binário
```
Frame Format:
┌──────────┬──────────┬─────────────────┐
│ Length(4) │ Type(1)  │ Payload(N)      │
│ uint32   │ uint8    │ variable        │
└──────────┴──────────┴─────────────────┘

ProduceRequest:
┌───────────────┬──────────┬───────────┬──────────┬──────────┐
│ TopicLen(2)   │ Topic(N) │ KeyLen(4) │ Key(N)   │ Value(N) │
│ uint16        │ string   │ uint32    │ bytes    │ bytes    │
└───────────────┴──────────┴───────────┴──────────┴──────────┘
```

### 2. Storage (WAL Segment)
```
Message on disk:
┌──────────┬────────────┬──────────┬────────┬──────────┬────────┬───────┐
│ Offset(8)│ Timestamp(8│ KeyLen(4)│ Key(N) │ ValLen(4)│ Val(N) │ CRC(4)│
│ uint64   │ int64 unix │ uint32   │ bytes  │ uint32   │ bytes  │ crc32 │
└──────────┴────────────┴──────────┴────────┴──────────┴────────┴───────┘
```
Use `os.File` com `O_APPEND|O_WRONLY` para writes e `mmap` ou `Seek+Read` para reads.

### 3. Concorrência
- Cada **partition** é uma goroutine com seu own channel
- **Consumer Group Coordinator** é uma goroutine que gerencia assignments
- **TCP Server**: 1 goroutine per connection
- **Replication**: goroutine por follower fazendo streaming
- Use `sync.RWMutex` para metadata, **channels** para data flow

### 4. Consumer Group Rebalancing
```
Quando consumer entra/sai:
1. Coordinator detecta mudança (heartbeat timeout ou novo subscribe)
2. Pausa todos os consumers do grupo
3. Recalcula assignments (partitions / num_consumers)
4. Notifica cada consumer das novas partitions
5. Consumers retomam de seus offsets
```

---

## 🧪 Benchmark Mínimo

Seu broker deve atingir:
- **100K+ mensagens/segundo** (produtor, msg de 100 bytes)
- **Latência p99 < 5ms** para produce
- **Zero perda de mensagens** com ack habilitado

Script de benchmark:
```go
// cmd/benchmark/main.go
// Lançar N producers em goroutines
// Cada um envia M mensagens
// Medir throughput e latência
// Verificar que consumers receberam TUDO
```

---

## 🧠 Conceitos Go Aplicados

| Conceito | Onde no projeto |
|----------|----------------|
| **TCP/Networking** | `net.Listen`, `net.Conn`, protocol framing |
| **encoding/binary** | Protocolo binário customizado |
| **Goroutines** | 1/connection, 1/partition, coordinator |
| **Channels** | Partition queues, consumer delivery |
| **sync.Mutex/RWMutex** | Metadata, offset tracking |
| **context.Context** | Timeouts, graceful shutdown |
| **io.Reader/Writer** | Storage engine, network streams |
| **os.File + Seek** | WAL, segment files, index |
| **hash/crc32** | Integridade de dados em disco |
| **Interfaces** | Storage, Protocol, Compressor |
| **Generics** | Result types, collections |
| **Testes** | Integration tests, benchmarks |
| **APIs** | Admin HTTP API |

---

## 🚀 Ordem de Implementação

```
Semana 1: Protocol + Storage
├── Definir frame format e message types
├── Implementar encoder/decoder com testes
├── Implementar WAL (append, read by offset)
├── Implementar segmentos e index
└── Benchmark: write/read 1M messages

Semana 2: Broker Core + Networking
├── TCP server que aceita conexões
├── Topic e Partition management
├── Producer flow: receive → route → store
├── Consumer flow: subscribe → consume → ack
└── Testes de integração producer↔broker↔consumer

Semana 3: Consumer Groups + Admin
├── Consumer Group Coordinator
├── Partition assignment (round-robin)
├── Offset commit/tracking
├── Rebalancing
├── Admin HTTP API
└── CLI

Semana 4+: Distribuído (se sobreviver)
├── Multi-broker cluster
├── Replication protocol
├── Leader election
└── ISR tracking
```

---

## ✅ Critérios de Sucesso

- [ ] Producer envia mensagens via TCP com protocolo binário
- [ ] Mensagens são persistidas em WAL com integridade (CRC32)
- [ ] Consumer Groups consomem com partition assignment correto
- [ ] Rebalancing funciona quando consumer entra/sai
- [ ] Offsets são tracked — consumer retoma de onde parou após restart
- [ ] 100K+ msg/s no benchmark
- [ ] Dados sobrevivem restart do broker (persistência real)
- [ ] Testes com 80%+ coverage
- [ ] Zero race conditions (`go test -race` passa limpo)
- [ ] `golangci-lint` passa sem erros

---

> 🏆 **Este projeto é brutalmente difícil.** Message brokers são uma das peças mais complexas
> de infraestrutura que existem. Se você implementar até a Fase 2, você está acima de 90%
> dos desenvolvedores Go. Fase 3 é territory de Staff Engineer.
>
> O Kafka levou anos e centenas de engenheiros. Você vai fazer uma versão simplificada
> sozinho. Isso é exatamente o tipo de coisa que impressiona em entrevistas e portfólio.

---

[← APIs](../modulo-14-apis/README.md) | [Voltar ao índice](../README.md)
