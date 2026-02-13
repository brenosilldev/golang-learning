package main

// ============================================================================
// EXERCÍCIO 16 — Banco de Dados
// ============================================================================
//
// Exercício 16.1 — CRUD com database/sql
// Crie uma aplicação de contatos usando database/sql + SQLite:
//   - Tabela: contatos (id, nome, telefone, email, grupo, criado_em)
//   - Funções: CriarContato, BuscarPorID, ListarTodos, Atualizar, Deletar
//   - Função: BuscarPorGrupo(grupo string) → []Contato
//   - Use Prepared Statements para TODAS as queries
//   - Trate sql.ErrNoRows corretamente
//
// Exercício 16.2 — Transactions
// Crie um sistema bancário simples:
//   - Tabela: contas (id, titular, saldo)
//   - Tabela: transferencias (id, de_conta, para_conta, valor, data)
//   - Função: Transferir(de, para, valor) que:
//     a) Verifica saldo suficiente
//     b) Debita de uma conta
//     c) Credita na outra
//     d) Registra a transferência
//     e) Tudo dentro de uma TRANSACTION
//   - Se qualquer passo falhar, faz ROLLBACK
//
// Exercício 16.3 — GORM (ORM)
// Refaça o exercício 16.1 usando GORM:
//   - Models com tags gorm
//   - AutoMigrate
//   - CRUD usando db.Create, db.First, db.Find, db.Update, db.Delete
//   - Adicione uma relação: Contato pertence a um Grupo (model separado)
//   - Use Preload para carregar o grupo junto
//
// Exercício 16.4 — Repository Pattern
// Crie uma interface Repository para desacoplar o banco:
//   type ContatoRepository interface {
//       Criar(c *Contato) error
//       BuscarPorID(id int) (*Contato, error)
//       Listar() ([]Contato, error)
//       Atualizar(c *Contato) error
//       Deletar(id int) error
//   }
// Implemente com SQLite (SqliteContatoRepo).
// Implemente com memória (MemoryContatoRepo) para testes.
// Escreva testes usando a implementação em memória.
//
// ============================================================================

func main() {
	// TODO: Exercício 16.1 — CRUD com database/sql

	// TODO: Exercício 16.2 — Transactions

	// TODO: Exercício 16.3 — GORM

	// TODO: Exercício 16.4 — Repository Pattern
}
