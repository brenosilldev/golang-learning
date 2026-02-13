# Módulo 09 — Tratamento de Erros

[← Interfaces](../modulo-08-interfaces/README.md) | [Próximo: Pacotes e Módulos →](../modulo-10-pacotes-modulos/README.md)

---

## 📖 Teoria

### O pattern fundamental do Go
```go
resultado, err := funcaoQuePodeFalhar()
if err != nil {
    // tratar o erro
    return err
}
// usar resultado
```

Go **não tem exceções**. Erros são valores retornados — simples e explícito.

### Criando erros
```go
errors.New("algo deu errado")           // erro simples
fmt.Errorf("falha no item %d: %w", id, err) // com contexto + wrap
```

### Custom errors
```go
type NotFoundError struct {
    Recurso string
    ID      int
}
func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s #%d não encontrado", e.Recurso, e.ID)
}
```

### errors.Is e errors.As (Go 1.13+)
```go
errors.Is(err, ErrNotFound)               // compara erro
errors.As(err, &target)                     // extrai tipo de erro
```

### panic / recover — APENAS para erros irrecuperáveis
```go
panic("bug impossível")  // nunca use para erros de negócio!
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo09_erros.go` | Todos os patterns de erro |
| `exercicios/ex09_erros.go` | 🏋️ Exercícios |

---

[← Interfaces](../modulo-08-interfaces/README.md) | [Próximo: Pacotes e Módulos →](../modulo-10-pacotes-modulos/README.md)
