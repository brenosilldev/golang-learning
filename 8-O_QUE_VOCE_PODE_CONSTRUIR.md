# O que você pode construir com este curso

Este documento resume **o que cada módulo habilita na prática**, **como combinar módulos** para projetos reais e **até onde a trilha completa pode levar** — do primeiro `main` até sistemas distribuídos de nível de portfólio sênior.

**Leitura complementar:** [curso/README.md](./curso/README.md) (índice da trilha), [2-GUIA_DE_ESTUDO.md](./2-GUIA_DE_ESTUDO.md) (ritmo e revisão), [curso/CARREIRA.md](./curso/CARREIRA.md) (mercado e posicionamento).

---

## Visão geral do repositório

O repositório organiza **25 módulos** de teoria com exemplos e exercícios, **dois projetos intermediários** (CLI e API) e **um projeto final** (NexusMQ). A progressão segue a ordem numérica: cada módulo assume o conhecimento dos anteriores.

Em termos de **resultado**, você sai com capacidade de:

- Escrever programas Go idiomáticos (tipagem, erros, pacotes, testes).
- Construir **CLIs** e **APIs** que poderiam ir para produção com ajustes de ambiente.
- Entender e aplicar **concorrência**, **observabilidade**, **containers**, **Kubernetes** e **mensageria**.
- Ter um **projeto âncora** (NexusMQ) que demonstra protocolo de rede, persistência e replicação.

---

## Por módulo: o que você passa a conseguir fazer

Abaixo, “o que você pode criar” significa **tipos de programa ou peça de sistema** que fazem sentido com *apenas* aquele módulo mais o que veio antes — não é obrigatório que cada módulo tenha um “produto” isolado; muitos são alicerces.

### Módulo 01 — Introdução ao Go

**Você domina:** instalação, `go run` / `go build`, estrutura `package main` e `func main`, fluxo básico de trabalho com o toolchain.

**Pode criar:** programas de uma única função que imprimem saída, leem argumentos simples e usam a biblioteca padrão para experimentos (por exemplo, “hello world” evoluído com formatação).

**Por que importa:** sem isso, nada do restante roda de forma confortável no seu ambiente.

---

### Módulo 02 — Variáveis e tipos

**Você domina:** `var`, `:=`, tipos numéricos, `string`, `bool`, conversões explícitas, constantes, `iota`, zero values.

**Pode criar:** programas que manipulam dados primitivos com segurança (cálculos, flags booleanas, textos), sem ainda estruturas grandes de dados.

**Combina bem com:** M03 (decidir fluxo com base nesses valores).

---

### Módulo 03 — Controle de fluxo

**Você domina:** `if` com statement inicial, `switch`, `for` (único loop), `range`, `break`/`continue`.

**Pode criar:** menus em texto, simulações simples, parsers rudimentares de linha de comando, loops sobre dados que ainda não estão em estruturas avançadas.

---

### Módulo 04 — Coleções

**Você domina:** arrays, slices (`append`, `cap`/`len`, compartilhamento de backing array), maps, padrão de “set” com `map[T]struct{}`.

**Pode criar:** listas dinâmicas, contagem de frequência, índices em memória, filtros e agregações básicas — base de quase qualquer CLI ou API.

---

### Módulo 05 — Funções

**Você domina:** múltiplos retornos (incluindo `(T, error)`), funções como valores, closures, `defer`, variadic.

**Pode criar:** funções reutilizáveis, pequenas bibliotecas internas, padrão de erro explícito em toda chamada — o estilo “Go real” começa aqui.

---

### Módulo 06 — Ponteiros

**Você domina:** `&`, `*`, passagem eficiente de structs grandes, `nil`, diferença entre cópia e mutação.

**Pode criar:** APIs internas que alteram estado sem copiar structs inteiras; base para métodos com pointer receiver (M07).

---

### Módulo 07 — Structs e métodos

**Você domina:** tipos próprios, métodos, composição em vez de herança, struct tags (ex.: JSON).

**Pode criar:** modelos de domínio (`User`, `Task`, `Order`), validação encapsulada, serialização JSON para arquivos ou HTTP — **objetos de negócio** no estilo Go.

---

### Módulo 08 — Interfaces

**Você domina:** satisfação implícita, polimorfismo, type assertion, `io.Reader`/`io.Writer` e o mantra “accept interfaces, return structs”.

**Pode criar:** camadas desacopladas (ex.: “armazenamento” plugável: arquivo vs memória), testes com doubles, pipelines de I/O com a stdlib.

---

### Módulo 09 — Tratamento de erros

**Você domina:** `error`, erros sentinela, `%w` e cadeias, `errors.Is` / `errors.As`, panic/recover com uso disciplinado.

**Pode criar:** fluxos robustos em CLI e libs, mensagens de falha rastreáveis, APIs que não “engolem” erro.

---

### Módulo 10 — Pacotes e módulos

**Você domina:** `go mod`, imports, visibilidade exportada, layout `cmd/` / `internal/` / `pkg/`.

**Pode criar:** **projetos com mais de um pacote**, bibliotecas reutilizáveis dentro do repo, binários separados (`cmd/...`) — pré-requisito de qualquer aplicação média.

**Marco:** com M01–M10 você está pronto para o projeto **[gotas (CLI)](./curso/projeto-intermediario-cli/README.md)**.

---

### Módulo 11 — Concorrência

**Você domina:** goroutines, channels, `select`, `sync.WaitGroup`, mutex, `context`, padrões (worker pool, pipeline, fan-out/fan-in), graceful shutdown.

**Pode criar:** scrapers paralelos com limite de taxa, jobs em background, servidores que cancelam trabalho por timeout, **throughput** sem threads pesadas do SO.

---

### Módulo 12 — Generics

**Você domina:** type parameters, constraints (`comparable`, `constraints.Ordered`), funções e tipos genéricos.

**Pode criar:** utilitários type-safe (mínimo/máximo, conjuntos, caches genéricos) sem duplicar código nem cair em `any` demais.

---

### Módulo 13 — Testes

**Você domina:** `testing`, table-driven tests, benchmarks, `-race`, coverage, mocks via interfaces.

**Pode criar:** regressão segura em refatorações, métricas de performance, **confiança** para colocar código em CI (M25).

**Marco:** com M01–M13 você está pronto para **[linkvault (API REST)](./curso/projeto-intermediario-api/README.md)**.

---

### Módulo 14 — APIs

**Você domina:** `net/http`, middleware, JSON, camadas handler/service/repository, frameworks (Gin, Fiber) como opção.

**Pode criar:** **APIs REST** consumíveis por frontends ou outros serviços, com roteamento, validação e extensão por middleware.

---

### Módulo 15 — gRPC e Protobuf

**Você domina:** `.proto`, codegen, unary e streaming, interceptors.

**Pode criar:** **contratos binários** entre microserviços, serviços internos de baixa latência, integração com ecossistemas que já falam gRPC (muitos em Go).

---

### Módulo 16 — Banco de dados

**Você domina:** `database/sql`, pool, transações, padrão repository, GORM/sqlx/sqlc no mapa de decisão.

**Pode criar:** serviços com **persistência real** (CRUD, relatórios, consistência transacional).

---

### Módulo 17 — Docker

**Você domina:** multi-stage build, imagens mínimas, `docker compose`, boas práticas de cache de build.

**Pode criar:** **imagens de deploy** reproduzíveis, ambientes locais com app + dependências, base para Kubernetes.

---

### Módulo 18 — Sistemas distribuídos (fundamentos)

**Você domina:** CAP, modelos de consistência, relógios lógicos, **circuit breaker**, **retry com backoff**, idempotência.

**Pode criar:** clientes HTTP resilientes, políticas de retry em workers, desenho consciente de trade-offs em sistemas multi-nó.

---

### Módulo 19 — Consenso e Raft

**Você domina:** papéis follower/candidate/leader, log replication, quorum, uso de **etcd** como sistema real.

**Pode criar:** experimentos de líder eleito, compreensão de como **Kubernetes e etcd** se encaixam; bases para KV coordenado e locks distribuídos (tópicos do próprio módulo).

---

### Módulo 20 — Message queues e event-driven

**Você domina:** Kafka, NATS, padrões CQRS/Saga/outbox, filas e tópicos.

**Pode criar:** **consumidores e produtores** de eventos, pipelines assíncronos, desacoplamento entre serviços (o “sistema não cai tudo junto” quando um parceiro falha).

---

### Módulo 21 — Observabilidade e SRE

**Você domina:** Prometheus, OpenTelemetry, Jaeger, logs estruturados, SLI/SLO.

**Pode criar:** endpoints `/metrics`, traces entre serviços, **dashboards** e alertas — o mínimo esperado em produção séria.

---

### Módulo 22 — Cloud Native e Kubernetes

**Você domina:** manifests, health checks, recursos, Helm, operadores (visão geral), app “K8s-ready”.

**Pode criar:** deploy em cluster, rolling updates, autoscaling orientado a métricas, separação config/código.

---

### Módulo 23 — Segurança

**Você domina:** bcrypt, JWT, RBAC em middleware, TLS/mTLS, rate limiting, OWASP.

**Pode criar:** **APIs autenticadas**, comunicação serviço-a-serviço mais segura, endurecimento contra abusos comuns.

---

### Módulo 24 — Performance e profiling

**Você domina:** `pprof`, benchmarks, escape analysis, `sync.Pool`, percentis de latência.

**Pode criar:** otimizações **baseadas em medição**, redução de alocações, diagnóstico de gargalos em produção.

---

### Módulo 25 — CI/CD (GitHub Actions)

**Você domina:** pipeline de testes, lint, coverage, imagens Docker multi-arch, `govulncheck`.

**Pode criar:** **fluxo de entrega contínua** que impede regressões óbvias e publica artefatos prontos para deploy.

---

## Projetos guiados do repositório

| Projeto | Quando fazer | O que demonstra |
|--------|----------------|-----------------|
| **[gotas](./curso/projeto-intermediario-cli/README.md)** | Após M10 | CLI estruturada, structs, interfaces (ex.: storage), erros, pacotes, I/O |
| **[linkvault](./curso/projeto-intermediario-api/README.md)** | Após M13 | API em camadas, concorrência, generics úteis, testes e shutdown gracioso |
| **[NexusMQ](./curso/projeto-final/README.md)** | Trilha avançada | Broker estilo Kafka: TCP binário, WAL, partições, consumer groups, replicação, admin HTTP/CLI, métricas, K8s |

Cada projeto é um **degrau de portfólio**: o primeiro prova fundamentos sólidos; o segundo prova “código de serviço”; o terceiro prova **arquitetura de sistema**.

---

## Combinações de módulos: o que passa a ficar ao seu alcance

### Bloco A — Fundamentos (M01–M06)

**Síntese:** linguagem procedural completa com dados em memória e funções.

**Você pode construir:** calculadoras, jogos de texto, pequenos utilitários que leem/escrevem pouco, protótipos de lógica sem OO avançada.

---

### Bloco B — Intermediário (M07–M10)

**Síntese:** modelagem + contratos + projeto organizado.

**Você pode construir:** bibliotecas com modelos claros, CLIs ou serviços pequenos **bem fatiados**, código que outra pessoa consegestender e estender.

**+ Projeto gotas:** uma CLI no nível de ferramentas que o ecossistema Go usa no dia a dia.

---

### Bloco C — Avançado core (M11–M13)

**Síntese:** paralelismo/concorrência seguros, reuso type-safe, qualidade verificável.

**Você pode construir:** workers, caches com goroutines, pipelines, qualquer serviço que precise **escalar por I/O** com testes que dão confiança.

**+ Projeto linkvault:** API que se aproxima do que empresas pedem em entrevistas (camadas, health, testes).

---

### Bloco D — Produção backend (M14–M17)

**Síntese:** expor contratos na rede + persistência + empacotar para rodar em qualquer lugar.

**Você pode construir:** **microserviço completo** (HTTP ou gRPC + banco + container), pronto para subir em um ambiente de staging.

**Exemplos de “stack mínima produtiva”:** REST + Postgres + Docker Compose; ou gRPC + SQL + imagem distroless.

---

### Bloco E — Sistemas distribuídos (M18–M22)

**Síntese:** falhas parciais, mensagens, métricas/traces, orquestração.

**Você pode construir:** plataforma de **vários serviços** com filas, políticas de resiliência, observabilidade e deploy em Kubernetes — o perfil **SRE / backend sênior** em muitas vagas.

---

### Bloco F — Endurecimento (M23–M25)

**Síntese:** segurança, performance com dados, automação de qualidade.

**Você pode construir:** o mesmo sistema dos blocos anteriores, porém **mais difícil de quebrar, invadir ou derrubar por regressão** — pronto para equipe e para auditoria básica de pipeline.

---

## O que você pode alcançar no fim da trilha (visão de carreira)

Esta trilha não é só “saber Go”: ela cobre o **pacote típico de um backend Go em produção** e ainda caminha em direção a **sistemas distribuídos** (mensageria, consenso, K8s).

De forma realista:

- **Com M01–M10 + gotas:** posicionamento forte como quem domina a linguagem e organização de projeto (junior/pleno entrando em Go).
- **Com M01–M13 + linkvault:** candidato a **pleno** em times que valorizam testes e API bem estruturada.
- **Com M14–17:** capacidade de entregar **serviço com banco e container** — expectativa de **pleno/sênior** dependendo da profundidade das práticas.
- **Com M18–22 + NexusMQ (ou parte dele):** narrativa de **sênior/staff** em entrevistas de system design orientadas a dados e infra.

O arquivo [2-GUIA_DE_ESTUDO.md](./2-GUIA_DE_ESTUDO.md) amarra isso a um **cronograma** e a checklists por nível (júnior, pleno, sênior, staff).

---

## Mapa rápido: “preciso de quais módulos para…?”

| Objetivo | Módulos principais |
|----------|-------------------|
| CLI profissional | M01–M10, opcional M13 |
| API REST | M01–M14, M16 |
| Microserviços internos (gRPC) | M01–M15, M16 |
| App em Docker / nuvem | M14–M17, M22 |
| Filas e eventos | M11, M14, M18, M20 |
| Métricas e traces | M14, M21 |
| Cluster Kubernetes | M17, M21–M22 |
| Auth e API segura | M14, M23 |
| Performance em produção | M13, M24 |
| Time com CI/CD | M13, M25 |
| Broker / log replicado (NexusMQ) | M01–M22 (ênfase M11, M18–M20), M21–M25 para fechar |

---

## Próximo passo

Abra o [README do curso](./curso/README.md) e siga os módulos em ordem; use este documento como **mapa de motivação** — quando souber *o que* cada etapa desbloqueia, fica mais fácil manter ritmo e escolher mini-projetos alinhados ao que você já estudou.
