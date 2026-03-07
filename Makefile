# Makefile — golang-learning
# Comandos de desenvolvimento para o curso
#
# Uso:
#   make help          → lista todos os comandos
#   make run-mod11     → roda o exemplo do módulo 11
#   make ex-mod18      → compila exercício do módulo 18
#   make up            → sobe toda a stack Docker
#   make up-kafka      → sobe só Kafka
#   make up-obs        → sobe observabilidade (Prometheus + Grafana + Jaeger)
#   make test-all      → roda todos os testes

.PHONY: help run-% ex-% up up-kafka up-obs up-db down test-all lint clean

# ─────────────────────────────────────────────────────────────────
# AJUDA
# ─────────────────────────────────────────────────────────────────

help: ## Mostra este menu de ajuda
	@echo ""
	@echo "  golang-learning — Comandos disponíveis"
	@echo ""
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ""

# ─────────────────────────────────────────────────────────────────
# RODAR MÓDULOS
# ─────────────────────────────────────────────────────────────────

run-mod01: ## Roda o exemplo do módulo 01 (Introdução)
	go run curso/modulo-01-introducao/exemplos/exemplo01_hello.go

run-mod11: ## Roda o exemplo do módulo 11 (Concorrência)
	go run curso/modulo-11-concorrencia/exemplos/exemplo11_concorrencia.go

# Template genérico para rodar exemplos: make run-modXX
run-%:
	@MODULE=$(shell echo $* | sed 's/mod//') && \
	find curso -name "exemplo*$$MODULE*.go" | head -1 | xargs -I {} go run {}

# ─────────────────────────────────────────────────────────────────
# EXERCÍCIOS (compilação)
# ─────────────────────────────────────────────────────────────────

ex-mod18: ## Compila exercício do módulo 18 (Sistemas Distribuídos)
	go build ./curso/modulo-18-sistemas-distribuidos/exercicios/...

ex-mod19: ## Compila exercício do módulo 19 (Raft)
	go build ./curso/modulo-19-consensus-raft/exercicios/...

ex-mod20: ## Compila exercício do módulo 20 (Message Queues)
	go build ./curso/modulo-20-message-queues/exercicios/...

ex-mod21: ## Compila exercício do módulo 21 (Observabilidade)
	go build ./curso/modulo-21-observabilidade/exercicios/...

ex-mod22: ## Compila exercício do módulo 22 (Cloud Native)
	go build ./curso/modulo-22-cloud-native-k8s/exercicios/...

ex-all: ex-mod18 ex-mod19 ex-mod20 ex-mod21 ex-mod22 ## Compila todos os exercícios 18-22

# ─────────────────────────────────────────────────────────────────
# TESTES
# ─────────────────────────────────────────────────────────────────

test-all: ## Roda todos os testes com -race
	go test -race -count=1 ./...

test-bench: ## Roda benchmarks do módulo 13
	go test -bench=. -benchmem ./curso/modulo-13-testes/...

cover: ## Gera relatório de cobertura
	go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# ─────────────────────────────────────────────────────────────────
# QUALIDADE DE CÓDIGO
# ─────────────────────────────────────────────────────────────────

lint: ## Roda golangci-lint (instale: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest)
	golangci-lint run ./...

fmt: ## Formata todo o código Go
	go fmt ./...

vet: ## Roda go vet
	go vet ./...

check: fmt vet lint ## Roda fmt + vet + lint

# ─────────────────────────────────────────────────────────────────
# DOCKER COMPOSE — STACK LOCAL
# ─────────────────────────────────────────────────────────────────

up: ## Sobe a stack completa (Postgres, NATS, Prometheus, Grafana, Jaeger, etcd)
	docker compose up -d
	@echo ""
	@echo "  ✅ Stack iniciada!"
	@echo ""
	@echo "  📊 Grafana:    http://localhost:3000  (admin/admin)"
	@echo "  📈 Prometheus: http://localhost:9090"
	@echo "  🔍 Jaeger:     http://localhost:16686"
	@echo "  🗄️  PgAdmin:    http://localhost:5050  (admin@golang.dev/admin)"
	@echo "  💬 NATS:       http://localhost:8222   (monitoring)"
	@echo "  🔑 etcd:       http://localhost:2379"
	@echo ""

up-kafka: ## Sobe Kafka + Kafka UI
	docker compose --profile kafka up -d
	@echo ""
	@echo "  ✅ Kafka iniciado!"
	@echo "  📊 Kafka UI: http://localhost:8080"
	@echo ""

up-obs: ## Sobe observabilidade (Prometheus + Grafana + Jaeger)
	docker compose up -d prometheus grafana jaeger
	@echo ""
	@echo "  ✅ Observabilidade iniciada!"
	@echo "  📊 Grafana:    http://localhost:3000  (admin/admin)"
	@echo "  📈 Prometheus: http://localhost:9090"
	@echo "  🔍 Jaeger:     http://localhost:16686"
	@echo ""

up-db: ## Sobe apenas o banco de dados (Postgres)
	docker compose up -d postgres
	@echo "  ✅ Postgres em localhost:5432 (golang/golang123/golang_learning)"

up-nats: ## Sobe apenas NATS
	docker compose up -d nats
	@echo "  ✅ NATS em localhost:4222 | Monitoring: http://localhost:8222"

up-tools: ## Sobe ferramentas (pgAdmin)
	docker compose --profile tools up -d

down: ## Para toda a stack Docker
	docker compose down

clean: ## Para a stack e remove volumes (CUIDADO: apaga dados!)
	docker compose down -v
	@echo "  ⚠️  Volumes removidos. Dados apagados."

# ─────────────────────────────────────────────────────────────────
# UTILITÁRIOS
# ─────────────────────────────────────────────────────────────────

mod-init: ## Inicializa go.mod (necessário para exercícios externos)
	go mod init github.com/brenosilldev/golang-learning

mod-tidy: ## Limpa e atualiza go.sum
	go mod tidy

.DEFAULT_GOAL := help
