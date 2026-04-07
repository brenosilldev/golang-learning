# ⚠️ Erros Comuns em Go — Do Iniciante ao Avançado

> Os erros que TODO mundo comete ao aprender Go. Leia antes de sofrer.

---

## Parte 1 — Erros de Iniciante

### 1. `declared and not used`

```go
// ❌ ERRO
func main() {
    x := 42 // declared and not used
}
```

Go **não permite** variáveis sem uso. Ou use, ou remova.

```go
// ✅ CORRETO
func main() {
    x := 42
    fmt.Println(x)
    
    // Se precisa ignorar temporariamente:
    _ = 42  // blank identifier
}
```

> **Por quê?** Go prioriza código limpo. Variável não usada = provável bug.

---

### 2. Confundir `:=` com `=`

```go
// ❌ ERRO — := cria NOVA variável (shadow!)
x := 10
if true {
    x := 20  // NOVA variável x dentro do if (shadow)
}
fmt.Println(x) // imprime 10, não 20!
```

```go
// ✅ CORRETO — = reatribui a MESMA variável
x := 10
if true {
    x = 20  // reatribui o x externo
}
fmt.Println(x) // imprime 20
```

> **Regra**: `:=` = criar nova variável. `=` = modificar existente. Dentro de blocos internos, use `=` se quer mudar a variável externa.

---

### 3. Goroutine que termina antes de executar

```go
// ❌ ERRO — programa termina antes da goroutine rodar
func main() {
    go fmt.Println("olá")  // pode nunca imprimir!
}
```

```go
// ✅ CORRETO — esperar com WaitGroup
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("olá")
    }()
    wg.Wait() // espera terminar
}
```

> **Por quê?** `main()` termina = programa todo morre, incluindo goroutines.

---

### 4. Nil pointer dereference

```go
// ❌ ERRO — acessar campo de ponteiro nil
var p *Pessoa   // p é nil
fmt.Println(p.Nome)  // PANIC: nil pointer dereference
```

```go
// ✅ CORRETO — verificar antes
var p *Pessoa
if p != nil {
    fmt.Println(p.Nome)
}

// ✅ Ou inicializar
p := &Pessoa{Nome: "Alice"}
fmt.Println(p.Nome)
```

> **Onde isso acontece mais**: retorno de funções que retornam ponteiro + error. Sempre verifique o error ANTES de usar o ponteiro.

---

### 5. Slice append gotcha

```go
// ❌ SURPRESA — append pode ou NÃO modificar o slice original
s := []int{1, 2, 3}
s2 := s[:2]      // s2 = [1, 2], compartilha memória com s
s2 = append(s2, 99)
fmt.Println(s)    // [1, 2, 99] ← s FOI modificado!!
```

```go
// ✅ SEGURO — usar three-index slice para isolar capacidade
s2 := s[:2:2]          // cap(s2) = 2 → append vai alocar novo array
s2 = append(s2, 99)
fmt.Println(s)          // [1, 2, 3] ← s não mudou

// ✅ Ou usar copy
s2 := make([]int, 2)
copy(s2, s[:2])
s2 = append(s2, 99)
```

> **Por quê?** Slices compartilham o array subjacente. `append` só aloca novo array quando a capacidade é ultrapassada.

---

### 6. Loop com goroutine capturando variável errada

```go
// ❌ ERRO CLÁSSICO — todas as goroutines veem o ÚLTIMO valor
for _, v := range []string{"a", "b", "c"} {
    go func() {
        fmt.Println(v)  // imprime "c" "c" "c"
    }()
}
```

```go
// ✅ CORRETO (Go 1.22+) — loop var tem escopo por iteração
// Em Go 1.22+, o código acima funciona corretamente!

// ✅ CORRETO (Go < 1.22) — passar como parâmetro
for _, v := range []string{"a", "b", "c"} {
    go func(val string) {
        fmt.Println(val)  // imprime "a" "b" "c"
    }(v)
}
```

> **Nota**: Go 1.22 corrigiu isso! Mas muitos projetos ainda rodam versões anteriores.

---

### 7. Map não é thread-safe

```go
// ❌ ERRO — race condition! Crash aleatório
m := make(map[string]int)
for i := 0; i < 100; i++ {
    go func(n int) {
        m[fmt.Sprint(n)] = n  // PANIC: concurrent map writes
    }(i)
}
```

```go
// ✅ CORRETO — usar sync.RWMutex (preferido para read-heavy)
var mu sync.RWMutex
m := make(map[string]int)

// Escrita:
mu.Lock()
m["chave"] = 1
mu.Unlock()

// Leitura:
mu.RLock()
v := m["chave"]
mu.RUnlock()
```

---

### 8. Esquecer de tratar error

```go
// ❌ ERRO — ignorar erro
arquivo, _ := os.Open("dados.txt")  // se falhar, arquivo é nil
dados := arquivo.Read(buf)           // PANIC!
```

```go
// ✅ CORRETO — SEMPRE verificar
arquivo, err := os.Open("dados.txt")
if err != nil {
    log.Fatal("erro ao abrir:", err)
}
defer arquivo.Close()
```

> **Regra absoluta em Go**: NUNCA ignore um `error`. Se realmente não importa, documente por quê: `_ = arquivo.Close() // erro de close é inofensivo aqui`

---

### 9. Defer em loop

```go
// ❌ PERIGOSO — defer só executa ao SAIR da função, não do loop
for _, nome := range arquivos {
    f, _ := os.Open(nome)
    defer f.Close()  // TODOS os arquivos ficam abertos até o fim da função!
}
```

```go
// ✅ CORRETO — usar função auxiliar
processarArquivo := func(nome string) {
    f, err := os.Open(nome)
    if err != nil {
        return
    }
    defer f.Close()  // fecha ao sair desta função
    // processar f
}
for _, nome := range arquivos {
    processarArquivo(nome)
}
```

---

### 10. String é imutável

```go
// ❌ ERRO — não pode modificar caractere individual
s := "hello"
s[0] = 'H'  // ERRO DE COMPILAÇÃO

// ✅ CORRETO — converter para []byte ou []rune
s := "hello"
b := []byte(s)
b[0] = 'H'
s = string(b)  // "Hello"
```

---

### 11. Range retorna cópia

```go
// ❌ SURPRESA — range cria CÓPIA do elemento
type Pessoa struct{ Nome string; Idade int }
pessoas := []Pessoa{{Nome: "Alice", Idade: 30}}

for _, p := range pessoas {
    p.Idade = 99  // modifica a CÓPIA, não o original!
}
fmt.Println(pessoas[0].Idade)  // ainda 30
```

```go
// ✅ CORRETO — usar índice
for i := range pessoas {
    pessoas[i].Idade = 99  // modifica o original
}
```

---

### 12. Interface nil vs valor nil

```go
// ❌ SURPRESA — interface com tipo mas valor nil NÃO é nil
var p *Pessoa = nil
var i interface{} = p

fmt.Println(p == nil)  // true
fmt.Println(i == nil)  // FALSE!! 😱
```

> **Por quê?** Uma interface tem dois campos: `(tipo, valor)`. Quando `i = p`, o tipo é `*Pessoa` e o valor é `nil`. A interface só é nil quando **ambos** são nil.

```go
// ❌ O mesmo problema aparece em funções que retornam error
func podeRetornarNil() error {
    var err *MyError = nil
    return err  // retorna interface com tipo *MyError + valor nil
                // quem chamar: if err != nil → TRUE! Bug silencioso!
}

// ✅ CORRETO
func podeRetornarNil() error {
    return nil  // retorna interface nil de verdade
}
```

---

### 13. Imports circulares

```go
// ❌ ERRO — pacote A importa B, e B importa A
// package a → import "b"
// package b → import "a"
// ERRO: import cycle not allowed
```

```go
// ✅ SOLUÇÃO — extrair interface para pacote separado
// package models → define interfaces e tipos comuns
// package a → importa models
// package b → importa models
```

---

## Parte 2 — Erros Intermediários

### 14. sync.WaitGroup: Add depois de Go

```go
// ❌ ERRO — Add chamado dentro da goroutine (race condition!)
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    go func(n int) {
        wg.Add(1)       // pode executar DEPOIS de wg.Wait()!
        defer wg.Done()
        processar(n)
    }(i)
}
wg.Wait() // pode retornar antes de todas goroutines rodarem
```

```go
// ✅ CORRETO — Add ANTES de lançar a goroutine
var wg sync.WaitGroup
for i := 0; i < 10; i++ {
    wg.Add(1)           // Add no goroutine pai, ANTES do go
    go func(n int) {
        defer wg.Done()
        processar(n)
    }(i)
}
wg.Wait()
```

> **Regra**: `wg.Add(n)` sempre antes de `go func()`. O `Done()` vai dentro da goroutine.

---

### 15. Fechar channel fechado — panic!

```go
// ❌ PÂNICO — fechar um channel já fechado
ch := make(chan int)
close(ch)
close(ch) // PANIC: close of closed channel

// ❌ PÂNICO — enviar para channel fechado
close(ch)
ch <- 1  // PANIC: send on closed channel
```

```go
// ✅ PADRÃO — quem produz fecha, quem consome apenas lê
func produtor(ch chan<- int) {
    defer close(ch)  // fecha uma vez, ao terminar de produzir
    for i := 0; i < 10; i++ {
        ch <- i
    }
}

func consumidor(ch <-chan int) {
    for v := range ch {  // range termina automaticamente quando ch fecha
        fmt.Println(v)
    }
}
```

> **Regra**: apenas o **produtor** fecha o channel. O consumidor nunca fecha.

---

### 16. Context: values como chave `string` — colisão silenciosa

```go
// ❌ RUIM — chave string pode colidir entre pacotes
ctx = context.WithValue(ctx, "userID", 42)

// ❌ RUIM — qualquer pacote pode sobrescrever a mesma chave "userID"
ctx = context.WithValue(ctx, "userID", "outro valor")
```

```go
// ✅ CORRETO — usar tipo privado como chave (impossível colidir)
type contextKey string
const userIDKey contextKey = "userID"

ctx = context.WithValue(ctx, userIDKey, 42)

// Quem não importar o tipo nunca vai conseguir ler ou sobrescrever
```

> **Regra**: nunca use `string` ou `int` como chave de context. Use um tipo definido no seu pacote.

---

### 17. Goroutine leak — goroutine que nunca termina

```go
// ❌ GOROUTINE LEAK — goroutine bloqueada para sempre
func buscarDados(url string) <-chan Result {
    ch := make(chan Result)
    go func() {
        result := http.Get(url)
        ch <- result  // se ninguém ler, a goroutine fica presa aqui para sempre
    }()
    return ch
}

// Se o chamador não ler o channel, a goroutine vaza.
```

```go
// ✅ CORRETO — usar context para cancelamento
func buscarDados(ctx context.Context, url string) <-chan Result {
    ch := make(chan Result, 1)  // buffer de 1 evita bloqueio
    go func() {
        req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
        result, err := http.DefaultClient.Do(req)
        select {
        case ch <- Result{resp: result, err: err}:
        case <-ctx.Done():
            // contexto cancelado, goroutine termina limpa
        }
    }()
    return ch
}
```

> **Como detectar**: use `go test -race` e `runtime.NumGoroutine()` em testes. Ou habilite pprof e observe `/debug/pprof/goroutine`.

---

### 18. errors.Is não funciona sem Unwrap

```go
// ❌ NÃO FUNCIONA — wrapping manual sem Unwrap
var ErrNotFound = errors.New("not found")

type DBError struct {
    ID  int
    Err error
}

func (e *DBError) Error() string {
    return fmt.Sprintf("db error for id %d: %v", e.ID, e.Err)
}

err := &DBError{ID: 42, Err: ErrNotFound}
fmt.Println(errors.Is(err, ErrNotFound))  // FALSE — sem Unwrap!
```

```go
// ✅ CORRETO — implementar Unwrap
func (e *DBError) Unwrap() error { return e.Err }

fmt.Println(errors.Is(err, ErrNotFound))  // TRUE ✅

// ✅ Ou usar fmt.Errorf com %w (faz Unwrap automaticamente)
err := fmt.Errorf("db error for id %d: %w", 42, ErrNotFound)
fmt.Println(errors.Is(err, ErrNotFound))  // TRUE ✅
```

---

### 19. json.Unmarshal em ponteiro nil

```go
// ❌ ERRO — desserializar em ponteiro nil
var p *Pessoa
json.Unmarshal(data, p)  // silenciosamente ignora! p continua nil
fmt.Println(p.Nome)       // PANIC
```

```go
// ✅ CORRETO — ponteiro para struct alocada
var p Pessoa
json.Unmarshal(data, &p)  // &p é o endereço de uma struct real

// ✅ Ou com ponteiro de ponteiro
p := new(Pessoa)
json.Unmarshal(data, p)
```

---

### 20. Mutex copiado por valor — deadlock silencioso

```go
// ❌ ERRO — sync.Mutex nunca deve ser copiado
type Cache struct {
    mu sync.Mutex
    m  map[string]string
}

func processar(c Cache) {  // Cache passado por VALOR = mutex copiado!
    c.mu.Lock()
    defer c.mu.Unlock()
    // ...
}
```

```go
// ✅ CORRETO — sempre usar ponteiro para structs com mutex
func processar(c *Cache) {  // ponteiro: mutex não é copiado
    c.mu.Lock()
    defer c.mu.Unlock()
    // ...
}
```

> **go vet** detecta isso! Sempre rode `go vet ./...` antes de fazer PR.

---

## Parte 3 — Erros Avançados

### 21. Interface nil retornada de função — o bug mais traiçoeiro do Go

Este é uma variação mais sutil do erro #12, que aparece em código de produção real:

```go
// ❌ BUG CLÁSSICO EM PRODUÇÃO
type MyError struct{ msg string }
func (e *MyError) Error() string { return e.msg }

func validar(x int) error {
    var err *MyError   // tipo concreto nil
    if x < 0 {
        err = &MyError{msg: "negativo"}
    }
    return err  // ← retorna interface (type=*MyError, value=nil) quando x >= 0
}

if err := validar(5); err != nil {  // TRUE mesmo com x=5 !!
    log.Fatal(err)  // programa morre sem motivo
}
```

```go
// ✅ REGRA: retorne nil explícito, não um ponteiro tipado nil
func validar(x int) error {
    if x < 0 {
        return &MyError{msg: "negativo"}
    }
    return nil  // nil interface = (type=nil, value=nil)
}
```

---

### 22. Cache stampede — goroutines duplicando trabalho caro

```go
// ❌ PROBLEMA — múltiplas goroutines fazem a mesma busca cara simultaneamente
var (
    cacheMu sync.RWMutex
    cache   = make(map[string]string)
)

func buscar(key string) string {
    cacheMu.RLock()
    if v, ok := cache[key]; ok {
        cacheMu.RUnlock()
        return v
    }
    cacheMu.RUnlock()

    // 1000 goroutines simultâneas todas chegam aqui e fazem a mesma query!
    result := queryDB(key) // query cara

    cacheMu.Lock()
    cache[key] = result
    cacheMu.Unlock()
    return result
}
```

```go
// ✅ SOLUÇÃO — singleflight garante que apenas UMA goroutine executa por chave
import "golang.org/x/sync/singleflight"

var (
    cacheMu sync.RWMutex
    cache   = make(map[string]string)
    sf      singleflight.Group
)

func buscar(key string) string {
    cacheMu.RLock()
    if v, ok := cache[key]; ok {
        cacheMu.RUnlock()
        return v
    }
    cacheMu.RUnlock()

    // singleflight: apenas uma goroutine executa, as outras esperam o resultado
    result, _, _ := sf.Do(key, func() (any, error) {
        v := queryDB(key)
        cacheMu.Lock()
        cache[key] = v
        cacheMu.Unlock()
        return v, nil
    })
    return result.(string)
}
```

---

### 23. Slice de ponteiros vs slice de valores — quem o GC vai coletar?

```go
// ❌ SURPRESA — slice de ponteiros mantém todos os objetos vivos
var cache []*Pessoa
cache = append(cache, &Pessoa{Nome: "Alice"})
cache = append(cache, &Pessoa{Nome: "Bob"})

// Remover Alice:
cache = cache[1:]  // cache aponta para Bob, mas o array subjacente
                   // ainda tem referência para Alice → GC não coleta!
```

```go
// ✅ CORRETO — nil explícito libera a referência
cache[0] = nil      // libera Alice para o GC
cache = cache[1:]   // agora Alice pode ser coletada
```

---

### 24. select com default — busy loop acidental

```go
// ❌ PROBLEMA — consome 100% de CPU sem pausa
func consumir(ch <-chan int) {
    for {
        select {
        case v := <-ch:
            processar(v)
        default:
            // ← executa constantemente quando ch está vazio
            // CPU vai a 100%!
        }
    }
}
```

```go
// ✅ CORRETO — sem default: bloqueia até ter dado
func consumir(ch <-chan int) {
    for v := range ch {  // bloqueia elegantemente
        processar(v)
    }
}

// ✅ CORRETO — com timeout se precisar de não-bloqueio
func consumirComTimeout(ch <-chan int) {
    for {
        select {
        case v := <-ch:
            processar(v)
        case <-time.After(5 * time.Second):
            log.Println("timeout esperando dado")
            return
        }
    }
}
```

> **Regra**: `select { default: }` é válido para **polling não-bloqueante pontual**. Nunca em um loop sem sleep.

---

### 25. time.After em loop — goroutine e timer leak

```go
// ❌ LEAK — time.After cria novo timer e goroutine a cada iteração
for {
    select {
    case v := <-ch:
        processar(v)
    case <-time.After(time.Second):  // novo timer a cada volta do loop!
        // timers anteriores não cancelados ficam vivos por 1 segundo
        log.Println("timeout")
    }
}
```

```go
// ✅ CORRETO — reutilizar timer com Reset
timer := time.NewTimer(time.Second)
defer timer.Stop()

for {
    timer.Reset(time.Second)
    select {
    case v := <-ch:
        processar(v)
    case <-timer.C:
        log.Println("timeout")
    }
}
```

---

### 26. Benchmark sem b.ReportAllocs — alocações invisíveis

```go
// ❌ INCOMPLETO — não mede alocações
func BenchmarkMinhaFunc(b *testing.B) {
    for i := 0; i < b.N; i++ {
        minhaFunc()
    }
}
```

```bash
# ✅ CORRETO — sempre use -benchmem
go test -bench=. -benchmem -count=5

# Resultado mostra:
# BenchmarkMinhaFunc-8    1000000    1234 ns/op    256 B/op    3 allocs/op
#                                                   ↑               ↑
#                                          bytes por op    alocações por op
```

```go
// ✅ Ou forçar no código do benchmark
func BenchmarkMinhaFunc(b *testing.B) {
    b.ReportAllocs()
    for i := 0; i < b.N; i++ {
        minhaFunc()
    }
}
```

---

### 27. init() — ordem não garantida entre pacotes

```go
// ❌ ARMADILHA — init() de pacotes diferentes têm ordem não determinística
// package database/init() pode rodar antes ou depois de package config/init()
// Se database depende de config estar pronto → bug intermitente

// ❌ init() que falha silenciosamente
func init() {
    db, err := sql.Open("postgres", os.Getenv("DB_URL"))
    if err != nil {
        log.Println(err)  // loga mas continua! db é nil
    }
    globalDB = db
}
```

```go
// ✅ MELHOR — inicialização explícita em main() ou via construtor
func main() {
    db, err := database.Connect(os.Getenv("DB_URL"))
    if err != nil {
        log.Fatalf("falha ao conectar ao banco: %v", err)
    }
    defer db.Close()
    
    app := NewApp(db)
    app.Run()
}
```

> **Regra**: use `init()` apenas para registro (como `sql.Register`, `image.RegisterFormat`). Nunca para inicialização com estado mutável.

---

### 28. Conversão de tipo em interface — panic com tipo errado

```go
// ❌ PANIC — type assertion sem verificação
var i interface{} = "hello"
n := i.(int)  // PANIC: interface conversion: interface {} is string, not int
```

```go
// ✅ CORRETO — two-value form (comma ok)
n, ok := i.(int)
if !ok {
    // trata o caso graciosamente
    log.Printf("esperava int, recebi %T", i)
    return
}

// ✅ Ou type switch para múltiplos tipos
switch v := i.(type) {
case int:
    fmt.Printf("int: %d\n", v)
case string:
    fmt.Printf("string: %s\n", v)
default:
    fmt.Printf("tipo desconhecido: %T\n", v)
}
```

---

## 🧠 Resumo Rápido

| # | Erro | Solução em 1 frase |
|---|------|-------------------|
| 1 | declared and not used | Use ou delete a variável |
| 2 | := vs = | := cria nova, = reatribui existente |
| 3 | Goroutine não executa | Use WaitGroup ou channel |
| 4 | Nil pointer | Verifique error antes de usar ponteiro |
| 5 | Slice append | Use `s[:n:n]` ou copy() para independência |
| 6 | Loop + goroutine | Passe variável como parâmetro (ou use Go 1.22+) |
| 7 | Map + goroutines | Use sync.RWMutex |
| 8 | Error ignorado | SEMPRE trate com if err != nil |
| 9 | Defer em loop | Encapsule em função auxiliar |
| 10 | String imutável | Converta para []byte |
| 11 | Range copia | Use índice para modificar |
| 12 | Interface nil | Cuidado com (tipo, nil) vs nil |
| 13 | Import circular | Extraia interfaces para pacote intermediário |
| 14 | WaitGroup.Add tarde | Add() antes do go, Done() dentro |
| 15 | Close de channel fechado | Só o produtor fecha; use defer close |
| 16 | Context com chave string | Use tipo privado como chave |
| 17 | Goroutine leak | Use context + select + ctx.Done() |
| 18 | errors.Is sem Unwrap | Implemente Unwrap() ou use %w |
| 19 | json.Unmarshal em nil | Passe endereço de struct alocada |
| 20 | Mutex copiado | Sempre use ponteiro para structs com mutex |
| 21 | Interface nil retornada | Retorne nil explícito, não ponteiro tipado nil |
| 22 | Cache stampede | Use singleflight.Group |
| 23 | Slice de ponteiros | nil explícito antes de remover elemento |
| 24 | select default em loop | Sem default: bloqueia. Com default: use sleep |
| 25 | time.After em loop | Reutilize timer com NewTimer + Reset |
| 26 | Benchmark sem -benchmem | Sempre `-benchmem` ou `b.ReportAllocs()` |
| 27 | init() com estado | Inicialize explicitamente em main() |
| 28 | Type assertion sem ok | Sempre use `v, ok := i.(T)` |
