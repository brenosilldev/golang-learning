# Módulo 08 — Interfaces

[← Structs](../modulo-07-structs/README.md) | [Próximo: Tratamento de Erros →](../modulo-09-tratamento-erros/README.md)

---

## 📖 Teoria

### Interfaces são IMPLÍCITAS em Go
Não precisa declarar `implements`. Se o tipo tem os métodos, ele implementa a interface.

```go
type Animal interface {
    Falar() string
}

type Cachorro struct{}
func (c Cachorro) Falar() string { return "Au au!" }
// Cachorro implementa Animal automaticamente!
```

### Interface vazia `any` (alias de `interface{}`)
Aceita qualquer tipo:
```go
func imprimir(v any) { fmt.Println(v) }
```

### Type assertion e type switch
```go
var i interface{} = "hello"
s := i.(string)                    // type assertion
s, ok := i.(string)                // safe assertion

switch v := i.(type) {             // type switch
case string: fmt.Println("string:", v)
case int:    fmt.Println("int:", v)
}
```

### Interfaces da stdlib
| Interface | Métodos | Uso |
|-----------|---------|-----|
| `fmt.Stringer` | `String() string` | Representação em texto |
| `error` | `Error() string` | Tratamento de erros |
| `io.Reader` | `Read([]byte) (int, error)` | Ler dados |
| `io.Writer` | `Write([]byte) (int, error)` | Escrever dados |
| `sort.Interface` | `Len()`, `Less()`, `Swap()` | Ordenação customizada |

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo08_interfaces.go` | Interfaces, polimorfismo, type assertion, stdlib |
| `exercicios/ex08_interfaces.go` | 🏋️ Exercícios |

---

[← Structs](../modulo-07-structs/README.md) | [Próximo: Tratamento de Erros →](../modulo-09-tratamento-erros/README.md)
