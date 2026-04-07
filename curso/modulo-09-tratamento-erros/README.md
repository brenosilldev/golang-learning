# Módulo 09 — Tratamento de Erros

[← Interfaces](../modulo-08-interfaces/README.md) | [Próximo: Pacotes e Módulos →](../modulo-10-pacotes-modulos/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que Go não tem exceções (try/catch)?
> 2. Qual a diferença entre `errors.Is` e `errors.As`?
> 3. Quando usar `panic` é aceitável?

---

## 1. O Pattern Fundamental do Go

```go
resultado, err := funcaoQuePodeFalhar()
if err != nil {
    return fmt.Errorf("contexto do que falhou: %w", err)
}
// usar resultado com segurança
```

Go **não tem exceções**. Erros são **valores** retornados — simples, explícito e sem mágica.

**Por que essa decisão de design?**
- Exceções criam fluxo de controle **invisível** (o código pula para um catch distante)
- Em Go, todo caminho de erro é **visível** no código — você vê exatamente onde cada falha é tratada
- Isso é mais verboso, mas previne bugs sutis em sistemas distribuídos

---

## 2. Criando Erros — 4 Formas

### Sentinela (constante global comparável)
```go
var (
    ErrNaoEncontrado = errors.New("recurso não encontrado")
    ErrSemPermissao  = errors.New("sem permissão")
)

// Uso: retornar quando a condição exata é conhecida
func buscar(id int) (*User, error) {
    if id < 0 {
        return nil, ErrNaoEncontrado
    }
    return &User{}, nil
}
```

### Com contexto (fmt.Errorf)
```go
func buscarUsuario(id int) (*User, error) {
    user, err := db.Get(id)
    if err != nil {
        // %w "embrulha" o erro original — preserva a cadeia
        return nil, fmt.Errorf("buscarUsuario(id=%d): %w", id, err)
    }
    return user, nil
}
```

### Custom error type (quando você precisa de dados extras)
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validação: campo '%s' — %s", e.Field, e.Message)
}

// Custom error com Unwrap (para errors.Is funcionar na cadeia)
type OperationError struct {
    Op  string
    Err error
}

func (e *OperationError) Error() string {
    return fmt.Sprintf("operação '%s': %v", e.Op, e.Err)
}

func (e *OperationError) Unwrap() error {
    return e.Err // permite errors.Is/As encontrar o erro original
}
```

### errors.Join (Go 1.20+) — Múltiplos erros simultâneos
```go
func validarFormulario(nome, email string) error {
    var errs []error

    if nome == "" {
        errs = append(errs, &ValidationError{Field: "nome", Message: "obrigatório"})
    }
    if !strings.Contains(email, "@") {
        errs = append(errs, &ValidationError{Field: "email", Message: "formato inválido"})
    }

    return errors.Join(errs...) // retorna nil se slice está vazio
}

// O erro retornado suporta errors.Is e errors.As para CADA erro individual
err := validarFormulario("", "invalido")
var valErr *ValidationError
if errors.As(err, &valErr) {
    fmt.Println("primeiro campo com erro:", valErr.Field)
}
```

> **`errors.Join` é essencial** quando você valida múltiplos campos e quer retornar **todos** os erros de uma vez (não só o primeiro).

---

## 3. errors.Is e errors.As — Inspecionando a Cadeia

Quando erros são embrulhados com `%w`, forma-se uma **cadeia**:

```
OperationError("buscar") → fmt.Errorf("repo: %w") → ErrNaoEncontrado
```

### errors.Is — "Este erro É (ou contém) X?"

```go
err := buscarUsuario(-1) // retorna cadeia embrulhada

// Percorre TODA a cadeia procurando ErrNaoEncontrado
if errors.Is(err, ErrNaoEncontrado) {
    // SIM — em algum nível da cadeia tem ErrNaoEncontrado
    fmt.Println("usuário não existe")
}
```

### errors.As — "Este erro contém um tipo X? Se sim, me dê ele."

```go
var valErr *ValidationError
if errors.As(err, &valErr) {
    // SIM — extraiu o ValidationError da cadeia
    fmt.Printf("campo com erro: %s\n", valErr.Field)
}
```

### Diferença prática

| Função | Pergunta que responde | Quando usar |
|--------|----------------------|-------------|
| `errors.Is(err, target)` | "O erro é **igual** a este valor?" | Comparar com sentinelas (`ErrNotFound`, `sql.ErrNoRows`) |
| `errors.As(err, &target)` | "O erro é **deste tipo**? Me dê o valor." | Extrair dados do erro (`ValidationError.Field`) |

---

## 4. Error Wrapping — Adicionando Contexto Sem Perder a Cadeia

```go
// ❌ ERRADO — erro original perdido, sem contexto
func handler(w http.ResponseWriter, r *http.Request) {
    _, err := service.Buscar(42)
    if err != nil {
        http.Error(w, "erro", 500) // QUEM chamou o quê? Qual ID? Impossível debugar
    }
}

// ✅ CORRETO — cada nível adiciona contexto com %w
func (repo *UserRepo) BuscarPorID(id int) (*User, error) {
    row := repo.db.QueryRow("SELECT ...", id)
    if err := row.Scan(&user); err != nil {
        return nil, fmt.Errorf("UserRepo.BuscarPorID(%d): %w", id, err)
    }
    return &user, nil
}

func (svc *UserService) Buscar(id int) (*User, error) {
    user, err := svc.repo.BuscarPorID(id)
    if err != nil {
        return nil, fmt.Errorf("UserService.Buscar: %w", err)
    }
    return user, nil
}

// Resultado: "UserService.Buscar: UserRepo.BuscarPorID(42): sql: no rows"
// Cada nível diz ONDE o erro aconteceu, e errors.Is(err, sql.ErrNoRows) funciona
```

### Regra de ouro: trate OU embrulhe, nunca ambos

```go
// ❌ ERRADO — loga E retorna (trata duas vezes)
if err != nil {
    log.Printf("erro: %v", err)  // quem chamou vai logar DE NOVO
    return err
}

// ✅ CORRETO (nível baixo) — embrulha e retorna
if err != nil {
    return fmt.Errorf("processar pedido %d: %w", id, err)
}

// ✅ CORRETO (nível mais alto, ex: handler HTTP) — trata definitivamente
if err != nil {
    log.Printf("erro: %v", err)
    writeError(w, http.StatusInternalServerError, "erro interno")
    // NÃO retorna o erro — já tratou
}
```

---

## 5. Erros em Goroutines — Armadilha Comum

Goroutines não retornam erros. Se uma goroutine falha silenciosamente, você perde a informação.

```go
// ❌ ERRADO — erro engolido pela goroutine
go func() {
    err := processar(item)
    if err != nil {
        log.Println(err) // loga, mas quem lançou a goroutine não sabe
    }
}()

// ✅ FIX 1: channel de erros
errCh := make(chan error, 1)
go func() {
    errCh <- processar(item)
}()
if err := <-errCh; err != nil {
    // agora o chamador sabe
}

// ✅ FIX 2: errgroup (melhor para N goroutines)
g, ctx := errgroup.WithContext(ctx)
g.Go(func() error {
    return processar(item) // erro propagado automaticamente
})
if err := g.Wait(); err != nil {
    // primeiro erro de qualquer goroutine
}
```

---

## 6. Árvore de Decisão — Qual Tipo de Erro Usar?

```
Você precisa retornar um erro?
│
├── É uma condição conhecida e fixa? (ex: "não encontrado", "já existe")
│   └── Use SENTINEL ERROR (var ErrNotFound = errors.New(...))
│
├── O chamador precisa de dados extras do erro? (ex: qual campo falhou)
│   └── Use CUSTOM ERROR TYPE (type ValidationError struct{...})
│
├── Precisa retornar MÚLTIPLOS erros de uma vez? (ex: validação de formulário)
│   └── Use errors.Join (Go 1.20+)
│
├── Quer adicionar contexto a um erro existente?
│   └── Use fmt.Errorf("contexto: %w", err)
│
└── É um erro simples sem necessidade especial?
    └── Use errors.New("descrição")
```

---

## 7. panic e recover — Quando Usar (e Quando NÃO)

```go
// ✅ QUANDO USAR panic:
// 1. Bug impossível no código (assertion falhou)
// 2. Inicialização falhou (DB não conectou, config inválida)
// 3. Invariante violada (ex: enum com valor inesperado)

func MustParse(s string) Config {
    cfg, err := Parse(s)
    if err != nil {
        panic(fmt.Sprintf("config inválida: %v", err)) // falha na inicialização
    }
    return cfg
}

// ❌ QUANDO NÃO USAR panic:
// - Arquivo não encontrado → retorne erro
// - Input do usuário inválido → retorne erro
// - Serviço externo indisponível → retorne erro
// - Qualquer situação "esperada" em produção → retorne erro

// recover: captura panic em defer (usado em middleware HTTP para não derrubar o servidor)
func RecoverMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if err := recover(); err != nil {
                log.Printf("PANIC: %v\n%s", err, debug.Stack())
                http.Error(w, "erro interno", 500)
            }
        }()
        next.ServeHTTP(w, r)
    })
}
```

> **Convenção Go**: funções que fazem panic no nome começam com `Must` (ex: `regexp.MustCompile`, `template.Must`). Isso sinaliza que a função pode panic se o input for inválido.

---

## ✅ Checklist de Tratamento de Erros para Produção

- [ ] Todo `if err != nil` **faz algo útil** (não ignora o erro)
- [ ] Erros são **embrulhados com contexto** (`fmt.Errorf("...: %w", err)`)
- [ ] Erros são tratados **uma vez só** (não loga E retorna)
- [ ] Sentinel errors são usados para condições **conhecidas e fixas**
- [ ] Custom error types são usados quando o chamador precisa de **dados extras**
- [ ] Goroutines **propagam erros** (via channel ou errgroup)
- [ ] `panic` só é usado para **bugs impossíveis** ou falha de inicialização
- [ ] Handler HTTP tem `recover` para panics não derrubarem o servidor
- [ ] `errors.Is` e `errors.As` são usados para inspecionar erros (não `err.Error() == "..."`)

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo09_erros.go` | Todos os patterns de erro |
| `exercicios/ex09_erros.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Sentinel Errors
Crie um pacote `storage` com operações CRUD que retornam sentinel errors (`ErrNotFound`, `ErrAlreadyExists`, `ErrInvalidInput`). No `main`, use `errors.Is` para tratar cada caso.

### 🟡 2. Custom Error com Contexto
Crie um `ValidationError` com campos `Field`, `Value` e `Message`. Implemente uma função `ValidarUsuario(nome, email, idade)` que retorna `errors.Join` com todos os erros de validação. No `main`, use `errors.As` para extrair cada campo com problema.

### 🟡 3. Error Wrapping em Camadas
Simule 3 camadas (repo → service → handler) onde cada uma embrulha o erro com contexto. Verifique que `errors.Is` ainda encontra o erro sentinela original através de 3 níveis de wrapping.

### 🔴 4. Error Handling em Goroutines
Crie 5 goroutines que processam itens. Algumas falham aleatoriamente. Use `errgroup` para coletar o primeiro erro e cancelar as demais. Compare com a abordagem de channel de erros.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Interfaces](../modulo-08-interfaces/README.md) | [Próximo: Pacotes e Módulos →](../modulo-10-pacotes-modulos/README.md)
