# Módulo 07 — Structs e Métodos

[← Ponteiros](../modulo-06-ponteiros/README.md) | [Próximo: Interfaces →](../modulo-08-interfaces/README.md)

---

## 📖 Teoria

### Structs — Agrupar dados
```go
type Pessoa struct {
    Nome  string
    Idade int
}
p := Pessoa{Nome: "Go", Idade: 15}
```

### Métodos — Funções associadas a um tipo
```go
// Value receiver — NÃO modifica
func (p Pessoa) Apresentar() string { return p.Nome }

// Pointer receiver — PODE modificar
func (p *Pessoa) Aniversario() { p.Idade++ }
```

### Composição (embedding) — "Herança" do Go
```go
type Animal struct { Nome string }
type Cachorro struct {
    Animal       // embedding — Cachorro "herda" Animal
    Raca string
}
c := Cachorro{Animal: Animal{Nome: "Rex"}, Raca: "Labrador"}
c.Nome // funciona direto!
```

### JSON Tags
```go
type User struct {
    Nome  string `json:"nome"`
    Email string `json:"email,omitempty"`
}
```

### Pattern: Construtor `NewXxx`
```go
func NewPessoa(nome string, idade int) *Pessoa {
    return &Pessoa{Nome: nome, Idade: idade}
}
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo07_structs.go` | Structs, métodos, composição, JSON |
| `exercicios/ex07_structs.go` | 🏋️ Exercícios |

---

[← Ponteiros](../modulo-06-ponteiros/README.md) | [Próximo: Interfaces →](../modulo-08-interfaces/README.md)
