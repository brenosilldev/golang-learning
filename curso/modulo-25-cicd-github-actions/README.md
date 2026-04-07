# Módulo 25 — CI/CD com GitHub Actions

[← Performance & Profiling](../modulo-24-performance-profiling/README.md) | [Projeto Final →](../projeto-final/README.md)

---

> **Antes de ler — tente responder:**
> 1. O que é CI/CD e qual o benefício concreto para um projeto Go?
> 2. Por que rodar `go test -race` no CI é diferente de rodar localmente?
> 3. O que é um Docker multi-arch build e por que importa para deploy em ARM (AWS Graviton, Apple Silicon)?

---

## 1. Por que CI/CD é Obrigatório em 2026

```
Sem CI/CD:                    Com CI/CD:
  Dev push → merge → deploy     Dev push → testes automáticos
  (break em prod às 2h da       → lint → build → docker push
   madrugada)                   → deploy automático
                                → rollback automático se falhar
```

**O que o mercado espera:**
- Testes rodando automaticamente em todo PR
- Lint (`staticcheck`, `golangci-lint`) bloqueando merge com código ruim
- Docker image buildada e publicada automaticamente no merge
- Deploy para staging/produção sem intervenção manual
- Coverage report comentado no PR

---

## 2. Estrutura de Arquivos

```
.github/
└── workflows/
    ├── ci.yml          ← roda em todo PR e push
    ├── release.yml     ← roda ao criar tag (deploy)
    └── security.yml    ← scan de vulnerabilidades (agendado)
```

---

## 3. Pipeline CI Completo

```yaml
# .github/workflows/ci.yml
name: CI

on:
  push:
    branches: [main, develop]
  pull_request:
    branches: [main]

# Cancela runs anteriores do mesmo branch/PR (economiza créditos)
concurrency:
  group: ${{ github.workflow }}-${{ github.ref }}
  cancel-in-progress: true

jobs:
  # ─────────────────────────────────────────────
  # Job 1: Lint — verifica qualidade do código
  # ─────────────────────────────────────────────
  lint:
    name: Lint
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true  # cache de módulos Go automático

      - name: golangci-lint
        uses: golangci/golangci-lint-action@v6
        with:
          version: latest
          args: --timeout=5m

  # ─────────────────────────────────────────────
  # Job 2: Test — testa em múltiplas versões Go
  # ─────────────────────────────────────────────
  test:
    name: Test (Go ${{ matrix.go-version }})
    runs-on: ubuntu-latest
    strategy:
      matrix:
        go-version: ["1.21", "1.22"]  # testa compatibilidade
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-go@v5
        with:
          go-version: ${{ matrix.go-version }}
          cache: true

      - name: Download dependencies
        run: go mod download

      - name: Run tests with race detector
        run: go test -race -coverprofile=coverage.out -covermode=atomic ./...

      - name: Check coverage threshold
        run: |
          COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')
          echo "Coverage: ${COVERAGE}%"
          if (( $(echo "$COVERAGE < 70" | bc -l) )); then
            echo "❌ Coverage ${COVERAGE}% is below threshold of 70%"
            exit 1
          fi
          echo "✅ Coverage ${COVERAGE}% meets threshold"

      - name: Upload coverage to Codecov
        uses: codecov/codecov-action@v4
        with:
          file: coverage.out
          token: ${{ secrets.CODECOV_TOKEN }}

  # ─────────────────────────────────────────────
  # Job 3: Build — verifica que compila para todos OS
  # ─────────────────────────────────────────────
  build:
    name: Build
    runs-on: ubuntu-latest
    needs: [lint, test]  # só builda se lint e test passaram
    strategy:
      matrix:
        include:
          - goos: linux
            goarch: amd64
          - goos: linux
            goarch: arm64
          - goos: darwin
            goarch: arm64
          - goos: windows
            goarch: amd64
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
          cache: true

      - name: Build ${{ matrix.goos }}/${{ matrix.goarch }}
        env:
          GOOS: ${{ matrix.goos }}
          GOARCH: ${{ matrix.goarch }}
          CGO_ENABLED: "0"
        run: |
          go build \
            -ldflags="-s -w -X main.version=${{ github.sha }}" \
            -o bin/server-${{ matrix.goos }}-${{ matrix.goarch }} \
            ./cmd/api

      - name: Upload binary
        uses: actions/upload-artifact@v4
        with:
          name: server-${{ matrix.goos }}-${{ matrix.goarch }}
          path: bin/server-${{ matrix.goos }}-${{ matrix.goarch }}
          retention-days: 7
```

---

## 4. Pipeline de Release com Docker

```yaml
# .github/workflows/release.yml
name: Release

on:
  push:
    tags:
      - "v*"  # dispara ao criar tag v1.0.0, v2.1.3, etc.

env:
  REGISTRY: ghcr.io
  IMAGE_NAME: ${{ github.repository }}

jobs:
  docker:
    name: Build & Push Docker
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write  # permissão para push no GHCR

    steps:
      - uses: actions/checkout@v4

      # Setup buildx para multi-arch (linux/amd64 + linux/arm64)
      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      # Login no GitHub Container Registry
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ${{ env.REGISTRY }}
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}

      # Extrai tags automaticamente (v1.2.3 → latest, 1.2.3, 1.2, 1)
      - name: Extract metadata
        id: meta
        uses: docker/metadata-action@v5
        with:
          images: ${{ env.REGISTRY }}/${{ env.IMAGE_NAME }}
          tags: |
            type=semver,pattern={{version}}
            type=semver,pattern={{major}}.{{minor}}
            type=sha,prefix=sha-

      # Build + push multi-arch com cache
      - name: Build and push
        uses: docker/build-push-action@v5
        with:
          context: .
          platforms: linux/amd64,linux/arm64  # suporta AWS Graviton + x86
          push: true
          tags: ${{ steps.meta.outputs.tags }}
          labels: ${{ steps.meta.outputs.labels }}
          cache-from: type=gha          # cache do GitHub Actions
          cache-to: type=gha,mode=max
          build-args: |
            VERSION=${{ github.ref_name }}
            GIT_COMMIT=${{ github.sha }}
```

---

## 5. golangci-lint — Configuração Profissional

```yaml
# .golangci.yml — configuração do linter
linters:
  enable:
    - errcheck        # verificar erros retornados não tratados
    - gosimple        # simplificações de código
    - govet           # problemas comuns (shadowing, printf format, etc.)
    - ineffassign     # atribuições sem uso
    - staticcheck     # análise estática avançada
    - unused          # código não usado
    - gofmt           # formatação obrigatória
    - goimports       # imports organizados
    - gocritic        # sugestões de qualidade
    - noctx           # força uso de context em funções HTTP/DB
    - bodyclose       # garante que resp.Body é fechado
    - sqlcloserows    # garante que sql.Rows é fechado
    - nilerr          # retornar nil quando err != nil
    - exhaustive      # switch deve cobrir todos os cases de enum

linters-settings:
  errcheck:
    check-type-assertions: true  # verifica type assertions
  govet:
    check-shadowing: true
  staticcheck:
    checks: ["all"]

issues:
  exclude-rules:
    # Permite erro ignorado em defer Close (aceitável)
    - linters: [errcheck]
      text: "Error return value of .*(Close|Flush|Write)"
      source: "defer"
```

```bash
# Rodar localmente antes de push
golangci-lint run ./...
golangci-lint run --fix ./...  # corrige o que puder automaticamente
```

---

## 6. Makefile — Automatizando Tarefas Locais

```makefile
# Makefile
.PHONY: all test lint build docker clean

# Variáveis
BINARY_NAME=server
VERSION=$(shell git describe --tags --always --dirty)
LDFLAGS=-ldflags "-s -w -X main.version=$(VERSION)"

all: lint test build

## test: roda testes com race detector e coverage
test:
	go test -race -cover -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report em coverage.html"

## lint: roda todos os linters
lint:
	golangci-lint run ./...
	go vet ./...

## build: compila o servidor
build:
	CGO_ENABLED=0 go build $(LDFLAGS) -o bin/$(BINARY_NAME) ./cmd/api

## docker: builda a imagem Docker
docker:
	docker build \
		--build-arg VERSION=$(VERSION) \
		--build-arg GIT_COMMIT=$(shell git rev-parse --short HEAD) \
		-t $(BINARY_NAME):$(VERSION) \
		-t $(BINARY_NAME):latest \
		.

## dev: roda com live reload (precisa do air)
dev:
	air -c .air.toml

## tidy: limpa e verifica dependências
tidy:
	go mod tidy
	go mod verify

## bench: roda benchmarks
bench:
	go test -bench=. -benchmem -count=3 ./...

## pprof-cpu: profile de CPU por 30s
pprof-cpu:
	go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=30

## clean: remove artefatos
clean:
	rm -rf bin/ coverage.out coverage.html

help:
	@grep -E '^## ' Makefile | sed 's/## //'
```

---

## 7. Scan de Segurança Automático

```yaml
# .github/workflows/security.yml
name: Security

on:
  schedule:
    - cron: "0 6 * * 1"  # todo domingo às 6h
  push:
    branches: [main]

jobs:
  govulncheck:
    name: Vulnerability Check
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Install govulncheck
        run: go install golang.org/x/vuln/cmd/govulncheck@latest

      - name: Check for vulnerabilities
        run: govulncheck ./...

  trivy:
    name: Container Security Scan
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - name: Build image for scanning
        run: docker build -t app:scan .

      - name: Trivy vulnerability scanner
        uses: aquasecurity/trivy-action@master
        with:
          image-ref: "app:scan"
          format: "sarif"
          output: "trivy-results.sarif"
          severity: "CRITICAL,HIGH"
          exit-code: "1"  # falha o job se encontrar vuln crítica

      - name: Upload Trivy scan results
        uses: github/codeql-action/upload-sarif@v3
        with:
          sarif_file: trivy-results.sarif
```

---

## ✅ Checklist de CI/CD para Produção

- [ ] `go test -race ./...` no CI — detecta race conditions que não aparecem localmente
- [ ] `golangci-lint` bloqueando merge com código problemático
- [ ] Coverage threshold configurado (sugestão: 70% mínimo)
- [ ] Build multi-arch (amd64 + arm64) para suportar AWS Graviton
- [ ] Docker multi-stage com cache do GitHub Actions
- [ ] `govulncheck` rodando semanalmente para dependências com CVEs
- [ ] Secrets nunca hardcoded — usar `${{ secrets.NOME }}`
- [ ] `concurrency` configurado para cancelar runs antigas do mesmo branch

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/.github/workflows/ci.yml` | Pipeline CI completo |
| `exemplos/.github/workflows/release.yml` | Pipeline de release com Docker |
| `exemplos/.golangci.yml` | Configuração de linters |
| `exemplos/Makefile` | Tarefas locais automatizadas |
| `exercicios/ex25_cicd.md` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Pipeline Básico
Crie um repositório GitHub com uma API Go simples. Configure um workflow que roda `go test -race ./...` e `go vet ./...` em todo PR. Faça um PR que quebra um teste e verifique que o merge é bloqueado.

### 🟡 2. Coverage Report no PR
Adicione coverage ao pipeline. Configure o Codecov (gratuito para projetos open source). Verifique que o bot comenta o coverage no PR. Adicione threshold de 70% que falha o build.

### 🟡 3. Docker Multi-Arch
Configure o workflow de release que builda imagem para `linux/amd64` e `linux/arm64`. Publique no GitHub Container Registry (GHCR). Verifique com `docker manifest inspect` que ambas as plataformas estão disponíveis.

### 🔴 4. Pipeline Completo com Deploy
Monte o pipeline completo: lint → test → build → docker push → deploy para staging (usando `ssh` ou `kubectl`). Configure ambientes no GitHub (staging, production) com revisão obrigatória antes de produção.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Performance & Profiling](../modulo-24-performance-profiling/README.md) | [Projeto Final →](../projeto-final/README.md)
