# Módulo 04 — Coleções

[← Controle de Fluxo](../modulo-03-controle-fluxo/README.md) | [Próximo: Funções →](../modulo-05-funcoes/README.md)

---

## 📖 Teoria

### Arrays — Tamanho fixo
```go
var arr [5]int                    // [0 0 0 0 0]
arr2 := [3]string{"a", "b", "c"} // tamanho fixo = 3
```
Arrays são **raramente usados** em Go. Use **slices**.

### Slices — Tamanho dinâmico (o que você vai usar 99% do tempo)
```go
s := []int{1, 2, 3}           // literal
s2 := make([]int, 5)          // com make (len=5, cap=5)
s3 := make([]int, 0, 10)      // len=0, cap=10
s = append(s, 4, 5)           // adicionar elementos
sub := s[1:3]                 // slice do slice [2, 3]
```

**len vs cap**: `len` = itens usados, `cap` = capacidade alocada. Quando `len == cap`, `append` aloca um novo array maior (geralmente 2x).

### Maps — Chave/Valor
```go
m := map[string]int{"go": 1, "rust": 2}
m["python"] = 3                   // inserir
valor, ok := m["go"]              // leitura segura
delete(m, "rust")                 // deletar
```

### Sets — Go não tem set, mas usa `map[T]struct{}`
```go
set := map[string]struct{}{}
set["item"] = struct{}{}
_, existe := set["item"] // true
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo04_colecoes.go` | Arrays, slices, maps, operações |
| `exercicios/ex04_colecoes.go` | 🏋️ Exercícios |

---

[← Controle de Fluxo](../modulo-03-controle-fluxo/README.md) | [Próximo: Funções →](../modulo-05-funcoes/README.md)
