# Módulo 07 — Structs e Métodos

[← Ponteiros](../modulo-06-ponteiros/README.md) | [Próximo: Interfaces →](../modulo-08-interfaces/README.md)

---

> **Antes de ler — tente responder:**
> 1. Go não tem herança. Como você reutiliza comportamento entre tipos?
> 2. Qual a diferença entre `value receiver` e `pointer receiver` na prática?
> 3. O que são struct tags e onde são usadas?

---

## 1. Structs — Agrupando Dados Relacionados

```go
// Definição
type Pessoa struct {
    Nome   string
    Idade  int
    Email  string
}

// Inicialização — prefira com nome dos campos (não depende da ordem)
p1 := Pessoa{Nome: "Alice", Idade: 30, Email: "alice@go.dev"}
p2 := Pessoa{"Bob", 25, "bob@go.dev"}  // frágil — quebra se adicionar campo

// Structs são valores — copiadas inteiras na atribuição
p3 := p1        // cópia completa
p3.Nome = "Charlie"
fmt.Println(p1.Nome) // "Alice" — não mudou

// Ponteiro para struct — trabalha com o original
p4 := &p1
p4.Nome = "Charlie"  // Go faz auto-dereference: (*p4).Nome
fmt.Println(p1.Nome) // "Charlie" — mudou!
```

### Struct anônima — para uso único

```go
// Útil em testes e respostas JSON únicas
config := struct {
    Host string
    Port int
}{
    Host: "localhost",
    Port: 8080,
}

// Em testes table-driven
tests := []struct {
    input    string
    expected int
}{
    {"hello", 5},
    {"", 0},
}
```

---

## 2. Métodos — Funções Associadas a Tipos

```go
type Retangulo struct {
    Largura float64
    Altura  float64
}

// Value receiver — não modifica o estado, recebe CÓPIA
func (r Retangulo) Area() float64 {
    return r.Largura * r.Altura
}

func (r Retangulo) Perimetro() float64 {
    return 2 * (r.Largura + r.Altura)
}

// Pointer receiver — modifica o estado, recebe PONTEIRO
func (r *Retangulo) Escalar(fator float64) {
    r.Largura *= fator
    r.Altura *= fator
}

// Uso — Go faz conversão automática entre valor e ponteiro
r := Retangulo{Largura: 5, Altura: 3}
fmt.Println(r.Area())       // 15 — value receiver
r.Escalar(2)                // Go converte automaticamente para (&r).Escalar(2)
fmt.Println(r.Area())       // 60

// Também funciona com ponteiro
rp := &Retangulo{Largura: 5, Altura: 3}
rp.Area()    // Go converte: (*rp).Area()
rp.Escalar(2)
```

### Regra de ouro para receivers

```go
// ✅ Use pointer receiver quando:
// 1. O método modifica o estado
func (u *User) SetEmail(email string) { u.Email = email }

// 2. A struct é grande (evita cópia)
func (c *Config) Validate() error { /* ... */ }

// 3. Consistência — se um método usa pointer, TODOS devem usar pointer

// ✅ Value receiver é ok quando:
// 1. A struct é pequena e método só lê
type Point struct{ X, Y float64 }
func (p Point) Distance() float64 { return math.Sqrt(p.X*p.X + p.Y*p.Y) }

// ❌ NUNCA misture sem motivo — causa confusão com interfaces
```

---

## 3. Constructor Pattern — `NewXxx`

Go não tem construtores, mas o padrão `NewXxx` é universal:

```go
type Server struct {
    host    string
    port    int
    timeout time.Duration
    db      *sql.DB
}

// Constructor com validação e defaults
func NewServer(host string, port int, db *sql.DB) (*Server, error) {
    if host == "" {
        return nil, errors.New("host não pode ser vazio")
    }
    if port < 1 || port > 65535 {
        return nil, fmt.Errorf("porta inválida: %d", port)
    }
    return &Server{
        host:    host,
        port:    port,
        timeout: 30 * time.Second, // default
        db:      db,
    }, nil
}

// Functional options — para structs com muitas opções opcionais
type Option func(*Server)

func WithTimeout(d time.Duration) Option {
    return func(s *Server) { s.timeout = d }
}

func NewServerWithOptions(host string, port int, opts ...Option) *Server {
    s := &Server{host: host, port: port, timeout: 30 * time.Second}
    for _, opt := range opts {
        opt(s)
    }
    return s
}

// Uso elegante
srv := NewServerWithOptions("localhost", 8080,
    WithTimeout(60*time.Second),
)
```

---

## 4. Composição (Embedding) — A "Herança" do Go

Go não tem herança. Em vez disso, usa **composição** — um tipo pode "embutir" outro e herdar seus métodos:

```go
type Animal struct {
    Nome string
    Peso float64
}

func (a Animal) Respirar() string {
    return fmt.Sprintf("%s está respirando", a.Nome)
}

// Cachorro EMBUTE Animal (não herda — compõe)
type Cachorro struct {
    Animal         // embedding — sem nome de campo
    Raca   string
}

func (c Cachorro) Latir() string {
    return "Au au!"
}

// Uso — campos e métodos de Animal são promovidos
c := Cachorro{
    Animal: Animal{Nome: "Rex", Peso: 25.0},
    Raca:   "Labrador",
}
fmt.Println(c.Nome)       // "Rex" — promovido de Animal
fmt.Println(c.Respirar()) // "Rex está respirando" — método promovido
fmt.Println(c.Latir())    // "Au au!" — método próprio
fmt.Println(c.Animal.Nome) // também funciona — acesso explícito

// Override de método
func (c Cachorro) Respirar() string {
    return fmt.Sprintf("%s respira e late", c.Nome)
}
// Agora c.Respirar() usa o do Cachorro, não do Animal
```

### Composição vs Herança

```
Herança (Java/Python):     Composição (Go):
class Dog extends Animal   type Dog struct {
                               Animal
                           }

Diferença crítica:
- Herança: Dog É um Animal
- Composição: Dog TEM um Animal (e promove seus métodos)

Implicação prática:
- Em Go, Dog NÃO satisfaz automaticamente interfaces que Animal satisfaz
  (a menos que os métodos sejam promovidos)
```

---

## 5. Struct Tags — Metadados para Serialização

```go
type User struct {
    ID        int       `json:"id"`
    Nome      string    `json:"nome"`
    Email     string    `json:"email,omitempty"`   // omitempty: omite se zero value
    Senha     string    `json:"-"`                 // - : nunca serializa
    CriadoEm time.Time `json:"criado_em"`
    
    // Tags para múltiplos frameworks ao mesmo tempo
    Cidade    string    `json:"cidade" db:"cidade" validate:"required,min=2"`
}

u := User{ID: 1, Nome: "Alice", Email: "alice@go.dev"}
data, _ := json.Marshal(u)
// {"id":1,"nome":"Alice","email":"alice@go.dev","criado_em":"0001-01-01T00:00:00Z"}
// Senha não aparece (tag "-")
// Cidade não aparece (omitempty + zero value)

// Desserialização
var u2 User
json.Unmarshal(data, &u2)
```

---

## ✅ Checklist de Structs

- [ ] Constructors (`NewXxx`) validam e retornam `(*T, error)`
- [ ] Campos privados (minúscula) são protegidos — acesso apenas via métodos
- [ ] Receivers são consistentes em toda a struct (todos pointer OU todos value)
- [ ] Struct tags usam `omitempty` para campos opcionais e `"-"` para campos sensíveis (senha, token)
- [ ] Composição via embedding preferida à duplicação de código
- [ ] Functional options para structs com muitas configurações opcionais

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo07_structs.go` | Structs, métodos, composição, JSON tags, functional options |
| `exercicios/ex07_structs.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Conta Bancária
Crie uma struct `ContaBancaria` com campos privados (`saldo`, `titular`). Implemente métodos `Depositar`, `Sacar` (com validação de saldo), `Saldo() float64` e `String() string`. Teste que não é possível acessar o saldo diretamente.

### 🟡 2. Composição de Formas Geométricas
Crie structs `Ponto`, `Circulo` (com Ponto central) e `Retangulo` (com dois Pontos). Implemente métodos `Area()`, `Perimetro()` e `Contem(p Ponto) bool`. Demonstre a composição vs duplicação.

### 🟡 3. HTTP Server com Functional Options
Crie um tipo `HTTPServer` com configurações: `host`, `port`, `readTimeout`, `writeTimeout`, `maxConnections`. Implemente o padrão `WithXxx` para cada opção. O constructor deve ter defaults sensatos.

### 🔴 4. Serialização Customizada
Implemente `MarshalJSON()` e `UnmarshalJSON()` para uma struct `Dinheiro` que armazena centavos internamente mas serializa como `"R$ 19,99"`. Use struct tags para controlar quais campos são expostos.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Ponteiros](../modulo-06-ponteiros/README.md) | [Próximo: Interfaces →](../modulo-08-interfaces/README.md)
