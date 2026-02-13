# ========================
# MULTI-STAGE BUILD (recomendado)
# Resultado: ~10MB
# ========================

# --- Estágio 1: Build ---
FROM golang:1.22-alpine AS builder

# Instalar certificados SSL (para HTTPS em runtime)
RUN apk --no-cache add ca-certificates

WORKDIR /app

# Cache de dependências (só rebuilda se go.mod mudar)
COPY go.mod go.sum ./
RUN go mod download

# Copiar código e buildar
COPY . .

# CGO_ENABLED=0 = binário estático (sem glibc)
# -ldflags="-w -s" = remove debug info (menor)
# -ldflags="-X main.version=..." = injeta versão
ARG VERSION=dev
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build -ldflags="-w -s -X main.version=${VERSION}" \
    -o /api .

# --- Estágio 2: Runtime ---
FROM alpine:3.19

# Certificados SSL para requests HTTPS
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

# Binário compilado
COPY --from=builder /api /api

# Usuário não-root (segurança)
RUN adduser -D -g '' appuser
USER appuser

EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/health || exit 1

ENTRYPOINT ["/api"]
