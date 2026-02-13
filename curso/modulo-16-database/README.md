# Módulo 16 — Banco de Dados

[← gRPC](../modulo-15-grpc/README.md) | [Próximo: Docker →](../modulo-17-docker/README.md)

---

## 📖 Teoria

### 3 abordagens para DB em Go

| Abordagem | Pacote | Quando usar |
|-----------|--------|------------|
| **database/sql** (raw) | stdlib | Controle total, queries SQL puras |
| **sqlx** | `jmoiron/sqlx` | database/sql + scan em structs |
| **GORM** | `gorm.io/gorm` | ORM completo (Rails-like) |

### database/sql — A base de tudo
```go
db, err := sql.Open("sqlite3", "app.db")
defer db.Close()

// Exec (INSERT, UPDATE, DELETE)
db.Exec("INSERT INTO users (name) VALUES (?)", "Alice")

// Query (SELECT múltiplos)
rows, _ := db.Query("SELECT id, name FROM users")
defer rows.Close()
for rows.Next() {
    var id int; var name string
    rows.Scan(&id, &name)
}

// QueryRow (SELECT único)
var name string
db.QueryRow("SELECT name FROM users WHERE id = ?", 1).Scan(&name)
```

### Prepared Statements (evita SQL injection)
```go
stmt, _ := db.Prepare("SELECT * FROM users WHERE email = ?")
defer stmt.Close()
stmt.QueryRow("alice@test.com").Scan(&user)
```

### Transactions
```go
tx, _ := db.Begin()
tx.Exec("UPDATE accounts SET balance = balance - 100 WHERE id = 1")
tx.Exec("UPDATE accounts SET balance = balance + 100 WHERE id = 2")
if err != nil { tx.Rollback() } else { tx.Commit() }
```

### GORM — ORM completo
```go
type User struct {
    gorm.Model           // ID, CreatedAt, UpdatedAt, DeletedAt
    Name  string
    Email string `gorm:"uniqueIndex"`
}

db.AutoMigrate(&User{})
db.Create(&User{Name: "Alice"})
db.First(&user, 1)
db.Where("name = ?", "Alice").Find(&users)
db.Model(&user).Update("name", "Bob")
db.Delete(&user, 1)
```

### Migrations
```bash
# Com golang-migrate
go install github.com/golang-migrate/migrate/v4/cmd/migrate@latest
migrate create -ext sql -dir migrations -seq create_users
migrate -path migrations -database "sqlite3://app.db" up
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo16_database_sql.go` | database/sql com SQLite |
| `exemplos/exemplo16_gorm.go` | GORM com SQLite |
| `exercicios/ex16_database.go` | 🏋️ Exercícios |

---

[← gRPC](../modulo-15-grpc/README.md) | [Próximo: Docker →](../modulo-17-docker/README.md)
