# Módulo 10 — Pacotes e Módulos

[← Tratamento de Erros](../modulo-09-tratamento-erros/README.md) | [Próximo: Concorrência →](../modulo-11-concorrencia/README.md)

---

> **Antes de ler — tente responder:**
> 1. Qual a diferença entre um `package` e um `module` em Go?
> 2. Por que o Go tem a pasta `internal/` com proteção no compilador?
> 3. O que torna uma função/variável "exportada" em Go?

---

## 1. Package vs Module — A Distinção Fundamental

```
MODULE = um projeto inteiro (tem go.mod)
         github.com/alice/myapp

PACKAGE = um diretório com arquivos .go
          github.com/alice/myapp/internal/handler
          github.com/alice/myapp/pkg/cache

Regra:
- Um módulo contém múltiplos pacotes
- Um pacote = um diretório
- Todos os .go no mesmo diretório compartilham o mesmo package name
```

---

## 2. go.mod — O Coração do Projeto

```bash
go mod init github.com/alice/myapp  # cria o go.mod
```

```
module github.com/alice/myapp

go 1.22

require (
    github.com/gin-gonic/gin v1.9.1
    golang.org/x/sync v0.6.0
)

require (
    // dependências indiretas — gerenciadas automaticamente
    github.com/bytedance/sonic v1.11.2 // indirect
    ...
)
```

### Comandos essenciais

```bash
go mod init github.com/user/projeto   # inicia módulo
go mod tidy                           # adiciona deps faltando, remove não usadas
go get github.com/gin-gonic/gin       # adiciona dependência
go get github.com/gin-gonic/gin@v1.9.1 # versão específica
go get github.com/gin-gonic/gin@latest # última versão
go list -m all                        # lista todas as dependências
go mod download                       # baixa deps para cache
```

---

## 3. Visibilidade — Exportado vs Privado

Em Go, a visibilidade é controlada pela **capitalização** do nome. Não há keywords `public`/`private`:

```go
package usuario

// Exportado (maiúscula) — acessível de fora do pacote
type User struct {
    ID    int    // exportado
    Email string // exportado
    senha string // privado — só acessível dentro do pacote usuario
}

func NewUser(email, senha string) *User {  // exportado
    return &User{Email: email, senha: hashSenha(senha)}
}

func (u *User) ValidarSenha(s string) bool {  // exportado
    return verificar(s, u.senha)  // acessa campo privado
}

func hashSenha(s string) string { /* ... */ }  // privado
func verificar(s, h string) bool { /* ... */ } // privado
```

### Por que campos privados importam

```go
// ❌ Tudo público — qualquer código pode colocar o objeto em estado inválido
type Conta struct {
    Saldo float64
}
c.Saldo = -1000  // inválido, mas possível

// ✅ Campo privado — invariante garantida pelo construtor e métodos
type Conta struct {
    saldo float64  // privado
}

func (c *Conta) Sacar(valor float64) error {
    if valor > c.saldo {
        return errors.New("saldo insuficiente")
    }
    c.saldo -= valor
    return nil
}
// Não há como colocar saldo negativo por fora
```

---

## 4. Estrutura de Projeto — O Padrão da Indústria

```
github.com/alice/myapp/
├── cmd/                    # executáveis (main packages)
│   ├── api/
│   │   └── main.go         # go run ./cmd/api
│   └── worker/
│       └── main.go         # go run ./cmd/worker
│
├── internal/               # código PRIVADO ao módulo
│   ├── handler/            # handlers HTTP
│   │   ├── user.go
│   │   └── user_test.go
│   ├── service/            # lógica de negócio
│   │   └── user.go
│   └── repository/         # acesso a dados
│       ├── user.go
│       └── memory.go       # implementação in-memory para testes
│
├── pkg/                    # código PÚBLICO (pode ser importado externamente)
│   └── validate/           # validadores genéricos
│       └── validate.go
│
├── go.mod
└── go.sum
```

### internal/ — Proteção no Nível do Compilador

```go
// ❌ ERRO DE COMPILAÇÃO se código externo tentar importar internal/
// github.com/outro/projeto não pode importar github.com/alice/myapp/internal/handler

// ✅ Apenas código dentro de github.com/alice/myapp pode importar internal/
// cmd/api/main.go pode importar internal/handler
// internal/handler pode importar internal/service
```

> **Por que isso importa?** `internal/` é como você cria uma API pública do seu módulo. Tudo que você não quer que terceiros usem vai em `internal/`. É uma promessa: "não vou garantir compatibilidade para isso."

---

## 5. Imports — Organização e Alias

```go
package main

import (
    // 1. stdlib — sem caminho de domínio
    "fmt"
    "net/http"
    "strings"

    // (linha em branco separa grupos — goimports faz isso automaticamente)

    // 2. dependências externas
    "github.com/gin-gonic/gin"
    "golang.org/x/sync/errgroup"

    // 3. pacotes internos do projeto
    "github.com/alice/myapp/internal/handler"
    "github.com/alice/myapp/pkg/validate"
)

// Alias para evitar conflito de nomes
import (
    crand "crypto/rand"   // alias: usa como crand.Read(...)
    mrand "math/rand"     // alias: usa como mrand.Intn(...)
)

// Blank import — importa só pelos side effects (init())
import _ "github.com/lib/pq" // registra driver PostgreSQL
```

---

## 6. init() — Inicialização do Pacote

```go
package db

var pool *sql.DB

// init() é chamado automaticamente quando o pacote é importado
// Ordem: variáveis → init() → main()
func init() {
    var err error
    pool, err = sql.Open("postgres", os.Getenv("DATABASE_URL"))
    if err != nil {
        log.Fatal("não conseguiu conectar ao banco:", err)
    }
}

// ⚠️ Use init() com moderação:
// - Dificulta testes (estado global implícito)
// - Ordem de init entre pacotes pode ser confusa
// ✅ Prefira: construtor explícito chamado no main
```

---

## ✅ Checklist de Pacotes e Módulos

- [ ] Módulo inicializado com `go mod init` (caminho de módulo com domínio)
- [ ] `go mod tidy` rodado após adicionar/remover dependências
- [ ] Código de uso interno em `internal/` (nunca em `pkg/` por padrão)
- [ ] Imports organizados em grupos: stdlib, externas, internas
- [ ] Nomes de pacotes são substantivos curtos, sem sublinhado (`handler`, não `http_handler`)
- [ ] Campos e funções privadas protegem invariantes do tipo

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo10_pacotes.go` | Organização, imports, visibilidade, init |
| `exercicios/ex10_pacotes.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Estrutura de Projeto
Crie a estrutura de diretórios para um módulo `github.com/seunome/todo-api` com:
- `cmd/api/main.go` — ponto de entrada
- `internal/handler/todo.go` — handlers HTTP
- `internal/service/todo.go` — lógica de negócio
- `internal/repository/todo.go` — armazenamento
Verifique que `internal/` não pode ser importado de fora.

### 🟡 2. Pacote de Validação
Crie um pacote `pkg/validate` com funções exportadas: `Email(s string) bool`, `CPF(s string) bool`, `Required(s string) bool`. Escreva testes no arquivo `validate_test.go`. Importe e use no `cmd/`.

### 🟡 3. Injeção de Dependências Manual
Refatore uma função que acessa banco diretamente para receber uma interface como parâmetro. Use o padrão `internal/repository` com implementação real e implementação em memória. Demonstre que o service funciona com ambas.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Tratamento de Erros](../modulo-09-tratamento-erros/README.md) | [Próximo: Concorrência →](../modulo-11-concorrencia/README.md)
