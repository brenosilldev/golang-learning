// Módulo 19 — Exercícios: Consenso & Raft
//
// 🧠 Antes de começar: sem olhar o README, tente responder:
//   1. Por que precisamos de consenso em sistemas distribuídos?
//   2. O que é "quorum" e por que N/2+1 nós?
//   3. O que acontece com o cluster se o líder cair?
//
// Instruções:
//   - Implemente cada função onde está marcado TODO
//   - Rode com: go run ex19_consensus_raft.go
//   - Experimente matar nós e observe a eleição acontecer

package main

import (
	"context"
	"fmt"
	"log"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// EXERCÍCIO 1 — Leader Election Simplificada (🟢 Fácil)
// =============================================================================
// Implemente um sistema de eleição de líder simples:
//   - Cada nó tem um ID e um term (mandato)
//   - Nós votam no candidato com maior term que não votaram ainda
//   - Quem receber maioria (N/2 + 1) de votos é eleito líder
//   - Líder envia heartbeats a cada 50ms para evitar re-eleições

type NodeRole int

const (
	Follower NodeRole = iota
	Candidate
	Leader
)

func (r NodeRole) String() string {
	return [...]string{"Follower", "Candidate", "Leader"}[r]
}

type LeaderElection struct {
	mu          sync.Mutex
	id          int
	role        NodeRole
	currentTerm int
	votedFor    int // -1 = não votou
	peers       []*LeaderElection
	heartbeatCh chan struct{}
}

func NewLeaderElection(id int) *LeaderElection {
	// TODO: inicialize com role=Follower, currentTerm=0, votedFor=-1
	panic("implemente NewLeaderElection")
}

func (n *LeaderElection) RequestVote(term, candidateID int) bool {
	// TODO: conceda o voto se:
	//   - term >= currentTerm
	//   - ainda não votou neste term (votedFor == -1 ou votedFor == candidateID)
	// Atualize currentTerm e votedFor se voto for concedido
	panic("implemente RequestVote")
}

func (n *LeaderElection) startElection() {
	// TODO: incremente o term, vote em si mesmo, peça votos aos peers em paralelo
	// Se ganhar quorum, torne-se Leader e comece a enviar heartbeats
	panic("implemente startElection")
}

func (n *LeaderElection) Run(ctx context.Context) {
	// TODO: loop principal:
	// - Se Follower: aguarde heartbeat. Se timeout (150-300ms aleatório), inicie eleição
	// - Se Leader: envie heartbeatCh para todos os peers a cada 50ms
	panic("implemente Run")
}

// =============================================================================
// EXERCÍCIO 2 — Distributed Lock com etcd (🟡 Médio)
// =============================================================================
// Simule um distributed lock sem etcd real (para não precisar de infraestrutura).
// Implemente a semântica correta: somente 1 holder por vez, com TTL automático.

type DistributedLock struct {
	mu        sync.Mutex
	holder    string
	expiresAt time.Time
	ttl       time.Duration
}

func NewDistributedLock(ttl time.Duration) *DistributedLock {
	// TODO: inicialize
	panic("implemente NewDistributedLock")
}

// TryLock tenta adquirir o lock. Retorna true se conseguiu, false se já está ocupado.
func (dl *DistributedLock) TryLock(clientID string) bool {
	// TODO: se não há holder ou o lock expirou, adquira para clientID
	panic("implemente TryLock")
}

// Lock bloqueia até conseguir o lock ou o contexto ser cancelado.
func (dl *DistributedLock) Lock(ctx context.Context, clientID string) error {
	// TODO: tente TryLock em loop com backoff até conseguir ou ctx cancelado
	panic("implemente Lock")
}

// Unlock libera o lock. Só o holder pode liberar.
func (dl *DistributedLock) Unlock(clientID string) error {
	// TODO: verifique se clientID é o holder antes de liberar
	panic("implemente Unlock")
}

// Renew renova o TTL do lock (keepalive — evita expiração durante uso longo)
func (dl *DistributedLock) Renew(clientID string) error {
	// TODO: renove expiresAt se clientID for o holder
	panic("implemente Renew")
}

// =============================================================================
// EXERCÍCIO 3 — Log Replication Simplificado (🔴 Difícil)
// =============================================================================
// Implemente a replicação de log do Raft de forma simplificada:
//   - Líder aceita entradas e as replica para todos os followers
//   - Uma entrada é "committed" quando a maioria confirmar
//   - Followers atrasados devem receber entradas que faltam

type LogEntry struct {
	Index   int
	Term    int
	Command string
}

type ReplicatedLog struct {
	mu          sync.Mutex
	entries     []LogEntry
	commitIndex int
	peers       []*ReplicatedLog
	isLeader    bool
	nextIndex   map[int]int // peerID → próxima entrada a enviar
}

func NewLeaderLog() *ReplicatedLog {
	// TODO: inicialize como líder com log vazio
	panic("implemente NewLeaderLog")
}

func NewFollowerLog() *ReplicatedLog {
	// TODO: inicialize como follower com log vazio
	panic("implemente NewFollowerLog")
}

// Append adiciona entrada ao log do líder e replica para followers
func (rl *ReplicatedLog) Append(command string) (int, error) {
	// TODO:
	// 1. Adicione a entrada ao log com o índice correto
	// 2. Replique para todos os followers em paralelo
	// 3. Quando quorum confirmar, incremente commitIndex
	// 4. Retorne o índice da entrada committed
	panic("implemente Append")
}

// AppendEntries é chamado pelo líder para replicar entradas
func (rl *ReplicatedLog) AppendEntries(entries []LogEntry, leaderCommit int) bool {
	// TODO:
	// 1. Adicione as entradas ao log local
	// 2. Se leaderCommit > commitIndex, atualize commitIndex
	panic("implemente AppendEntries")
}

// GetCommitted retorna todas as entradas committed
func (rl *ReplicatedLog) GetCommitted() []LogEntry {
	// TODO: retorne entries[0:commitIndex+1]
	panic("implemente GetCommitted")
}

// =============================================================================
// MAIN — Testa todos os exercícios
// =============================================================================

func testLeaderElection() {
	fmt.Println("\n=== LEADER ELECTION ===")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	nodes := make([]*LeaderElection, 5)
	for i := range nodes {
		nodes[i] = NewLeaderElection(i)
	}
	for i, n := range nodes {
		for j, peer := range nodes {
			if i != j {
				n.peers = append(n.peers, peer)
			}
		}
	}

	for _, n := range nodes {
		go n.Run(ctx)
	}

	<-ctx.Done()

	leaders := 0
	for _, n := range nodes {
		if n.role == Leader {
			leaders++
			fmt.Printf("Node %d é o LÍDER (term %d)\n", n.id, n.currentTerm)
		}
	}
	fmt.Printf("Total de líderes: %d (esperado: 1)\n", leaders)
}

func testDistributedLock() {
	fmt.Println("\n=== DISTRIBUTED LOCK ===")
	lock := NewDistributedLock(500 * time.Millisecond)
	var wg sync.WaitGroup
	results := make(chan string, 5)

	for i := 0; i < 5; i++ {
		wg.Add(1)
		clientID := fmt.Sprintf("client-%d", i)
		go func(id string) {
			defer wg.Done()
			ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
			defer cancel()
			if err := lock.Lock(ctx, id); err != nil {
				results <- fmt.Sprintf("%s: falhou (%v)", id, err)
				return
			}
			results <- fmt.Sprintf("%s: adquiriu o lock", id)
			time.Sleep(100 * time.Millisecond)
			lock.Unlock(id)
		}(clientID)
	}

	wg.Wait()
	close(results)
	for r := range results {
		fmt.Println(r)
	}
}

func testReplicatedLog() {
	fmt.Println("\n=== REPLICATED LOG ===")
	leader := NewLeaderLog()
	followers := []*ReplicatedLog{NewFollowerLog(), NewFollowerLog()}
	leader.peers = followers

	commands := []string{"SET x=1", "SET y=2", "DEL x", "SET z=3"}
	for _, cmd := range commands {
		idx, err := leader.Append(cmd)
		if err != nil {
			log.Printf("Erro ao adicionar '%s': %v", cmd, err)
			continue
		}
		fmt.Printf("Committed '%s' no índice %d\n", cmd, idx)
	}

	fmt.Printf("\nLíder committed: %v\n", leader.GetCommitted())
	for i, f := range followers {
		fmt.Printf("Follower %d committed: %v\n", i, f.GetCommitted())
	}
}

func main() {
	rand.Seed(time.Now().UnixNano())
	testLeaderElection()
	testDistributedLock()
	testReplicatedLog()
	fmt.Println("\n✅ Exercícios do Módulo 19 concluídos!")
}
