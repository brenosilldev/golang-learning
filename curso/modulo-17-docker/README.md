# Módulo 17 — Docker para Go

[← Banco de Dados](../modulo-16-database/README.md) | [Projeto Final →](../projeto-final/README.md)

---

## 📖 Por que Docker + Go é perfeito?

Go compila para um **binário estático** — sem runtime, sem dependências. Isso significa imagens Docker **absurdamente pequenas**:

| Linguagem | Imagem típica | Com multi-stage |
|-----------|--------------|-----------------|
| Node.js | ~900MB | ~150MB |
| Python | ~800MB | ~100MB |
| Java | ~600MB | ~200MB |
| **Go** | ~800MB | **~5-10MB** 🔥 |

---

## 🔧 Os 3 níveis de Dockerfile para Go

### Nível 1 — Básico (ruim, ~800MB)
```dockerfile
FROM golang:1.22
WORKDIR /app
COPY . .
RUN go build -o main .
CMD ["./main"]
```
❌ Inclui o compilador Go inteiro na imagem final

### Nível 2 — Multi-stage (bom, ~10MB)
```dockerfile
# Estágio 1: Build
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o main .

# Estágio 2: Runtime (só o binário)
FROM alpine:3.19
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /app/main .
EXPOSE 8080
CMD ["./main"]
```
✅ Builder descartado — imagem final só tem Alpine + binário

### Nível 3 — Scratch (perfeito, ~5MB)
```dockerfile
FROM golang:1.22-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o main .

# Imagem FROM SCRATCH — literalmente VAZIA
FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /app/main /main
EXPOSE 8080
ENTRYPOINT ["/main"]
```
✅ `scratch` = imagem vazia, só seu binário
✅ `-ldflags="-w -s"` remove debug info (~30% menor)
✅ `CGO_ENABLED=0` = binário 100% estático

---

## 🐳 Docker Compose — API + Banco

```yaml
version: "3.8"

services:
  api:
    build: .
    ports:
      - "8080:8080"
    environment:
      - DATABASE_URL=postgres://app:secret@db:5432/mydb?sslmode=disable
      - PORT=8080
    depends_on:
      db:
        condition: service_healthy
    restart: unless-stopped

  db:
    image: postgres:16-alpine
    environment:
      POSTGRES_USER: app
      POSTGRES_PASSWORD: secret
      POSTGRES_DB: mydb
    volumes:
      - pgdata:/var/lib/postgresql/data
    ports:
      - "5432:5432"
    healthcheck:
      test: ["CMD-SHELL", "pg_isready -U app"]
      interval: 5s
      timeout: 5s
      retries: 5

volumes:
  pgdata:
```

---

## 🛠️ Dicas avançadas

### Cache de dependências
```dockerfile
# Copie go.mod/go.sum ANTES do código
COPY go.mod go.sum ./
RUN go mod download
# Agora copie o código — só rebuilda se código mudar
COPY . .
```

### Build arguments
```dockerfile
ARG VERSION=dev
RUN go build -ldflags="-X main.version=${VERSION}" -o main .
```
```bash
docker build --build-arg VERSION=1.2.3 -t myapp:1.2.3 .
```

### .dockerignore
```
.git
.env
*.md
Dockerfile
docker-compose.yml
vendor/
tmp/
```

### Health check no Go
```go
http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    w.WriteHeader(http.StatusOK)
    w.Write([]byte(`{"status": "ok"}`))
})
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/Dockerfile` | Multi-stage build otimizado |
| `exemplos/Dockerfile.scratch` | Imagem scratch (~5MB) |
| `exemplos/docker-compose.yml` | API + PostgreSQL |
| `exemplos/main.go` | API simples para dockerizar |
| `exemplos/.dockerignore` | Arquivos para ignorar |
| `exercicios/ex17_docker.md` | 🏋️ Exercícios |

---

[← Banco de Dados](../modulo-16-database/README.md) | [Projeto Final →](../projeto-final/README.md)
