# 🗳️ Módulo 19 — Consenso Distribuído & Raft

> **Nível**: Avançado | **Pré-requisito**: Módulo 18 (Fundamentos de Sistemas Distribuídos)

---

## 🤔 O Problema do Consenso

Como fazer **múltiplos nós** concordarem em um único valor quando qualquer um deles pode falhar a qualquer momento?

Esse problemas aparece em todo lugar:
- **Quem é o líder?** (leader election)
- **Em que ordem aplicar as operações?** (log replication)  
- **Qual é o valor "oficial" de uma chave?** (distributed KV store)

O **Teorema FLP** (Fischer, Lynch, Paterson, 1985) prova que consenso determinístico é **impossível** num sistema assíncrono com mesmo uma falha. Por isso algoritmos como Raft adicionam **timeouts** para sair do problema.

---

## 📖 O Algoritmo Raft

Raft foi projetado para ser **compreensível** — diferente de Paxos, que é famoso por ser difícil de entender.

### Os 3 papéis de um nó Raft:

```
┌──────────┐     election timeout     ┌───────────┐
│ FOLLOWER │  ─────────────────────►  │ CANDIDATE │
└──────────┘                          └─────┬─────┘
     ▲                                      │
     │       recebe heartbeat              │ ganha eleição
     │       do líder                      ▼
     │                               ┌──────────┐
     └──────────────────────────────►│  LEADER  │
                                     └──────────┘
```

### Como funciona uma eleição:

1. **Follower** não recebe heartbeat no timeout → vira **Candidate**
2. Incrementa `term` (mandato), vota em si mesmo, pede votos aos outros
3. Ganha maioria (quorum) → vira **Leader**
4. Líder envia heartbeats regulares para evitar novas eleições

### Replicação de Log:

```
Leader:   [1:set x=1] [2:set y=2] [3:del x] ← recebe novas entradas
             ↓           ↓           ↓
Follower: [1:set x=1] [2:set y=2] [3:del x] ← replicado
Follower: [1:set x=1] [2:set y=2]            ← atrasado, será atualizado
```

Uma entrada é **committed** quando mais de 50% dos nós a têm.

---

## 💻 Implementação Simplificada de Raft em Go

```go
package main

import (
    "context"
    "fmt"
    "math/rand"
    "sync"
    "time"
)

type NodeState int

const (
    Follower NodeState = iota
    Candidate
    Leader
)

func (s NodeState) String() string {
    return [...]string{"Follower", "Candidate", "Leader"}[s]
}

type LogEntry struct {
    Term    int
    Index   int
    Command string
}

type RaftNode struct {
    mu sync.Mutex

    id    int
    peers []*RaftNode

    // Estado persistente (em produção, salvar em disco)
    currentTerm int
    votedFor    int // -1 = nenhum
    log         []LogEntry

    // Estado volátil
    state       NodeState
    commitIndex int
    lastApplied int

    // Channels de comunicação
    heartbeatCh chan struct{}
    voteCh      chan struct{}
}

func NewRaftNode(id int) *RaftNode {
    return &RaftNode{
        id:          id,
        currentTerm: 0,
        votedFor:    -1,
        state:       Follower,
        heartbeatCh: make(chan struct{}, 1),
        voteCh:      make(chan struct{}, 1),
    }
}

func (n *RaftNode) electionTimeout() time.Duration {
    // Timeout aleatório entre 150-300ms para evitar eleições simultâneas
    return time.Duration(150+rand.Intn(150)) * time.Millisecond
}

func (n *RaftNode) Run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        default:
        }

        n.mu.Lock()
        state := n.state
        n.mu.Unlock()

        switch state {
        case Follower:
            n.runFollower(ctx)
        case Candidate:
            n.runCandidate(ctx)
        case Leader:
            n.runLeader(ctx)
        }
    }
}

func (n *RaftNode) runFollower(ctx context.Context) {
    timeout := time.NewTimer(n.electionTimeout())
    defer timeout.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-n.heartbeatCh:
            // Recebeu heartbeat, reinicia timer
            timeout.Reset(n.electionTimeout())
        case <-timeout.C:
            // Timeout — inicia eleição
            fmt.Printf("Node %d: timeout como Follower, iniciando eleição\n", n.id)
            n.mu.Lock()
            n.state = Candidate
            n.mu.Unlock()
            return
        }
    }
}

func (n *RaftNode) runCandidate(ctx context.Context) {
    n.mu.Lock()
    n.currentTerm++
    n.votedFor = n.id // vota em si mesmo
    term := n.currentTerm
    n.mu.Unlock()

    fmt.Printf("Node %d: iniciando eleição para term %d\n", n.id, term)

    votes := 1 // voto próprio
    var mu sync.Mutex
    var wg sync.WaitGroup

    for _, peer := range n.peers {
        peer := peer
        wg.Add(1)
        go func() {
            defer wg.Done()
            granted := peer.RequestVote(term, n.id)
            if granted {
                mu.Lock()
                votes++
                mu.Unlock()
            }
        }()
    }
    wg.Wait()

    n.mu.Lock()
    defer n.mu.Unlock()

    quorum := (len(n.peers)+1)/2 + 1
    if votes >= quorum {
        fmt.Printf("Node %d: eleito Leader com %d votos (quorum=%d)\n", n.id, votes, quorum)
        n.state = Leader
    } else {
        fmt.Printf("Node %d: perdeu eleição (%d votos), voltando a Follower\n", n.id, votes)
        n.state = Follower
    }
}

func (n *RaftNode) runLeader(ctx context.Context) {
    ticker := time.NewTicker(50 * time.Millisecond) // heartbeat interval
    defer ticker.Stop()

    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            n.sendHeartbeats()
        }
    }
}

func (n *RaftNode) RequestVote(term, candidateID int) bool {
    n.mu.Lock()
    defer n.mu.Unlock()

    if term < n.currentTerm {
        return false
    }
    if term > n.currentTerm {
        n.currentTerm = term
        n.votedFor = -1
        n.state = Follower
    }
    if n.votedFor == -1 || n.votedFor == candidateID {
        n.votedFor = candidateID
        fmt.Printf("Node %d: votou em %d para term %d\n", n.id, candidateID, term)
        return true
    }
    return false
}

func (n *RaftNode) AppendEntries(term, leaderID int, entries []LogEntry) bool {
    n.mu.Lock()
    defer n.mu.Unlock()

    if term < n.currentTerm {
        return false
    }
    n.currentTerm = term
    n.state = Follower

    // Entrega heartbeat não-bloqueante
    select {
    case n.heartbeatCh <- struct{}{}:
    default:
    }

    if len(entries) > 0 {
        n.log = append(n.log, entries...)
        fmt.Printf("Node %d: replicou %d entradas do Leader %d\n", n.id, len(entries), leaderID)
    }
    return true
}

func (n *RaftNode) sendHeartbeats() {
    n.mu.Lock()
    term := n.currentTerm
    id := n.id
    n.mu.Unlock()

    for _, peer := range n.peers {
        go peer.AppendEntries(term, id, nil)
    }
}

// AppendLog é chamado pelo cliente para adicionar uma entrada ao log
func (n *RaftNode) AppendLog(command string) error {
    n.mu.Lock()
    defer n.mu.Unlock()

    if n.state != Leader {
        return fmt.Errorf("não sou o líder")
    }

    entry := LogEntry{
        Term:    n.currentTerm,
        Index:   len(n.log),
        Command: command,
    }
    n.log = append(n.log, entry)
    fmt.Printf("Node %d (Leader): adicionou ao log: %s\n", n.id, command)

    // Em producão: replicar para followers e aguardar quorum antes de retornar
    for _, peer := range n.peers {
        go peer.AppendEntries(n.currentTerm, n.id, []LogEntry{entry})
    }
    return nil
}

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel()

    // Cria cluster com 3 nós
    nodes := make([]*RaftNode, 3)
    for i := range nodes {
        nodes[i] = NewRaftNode(i)
    }
    // Conecta peers
    for i, n := range nodes {
        for j, peer := range nodes {
            if i != j {
                n.peers = append(n.peers, peer)
            }
        }
    }

    // Inicia todos os nós
    for _, n := range nodes {
        go n.Run(ctx)
    }

    <-ctx.Done()
    fmt.Println("\nSimulação encerrada.")
    for _, n := range nodes {
        fmt.Printf("Node %d: term=%d, state=%s, log=%v\n",
            n.id, n.currentTerm, n.state, n.log)
    }
}
```

---

## 🗄️ etcd — Raft em Produção

`etcd` é o armazenamento distribuído do Kubernetes. Ele usa Raft internamente e expõe uma API KV simples.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    clientv3 "go.etcd.io/etcd/client/v3"
)

func main() {
    // Conectar ao etcd
    cli, err := clientv3.New(clientv3.Config{
        Endpoints:   []string{"localhost:2379"},
        DialTimeout: 5 * time.Second,
    })
    if err != nil {
        log.Fatal(err)
    }
    defer cli.Close()

    ctx := context.Background()

    // PUT — escrever com TTL (lease de 10 segundos)
    lease, _ := cli.Grant(ctx, 10)
    _, err = cli.Put(ctx, "/services/api-gateway", "192.168.1.100:8080",
        clientv3.WithLease(lease.ID))
    if err != nil {
        log.Fatal(err)
    }

    // GET — leitura linearizável
    resp, err := cli.Get(ctx, "/services/api-gateway")
    if err != nil {
        log.Fatal(err)
    }
    for _, kv := range resp.Kvs {
        fmt.Printf("Key: %s → Value: %s\n", kv.Key, kv.Value)
    }

    // WATCH — reagir a mudanças (service discovery)
    fmt.Println("Watching /services/...")
    watchCh := cli.Watch(ctx, "/services/", clientv3.WithPrefix())
    go func() {
        for watchResp := range watchCh {
            for _, event := range watchResp.Events {
                fmt.Printf("Evento: %s %s=%s\n",
                    event.Type, event.Kv.Key, event.Kv.Value)
            }
        }
    }()

    // Transação (compare-and-swap) — atomica e linearizável
    txnResp, err := cli.Txn(ctx).
        If(clientv3.Compare(clientv3.Value("/lock/leader"), "=", "")).
        Then(clientv3.OpPut("/lock/leader", "node-1")).
        Else(clientv3.OpGet("/lock/leader")).
        Commit()
    if err != nil {
        log.Fatal(err)
    }
    if txnResp.Succeeded {
        fmt.Println("Lock adquirido!")
    } else {
        fmt.Printf("Lock já existe: %s\n", txnResp.Responses[0].GetResponseRange().Kvs[0].Value)
    }

    time.Sleep(2 * time.Second)
}
```

---

## 🏗️ Onde Raft é Usado

| Sistema | Uso |
|---------|-----|
| **etcd** | Armazenamento de config do Kubernetes |
| **CockroachDB** | Replicação de dados entre ranges |
| **TiKV** | Storage distribuído do TiDB |
| **Consul** | Service discovery da HashiCorp |
| **NexusMQ** | Replicação de mensagens (nosso projeto final) |

---

## 📋 Exercícios

### 🟢 1. Visualizador de Eleição
Adicione logs detalhados ao código de Raft acima e visualize uma eleição completa com 5 nós, um deles "morrendo" durante a eleição.

### 🟡 2. Distributed Lock com etcd
Implemente um `DistributedLock` que:
- Usa etcd com lease para garantir que o lock expira automaticamente
- Suporta `TryLock()` (não-bloqueante) e `Lock()` (bloqueante com watch)
- Renova o lease automaticamente em background (keepalive)

### 🔴 3. Log Replication completo
Estenda o `RaftNode` para:
- Implementar `nextIndex[]` e `matchIndex[]` por follower
- Replicar corretamente entradas para followers atrasados
- Marcar entradas como committed quando maioria confirmar

---

> **Próximo**: [Módulo 20 — Message Queues & Event-Driven](../modulo-20-message-queues/README.md)
