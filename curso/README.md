# 🚀 Curso Completo de Go — Do Básico ao Avançado

> **Autor**: Gerado como material de estudo personalizado  
> **Pré-requisito**: Experiência prévia com programação (3+ anos)  
> **Go Version**: 1.21+

---

## 🤔 Por que Go?

Go (ou Golang) foi criada pelo Google em 2009 por **Robert Griesemer**, **Rob Pike** e **Ken Thompson** — lendas da computação (criadores de Unix, UTF-8, etc.). Ela nasceu para resolver problemas reais de escala do Google.

### Motivos para usar Go:

| Característica | Detalhe |
|---|---|
| ⚡ **Performance** | Compilada para código nativo, tão rápida quanto C/C++ em muitos cenários |
| 🧵 **Concorrência nativa** | Goroutines são leves (~2KB de stack) vs threads do OS (~1MB). Milhões de goroutines simultâneas |
| 📦 **Binário único** | Compila para um único executável sem dependências. Deploy trivial |
| 🔒 **Type-safe** | Tipagem estática com inferência de tipos. Erros pegos em compilação |
| 🛠️ **Tooling incrível** | `go fmt`, `go vet`, `go test`, `go doc`, `go mod` — tudo built-in |
| 📖 **Simplicidade** | 25 palavras-chave. Sem herança, sem exceções, sem magia |
| 🏢 **Adoção massiva** | Docker, Kubernetes, Terraform, Prometheus, CockroachDB, Hugo — todos em Go |

### O que posso construir com Go?

- 🌐 **APIs REST e gRPC** — Alta performance, ideal para microserviços
- 🔧 **CLIs** — Ferramentas de linha de comando (Cobra, Viper)
- ☁️ **Cloud & DevOps** — SDKs da AWS, GCP, Azure são Go-first
- 🐳 **Containers** — Docker foi escrito em Go
- 📊 **Sistemas distribuídos** — Consenso, message brokers, databases
- 🕸️ **Web scraping** — Colly, chromedp
- 🤖 **Automação** — Scripts que precisam de performance
- 🎮 **Game servers** — Backends de jogos multiplayer

### Go vs Outras Linguagens

```
Go vs Java     → Menos verboso, sem JVM, compilação mais rápida
Go vs Python   → 10-100x mais rápido, tipagem estática, concorrência real
Go vs Rust     → Mais simples, GC (sem borrow checker), compilação mais rápida
Go vs Node.js  → Concorrência real (não event loop), tipagem estática, melhor performance
```

---

## 📚 Trilha do Curso

### 🟢 Fundamentos (Módulos 01-06)
| # | Módulo | Tópicos |
|---|--------|---------|
| 01 | [Introdução ao Go](./modulo-01-introducao/README.md) | Instalação, primeiro programa, tooling |
| 02 | [Variáveis e Tipos](./modulo-02-variaveis-tipos/README.md) | var, :=, tipos, constantes, zero values |
| 03 | [Controle de Fluxo](./modulo-03-controle-fluxo/README.md) | if, switch, for, range, break/continue |
| 04 | [Coleções](./modulo-04-colecoes/README.md) | Arrays, slices, maps, sets |
| 05 | [Funções](./modulo-05-funcoes/README.md) | Múltiplos retornos, closures, defer, variadic |
| 06 | [Ponteiros](./modulo-06-ponteiros/README.md) | &, *, valor vs referência, nil |

### 🟡 Intermediário (Módulos 07-10)
| # | Módulo | Tópicos |
|---|--------|---------|
| 07 | [Structs e Métodos](./modulo-07-structs/README.md) | Structs, métodos, composição, JSON tags |
| 08 | [Interfaces](./modulo-08-interfaces/README.md) | Polimorfismo, type assertion, stdlib interfaces |
| 09 | [Tratamento de Erros](./modulo-09-tratamento-erros/README.md) | error, custom errors, errors.Is/As, panic/recover |
| 10 | [Pacotes e Módulos](./modulo-10-pacotes-modulos/README.md) | go mod, visibilidade, organização de projeto |

### 🔴 Avançado (Módulos 11-13)
| # | Módulo | Tópicos |
|---|--------|---------|
| 11 | [Concorrência](./modulo-11-concorrencia/README.md) | Goroutines, channels, select, patterns, context |
| 12 | [Generics](./modulo-12-generics/README.md) | Type parameters, constraints, uso prático |
| 13 | [Testes](./modulo-13-testes/README.md) | testing, benchmarks, table-driven, testify |

### 🏗️ APIs & Produção (Módulos 14-17)
| # | Módulo | Tópicos |
|---|--------|---------|
| 14 | [Construindo APIs](./modulo-14-apis/README.md) | net/http puro, Gin, Fiber, middleware, patterns |
| 15 | [gRPC e Protobuf](./modulo-15-grpc/README.md) | Protocol Buffers, server/client, streaming, interceptors |
| 16 | [Banco de Dados](./modulo-16-database/README.md) | database/sql, GORM, transactions, repository pattern |
| 17 | [Docker para Go](./modulo-17-docker/README.md) | Multi-stage build, scratch, compose, CI/CD |

### 🎯 Projeto Final
| # | Módulo | Tópicos |
|---|--------|---------|
| 🔥 | [NexusMQ — Message Broker](./projeto-final/README.md) | Protocolo TCP, WAL, consumer groups, replicação |

### 💼 Carreira
| Arquivo | Conteúdo |
|---------|----------|
| 📊 | [Seu nível e como se posicionar](./CARREIRA.md) | Salários, vagas, portfólio, próximos passos |

---

## 🏃 Como Seguir o Curso

1. **Leia o README** de cada módulo — contém toda a teoria com exemplos
2. **Rode os exemplos** — Copie e execute para ver o comportamento
3. **Faça os exercícios** — Cada módulo tem exercícios progressivos
4. **Não pule módulos** — Cada um depende dos anteriores
5. **No final, faça o projeto** — É o verdadeiro teste do seu aprendizado

```bash
# Para rodar qualquer exemplo
go run curso/modulo-XX/exercicios/exXX.go

# Para rodar testes
go test ./curso/modulo-13-testes/exercicios/...
```

---

> 💡 **Dica**: Go preza pela simplicidade. Se algo parece complicado demais, provavelmente existe uma forma mais simples de fazer. Essa mentalidade é o coração da linguagem.
