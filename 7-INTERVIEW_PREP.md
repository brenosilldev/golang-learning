# 🎯 Preparação para Entrevistas Go — 2026

> Perguntas reais que aparecem em entrevistas técnicas para vagas Go. Organized by: Go-específico, System Design, e LeetCode patterns em Go.

---

## Parte 1 — Perguntas Técnicas Go (mais frequentes)

### 🟢 Nível Júnior/Pleno

**Q: Qual a diferença entre `var x int` e `x := 0`?**
> `var x int` pode ser usado em qualquer escopo (incluindo package-level) e deixa claro o tipo. `:=` só funciona dentro de funções e infere o tipo. Para zero values intencionais, `var x int` é mais idiomático.

**Q: O que é um zero value? Cite 3 exemplos.**
> Cada tipo tem um valor padrão sem precisar inicializar explicitamente: `int` → `0`, `string` → `""`, `bool` → `false`, `*T` → `nil`, `slice` → `nil`, `map` → `nil`.

**Q: Por que Go não tem exceções? Como erros funcionam?**
> Erros são valores de retorno — explícitos, visíveis, sem fluxo de controle "mágico". Isso torna o caminho de erro rastreável. O tipo `error` é uma interface com um único método `Error() string`.

**Q: O que é o `defer` e em que ordem múltiplos defers executam?**
> `defer` adia a execução para quando a função retornar (LIFO — último declarado executa primeiro). Garante limpeza mesmo com `panic` ou `return` antecipado.

**Q: Qual a diferença entre um slice nil e um slice vazio?**
```go
var s1 []int         // nil slice: s1 == nil → true, len=0, cap=0
s2 := []int{}        // empty slice: s2 == nil → false, len=0, cap=0
s3 := make([]int, 0) // empty slice: s3 == nil → false

// append funciona em ambos — mas json.Marshal difere!
json.Marshal(s1) // "null"
json.Marshal(s2) // "[]"
```

**Q: O que acontece quando você faz append além da capacidade?**
> Go aloca um novo array (geralmente 2x a capacidade anterior), copia os elementos e retorna um novo slice. O slice original não é modificado. Por isso SEMPRE reatribua: `s = append(s, v)`.

**Q: Explique goroutines vs threads do OS.**
> Goroutines são leves (~2KB de stack inicial, cresce dinamicamente), gerenciadas pelo runtime Go via M:N scheduling (N goroutines em M threads OS). Uma thread OS pode executar milhares de goroutines. Criar uma goroutine é ~100x mais barato que uma thread OS.

**Q: Qual a diferença entre channel buffered e unbuffered?**
```go
ch1 := make(chan int)    // unbuffered: send bloqueia até haver receiver
ch2 := make(chan int, 5) // buffered: send só bloqueia quando buffer cheio
```
> Unbuffered garante sincronização — sender e receiver se encontram. Buffered desacopla: sender avança sem esperar receiver.

**Q: O que é uma interface em Go? Como ela é implementada?**
> Interface é implícita — qualquer tipo que tem os métodos certos a implementa sem declarar `implements`. Internamente uma interface é `(tipo, ponteiro para valor)`. Uma interface nil tem ambos como nil.

---

### 🟡 Nível Pleno/Sênior

**Q: Explique o race condition do mapa e como corrigir.**
```go
// ❌ concurrent map read and map write → fatal error em runtime
m := map[string]int{}
go func() { m["a"] = 1 }()
go func() { _ = m["a"] }()

// ✅ Opção 1: sync.RWMutex
var mu sync.RWMutex
mu.Lock(); m["a"] = 1; mu.Unlock()

// ✅ Opção 2: sync.Map (ótimo para write-once, read-many)
var sm sync.Map
sm.Store("a", 1)
v, _ := sm.Load("a")
```

**Q: Quando usar goroutines com `errgroup` vs `sync.WaitGroup`?**
> `WaitGroup` quando goroutines são fire-and-forget (não retornam erro). `errgroup` quando podem falhar — propaga o primeiro erro e cancela as demais via context. Em produção, `errgroup` aparece em ~80% dos casos.

**Q: O que é goroutine leak? Como detectar?**
> Goroutine que bloqueia para sempre (channel sem receiver, context nunca cancelado). O GC não coleta goroutines. Detectar: `runtime.NumGoroutine()` crescendo sem parar, ou `go.uber.org/goleak` em testes.

**Q: Explique a regra "accept interfaces, return structs".**
> Aceitar interfaces = flexibilidade para quem chama (pode passar qualquer tipo que satisfaça). Retornar structs = clareza sobre o que foi criado (evita que o chamador precise fazer type assertion). Retornar interfaces esconde informação e dificulta evolução da API.

**Q: O que é context.Context e quando cancelar?**
> Propaga cancelamento, deadlines e valores pela cadeia de goroutines. Sempre `defer cancel()` imediatamente após criar um context derivado. Usar `r.Context()` em handlers HTTP para cancelar queries no banco quando o cliente desconectar.

**Q: Qual a diferença entre `errors.Is` e `errors.As`?**
```go
// errors.Is: "o erro É (ou contém) este valor específico?"
errors.Is(err, sql.ErrNoRows) // verdadeiro mesmo embrulhado com %w

// errors.As: "o erro contém um valor deste TIPO? Me dê ele."
var valErr *ValidationError
errors.As(err, &valErr) // extrai ValidationError da cadeia de wrapping
```

**Q: Por que `interface{}` nil não é nil?**
```go
var p *Pessoa = nil
var i interface{} = p
fmt.Println(i == nil) // FALSE!
// Interface = (tipo, valor). Aqui tipo=*Pessoa, valor=nil → interface não é nil
// Interface só é nil quando AMBOS são nil
```

**Q: Como funciona o escape analysis do compilador Go?**
> O compilador decide se uma variável vai para a stack (rápido, GC grátis) ou heap (GC precisa coletar). Se a variável "escapa" da função (retornada como ponteiro, capturada por closure, enviada para interface), vai para o heap. Ver com: `go build -gcflags="-m" .`

---

### 🔴 Nível Sênior/Staff

**Q: Explique o Go memory model e acontece-antes (happens-before).**
> Go tem um memory model formal. Operações em goroutines diferentes só têm garantia de ordem quando há uma relação happens-before explícita: `sync.Mutex`, `channel send/receive`, `sync.Once`, `WaitGroup.Wait()`. Sem isso, duas goroutines podem ver memória em estados diferentes.

**Q: Como `sync.Once` garante inicialização única thread-safe?**
> Internamente usa atomic + mutex: verifica atomicamente se já foi executado. Se não, adquire mutex, verifica novamente (double-checked locking) e executa. Subsequentes chamadas são apenas uma leitura atômica — quase zero custo.

**Q: O que é work stealing no scheduler Go?**
> O scheduler Go (GMP: Goroutines, Machine threads, Processors) usa work stealing: quando um P (processor) fica sem goroutines, ele "rouba" da fila de outro P. Isso mantém CPUs ocupadas sem precisar de locks globais na fila de goroutines.

**Q: Quando usaria `sync/atomic` vs `sync.Mutex`?**
> `atomic` para operações simples em tipos primitivos (contador, flag booleana) — lock-free, mais rápido. `Mutex` quando a seção crítica envolve múltiplas variáveis ou lógica complexa. Atomic garante atomicidade de uma operação; Mutex garante exclusividade de uma seção.

---

## Parte 2 — System Design (Go-flavored)

### Como responder perguntas de design

```
Framework RESHADED:
R — Requirements (requisitos funcionais e não-funcionais)
E — Estimation (escala: QPS, storage, bandwidth)
S — Storage (qual banco, por quê)
H — High-level design (diagrama de blocos)
A — API design (endpoints/RPCs)
D — Detailed design (componentes críticos)
E — Edge cases (falhas, timeouts, idempotência)
D — Decisões (trade-offs explicados)
```

### Perguntas frequentes + resposta Go-específica

**"Design a rate limiter"**
```go
// Token bucket com golang.org/x/time/rate
type RateLimiter struct {
    limiters map[string]*rate.Limiter
    mu       sync.Mutex
    rate     rate.Limit
    burst    int
}

func (rl *RateLimiter) Allow(key string) bool {
    rl.mu.Lock()
    l, ok := rl.limiters[key]
    if !ok {
        l = rate.NewLimiter(rl.rate, rl.burst)
        rl.limiters[key] = l
    }
    rl.mu.Unlock()
    return l.Allow()
}
// Em produção: Redis com EVAL para rate limiting distribuído
```

**"Design a worker pool"**
```go
func NewWorkerPool(workers int, jobs <-chan Job) <-chan Result {
    results := make(chan Result)
    var wg sync.WaitGroup
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                results <- process(job)
            }
        }()
    }
    go func() { wg.Wait(); close(results) }()
    return results
}
// Mencionar: errgroup.SetLimit como alternativa moderna
```

**"Design a cache"**
> Mencionar: `sync.Map` para read-heavy, LRU com `container/list` + map para bounded cache, TTL com goroutine de limpeza ou `time.AfterFunc`. Em produção: Redis com `go-redis`. Falar sobre cache stampede e singleflight.

**"Como você faria graceful shutdown?"**
```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
<-quit
ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
defer cancel()
server.Shutdown(ctx) // para de aceitar, aguarda em andamento
```

---

## Parte 3 — LeetCode Patterns em Go

> Não memorize soluções. Memorize **patterns**. Cada pattern resolve dezenas de problemas.

### Pattern 1: Two Pointers

```go
// Template
func twoPointers(s []int) {
    left, right := 0, len(s)-1
    for left < right {
        // lógica com s[left] e s[right]
        if condição {
            left++
        } else {
            right--
        }
    }
}

// Exemplo: Two Sum em array ordenado
func twoSum(nums []int, target int) []int {
    l, r := 0, len(nums)-1
    for l < r {
        sum := nums[l] + nums[r]
        if sum == target {
            return []int{l + 1, r + 1}
        } else if sum < target {
            l++
        } else {
            r--
        }
    }
    return nil
}
```

### Pattern 2: Sliding Window

```go
// Template: janela de tamanho variável
func slidingWindow(s string, k int) int {
    freq := make(map[byte]int)
    left, result := 0, 0

    for right := 0; right < len(s); right++ {
        freq[s[right]]++

        for /* condição de janela inválida */ {
            freq[s[left]]--
            if freq[s[left]] == 0 {
                delete(freq, s[left])
            }
            left++
        }
        result = max(result, right-left+1)
    }
    return result
}

// Exemplo: Longest substring without repeating characters
func lengthOfLongestSubstring(s string) int {
    seen := make(map[byte]int) // char → último índice
    left, best := 0, 0
    for right := 0; right < len(s); right++ {
        if idx, ok := seen[s[right]]; ok && idx >= left {
            left = idx + 1
        }
        seen[s[right]] = right
        if right-left+1 > best {
            best = right - left + 1
        }
    }
    return best
}
```

### Pattern 3: HashMap para frequência/lookup O(1)

```go
// Template
func hashMapPattern(nums []int) bool {
    seen := make(map[int]bool)
    for _, n := range nums {
        complement := /* o que estou procurando */
        if seen[complement] {
            return true
        }
        seen[n] = true
    }
    return false
}

// Exemplo: Two Sum
func twoSumHash(nums []int, target int) []int {
    seen := make(map[int]int) // valor → índice
    for i, n := range nums {
        complement := target - n
        if j, ok := seen[complement]; ok {
            return []int{j, i}
        }
        seen[n] = i
    }
    return nil
}
```

### Pattern 4: BFS com fila (grafos, árvores por nível)

```go
// Template BFS
func bfs(start int, graph map[int][]int) []int {
    visited := make(map[int]bool)
    queue := []int{start}
    visited[start] = true
    var result []int

    for len(queue) > 0 {
        node := queue[0]
        queue = queue[1:]
        result = append(result, node)

        for _, neighbor := range graph[node] {
            if !visited[neighbor] {
                visited[neighbor] = true
                queue = append(queue, neighbor)
            }
        }
    }
    return result
}
```

### Pattern 5: DFS com recursão/stack

```go
// Template DFS recursivo
func dfs(node int, graph map[int][]int, visited map[int]bool) {
    if visited[node] {
        return
    }
    visited[node] = true
    // processa node
    for _, neighbor := range graph[node] {
        dfs(neighbor, graph, visited)
    }
}

// Exemplo: Number of Islands
func numIslands(grid [][]byte) int {
    if len(grid) == 0 {
        return 0
    }
    count := 0
    for i := range grid {
        for j := range grid[i] {
            if grid[i][j] == '1' {
                count++
                sinkIsland(grid, i, j)
            }
        }
    }
    return count
}

func sinkIsland(grid [][]byte, i, j int) {
    if i < 0 || i >= len(grid) || j < 0 || j >= len(grid[0]) || grid[i][j] != '1' {
        return
    }
    grid[i][j] = '0'
    sinkIsland(grid, i+1, j)
    sinkIsland(grid, i-1, j)
    sinkIsland(grid, i, j+1)
    sinkIsland(grid, i, j-1)
}
```

### Pattern 6: Binary Search

```go
// Template
func binarySearch(nums []int, target int) int {
    left, right := 0, len(nums)-1
    for left <= right {
        mid := left + (right-left)/2 // evita overflow
        if nums[mid] == target {
            return mid
        } else if nums[mid] < target {
            left = mid + 1
        } else {
            right = mid - 1
        }
    }
    return -1
}

// Variante: encontrar primeiro True em condição monotônica
func firstTrue(n int, condition func(int) bool) int {
    left, right := 0, n
    for left < right {
        mid := left + (right-left)/2
        if condition(mid) {
            right = mid
        } else {
            left = mid + 1
        }
    }
    return left
}
```

### Pattern 7: Heap / Priority Queue

```go
import "container/heap"

// Min-heap de ints
type MinHeap []int
func (h MinHeap) Len() int           { return len(h) }
func (h MinHeap) Less(i, j int) bool { return h[i] < h[j] }
func (h MinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *MinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *MinHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

// Uso — K largest elements
func kLargest(nums []int, k int) []int {
    h := &MinHeap{}
    heap.Init(h)
    for _, n := range nums {
        heap.Push(h, n)
        if h.Len() > k {
            heap.Pop(h) // remove o menor
        }
    }
    return *h // contém os K maiores
}
```

### Pattern 8: Dynamic Programming

```go
// Template: DP com memoização (top-down)
func dpMemo(n int, memo map[int]int) int {
    if n <= 1 {
        return n
    }
    if v, ok := memo[n]; ok {
        return v
    }
    result := dpMemo(n-1, memo) + dpMemo(n-2, memo)
    memo[n] = result
    return result
}

// Template: DP iterativo (bottom-up) — geralmente mais eficiente
func dpIterative(nums []int) int {
    n := len(nums)
    dp := make([]int, n+1)
    dp[0] = /* base case */
    for i := 1; i <= n; i++ {
        dp[i] = /* transição usando dp[i-1], dp[i-2], ... */
        _ = dp[i]
    }
    return dp[n]
}
```

---

## Parte 4 — Checklist Antes da Entrevista

### Semana antes
- [ ] Reler módulos 08 (interfaces), 09 (erros), 11 (concorrência) — os mais perguntados
- [ ] Praticar 2-3 problemas de cada pattern do LeetCode
- [ ] Fazer 1 system design completo (rate limiter, cache, ou URL shortener)

### No dia
- [ ] Falar em voz alta enquanto resolve — entrevistadores avaliam processo, não só resposta
- [ ] Começar com força bruta, depois otimizar — nunca ficar em silêncio
- [ ] Perguntar edge cases antes de codificar: "e se o array for vazio?"
- [ ] Em Go: sempre mencionar tratamento de erro, contexto, e goroutine safety quando relevante

### Perguntas para fazer ao entrevistador
- "Qual o maior desafio técnico que o time enfrentou recentemente?"
- "Como vocês lidam com observabilidade em produção?"
- "Qual o processo de code review?"

---

## Parte 5 — Perguntas Comportamentais (STAR)

> Empresas internacionais usam muito STAR: **S**ituação, **T**arefa, **A**ção, **R**esultado.

Prepare histórias para:
- **"Me fale de um sistema que você construiu do zero"** → NexusMQ, API production-ready
- **"Como você lidou com um bug difícil?"** → race condition, goroutine leak, memory leak
- **"Como você melhora performance?"** → pprof, benchmark antes/depois, qual foi a melhoria
- **"Como você trabalha em equipe remota?"** → PRs, code review, documentação
