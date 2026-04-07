# Módulo 24 — Performance & Profiling

[← Segurança](../modulo-23-seguranca/README.md) | [Próximo: CI/CD →](../modulo-25-cicd-github-actions/README.md)

---

> **Antes de ler — tente responder:**
> 1. O que é o `pprof` e como você acessa o perfil de CPU de um servidor em produção?
> 2. O que é escape analysis e por que importa para performance?
> 3. Qual a diferença entre latência de p50, p95 e p99?

---

## 1. Por que Performance Importa em Go

Go é rápido por design — mas código Go mal escrito pode ser 10-100x mais lento que necessário. As causas mais comuns:

```
Alocações desnecessárias → GC pressure → pausas de latência
Goroutine leaks          → memória crescendo sem parar
Locks de alta contention  → goroutines bloqueadas esperando
Algoritmos errados        → O(n²) onde O(n log n) bastava
```

A regra de ouro:
> **Meça primeiro. Otimize depois. Nunca otimize sem dados.**

---

## 2. Benchmarks — Medindo o Que Importa

```go
package main_test

import (
    "strings"
    "testing"
)

// Benchmark básico
func BenchmarkConcatPlus(b *testing.B) {
    for i := 0; i < b.N; i++ {
        s := ""
        for j := 0; j < 100; j++ {
            s += "hello"  // ❌ aloca nova string a cada iteração
        }
        _ = s
    }
}

func BenchmarkConcatBuilder(b *testing.B) {
    for i := 0; i < b.N; i++ {
        var sb strings.Builder
        sb.Grow(500) // ✅ aloca uma vez
        for j := 0; j < 100; j++ {
            sb.WriteString("hello")
        }
        _ = sb.String()
    }
}

// Executar:
// go test -bench=. -benchmem -count=5
// -benchmem: mostra alocações de memória por operação
// -count=5:  roda 5 vezes para resultados estáveis
```

### Lendo os resultados

```
BenchmarkConcatPlus-8      50000    24100 ns/op    53248 B/op    99 allocs/op
BenchmarkConcatBuilder-8  500000     2150 ns/op      512 B/op     2 allocs/op
                                      ↑               ↑            ↑
                              11x mais rápido    100x menos bytes  50x menos alocações
```

### b.ResetTimer, b.StopTimer, b.StartTimer

```go
func BenchmarkComSetup(b *testing.B) {
    // Setup não conta no benchmark
    dados := gerarDados(10000)
    b.ResetTimer() // reseta após setup

    for i := 0; i < b.N; i++ {
        b.StopTimer()
        copia := clonar(dados) // preparo por iteração não conta
        b.StartTimer()

        processar(copia)       // só isso é medido
    }
}
```

### Comparar dois benchmarks com benchstat

```bash
go test -bench=. -benchmem -count=10 > antes.txt
# ... faz a otimização ...
go test -bench=. -benchmem -count=10 > depois.txt
go install golang.org/x/perf/cmd/benchstat@latest
benchstat antes.txt depois.txt
# Mostra delta percentual com significância estatística
```

---

## 3. pprof — O Perfilador do Go

`pprof` é o profiler built-in do Go. Não precisa instalar nada.

### Tipos de perfil

| Perfil | O que mede | Quando usar |
|--------|-----------|-------------|
| **cpu** | Onde o CPU gasta tempo | Função lenta, alto uso de CPU |
| **heap** | Alocações de memória vivas | Memória crescendo, OOM |
| **goroutine** | Todas as goroutines ativas | Goroutine leaks |
| **allocs** | Todas as alocações (incluindo coletadas) | Muita pressão no GC |
| **mutex** | Contention em locks | Goroutines esperando mutex |
| **block** | Bloqueios em channels/locks | Baixo throughput |

### Profiling em testes (mais fácil)

```bash
# CPU profile
go test -bench=BenchmarkMinhafuncao -cpuprofile=cpu.out
go tool pprof cpu.out

# Memory profile
go test -bench=BenchmarkMinhafuncao -memprofile=mem.out
go tool pprof mem.out

# Interface web interativa (melhor visualização)
go tool pprof -http=:8080 cpu.out
```

### Profiling em servidor HTTP (produção)

```go
package main

import (
    "net/http"
    _ "net/http/pprof" // ← import com side-effect: registra handlers /debug/pprof/
    "log"
)

func main() {
    // ⚠️ NUNCA exponha pprof publicamente — só em rede interna ou com auth
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()

    // ... seu servidor principal ...
}
```

```bash
# Coletar 30 segundos de CPU profile do servidor
go tool pprof http://localhost:6060/debug/pprof/profile?seconds=30

# Goroutines ativas agora
go tool pprof http://localhost:6060/debug/pprof/goroutine

# Heap atual
go tool pprof http://localhost:6060/debug/pprof/heap

# Interface web completa
go tool pprof -http=:8080 http://localhost:6060/debug/pprof/profile?seconds=10
```

### Interpretando o flame graph

```
▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓▓ main.processar (40%)
    ▓▓▓▓▓▓▓▓▓▓▓▓ json.Marshal (25%)
        ▓▓▓▓▓▓ encoding/json.marshal (20%)
    ▓▓▓▓ db.Query (10%)
    ▓▓▓ strings.Builder (5%)

Largura = % do tempo total
Altura = profundidade do call stack
Hotspot = função larga no topo da pilha
```

---

## 4. Escape Analysis — Stack vs Heap

O compilador decide onde cada variável vai. Stack = rápido (sem GC). Heap = GC precisa coletar.

```bash
# Ver quais variáveis escapam para o heap
go build -gcflags="-m" ./...
# ou mais detalhado:
go build -gcflags="-m -m" ./... 2>&1 | grep escape
```

```go
// ✅ Fica na stack (sem alocação)
func soma(a, b int) int {
    result := a + b  // result fica na stack
    return result
}

// ❌ Escapa para heap (alocação)
func novoInt(n int) *int {
    return &n  // &n escapa: go move para heap automaticamente
}

// ❌ Escapa para heap (interface boxing)
func imprimir(v interface{}) {
    fmt.Println(v) // v é boxed em interface → heap
}

// ✅ Não escapa (evita interface{})
func imprimirInt(n int) {
    fmt.Println(n) // fmt.Println recebe ...any, mas int pequeno pode ser inline
}
```

### Reduzindo alocações

```go
// ❌ Aloca nova slice a cada chamada
func filtrar(nums []int) []int {
    var result []int
    for _, n := range nums {
        if n > 0 {
            result = append(result, n)
        }
    }
    return result
}

// ✅ Reutiliza buffer fornecido pelo chamador
func filtrarInto(dst []int, nums []int) []int {
    dst = dst[:0] // reutiliza capacidade sem alocar
    for _, n := range nums {
        if n > 0 {
            dst = append(dst, n)
        }
    }
    return dst
}

// ✅ sync.Pool — reutiliza objetos pesados entre chamadas
var bufferPool = sync.Pool{
    New: func() any { return make([]byte, 0, 4096) },
}

func processarRequest(data []byte) []byte {
    buf := bufferPool.Get().([]byte)
    defer bufferPool.Put(buf[:0])

    buf = append(buf, data...)
    // ... processa ...
    result := make([]byte, len(buf)) // copia o resultado
    copy(result, buf)
    return result
}
```

---

## 5. Métricas de Latência — p50, p95, p99

```
p50 = 50% das requisições terminam abaixo deste valor (mediana)
p95 = 95% das requisições terminam abaixo deste valor
p99 = 99% das requisições terminam abaixo deste valor

Exemplo:
  p50 = 5ms   → a maioria é rápida
  p95 = 50ms  → 1 em 20 é lenta
  p99 = 200ms → 1 em 100 é muito lenta (causa timeouts em cascata)

Por que p99 importa mais:
  Em um sistema com 10 serviços em série, p99 de cada um compõe:
  0.99^10 = 0.90 → 10% das requisições atingem a cauda longa
```

```go
// Implementando um histograma manual simples
type Histogram struct {
    mu      sync.Mutex
    buckets []float64 // limites em ms
    counts  []int64
}

func NewHistogram(buckets []float64) *Histogram {
    return &Histogram{
        buckets: buckets,
        counts:  make([]int64, len(buckets)+1),
    }
}

func (h *Histogram) Observe(ms float64) {
    h.mu.Lock()
    defer h.mu.Unlock()
    for i, b := range h.buckets {
        if ms <= b {
            h.counts[i]++
            return
        }
    }
    h.counts[len(h.buckets)]++ // overflow bucket
}

// Na prática: use Prometheus histogram — já faz isso e exporta p50/p95/p99
```

---

## 6. Otimizações Comuns e Seus Impactos

### 1. Pré-alocar slices quando o tamanho é conhecido

```go
// ❌ Re-aloca várias vezes
var result []string
for _, user := range users {
    result = append(result, user.Nome)
}

// ✅ Aloca uma vez
result := make([]string, 0, len(users))
for _, user := range users {
    result = append(result, user.Nome)
}
```

### 2. strings.Builder para concatenação em loop

```go
// ❌ O(n²) alocações
s := ""
for _, palavra := range palavras {
    s += palavra + " "
}

// ✅ O(n) com uma alocação
var sb strings.Builder
sb.Grow(estimativa)
for _, palavra := range palavras {
    sb.WriteString(palavra)
    sb.WriteByte(' ')
}
s := sb.String()
```

### 3. Evitar conversões string↔[]byte desnecessárias

```go
// ❌ Aloca nova string a cada chamada
func temPrefixo(data []byte, prefix string) bool {
    return strings.HasPrefix(string(data), prefix) // conversão desnecessária
}

// ✅ Sem alocação
func temPrefixo(data []byte, prefix string) bool {
    return bytes.HasPrefix(data, []byte(prefix))
}
// Ou usar: bytes.HasPrefix, bytes.Contains, etc.
```

### 4. Reutilizar objetos JSON

```go
// ❌ Cria novo decoder a cada chamada
func decodificar(r io.Reader) (*Payload, error) {
    var p Payload
    return &p, json.NewDecoder(r).Decode(&p)
}

// ✅ Para hot paths: usar json.Unmarshal em vez de Decoder (menos alocações)
// ✅ Para muito volume: considerar sonic ou jsoniter
import jsoniter "github.com/json-iterator/go"
var json = jsoniter.ConfigCompatibleWithStandardLibrary
```

---

## 7. Go Trace — Visualizando o Runtime

```bash
# Gerar trace de 5 segundos
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out
go tool trace trace.out
# Abre browser com:
# - Goroutine schedule (cada goroutine ao longo do tempo)
# - GC pauses
# - Heap size ao longo do tempo
# - Syscalls e network I/O
```

---

## ✅ Checklist de Performance

- [ ] **Meça antes de otimizar** — tenha benchmark mostrando o problema
- [ ] **pprof CPU** para funções lentas, **pprof heap** para memória crescendo
- [ ] **`-gcflags="-m"`** para ver o que escapa para o heap desnecessariamente
- [ ] **`go test -race`** antes de qualquer otimização concorrente
- [ ] Slices pré-alocados com `make([]T, 0, cap)` quando tamanho é conhecido
- [ ] `strings.Builder` com `Grow()` para concatenação em loop
- [ ] `sync.Pool` para objetos pesados que são criados e descartados frequentemente
- [ ] `benchstat` para comparar melhorias com significância estatística

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo24_bench.go` | Benchmarks comparativos |
| `exemplos/exemplo24_pprof.go` | Servidor com pprof habilitado |
| `exercicios/ex24_performance.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Benchmark de Estruturas
Escreva benchmarks comparando: `map[string]int` vs `[]struct{key, val}` vs `sync.Map` para 1000 entradas com 80% leituras. Qual é mais rápido para read-heavy? Use `-benchmem` para comparar alocações.

### 🟡 2. Encontrar Hotspot com pprof
Escreva um programa que processa 1 milhão de registros JSON. Habilite pprof HTTP. Gere um CPU profile e identifique a função que consome mais CPU. Otimize e compare com `benchstat`.

### 🟡 3. Reduzir Alocações
Implemente uma função `ProcessLines(reader io.Reader) []string` que filtra linhas com mais de 10 chars. Versão 1: simples com alocações. Versão 2: com `bufio.Scanner` e slice pré-alocado. Compare alocações com `-benchmem`.

### 🔴 4. Pool de Objetos
Implemente um `JSONProcessor` que usa `sync.Pool` para reutilizar buffers de parse JSON. Demonstre com benchmark que a versão com Pool reduz alocações em carga alta (1000 requisições simultâneas).

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Segurança](../modulo-23-seguranca/README.md) | [Próximo: CI/CD →](../modulo-25-cicd-github-actions/README.md)
