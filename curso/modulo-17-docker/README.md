# Módulo 17 — Docker para Go

[← Banco de Dados](../modulo-16-database/README.md) | [Próximo: Sistemas Distribuídos →](../modulo-18-sistemas-distribuidos/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que imagens Go com multi-stage build chegam a ~5MB enquanto Node.js fica em ~150MB?
> 2. O que é `CGO_ENABLED=0` e por que é necessário para Docker?
> 3. Por que copiar `go.mod` e `go.sum` ANTES do código fonte aproveita melhor o cache do Docker?

---

## 1. Por que Go + Docker é a Combinação Perfeita

Go compila para um **binário estático** — sem runtime, sem interpreter, sem dependências do sistema. Isso permite imagens Docker **absurdamente pequenas**:

| Linguagem | Imagem típica | Com multi-stage |
|-----------|--------------|-----------------|
| Node.js | ~900MB | ~150MB |
| Python | ~800MB | ~100MB |
| Java | ~600MB | ~200MB |
| **Go** | ~800MB (só o builder) | **~5-15MB** |

**Por que Go consegue imagens tão pequenas?**
- Binário estático — não precisa de runtime
- `CGO_ENABLED=0` — sem dependência de bibliotecas C do sistema
- `-ldflags="-s -w"` — remove debug info e symbol table

---

## 2. Os 3 Níveis de Dockerfile

### Nível 1 — Básico (NÃO faça isso em produção, ~800MB)

```dockerfile
FROM golang:1.22
WORKDIR /app
COPY . .
RUN go build -o server .
CMD ["./server"]
```
❌ Inclui o compilador Go, cache de módulos e fonte inteiro na imagem final.

### Nível 2 — Multi-stage com Alpine (~15MB)

```dockerfile
# ──── Estágio 1: Build ────
FROM golang:1.22-alpine AS builder
WORKDIR /app

# Copie dependências PRIMEIRO para aproveitar o cache do Docker
# Se só o código mudar (não go.mod/go.sum), essa layer fica em cache
COPY go.mod go.sum ./
RUN go mod download

COPY . .
# CGO_ENABLED=0 = binário 100% estático, sem libc do sistema
RUN CGO_ENABLED=0 GOOS=linux go build -o server ./cmd/api

# ──── Estágio 2: Runtime ────
FROM alpine:3.19
RUN apk --no-cache add ca-certificates  # para HTTPS
WORKDIR /app
COPY --from=builder /app/server .
EXPOSE 8080
CMD ["./server"]
```

### Nível 3 — Distroless (produção segura, ~10MB, recomendado)

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w" \          # -s: sem symbol table, -w: sem debug info (~30% menor)
    -o server ./cmd/api

# distroless: sem shell, sem pacotes desnecessários — superfície de ataque mínima
FROM gcr.io/distroless/static-debian12:nonroot
COPY --from=builder /app/server /server
# Não rodar como root — segurança
USER nonroot:nonroot
EXPOSE 8080
ENTRYPOINT ["/server"]
```

### Nível 4 — Scratch (mínimo absoluto, ~5MB)

```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
# Precisa dos certs para HTTPS
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o server ./cmd/api

FROM scratch  # imagem VAZIA — literalmente nada
# Copia certificados CA para HTTPS funcionar
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/server /server
EXPOSE 8080
ENTRYPOINT ["/server"]
```

---

## 3. .dockerignore — O Que NÃO Copiar

```
# .dockerignore — evita enviar arquivos desnecessários para o build context
.git
.github
*.md
.env
.env.*
Dockerfile
docker-compose.yml
docker-compose.*.yml
vendor/
tmp/
*.log
coverage.out
*.test
```

> **Impacto**: sem `.dockerignore`, `docker build` envia tudo para o daemon. Em projetos grandes, isso pode ser centenas de MB — build lento e possível vazamento de credenciais.

---

## 4. Docker Compose — Desenvolvimento Local com Dependências

```yaml
# docker-compose.yml
version: "3.8"

services:
  api:
    build:
      context: .
      dockerfile: Dockerfile
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://app:secret@db:5432/mydb?sslmode=disable
      - PORT=8080
      - ENV=development
    depends_on:
      db:
        condition: service_healthy  # aguarda banco estar pronto
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: mydb
    ports:
      - "5432:5432"
    volumes:
      - pgdata:/var/lib/postgresql/data
      - ./migrations:/docker-entrypoint-initdb.d  # roda scripts na inicialização
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

```bash
docker compose up -d          # sobe tudo em background
docker compose logs -f api    # segue os logs da api
docker compose down           # para e remove containers
docker compose down -v        # para, remove e apaga volumes
docker compose exec api sh    # abre shell no container (use sh, não bash — Alpine)
```

---

## 5. Build Arguments e Versioning

```dockerfile
ARG VERSION=dev
ARG GIT_COMMIT=unknown

RUN CGO_ENABLED=0 go build \
    -ldflags="-s -w -X main.version=${VERSION} -X main.gitCommit=${GIT_COMMIT}" \
    -o server .
```

```go
// No código Go
var (
    version   = "dev"     // sobrescrito pelo build
    gitCommit = "unknown" // sobrescrito pelo build
)

func main() {
    log.Printf("iniciando versão %s (%s)", version, gitCommit)
}
```

```bash
docker build \
  --build-arg VERSION=1.2.3 \
  --build-arg GIT_COMMIT=$(git rev-parse --short HEAD) \
  -t myapp:1.2.3 .
```

---

## 6. Desenvolvimento com Live Reload — Air

Para desenvolvimento local sem rebuild manual a cada alteração:

```bash
# Instalar air (live reload para Go)
go install github.com/air-verse/air@latest

# docker-compose.dev.yml — sobrescreve o compose para dev
```

```yaml
# docker-compose.dev.yml
services:
  api:
    build:
      target: builder  # usa o estágio de build, não o final
    volumes:
      - .:/app         # monta o código local
    command: air       # usa air para live reload
```

```bash
docker compose -f docker-compose.yml -f docker-compose.dev.yml up
```

---

## ✅ Checklist de Docker para Produção

- [ ] **Multi-stage build** — imagem final não contém o compilador
- [ ] **`CGO_ENABLED=0`** — binário estático (funciona em qualquer Linux)
- [ ] **`-ldflags="-s -w"`** — remove debug info para imagem menor
- [ ] **`.dockerignore`** configurado — evita vazar `.env` ou credentials
- [ ] **`go.mod`/`go.sum` copiados antes do código** — cache de dependências
- [ ] **Não rodar como root** — use `USER nonroot` (distroless) ou `adduser`
- [ ] **`COPY --from=builder`** copia apenas o binário final
- [ ] **Health check** configurado no compose e no K8s manifest
- [ ] **Versão e commit** injetados no binário via `-ldflags -X`

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/Dockerfile` | Multi-stage build com Alpine |
| `exemplos/Dockerfile.distroless` | Build com distroless (recomendado) |
| `exemplos/docker-compose.yml` | API + PostgreSQL |
| `exemplos/main.go` | API simples para dockerizar |
| `exemplos/.dockerignore` | Arquivos para ignorar |
| `exercicios/ex17_docker.md` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Dockerizar uma API
Pegue a API do módulo 14 e crie um Dockerfile multi-stage. Compare o tamanho da imagem com e sem multi-stage (`docker images`). Use `.dockerignore` correto.

### 🟡 2. Docker Compose Completo
Crie um `docker-compose.yml` com API + PostgreSQL + Redis. Configure health checks, variáveis de ambiente e volumes. A API deve aguardar o banco estar saudável antes de iniciar.

### 🟡 3. Versioning no Binário
Injete versão e git commit no binário via `-ldflags -X`. Crie um endpoint `/version` que retorna essas informações em JSON. Demonstre que a versão muda conforme o build argument.

### 🔴 4. Build Otimizado com Cache
Configure um pipeline de CI (GitHub Actions ou similar) que:
- Usa cache de layers Docker
- Pusha para um registry (Docker Hub ou GHCR)
- Compara tamanho da imagem antes e depois das otimizações
- Falha o build se a imagem for maior que um limite definido

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Banco de Dados](../modulo-16-database/README.md) | [Próximo: Sistemas Distribuídos →](../modulo-18-sistemas-distribuidos/README.md)
