# Módulo 12 — Generics

[← Concorrência](../modulo-11-concorrencia/README.md) | [Próximo: Testes →](../modulo-13-testes/README.md)

---

## 📖 Teoria

### Generics (Go 1.18+)
Funções e tipos que trabalham com múltiplos tipos:

```go
func Map[T any, U any](slice []T, fn func(T) U) []U {
    result := make([]U, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}
```

### Constraints (restrições de tipo)
```go
type Numero interface {
    int | int32 | int64 | float32 | float64
}

func Soma[T Numero](a, b T) T { return a + b }
```

### Constraints built-in
| Constraint | Significado |
|-----------|------------|
| `any` | Qualquer tipo (= `interface{}`) |
| `comparable` | Tipos comparáveis com `==` |
| `constraints.Ordered` | Tipos ordenáveis com `<` `>` |

### Structs genéricas
```go
type Stack[T any] struct { items []T }
func (s *Stack[T]) Push(item T) { s.items = append(s.items, item) }
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo12_generics.go` | Funções e tipos genéricos |
| `exercicios/ex12_generics.go` | 🏋️ Exercícios |

---

[← Concorrência](../modulo-11-concorrencia/README.md) | [Próximo: Testes →](../modulo-13-testes/README.md)
