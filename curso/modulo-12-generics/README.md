# Módulo 12 — Generics

[← Concorrência](../modulo-11-concorrencia/README.md) | [Próximo: Testes →](../modulo-13-testes/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que usaríamos generics em vez de `interface{}`/`any`?
> 2. O que é uma "constraint" em generics Go?
> 3. Quando generics NÃO são a solução certa?

---

## 1. O Problema que Generics Resolve

Antes de Go 1.18, você precisava duplicar código para cada tipo:

```go
// ❌ Antes: um para cada tipo
func MinInt(a, b int) int {
    if a < b { return a }
    return b
}
func MinFloat64(a, b float64) float64 {
    if a < b { return a }
    return b
}

// ❌ Ou usar interface{} — perde type safety, precisa de type assertion
func Min(a, b interface{}) interface{} {
    // como comparar? não dá com <
}

// ✅ Go 1.18+: generics — uma função, múltiplos tipos, type-safe
func Min[T constraints.Ordered](a, b T) T {
    if a < b { return a }
    return b
}

Min(3, 5)         // int
Min(3.14, 2.71)   // float64
Min("abc", "xyz") // string
```

---

## 2. Sintaxe — Type Parameters

```go
// [T any] — T é o type parameter, any é a constraint (qualquer tipo)
func Contem[T comparable](slice []T, item T) bool {
    for _, v := range slice {
        if v == item {
            return true
        }
    }
    return false
}

// Múltiplos type parameters
func Map[T any, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Uso — Go infere os tipos (você raramente precisa especificar)
nums := []int{1, 2, 3, 4, 5}
strs := Map(nums, func(n int) string {
    return strconv.Itoa(n)
}) // []string{"1", "2", "3", "4", "5"}

// Especificando tipos explicitamente (raramente necessário)
strs2 := Map[int, string](nums, strconv.Itoa)
```

---

## 3. Constraints — O Que o Tipo Pode Fazer

Constraints definem quais operações são permitidas sobre o type parameter:

```go
import "golang.org/x/exp/constraints"

// Constraints built-in
any          // qualquer tipo (= interface{})
comparable   // tipos que suportam == e != (int, string, struct sem slice/map/func)

// Constraints do pacote constraints
constraints.Ordered    // int, float, string (suportam <, >, <=, >=)
constraints.Integer    // int, int8, ..., uint, uint8...
constraints.Float      // float32, float64
constraints.Signed     // inteiros com sinal
constraints.Unsigned   // inteiros sem sinal

// Constraint customizada — union de tipos
type Numero interface {
    int | int32 | int64 | float32 | float64
}

func Soma[T Numero](valores ...T) T {
    var total T
    for _, v := range valores {
        total += v
    }
    return total
}

Soma(1, 2, 3)         // int → 6
Soma(1.1, 2.2, 3.3)   // float64 → 6.6

// Constraint com método
type Stringer interface {
    String() string
}

func Imprimir[T Stringer](items []T) {
    for _, item := range items {
        fmt.Println(item.String())
    }
}
```

---

## 4. Tipos Genéricos — Structs e Interfaces

```go
// Stack genérica — funciona com qualquer tipo
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(item T) {
    s.items = append(s.items, item)
}

func (s *Stack[T]) Pop() (T, bool) {
    var zero T
    if len(s.items) == 0 {
        return zero, false
    }
    item := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return item, true
}

func (s *Stack[T]) Len() int { return len(s.items) }

// Uso
intStack := Stack[int]{}
intStack.Push(1)
intStack.Push(2)
v, ok := intStack.Pop() // v=2, ok=true

strStack := Stack[string]{}
strStack.Push("hello")
strStack.Push("world")

// Resultado genérico com erro
type Result[T any] struct {
    Value T
    Err   error
}

func BuscarUsuario(id int) Result[*User] {
    // ...
    return Result[*User]{Value: user, Err: nil}
}
```

---

## 5. Funções Utilitárias — A Stdlib go 1.21+

Go 1.21 adicionou pacotes genéricos na stdlib:

```go
import (
    "slices"
    "maps"
)

// slices — operações em slices de qualquer tipo comparável/ordenável
nums := []int{3, 1, 4, 1, 5, 9, 2, 6}
slices.Sort(nums)                  // ordena
max := slices.Max(nums)            // 9
min := slices.Min(nums)            // 1
idx := slices.Index(nums, 5)       // índice de 5
found := slices.Contains(nums, 7)  // false
equal := slices.Equal([]int{1,2}, []int{1,2}) // true

// maps — operações em maps
m := map[string]int{"a": 1, "b": 2, "c": 3}
keys := maps.Keys(m)      // []string{"a","b","c"} (ordem não garantida)
vals := maps.Values(m)    // []int{1,2,3}
clone := maps.Clone(m)    // cópia independente
```

---

## 6. Quando NÃO Usar Generics

Generics são poderosos mas têm custo — código mais complexo, tempo de compilação maior, às vezes mais difícil de depurar.

```
Use generics quando:
✅ A mesma lógica se repete identicamente para N tipos diferentes
✅ Você está construindo estruturas de dados reutilizáveis (Stack, Queue, Set)
✅ Funções utilitárias (Map, Filter, Reduce) que operam sobre slices/maps

NÃO use generics quando:
❌ Interfaces comuns resolvem o problema (é mais idiomático em Go)
❌ Você tem apenas 1-2 tipos concretos (duplicar é mais simples)
❌ O comportamento difere entre tipos (use interfaces com métodos)
❌ Para simplificar código que já está simples
```

```go
// ❌ Desnecessário — interface resolve melhor
func ImprimirGenerico[T fmt.Stringer](item T) {
    fmt.Println(item.String())
}
// ✅ Mais simples com interface diretamente
func Imprimir(item fmt.Stringer) {
    fmt.Println(item.String())
}

// ✅ Generics fazem sentido aqui (sem duplicação, type-safe)
func Keys[K comparable, V any](m map[K]V) []K {
    result := make([]K, 0, len(m))
    for k := range m {
        result = append(result, k)
    }
    return result
}
```

---

## ✅ Checklist de Generics

- [ ] Type parameters têm nomes descritivos quando não óbvio (`T` para geral, `K`/`V` para chave/valor)
- [ ] Constraints usam o tipo mais restrito que ainda funciona (não `any` quando `comparable` é suficiente)
- [ ] Verifico se uma interface comum resolve o problema antes de usar generics
- [ ] Uso pacotes `slices` e `maps` da stdlib (Go 1.21+) antes de reimplementar

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo12_generics.go` | Funções genéricas, tipos genéricos, constraints |
| `exercicios/ex12_generics.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Funções Utilitárias
Implemente versões genéricas de: `Contains[T comparable](s []T, v T) bool`, `Filter[T any](s []T, fn func(T) bool) []T`, `Map[T, U any](s []T, fn func(T) U) []U`. Escreva testes com `int`, `string` e `struct`.

### 🟡 2. Set Genérico Thread-Safe
Implemente `Set[T comparable]` com métodos `Add`, `Remove`, `Contains`, `Len`, `Union`, `Intersection`. Use `sync.RWMutex` para thread-safety. Demonstre que funciona com `int` e `string`.

### 🟡 3. Result[T] — Error Handling Funcional
Crie `Result[T any]` com `Value T` e `Err error`. Implemente métodos `Map[U](fn func(T) U) Result[U]`, `FlatMap[U](fn func(T) Result[U]) Result[U]`, `Unwrap() T` (panic se err), `UnwrapOr(def T) T`. Encadeie transformações de forma elegante.

### 🔴 4. Cache LRU Genérico
Implemente um `LRUCache[K comparable, V any]` com tamanho fixo usando uma lista duplamente encadeada e um map. Métodos: `Get(key K) (V, bool)`, `Put(key K, value V)`, `Len() int`. Escreva benchmarks comparando com `sync.Map`.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Concorrência](../modulo-11-concorrencia/README.md) | [Próximo: Testes →](../modulo-13-testes/README.md)
