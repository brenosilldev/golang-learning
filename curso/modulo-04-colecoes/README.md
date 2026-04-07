# Módulo 04 — Coleções

[← Controle de Fluxo](../modulo-03-controle-fluxo/README.md) | [Próximo: Funções →](../modulo-05-funcoes/README.md)

---

> **Antes de ler — tente responder:**
> 1. O que é a diferença entre `len` e `cap` em um slice?
> 2. O que acontece quando você modifica um sub-slice? O original muda?
> 3. Por que `map[string]struct{}` é preferível a `map[string]bool` para sets?

---

## 1. Arrays — Tamanho Fixo (raramente usado diretamente)

```go
var arr [5]int                    // [0 0 0 0 0] — zero values
arr2 := [3]string{"a", "b", "c"} // literal
arr3 := [...]int{1, 2, 3, 4}     // ... = Go conta o tamanho automaticamente (4)

// Arrays são valores, não referências!
a := [3]int{1, 2, 3}
b := a        // CÓPIA COMPLETA — b é independente de a
b[0] = 99
fmt.Println(a) // [1 2 3] — não mudou
fmt.Println(b) // [99 2 3]
```

> **Na prática**: arrays são usados internamente (como backing storage de slices), mas você raramente vai declarar um array diretamente. Use **slices**.

---

## 2. Slices — A Estrutura Mais Importante do Go

Um slice é uma **janela** sobre um array. Ele tem 3 componentes internos:

```
slice
├── pointer → aponta para o início no array
├── len     → quantos elementos o slice "vê"
└── cap     → quantos elementos o array tem a partir do ponteiro
```

```go
// Criação
s1 := []int{1, 2, 3, 4, 5}     // literal (len=5, cap=5)
s2 := make([]int, 5)            // make(tipo, len) — todos zeros
s3 := make([]int, 3, 10)        // make(tipo, len, cap) — len=3, cap=10
var s4 []int                    // nil slice — len=0, cap=0
```

### append — Como funciona internamente

```go
s := make([]int, 0, 3) // len=0, cap=3

s = append(s, 1)       // len=1, cap=3
s = append(s, 2)       // len=2, cap=3
s = append(s, 3)       // len=3, cap=3
s = append(s, 4)       // len=4, cap=6 ← Go alocou novo array (2x)

// SEMPRE reatribua o resultado de append!
// append pode retornar um slice DIFERENTE se realocar
s = append(s, 5) // ✅ correto
append(s, 5)     // ❌ errado — descarta o resultado
```

### ⚠️ Armadilha Crítica: Sub-slices compartilham memória

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:3]  // sub = [2, 3], compartilha memória com original

sub[0] = 99
fmt.Println(original) // [1 99 3 4 5] ← original foi modificado!
fmt.Println(sub)      // [99 3]

// ✅ Para uma cópia independente, use copy ou append
copia := append([]int{}, original[1:3]...)  // cópia independente
copia2 := make([]int, len(original[1:3]))
copy(copia2, original[1:3])                  // também independente
```

### ⚠️ Armadilha: append em sub-slice pode corromper o original

```go
a := []int{1, 2, 3, 4, 5}
b := a[:3] // b=[1,2,3], len=3, cap=5 (ainda tem espaço no array original!)

b = append(b, 99)      // len<cap, então SOBRESCREVE a[3]!
fmt.Println(a)         // [1 2 3 99 5] ← a foi corrompido silenciosamente!

// ✅ Use o three-index slice para limitar o cap
b = a[:3:3] // len=3, cap=3 — próximo append aloca novo array
b = append(b, 99)
fmt.Println(a)         // [1 2 3 4 5] ← original protegido
```

### Operações essenciais com slices

```go
// Copiar
src := []int{1, 2, 3}
dst := make([]int, len(src))
n := copy(dst, src)          // retorna quantos elementos foram copiados

// Deletar elemento por índice (mantém ordem)
s := []int{1, 2, 3, 4, 5}
i := 2
s = append(s[:i], s[i+1:]...) // remove elemento no índice 2 → [1, 2, 4, 5]

// Deletar sem manter ordem (mais eficiente)
s[i] = s[len(s)-1]
s = s[:len(s)-1]

// Inserir no meio
pos := 2
val := 99
s = append(s[:pos], append([]int{val}, s[pos:]...)...)

// Reverter
for i, j := 0, len(s)-1; i < j; i, j = i+1, j-1 {
    s[i], s[j] = s[j], s[i]
}
```

---

## 3. Maps — Chave/Valor com Hash Table

```go
// Criação — SEMPRE inicialize com make antes de escrever
m := make(map[string]int)
m["go"] = 1
m["rust"] = 2

// Literal (inicializado)
m2 := map[string]int{
    "go":   1,
    "rust": 2,
}

// Leitura SEGURA — sempre use a forma de dois valores
valor, ok := m["go"]     // ok=true se existe
if !ok {
    // chave não existe
}

// Leitura simples — retorna zero value se não existe
n := m["inexistente"]    // n = 0 (zero value de int)

// Deletar
delete(m, "rust")

// Verificar existência
if _, ok := m["go"]; ok {
    fmt.Println("existe")
}

// Iterar (ordem NÃO garantida)
for k, v := range m {
    fmt.Printf("%s: %d\n", k, v)
}
```

### ⚠️ Maps não são thread-safe

```go
// ❌ RACE CONDITION — duas goroutines lendo/escrevendo ao mesmo tempo
go func() { m["a"] = 1 }()
go func() { m["b"] = 2 }()

// ✅ Use sync.RWMutex ou sync.Map para acesso concorrente
var mu sync.RWMutex
mu.Lock()
m["a"] = 1
mu.Unlock()
```

---

## 4. Sets — Go não tem set nativo

```go
// Use map[T]struct{} — struct{} ocupa 0 bytes!
tipo Set[T comparable] = map[T]struct{}

set := make(map[string]struct{})

// Adicionar
set["golang"] = struct{}{}

// Verificar
_, existe := set["golang"]

// Remover
delete(set, "golang")

// Por que struct{} e não bool?
// map[string]bool{"a": true}   → bool ocupa 1 byte por entrada
// map[string]struct{}{"a": {}} → struct{} ocupa 0 bytes
// Para sets com milhões de entradas, a diferença é significativa
```

---

## 5. Operações Funcionais em Slices (Go 1.21+)

O pacote `slices` da stdlib oferece operações comuns:

```go
import "slices"

s := []int{3, 1, 4, 1, 5, 9, 2, 6}

slices.Sort(s)                   // [1 1 2 3 4 5 6 9]
idx := slices.Index(s, 5)       // índice de 5 = 4
encontrado := slices.Contains(s, 7) // false
s = slices.Compact(s)           // remove duplicatas adjacentes → [1 2 3 4 5 6 9]
s = slices.Reverse(s)           // inverte no local

// Comparar dois slices
a := []int{1, 2, 3}
b := []int{1, 2, 3}
fmt.Println(slices.Equal(a, b)) // true
```

---

## ✅ Checklist de Coleções

- [ ] Nunca escrevo em map `nil` (inicialize com `make`)
- [ ] Sempre reatribuo o resultado de `append` (`s = append(s, ...)`)
- [ ] Sub-slices que precisam ser independentes usam `copy` ou `append([]T{}, ...)`
- [ ] Para acesso concorrente a maps, uso `sync.RWMutex` ou `sync.Map`
- [ ] Para sets, uso `map[T]struct{}` (não `map[T]bool`)
- [ ] Verifico existência de chave em maps com a forma de dois valores `v, ok := m[k]`

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo04_colecoes.go` | Arrays, slices, maps, operações, armadilhas |
| `exercicios/ex04_colecoes.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Frequência de Palavras
Dado um texto, use um `map[string]int` para contar a frequência de cada palavra. Imprima as 5 palavras mais frequentes. Normalize para minúsculas e remova pontuação.

### 🟢 2. Deduplicar Slice
Escreva uma função `Unique(s []string) []string` que remove duplicatas preservando a ordem de primeira aparição. Use um `map[string]struct{}` como conjunto visitado.

### 🟡 3. Stack e Queue com Slices
Implemente uma `Stack[T]` (LIFO) e uma `Queue[T]` (FIFO) usando slices. Cada uma deve ter: `Push/Enqueue`, `Pop/Dequeue`, `Peek`, `Len`, `IsEmpty`. Verifique o comportamento quando vazia.

### 🔴 4. Anagrama
Escreva uma função que agrupa palavras que são anagramas entre si. Exemplo: `["eat","tea","tan","ate","nat","bat"]` → `[["eat","tea","ate"],["tan","nat"],["bat"]]`. Use um map com a palavra ordenada como chave.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Controle de Fluxo](../modulo-03-controle-fluxo/README.md) | [Próximo: Funções →](../modulo-05-funcoes/README.md)
