# 📊 Guia de Estudo — Como Aprender Go de Forma Eficiente

> Estruturado com princípios de neurociência cognitiva aplicados ao aprendizado de programação.

---

## 🧠 Os 7 Princípios Neurocientíficos deste Curso

A pesquisa em neurociência mostra que a maioria das pessoas aprende programação de forma **errada** — lendo passivamente, executando tutoriais passo a passo sem pensar, e confundindo familiaridade com conhecimento real.

Este guia aplica o que funciona de verdade:

| # | Princípio | O que é | Como está no curso |
|---|-----------|---------|-------------------|
| 1 | **Retrieval Practice** | Testar memória é 3x mais eficaz que reler | Perguntas no início de cada módulo |
| 2 | **Spaced Repetition** | Revisar em intervalos crescentes fixa mais | Calendário de revisão abaixo |
| 3 | **Interleaving** | Misturar tópicos melhora retenção de longo prazo | Exercícios que combinam módulos anteriores |
| 4 | **Elaborative Interrogation** | Perguntar "por quê funciona assim?" ativa mais circuitos | Seções "Por que Go faz assim?" nos módulos |
| 5 | **Dificuldade Desejável** | Exercícios um nível acima do conforto criam crescimento | Exercícios 🟡 e 🔴 propositalmente difíceis |
| 6 | **Dual Coding** | Texto + diagrama ativa mais regiões do cérebro | Diagramas ASCII em todo módulo de concorrência |
| 7 | **Efeito de Geração** | Escrever código do zero fixa 5x mais que copiar | Skeletons com TODO — você completa, não copia |

---

## 📅 Protocolo de Spaced Repetition

### Como usar:
Após estudar um módulo, marque a data. Revise nos intervalos abaixo.

```
Dia 0  → Estuda o módulo (lê + faz exercícios)
Dia 1  → Revisão rápida (5 min: tente explicar em voz alta sem olhar)
Dia 4  → Revisão de exercícios (refaz os exercícios sem ver a solução)
Dia 10 → Revisão de conceitos (releia apenas os títulos, reconstrua de cabeça)
Dia 21 → Revisão final (tente fazer um exercício novo do mesmo tópico)
```

> **Por que isso funciona**: o esquecimento não é bug, é feature. Cada vez que você recupera algo da memória com esforço, o traço neural fica mais forte. Revisar antes de esquecer completamente é o ponto ótimo.

---

## 🗓️ Plano de Estudo Semanal (2-3h/dia)

### Semanas 1-2 — Fundamentos (Módulos 01-06)

| Dia | Fazer | Revisar |
|-----|-------|---------|
| Seg | Módulo 01 (Intro) | — |
| Ter | Módulo 02 (Variáveis) | M01 |
| Qua | Módulo 03 (Controle) | M01, M02 |
| Qui | Módulo 04 (Coleções) | M02, M03 |
| Sex | Módulo 05 (Funções) | M03, M04 |
| Sáb | Módulo 06 (Ponteiros) | M04, M05 |
| Dom | **Projeto integrador**: calculadora com todos os conceitos | M01-05 |

### Semanas 3-4 — Intermediário (Módulos 07-10)

| Dia | Fazer | Revisar |
|-----|-------|---------|
| Seg | Módulo 07 (Structs) | M05, M06 |
| Ter | Módulo 08 (Interfaces) | M06, M07 |
| Qua | Módulo 09 (Erros) | M07, M08 |
| Qui | Módulo 10 (Pacotes) | M08, M09 |
| Sex | **Mini projeto**: CLI de gerenciamento de tarefas | M07-10 |
| Sáb | Refatorar o CLI com boas práticas | M07-10 revisão |
| Dom | Descanso — consolide dormindo (memória é fixada no sono) |

### Semanas 5-6 — Avançado (Módulos 11-13)

| Dia | Fazer | Revisar |
|-----|-------|---------|
| Seg | Módulo 11 (Concorrência) — goroutines + channels | M09, M10 |
| Ter | Módulo 11 — padrões (worker pool, pipeline, fan-out) | M11 dia anterior |
| Qua | Módulo 12 (Generics) | M11 |
| Qui | Módulo 13 (Testes) | M12 |
| Sex | **Mini projeto**: scraper concorrente com rate limiting | M11-13 |
| Sáb | Adicionar testes ao mini projeto (coverage 80%+) | M13 |
| Dom | Revisão geral M01-13 — tente explicar tudo em 30 min |

### Semanas 7-9 — Produção (Módulos 14-17)

| Dia | Fazer | Revisar |
|-----|-------|---------|
| Seg-Ter | Módulo 14 (APIs) | M11, M13 |
| Qua | Módulo 15 (gRPC) | M14 |
| Qui | Módulo 16 (Banco de Dados) | M14, M15 |
| Sex | Módulo 17 (Docker) | M16 |
| Sáb | **Mini projeto**: API REST + banco + Docker compose | M14-17 |
| Dom | Deploy local, teste de carga, observação de comportamento |

### Semanas 10-13 — Sistemas Distribuídos (Módulos 18-22)

| Dia | Fazer | Revisar |
|-----|-------|---------|
| Seg | Módulo 18 (Fundamentos Dist.) | M11 (concorrência) |
| Ter | Módulo 19 (Raft) | M18 |
| Qua | Módulo 20 (Kafka/NATS) | M19 |
| Qui | Módulo 21 (Observabilidade) | M17 (Docker) |
| Sex | Módulo 22 (Kubernetes) | M21 |
| Sáb | **Mini projeto**: API com circuit breaker + Prometheus | M18-21 |
| Dom | Deploy no kind (K8s local) com health checks |

### Semanas 14-17 — Projeto Final (NexusMQ)

Veja o README do projeto-final para o cronograma por semana.

---

## ⏱️ Estimativa de Tempo

| Perfil | Ritmo | Total |
|--------|-------|-------|
| Dev experiente (3+ anos), 3h/dia | Intenso | 10-12 semanas |
| Dev pleno (1-2 anos), 2h/dia | Moderado | 16-18 semanas |
| Estudante, 1h/dia | Contínuo | 22-26 semanas |
| Intensivo full-time (8h/dia) | Máximo | 4-6 semanas |

---

## 📈 Curva de Dificuldade

```
Dificuldade
    │
 10 │                                                        ████ NexusMQ
    │                                                   ████
  8 │                                              ████ K8s / etcd
    │                                         ████ OpenTelemetry / Raft
  6 │                                    ████ Kafka / Circuit Breaker
    │                               ████ Concorrência / gRPC
  4 │                          ████ Interfaces / Banco de Dados
    │                     ████ Structs / APIs
  2 │               ████  Funções / Controle de Fluxo
    │          ████ Variáveis / Introdução
  0 └──────────────────────────────────────────────────────────▶
    01   03   05   07   09   11   13   15   17   19   21   NexusMQ
```

### Picos onde a maioria trava (e o que fazer)

| Módulo | Por que trava | Estratégia |
|--------|---------------|-----------|
| **06 — Ponteiros** | Abstrato, sem visualização | Desenhe: caixas de memória com endereços |
| **08 — Interfaces** | Polimorfismo implícito | Implemente 5 tipos diferentes satisfazendo a mesma interface |
| **11 — Concorrência** | Race conditions, deadlocks | Rode com `-race`. Quebre o código. Conserte. Repita. |
| **19 — Raft** | Muitos estados simultâneos | Rode a simulação com logs verbosos, desenhe o diagrama de estados |
| **20 — Kafka** | Setup complexo | Use `docker compose up kafka` do repositório |
| **NexusMQ** | Tudo junto | Implemente uma fase de cada vez. Fase 1 já é portfólio válido. |

---

## 🔬 Protocolo de Estudo por Módulo

**Para cada módulo, siga esta sequência — não pule etapas:**

```
1. ANTES de ler (5 min)
   → Tente responder as perguntas do início do README sem olhar
   → Escreva o que você já sabe no papel

2. LEITURA ATIVA (30-60 min)
   → Leia uma seção → feche o arquivo → explique em voz alta
   → Rode CADA exemplo. Modifique o código. Quebre. Conserte.

3. EXERCÍCIOS (30-60 min)
   → Comece pelo arquivo skeleton (ex. ex11_concorrencia.go)
   → Tente sem ver a solução por 20 minutos antes de buscar ajuda

4. SÍNTESE (10 min)
   → Escreva 3 frases explicando o que aprendeu
   → O que foi surpresa? O que confirmou algo que já sabia?

5. ENSINO (opcional, mas poderoso)
   → Explique o conceito para alguém (colega, chat de devs, ou voz alta)
   → Se travar ao explicar, essa é sua lacuna real
```

---

## 🏋️ Análise dos Exercícios

### Distribuição por dificuldade

| Nível | Quantidade | Onde encontrar |
|-------|-----------|---------------|
| 🟢 Fácil | ~25 exercícios | Módulos 01-06 |
| 🟡 Médio | ~30 exercícios | Módulos 07-17 |
| 🔴 Difícil | ~15 exercícios | Módulos 11, 14, 18-22 |
| ⚫ Expert | 1 projeto | NexusMQ |

### Exercícios mais valiosos para o mercado

| Exercício | Por que é importante |
|-----------|---------------------|
| **11: Worker Pool com Graceful Shutdown** | Padrão em TODO serviço Go de produção |
| **13: Testes + Benchmarks** | Diferencial em entrevistas técnicas |
| **14: API com net/http (sem framework)** | Base de 80% das vagas Go |
| **18: Retry Middleware HTTP** | Skill de SRE fundamental |
| **19: Distributed Lock com etcd** | Aplicável em produção diretamente |
| **20: Dead Letter Queue** | Padrão essencial em event-driven |
| **21: Trace Propagation** | Diferencia sênior de pleno em entrevistas |

---

## ✅ Checklist de Conclusão por Nível

### Nível Júnior (Módulos 01-10)
- [ ] Consigo explicar a diferença entre slice e array
- [ ] Sei implementar uma interface sem consultar documentação
- [ ] Consigo criar e usar custom errors com `errors.Is` e `errors.As`
- [ ] Sei organizar um projeto Go com múltiplos pacotes

### Nível Pleno (Módulos 01-13)
- [ ] Consigo usar goroutines + channels sem race conditions
- [ ] Sei escrever table-driven tests com benchmarks
- [ ] Entendo interfaces implícitas e consigo desenhar o diagrama
- [ ] Sei escrever um Worker Pool funcional do zero

### Nível Sênior (Módulos 01-17)
- [ ] Consigo criar uma API REST com net/http sem framework
- [ ] Sei fazer deploy de uma app Go com Docker multi-stage
- [ ] Entendo o que meu app faz em produção (logs, métricas, health checks)

### Nível Staff (Módulos 01-22 + NexusMQ)
- [ ] Consigo explicar CAP Theorem com trade-offs reais de banco de dados
- [ ] Sei implementar Circuit Breaker e Retry com backoff
- [ ] Instrumentei uma API com Prometheus e OpenTelemetry
- [ ] Entendo como Kubernetes decide quando reiniciar ou remover um pod

---

## 💡 Dicas Reforçadas pela Neurociência

1. **Durma depois de estudar** — a consolidação de memória acontece no sono. Não estude até às 3h da manhã.
2. **Estude em sessões de 25-45 min** com pausas (Pomodoro) — atenção sustentada tem limite físico.
3. **Escreva à mão** conceitos difíceis (ponteiros, goroutines) — ativa mais áreas cerebrais que digitar.
4. **Ensine** — se não consegue explicar, não aprendeu. Use o Rubber Duck Debugging.
5. **Erro é aprendizado** — quando o código quebra e você conserta, o traço neural fica mais forte do que quando funciona de primeira.
6. **Git diário** — veja seu progresso. Dopamina do progresso visível é real.
7. **Não leia sem codar** — familiaridade com o código de alguém não é competência sua.

---

> 🏆 Ao completar todos os 22 módulos + NexusMQ com esse protocolo de estudo, você terá fixado o conhecimento de forma muito mais durável do que seguindo uma trilha passiva. Estima-se 80-90% de retenção de longo prazo vs 10-20% de leitura passiva.
