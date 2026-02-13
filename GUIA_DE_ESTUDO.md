# 📊 Guia de Estudo — Análise Completa do Repositório

> Avaliação detalhada do conteúdo, tempo estimado, dificuldade e recomendações para quem vai estudar.

---

## 🎯 Visão Geral

Este repositório contém **dois cursos completos** que formam um caminho do zero ao profissional:

| Curso | Módulos | Exemplos | Exercícios | Projeto Final | Público-alvo |
|-------|---------|----------|-----------|---------------|-------------|
| 📘 **Go** | 17 | 30+ arquivos | 14 sets | NexusMQ | Devs que querem aprender Go |
| 🌐 **Web3** | 7 | 12 arquivos | 7 sets | ChainPulse | Devs Go que querem entrar em blockchain |

**Total**: 24 módulos, 42+ exemplos executáveis, 21 conjuntos de exercícios, 2 projetos finais.

---

## ⏱️ Estimativa de Tempo

### Curso Go (Total: 6-10 semanas)

| Trilha | Módulos | Tempo estimado | Ritmo |
|--------|---------|---------------|-------|
| 🟢 Fundamentos (01-06) | Variáveis → Ponteiros | **1-2 semanas** | ~2-3h/dia |
| 🟡 Intermediário (07-10) | Structs → Pacotes | **1-2 semanas** | ~2-3h/dia |
| 🔴 Avançado (11-13) | Concorrência, Generics, Testes | **1-2 semanas** | ~3h/dia |
| 🏗️ Produção (14-17) | APIs, gRPC, DB, Docker | **2-3 semanas** | ~3h/dia |
| 🎯 NexusMQ | Projeto Final | **3-4 semanas** | ~4h/dia |

| Perfil do estudante | Tempo total |
|---------------------|-------------|
| Dev experiente (3+ anos), estuda 3h/dia | **6-8 semanas** |
| Dev júnior (1-2 anos), estuda 2h/dia | **10-14 semanas** |
| Estudante, estuda 1h/dia | **16-20 semanas** |
| Intensivo full-time (8h/dia) | **3-4 semanas** |

### Curso Web3 (Total: 4-8 semanas)

| Trilha | Módulos | Tempo estimado | Ritmo |
|--------|---------|---------------|-------|
| 🔐 Fundamentos (01-02) | Cripto, Blockchain do Zero | **1-2 semanas** | ~3h/dia |
| ⛓️ Ethereum (03-04) | ethclient, Smart Contracts | **1-2 semanas** | ~3h/dia |
| 💰 Aplicações (05-06) | Tokens, DeFi | **1 semana** | ~2-3h/dia |
| 🔍 Infraestrutura (07) | Indexer | **1 semana** | ~3h/dia |
| 🎯 ChainPulse | Projeto Final | **2-3 semanas** | ~4h/dia |

### Tempo total (ambos os cursos)

| Perfil | Go | Web3 | Total |
|--------|-----|------|-------|
| Dev experiente, 3h/dia | 6-8 sem | 4-6 sem | **10-14 semanas** |
| Dev júnior, 2h/dia | 10-14 sem | 6-8 sem | **16-22 semanas** |
| Intensivo full-time | 3-4 sem | 2-3 sem | **5-7 semanas** |

---

## 📈 Curva de Dificuldade

```
Dificuldade
    │
 10 │                                              ████ NexusMQ
    │                                         ████
  8 │                                    ████ ChainPulse
    │                               ████
  6 │                          ████ Concorrência / DeFi / Indexer
    │                     ████ Interfaces / Smart Contracts
  4 │                ████ Structs / Ethereum / Tokens
    │           ████ Funções / Criptografia / Blockchain
  2 │      ████ Controle de Fluxo
    │ ████ Introdução / Variáveis
  0 └──────────────────────────────────────────────────────▶
    Go 01  03   05   07   09   11   13  15  17  W3-01  03   05   07  Projetos
```

### Picos de dificuldade (onde a maioria trava)

| Módulo | Por que é difícil | Dica para superar |
|--------|------------------|------------------|
| **Go 06 — Ponteiros** | Conceito abstrato, memória | Desenhe diagramas, rode com debugger |
| **Go 08 — Interfaces** | Polimorfismo implícito | Releia, implemente os exercícios 2x |
| **Go 11 — Concorrência** | Race conditions, deadlocks | Comece com exemplos simples, adicione complexidade |
| **Web3 02 — Blockchain** | Muitas peças se conectando | Rode o exemplo, modifique valores, observe |
| **Web3 04 — Smart Contracts** | Solidity + Go + ferramentas | Siga o setup exatamente, não pule passos |
| **Web3 06 — DeFi** | Matemática financeira + MEV | Foque no AMM primeiro, o resto vem depois |

---

## 🏋️ Análise dos Exercícios

### Distribuição por dificuldade

| Nível | Quantidade | Onde encontrar |
|-------|-----------|---------------|
| 🟢 Fácil | ~25 exercícios | Módulos 01-06, Web3 01 |
| 🟡 Médio | ~30 exercícios | Módulos 07-13, Web3 02-05 |
| 🔴 Difícil | ~15 exercícios | Módulos 14-17, Web3 06-07 |
| ⚫ Expert | 2 projetos | NexusMQ, ChainPulse |

### Exercícios mais valiosos para o aprendizado

| Exercício | Por que é importante |
|-----------|---------------------|
| **11.3 — Worker Pool** | Padrão que você vai usar em TODO projeto Go real |
| **12.1 — Funções genéricas (Unique, Reverse)** | Entende generics de verdade |
| **13.3 — Testes + Benchmarks** | Diferencial em entrevistas |
| **14.1 — API com net/http** | Base de 80% dos jobs Go |
| **W3 1.4 — Multi-Sig Wallet** | Entende assinatura digital na prática |
| **W3 2.5 — Rede P2P** | O exercício mais difícil — simula consenso real |
| **W3 6.3 — Flash Loan** | Entende o conceito mais louco do DeFi |

---

## 🔍 Análise dos Exemplos

### O que roda sem nenhuma dependência

Estes exemplos funcionam com `go run` direto, sem instalar nada:

| Arquivo | Módulo | O que demonstra |
|---------|--------|----------------|
| `curso/modulo-01-*/exemplos/*.go` até `modulo-13-*` | Go 01-13 | Toda a linguagem Go |
| `web3/01-criptografia/exemplos/hashing.go` | Web3 01 | SHA-256, mining |
| `web3/01-criptografia/exemplos/wallet.go` | Web3 01 | ECDSA, assinaturas |
| `web3/01-criptografia/exemplos/merkle.go` | Web3 01 | Merkle Tree |
| `web3/02-blockchain-do-zero/exemplos/blockchain.go` | Web3 02 | Blockchain completa |
| `web3/05-tokens-nfts/exemplos/tokens.go` | Web3 05 | ERC-20, ERC-721 |
| `web3/06-defi/exemplos/amm.go` | Web3 06 | Uniswap simulado |
| `web3/07-indexer/exemplos/indexer.go` | Web3 07 | Indexador de blockchain |

### O que precisa de setup

| Arquivo | Dependência | Setup |
|---------|-------------|-------|
| `curso/modulo-14-*/api-gin/*` | Gin framework | `go get github.com/gin-gonic/gin` |
| `curso/modulo-14-*/api-fiber/*` | Fiber framework | `go get github.com/gofiber/fiber/v2` |
| `curso/modulo-15-*/exemplos/*` | gRPC + protoc | `go get google.golang.org/grpc` |
| `curso/modulo-16-*/exemplos/*` | SQLite driver | `go get github.com/mattn/go-sqlite3` |
| `web3/03-ethereum/exemplos/*` | go-ethereum + Ganache | `go get github.com/ethereum/go-ethereum` |
| `web3/04-smart-contracts/*` | solc + abigen + Ganache  | Ver README do módulo |

---

## 🎯 Análise dos Projetos Finais

### NexusMQ (Curso Go)

| Aspecto | Avaliação |
|---------|-----------|
| **Nível** | Sênior / Staff Engineer |
| **Tempo** | 3-4 semanas (dedicação forte) |
| **Diferencial** | Pouquíssimos devs constroem algo assim — impressiona em entrevistas |
| **Conceitos Go** | Goroutines, channels, TCP, binary protocol, WAL, mutex |
| **Risco** | Alto — pode travar em replicação/consenso. Faça fase por fase. |
| **Impacto no portfólio** | ⭐⭐⭐⭐⭐ — Mostra domínio de sistemas distribuídos |

### ChainPulse (Curso Web3)

| Aspecto | Avaliação |
|---------|-----------|
| **Nível** | Pleno / Sênior Web3 |
| **Tempo** | 2-3 semanas |
| **Diferencial** | Todo protocolo precisa de um explorer — skill comercializável |
| **Conceitos Go** | HTTP, WebSocket, PostgreSQL, goroutines, Docker |
| **Risco** | Médio — depende de setup Ganache/Ethereum funcionar |
| **Impacto no portfólio** | ⭐⭐⭐⭐ — Demonstra Go + blockchain + infra |

---

## 🧭 Caminhos de Estudo Recomendados

### Caminho 1: "Quero ser Go Developer" (Backend/Infra)
```
Go 01-13 → Go 14 (APIs) → Go 16 (DB) → Go 17 (Docker) → NexusMQ
Tempo: ~8-10 semanas | Pular: Web3, gRPC (voltar depois)
```

### Caminho 2: "Quero entrar em Web3"
```
Go 01-13 → Web3 01-07 → ChainPulse
Tempo: ~10-12 semanas | Pular: Go 14-17 (voltar depois)
```

### Caminho 3: "Quero tudo" (mais completo)
```
Go 01-17 → Web3 01-07 → NexusMQ ou ChainPulse
Tempo: ~14-18 semanas | Não pular nada
```

### Caminho 4: "Tenho pressa" (intensivo)
```
Go 01-06 (rápido) → Go 07-11 → Go 14 → Web3 01-02 → Web3 05-06
Tempo: ~4-5 semanas (full-time) | Pular: Generics, Testes, gRPC, Docker, Indexer
```

---

## ✅ Checklist de Conclusão

Use isso para acompanhar seu progresso:

### Curso Go
- [ ] Consigo explicar a diferença entre slice e array
- [ ] Sei usar goroutines + channels sem race conditions
- [ ] Consigo criar uma API REST com net/http sem framework
- [ ] Sei escrever table-driven tests
- [ ] Entendo interfaces implícitas e por que são poderosas
- [ ] Consigo dockerizar um app Go com multi-stage build

### Curso Web3
- [ ] Sei implementar SHA-256 hashing e entendo as propriedades
- [ ] Consigo explicar como ECDSA prova autoria de uma transação
- [ ] Sei como Proof of Work funciona (implementei!)
- [ ] Entendo x*y=k (AMM) e consigo calcular slippage
- [ ] Sei a diferença entre ERC-20 e ERC-721
- [ ] Consigo conectar Go ao Ethereum via ethclient

### Nível alcançado ao concluir tudo

| Se completou... | Seu nível |
|-----------------|-----------|
| Go 01-10 | Júnior Go |
| Go 01-13 + exercícios | Pleno Go |
| Go 01-17 + NexusMQ | Sênior Go |
| Go 01-13 + Web3 01-07 + ChainPulse | Pleno Go + Júnior Web3 |
| Tudo + ambos os projetos | Sênior Go + Pleno Web3 |

---

## 💡 Dicas Finais

1. **Não leia sem codar**. Rode TODOS os exemplos. Modifique. Quebre. Conserte.
2. **Não pule exercícios**. Eles são onde o aprendizado real acontece.
3. **Faça um exercício por dia no mínimo**. Consistência > intensidade.
4. **Cometa no GitHub diariamente**. Mostra progresso e disciplina.
5. **Se travou > 30 minutos**, releia o README e o exemplo correspondente.
6. **O projeto final é negociável**. Se NexusMQ parecer impossível, faça a Fase 1 só.
7. **Web3 sem Go é incompleto**. Termine pelo menos Go 01-11 antes do Web3.

---

> 🏆 Ao completar ambos os cursos, você terá um portfólio com 24 módulos estudados, 70+ exercícios feitos e 2 projetos de nível profissional. Isso te coloca à frente de 90% dos candidatos em processos seletivos para Go e Web3.
