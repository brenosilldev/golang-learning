# 🌐 Projeto Intermediário 2 — API REST Production-Ready

> **Quando fazer**: após completar os Módulos 01–13
> **Tempo estimado**: 5–7 dias
> **Objetivo**: consolidar concorrência, generics e testes antes dos módulos de produção

---

## Por que este projeto?

Você terminou concorrência, generics e testes — os tópicos mais difíceis do Go. Agora é hora de **construir algo que poderia ir para produção de verdade**. Este projeto serve como portfólio e como ponte para o Módulo 14 (APIs).

A diferença entre um projeto de estudante e um projeto de profissional está nos detalhes que este projeto exige: graceful shutdown, testes com cobertura real, middleware stack, validação correta.

---

## O que você vai construir

Uma API REST de **gerenciamento de links/bookmarks** chamada **`linkvault`**:

```bash
# Criar bookmark
POST /api/v1/bookmarks
{
  "url": "https://go.dev",
  "title": "Go Language",
  "tags": ["golang", "docs"]
}

# Listar com filtro e paginação
GET /api/v1/bookmarks?tag=golang&page=1&per_page=20

# Busca full-text
GET /api/v1/bookmarks/search?q=go+language

# Bookmark individual
GET  /api/v1/bookmarks/:id
PUT  /api/v1/bookmarks/:id
DELETE /api/v1/bookmarks/:id

# Tags
GET /api/v1/tags          ← lista todas as tags com contagem

# Health checks (obrigatório para K8s)
GET /health/live
GET /health/ready

# Métricas (opcional avançado)
GET /metrics
```

---

## Estrutura do Projeto

```
linkvault/
├── cmd/
│   └── api/
│       └── main.go              ← wiring: cria deps e sobe servidor
├── internal/
│   ├── bookmark/
│   │   ├── bookmark.go          ← struct + validação de negócio
│   │   └── bookmark_test.go
│   ├── handler/
│   │   ├── bookmark.go          ← HTTP handlers
│   │   ├── bookmark_test.go     ← testes com httptest
│   │   └── middleware.go        ← logging, recover, cors, rate limit
│   ├── service/
│   │   ├── bookmark.go          ← lógica de negócio
│   │   └── bookmark_test.go     ← testes com mock do repository
│   └── repository/
│       ├── repository.go        ← interface
│       ├── memory.go            ← implementação in-memory (testes)
│       └── memory_test.go
├── pkg/
│   └── validate/
│       ├── validate.go          ← validadores reutilizáveis
│       └── validate_test.go
├── go.mod
├── go.sum
└── Makefile
```

---

## Especificação Técnica

### Estrutura de Domínio

```go
// internal/bookmark/bookmark.go
type Bookmark struct {
    ID          string    `json:"id"`           // UUID
    URL         string    `json:"url"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Tags        []string  `json:"tags"`
    CreatedAt   time.Time `json:"created_at"`
    UpdatedAt   time.Time `json:"updated_at"`
}

type CreateInput struct {
    URL         string   `json:"url"`
    Title       string   `json:"title"`
    Description string   `json:"description"`
    Tags        []string `json:"tags"`
}

func (i CreateInput) Validate() error {
    var errs []error
    if strings.TrimSpace(i.URL) == "" {
        errs = append(errs, &ValidationError{Field: "url", Msg: "obrigatório"})
    }
    if _, err := url.ParseRequestURI(i.URL); err != nil {
        errs = append(errs, &ValidationError{Field: "url", Msg: "formato inválido"})
    }
    if strings.TrimSpace(i.Title) == "" {
        errs = append(errs, &ValidationError{Field: "title", Msg: "obrigatório"})
    }
    return errors.Join(errs...)
}
```

### Interface Repository

```go
// internal/repository/repository.go
type Repository interface {
    Create(ctx context.Context, b *bookmark.Bookmark) error
    FindByID(ctx context.Context, id string) (*bookmark.Bookmark, error)
    List(ctx context.Context, filter Filter) ([]*bookmark.Bookmark, int, error)
    Update(ctx context.Context, b *bookmark.Bookmark) error
    Delete(ctx context.Context, id string) error
    Search(ctx context.Context, query string) ([]*bookmark.Bookmark, error)
    AllTags(ctx context.Context) (map[string]int, error)
}

type Filter struct {
    Tag     string
    Page    int
    PerPage int
}
```

### Middleware Stack Obrigatória

```go
// Ordem de aplicação (externo → interno):
// 1. Recover    ← captura panics, nunca derruba o servidor
// 2. RequestID  ← injeta X-Request-ID para rastreabilidade
// 3. Logger     ← loga método, path, status, duração, request-id
// 4. RateLimit  ← máx 100 req/min por IP (use sync.Map + rate.Limiter)
// 5. CORS       ← headers permissivos em dev, restritivos em prod
```

### Resposta de Erro Padronizada

```go
// TODOS os erros da API usam este formato
type APIError struct {
    Code    int    `json:"code"`
    Message string `json:"message"`
    Details any    `json:"details,omitempty"` // para erros de validação
}

// Exemplo de resposta para validação:
// HTTP 422 Unprocessable Entity
{
    "code": 422,
    "message": "dados inválidos",
    "details": [
        {"field": "url", "message": "obrigatório"},
        {"field": "title", "message": "obrigatório"}
    ]
}
```

---

## Requisitos por Nível

### 🟢 Versão Mínima (obrigatória)
- [ ] CRUD completo com armazenamento in-memory thread-safe (`sync.RWMutex`)
- [ ] Validação de input com erros estruturados
- [ ] Middleware: Recover + Logger + RequestID
- [ ] Paginação em `GET /bookmarks`
- [ ] Health check em `/health/live`
- [ ] Testes de handler com `httptest` — todos os endpoints
- [ ] `go test -race ./...` passa limpo

### 🟡 Versão Intermediária
- [ ] Busca por texto em título e descrição
- [ ] Endpoint `GET /tags` com contagem por tag
- [ ] Rate limiting por IP (100 req/min)
- [ ] CORS configurável via variável de ambiente
- [ ] Health check `/health/ready` que verifica dependências
- [ ] Graceful shutdown com `server.Shutdown(ctx)`
- [ ] Coverage ≥ 75% com `go test -cover`

### 🔴 Versão Avançada
- [ ] Trocar in-memory por SQLite (`modernc.org/sqlite` — sem CGO)
- [ ] Migrations versionadas com `golang-migrate`
- [ ] Cache de listagens com `sync.Map` + TTL (invalida no Create/Update/Delete)
- [ ] Export/Import de bookmarks em JSON e Netscape HTML (formato padrão de bookmarks)
- [ ] Prometheus metrics: total requests, latência p50/p95/p99, bookmarks criados
- [ ] Dockerfile multi-stage + docker-compose com SQLite persistido em volume

---

## Guia de Implementação

### Ordem recomendada (domain-first)

```
1. bookmark.go    → struct + validação (sem HTTP, sem DB)
2. repository.go  → interface + MemoryRepo
3. service.go     → lógica de negócio usando Repository
4. handler.go     → HTTP usando Service
5. middleware.go  → wrappers de handler
6. main.go        → conecta tudo
```

**Por que essa ordem?** Você consegue testar cada camada isoladamente antes de montar a próxima. Quando chegar no handler, o service já está testado.

### Dica para testes de handler

```go
func TestHandlerCreate(t *testing.T) {
    // Cria handler com repository in-memory (sem mock)
    repo := repository.NewMemoryRepo()
    svc := service.NewBookmarkService(repo)
    h := handler.NewBookmarkHandler(svc)

    tests := []struct {
        name       string
        body       string
        wantStatus int
    }{
        {"criação válida", `{"url":"https://go.dev","title":"Go"}`, 201},
        {"url vazia", `{"url":"","title":"Go"}`, 422},
        {"json inválido", `{invalido}`, 400},
        {"url inválida", `{"url":"nao-e-url","title":"Go"}`, 422},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            req := httptest.NewRequest("POST", "/api/v1/bookmarks",
                strings.NewReader(tt.body))
            req.Header.Set("Content-Type", "application/json")
            rec := httptest.NewRecorder()

            h.Create(rec, req)

            if rec.Code != tt.wantStatus {
                t.Errorf("status = %d, want %d\nbody: %s",
                    rec.Code, tt.wantStatus, rec.Body)
            }
        })
    }
}
```

---

## Conceitos dos Módulos 01–13 Usados

| Conceito | Onde aparece |
|----------|-------------|
| **M04** Slices/Maps | `[]string` para tags, `map[string]int` para contagem |
| **M05** Closures | Middleware pattern: `func(http.Handler) http.Handler` |
| **M06** Ponteiros | `*Bookmark` nos repositórios, `*time.Time` opcional |
| **M07** Structs | Todas as structs de domínio, DTOs |
| **M08** Interfaces | `Repository`, `Service` — permite mocks em testes |
| **M09** Erros | `errors.Join` para validação, wrapping por camada |
| **M10** Pacotes | Estrutura `cmd/internal/pkg`, visibilidade correta |
| **M11** Concorrência | `sync.RWMutex` no MemoryRepo, rate limiter |
| **M13** Testes | `httptest`, table-driven, fakes por injeção de deps |

---

## Critérios de Avaliação

- [ ] Todos os endpoints retornam Content-Type `application/json`
- [ ] Erros têm formato consistente (nunca retorna HTML ou texto cru)
- [ ] `go test -race ./...` passa sem data races
- [ ] Trocar MemoryRepo por SQLiteRepo não muda nenhuma linha de service ou handler
- [ ] O servidor faz graceful shutdown (não mata requests em andamento)
- [ ] A API tem `/health/live` e `/health/ready` funcionando

---

## Dica de Portfólio

Este projeto demonstra que você sabe:
- Arquitetura em camadas (handler → service → repository)
- Injeção de dependências sem framework
- Testes com coverage real
- Go idiomático (interfaces, errors, context)

Adicione ao GitHub com README em inglês, instruções `docker-compose up`, e exemplo de uso com `curl`. É o projeto que você vai mencionar em entrevistas quando perguntarem "me mostre código que você escreveu".

> **Próximo passo**: ao terminar, você está pronto para o Módulo 14 (APIs) onde vai adicionar Gin, Fiber e patterns mais avançados.
