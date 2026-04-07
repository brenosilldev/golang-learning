# Módulo 06 — Ponteiros

[← Funções](../modulo-05-funcoes/README.md) | [Próximo: Structs →](../modulo-07-structs/README.md)

---

> **Antes de ler — tente responder:**
> 1. Qual a diferença entre passar por valor e por referência em Go?
> 2. Quando devo usar ponteiro receiver vs value receiver em métodos?
> 3. O que acontece se você desreferenciar um ponteiro `nil`?

---

## 1. O Que São Ponteiros

Um ponteiro armazena o **endereço de memória** de uma variável. Em vez de copiar o valor, você passa o endereço — quem recebe pode modificar o original.

```go
x := 42
p := &x          // p = endereço de x (ex: 0xc000018080)
fmt.Println(p)   // 0xc000018080 (o endereço)
fmt.Println(*p)  // 42 (o valor — desreferência)

*p = 100         // modifica x através do ponteiro
fmt.Println(x)   // 100
```

| Operador | Nome | Significado |
|----------|------|-------------|
| `&x` | address-of | "me dê o endereço de x" |
| `*p` | dereference | "me dê o valor que p aponta" |
| `*int` | pointer type | "tipo: ponteiro para int" |

---

## 2. Passagem por Valor vs Referência

Go passa **tudo por valor** por padrão. Ponteiros são a forma de passar "por referência":

```go
// Por VALOR — recebe uma cópia
func dobrarValor(n int) {
    n *= 2
    // n aqui é uma cópia — não afeta o original
}

// Por REFERÊNCIA — recebe o endereço
func dobrarRef(n *int) {
    *n *= 2  // modifica o original via ponteiro
}

func main() {
    x := 10
    dobrarValor(x)
    fmt.Println(x)  // 10 — não mudou

    dobrarRef(&x)
    fmt.Println(x)  // 20 — mudou!
}
```

### Quando structs são copiadas inteiras

```go
type Config struct {
    Host     string
    Port     int
    MaxConns int
    // ... 20 campos
}

// ❌ Copia a struct inteira na chamada (ineficiente para structs grandes)
func inicializar(cfg Config) {
    cfg.MaxConns = 100 // modifica a cópia — original não muda
}

// ✅ Passa ponteiro — sem cópia, pode modificar o original
func inicializar(cfg *Config) {
    cfg.MaxConns = 100 // modifica o original
}
```

---

## 3. Nil Pointer — A Armadilha Mais Comum

```go
var p *int // p é nil (não aponta para nada)
fmt.Println(p)  // <nil>
fmt.Println(*p) // PANIC: runtime error: invalid memory address or nil pointer dereference

// ✅ Sempre verifique nil antes de desreferenciar
if p != nil {
    fmt.Println(*p)
}

// Structs com ponteiro nil — MUITO comum em código real
type Node struct {
    Value int
    Next  *Node
}

var head *Node
head.Value // PANIC! head é nil

// ✅ Verifique antes de acessar
if head != nil {
    fmt.Println(head.Value)
}
```

---

## 4. new vs make — Qual Usar?

```go
// new(T) — aloca zero value de T, retorna *T
p := new(int)       // *p = 0
s := new(string)    // *s = ""
cfg := new(Config)  // todos os campos zerados

// make(T, ...) — APENAS para slice, map e channel
// Inicializa a estrutura interna (não apenas aloca)
sl := make([]int, 5)         // slice com len=5, cap=5
m := make(map[string]int)    // map inicializado (pronto para uso)
ch := make(chan int, 10)      // channel com buffer de 10

// Regra prática:
// - Para structs: use literal {} ou &Config{...}
// - Para slice/map/chan: use make()
// - new() raramente é necessário

// Preferido para structs:
cfg2 := &Config{Host: "localhost", Port: 8080}
// ao invés de:
cfg3 := new(Config)
cfg3.Host = "localhost"
cfg3.Port = 8080
```

---

## 5. Ponteiros em Structs — O Pattern Mais Comum

```go
type User struct {
    ID    int
    Nome  string
    Email string
}

// Constructor retorna ponteiro (evita cópia, permite modificação)
func NewUser(nome, email string) *User {
    return &User{Nome: nome, Email: email}
}

// Go faz auto-dereference em campos de struct
u := NewUser("Alice", "alice@go.dev")
fmt.Println(u.Nome)   // funciona! Go faz (*u).Nome automaticamente
u.Nome = "Bob"        // modifica o original via ponteiro

// Ponteiro para campo
email := &u.Email     // ponteiro para o campo Email
*email = "bob@go.dev" // modifica diretamente
fmt.Println(u.Email)  // "bob@go.dev"
```

---

## 6. Quando Usar (e Não Usar) Ponteiros

```
Devo usar ponteiro?
│
├── Preciso MODIFICAR o valor original?
│   └── Sim → use ponteiro
│
├── A struct é GRANDE (muitos campos)?
│   └── Sim → use ponteiro (evita cópia cara)
│
├── Precisa ser NULLABLE (pode não existir)?
│   └── Sim → use *T (nil = não existe)
│
├── É um tipo PEQUENO (int, bool, small struct)?
│   └── Não → use por valor (mais eficiente)
│
└── Para receivers de métodos:
    ├── Método MODIFICA o estado → pointer receiver
    └── Método só LÊ → value receiver (mas seja consistente!)
```

```go
// ✅ Pointer receiver quando modifica estado
func (u *User) SetNome(nome string) {
    u.Nome = nome
}

// ✅ Value receiver quando só lê
func (u User) String() string {
    return fmt.Sprintf("%s <%s>", u.Nome, u.Email)
}

// ⚠️ Regra de ouro: seja consistente!
// Se um método de User usa pointer receiver,
// TODOS os métodos de User devem usar pointer receiver
// (misturar pode causar bugs sutis com interfaces)
```

---

## 7. Armadilha: Retornar Ponteiro para Variável Local

```go
// ✅ Seguro em Go — o compilador garante que a variável sobrevive
// Go usa escape analysis: se a variável "escapa" da função, vai para o heap
func newInt(n int) *int {
    n := n  // variável local
    return &n  // Go move para o heap automaticamente — seguro!
}

// ✅ Pattern comum em Go
func NewConfig() *Config {
    cfg := Config{Host: "localhost", Port: 8080}
    return &cfg  // safe — Go faz escape analysis
}
```

---

## ✅ Checklist de Ponteiros

- [ ] Verifico `nil` antes de desreferenciar ponteiros vindos de fora
- [ ] Uso `*T` quando o valor pode não existir (zero value não faria sentido)
- [ ] Receivers de métodos são consistentes: todos pointer OU todos value
- [ ] Uso `make` para slice/map/chan (não `new`)
- [ ] Funções que modificam o receptor usam pointer receiver

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo06_ponteiros.go` | Ponteiros, &, *, nil, passagem por referência, new vs make |
| `exercicios/ex06_ponteiros.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Troca de Valores
Escreva uma função `Swap(a, b *int)` que troca os valores de duas variáveis via ponteiro. Por que não dá para fazer isso com passagem por valor?

### 🟢 2. Lista Encadeada
Implemente uma lista encadeada simples com `*Node`. Adicione métodos `Append(val int)`, `Len() int` e `Print()`. Demonstre o comportamento quando a lista está vazia (head == nil).

### 🟡 3. Optional Pattern
Em Go, `*T` pode ser usado como "optional" (similar a `Option<T>` em Rust). Implemente uma struct `Config` onde campos opcionais são `*string` ou `*int`. Crie funções auxiliares `StringPtr(s string) *string` e `IntPtr(n int) *int` para facilitar a criação.

### 🟡 4. Benchmark: Valor vs Ponteiro
Escreva dois benchmarks (`go test -bench=.`): um que passa uma struct grande (10+ campos) por valor, outro por ponteiro. Observe a diferença de performance com `go test -bench=. -benchmem`.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Funções](../modulo-05-funcoes/README.md) | [Próximo: Structs →](../modulo-07-structs/README.md)
