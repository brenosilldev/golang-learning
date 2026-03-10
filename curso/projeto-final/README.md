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

## 📖 Glossário e Conceitos Fundamentais

> Antes de codar, entenda **o que** você vai construir e **por que** existe. Cada conceito abaixo aparece no projeto — esta seção é sua bússola de referência.

---

### 🗂️ Message Broker — O que é e por que existe?

Um **Message Broker** é um intermediário que desacopla quem produz dados de quem os consome.

```
SEM broker:
  Serviço A ──────────────────────────────► Serviço B
             (se B cair, A perde os dados)

COM broker:
  Serviço A ──► [NexusMQ] ──► Serviço B
                    │         (B pode processar quando voltar)
                    └──► Serviço C
                         (vários consumidores independentes)
```

**Por que isso importa?**
- Serviços ficam **independentes** — A não precisa saber onde B está
- **Buffering natural** — o broker absorve picos de tráfego
- **Replay** — consumidores podem reler mensagens do passado
- É o coração de arquiteturas **event-driven** (Netflix, Uber, Nubank usam Kafka, que é o modelo do NexusMQ)

---

### 📦 Topic, Partition e Offset

**Topic** é como uma "pasta" ou "canal" de mensagens. Exemplo: `topic: pedidos`, `topic: pagamentos`.

**Partition** é a subdivisão de um topic para permitir paralelismo:
```
Topic "pedidos" com 3 partitions:
  Partition 0: [ msg1, msg4, msg7 ]  ← mensagens com key hash % 3 == 0
  Partition 1: [ msg2, msg5, msg8 ]
  Partition 2: [ msg3, msg6, msg9 ]
```
- Mensagens com a **mesma key** sempre vão para a **mesma partition** (garante ordem)
- Sem key → **round-robin** (distribuição equilibrada)

**Offset** é o número sequencial de cada mensagem dentro de uma partition. Como um índice de array — começa em 0 e só cresce. O consumer usa o offset para saber "li até aqui".

```
Partition 0: offset 0, offset 1, offset 2, offset 3...
                ↑ consumer está aqui (leu até offset 2)
```

---

### 📝 WAL — Write-Ahead Log

O **WAL** (Write-Ahead Log) é a técnica de persistência mais importante em bancos e sistemas distribuídos.

**A regra fundamental:** *antes de confirmar qualquer operação, escreva no log em disco.*

```
Mensagem chega:
  1. Escreve no arquivo WAL (append-only) ← PRIMEIRO
  2. Confirma pro producer que recebeu    ← SEGUNDO
  3. Serve para consumers lerem
```

**Por que append-only?**
- Operação mais rápida em disco — o head do HD não precisa se mover
- SSD tem throughput máximo em writes sequenciais
- Não tem fragmentação — simples de implementar e recuperar após crash

**Formato em disco no NexusMQ:**
```
┌──────────┬────────────┬──────────┬────────┬──────────┬────────┬───────┐
│ Offset(8)│ Timestamp(8│ KeyLen(4)│ Key(N) │ ValLen(4)│ Val(N) │ CRC(4)│
└──────────┴────────────┴──────────┴────────┴──────────┴────────┴───────┘
```
- **CRC32** = checksum para detectar corrupção (se o disco falhar no meio de um write)
- **Offset** = posição sequencial da mensagem
- **Segmentos**: quando o arquivo chega em X MB, fecha e abre um novo (facilita compactação e busca)

**Em Go:**
```go
// Append-only write: O_APPEND garante que nunca sobrescrevemos
file, _ := os.OpenFile("segment-000.log", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0644)

// Para leitura por offset, usamos Seek para ir direto à posição:
file.Seek(indexedPosition, io.SeekStart)
```

---

### 🔌 Protocolo Binário sobre TCP

**Por que não usar HTTP/JSON?**
- JSON = texto = desperdício de bytes. `{"offset": 12345}` são 17 bytes. Em binário: 4 bytes.
- HTTP tem overhead de headers (centenas de bytes por request)
- Para 100K msg/s, cada byte importa

**Como funciona:** você define seu próprio "idioma" entre producer, consumer e broker:

```
Frame (envelope de toda mensagem):
┌──────────────┬────────────┬─────────────────────┐
│  Length (4B) │  Type (1B) │  Payload (N bytes)  │
└──────────────┴────────────┴─────────────────────┘
     uint32       uint8             variável

Types:
  0x01 = PRODUCE    (producer → broker: "guarda isso")
  0x02 = CONSUME    (consumer → broker: "me dê mensagens")
  0x03 = ACK        (broker → producer: "guardei")
  0x04 = SUBSCRIBE  (consumer → broker: "quero este topic")
  0x05 = ERROR      (qualquer direção: "deu errado")
  0x06 = HEARTBEAT  (broker ↔ broker: "ainda estou vivo")
```

**Em Go, `encoding/binary` lida com isso:**
```go
// Escrever um uint32 em big-endian (padrão de rede):
binary.Write(conn, binary.BigEndian, uint32(len(payload)))

// Ler:
var length uint32
binary.Read(conn, binary.BigEndian, &length)
```

---

### 👥 Consumer Groups — Como o Load Balancing funciona

Um **Consumer Group** permite que vários consumers cooperem para processar um topic:

```
Topic "pedidos" (3 partitions) + Consumer Group "processador" (2 consumers):

  Partition 0 ──► Consumer A  (consumer A processa P0 e P1)
  Partition 1 ──► Consumer A
  Partition 2 ──► Consumer B  (consumer B processa P2)
```

**Regra de ouro:** cada partition é atribuída a **no máximo UM** consumer por grupo.
Isso garante que cada mensagem seja processada **exatamente uma vez** dentro do grupo.

**Rebalancing** acontece quando consumers entram ou saem:
```
Consumer C entra no grupo:
  Antes: A → P0, P1 | B → P2
  Após:  A → P0     | B → P2 | C → P1
```

O **Coordinator** (uma goroutine central) detecta isso via heartbeats e recalcula os assignments.

---

### 🔄 Replicação Leader/Follower (Fase 3)

Para tolerância a falhas, cada partition tem uma cópia em múltiplos nós:

```
Cluster de 3 brokers:
  Broker 1: Leader de P0, P1   |  Follower de P2
  Broker 2: Follower de P0     |  Leader de P2
  Broker 3: Follower de P0, P1 |  Follower de P2
```

- **Writes** sempre vão para o **Leader** — ele replica para os Followers
- **ISR (In-Sync Replicas)** = lista de Followers que estão atualizados
- Se o Leader cair, um Follower do ISR vira o novo Leader (**leader election**)

**Por que ISR importa?** Uma mensagem só é considerada "committed" quando todas as réplicas ISR a confirmaram. Isso garante que se o leader cair, nenhuma mensagem confirmada será perdida.

---

### 🔐 CRC32 — Integridade de Dados

**CRC32** (Cyclic Redundancy Check) é um checksum de 4 bytes calculado sobre os dados da mensagem. Quando você lê do disco, recalcula o CRC e compara:

```go
import "hash/crc32"

// Ao escrever:
checksum := crc32.ChecksumIEEE(messageBytes)
binary.Write(file, binary.BigEndian, checksum)

// Ao ler:
var storedCRC uint32
binary.Read(file, binary.BigEndian, &storedCRC)
computed := crc32.ChecksumIEEE(messageBytes)
if computed != storedCRC {
    return errors.New("dados corrompidos no disco")
}
```

Isso protege contra falhas de hardware (bit flip, write parcial durante queda de energia).

---

### 📡 Context e Graceful Shutdown

Em sistemas de produção, quando você mata o processo (`SIGTERM`), você não pode simplesmente parar no meio de uma operação. O **Graceful Shutdown** garante que:
1. Para de aceitar novas conexões
2. Termina de processar as mensagens em andamento
3. Fecha arquivos WAL com `Sync()` (flush para disco)
4. Confirma pendentes para producers

```go
// O padrão em Go:
ctx, cancel := context.WithCancel(context.Background())

// Em outra goroutine, aguarda sinal do OS:
sigChan := make(chan os.Signal, 1)
signal.Notify(sigChan, syscall.SIGTERM, syscall.SIGINT)
go func() {
    <-sigChan
    cancel() // propaga cancelamento para todas as goroutines
}()

// Cada componente respeita o context:
select {
case msg := <-partition.queue:
    store(msg)
case <-ctx.Done():
    flush(); return // saída limpa
}
```

---

### ⚡ Por que 100K msg/s é o benchmark mínimo?

Para referência:
- Kafka (bem configurado): 1-2 milhões msg/s
- NATS: 10-20 milhões msg/s
- RabbitMQ: 50K-100K msg/s

**100K msg/s** com Go é perfeitamente razoável usando:
- Batching (agrupa mensagens antes de escrever no disco)
- Buffer de canais (`make(chan Message, 1000)`)
- `bufio.Writer` (reduz chamadas de sistema de write)
- `O_SYNC` vs `Sync()` manual (controla quando dar flush)

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
