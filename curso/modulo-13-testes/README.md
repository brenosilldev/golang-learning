# Módulo 13 — Testes

[← Generics](../modulo-12-generics/README.md) | [Próximo: APIs →](../modulo-14-apis/README.md)

---

> **Antes de ler — tente responder:**
> 1. O que é um "table-driven test" e por que é o padrão em Go?
> 2. Qual a diferença entre `t.Error` e `t.Fatal`?
> 3. Como você testa código que depende de um banco de dados?

---

## 1. Convenções de Teste em Go

```
arquivo.go          → arquivo_test.go       (mesmo diretório)
func Soma(...)      → func TestSoma(...)    (prefixo Test)
func BenchmarkSoma  → benchmarks
func ExampleSoma    → exemplos executáveis na documentação
```

```bash
go test ./...                    # todos os testes
go test ./internal/service/...   # pacote específico
go test -v ./...                 # verbose — mostra cada teste
go test -run TestSoma ./...      # só testes que combinam com "TestSoma"
go test -race ./...              # detecta race conditions
go test -cover ./...             # coverage
go test -bench=. -benchmem ./... # benchmarks com memória
```

---

## 2. Estrutura Básica — testing.T

```go
package soma_test  // pacote separado — testa a API pública (black-box)
// ou: package soma — acessa internals (white-box)

import (
    "testing"
    "github.com/alice/myapp/internal/soma"
)

func TestSomaSimples(t *testing.T) {
    resultado := soma.Somar(2, 3)
    if resultado != 5 {
        t.Errorf("Somar(2, 3) = %d; esperado 5", resultado)
    }
}

// t.Error  — registra falha, continua o teste
// t.Errorf — registra falha com formato, continua
// t.Fatal  — registra falha, PARA o teste imediatamente
// t.Fatalf — registra falha com formato, PARA
// t.Log    — log (só aparece com -v ou se teste falhar)
// t.Helper() — marca função como helper (stacktrace correto)
```

---

## 3. Table-Driven Tests — O Pattern Padrão do Go

Table-driven tests são o jeito idiomático de cobrir múltiplos cenários sem duplicar código:

```go
func TestDividir(t *testing.T) {
    tests := []struct {
        name      string
        a, b      float64
        want      float64
        wantErr   bool
    }{
        {
            name: "divisão normal",
            a: 10, b: 2,
            want: 5,
        },
        {
            name: "divisão por zero",
            a: 10, b: 0,
            wantErr: true,
        },
        {
            name: "números negativos",
            a: -10, b: 2,
            want: -5,
        },
        {
            name: "resultado decimal",
            a: 1, b: 3,
            want: 0.3333333333333333,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {  // ← subteste com nome
            got, err := Dividir(tt.a, tt.b)

            if (err != nil) != tt.wantErr {
                t.Fatalf("Dividir() error = %v, wantErr %v", err, tt.wantErr)
            }
            if !tt.wantErr && got != tt.want {
                t.Errorf("Dividir() = %v, want %v", got, tt.want)
            }
        })
    }
}

// Rodar apenas um caso: go test -run TestDividir/divisão_por_zero
```

**Vantagens:**
- Um novo caso = uma nova linha na tabela
- Nomes descritivos tornam falhas fáceis de identificar
- `t.Run` cria sub-testes — podem ser rodados individualmente

---

## 4. Helpers e Setup/Teardown

```go
// t.Helper() marca a função como helper
// O stacktrace aponta para quem chamou o helper, não para dentro dele
func assertEqual(t *testing.T, got, want interface{}) {
    t.Helper()  // ← essencial para helper functions
    if got != want {
        t.Errorf("got %v, want %v", got, want)
    }
}

// Setup e teardown com t.Cleanup
func TestComSetup(t *testing.T) {
    // Setup
    db := criarBancoDeTeste(t)
    t.Cleanup(func() {
        db.Close()  // executado quando o teste terminar, mesmo se falhar
    })

    // Testa usando db...
}

// TestMain — setup/teardown global para o pacote
func TestMain(m *testing.M) {
    // Roda antes de todos os testes
    setup()

    // Roda os testes
    code := m.Run()

    // Roda depois de todos os testes
    teardown()

    os.Exit(code)
}
```

---

## 5. Testando com Dependências — Mocks e Fakes

```go
// Interface (definida no código de produção)
type UserRepository interface {
    FindByID(ctx context.Context, id int) (*User, error)
    Save(ctx context.Context, user *User) error
}

// Fake (implementação in-memory para testes — simples e rápida)
type FakeUserRepo struct {
    users map[int]*User
    mu    sync.RWMutex
}

func NewFakeUserRepo() *FakeUserRepo {
    return &FakeUserRepo{users: make(map[int]*User)}
}

func (f *FakeUserRepo) FindByID(_ context.Context, id int) (*User, error) {
    f.mu.RLock()
    defer f.mu.RUnlock()
    if u, ok := f.users[id]; ok {
        return u, nil
    }
    return nil, ErrNotFound
}

func (f *FakeUserRepo) Save(_ context.Context, user *User) error {
    f.mu.Lock()
    defer f.mu.Unlock()
    f.users[user.ID] = user
    return nil
}

// Mock (verifica que foi chamado corretamente)
type MockUserRepo struct {
    FindByIDCalled bool
    FindByIDArg    int
    FindByIDReturn *User
    FindByIDErr    error
}

func (m *MockUserRepo) FindByID(_ context.Context, id int) (*User, error) {
    m.FindByIDCalled = true
    m.FindByIDArg = id
    return m.FindByIDReturn, m.FindByIDErr
}

// Teste usando o mock
func TestUserService_Buscar(t *testing.T) {
    repo := &MockUserRepo{
        FindByIDReturn: &User{ID: 1, Nome: "Alice"},
    }
    svc := NewUserService(repo)

    user, err := svc.Buscar(context.Background(), 1)

    if err != nil {
        t.Fatalf("erro inesperado: %v", err)
    }
    if !repo.FindByIDCalled {
        t.Error("FindByID não foi chamado")
    }
    if repo.FindByIDArg != 1 {
        t.Errorf("FindByID chamado com %d, esperava 1", repo.FindByIDArg)
    }
    if user.Nome != "Alice" {
        t.Errorf("user.Nome = %s, esperava Alice", user.Nome)
    }
}
```

---

## 6. Benchmarks

```go
func BenchmarkSoma(b *testing.B) {
    // b.N é ajustado automaticamente pelo Go para obter resultados estáveis
    for i := 0; i < b.N; i++ {
        Soma(1, 2)
    }
}

// Benchmark com setup (não conta no tempo)
func BenchmarkProcessar(b *testing.B) {
    dados := gerarDadosTeste(10000) // setup não entra no benchmark
    b.ResetTimer()                   // reseta o timer após o setup

    for i := 0; i < b.N; i++ {
        Processar(dados)
    }
}

// Rodar: go test -bench=. -benchmem
// Resultado:
// BenchmarkSoma-8    1000000000    0.3 ns/op
// BenchmarkProcessar-8  50000    24521 ns/op    8192 B/op    2 allocs/op
//                                                ↑ bytes alocados por op
```

---

## 7. Coverage — Quanto do Código é Testado

```bash
# Mostrar % de coverage
go test -cover ./...

# Gerar relatório visual (HTML)
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out

# Coverage por função
go tool cover -func=coverage.out
```

> **Meta prática**: 80% de coverage é um bom alvo. 100% pode ser desnecessário — foco em testar a lógica de negócio crítica, não getters triviais.

---

## ✅ Checklist de Testes para Produção

- [ ] Funções de produção são cobertas por **table-driven tests** com casos de sucesso E erro
- [ ] Dependências externas são injetadas via interfaces (testável com fakes/mocks)
- [ ] `t.Helper()` em toda função auxiliar de teste
- [ ] `t.Cleanup()` para limpeza de recursos (não `defer` solto)
- [ ] `go test -race ./...` passa sem data races
- [ ] Coverage acima de 70-80% nas partes críticas
- [ ] Benchmarks escritos para código sensível a performance

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo13_codigo.go` | Código que será testado |
| `exemplos/exemplo13_codigo_test.go` | Testes completos demonstrando todos os patterns |
| `exercicios/ex13_testes_test.go` | 🏋️ Exercícios de teste |

---

## 📋 Exercícios

### 🟢 1. Table-Driven Tests
Escreva table-driven tests para uma função `Validar(email string) bool`. Cubra: email válido, sem @, domínio inválido, string vazia, espaços. Use `t.Run` com nomes descritivos.

### 🟡 2. Service com Mock de Repository
Crie `UserService` que depende de `UserRepository` (interface). Escreva testes usando um `MockUserRepository` que verifica: quantas vezes foi chamado, com quais argumentos, e permite configurar o retorno.

### 🟡 3. Benchmark de Strings
Compare a performance de: concatenação com `+` em loop vs `strings.Builder` vs `fmt.Sprintf`. Use `b.ReportAllocs()`. Qual é a diferença em alocações?

### 🔴 4. Teste de Integração com Banco
Configure uma suite de testes que:
- Usa SQLite em memória (`:memory:`) para testes de integração
- Cria o schema no `TestMain`
- Cada teste começa com banco limpo usando `t.Cleanup`
- Testa o `UserRepository` real (não mock) contra banco real

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Generics](../modulo-12-generics/README.md) | [Próximo: APIs →](../modulo-14-apis/README.md)
