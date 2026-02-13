# Módulo 11 — Concorrência

[← Pacotes](../modulo-10-pacotes-modulos/README.md) | [Próximo: Generics →](../modulo-12-generics/README.md)

---

## 📖 Teoria

### Goroutines
Leves (~2KB), gerenciadas pelo runtime do Go (não pelo OS):
```go
go minhaFuncao()         // lança em paralelo
go func() { ... }()     // goroutine anônima
```

### Channels — Comunicação entre goroutines
```go
ch := make(chan int)       // unbuffered (bloqueia)
ch := make(chan int, 10)   // buffered (cap 10)
ch <- 42                   // enviar
valor := <-ch              // receber
close(ch)                  // fechar
```

### Select — "switch" para channels
```go
select {
case msg := <-ch1: ...
case ch2 <- valor: ...
case <-time.After(5*time.Second): ... // timeout
default: ...                           // non-blocking
}
```

### sync.WaitGroup — Esperar goroutines terminarem
```go
var wg sync.WaitGroup
wg.Add(1)
go func() { defer wg.Done(); ... }()
wg.Wait()
```

### Patterns
- **Fan-out**: 1 producer → N workers
- **Fan-in**: N producers → 1 consumer
- **Worker Pool**: N goroutines processando uma fila
- **Pipeline**: Estágio A → Channel → Estágio B → Channel → Estágio C

### context.Context — Cancelamento e timeouts
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo11_concorrencia.go` | Goroutines, channels, select, patterns |
| `exercicios/ex11_concorrencia.go` | 🏋️ Exercícios |

---

[← Pacotes](../modulo-10-pacotes-modulos/README.md) | [Próximo: Generics →](../modulo-12-generics/README.md)
