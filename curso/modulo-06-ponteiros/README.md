# Módulo 06 — Ponteiros

[← Funções](../modulo-05-funcoes/README.md) | [Próximo: Structs →](../modulo-07-structs/README.md)

---

## 📖 Teoria

### O que são ponteiros?
Um ponteiro armazena o **endereço de memória** de uma variável.

```go
x := 42
p := &x          // p aponta para x
fmt.Println(*p)  // 42 (lê o valor)
*p = 100         // altera x via ponteiro
```

| Operador | Significado |
|----------|------------|
| `&x` | Endereço de x ("referência") |
| `*p` | Valor apontado por p ("desreferência") |
| `*int` | Tipo "ponteiro para int" |

### Passagem por valor vs referência
```go
// Por VALOR — recebe uma cópia
func dobrarValor(n int) { n *= 2 }  // NÃO altera o original

// Por REFERÊNCIA — recebe o ponteiro
func dobrarRef(n *int) { *n *= 2 }  // ALTERA o original
```

### Quando usar ponteiros?
- ✅ Quando precisa **modificar** o valor original
- ✅ Quando a struct é **grande** (evita cópia)
- ✅ Receivers de métodos que modificam estado
- ❌ Para tipos pequenos (`int`, `bool`) — geralmente não vale

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo06_ponteiros.go` | Ponteiros, &, *, nil, passagem por referência |
| `exercicios/ex06_ponteiros.go` | 🏋️ Exercícios |

---

[← Funções](../modulo-05-funcoes/README.md) | [Próximo: Structs →](../modulo-07-structs/README.md)
