// Módulo 22 — Exercícios: Cloud Native & Kubernetes
//
// 🧠 Antes de começar:
//   1. Por que uma Liveness probe diferente de Readiness probe?
//   2. O que é graceful shutdown e por que o Kubernetes precisa disso?
//   3. O que acontece com as requests em andamento quando um pod recebe SIGTERM?
//
// Rode com: go run ex22_cloud_native.go
// Teste o shutdown: Ctrl+C (SIGINT) e observe o graceful shutdown

package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

// Referências que você vai usar no Exercício 2 (GracefulServer.Run):
// signal.NotifyContext, signal.Notify, syscall.SIGTERM, syscall.SIGINT
var _ = signal.NotifyContext
var _ = syscall.SIGTERM

// =============================================================================
// EXERCÍCIO 1 — Health Check Server (🟢 Fácil)
// =============================================================================
// Implemente um servidor HTTP com 3 endpoints de saúde:
//   - GET /health/live  → sempre 200 (app está viva)
//   - GET /health/ready → 200 se pronta, 503 se não
//   - GET /metrics      → saída de métricas simples (texto)
//
// O servidor deve:
//   - Começar no estado "not ready" (initializing)
//   - Ficar "ready" após 2 segundos (simula inicialização)
//   - Contar requests recebidas e erros

type HealthServer struct {
	ready          atomic.Bool
	requestCount   atomic.Int64
	errorCount     atomic.Int64
	version        string
	startedAt      time.Time
}

func NewHealthServer(version string) *HealthServer {
	// TODO: inicialize e defina startedAt = time.Now()
	panic("implemente NewHealthServer")
}

func (s *HealthServer) SimulateInit(delay time.Duration) {
	// TODO: aguarde delay e então defina ready = true
	// Use goroutine para não bloquear
	panic("implemente SimulateInit")
}

func (s *HealthServer) Handler() http.Handler {
	// TODO: crie um http.ServeMux com os 3 endpoints
	// /health/live → sempre 200 JSON {"status":"ok","version":...}
	// /health/ready → 200 se ready, 503 se não, com uptime
	// /metrics → texto com request_count, error_count, uptime_seconds
	// Qualquer outra rota → 404 e incrementa errorCount
	panic("implemente Handler")
}

// =============================================================================
// EXERCÍCIO 2 — Graceful Shutdown (🟡 Médio)
// =============================================================================
// Implemente um servidor HTTP com graceful shutdown correto.
// O servidor deve:
//   1. Receber SIGTERM ou SIGINT
//   2. Parar de aceitar novas conexões (readiness probe retorna 503)
//   3. Aguardar requests em andamento terminarem (até drainTimeout)
//   4. Logar cada etapa do shutdown
//
// Simule uma request longa para testar que ela termina antes do shutdown.

type GracefulServer struct {
	server       *http.Server
	healthServer *HealthServer
	logger       *slog.Logger
	drainTimeout time.Duration
	// TODO: adicione campo para contar requests ativas
}

func NewGracefulServer(addr string, drainTimeout time.Duration) *GracefulServer {
	// TODO: configure o http.Server com timeouts adequados para K8s:
	// ReadTimeout: 10s, WriteTimeout: 30s, IdleTimeout: 120s
	panic("implemente NewGracefulServer")
}

func (gs *GracefulServer) Run() error {
	// TODO: Implemente o ciclo de vida completo:
	//
	// 1. Inicie o servidor em goroutine
	// 2. Aguarde SIGTERM ou SIGINT via signal.NotifyContext ou signal.Notify
	// 3. Ao receber sinal:
	//    a. Logue "recebeu sinal, iniciando shutdown"
	//    b. Defina healthServer.ready = false (readiness probe vira 503)
	//    c. Aguarde 5s (deregistration delay — K8s remove do LB)
	//    d. Chame server.Shutdown(ctx) com drainTimeout
	//    e. Logue "shutdown concluído"
	// 4. Retorne erro se o servidor falhar (não por shutdown normal)
	panic("implemente GracefulServer.Run")
}

// =============================================================================
// EXERCÍCIO 3 — Config from Environment (🟢 Fácil)
// =============================================================================
// Implemente um carregador de configuração que:
//   - Lê variáveis de ambiente
//   - Usa valores padrão quando não definidas
//   - Valida configurações obrigatórias
//   - Nunca expõe valores sensíveis no log (mascarar senhas/tokens)

type AppConfig struct {
	Port        string
	DatabaseURL string // obrigatório
	AppVersion  string
	LogLevel    string
	MaxWorkers  int
	Timeout     time.Duration
}

type ConfigError struct {
	Missing []string
}

func (e *ConfigError) Error() string {
	return fmt.Sprintf("configurações obrigatórias faltando: %v", e.Missing)
}

func LoadConfig() (AppConfig, error) {
	// TODO: leia as variáveis de ambiente abaixo com os defaults indicados:
	//   PORT          → default "8080"
	//   DATABASE_URL  → obrigatório (sem default)
	//   APP_VERSION   → default "dev"
	//   LOG_LEVEL     → default "info"
	//   MAX_WORKERS   → default 10, parse como int
	//   TIMEOUT       → default "30s", parse como duration
	//
	// Se DATABASE_URL não estiver definida, retorne ConfigError
	panic("implemente LoadConfig")
}

func (c AppConfig) LogSafe() map[string]string {
	// TODO: retorne um mapa com as configs para logging,
	// mascarando DATABASE_URL (mostre apenas os primeiros 20 chars + "***")
	panic("implemente LogSafe")
}

// =============================================================================
// EXERCÍCIO 4 — Middleware Stack para K8s (🔴 Difícil)
// =============================================================================
// Implemente uma stack de middlewares typical de produção no K8s:
//   1. Request ID: gera UUID e adiciona ao header X-Request-ID e ao contexto
//   2. Logger: loga método, path, status, duração (structured logging)
//   3. Recover: captura panics e retorna 500 (sem derrubar o servidor)
//   4. Timeout: cancela a request após timeout configurável

type MiddlewareStack struct {
	middlewares []func(http.Handler) http.Handler
}

func NewMiddlewareStack() *MiddlewareStack {
	return &MiddlewareStack{}
}

func (ms *MiddlewareStack) Use(middleware func(http.Handler) http.Handler) *MiddlewareStack {
	ms.middlewares = append(ms.middlewares, middleware)
	return ms
}

func (ms *MiddlewareStack) Apply(handler http.Handler) http.Handler {
	// TODO: aplique os middlewares de trás para frente
	// O primeiro Use() é o mais externo (wrapper mais externo)
	panic("implemente Apply")
}

// RequestIDMiddleware adiciona X-Request-ID único a cada request
func RequestIDMiddleware() func(http.Handler) http.Handler {
	// TODO: gere um ID único (pode usar time.Now().UnixNano() formatado)
	// Adicione ao header de response e ao contexto via context.WithValue
	panic("implemente RequestIDMiddleware")
}

// LoggerMiddleware loga cada request com slog
func LoggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	// TODO: logue: request_id, method, path, status, duration
	panic("implemente LoggerMiddleware")
}

// RecoverMiddleware captura panics e retorna 500
func RecoverMiddleware() func(http.Handler) http.Handler {
	// TODO: use defer/recover, logue o panic, retorne 500
	panic("implemente RecoverMiddleware")
}

// TimeoutMiddleware cancela a request após duration
func TimeoutMiddleware(timeout time.Duration) func(http.Handler) http.Handler {
	// TODO: use context.WithTimeout e o mecanismo de cancelamento do Go
	// Se o handler demorar mais que timeout, retorne 503 "request timeout"
	panic("implemente TimeoutMiddleware")
}

// =============================================================================
// MAIN
// =============================================================================

func testHealthServer() {
	fmt.Println("\n=== HEALTH SERVER ===")
	hs := NewHealthServer("v1.0.0")
	hs.SimulateInit(500 * time.Millisecond)

	client := &http.Client{}
	server := &http.Server{Addr: ":18080", Handler: hs.Handler()}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
			fmt.Println("server error:", err)
		}
	}()

	time.Sleep(100 * time.Millisecond)

	// Deve retornar 503 (ainda inicializando)
	resp, _ := client.Get("http://localhost:18080/health/ready")
	fmt.Printf("Ready imediato: %d (esperado: 503)\n", resp.StatusCode)

	time.Sleep(600 * time.Millisecond)

	// Deve retornar 200
	resp, _ = client.Get("http://localhost:18080/health/ready")
	fmt.Printf("Ready após init: %d (esperado: 200)\n", resp.StatusCode)

	resp, _ = client.Get("http://localhost:18080/metrics")
	fmt.Printf("Metrics: %d\n", resp.StatusCode)

	server.Shutdown(context.Background())
	wg.Wait()
}

func testConfig() {
	fmt.Println("\n=== CONFIG FROM ENVIRONMENT ===")

	// Sem DATABASE_URL — deve falhar
	_, err := LoadConfig()
	if err != nil {
		fmt.Printf("Erro esperado: %v\n", err)
	}

	// Com DATABASE_URL
	os.Setenv("DATABASE_URL", "postgres://user:secret@localhost:5432/db")
	os.Setenv("PORT", "9090")
	os.Setenv("MAX_WORKERS", "20")
	defer os.Unsetenv("DATABASE_URL")
	defer os.Unsetenv("PORT")
	defer os.Unsetenv("MAX_WORKERS")

	cfg, err := LoadConfig()
	if err != nil {
		fmt.Println("Erro inesperado:", err)
		return
	}

	fmt.Printf("Config: port=%s, workers=%d\n", cfg.Port, cfg.MaxWorkers)
	logSafe := cfg.LogSafe()
	fmt.Printf("DB (mascarada): %s\n", logSafe["database_url"])
}

func testGracefulServer() {
	fmt.Println("\n=== GRACEFUL SHUTDOWN ===")
	fmt.Println("Inicie e pressione Ctrl+C para testar graceful shutdown")
	fmt.Println("(Aguardando 3 segundos antes de continuar...)")

	// Para não bloquear o main nos outros testes, apenas valida a estrutura
	gs := NewGracefulServer(":18081", 10*time.Second)
	if gs == nil {
		fmt.Println("ERRO: NewGracefulServer retornou nil")
	} else {
		fmt.Println("GracefulServer criado com sucesso ✅")
	}
}

func main() {
	testHealthServer()
	testConfig()
	testGracefulServer()
	fmt.Println("\n✅ Exercícios do Módulo 22 concluídos!")
}
