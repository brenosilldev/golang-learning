# Módulo 13 — Testes

[← Generics](../modulo-12-generics/README.md) | [Próximo: APIs →](../modulo-14-apis/README.md)

---

## 📖 Teoria

### Convenções de teste em Go
```
arquivo.go      → arquivo_test.go     (mesmo pacote)
func Soma(...)  → func TestSoma(...)   (prefixo Test)
```

### testing.T
```go
func TestSoma(t *testing.T) {
    resultado := Soma(2, 3)
    if resultado != 5 {
        t.Errorf("Soma(2,3) = %d, esperado 5", resultado)
    }
}
```

### Table-driven tests (pattern mais importante!)
```go
func TestSoma(t *testing.T) {
    tests := []struct{ a, b, esperado int }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }
    for _, tt := range tests {
        resultado := Soma(tt.a, tt.b)
        if resultado != tt.esperado {
            t.Errorf("Soma(%d,%d) = %d, esperado %d", tt.a, tt.b, resultado, tt.esperado)
        }
    }
}
```

### Subtestes com t.Run
```go
t.Run("nome_do_teste", func(t *testing.T) { ... })
```

### Benchmarks
```go
func BenchmarkSoma(b *testing.B) {
    for i := 0; i < b.N; i++ { Soma(1, 2) }
}
// go test -bench=. -benchmem
```

### Coverage
```bash
go test -cover ./...
go test -coverprofile=cover.out && go tool cover -html=cover.out
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo13_codigo.go` | Código que será testado |
| `exemplos/exemplo13_codigo_test.go` | Testes completos demonstrando patterns |
| `exercicios/ex13_testes_test.go` | 🏋️ Exercícios de teste |

---

[← Generics](../modulo-12-generics/README.md) | [Próximo: APIs →](../modulo-14-apis/README.md)
