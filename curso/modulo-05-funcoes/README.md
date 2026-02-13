# Módulo 05 — Funções

[← Coleções](../modulo-04-colecoes/README.md) | [Próximo: Ponteiros →](../modulo-06-ponteiros/README.md)

---

## 📖 Teoria

### Funções básicas
```go
func soma(a, b int) int {
    return a + b
}
```

### Múltiplos retornos
```go
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("divisão por zero")
    }
    return a / b, nil
}
```

### Retornos nomeados
```go
func retornoNomeado(a, b int) (soma int, diff int) {
    soma = a + b
    diff = a - b
    return // naked return
}
```

### Funções variádicas
```go
func somaTotal(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}
somaTotal(1, 2, 3, 4, 5) // 15
```

### Closures
```go
func contador() func() int {
    i := 0
    return func() int {
        i++
        return i
    }
}
c := contador()
c() // 1
c() // 2
```

### defer
Executa ao final da função, na ordem LIFO (pilha):
```go
func exemplo() {
    defer fmt.Println("3º")
    defer fmt.Println("2º")
    defer fmt.Println("1º")
}
// Imprime: 1º, 2º, 3º
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo05_funcoes.go` | Todos os tipos de funções |
| `exercicios/ex05_funcoes.go` | 🏋️ Exercícios |

---

[← Coleções](../modulo-04-colecoes/README.md) | [Próximo: Ponteiros →](../modulo-06-ponteiros/README.md)
