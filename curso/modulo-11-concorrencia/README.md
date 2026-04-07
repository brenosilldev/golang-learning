# Módulo 11 — Concorrência em Go

[← Pacotes](../modulo-10-pacotes-modulos/README.md) | [Próximo: Generics →](../modulo-12-generics/README.md)

---

> **Antes de ler**: tente responder sem olhar a resposta.
> 1. Qual a diferença entre paralelo e concorrente?
> 2. O que acontece se duas goroutines escrevem no mesmo mapa ao mesmo tempo?
> 3. Para que serve `context.Context`?
>
> *(guarde suas respostas — você vai encontrá-las ao longo do módulo)*

---

## 🧠 Concorrência vs Paralelismo

```
Concorrência: ESTRUTURA do código (lidar com múltiplas tarefas)
Paralelismo:  EXECUÇÃO simultânea (rodar no mesmo instante)

Concorrente mas NÃO paralelo (1 CPU):
  ──A──A──B──B──A──B──B──A──

Paralelo E concorrente (2 CPUs):
  ──A──A──A──A──
  ──B──B──B──B──
```

> **"Concorrência é sobre estrutura, paralelismo é sobre execução."** — Rob Pike

---

## 1. Goroutines

Goroutines são funções executadas **concorrentemente** pelo runtime Go. Custam ~2KB de stack (vs ~1MB de uma thread OS).

```go
package main

import (
    "fmt"
    "time"
)

func tarefa(id int) {
    fmt.Printf("goroutine %d iniciou\n", id)
    time.Sleep(time.Second)
    fmt.Printf("goroutine %d terminou\n", id)
}

func main() {
    go tarefa(1) // lança em background
    go tarefa(2)
    go tarefa(3)

    // ⚠️ Problema: main pode terminar antes das goroutines!
    time.Sleep(2 * time.Second) // solução ruim — use WaitGroup
}
```

### sync.WaitGroup — Coordenação correta

```go
package main

import (
    "fmt"
    "sync"
)

func worker(id int, wg *sync.WaitGroup) {
    defer wg.Done() // SEMPRE com defer — garante execução mesmo com panic
    fmt.Printf("worker %d trabalhando\n", id)
}

func main() {
    var wg sync.WaitGroup

    for i := 1; i <= 5; i++ {
        wg.Add(1)       // incrementa ANTES de lançar a goroutine
        go worker(i, &wg)
    }

    wg.Wait() // bloqueia até todos chamarem Done()
    fmt.Println("todos terminaram")
}
```

---

## 2. Channels — Comunicação sem compartilhamento de memória

> **Filosofia Go**: *"Don't communicate by sharing memory; share memory by communicating."*

```go
// Unbuffered: sender bloqueia até receiver estar pronto
ch := make(chan int)

// Buffered: sender só bloqueia quando buffer está cheio
ch := make(chan int, 10)

ch <- 42      // enviar
v := <-ch     // receber
close(ch)     // sinaliza que não virão mais valores
```

### Iterando um channel com range

```go
package main

import "fmt"

func gerador(nums ...int) <-chan int { // retorna read-only channel
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out) // ESSENCIAL: sem close, o range nunca termina
    }()
    return out
}

func main() {
    for n := range gerador(2, 3, 5, 7, 11, 13) {
        fmt.Println(n)
    }
}
```

### Direção de channels (type safety)

```go
func producer(out chan<- int) { // só pode enviar
    out <- 42
}

func consumer(in <-chan int) { // só pode receber
    fmt.Println(<-in)
}
```

---

## 3. Select — Multiplexação de Channels

`select` é como um `switch` para channels. Escolhe o case que está pronto, aleatoriamente se múltiplos estiverem.

```go
package main

import (
    "fmt"
    "time"
)

func main() {
    ch1 := make(chan string)
    ch2 := make(chan string)

    go func() {
        time.Sleep(1 * time.Second)
        ch1 <- "um"
    }()
    go func() {
        time.Sleep(2 * time.Second)
        ch2 <- "dois"
    }()

    for i := 0; i < 2; i++ {
        select {
        case msg := <-ch1:
            fmt.Println("ch1:", msg)
        case msg := <-ch2:
            fmt.Println("ch2:", msg)
        case <-time.After(3 * time.Second): // timeout pattern
            fmt.Println("timeout!")
        }
    }
}
```

### Non-blocking com default

```go
select {
case msg := <-ch:
    fmt.Println("recebeu:", msg)
default:
    fmt.Println("nada disponível agora") // não bloqueia
}
```

---

## 4. Padrões de Concorrência

### 🔷 Fan-Out: 1 produtor → N workers

```go
package main

import (
    "fmt"
    "sync"
)

func fanOut(input <-chan int, numWorkers int) []<-chan int {
    outputs := make([]<-chan int, numWorkers)
    for i := 0; i < numWorkers; i++ {
        out := make(chan int)
        outputs[i] = out
        go func(out chan<- int) {
            for v := range input {
                out <- v * v // processa
            }
            close(out)
        }(out)
    }
    return outputs
}

func main() {
    input := make(chan int)
    go func() {
        for i := 1; i <= 10; i++ {
            input <- i
        }
        close(input)
    }()

    workers := fanOut(input, 3)
    
    var wg sync.WaitGroup
    for i, out := range workers {
        wg.Add(1)
        go func(id int, ch <-chan int) {
            defer wg.Done()
            for v := range ch {
                fmt.Printf("worker %d: %d\n", id, v)
            }
        }(i, out)
    }
    wg.Wait()
}
```

### 🔷 Fan-In: N produtores → 1 consumidor

```go
package main

import (
    "fmt"
    "sync"
)

func fanIn(channels ...<-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup

    forward := func(ch <-chan int) {
        defer wg.Done()
        for v := range ch {
            out <- v
        }
    }

    wg.Add(len(channels))
    for _, ch := range channels {
        go forward(ch)
    }

    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

### 🔷 Worker Pool — O padrão mais usado em produção

```go
package main

import (
    "fmt"
    "sync"
    "time"
)

type Job struct {
    ID    int
    Input string
}

type Result struct {
    JobID  int
    Output string
}

func workerPool(numWorkers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup

    for i := 0; i < numWorkers; i++ {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            for job := range jobs {
                // simula trabalho pesado
                time.Sleep(10 * time.Millisecond)
                results <- Result{
                    JobID:  job.ID,
                    Output: fmt.Sprintf("worker%d processou: %s", workerID, job.Input),
                }
            }
        }(i)
    }

    // fecha results quando todos os workers terminarem
    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}

func main() {
    jobs := make(chan Job, 100)
    results := workerPool(5, jobs) // 5 workers concorrentes

    // produz 20 jobs
    go func() {
        for i := 0; i < 20; i++ {
            jobs <- Job{ID: i, Input: fmt.Sprintf("tarefa-%d", i)}
        }
        close(jobs)
    }()

    // consome resultados
    for result := range results {
        fmt.Println(result.Output)
    }
}
```

### 🔷 Pipeline — Estágios encadeados

```go
package main

import "fmt"

// Estágio 1: gera números
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Estágio 2: eleva ao quadrado
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Estágio 3: filtra pares
func filterEven(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

func main() {
    // Pipeline: generate → square → filterEven
    source := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squared := square(source)
    evens := filterEven(squared)

    for n := range evens {
        fmt.Println(n) // 4, 16, 36, 64, 100
    }
}
```

---

## 5. Sincronização com sync

### Mutex — Exclusão mútua

```go
package main

import (
    "fmt"
    "sync"
)

type SafeCounter struct {
    mu    sync.RWMutex
    count map[string]int
}

func (c *SafeCounter) Inc(key string) {
    c.mu.Lock() // lock de escrita (exclusivo)
    defer c.mu.Unlock()
    c.count[key]++
}

func (c *SafeCounter) Value(key string) int {
    c.mu.RLock() // lock de leitura (compartilhado — múltiplos leitores OK)
    defer c.mu.RUnlock()
    return c.count[key]
}

func main() {
    counter := &SafeCounter{count: make(map[string]int)}
    var wg sync.WaitGroup

    for i := 0; i < 1000; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            counter.Inc("requests")
        }()
    }
    wg.Wait()
    fmt.Println("total:", counter.Value("requests")) // sempre 1000
}
```

### sync.Once — Executa exatamente uma vez

```go
var (
    instance *Database
    once     sync.Once
)

func GetDB() *Database {
    once.Do(func() {
        instance = &Database{} // conexão cara — feita só uma vez
    })
    return instance
}
```

### sync.Map — Mapa thread-safe sem mutex manual

```go
var cache sync.Map

cache.Store("key", "value")

if val, ok := cache.Load("key"); ok {
    fmt.Println(val)
}

cache.Range(func(key, value any) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true // continua iterando
})
```

---

## 6. context.Context — Cancelamento e Deadlines

`Context` propaga cancelamento e timeouts pela cadeia de goroutines.

```go
package main

import (
    "context"
    "fmt"
    "time"
)

func buscaDados(ctx context.Context, id int) (string, error) {
    select {
    case <-time.After(2 * time.Second): // simula operação lenta
        return fmt.Sprintf("dados-%d", id), nil
    case <-ctx.Done():
        return "", ctx.Err() // context.DeadlineExceeded ou context.Canceled
    }
}

func main() {
    // Timeout de 1 segundo — buscaDados demora 2s, então vai cancelar
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
    defer cancel() // SEMPRE defer cancel() para liberar recursos

    resultado, err := buscaDados(ctx, 42)
    if err != nil {
        fmt.Println("erro:", err) // context deadline exceeded
        return
    }
    fmt.Println(resultado)
}
```

### Passando valores no contexto (com moderação)

```go
type contextKey string

const userIDKey contextKey = "userID"

// Adicionar ao contexto
ctx = context.WithValue(ctx, userIDKey, "user-123")

// Recuperar do contexto
if userID, ok := ctx.Value(userIDKey).(string); ok {
    fmt.Println("user:", userID)
}
```

> **Regra**: use Context para cancelamento/deadlines/valores de escopo de requisição.
> **NUNCA** para passar parâmetros opcionais de funções.

---

## 7. Armadilhas Comuns (Race Conditions)

```go
// ❌ RACE CONDITION — variável compartilhada sem proteção
counter := 0
var wg sync.WaitGroup
for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        counter++ // CORRIDA! leitura + incremento + escrita não é atômica
    }()
}

// ✅ FIX 1: atomic
var atomicCounter int64
atomic.AddInt64(&atomicCounter, 1)

// ✅ FIX 2: channel como mutex
counter := 0
ch := make(chan int, 1)
ch <- 0 // coloca token inicial
go func() {
    v := <-ch
    ch <- v + 1
}()

// Detectar race conditions em testes:
// go test -race ./...
// go run -race main.go
```

---

## 8. errgroup — O Padrão de Produção para Goroutines Paralelas

`sync.WaitGroup` é a base, mas em produção você quase sempre quer **propagar erros** e **cancelar goroutines irmãs quando uma falha**. Para isso existe `golang.org/x/sync/errgroup`:

```go
package main

import (
    "context"
    "fmt"
    "net/http"

    "golang.org/x/sync/errgroup"
)

func main() {
    urls := []string{
        "https://go.dev",
        "https://pkg.go.dev",
        "https://nao-existe.invalido", // vai falhar
    }

    // errgroup cria um context derivado que cancela quando qualquer goroutine retorna erro
    g, ctx := errgroup.WithContext(context.Background())

    resultados := make([]int, len(urls))

    for i, url := range urls {
        i, url := i, url // captura loop vars (desnecessário a partir do Go 1.22)
        g.Go(func() error {
            req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
            if err != nil {
                return fmt.Errorf("criar request %s: %w", url, err)
            }
            resp, err := http.DefaultClient.Do(req)
            if err != nil {
                return fmt.Errorf("fetch %s: %w", url, err)
            }
            defer resp.Body.Close()
            resultados[i] = resp.StatusCode
            return nil
        })
    }

    // Wait bloqueia até TODAS terminarem e retorna o PRIMEIRO erro
    if err := g.Wait(); err != nil {
        fmt.Println("erro:", err)
    }
    fmt.Println("status codes:", resultados)
}
```

### errgroup vs WaitGroup — Quando usar cada um

| Cenário | Use |
|---------|-----|
| Goroutines que **não retornam erro** (fire-and-forget) | `sync.WaitGroup` |
| Goroutines que **podem falhar** e você quer o primeiro erro | `errgroup` |
| Goroutines que **devem cancelar** as outras se uma falhar | `errgroup.WithContext` |
| **Limitar** o número de goroutines simultâneas | `errgroup.SetLimit(n)` (Go 1.20+) |

```go
// Limitar concorrência com errgroup (substitui semáforo manual)
g := new(errgroup.Group)
g.SetLimit(10) // no máximo 10 goroutines simultâneas

for _, item := range items {
    item := item
    g.Go(func() error {
        return processar(item) // só 10 rodando ao mesmo tempo
    })
}
err := g.Wait()
```

> **No mercado**: `errgroup` aparece em ~80% dos projetos Go de produção. Se um entrevistador perguntar "como você gerencia N goroutines que podem falhar?", a resposta é `errgroup`.

---

## 9. Goroutine Leaks — O Bug Silencioso Mais Comum

Uma goroutine leak acontece quando uma goroutine fica bloqueada para sempre e nunca termina. O garbage collector do Go **não** coleta goroutines — elas ficam consumindo memória até o processo morrer.

### Leak 1: Channel sem receptor

```go
// ❌ LEAK — ninguém lê do channel, goroutine bloqueia para sempre
func fazerAlgo() {
    ch := make(chan int)
    go func() {
        resultado := calcular()
        ch <- resultado // BLOQUEIA PARA SEMPRE se ninguém ler
    }()
    // ch nunca é lido → goroutine nunca termina → LEAK
}

// ✅ FIX: use channel buffered quando o resultado pode ser descartado
func fazerAlgo() {
    ch := make(chan int, 1) // buffer de 1: envio não bloqueia
    go func() {
        ch <- calcular() // não bloqueia mesmo sem receptor
    }()
    // se ninguém ler, a goroutine termina normalmente
}
```

### Leak 2: Goroutine ignorando cancelamento

```go
// ❌ LEAK — goroutine não respeita o context
func buscar(ctx context.Context) {
    go func() {
        for {
            // loop infinito sem checar ctx.Done()
            time.Sleep(time.Second)
            trabalhar()
        }
    }()
}

// ✅ FIX: SEMPRE verifique ctx.Done() em loops
func buscar(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                return // sai limpo
            case <-time.After(time.Second):
                trabalhar()
            }
        }
    }()
}
```

### Leak 3: Producer sem consumer (ou vice-versa)

```go
// ❌ LEAK — producer envia mas ninguém consome
func processar(ctx context.Context) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for i := 0; ; i++ {
            select {
            case out <- i:
            case <-ctx.Done():
                return // SEM ISSO: se o consumer desistir, producer trava
            }
        }
    }()
    return out
}
```

### Como detectar leaks

```bash
# Em testes: use goleak (Uber)
go get go.uber.org/goleak

# No teste:
func TestMain(m *testing.M) {
    goleak.VerifyTestMain(m)
}

# Em runtime: monitore o número de goroutines
runtime.NumGoroutine() // se cresce sem parar, tem leak
```

> **Regra de ouro**: toda goroutine que você cria deve ter um **caminho claro de saída**. Se não tem, você tem um leak.

---

## 10. Channel Axioms — Comportamento que Cai em Entrevista

Channels em Go têm comportamentos específicos que não são óbvios. Decorar essa tabela evita bugs sutis:

| Operação | Channel nil | Channel fechado | Channel aberto |
|----------|------------|-----------------|----------------|
| `ch <- v` (enviar) | **Bloqueia para sempre** | **PANIC** | Bloqueia até ter receptor (ou cabe no buffer) |
| `<-ch` (receber) | **Bloqueia para sempre** | Retorna zero value (sem bloquear) | Bloqueia até ter valor |
| `close(ch)` | **PANIC** | **PANIC** | Fecha normalmente |
| `len(ch)` | 0 | Itens restantes no buffer | Itens no buffer |

### Implicações práticas

```go
// Channel nil bloqueia para sempre — útil para desabilitar cases no select
func mesclar(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for ch1 != nil || ch2 != nil {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil // desabilita este case — select ignora nil channels
                    continue
                }
                out <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

```go
// Receber de channel fechado retorna zero value — use a segunda variável
v, ok := <-ch
if !ok {
    // channel foi fechado, v é zero value do tipo
}

// Fechar channel fechado = PANIC — feche apenas uma vez, e sempre pelo producer
// Nunca feche pelo lado do consumer
```

> **Pergunta clássica de entrevista**: "O que acontece se você enviar em um channel nil?" → Bloqueia para sempre. "E se fechar um channel nil?" → Panic.

---

## 11. Deadlocks — Como Identificar e Evitar

Deadlock acontece quando goroutines esperam umas pelas outras em ciclo e nenhuma consegue progredir.

### Deadlock 1: Goroutine única bloqueada

```go
// ❌ DEADLOCK — main envia mas ninguém recebe
func main() {
    ch := make(chan int)
    ch <- 42 // bloqueia aqui — não tem nenhuma goroutine para receber
    fmt.Println(<-ch)
}
// fatal error: all goroutines are asleep - deadlock!
```

### Deadlock 2: Dependência circular entre channels

```go
// ❌ DEADLOCK — cada goroutine espera a outra
func main() {
    ch1 := make(chan int)
    ch2 := make(chan int)

    go func() {
        v := <-ch1   // espera ch1
        ch2 <- v + 1 // só depois envia em ch2
    }()

    go func() {
        v := <-ch2   // espera ch2
        ch1 <- v + 1 // só depois envia em ch1
    }()

    time.Sleep(time.Second) // ambas goroutines travadas eternamente
}
```

### Deadlock 3: Mutex re-entrante (Go NÃO suporta)

```go
// ❌ DEADLOCK — Go mutex não é re-entrante
var mu sync.Mutex

func A() {
    mu.Lock()
    defer mu.Unlock()
    B() // chama B que tenta Lock de novo → DEADLOCK
}

func B() {
    mu.Lock() // bloqueia para sempre — já está locked pela mesma goroutine
    defer mu.Unlock()
}

// ✅ FIX: separe a lógica interna (sem lock) da pública (com lock)
func A() {
    mu.Lock()
    defer mu.Unlock()
    bInternal() // versão interna sem lock
}

func B() {
    mu.Lock()
    defer mu.Unlock()
    bInternal()
}

func bInternal() {
    // lógica sem lock
}
```

### Como prevenir deadlocks

1. **Sempre adquira locks na mesma ordem** em todo o código
2. **Nunca segure um lock enquanto espera um channel** (ou vice-versa)
3. **Use `context.WithTimeout`** para operações que podem bloquear
4. **Prefira channels a mutexes** quando possível — channels têm detecção de deadlock built-in
5. **Rode `go vet`** — detecta alguns deadlocks em tempo de compilação

---

## 12. Semaphore Pattern — Limitar Concorrência

Quando você precisa limitar o número de goroutines simultâneas sem `errgroup.SetLimit`:

```go
func processarComLimite(items []string, maxConcorrencia int) {
    sem := make(chan struct{}, maxConcorrencia) // semáforo
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)
        sem <- struct{}{} // adquire vaga (bloqueia se cheio)

        go func(item string) {
            defer wg.Done()
            defer func() { <-sem }() // libera vaga

            processar(item)
        }(item)
    }

    wg.Wait()
}
```

> Esse pattern é usado quando: (a) você não quer importar `errgroup`, (b) quer controle fino sobre o semáforo, ou (c) o processamento é fire-and-forget.

---

## ✅ Checklist de Concorrência para Produção

Antes de mandar código concorrente para produção, verifique:

- [ ] Toda goroutine tem **caminho claro de saída** (via `context`, `close(ch)` ou `return`)
- [ ] `go test -race ./...` **passa limpo** (sem data races detectadas)
- [ ] Channels **unbuffered** são usados intencionalmente (não por esquecimento)
- [ ] `context.Context` é propagado da request até a goroutine mais interna
- [ ] Mutexes são adquiridos **sempre na mesma ordem** em todo o código
- [ ] Goroutines de background **logam erros** (não engolem silenciosamente)
- [ ] Worker pools têm **graceful shutdown** (drenam jobs antes de sair)
- [ ] `errgroup` é usado quando goroutines **podem retornar erros**
- [ ] Número de goroutines é **limitado** (semáforo ou `errgroup.SetLimit`)
- [ ] `runtime.NumGoroutine()` é **monitorado** em produção (detecta leaks)

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo11_concorrencia.go` | Workers, channels, select, patterns |
| `exemplos/exemplo11_patterns.go` | Fan-out, fan-in, pipeline completo |
| `exercicios/ex11_concorrencia.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Pipeline de processamento de texto
Implemente um pipeline de 3 estágios: `lê linhas` → `converte para maiúsculo` → `filtra linhas com mais de 10 chars`.

### 🟡 2. Rate Limiter com Ticker
Implemente um rate limiter que processa no máximo N requisições por segundo usando `time.Ticker` e um channel como semáforo.

### 🟡 3. Timeout Gracioso
Crie uma função `buscaComTimeout(ctx, id)` que:
- Tenta buscar dados de 3 "serviços" em paralelo
- Retorna o primeiro que responder
- Cancela os outros via context

### 🟡 4. errgroup com Limite
Use `errgroup` para fazer fetch de 20 URLs com no máximo 5 goroutines simultâneas (`SetLimit`). Retorne o primeiro erro, ou todos os status codes se tudo der certo.

### 🔴 5. Worker Pool com Graceful Shutdown
Estenda o Worker Pool para:
- Aceitar sinal SIGTERM/SIGINT via `os/signal`
- Parar de aceitar novos jobs quando receber o sinal
- Aguardar workers em andamento terminarem antes de encerrar

### 🔴 6. Detector de Goroutine Leaks
Crie uma função `MonitorGoroutines(ctx, intervalo)` que:
- A cada `intervalo`, imprime `runtime.NumGoroutine()`
- Se o número crescer 3x em relação ao inicial, loga alerta
- Para quando `ctx` for cancelado

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Pacotes](../modulo-10-pacotes-modulos/README.md) | [Próximo: Generics →](../modulo-12-generics/README.md)
