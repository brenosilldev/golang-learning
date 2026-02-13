package main

import (
	"database/sql"
	"fmt"
	"log"
	"time"

	_ "github.com/mattn/go-sqlite3" // Driver SQLite — o _ importa apenas o init()
)

// ============================================================================
// MÓDULO 16 — Banco de Dados com database/sql
// ============================================================================
//
// database/sql é a interface PADRÃO do Go para bancos de dados.
// Ela funciona com qualquer banco: SQLite, PostgreSQL, MySQL, etc.
// Você só troca o DRIVER.
//
// Instale o driver SQLite:
//   go get github.com/mattn/go-sqlite3
//
// Rode com: go run exemplo16_database_sql.go
// ============================================================================

// --- Model ---
type Usuario struct {
	ID       int
	Nome     string
	Email    string
	Ativo    bool
	CriadoEm time.Time
}

func main() {
	// ==========================================================
	fmt.Println("=== CONECTAR AO BANCO ===")

	// sql.Open NÃO conecta — apenas valida o driver
	db, err := sql.Open("sqlite3", "./exemplo.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	// Ping verifica a conexão de verdade
	if err := db.Ping(); err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Conectado ao SQLite!")

	// Configurar pool de conexões
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)
	db.SetConnMaxLifetime(5 * time.Minute)

	// ==========================================================
	fmt.Println("\n=== CRIAR TABELA ===")

	_, err = db.Exec(`
		CREATE TABLE IF NOT EXISTS usuarios (
			id        INTEGER PRIMARY KEY AUTOINCREMENT,
			nome      TEXT NOT NULL,
			email     TEXT UNIQUE NOT NULL,
			ativo     BOOLEAN DEFAULT true,
			criado_em DATETIME DEFAULT CURRENT_TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("  Tabela criada!")

	// Limpar dados anteriores (para reexecução)
	db.Exec("DELETE FROM usuarios")

	// ==========================================================
	fmt.Println("\n=== INSERT ===")

	// Exec retorna Result (LastInsertId, RowsAffected)
	result, err := db.Exec(
		"INSERT INTO usuarios (nome, email) VALUES (?, ?)",
		"Alice", "alice@example.com",
	)
	if err != nil {
		log.Fatal(err)
	}

	id, _ := result.LastInsertId()
	fmt.Printf("  Inserido: ID=%d\n", id)

	// Inserir mais
	nomes := []struct{ nome, email string }{
		{"Bob", "bob@example.com"},
		{"Carol", "carol@example.com"},
		{"Dave", "dave@example.com"},
	}
	for _, n := range nomes {
		db.Exec("INSERT INTO usuarios (nome, email) VALUES (?, ?)", n.nome, n.email)
	}
	fmt.Println("  4 usuários inseridos!")

	// ==========================================================
	fmt.Println("\n=== SELECT (QueryRow — um registro) ===")

	var usuario Usuario
	err = db.QueryRow(
		"SELECT id, nome, email, ativo FROM usuarios WHERE id = ?", 1,
	).Scan(&usuario.ID, &usuario.Nome, &usuario.Email, &usuario.Ativo)

	if err == sql.ErrNoRows {
		fmt.Println("  Não encontrado")
	} else if err != nil {
		log.Fatal(err)
	} else {
		fmt.Printf("  Encontrado: %+v\n", usuario)
	}

	// ==========================================================
	fmt.Println("\n=== SELECT (Query — múltiplos registros) ===")

	rows, err := db.Query("SELECT id, nome, email, ativo FROM usuarios")
	if err != nil {
		log.Fatal(err)
	}
	defer rows.Close() // SEMPRE fechar rows!

	fmt.Println("  Todos os usuários:")
	for rows.Next() {
		var u Usuario
		if err := rows.Scan(&u.ID, &u.Nome, &u.Email, &u.Ativo); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("    [%d] %s <%s> ativo=%t\n", u.ID, u.Nome, u.Email, u.Ativo)
	}

	// Verificar erros do iterator
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}

	// ==========================================================
	fmt.Println("\n=== UPDATE ===")

	result, err = db.Exec(
		"UPDATE usuarios SET ativo = ? WHERE nome = ?",
		false, "Bob",
	)
	if err != nil {
		log.Fatal(err)
	}
	affected, _ := result.RowsAffected()
	fmt.Printf("  Atualizado: %d registro(s)\n", affected)

	// ==========================================================
	fmt.Println("\n=== PREPARED STATEMENTS ===")

	// Prepared = compilado uma vez, executado N vezes
	// Protege contra SQL injection
	stmt, err := db.Prepare("SELECT nome, email FROM usuarios WHERE ativo = ?")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	rows2, _ := stmt.Query(true)
	defer rows2.Close()

	fmt.Println("  Usuários ativos:")
	for rows2.Next() {
		var nome, email string
		rows2.Scan(&nome, &email)
		fmt.Printf("    %s <%s>\n", nome, email)
	}

	// ==========================================================
	fmt.Println("\n=== TRANSACTIONS ===")

	tx, err := db.Begin()
	if err != nil {
		log.Fatal(err)
	}

	// Dentro da transação: ou TUDO funciona, ou NADA
	_, err1 := tx.Exec("UPDATE usuarios SET nome = 'Alice Updated' WHERE id = 1")
	_, err2 := tx.Exec("UPDATE usuarios SET nome = 'Carol Updated' WHERE id = 3")

	if err1 != nil || err2 != nil {
		tx.Rollback() // Desfaz tudo
		fmt.Println("  Transaction: ROLLBACK")
	} else {
		tx.Commit() // Confirma tudo
		fmt.Println("  Transaction: COMMIT ✓")
	}

	// Verificar
	var nome string
	db.QueryRow("SELECT nome FROM usuarios WHERE id = 1").Scan(&nome)
	fmt.Printf("  Alice agora é: %s\n", nome)

	// ==========================================================
	fmt.Println("\n=== DELETE ===")

	result, _ = db.Exec("DELETE FROM usuarios WHERE id = ?", 4)
	affected, _ = result.RowsAffected()
	fmt.Printf("  Deletado: %d registro(s)\n", affected)

	// Contagem final
	var count int
	db.QueryRow("SELECT COUNT(*) FROM usuarios").Scan(&count)
	fmt.Printf("  Total de usuários: %d\n", count)
}
