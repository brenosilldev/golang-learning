# 📦 Módulo 18 — Fundamentos de Sistemas Distribuídos

> **Nível**: Avançado | **Pré-requisito**: Módulos 01-17 (especialmente concorrência e Go para produção)

---

## 🤔 O que é um Sistema Distribuído?

Um sistema distribuído é um conjunto de computadores independentes que **parecem um único sistema** para o usuário. A lógica roda em múltiplas máquinas que se comunicam via rede.

**Exemplos do mundo real**: Kubernetes, Kafka, CockroachDB, Redis Cluster, Cassandra, DynamoDB.

**Por que Go domina esse espaço?**
- Goroutines tornam fácil modelar múltiplos nós simultaneamente
- Channels são ideais para comunicação entre componentes
- A stdlib tem TCP, TLS, HTTP/2 e encoding binário de primeira classe
- Binário único simplifica o deploy em múltiplos nós

---

## ⚠️ O Teorema CAP

Todo sistema distribuído tem **três propriedades** desejáveis, mas só pode garantir **duas ao mesmo tempo**:

```
        Consistency (Consistência)
           /\
          /  \
         /    \
        /  CA  \
       /--------\
      / CP |  AP \
     /     |      \
    /-------|-------\
   Partition Tolerance
   (Tolerância a Partições)
```

| Combinação | Exemplo | Quando usar |
|-----------|---------|-------------|
| **CP** (Consistência + Partição) | CockroachDB, etcd, Zookeeper | Sistemas financeiros, onde dado errado é pior que indisponibilidade |
| **AP** (Disponibilidade + Partição) | Cassandra, DynamoDB, CouchDB | E-commerce, redes sociais, onde disponibilidade > consistência estrita |
| **CA** (Consistência + Disponibilidade) | Bancos relacionais tradicionais | Não funciona em redes reais (partições são inevitáveis) |

> **Na prática**: toda rede sofre partições. A escolha real é entre **CP** e **AP**.

---

## 🕐 Modelos de Consistência

Do mais forte ao mais fraco:

### 1. Linearizabilidade (Strict Consistency)
Toda operação parece acontecer instantaneamente num único ponto no tempo global. A mais forte e a mais cara.

```go
// Exemplo: etcd com leitura linearizável
resp, err := cli.Get(ctx, "key", clientv3.WithSerializable())
```

### 2. Consistência Sequencial
Operações de um único cliente aparecem na ordem em que foram feitas, mas clientes diferentes podem ver ordens distintas.

### 3. Eventual Consistency (Consistência Eventual)
Se parar de escrever, eventualmente todos os nós convergem para o mesmo valor. Cassandra e DynamoDB usam isso.

```go
// Simulação de eventual consistency
type ReplicatedValue struct {
    value     string
    timestamp int64 // last-write-wins
    mu        sync.RWMutex
}

func (r *ReplicatedValue) Write(val string) {
    r.mu.Lock()
    defer r.mu.Unlock()
    now := time.Now().UnixNano()
    if now > r.timestamp { // last-write-wins
        r.value = val
        r.timestamp = now
    }
}
```

---

## 🕰️ Clocks Lógicos — Lamport Timestamps

Em sistemas distribuídos, **não existe relógio global**. O NTP só garante ~100ms de precisão. Para ordenar eventos, usamos clocks lógicos.

### Regra de Lamport
1. Incrementa o clock antes de enviar qualquer mensagem
2. Ao receber uma mensagem, `clock = max(local, received) + 1`

```go
package main

import (
    "fmt"
    "sync"
    "sync/atomic"
)

type LamportClock struct {
    value int64
}

func (lc *LamportClock) Tick() int64 {
    return atomic.AddInt64(&lc.value, 1)
}

func (lc *LamportClock) Update(received int64) int64 {
    for {
        current := atomic.LoadInt64(&lc.value)
        newVal := received + 1
        if current >= newVal {
            newVal = current + 1
        }
        if atomic.CompareAndSwapInt64(&lc.value, current, newVal) {
            return newVal
        }
    }
}

func (lc *LamportClock) Get() int64 {
    return atomic.LoadInt64(&lc.value)
}

func main() {
    var wg sync.WaitGroup
    nodeA := &LamportClock{}
    nodeB := &LamportClock{}

    wg.Add(2)

    go func() {
        defer wg.Done()
        t1 := nodeA.Tick()
        fmt.Printf("Node A: evento local em t=%d\n", t1)
        // "Envia" mensagem para B com timestamp
        sentAt := nodeA.Tick()
        fmt.Printf("Node A: envia msg com t=%d\n", sentAt)
        nodeB.Update(sentAt) // B recebe
    }()

    go func() {
        defer wg.Done()
        t1 := nodeB.Tick()
        fmt.Printf("Node B: evento local em t=%d\n", t1)
    }()

    wg.Wait()
    fmt.Printf("Clock A final: %d | Clock B final: %d\n", nodeA.Get(), nodeB.Get())
}
```

---

## 🔁 Idempotência

Uma operação é **idempotente** se executá-la múltiplas vezes produz o mesmo resultado que executar uma vez.

**Por que importa?** Em redes instáveis, você SEMPRE fará retries. Se a operação não for idempotente, vai criar dados duplicados.

```go
// ❌ NÃO idempotente — cada chamada adiciona saldo
func (db *DB) AddBalance(userID string, amount float64) error {
    _, err := db.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, userID)
    return err
}

// ✅ Idempotente — usa idempotency key
func (db *DB) AddBalanceIdempotent(txID, userID string, amount float64) error {
    // Verifica se transação já foi processada
    var exists bool
    db.QueryRow("SELECT EXISTS(SELECT 1 FROM transactions WHERE id = ?)", txID).Scan(&exists)
    if exists {
        return nil // já processado, retorna sem erro
    }
    
    tx, _ := db.Begin()
    tx.Exec("INSERT INTO transactions (id, user_id, amount) VALUES (?, ?, ?)", txID, userID, amount)
    tx.Exec("UPDATE accounts SET balance = balance + ? WHERE id = ?", amount, userID)
    return tx.Commit()
}
```

---

## 🔄 Retry com Backoff Exponencial

Nunca faça retry com intervalo fixo — vai derrubar o servidor no momento em que ele está tentando se recuperar.

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "math"
    "math/rand"
    "time"
)

type RetryConfig struct {
    MaxAttempts int
    BaseDelay   time.Duration
    MaxDelay    time.Duration
    Multiplier  float64
    Jitter      bool // adiciona aleatoriedade para evitar thundering herd
}

var DefaultRetry = RetryConfig{
    MaxAttempts: 5,
    BaseDelay:   100 * time.Millisecond,
    MaxDelay:    30 * time.Second,
    Multiplier:  2.0,
    Jitter:      true,
}

func WithRetry(ctx context.Context, cfg RetryConfig, fn func() error) error {
    var lastErr error
    for attempt := 0; attempt < cfg.MaxAttempts; attempt++ {
        if err := ctx.Err(); err != nil {
            return fmt.Errorf("contexto cancelado: %w", err)
        }

        lastErr = fn()
        if lastErr == nil {
            return nil
        }

        if attempt == cfg.MaxAttempts-1 {
            break
        }

        delay := time.Duration(float64(cfg.BaseDelay) * math.Pow(cfg.Multiplier, float64(attempt)))
        if delay > cfg.MaxDelay {
            delay = cfg.MaxDelay
        }
        if cfg.Jitter {
            // ±30% de jitter para evitar thundering herd
            jitter := time.Duration(rand.Float64() * float64(delay) * 0.3)
            delay += jitter
        }

        fmt.Printf("Tentativa %d falhou: %v. Aguardando %v...\n", attempt+1, lastErr, delay)
        select {
        case <-time.After(delay):
        case <-ctx.Done():
            return ctx.Err()
        }
    }
    return fmt.Errorf("todas as %d tentativas falharam: %w", cfg.MaxAttempts, lastErr)
}

// Exemplo de uso
func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    attempts := 0
    err := WithRetry(ctx, DefaultRetry, func() error {
        attempts++
        if attempts < 4 {
            return errors.New("serviço indisponível")
        }
        return nil // sucesso na 4ª tentativa
    })

    if err != nil {
        fmt.Printf("Erro final: %v\n", err)
    } else {
        fmt.Printf("Sucesso após %d tentativas!\n", attempts)
    }
}
```

---

## 🔌 Circuit Breaker

Protege seu serviço de continuar chamando um serviço que definitivamente está fora. Evita cascade failures.

```
CLOSED ──[falhas > limite]──► OPEN
  ▲                             │
  │                     [timeout expirou]
  │                             ▼
  └────[sucesso]────── HALF-OPEN
```

```go
package main

import (
    "errors"
    "fmt"
    "sync"
    "time"
)

type State int

const (
    StateClosed State = iota   // operando normalmente
    StateOpen                  // bloqueando chamadas
    StateHalfOpen              // testando recuperação
)

var ErrCircuitOpen = errors.New("circuit breaker aberto")

type CircuitBreaker struct {
    mu           sync.Mutex
    state        State
    failures     int
    maxFailures  int
    timeout      time.Duration
    lastFailure  time.Time
    successCount int
}

func NewCircuitBreaker(maxFailures int, timeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        maxFailures: maxFailures,
        timeout:     timeout,
        state:       StateClosed,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    cb.mu.Lock()
    defer cb.mu.Unlock()

    switch cb.state {
    case StateOpen:
        if time.Since(cb.lastFailure) > cb.timeout {
            fmt.Println("Circuit: OPEN → HALF-OPEN")
            cb.state = StateHalfOpen
            cb.successCount = 0
        } else {
            return ErrCircuitOpen
        }

    case StateHalfOpen:
        // Em half-open, só deixa uma chamada passar por vez

    case StateClosed:
        // Opera normalmente
    }

    err := fn()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.state == StateHalfOpen || cb.failures >= cb.maxFailures {
            fmt.Printf("Circuit: → OPEN (falhas: %d)\n", cb.failures)
            cb.state = StateOpen
        }
        return err
    }

    // Sucesso
    if cb.state == StateHalfOpen {
        cb.successCount++
        if cb.successCount >= 3 { // 3 sucessos consecutivos para fechar
            fmt.Println("Circuit: HALF-OPEN → CLOSED ✅")
            cb.state = StateClosed
            cb.failures = 0
        }
    } else {
        cb.failures = 0
    }
    return nil
}

func main() {
    cb := NewCircuitBreaker(3, 5*time.Second)

    // Simula falhas
    for i := 0; i < 6; i++ {
        err := cb.Execute(func() error {
            return errors.New("timeout")
        })
        fmt.Printf("Chamada %d: %v\n", i+1, err)
    }
}
```

---

## 📋 Exercícios

### 🟢 1. Simulação de CAP
Implemente um KV store simples que pode ser configurado como CP ou AP:
- Modo CP: retorna erro se detecta inconsistência
- Modo AP: retorna valor possivelmente desatualizado mas nunca falha

### 🟡 2. Vector Clocks
Implemente Vector Clocks (extensão de Lamport) que conseguem detectar causalidade entre eventos em múltiplos nós simultaneamente.

### 🟡 3. Retry Middleware HTTP
Crie um `http.RoundTripper` customizado que:
- Faz retry automático em erros 5xx e timeouts
- Implementa backoff exponencial com jitter
- Respeita o contexto (cancela corretamente)

### 🔴 4. Circuit Breaker Genérico com Prometheus
Estenda o `CircuitBreaker` para:
- Exportar métricas Prometheus (estado atual, contagem de falhas, latência)
- Ser genérico com type parameters
- Suportar múltiplos serviços downstream simultaneamente

---

> **Próximo**: [Módulo 19 — Consenso & Raft](../modulo-19-consensus-raft/README.md) — como garantir que múltiplos nós concordam em um único valor mesmo com falhas.
