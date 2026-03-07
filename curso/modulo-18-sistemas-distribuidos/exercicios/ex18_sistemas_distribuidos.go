// Módulo 18 — Exercícios: Fundamentos de Sistemas Distribuídos
//
// 🧠 Antes de começar: sem olhar o README, tente responder:
//   1. O que o CAP Theorem diz que é impossível garantir simultaneamente?
//   2. Por que backoff exponencial é melhor que retry com intervalo fixo?
//   3. O que torna uma operação "idempotente"?
//
// Instruções:
//   - Implemente cada função onde está marcado TODO
//   - Rode com: go run ex18_sistemas_distribuidos.go
//   - Teste race conditions com: go run -race ex18_sistemas_distribuidos.go
//   - NÃO copie o código do README — escreva do zero, consultando apenas a assinatura

package main

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"sync"
	"time"
)

// =============================================================================
// EXERCÍCIO 1 — Circuit Breaker (🟢 Fácil)
// =============================================================================
// Implemente um Circuit Breaker com os 3 estados: Closed, Open, HalfOpen.
// O breaker deve:
//   - Abrir após `maxFailures` falhas consecutivas
//   - Tentar fechar após `resetTimeout` segundos no estado Open
//   - Fechar de volta após `successThreshold` sucessos consecutivos no HalfOpen

type CBState int

const (
	CBClosed CBState = iota
	CBOpen
	CBHalfOpen
)

type CircuitBreaker struct {
	mu               sync.Mutex
	state            CBState
	failures         int
	successes        int
	maxFailures      int
	successThreshold int
	resetTimeout     time.Duration
	lastFailure      time.Time
}

var ErrCircuitOpen = errors.New("circuit breaker: serviço indisponível")

func NewCircuitBreaker(maxFailures, successThreshold int, resetTimeout time.Duration) *CircuitBreaker {
	// TODO: inicialize e retorne um novo CircuitBreaker
	panic("implemente NewCircuitBreaker")
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
	// TODO: implemente a lógica de transição de estados
	// Dica: verifique o estado atual, execute fn(), e transicione baseado no resultado
	panic("implemente CircuitBreaker.Execute")
}

func (cb *CircuitBreaker) State() string {
	// TODO: retorne o estado como string ("Closed", "Open", "HalfOpen")
	panic("implemente CircuitBreaker.State")
}

// =============================================================================
// EXERCÍCIO 2 — Retry com Backoff Exponencial + Jitter (🟡 Médio)
// =============================================================================
// Implemente a função WithRetry que:
//   - Tenta executar fn() até maxAttempts vezes
//   - Espera baseDelay * 2^attempt entre tentativas
//   - Adiciona jitter aleatório de ±30% para evitar thundering herd
//   - Respeita o contexto (para se ctx for cancelado)
//   - Retorna erro descritivo com número de tentativas esgotadas

type RetryConfig struct {
	MaxAttempts int
	BaseDelay   time.Duration
	MaxDelay    time.Duration
}

var DefaultRetryConfig = RetryConfig{
	MaxAttempts: 5,
	BaseDelay:   100 * time.Millisecond,
	MaxDelay:    30 * time.Second,
}

func WithRetry(ctx context.Context, cfg RetryConfig, fn func() error) error {
	// TODO: implemente o retry com backoff exponencial
	// Fórmula: delay = min(baseDelay * 2^attempt, maxDelay) * (1 ± 0.3 * random)
	panic("implemente WithRetry")
}

// =============================================================================
// EXERCÍCIO 3 — KV Store com Eventual Consistency + Last-Write-Wins (🟡 Médio)
// =============================================================================
// Simule um KV store distribuído com 3 réplicas.
// Escrita: vai para todas as réplicas com um atraso aleatório (simula rede)
// Leitura: retorna o valor com o timestamp mais recente (last-write-wins)
// Objetivo: observar como réplicas ficam temporariamente inconsistentes

type VersionedValue struct {
	Value     string
	Timestamp int64 // Unix nano
	NodeID    int
}

type Replica struct {
	id   int
	mu   sync.RWMutex
	data map[string]VersionedValue
}

func NewReplica(id int) *Replica {
	// TODO: inicialize a replica
	panic("implemente NewReplica")
}

func (r *Replica) Write(key, value string) {
	// TODO: escreva com timestamp atual
	// Aplique last-write-wins: só atualiza se o novo timestamp for maior
	panic("implemente Replica.Write")
}

func (r *Replica) Read(key string) (VersionedValue, bool) {
	// TODO: leia o valor atual
	panic("implemente Replica.Read")
}

type DistributedKV struct {
	replicas []*Replica
}

func NewDistributedKV(numReplicas int) *DistributedKV {
	// TODO: crie numReplicas réplicas
	panic("implemente NewDistributedKV")
}

func (kv *DistributedKV) Write(key, value string) {
	// TODO: escreva em todas as réplicas com atraso aleatório (0-100ms)
	// Use goroutines — não espere cada réplica terminar
	panic("implemente DistributedKV.Write")
}

func (kv *DistributedKV) Read(key string) string {
	// TODO: leia de todas as réplicas e retorne o valor com maior timestamp
	panic("implemente DistributedKV.Read")
}

// =============================================================================
// EXERCÍCIO 4 — Idempotency Key Middleware (🔴 Difícil)
// =============================================================================
// Implemente um middleware que garante idempotência para qualquer função.
// Se a mesma idempotency key for usada, retorna o resultado cacheado.
// O cache deve expirar após `ttl` para evitar crescimento ilimitado.

type IdempotencyCache struct {
	// TODO: defina os campos necessários (cache, mutex, ttl)
}

type CachedResult struct {
	Result interface{}
	Err    error
	CachedAt time.Time
}

func NewIdempotencyCache(ttl time.Duration) *IdempotencyCache {
	// TODO: implemente
	panic("implemente NewIdempotencyCache")
}

func (c *IdempotencyCache) Execute(key string, fn func() (interface{}, error)) (interface{}, error) {
	// TODO: se key existe no cache e não expirou, retorne resultado cacheado
	// Se não existe, execute fn(), armazene no cache, retorne o resultado
	panic("implemente IdempotencyCache.Execute")
}

// =============================================================================
// MAIN — Testa todos os exercícios
// =============================================================================

func testCircuitBreaker() {
	fmt.Println("\n=== CIRCUIT BREAKER ===")
	cb := NewCircuitBreaker(3, 2, 2*time.Second)

	// Simula falhas para abrir o breaker
	for i := 0; i < 5; i++ {
		err := cb.Execute(func() error {
			if rand.Float32() < 0.7 {
				return errors.New("serviço fora")
			}
			return nil
		})
		fmt.Printf("Chamada %d: estado=%s, err=%v\n", i+1, cb.State(), err)
	}
}

func testRetry() {
	fmt.Println("\n=== RETRY COM BACKOFF ===")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	attempts := 0
	err := WithRetry(ctx, DefaultRetryConfig, func() error {
		attempts++
		if attempts < 4 {
			return fmt.Errorf("tentativa %d falhou", attempts)
		}
		return nil
	})
	fmt.Printf("Resultado após %d tentativas: err=%v\n", attempts, err)
}

func testDistributedKV() {
	fmt.Println("\n=== DISTRIBUTED KV (Eventual Consistency) ===")
	kv := NewDistributedKV(3)
	kv.Write("user:1", "Alice")

	// Leituras imediatas podem retornar valores desatualizados (consistência eventual)
	for i := 0; i < 3; i++ {
		fmt.Printf("Leitura imediata %d: %s\n", i, kv.Read("user:1"))
	}

	// Após convergência
	time.Sleep(200 * time.Millisecond)
	fmt.Printf("Após convergência: %s\n", kv.Read("user:1"))
}

func testIdempotency() {
	fmt.Println("\n=== IDEMPOTENCY CACHE ===")
	cache := NewIdempotencyCache(5 * time.Minute)
	calls := 0

	fn := func() (interface{}, error) {
		calls++
		return fmt.Sprintf("resultado-%d", calls), nil
	}

	r1, _ := cache.Execute("txn-abc", fn)
	r2, _ := cache.Execute("txn-abc", fn) // mesma key — deve retornar r1
	r3, _ := cache.Execute("txn-xyz", fn) // key diferente — executa fn novamente

	fmt.Printf("r1=%v, r2=%v (igual a r1? %v), r3=%v\n", r1, r2, r1 == r2, r3)
	fmt.Printf("fn foi chamada %d vezes (esperado: 2)\n", calls)
}

func main() {
	rand.Seed(time.Now().UnixNano())
	_ = math.Pow // evita import não-usado se você não usar em WithRetry

	testCircuitBreaker()
	testRetry()
	testDistributedKV()
	testIdempotency()

	fmt.Println("\n✅ Todos os exercícios concluídos!")
}
