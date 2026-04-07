# Módulo 14 — Construindo APIs

[← Testes](../modulo-13-testes/README.md) | [Próximo: gRPC →](../modulo-15-grpc/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que Go é uma das linguagens mais usadas para APIs?
> 2. Qual a diferença entre um Handler e um Service?
> 3. O que acontece com as requests em andamento quando seu servidor recebe SIGTERM?

---

## 📖 Teoria

### 3 abordagens para APIs em Go

| Abordagem | Prós | Contras | Quando usar |
|-----------|------|---------|-------------|
| **`net/http` puro** | Zero dependências, controle total, stdlib estável | Mais verboso, roteamento manual | Aprendizado, APIs simples, bibliotecas |
| **Gin** | Mais popular, middleware rico, validação automática | Dependência externa | APIs REST de produção, equipes grandes |
| **Fiber** | API similar ao Express.js, ultra-rápido | Baseado em fasthttp (incompatível com net/http) | Migrando de Node.js, performance extrema |

> **No mercado**: ~60% das vagas Go pedem experiência com `net/http` puro. Gin aparece em ~30%. Saber **ambos** é o ideal — entenda a stdlib e saiba quando usar framework.

### Arquitetura de API — Camadas de Responsabilidade

```
Request
  │
  ▼
┌─────────────────────────────────────────────────────┐
│ MIDDLEWARE (logging, auth, CORS, rate limit, panic)  │
└──────────────────────┬──────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│ HANDLER (HTTP) — recebe request, valida input,       │
│                  chama service, retorna JSON response │
└──────────────────────┬───────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│ SERVICE (Negócio) — lógica de negócio pura           │
│                     NÃO sabe o que é HTTP            │
│                     recebe e retorna structs Go       │
└──────────────────────┬───────────────────────────────┘
                       │
                       ▼
┌──────────────────────────────────────────────────────┐
│ REPOSITORY (Dados) — acesso a banco, cache, APIs ext │
│                      NÃO sabe de regras de negócio   │
└──────────────────────────────────────────────────────┘
```

**Por que separar?**
- **Handler** muda quando a interface HTTP muda (novo campo no JSON, novo header)
- **Service** muda quando regra de negócio muda (desconto, validação)
- **Repository** muda quando o banco muda (Postgres → DynamoDB)
- Cada camada pode ser **testada isoladamente** com mocks

---

## 1. Graceful Shutdown — Obrigatório em Produção

Em Kubernetes (e qualquer orquestrador), quando seu pod é encerrado:
1. K8s envia **SIGTERM** ao processo
2. Você tem **30 segundos** (padrão) para encerrar
3. Se não encerrar, recebe **SIGKILL** (morte instantânea)

Se você não tratar SIGTERM, requests em andamento são **cortadas no meio** — dados parciais, transações incompletas, clientes recebendo erros.

```go
package main

import (
    "context"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/tarefas", handleTarefas)

    server := &http.Server{
        Addr:              ":8080",
        Handler:           mux,
        ReadHeaderTimeout: 10 * time.Second,  // previne Slowloris attack
        ReadTimeout:       30 * time.Second,
        WriteTimeout:      30 * time.Second,
        IdleTimeout:       120 * time.Second,
    }

    // Inicia o servidor em goroutine separada
    go func() {
        log.Printf("API rodando em %s", server.Addr)
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("servidor falhou: %v", err)
        }
    }()

    // Espera SIGTERM ou SIGINT
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
    sig := <-quit
    log.Printf("recebeu sinal %s, iniciando shutdown...", sig)

    // Dá até 15 segundos para requests em andamento terminarem
    ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
    defer cancel()

    if err := server.Shutdown(ctx); err != nil {
        log.Fatalf("shutdown forçado: %v", err)
    }
    log.Println("servidor encerrado com sucesso")
}
```

**O que `server.Shutdown(ctx)` faz:**
1. Para de aceitar **novas** conexões
2. Espera requests **em andamento** terminarem
3. Se o `ctx` expirar, força o encerramento

> **Pergunta de entrevista**: "Seu serviço Go está no K8s e recebe SIGTERM. O que acontece com as 50 requests em andamento?" A resposta correta envolve `server.Shutdown` com context timeout.

---

## 2. Respostas de Erro Estruturadas

APIs profissionais retornam erros em formato **consistente e previsível**:

```go
// Struct padrão de erro — use em TODA a API
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Detail  string `json:"detail,omitempty"`
}

// Helper que toda a API usa
func writeError(w http.ResponseWriter, status int, msg string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(APIError{
        Code:    status,
        Message: msg,
    })
}

// Helper para sucesso
func writeJSON(w http.ResponseWriter, status int, data any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(data)
}
```

### Validação de Input

```go
func handleCriar(w http.ResponseWriter, r *http.Request) {
    // 1. Limitar tamanho do body (previne DoS)
    r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // 1MB

    // 2. Decodificar
    var dto CriarTarefaDTO
    if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
        writeError(w, http.StatusBadRequest, "JSON inválido")
        return
    }

    // 3. Validar campos obrigatórios
    if strings.TrimSpace(dto.Titulo) == "" {
        writeError(w, http.StatusBadRequest, "campo 'titulo' é obrigatório")
        return
    }
    if len(dto.Titulo) > 200 {
        writeError(w, http.StatusBadRequest, "campo 'titulo' excede 200 caracteres")
        return
    }

    // 4. Chamar service (lógica de negócio)
    tarefa, err := service.Criar(r.Context(), dto)
    if err != nil {
        writeError(w, http.StatusInternalServerError, "erro ao criar tarefa")
        return
    }

    writeJSON(w, http.StatusCreated, tarefa)
}
```

---

## 3. Middleware — A Pilha de Produção

Middleware é **o pattern mais importante** em APIs Go. Todo request passa por uma cadeia de middlewares antes de chegar ao handler:

```go
// Middleware é uma função que envolve um handler
type Middleware func(http.Handler) http.Handler

// Encadear middlewares (o primeiro da lista é o mais externo)
func Chain(handler http.Handler, middlewares ...Middleware) http.Handler {
    for i := len(middlewares) - 1; i >= 0; i-- {
        handler = middlewares[i](handler)
    }
    return handler
}
```

### Middleware de Logging (produção)

```go
func LoggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()

        // Capturar o status code
        rec := &statusRecorder{ResponseWriter: w, statusCode: http.StatusOK}
        next.ServeHTTP(rec, r)

        log.Printf("%s %s %d %v",
            r.Method, r.URL.Path, rec.statusCode, time.Since(start))
    })
}

type statusRecorder struct {
    http.ResponseWriter
    statusCode int
}

func (r *statusRecorder) WriteHeader(code int) {
    r.statusCode = code
    r.ResponseWriter.WriteHeader(code)
}
```

### Middleware de Recover (nunca deixe um panic derrubar o servidor)

```go
func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("PANIC recuperado: %v\n%s", err, debug.Stack())
                writeError(w, http.StatusInternalServerError, "erro interno")
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

### Montando a stack

```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/api/tarefas", handleTarefas)

    handler := Chain(mux,
        RecoverMiddleware,   // 1º: captura panics
        LoggingMiddleware,   // 2º: loga request/response
        CORSMiddleware,      // 3º: headers CORS
    )

    server := &http.Server{Addr: ":8080", Handler: handler}
}
```

---

## 4. Testando APIs com httptest

`net/http/httptest` permite testar handlers **sem subir um servidor real**:

```go
package main

import (
    "net/http"
    "net/http/httptest"
    "strings"
    "testing"
)

func TestHandleListar(t *testing.T) {
    store := NovaStore()
    store.Criar(CriarTarefaDTO{Titulo: "Teste"})

    // Cria request e recorder (simula ResponseWriter)
    req := httptest.NewRequest("GET", "/api/tarefas", nil)
    rec := httptest.NewRecorder()

    store.handleListar(rec, req)

    // Verifica status code
    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; esperado %d", rec.Code, http.StatusOK)
    }

    // Verifica Content-Type
    ct := rec.Header().Get("Content-Type")
    if ct != "application/json" {
        t.Errorf("Content-Type = %s; esperado application/json", ct)
    }

    // Verifica body contém dados
    body := rec.Body.String()
    if !strings.Contains(body, "Teste") {
        t.Errorf("body não contém 'Teste': %s", body)
    }
}

func TestHandleCriar(t *testing.T) {
    store := NovaStore()

    body := strings.NewReader(`{"titulo": "Nova Tarefa", "descricao": "teste"}`)
    req := httptest.NewRequest("POST", "/api/tarefas", body)
    req.Header.Set("Content-Type", "application/json")
    rec := httptest.NewRecorder()

    store.handleCriar(rec, req)

    if rec.Code != http.StatusCreated {
        t.Errorf("status = %d; esperado %d", rec.Code, http.StatusCreated)
    }
}

// Table-driven test para cenários de erro
func TestHandleCriarErros(t *testing.T) {
    tests := []struct {
        name       string
        body       string
        wantStatus int
    }{
        {"json inválido", "{invalido", http.StatusBadRequest},
        {"titulo vazio", `{"titulo": ""}`, http.StatusBadRequest},
        {"titulo faltando", `{"descricao": "só desc"}`, http.StatusBadRequest},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            store := NovaStore()
            req := httptest.NewRequest("POST", "/api/tarefas",
                strings.NewReader(tt.body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            store.handleCriar(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d; esperado %d", rec.Code, tt.wantStatus)
            }
        })
    }
}
```

### Testando middleware

```go
func TestLoggingMiddleware(t *testing.T) {
    // Handler de teste que retorna 200
    inner := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })

    handler := LoggingMiddleware(inner)

    req := httptest.NewRequest("GET", "/test", nil)
    rec := httptest.NewRecorder()

    handler.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("status = %d; esperado 200", rec.Code)
    }
}
```

> **No mercado**: todo projeto Go sério testa APIs com `httptest`. Não é opcional — é como entrevistadores verificam se você sabe testar de verdade.

---

## 5. Patterns de Produção que Diferenciam Sênior de Júnior

### Request ID para rastreabilidade

```go
func RequestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        id := r.Header.Get("X-Request-ID")
        if id == "" {
            id = fmt.Sprintf("%d", time.Now().UnixNano())
        }

        // Propaga no contexto para uso em logs e chamadas downstream
        ctx := context.WithValue(r.Context(), "requestID", id)
        w.Header().Set("X-Request-ID", id)

        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Paginação (toda API real precisa)

```go
type PaginatedResponse struct {
    Data       any   `json:"data"`
    Total      int   `json:"total"`
    Page       int   `json:"page"`
    PerPage    int   `json:"per_page"`
    TotalPages int   `json:"total_pages"`
}

func parsePaginacao(r *http.Request) (page, perPage int) {
    page, _ = strconv.Atoi(r.URL.Query().Get("page"))
    perPage, _ = strconv.Atoi(r.URL.Query().Get("per_page"))

    if page < 1 {
        page = 1
    }
    if perPage < 1 || perPage > 100 {
        perPage = 20
    }
    return
}
```

### Health Check (K8s readiness/liveness)

```go
mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ok"}`))
})

mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
    if err := db.Ping(); err != nil {
        w.WriteHeader(http.StatusServiceUnavailable)
        w.Write([]byte(`{"status":"not ready"}`))
        return
    }
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status":"ready"}`))
})
```

---

## ✅ Checklist de API para Produção

Antes de fazer deploy de uma API Go:

- [ ] **Graceful shutdown** com `server.Shutdown(ctx)` + tratamento de SIGTERM
- [ ] **Timeouts** configurados: `ReadTimeout`, `WriteTimeout`, `IdleTimeout`, `ReadHeaderTimeout`
- [ ] **Middleware de recover** para panics não derrubarem o servidor
- [ ] **Logging estruturado** com request ID, método, path, status, duração
- [ ] **Respostas de erro consistentes** em JSON com código e mensagem
- [ ] **Validação de input** com limites de tamanho (`MaxBytesReader`)
- [ ] **CORS** configurado corretamente (não `*` em produção)
- [ ] **Health checks** `/health/live` e `/health/ready`
- [ ] **Testes com httptest** para todos os endpoints e cenários de erro
- [ ] **Context propagado** de `r.Context()` para services e repositórios

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `api-pura/main.go` | API CRUD completa com `net/http` |
| `api-gin/main.go` | API CRUD completa com Gin |
| `api-fiber/main.go` | API CRUD completa com Fiber |
| `exercicios/ex14_api.go` | 🏋️ Exercícios |

### Como rodar cada API:
```bash
# API pura (porta 8080)
go run curso/modulo-14-apis/api-pura/main.go

# API Gin (porta 8081) — precisa: go get github.com/gin-gonic/gin
go run curso/modulo-14-apis/api-gin/main.go

# API Fiber (porta 8082) — precisa: go get github.com/gofiber/fiber/v2
go run curso/modulo-14-apis/api-fiber/main.go
```

### Testando com curl:
```bash
# Listar todos
curl http://localhost:8080/api/tarefas

# Criar
curl -X POST http://localhost:8080/api/tarefas \
  -H "Content-Type: application/json" \
  -d '{"titulo": "Aprender Go", "descricao": "Terminar o curso"}'

# Buscar por ID
curl http://localhost:8080/api/tarefas/1

# Atualizar
curl -X PUT http://localhost:8080/api/tarefas/1 \
  -H "Content-Type: application/json" \
  -d '{"titulo": "Go Avançado", "concluida": true}'

# Deletar
curl -X DELETE http://localhost:8080/api/tarefas/1
```

---

## 📋 Exercícios

### 🟢 1. API de Notas (net/http puro)
Crie uma API REST para gerenciar notas/anotações:
- `GET /api/notas` — listar todas
- `POST /api/notas` — criar nova
- `GET /api/notas/{id}` — buscar por ID
- `PUT /api/notas/{id}` — atualizar
- `DELETE /api/notas/{id}` — deletar
- `GET /api/notas/busca?q=` — buscar por texto

### 🟡 2. Middleware Stack
Adicione à sua API:
- Logger com status code e duração
- Auth com header `X-API-Key`
- Recover que captura panics e retorna 500

### 🟡 3. Testes com httptest
Escreva testes para **todos** os endpoints da API usando `httptest`:
- Table-driven tests para cenários de sucesso e erro
- Teste cada middleware isoladamente
- Coverage mínimo de 80%

### 🔴 4. API Production-Ready
Combine tudo: API com graceful shutdown, health checks, request ID, paginação, rate limiting e testes. Este é o exercício mais valioso do curso para portfólio.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Testes](../modulo-13-testes/README.md) | [Próximo: gRPC →](../modulo-15-grpc/README.md)
