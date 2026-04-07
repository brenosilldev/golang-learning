# 🛠️ Projeto Intermediário 1 — CLI de Gerenciamento de Tarefas

> **Quando fazer**: após completar os Módulos 01–10
> **Tempo estimado**: 3–5 dias
> **Objetivo**: consolidar fundamentos antes de entrar em concorrência e APIs

---

## Por que este projeto?

Você aprendeu variáveis, tipos, coleções, funções, ponteiros, structs, interfaces, erros e pacotes. Mas provavelmente ainda não **combinou tudo isso em um projeto real**. Este projeto força você a usar tudo junto.

Go é a linguagem de CLIs do mercado: `kubectl`, `docker`, `terraform`, `gh`, `hugo` — todos são ferramentas de linha de comando escritas em Go. Saber construir uma CLI bem estruturada é diferencial real.

---

## O que você vai construir

Um gerenciador de tarefas em linha de comando chamado **`gotas`** (Go Task Manager):

```bash
# Criar tarefa
gotas add "Estudar Go" --priority alta --tag "estudo"
# ✅ Tarefa #1 criada: "Estudar Go"

# Listar tarefas
gotas list
# ID  TITULO              PRIORIDADE  STATUS      CRIADA EM
# 1   Estudar Go          alta        pendente    07/04/2026
# 2   Fazer exercícios    média       pendente    07/04/2026

# Listar com filtro
gotas list --status pendente --tag estudo

# Completar tarefa
gotas done 1
# ✅ Tarefa #1 marcada como concluída

# Detalhe de uma tarefa
gotas show 1

# Deletar
gotas delete 1

# Estatísticas
gotas stats
# Total: 10 | Pendentes: 3 | Concluídas: 7 | Taxa: 70%

# Exportar para JSON
gotas export --format json > tarefas.json

# Importar
gotas import tarefas.json
```

---

## Estrutura do Projeto

```
gotas/
├── cmd/
│   └── gotas/
│       └── main.go          ← ponto de entrada
├── internal/
│   ├── task/
│   │   ├── task.go          ← struct Task + validações
│   │   └── task_test.go
│   ├── storage/
│   │   ├── storage.go       ← interface Storage
│   │   ├── json_store.go    ← implementação em arquivo JSON
│   │   └── storage_test.go
│   └── cli/
│       ├── root.go          ← comando raiz + setup
│       ├── add.go           ← comando add
│       ├── list.go          ← comando list
│       ├── done.go          ← comando done
│       ├── delete.go        ← comando delete
│       ├── stats.go         ← comando stats
│       └── export.go        ← comando export/import
├── go.mod
├── go.sum
└── Makefile
```

---

## Especificação Técnica

### Struct Task

```go
type Priority string
const (
    PriorityLow    Priority = "baixa"
    PriorityMedium Priority = "média"
    PriorityHigh   Priority = "alta"
)

type Status string
const (
    StatusPending   Status = "pendente"
    StatusDone      Status = "concluída"
    StatusCancelled Status = "cancelada"
)

type Task struct {
    ID          int       `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description,omitempty"`
    Priority    Priority  `json:"priority"`
    Status      Status    `json:"status"`
    Tags        []string  `json:"tags,omitempty"`
    CreatedAt   time.Time `json:"created_at"`
    DoneAt      *time.Time `json:"done_at,omitempty"` // ponteiro = opcional
}
```

### Interface Storage

```go
type Storage interface {
    Save(task *task.Task) error
    FindByID(id int) (*task.Task, error)
    List(filter Filter) ([]*task.Task, error)
    Update(task *task.Task) error
    Delete(id int) error
    NextID() (int, error)
}

type Filter struct {
    Status   *task.Status   // nil = sem filtro
    Priority *task.Priority
    Tag      string
}
```

### Persistência em JSON

Os dados ficam em `~/.gotas/tasks.json`. Use `os.UserHomeDir()` para encontrar o diretório home.

```go
func defaultStorePath() string {
    home, err := os.UserHomeDir()
    if err != nil {
        home = "."
    }
    return filepath.Join(home, ".gotas", "tasks.json")
}
```

---

## Requisitos por Nível

### 🟢 Versão Mínima (obrigatória)
- [ ] Comandos: `add`, `list`, `done`, `delete`
- [ ] Persistência em arquivo JSON (`~/.gotas/tasks.json`)
- [ ] Flags básicas: `--priority`, `--tag`
- [ ] Saída formatada em tabela (sem biblioteca externa)
- [ ] Tratamento de erros com mensagens úteis ao usuário
- [ ] Testes para `task.Task` (validação) e `Storage`

### 🟡 Versão Intermediária
- [ ] Comando `stats` com porcentagem de conclusão
- [ ] Filtros em `list`: `--status`, `--tag`, `--priority`
- [ ] Comando `show` com detalhes completos de uma tarefa
- [ ] Cores no terminal (`\033[32m` para verde, etc.) sem biblioteca
- [ ] Confirmação antes de deletar (`Tem certeza? [s/N]`)
- [ ] Ordenação em `list`: `--sort created,priority`

### 🔴 Versão Avançada
- [ ] Comando `export`/`import` com suporte a JSON e CSV
- [ ] Subcomandos: `gotas tag add 1 "urgente"` / `gotas tag remove 1 "urgente"`
- [ ] Configuração em `~/.gotas/config.json` (prioridade padrão, formato de data)
- [ ] `gotas completion bash` gera autocomplete para bash/zsh
- [ ] Testes com 80%+ de coverage incluindo integration tests com storage real

---

## Guia de Implementação

### Passo 1 — Estrutura e dependências

```bash
mkdir gotas && cd gotas
go mod init github.com/seuusuario/gotas
# Não use cobra para este projeto — implemente o parsing manualmente
# Objetivo: entender como CLIs funcionam por baixo dos panos
# Se quiser usar cobra na versão avançada, ok
```

### Passo 2 — Comece pelo domínio (task.go)

```go
// Implemente primeiro sem CLI — só structs e métodos
// Escreva os testes ANTES de implementar (TDD)

func TestNewTask(t *testing.T) {
    t.Run("título vazio deve retornar erro", func(t *testing.T) {
        _, err := NewTask("", PriorityMedium)
        if err == nil {
            t.Error("esperava erro para título vazio")
        }
    })
}
```

### Passo 3 — Storage antes de CLI

```go
// Implemente JSONStorage primeiro
// Teste com arquivo temporário: os.CreateTemp("", "gotas-*.json")
func TestJSONStorage(t *testing.T) {
    f, _ := os.CreateTemp("", "gotas-*.json")
    defer os.Remove(f.Name())
    
    store := NewJSONStorage(f.Name())
    // ... testa Save, FindByID, List, etc.
}
```

### Passo 4 — CLI por último

```go
// main.go — parse manual de os.Args
func main() {
    if len(os.Args) < 2 {
        printHelp()
        os.Exit(1)
    }
    
    store := storage.NewJSONStorage(storage.DefaultPath())
    
    switch os.Args[1] {
    case "add":
        cli.RunAdd(store, os.Args[2:])
    case "list":
        cli.RunList(store, os.Args[2:])
    case "done":
        cli.RunDone(store, os.Args[2:])
    case "delete":
        cli.RunDelete(store, os.Args[2:])
    default:
        fmt.Fprintf(os.Stderr, "comando desconhecido: %s\n", os.Args[1])
        os.Exit(1)
    }
}
```

---

## Conceitos dos Módulos 01–10 Usados

| Conceito | Onde aparece no projeto |
|----------|------------------------|
| **M02** Tipos/iota | `Priority` e `Status` como tipos com constantes |
| **M03** Controle de fluxo | Switch de comandos, filtros em list |
| **M04** Slices/Maps | Lista de tarefas, mapa de tags |
| **M05** Funções/defer | Funções de comando, defer para fechar arquivo |
| **M06** Ponteiros | `*Task`, `*time.Time` para campos opcionais |
| **M07** Structs | `Task`, `Filter`, `Config` |
| **M08** Interfaces | `Storage` com múltiplas implementações |
| **M09** Erros | Sentinelas `ErrNotFound`, wrapping com contexto |
| **M10** Pacotes | Estrutura `cmd/internal`, visibilidade |

---

## Critérios de Avaliação

Após terminar, responda:

- [ ] O código compila sem warnings de `go vet`?
- [ ] `go test -race ./...` passa limpo?
- [ ] Erros têm mensagens úteis para o usuário (não stack trace)?
- [ ] A interface `Storage` permite trocar JSON por SQLite sem mudar a CLI?
- [ ] Um usuário consegue usar o CLI sem ler o código?

---

## Dica de Portfolio

Após terminar, faça um README bonito em inglês com:
- GIF ou screenshot do CLI em uso (`asciinema record`)
- Instruções de instalação: `go install github.com/seuusuario/gotas@latest`
- Badges: Go version, license, coverage

> **Próximo passo**: ao terminar este projeto, você está pronto para o Módulo 11 (Concorrência).
