# Módulo 08 — Interfaces

[← Structs](../modulo-07-structs/README.md) | [Próximo: Tratamento de Erros →](../modulo-09-tratamento-erros/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que Go não tem a keyword `implements`?
> 2. Qual o risco de usar `interface{}` (ou `any`) em todo lugar?
> 3. Por que `io.Reader` tem apenas 1 método?

---

## 1. Interfaces São Implícitas — O Poder do Go

Em Java/C#, você declara `class Dog implements Animal`. Em Go, **não precisa**. Se o tipo tem os métodos, ele implementa a interface automaticamente:

```go
type Animal interface {
    Falar() string
}

type Cachorro struct{}
func (c Cachorro) Falar() string { return "Au au!" }

type Gato struct{}
func (g Gato) Falar() string { return "Miau!" }

// Ambos implementam Animal — sem declarar "implements" em lugar nenhum
func fazerBarulho(a Animal) {
    fmt.Println(a.Falar())
}

func main() {
    fazerBarulho(Cachorro{}) // funciona
    fazerBarulho(Gato{})     // funciona
}
```

**Por que isso é poderoso?** Você pode criar interfaces **para tipos que não são seus**. Um tipo da stdlib pode satisfazer sua interface sem que os autores saibam:

```go
type Sizer interface {
    Len() int
}

// strings.Builder, bytes.Buffer, slices — todos implementam Sizer
// sem que os autores dessas libs tenham pensado nisso
```

---

## 2. A Regra de Ouro: "Accept Interfaces, Return Structs"

Este é o princípio de design mais importante para interfaces em Go:

```go
// ✅ CORRETO — aceita interface, retorna struct
func ProcessarDados(r io.Reader) (*Resultado, error) {
    data, err := io.ReadAll(r)
    if err != nil {
        return nil, err
    }
    return &Resultado{Data: data}, nil
}

// Funciona com arquivo, HTTP body, buffer, string... qualquer io.Reader
ProcessarDados(os.Stdin)
ProcessarDados(resp.Body)
ProcessarDados(bytes.NewReader(data))
ProcessarDados(strings.NewReader("teste"))
```

```go
// ❌ ERRADO — retorna interface (esconde o tipo real)
func NovoUsuario(nome string) UserInterface {
    return &User{Nome: nome}
}

// ✅ CORRETO — retorna o tipo concreto
func NovoUsuario(nome string) *User {
    return &User{Nome: nome}
}
```

**Por que?**
- **Aceitar interfaces** = flexibilidade para quem chama (pode passar qualquer tipo que satisfaça)
- **Retornar structs** = clareza sobre o que é retornado (o chamador vê os campos e métodos)

---

## 3. Interfaces Pequenas — A Filosofia Go

Go preza por interfaces com **poucos métodos** (idealmente 1-2). Isso maximiza a chance de tipos diferentes satisfazerem a interface:

```go
// Interfaces da stdlib — observe o tamanho
type Reader interface {
    Read(p []byte) (n int, err error) // 1 método
}

type Writer interface {
    Write(p []byte) (n int, err error) // 1 método
}

type Closer interface {
    Close() error // 1 método
}

// Composição: interfaces grandes são COMPOSTAS de pequenas
type ReadWriter interface {
    Reader
    Writer
}

type ReadWriteCloser interface {
    Reader
    Writer
    Closer
}
```

### Quantos tipos implementam cada interface?

| Tamanho da interface | Exemplo | Tipos que implementam |
|----------------------|---------|-----------------------|
| 1 método (`io.Reader`) | `Read([]byte) (int, error)` | **Centenas**: File, Buffer, Body, Conn, Pipe... |
| 2 métodos (`io.ReadWriter`) | Reader + Writer | Dezenas |
| 5+ métodos | Interface grande customizada | Provavelmente só 1-2 tipos |

> **Provérbio Go**: *"The bigger the interface, the weaker the abstraction."* — Rob Pike

```go
// ❌ Interface grande demais — provavelmente só 1 tipo implementa
type UserManager interface {
    Create(user *User) error
    Update(user *User) error
    Delete(id int) error
    FindByID(id int) (*User, error)
    FindByEmail(email string) (*User, error)
    List(page, size int) ([]User, error)
    Count() (int, error)
}

// ✅ Melhor: quebre em interfaces menores por capacidade
type UserReader interface {
    FindByID(ctx context.Context, id int) (*User, error)
    List(ctx context.Context) ([]User, error)
}

type UserWriter interface {
    Create(ctx context.Context, user *User) error
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}

// O handler de listagem só precisa de leitura
func NewListHandler(repo UserReader) http.Handler { ... }

// O handler de criação só precisa de escrita
func NewCreateHandler(repo UserWriter) http.Handler { ... }
```

---

## 4. Interfaces da stdlib — As Mais Importantes

| Interface | Métodos | Uso real | Quem implementa |
|-----------|---------|----------|-----------------|
| `fmt.Stringer` | `String() string` | Representação em texto para `fmt.Println` | Qualquer tipo que queira se "apresentar" |
| `error` | `Error() string` | Todo tratamento de erros em Go | Custom errors, sentinel errors |
| `io.Reader` | `Read([]byte) (int, error)` | Ler dados de qualquer fonte | File, Buffer, HTTP Body, Conn |
| `io.Writer` | `Write([]byte) (int, error)` | Escrever dados para qualquer destino | File, Buffer, HTTP ResponseWriter |
| `io.Closer` | `Close() error` | Liberar recursos | File, DB rows, HTTP Body |
| `json.Marshaler` | `MarshalJSON() ([]byte, error)` | Customizar serialização JSON | Tipos com formato especial |
| `http.Handler` | `ServeHTTP(w, r)` | Tratar requisições HTTP | Handlers, middleware, mux |

### io.Reader e io.Writer na prática

```go
// Função que funciona com QUALQUER fonte de dados
func contarLinhas(r io.Reader) (int, error) {
    scanner := bufio.NewScanner(r)
    count := 0
    for scanner.Scan() {
        count++
    }
    return count, scanner.Err()
}

// Funciona com arquivo
file, _ := os.Open("dados.txt")
n, _ := contarLinhas(file)

// Funciona com HTTP response
resp, _ := http.Get("https://example.com")
n, _ = contarLinhas(resp.Body)

// Funciona com string
n, _ = contarLinhas(strings.NewReader("linha1\nlinha2\nlinha3"))

// Funciona com buffer
var buf bytes.Buffer
buf.WriteString("teste\naqui")
n, _ = contarLinhas(&buf)
```

---

## 5. Type Assertion e Type Switch

### Type assertion (extrair o tipo concreto)

```go
var i interface{} = "hello"

// Assertion UNSAFE — panic se o tipo estiver errado
s := i.(string)

// Assertion SAFE — SEMPRE use esta
s, ok := i.(string)
if !ok {
    fmt.Println("não é string")
}
```

### Type switch (quando você precisa agir diferente por tipo)

```go
func descrever(v any) string {
    switch val := v.(type) {
    case string:
        return fmt.Sprintf("string de %d chars", len(val))
    case int:
        return fmt.Sprintf("inteiro: %d", val)
    case bool:
        return fmt.Sprintf("booleano: %t", val)
    case nil:
        return "nil"
    default:
        return fmt.Sprintf("tipo desconhecido: %T", val)
    }
}
```

> **Cuidado**: se você usa type switch com frequência, provavelmente deveria usar uma interface com métodos ao invés de `any`. Type switches são um sinal de que falta abstração.

---

## 6. Verificação em Tempo de Compilação

```go
// Garante que PostgresUserRepo implementa UserRepository
// Se não implementar, o código NÃO COMPILA
var _ UserRepository = (*PostgresUserRepo)(nil)

// Isso é o pattern de "compile-time interface check"
// Coloque no topo do arquivo da implementação
```

---

## 7. Interfaces para Testes (Mocks)

O poder real das interfaces aparece nos testes. Com uma interface, você substitui o banco real por um mock:

```go
// Interface
type EmailSender interface {
    Send(to, subject, body string) error
}

// Implementação real (produção)
type SMTPSender struct { /* config SMTP */ }
func (s *SMTPSender) Send(to, subject, body string) error {
    // envia email de verdade
}

// Mock (testes)
type MockSender struct {
    Calls []struct{ To, Subject, Body string }
    Err   error // configura erro simulado
}

func (m *MockSender) Send(to, subject, body string) error {
    m.Calls = append(m.Calls, struct{ To, Subject, Body string }{to, subject, body})
    return m.Err
}

// Teste
func TestRegistro(t *testing.T) {
    mock := &MockSender{}
    svc := NewUserService(mock) // injeta o mock

    svc.Registrar("alice@test.com", "Alice")

    if len(mock.Calls) != 1 {
        t.Errorf("esperava 1 email, enviou %d", len(mock.Calls))
    }
    if mock.Calls[0].To != "alice@test.com" {
        t.Errorf("email enviado para %s, esperava alice@test.com", mock.Calls[0].To)
    }
}
```

> **No mercado**: a habilidade de "desacoplar com interfaces para facilitar testes" é o que separa um dev Go pleno de um júnior. Entrevistadores perguntam isso com frequência.

---

## ✅ Checklist de Interfaces para Produção

- [ ] Interfaces são **pequenas** (1-3 métodos — se tem mais, considere quebrar)
- [ ] Interfaces são definidas **pelo consumidor** (quem usa), não pelo implementador
- [ ] Funções **aceitam interfaces, retornam structs**
- [ ] Type assertions usam a forma **safe** (`v, ok := i.(Type)`)
- [ ] Verificação de implementação em **tempo de compilação** (`var _ I = (*T)(nil)`)
- [ ] Dependências externas (banco, email, API) têm **interface para mocking** em testes
- [ ] `any` é usado com **moderação** — prefira interfaces tipadas

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo08_interfaces.go` | Interfaces, polimorfismo, type assertion, stdlib |
| `exercicios/ex08_interfaces.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Formas Geométricas
Crie a interface `Forma` com `Area()` e `Perimetro()`. Implemente para `Retangulo`, `Circulo` e `Triangulo`. Crie uma função `MaiorArea(formas []Forma) Forma` que retorna a forma com maior área.

### 🟡 2. io.Reader customizado
Crie um tipo `RepeatReader` que implementa `io.Reader` e repete uma string N vezes. Use com `io.Copy` e `bufio.Scanner` para provar que funciona com a stdlib.

### 🟡 3. Mock para testes
Crie uma interface `NotificationSender` com método `Send(userID, message string) error`. Implemente um `MockSender` que registra as chamadas. Escreva testes que verificam que o service chama `Send` com os argumentos corretos.

### 🔴 4. Interface Segregation
Refatore uma interface `Storage` com 8 métodos em 3 interfaces menores (`Reader`, `Writer`, `Deleter`). Mostre que cada handler só precisa da interface mínima.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Structs](../modulo-07-structs/README.md) | [Próximo: Tratamento de Erros →](../modulo-09-tratamento-erros/README.md)
