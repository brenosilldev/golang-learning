package main

import (
	"fmt"
)

// ============================================================================
// MÓDULO 16 — Banco de Dados com GORM
// ============================================================================
//
// GORM é o ORM mais popular do Go (como ActiveRecord/Eloquent/Prisma).
// Ele abstrai SQL — você trabalha com structs Go.
//
// Instale:
//   go get gorm.io/gorm
//   go get gorm.io/driver/sqlite
//
// Para usar PostgreSQL em vez de SQLite:
//   go get gorm.io/driver/postgres
//   dsn := "host=localhost user=app password=secret dbname=mydb port=5432"
//   db, _ := gorm.Open(postgres.Open(dsn), &gorm.Config{})
//
// Descomente o código abaixo e rode: go run exemplo16_gorm.go
// ============================================================================

// import (
// 	"log"
// 	"time"
//
// 	"gorm.io/driver/sqlite"
// 	"gorm.io/gorm"
// )

// --- Models (GORM usa tags para configurar o schema) ---
//
// type Categoria struct {
// 	ID   uint   `gorm:"primaryKey"`
// 	Nome string `gorm:"uniqueIndex;not null"`
// }
//
// type Produto struct {
// 	ID           uint      `gorm:"primaryKey"`
// 	Nome         string    `gorm:"not null;size:100"`
// 	Preco        float64   `gorm:"not null"`
// 	Estoque      int       `gorm:"default:0"`
// 	CategoriaID  uint      `gorm:"not null"`            // FK
// 	Categoria    Categoria `gorm:"foreignKey:CategoriaID"` // Relationship
// 	CriadoEm     time.Time `gorm:"autoCreateTime"`
// 	AtualizadoEm time.Time `gorm:"autoUpdateTime"`
// }
//
// func main() {
// 	// Conectar
// 	db, err := gorm.Open(sqlite.Open("gorm_exemplo.db"), &gorm.Config{})
// 	if err != nil {
// 		log.Fatal(err)
// 	}
//
// 	// AutoMigrate cria/atualiza tabelas baseado nas structs
// 	db.AutoMigrate(&Categoria{}, &Produto{})
//
// 	// Limpar
// 	db.Exec("DELETE FROM produtos")
// 	db.Exec("DELETE FROM categorias")
//
// 	// ==========================================================
// 	fmt.Println("=== CREATE ===")
//
// 	// Criar categorias
// 	cat1 := Categoria{Nome: "Eletrônicos"}
// 	cat2 := Categoria{Nome: "Livros"}
// 	db.Create(&cat1)
// 	db.Create(&cat2)
// 	fmt.Printf("  Categoria criada: ID=%d, Nome=%s\n", cat1.ID, cat1.Nome)
//
// 	// Criar produtos
// 	produtos := []Produto{
// 		{Nome: "MacBook Pro", Preco: 14999.99, Estoque: 10, CategoriaID: cat1.ID},
// 		{Nome: "iPhone 16", Preco: 7999.99, Estoque: 25, CategoriaID: cat1.ID},
// 		{Nome: "O Senhor dos Anéis", Preco: 49.90, Estoque: 100, CategoriaID: cat2.ID},
// 		{Nome: "Clean Code", Preco: 89.90, Estoque: 50, CategoriaID: cat2.ID},
// 	}
// 	db.Create(&produtos) // Batch insert!
// 	fmt.Printf("  %d produtos criados!\n", len(produtos))
//
// 	// ==========================================================
// 	fmt.Println("\n=== READ (Find, First, Where) ===")
//
// 	// Buscar por ID
// 	var prod Produto
// 	db.First(&prod, 1) // SELECT * FROM produtos WHERE id = 1
// 	fmt.Printf("  Primeiro: %s — R$%.2f\n", prod.Nome, prod.Preco)
//
// 	// Buscar com Where
// 	var caros []Produto
// 	db.Where("preco > ?", 1000).Find(&caros)
// 	fmt.Printf("  Produtos > R$1000: %d\n", len(caros))
//
// 	// Com Preload (carrega relationship)
// 	var prodComCat Produto
// 	db.Preload("Categoria").First(&prodComCat, 1)
// 	fmt.Printf("  %s — Categoria: %s\n", prodComCat.Nome, prodComCat.Categoria.Nome)
//
// 	// Listar todos com categoria
// 	var todos []Produto
// 	db.Preload("Categoria").Find(&todos)
// 	for _, p := range todos {
// 		fmt.Printf("    [%d] %s (R$%.2f) — %s\n", p.ID, p.Nome, p.Preco, p.Categoria.Nome)
// 	}
//
// 	// ==========================================================
// 	fmt.Println("\n=== UPDATE ===")
//
// 	// Atualizar um campo
// 	db.Model(&Produto{}).Where("id = ?", 1).Update("preco", 13999.99)
//
// 	// Atualizar múltiplos campos
// 	db.Model(&Produto{}).Where("id = ?", 2).Updates(map[string]interface{}{
// 		"preco":   6999.99,
// 		"estoque": 30,
// 	})
// 	fmt.Println("  Preços atualizados!")
//
// 	// ==========================================================
// 	fmt.Println("\n=== QUERIES AVANÇADAS ===")
//
// 	// Count
// 	var total int64
// 	db.Model(&Produto{}).Count(&total)
// 	fmt.Printf("  Total de produtos: %d\n", total)
//
// 	// Select específico + Order
// 	var nomes []struct{ Nome string; Preco float64 }
// 	db.Model(&Produto{}).Select("nome, preco").Order("preco DESC").Scan(&nomes)
// 	for _, n := range nomes {
// 		fmt.Printf("    %s — R$%.2f\n", n.Nome, n.Preco)
// 	}
//
// 	// Group By
// 	var stats []struct{ CategoriaID uint; Total int64; Media float64 }
// 	db.Model(&Produto{}).
// 		Select("categoria_id, COUNT(*) as total, AVG(preco) as media").
// 		Group("categoria_id").
// 		Scan(&stats)
// 	for _, s := range stats {
// 		fmt.Printf("    Cat %d: %d produtos, média R$%.2f\n", s.CategoriaID, s.Total, s.Media)
// 	}
//
// 	// ==========================================================
// 	fmt.Println("\n=== TRANSACTION ===")
//
// 	err = db.Transaction(func(tx *gorm.DB) error {
// 		if err := tx.Model(&Produto{}).Where("id = ?", 1).
// 			Update("estoque", gorm.Expr("estoque - ?", 1)).Error; err != nil {
// 			return err // Rollback automático
// 		}
// 		// Mais operações...
// 		return nil // Commit automático
// 	})
// 	if err != nil {
// 		fmt.Println("  Transaction falhou:", err)
// 	} else {
// 		fmt.Println("  Transaction OK ✓")
// 	}
//
// 	// ==========================================================
// 	fmt.Println("\n=== DELETE ===")
//
// 	// Soft delete (se model tiver gorm.Model com DeletedAt)
// 	// Hard delete
// 	db.Unscoped().Delete(&Produto{}, 4)
// 	fmt.Println("  Produto 4 deletado!")
//
// 	// Contagem final
// 	db.Model(&Produto{}).Count(&total)
// 	fmt.Printf("  Produtos restantes: %d\n", total)
// }

func main() {
	fmt.Println("============================================================")
	fmt.Println("  GORM — Código de referência")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Para rodar:")
	fmt.Println("  1. go get gorm.io/gorm")
	fmt.Println("  2. go get gorm.io/driver/sqlite")
	fmt.Println("  3. Descomente o código neste arquivo")
	fmt.Println("  4. go run exemplo16_gorm.go")
	fmt.Println()
	fmt.Println("GORM vs database/sql:")
	fmt.Println("  ✅ AutoMigrate — cria tabelas a partir de structs")
	fmt.Println("  ✅ Preload    — carrega relationships automaticamente")
	fmt.Println("  ✅ Chainable  — db.Where(...).Order(...).Find(...)")
	fmt.Println("  ✅ Hooks      — BeforeCreate, AfterUpdate, etc")
	fmt.Println("  ⚠️  Mais lento — overhead do ORM")
	fmt.Println("  ⚠️  Queries complexas — às vezes é melhor SQL puro")
}
