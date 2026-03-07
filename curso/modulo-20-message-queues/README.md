# 📨 Módulo 20 — Message Queues & Event-Driven Architecture

> **Nível**: Avançado | **Pré-requisito**: Módulos 11 (Concorrência), 14 (APIs), 18 (Fundamentos Distribuídos)

---

## 🤔 Por que Event-Driven?

Em arquiteturas monolíticas, os serviços chamam uns aos outros diretamente (HTTP/gRPC). Isso cria **acoplamento temporal** — se o serviço B estiver fora, o serviço A falha.

```
❌ Acoplamento direto (síncrono):
   A → HTTP → B → HTTP → C
   (B offline = A falha)

✅ Event-Driven (assíncrono):
   A → [Fila] → B
               → C
   (B offline = mensagem fica na fila, entregue quando B voltar)
```

**Casos de uso reais**:
- Notificações (email, SMS, push) — nunca bloquear o fluxo principal
- Processamento de pagamentos — garantia de entrega, auditoria
- Sincronização de dados entre microserviços
- Data pipelines (clickstream, logs, métricas)

---

## 📊 Comparativo das Principais Ferramentas

| | Kafka | NATS | RabbitMQ |
|---|---|---|---|
| **Modelo** | Log distribuído | Pub/Sub / Queue | AMQP (exchange/queue) |
| **Persistência** | Sim (por padrão, dias/semanas) | Opcional (JetStream) | Sim |
| **Throughput** | Muito alto (milhões/s) | Extremamente alto | Alto |
| **Latência** | ~5-10ms | ~1ms | ~1-5ms |
| **Ordenação** | Por partição | Não garantida | Por fila |
| **Replay** | Sim (offset rewind) | JetStream | Não (fila vazia depois de ACK) |
| **Caso de uso** | Data pipeline, event sourcing | IoT, microserviços low-latency | Workflows, filas de trabalho |
| **Go client** | `sarama`, `confluent-kafka-go` | `nats.go` | `amqp091-go` |

---

## 🔴 Apache Kafka com Go

### Conceitos essenciais:

```
Tópico: "pagamentos"
├── Partição 0: [msg1][msg5][msg9]  ← Producer escolhe por key
├── Partição 1: [msg2][msg6][msg10]
└── Partição 2: [msg3][msg4][msg7]

Consumer Group: "processador-pagamentos"
├── Consumer A → lê Partição 0
├── Consumer B → lê Partição 1
└── Consumer C → lê Partição 2
```

**Regra**: um offset só avança quando você faz ACK. Se o consumer crashar, recomeça do último offset confirmado.

### Producer:

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "time"

    "github.com/IBM/sarama"
)

type PagamentoEvent struct {
    ID        string    `json:"id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    Currency  string    `json:"currency"`
    CreatedAt time.Time `json:"created_at"`
}

func conectarKafka(brokers []string) sarama.SyncProducer {
    cfg := sarama.NewConfig()
    cfg.Producer.Return.Successes = true
    cfg.Producer.RequiredAcks = sarama.WaitForAll // todos os ISR confirmam
    cfg.Producer.Retry.Max = 5
    cfg.Producer.Idempotent = true // exactly-once semântica por partição
    cfg.Net.MaxOpenRequests = 1    // necessário para idempotência

    producer, err := sarama.NewSyncProducer(brokers, cfg)
    if err != nil {
        log.Fatalf("Erro ao criar producer: %v", err)
    }
    return producer
}

func publicarPagamento(producer sarama.SyncProducer, evento PagamentoEvent) error {
    payload, err := json.Marshal(evento)
    if err != nil {
        return fmt.Errorf("erro ao serializar: %w", err)
    }

    msg := &sarama.ProducerMessage{
        Topic: "pagamentos",
        Key:   sarama.StringEncoder(evento.UserID), // mesma key → mesma partição → ordem garantida por usuário
        Value: sarama.ByteEncoder(payload),
        Headers: []sarama.RecordHeader{
            {Key: []byte("event-type"), Value: []byte("payment.created")},
            {Key: []byte("schema-version"), Value: []byte("v1")},
        },
    }

    partition, offset, err := producer.SendMessage(msg)
    if err != nil {
        return fmt.Errorf("erro ao enviar: %w", err)
    }
    fmt.Printf("Publicado → partição=%d offset=%d\n", partition, offset)
    return nil
}

func main() {
    producer := conectarKafka([]string{"localhost:9092"})
    defer producer.Close()

    // Publica 10 eventos
    for i := 0; i < 10; i++ {
        err := publicarPagamento(producer, PagamentoEvent{
            ID:        fmt.Sprintf("pay-%d", i),
            UserID:    fmt.Sprintf("user-%d", i%3), // 3 usuários, ordem garantida por usuário
            Amount:    float64(i) * 100.0,
            Currency:  "BRL",
            CreatedAt: time.Now(),
        })
        if err != nil {
            log.Printf("Erro: %v", err)
        }
    }
}
```

### Consumer com Consumer Group:

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "os"
    "os/signal"
    "syscall"

    "github.com/IBM/sarama"
)

type ConsumerGroupHandler struct {
    ready chan bool
}

// Setup é chamado ao iniciar/reiniciar o consumer
func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
    close(h.ready)
    return nil
}

// Cleanup é chamado ao encerrar
func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error { return nil }

// ConsumeClaim processa mensagens de uma partição
func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
    for msg := range claim.Messages() {
        var evento PagamentoEvent
        if err := json.Unmarshal(msg.Value, &evento); err != nil {
            log.Printf("Erro ao deserializar: %v", err)
            session.MarkMessage(msg, "") // marca mesmo com erro (dead-letter em produção)
            continue
        }

        fmt.Printf("Processando pagamento %s de user %s: R$%.2f (partição=%d, offset=%d)\n",
            evento.ID, evento.UserID, evento.Amount, msg.Partition, msg.Offset)

        // Processa o evento...
        if err := processarPagamento(evento); err != nil {
            log.Printf("Erro ao processar %s: %v", evento.ID, err)
            // Em produção: enviar para dead-letter topic
        }

        // Faz ACK — só avança o offset depois de processar com sucesso
        session.MarkMessage(msg, "")
    }
    return nil
}

func processarPagamento(evt PagamentoEvent) error {
    // Lógica de negócio...
    return nil
}

func main() {
    cfg := sarama.NewConfig()
    cfg.Consumer.Group.Rebalance.GroupStrategies = []sarama.BalanceStrategy{
        sarama.NewBalanceStrategyRoundRobin(),
    }
    cfg.Consumer.Offsets.Initial = sarama.OffsetOldest // começa do início se não há offset salvo

    client, err := sarama.NewConsumerGroup([]string{"localhost:9092"}, "processador-pagamentos", cfg)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    ctx, cancel := context.WithCancel(context.Background())
    handler := &ConsumerGroupHandler{ready: make(chan bool)}

    go func() {
        for {
            if err := client.Consume(ctx, []string{"pagamentos"}, handler); err != nil {
                log.Printf("Erro no consumer: %v", err)
            }
            if ctx.Err() != nil {
                return
            }
            handler.ready = make(chan bool) // reset para próximo ciclo
        }
    }()

    <-handler.ready
    fmt.Println("Consumer pronto!")

    // Aguarda sinal de parada
    sigterm := make(chan os.Signal, 1)
    signal.Notify(sigterm, syscall.SIGTERM, syscall.SIGINT)
    <-sigterm
    cancel()
}
```

---

## 💚 NATS — Alta Performance com Go

NATS é um message broker escrito em Go, com foco em latência ultra-baixa. NATS JetStream adiciona persistência.

```go
package main

import (
    "fmt"
    "log"
    "time"

    "github.com/nats-io/nats.go"
)

func main() {
    nc, err := nats.Connect(nats.DefaultURL,
        nats.PingInterval(20*time.Second),
        nats.MaxPingsOutstanding(5),
        nats.ReconnectWait(2*time.Second),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer nc.Drain()

    // === Pub/Sub simples ===
    nc.Subscribe("pedidos.criados", func(msg *nats.Msg) {
        fmt.Printf("Pedido recebido: %s\n", msg.Data)
    })

    nc.Subscribe("pedidos.*", func(msg *nats.Msg) {
        fmt.Printf("Todos os pedidos — subject: %s, data: %s\n", msg.Subject, msg.Data)
    })

    nc.Publish("pedidos.criados", []byte(`{"id":"123","total":250.00}`))
    nc.Publish("pedidos.cancelados", []byte(`{"id":"124"}`))

    // === Request/Reply (RPC sobre NATS) ===
    nc.Subscribe("calcular.frete", func(msg *nats.Msg) {
        frete := calculaFrete(string(msg.Data))
        msg.Respond([]byte(frete))
    })

    resposta, err := nc.Request("calcular.frete", []byte("CEP:01310-100"), 2*time.Second)
    if err != nil {
        log.Printf("Timeout no request: %v", err)
    } else {
        fmt.Printf("Frete calculado: %s\n", resposta.Data)
    }

    // === JetStream — persistência e at-least-once delivery ===
    js, _ := nc.JetStream()

    // Cria stream
    js.AddStream(&nats.StreamConfig{
        Name:       "PEDIDOS",
        Subjects:   []string{"pedidos.>"},
        Retention:  nats.LimitsPolicy,
        MaxAge:     7 * 24 * time.Hour, // retém por 7 dias
        Replicas:   3,                  // 3 réplicas para HA
    })

    // Publica com confirmação de persistência
    ack, err := js.Publish("pedidos.criados", []byte(`{"id":"999"}`))
    if err != nil {
        log.Printf("Erro ao publicar: %v", err)
    } else {
        fmt.Printf("Persistido: stream=%s, seq=%d\n", ack.Stream, ack.Sequence)
    }

    // Consumer durável (sobrevive a restarts)
    js.Subscribe("pedidos.>", func(msg *nats.Msg) {
        fmt.Printf("JetStream recebeu: %s\n", msg.Data)
        msg.Ack() // ACK explícito
    }, nats.Durable("processador-pedidos"),
        nats.AckExplicit(),
        nats.MaxDeliver(3), // máximo 3 tentativas antes de dead-letter
    )

    time.Sleep(time.Second)
}

func calculaFrete(cep string) string {
    return fmt.Sprintf("R$25.00 para %s", cep)
}
```

---

## 🏗️ Padrões Arquiteturais

### Outbox Pattern (garantia de entrega)

Problema: publicar num banco e num broker é uma operação distribuída — pode falhar no meio.

```go
// ❌ Problemático: e se crashar entre o Commit e o Publish?
tx.Commit()
kafka.Publish(evento) // ← crash aqui = dado no banco mas não no Kafka

// ✅ Outbox Pattern: salva no banco junto com a transação
tx.Begin()
tx.Exec("INSERT INTO pedidos ...")
tx.Exec("INSERT INTO outbox (payload, topic) VALUES (?, ?)", jsonEvento, "pedidos.criados")
tx.Commit() // operação atômica — salva pedido E outbox juntos

// Worker separado lê outbox e publica no Kafka
func outboxWorker(db *sql.DB, producer sarama.SyncProducer) {
    for {
        rows, _ := db.Query("SELECT id, topic, payload FROM outbox WHERE published = false LIMIT 100")
        for rows.Next() {
            var id int64; var topic, payload string
            rows.Scan(&id, &topic, &payload)
            producer.SendMessage(&sarama.ProducerMessage{
                Topic: topic,
                Value: sarama.StringEncoder(payload),
            })
            db.Exec("UPDATE outbox SET published = true WHERE id = ?", id)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

### CQRS — Command Query Responsibility Segregation

```
Client → [Command] → Write Model (Postgres) → [Event] → Read Model (Elasticsearch)
                                                           ↓
Client ←←←←←←←←←←←←←←←←←← [Query] ←←←←←←←←←←←←←←←←←←←
```

---

## 📋 Exercícios

### 🟢 1. NATS Pub/Sub local
Rode um servidor NATS (`docker run nats`) e implemente:
- Publisher que envia métricas de CPU a cada segundo
- 3 subscribers com lógica diferente (log, alert, store)

### 🟡 2. Dead Letter Queue
Implemente um consumer Kafka que:
- Processa mensagens de `pedidos.criados`
- Após 3 falhas na mesma mensagem, move para `pedidos.dlq` com metadata de erro

### 🔴 3. Saga Pattern
Implemente o padrão Saga para um fluxo de checkout:
- `pagamento.criado` → reserva estoque → `estoque.reservado` → confirma pagamento
- Se qualquer etapa falhar, desfaz as anteriores (compensating transactions)

---

> **Próximo**: [Módulo 21 — Observabilidade & SRE](../modulo-21-observabilidade/README.md)
