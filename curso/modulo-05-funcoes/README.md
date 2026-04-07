# Módulo 05 — Funções

[← Coleções](../modulo-04-colecoes/README.md) | [Próximo: Ponteiros →](../modulo-06-ponteiros/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que Go retorna múltiplos valores ao invés de usar exceções?
> 2. O que é uma closure e qual variável ela "captura"?
> 3. Qual a ordem de execução de múltiplos `defer`?

---

## 1. Funções — Cidadãos de Primeira Classe

Em Go, funções são valores — podem ser atribuídas a variáveis, passadas como parâmetros e retornadas de outras funções.

```go
// Função básica
func soma(a, b int) int {
    return a + b
}

// Parâmetros do mesmo tipo: lista compacta
func somaCompacta(a, b, c int) int {
    return a + b + c
}

// Função como variável
var fn func(int, int) int = soma
resultado := fn(3, 4)  // 7

// Função anônima
dobro := func(x int) int { return x * 2 }
fmt.Println(dobro(5)) // 10
```

---

## 2. Múltiplos Retornos — O Pattern Fundamental do Go

Go tem múltiplos retornos — o mecanismo central para tratamento de erros sem exceções:

```go
// Convenção: erro sempre é o ÚLTIMO retorno
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("divisão por zero")
    }
    return a / b, nil
}

resultado, err := dividir(10, 2)
if err != nil {
    log.Fatal(err)
}
fmt.Println(resultado) // 5

// Descartar um retorno com _
valor, _ := dividir(10, 2) // ignora o erro — só faça isso se tem certeza
```

### Retornos nomeados — use com moderação

```go
// Retornos nomeados: úteis em funções curtas com lógica clara
func minMax(nums []int) (min, max int) {
    min, max = nums[0], nums[0]
    for _, n := range nums[1:] {
        if n < min { min = n }
        if n > max { max = n }
    }
    return // naked return — retorna min e max
}

// ❌ Evite em funções longas — naked return prejudica legibilidade
// Em funções longas, prefira return explícito:
func operacaoCompleta() (resultado int, err error) {
    // ... muito código ...
    return resultado, err // explícito: mais claro
}
```

---

## 3. Funções Variádicas

```go
// ...T aceita zero ou mais argumentos do tipo T
func somaTotal(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

somaTotal()          // 0
somaTotal(1, 2, 3)   // 6
somaTotal(1, 2, 3, 4, 5) // 15

// Expandir slice para variádica com ...
nums := []int{1, 2, 3, 4}
somaTotal(nums...)   // equivale a somaTotal(1, 2, 3, 4)

// Padrão comum: fmt.Println, append, etc.
// fmt.Println(a ...any)
// append(slice []T, elems ...T) []T
```

---

## 4. Closures — Capturando o Ambiente

Uma closure é uma função que "fecha sobre" variáveis do escopo onde foi criada. Ela mantém referência à variável, não uma cópia.

```go
// Contador com estado encapsulado
func novoContador() func() int {
    i := 0  // esta variável é "capturada" pela closure
    return func() int {
        i++
        return i
    }
}

c1 := novoContador()
c2 := novoContador()

fmt.Println(c1()) // 1
fmt.Println(c1()) // 2
fmt.Println(c2()) // 1 — estado independente de c1!
fmt.Println(c1()) // 3

// ⚠️ Closure captura a VARIÁVEL, não o valor!
funcs := make([]func(), 3)
for i := 0; i < 3; i++ {
    i := i  // ← cria nova variável a cada iteração (fix para Go < 1.22)
    funcs[i] = func() { fmt.Println(i) }
}
// Go 1.22+: o loop cria variável nova automaticamente
```

### Casos de uso práticos para closures

```go
// 1. Middleware HTTP
func withLogging(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
    })
}

// 2. Função que "lembra" configuração
func novoMultiplicador(fator int) func(int) int {
    return func(x int) int {
        return x * fator  // fator capturado da closure
    }
}
dobro := novoMultiplicador(2)
triplo := novoMultiplicador(3)
fmt.Println(dobro(5))  // 10
fmt.Println(triplo(5)) // 15

// 3. Memoização (cache de resultados)
func memoize(fn func(int) int) func(int) int {
    cache := make(map[int]int)
    return func(n int) int {
        if v, ok := cache[n]; ok {
            return v
        }
        result := fn(n)
        cache[n] = result
        return result
    }
}
```

---

## 5. defer — Limpeza Garantida

`defer` adia a execução para quando a função retornar — independente de como ela retorna (normal, `return`, `panic`).

```go
func lerArquivo(nome string) error {
    f, err := os.Open(nome)
    if err != nil {
        return err
    }
    defer f.Close()  // ← executado quando lerArquivo retornar, SEMPRE
    
    // processa o arquivo...
    return nil
}
```

### Ordem LIFO: último defer, primeiro a executar

```go
func exemplo() {
    defer fmt.Println("1º defer — executa por ÚLTIMO")
    defer fmt.Println("2º defer — executa no meio")
    defer fmt.Println("3º defer — executa PRIMEIRO")
    fmt.Println("corpo da função")
}
// Saída:
// corpo da função
// 3º defer — executa PRIMEIRO
// 2º defer — executa no meio
// 1º defer — executa por ÚLTIMO
```

### defer captura argumentos na hora da declaração

```go
func exemplo() {
    x := 10
    defer fmt.Println("defer capturou x =", x) // captura x=10 aqui
    x = 20
    fmt.Println("no corpo, x =", x)
}
// Saída:
// no corpo, x = 20
// defer capturou x = 10  ← argumento capturado no momento do defer
```

### defer com função anônima captura por referência

```go
func exemplo() {
    x := 10
    defer func() {
        fmt.Println("defer vê x =", x) // referência — vê o valor atual
    }()
    x = 20
}
// Saída: defer vê x = 20
```

### Uso com recover

```go
func safe(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("panic: %v", r)
        }
    }()
    fn()
    return nil
}
```

---

## 6. Funções de Ordem Superior — Programação Funcional em Go

```go
// Map: transforma cada elemento
func Map[T, U any](s []T, fn func(T) U) []U {
    result := make([]U, len(s))
    for i, v := range s {
        result[i] = fn(v)
    }
    return result
}

// Filter: filtra elementos
func Filter[T any](s []T, fn func(T) bool) []T {
    var result []T
    for _, v := range s {
        if fn(v) {
            result = append(result, v)
        }
    }
    return result
}

// Reduce: agrega todos os elementos
func Reduce[T, U any](s []T, init U, fn func(U, T) U) U {
    acc := init
    for _, v := range s {
        acc = fn(acc, v)
    }
    return acc
}

// Uso
nums := []int{1, 2, 3, 4, 5}
dobrados := Map(nums, func(n int) int { return n * 2 })     // [2 4 6 8 10]
pares := Filter(nums, func(n int) bool { return n%2 == 0 }) // [2 4]
soma := Reduce(nums, 0, func(acc, n int) int { return acc + n }) // 15
```

---

## ✅ Checklist de Funções

- [ ] Erro sempre é o **último** valor de retorno
- [ ] Funções usam `defer` para liberar recursos (arquivo, lock, conexão)
- [ ] Closures capturam referências — cuidado com variáveis de loop
- [ ] Funções variádicas usam `...T` e podem receber slices com `slice...`
- [ ] `panic/recover` são usados apenas em casos extremos (nunca para controle de fluxo normal)

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo05_funcoes.go` | Todos os tipos de funções, defer, closures |
| `exercicios/ex05_funcoes.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Calculadora Funcional
Crie um `map[string]func(float64, float64) float64` com as 4 operações básicas. Implemente uma função `Calcular(op string, a, b float64) (float64, error)` que usa esse mapa.

### 🟡 2. Pipeline Funcional
Implemente as funções `Map`, `Filter` e `Reduce` usando generics. Use-as para: transformar strings em maiúsculo, filtrar as que têm mais de 3 chars, e concatenar todas em uma única string.

### 🟡 3. Retry com Backoff
Escreva uma função `WithRetry(fn func() error, tentativas int) error` usando closures. Se `fn` retornar erro, tente novamente com delay crescente (1s, 2s, 4s...). Use `defer` para logar o resultado final.

### 🔴 4. Memoização Genérica
Implemente `Memoize[K comparable, V any](fn func(K) V) func(K) V` usando uma closure com map interno e `sync.RWMutex` para thread-safety. Demonstre com Fibonacci memoizado.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Coleções](../modulo-04-colecoes/README.md) | [Próximo: Ponteiros →](../modulo-06-ponteiros/README.md)
