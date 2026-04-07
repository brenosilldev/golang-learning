# Módulo 16 — Banco de Dados

[← gRPC](../modulo-15-grpc/README.md) | [Próximo: Docker →](../modulo-17-docker/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que `sql.Open` não conecta imediatamente ao banco?
> 2. O que acontece se você não fechar `rows` após um `Query`?
> 3. Por que transferências bancárias precisam de transactions?

---

## 1. Abordagens para Banco de Dados em Go

| Abordagem | Pacote | Prós | Contras | Quando usar |
|-----------|--------|------|---------|-------------|
| **database/sql** | stdlib | Zero deps, controle total, performático | Verboso, scan manual | APIs de alta performance, quando precisa de controle |
| **sqlx** | `jmoiron/sqlx` | Scan automático em structs, named queries | Ainda precisa escrever SQL | Melhor custo-benefício para a maioria dos projetos |
| **sqlc** | `sqlc.dev` | Type-safe, gera Go a partir de SQL | Setup inicial, migrations separadas | Projetos que valorizam type safety + SQL puro |
| **GORM** | `gorm.io/gorm` | ORM completo, migrations, hooks | Overhead, queries complexas difíceis | Prototipação rápida, CRUDs simples |

> **No mercado**: `database/sql` + `sqlx` é a combinação mais comum. `sqlc` está crescendo rápido em projetos novos. GORM é popular para MVPs mas controverso em produção pesada.

---

## 2. database/sql — A Base de Tudo

### Conectando (e o pool que você não vê)

```go
import (
    "database/sql"
    _ "github.com/lib/pq" // driver PostgreSQL (o _ registra o driver via init())
)

func main() {
    // sql.Open NÃO conecta — apenas valida o driver e cria o pool
    db, err := sql.Open("postgres",
        "host=localhost port=5432 user=app password=secret dbname=mydb sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // Ping verifica a conexão de verdade
    if err := db.Ping(); err != nil {
        log.Fatal("banco inacessível:", err)
    }
}
```

### Connection Pool — Configuração Essencial

`database/sql` gerencia um **pool de conexões** automaticamente. Você DEVE configurá-lo:

```go
db.SetMaxOpenConns(25)                 // máximo de conexões abertas simultâneas
db.SetMaxIdleConns(5)                  // conexões ociosas mantidas prontas
db.SetConnMaxLifetime(5 * time.Minute) // tempo máximo de vida de uma conexão
db.SetConnMaxIdleTime(1 * time.Minute) // tempo máximo ocioso
```

**Por que isso importa?**

| Configuração | Muito baixo | Muito alto |
|--------------|-------------|------------|
| `MaxOpenConns` | Lentidão (goroutines esperando) | Banco sobrecarregado |
| `MaxIdleConns` | Reconexões frequentes (lento) | Conexões ociosas consumindo memória |
| `ConnMaxLifetime` | Reconexões constantes | Conexões stale (banco pode derrubar) |

```
Regra prática para começar:
  MaxOpenConns = 25 (ajuste pelo número de CPUs do banco)
  MaxIdleConns = 5-10
  ConnMaxLifetime = 5 min
  ConnMaxIdleTime = 1 min
```

---

## 3. Context com Banco — Nunca Esqueça

Em produção, **toda operação de banco deve receber um context**. Sem context, uma query lenta trava a goroutine para sempre.

```go
// ❌ ERRADO — sem context, query pode travar indefinidamente
rows, err := db.Query("SELECT * FROM users WHERE name LIKE '%slow%'")

// ✅ CORRETO — context com timeout vindo da request HTTP
func handleListar(w http.ResponseWriter, r *http.Request) {
    // r.Context() já tem o timeout da request
    rows, err := db.QueryContext(r.Context(), "SELECT id, nome FROM users")
    if err != nil {
        // Se o client desconectar, o context cancela e a query para
        writeError(w, 500, "erro na consulta")
        return
    }
    defer rows.Close()
    // ...
}

// Para operações fora de handlers HTTP, crie um timeout explícito
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
_, err := db.ExecContext(ctx, "DELETE FROM logs WHERE created_at < $1", cutoff)
```

**Versões com Context (sempre prefira essas):**

| Sem Context | Com Context (use esta) |
|------------|----------------------|
| `db.Query(...)` | `db.QueryContext(ctx, ...)` |
| `db.QueryRow(...)` | `db.QueryRowContext(ctx, ...)` |
| `db.Exec(...)` | `db.ExecContext(ctx, ...)` |
| `db.Begin()` | `db.BeginTx(ctx, nil)` |
| `stmt.Query(...)` | `stmt.QueryContext(ctx, ...)` |

---

## 4. CRUD Completo com database/sql

```go
// Sempre verifique rows.Err() após o loop — erros de iteração não aparecem no Next()
rows, err := db.QueryContext(ctx, "SELECT id, nome, email FROM users")
if err != nil {
    return fmt.Errorf("listar users: %w", err)
}
defer rows.Close() // OBRIGATÓRIO — sem isso, conexão vaza do pool

var users []User
for rows.Next() {
    var u User
    if err := rows.Scan(&u.ID, &u.Nome, &u.Email); err != nil {
        return fmt.Errorf("scan user: %w", err)
    }
    users = append(users, u)
}
if err := rows.Err(); err != nil { // ← muita gente esquece isso
    return fmt.Errorf("iteração rows: %w", err)
}
```

### QueryRow para registro único

```go
var user User
err := db.QueryRowContext(ctx,
    "SELECT id, nome, email FROM users WHERE id = $1", id,
).Scan(&user.ID, &user.Nome, &user.Email)

switch {
case errors.Is(err, sql.ErrNoRows):
    return nil, ErrNaoEncontrado // sentinela do seu domínio
case err != nil:
    return nil, fmt.Errorf("buscar user %d: %w", id, err)
}
```

---

## 5. Transactions — Tudo ou Nada

```go
func Transferir(ctx context.Context, db *sql.DB, de, para int, valor float64) error {
    // BeginTx respeita o context
    tx, err := db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("iniciar tx: %w", err)
    }
    // defer Rollback é seguro — se Commit já foi chamado, Rollback é no-op
    defer tx.Rollback()

    // 1. Verificar saldo
    var saldo float64
    err = tx.QueryRowContext(ctx,
        "SELECT saldo FROM contas WHERE id = $1 FOR UPDATE", de,
    ).Scan(&saldo)
    if err != nil {
        return fmt.Errorf("consultar saldo: %w", err)
    }
    if saldo < valor {
        return fmt.Errorf("saldo insuficiente: %.2f < %.2f", saldo, valor)
    }

    // 2. Debitar
    _, err = tx.ExecContext(ctx,
        "UPDATE contas SET saldo = saldo - $1 WHERE id = $2", valor, de)
    if err != nil {
        return fmt.Errorf("debitar: %w", err)
    }

    // 3. Creditar
    _, err = tx.ExecContext(ctx,
        "UPDATE contas SET saldo = saldo + $1 WHERE id = $2", valor, para)
    if err != nil {
        return fmt.Errorf("creditar: %w", err)
    }

    // 4. Commit (se falhar, o defer Rollback desfaz tudo)
    if err := tx.Commit(); err != nil {
        return fmt.Errorf("commit: %w", err)
    }
    return nil
}
```

> `FOR UPDATE` no SELECT trava a linha — impede que outra transação leia o saldo enquanto esta transferência está em andamento (previne race condition).

---

## 6. Repository Pattern — Desacoplando o Banco

O Repository isola o acesso a dados atrás de uma interface. Isso permite trocar o banco (ou usar memória em testes) sem mudar a lógica de negócio.

```go
// Interface — o contrato
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    FindByID(ctx context.Context, id int) (*User, error)
    List(ctx context.Context) ([]User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
}

// Implementação com Postgres
type PostgresUserRepo struct {
    db *sql.DB
}

func NewPostgresUserRepo(db *sql.DB) *PostgresUserRepo {
    return &PostgresUserRepo{db: db}
}

func (r *PostgresUserRepo) FindByID(ctx context.Context, id int) (*User, error) {
    var u User
    err := r.db.QueryRowContext(ctx,
        "SELECT id, nome, email FROM users WHERE id = $1", id,
    ).Scan(&u.ID, &u.Nome, &u.Email)
    if errors.Is(err, sql.ErrNoRows) {
        return nil, ErrNaoEncontrado
    }
    if err != nil {
        return nil, fmt.Errorf("PostgresUserRepo.FindByID(%d): %w", id, err)
    }
    return &u, nil
}

// Implementação em memória (para testes)
type MemoryUserRepo struct {
    mu    sync.RWMutex
    users map[int]*User
    seq   int
}

func (r *MemoryUserRepo) FindByID(_ context.Context, id int) (*User, error) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    u, ok := r.users[id]
    if !ok {
        return nil, ErrNaoEncontrado
    }
    return u, nil
}

// O service recebe a INTERFACE — não sabe qual banco está por trás
type UserService struct {
    repo UserRepository // aceita Postgres, MySQL, memória — qualquer um
}

func NewUserService(repo UserRepository) *UserService {
    return &UserService{repo: repo}
}
```

**Benefícios:**
- Testes unitários sem banco real (usa `MemoryUserRepo`)
- Trocar Postgres por DynamoDB sem mudar `UserService`
- Cada implementação pode ser testada isoladamente

---

## 7. sqlc — Type-Safe SQL (Alternativa Moderna)

`sqlc` gera código Go a partir de queries SQL. Você escreve SQL puro e ganha type safety de graça:

```sql
-- queries.sql
-- name: GetUser :one
SELECT id, nome, email FROM users WHERE id = $1;

-- name: ListUsers :many
SELECT id, nome, email FROM users ORDER BY nome;

-- name: CreateUser :one
INSERT INTO users (nome, email) VALUES ($1, $2) RETURNING *;
```

```bash
# Gera código Go automaticamente
sqlc generate
```

```go
// Código gerado — tipado, com context, com error handling correto
user, err := queries.GetUser(ctx, 42)
users, err := queries.ListUsers(ctx)
newUser, err := queries.CreateUser(ctx, sqlc.CreateUserParams{
    Nome:  "Alice",
    Email: "alice@example.com",
})
```

> **Vantagem do sqlc**: erros de SQL são pegos em **tempo de compilação** (o sqlc valida contra o schema). Nenhum ORM faz isso.

---

## 8. Migrations — Evolução do Schema

```bash
# Instalar golang-migrate
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# Criar migration
migrate create -ext sql -dir migrations -seq create_users

# Isso cria 2 arquivos:
#   migrations/000001_create_users.up.sql    ← aplica a mudança
#   migrations/000001_create_users.down.sql  ← reverte a mudança
```

```sql
-- 000001_create_users.up.sql
CREATE TABLE users (
    id         SERIAL PRIMARY KEY,
    nome       VARCHAR(100) NOT NULL,
    email      VARCHAR(255) UNIQUE NOT NULL,
    criado_em  TIMESTAMP DEFAULT NOW()
);

-- 000001_create_users.down.sql
DROP TABLE IF EXISTS users;
```

```bash
# Aplicar
migrate -path migrations -database "postgres://app:secret@localhost/mydb?sslmode=disable" up

# Reverter última
migrate -path migrations -database "..." down 1
```

---

## ✅ Checklist de Banco de Dados para Produção

- [ ] **Connection pool** configurado (`MaxOpenConns`, `MaxIdleConns`, `ConnMaxLifetime`)
- [ ] **Context** propagado em toda operação (`QueryContext`, `ExecContext`, `BeginTx`)
- [ ] **rows.Close()** sempre chamado com `defer` imediatamente após `Query`
- [ ] **rows.Err()** verificado após o loop `Next()`
- [ ] **sql.ErrNoRows** tratado corretamente (não é um erro fatal)
- [ ] **Transactions** usadas para operações que precisam ser atômicas
- [ ] **defer tx.Rollback()** logo após `BeginTx` (seguro mesmo com Commit)
- [ ] **Prepared statements** ou parametrização para prevenir SQL injection
- [ ] **Migrations** versionadas (nunca ALTER TABLE diretamente em produção)
- [ ] **Repository pattern** para desacoplar banco da lógica de negócio

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo16_database_sql.go` | database/sql com SQLite |
| `exemplos/exemplo16_gorm.go` | GORM com SQLite |
| `exercicios/ex16_database.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. CRUD com database/sql
Crie uma aplicação de contatos usando database/sql + SQLite:
- Funções: `Criar`, `BuscarPorID`, `Listar`, `Atualizar`, `Deletar`
- Use `QueryContext`/`ExecContext` com context em **todas** as operações
- Trate `sql.ErrNoRows` corretamente

### 🟡 2. Transferência com Transaction
Crie um sistema bancário com transações:
- Verifique saldo suficiente, debite, credite, registre a transferência
- Tudo dentro de uma TRANSACTION com `FOR UPDATE`
- Se qualquer passo falhar, Rollback automático

### 🟡 3. Repository Pattern
Crie `UserRepository` (interface) com duas implementações:
- `PostgresUserRepo` (ou SQLite)
- `MemoryUserRepo` (para testes)
- Escreva testes usando a implementação em memória

### 🔴 4. Connection Pool Stress Test
Crie 100 goroutines que fazem queries simultâneas. Meça a latência com diferentes valores de `MaxOpenConns` (5, 25, 50). Observe como o pool afeta a performance.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← gRPC](../modulo-15-grpc/README.md) | [Próximo: Docker →](../modulo-17-docker/README.md)
