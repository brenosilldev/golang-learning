# Módulo 02 — Variáveis e Tipos

[← Voltar ao índice](../README.md) | [Próximo: Controle de Fluxo →](../modulo-03-controle-fluxo/README.md)

---

## 📖 Teoria

### Declaração de Variáveis
Go tem **duas formas** de declarar variáveis:

```go
// Forma 1: var (pode usar em qualquer escopo)
var nome string = "Go"
var idade int = 25

// Forma 2: := (apenas dentro de funções — mais usada!)
nome := "Go"
idade := 25
```

### Tipos Primitivos
| Tipo | Exemplo | Zero Value |
|------|---------|------------|
| `int`, `int8`, `int16`, `int32`, `int64` | `42` | `0` |
| `uint`, `uint8`, `uint16`, `uint32`, `uint64` | `42` | `0` |
| `float32`, `float64` | `3.14` | `0.0` |
| `string` | `"hello"` | `""` |
| `bool` | `true` / `false` | `false` |
| `byte` (alias de `uint8`) | `'A'` | `0` |
| `rune` (alias de `int32`) | `'世'` | `0` |

### Zero Values — Tudo tem valor padrão
Em Go, **variáveis nunca são "undefined"**. Cada tipo tem um **zero value**:
```go
var i int      // 0
var f float64  // 0.0
var s string   // ""
var b bool     // false
```

### Constantes e `iota`
```go
const Pi = 3.14159
const (
    Domingo = iota  // 0
    Segunda         // 1
    Terca           // 2
)
```

### Conversão de Tipos
Go **não faz conversão implícita** — você precisa ser explícito:
```go
i := 42
f := float64(i)       // int → float64
s := string(rune(65)) // int → rune → string ("A")
```

---

## 📂 Arquivos deste módulo

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo02_variaveis.go` | Declaração, tipos, zero values, conversão |
| `exercicios/ex02_tipos.go` | 🏋️ Exercícios para praticar |

---

[← Voltar ao índice](../README.md) | [Próximo: Controle de Fluxo →](../modulo-03-controle-fluxo/README.md)
