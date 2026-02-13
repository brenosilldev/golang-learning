package main

import (
	"fmt"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// MÓDULO 14 — API REST com Fiber
// ============================================================================
//
// Fiber é inspirado no Express.js — se você tem experiência com Node,
// vai se sentir em casa.
//
// Instale com: go get github.com/gofiber/fiber/v2
//
// ⚠️ Fiber usa fasthttp ao invés de net/http.
//     Isso significa: mais rápido, mas incompatível com middleware net/http.
//
// Rode com: go run main.go
// Teste com: curl http://localhost:8082/api/tarefas
// ============================================================================

// ⚠️ NOTA: Para rodar este arquivo, você precisa instalar o Fiber:
//   go get github.com/gofiber/fiber/v2
//
// Depois descomente as imports e o código abaixo.

// import "github.com/gofiber/fiber/v2"

type Tarefa struct {
	ID        int       `json:"id"`
	Titulo    string    `json:"titulo"`
	Descricao string    `json:"descricao,omitempty"`
	Concluida bool      `json:"concluida"`
	CriadaEm  time.Time `json:"criada_em"`
}

type TarefaStore struct {
	mu        sync.RWMutex
	tarefas   map[int]Tarefa
	proximoID int
}

func NovaStore() *TarefaStore {
	return &TarefaStore{tarefas: make(map[int]Tarefa), proximoID: 1}
}

// ============================================================================
// CÓDIGO FIBER (descomente após instalar)
// ============================================================================
//
// func main() {
// 	store := NovaStore()
// 	store.tarefas[1] = Tarefa{ID: 1, Titulo: "Aprender Go", CriadaEm: time.Now()}
// 	store.tarefas[2] = Tarefa{ID: 2, Titulo: "Usar Fiber", CriadaEm: time.Now()}
// 	store.proximoID = 3
//
// 	// Criar app Fiber
// 	app := fiber.New(fiber.Config{
// 		AppName: "Tarefas API v1.0",
// 	})
//
// 	// Grupo de rotas
// 	api := app.Group("/api")
//
// 	// GET /api/tarefas
// 	api.Get("/tarefas", func(c *fiber.Ctx) error {
// 		store.mu.RLock()
// 		defer store.mu.RUnlock()
//
// 		lista := make([]Tarefa, 0)
// 		for _, t := range store.tarefas {
// 			lista = append(lista, t)
// 		}
// 		return c.JSON(lista)
// 	})
//
// 	// POST /api/tarefas
// 	api.Post("/tarefas", func(c *fiber.Ctx) error {
// 		var body struct {
// 			Titulo    string `json:"titulo"`
// 			Descricao string `json:"descricao"`
// 		}
// 		if err := c.BodyParser(&body); err != nil {
// 			return c.Status(400).JSON(fiber.Map{"erro": "JSON inválido"})
// 		}
// 		if body.Titulo == "" {
// 			return c.Status(400).JSON(fiber.Map{"erro": "título obrigatório"})
// 		}
//
// 		store.mu.Lock()
// 		tarefa := Tarefa{
// 			ID: store.proximoID, Titulo: body.Titulo,
// 			Descricao: body.Descricao, CriadaEm: time.Now(),
// 		}
// 		store.tarefas[store.proximoID] = tarefa
// 		store.proximoID++
// 		store.mu.Unlock()
//
// 		return c.Status(201).JSON(tarefa)
// 	})
//
// 	// GET /api/tarefas/:id
// 	api.Get("/tarefas/:id", func(c *fiber.Ctx) error {
// 		id, err := c.ParamsInt("id")
// 		if err != nil {
// 			return c.Status(400).JSON(fiber.Map{"erro": "ID inválido"})
// 		}
//
// 		store.mu.RLock()
// 		tarefa, ok := store.tarefas[id]
// 		store.mu.RUnlock()
//
// 		if !ok {
// 			return c.Status(404).JSON(fiber.Map{"erro": "não encontrada"})
// 		}
// 		return c.JSON(tarefa)
// 	})
//
// 	// DELETE /api/tarefas/:id
// 	api.Delete("/tarefas/:id", func(c *fiber.Ctx) error {
// 		id, _ := c.ParamsInt("id")
//
// 		store.mu.Lock()
// 		if _, ok := store.tarefas[id]; !ok {
// 			store.mu.Unlock()
// 			return c.Status(404).JSON(fiber.Map{"erro": "não encontrada"})
// 		}
// 		delete(store.tarefas, id)
// 		store.mu.Unlock()
//
// 		return c.JSON(fiber.Map{"mensagem": "deletada"})
// 	})
//
// 	fmt.Println("🚀 API Fiber rodando em http://localhost:8082")
// 	log.Fatal(app.Listen(":8082"))
// }

func main() {
	_ = NovaStore()
	fmt.Println("============================================================")
	fmt.Println("  API FIBER — Código de referência")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Para rodar esta API com Fiber, siga os passos:")
	fmt.Println()
	fmt.Println("  1. Instale o Fiber:")
	fmt.Println("     go get github.com/gofiber/fiber/v2")
	fmt.Println()
	fmt.Println("  2. Descomente o código main() com Fiber neste arquivo")
	fmt.Println("  3. Comente este main() placeholder")
	fmt.Println("  4. Rode: go run main.go")
	fmt.Println()
	fmt.Println("Vantagens do Fiber:")
	fmt.Println("  ✅ API idêntica ao Express.js")
	fmt.Println("  ✅ Ultra-rápido (baseado em fasthttp)")
	fmt.Println("  ✅ c.BodyParser() para parse automático")
	fmt.Println("  ✅ c.ParamsInt() para path params tipados")
	fmt.Println("  ✅ fiber.Map{} para respostas rápidas")
	fmt.Println()
	fmt.Println("⚠️  Cuidado: Fiber usa fasthttp, não net/http.")
	fmt.Println("    Middleware padrão do Go NÃO funciona com Fiber.")

	_ = http.StatusOK
}
