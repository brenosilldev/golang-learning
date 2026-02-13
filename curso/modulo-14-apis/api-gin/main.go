package main

import (
	"fmt"
	"log"
	"net/http"
	"sync"
	"time"
)

// ============================================================================
// MÓDULO 14 — API REST com Gin
// ============================================================================
//
// Gin é o framework HTTP mais popular do Go.
// Instale com: go get github.com/gin-gonic/gin
//
// Conceitos:
//   - Router e grupos de rotas
//   - Binding automático (JSON → struct)
//   - Middleware built-in
//   - Path parameters (:id)
//   - Validação com tags
//
// Rode com: go run main.go
// Teste com: curl http://localhost:8081/api/tarefas
// ============================================================================

// ⚠️ NOTA: Para rodar este arquivo, você precisa instalar o Gin:
//   go get github.com/gin-gonic/gin
//
// Depois descomente as imports e o código abaixo.
// O código está comentado para não quebrar o build do projeto.

// import "github.com/gin-gonic/gin"

// --- Models (mesmos da API pura) ---

type Tarefa struct {
	ID        int       `json:"id"`
	Titulo    string    `json:"titulo" binding:"required"`
	Descricao string    `json:"descricao,omitempty"`
	Concluida bool      `json:"concluida"`
	CriadaEm  time.Time `json:"criada_em"`
}

type CriarTarefaDTO struct {
	Titulo    string `json:"titulo" binding:"required"`
	Descricao string `json:"descricao"`
}

type AtualizarTarefaDTO struct {
	Titulo    *string `json:"titulo"`
	Descricao *string `json:"descricao"`
	Concluida *bool   `json:"concluida"`
}

// --- Store ---

type TarefaStore struct {
	mu        sync.RWMutex
	tarefas   map[int]Tarefa
	proximoID int
}

func NovaStore() *TarefaStore {
	return &TarefaStore{tarefas: make(map[int]Tarefa), proximoID: 1}
}

// ============================================================================
// CÓDIGO GIN (descomente após instalar o gin)
// ============================================================================
//
// func main() {
// 	store := NovaStore()
//
// 	// Dados iniciais
// 	store.tarefas[1] = Tarefa{ID: 1, Titulo: "Aprender Go", CriadaEm: time.Now()}
// 	store.tarefas[2] = Tarefa{ID: 2, Titulo: "Usar Gin", CriadaEm: time.Now()}
// 	store.proximoID = 3
//
// 	// Criar router (gin.Default inclui Logger e Recovery middleware)
// 	r := gin.Default()
//
// 	// Grupo de rotas
// 	api := r.Group("/api")
// 	{
// 		api.GET("/tarefas", func(c *gin.Context) {
// 			store.mu.RLock()
// 			defer store.mu.RUnlock()
//
// 			lista := make([]Tarefa, 0)
// 			for _, t := range store.tarefas {
// 				lista = append(lista, t)
// 			}
// 			c.JSON(http.StatusOK, lista)
// 		})
//
// 		api.POST("/tarefas", func(c *gin.Context) {
// 			var dto CriarTarefaDTO
// 			// ShouldBindJSON faz parse E validação automática!
// 			if err := c.ShouldBindJSON(&dto); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
// 				return
// 			}
//
// 			store.mu.Lock()
// 			tarefa := Tarefa{
// 				ID:        store.proximoID,
// 				Titulo:    dto.Titulo,
// 				Descricao: dto.Descricao,
// 				CriadaEm:  time.Now(),
// 			}
// 			store.tarefas[store.proximoID] = tarefa
// 			store.proximoID++
// 			store.mu.Unlock()
//
// 			c.JSON(http.StatusCreated, tarefa)
// 		})
//
// 		// :id é path parameter — acessa com c.Param("id")
// 		api.GET("/tarefas/:id", func(c *gin.Context) {
// 			id, err := strconv.Atoi(c.Param("id"))
// 			if err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"erro": "ID inválido"})
// 				return
// 			}
//
// 			store.mu.RLock()
// 			tarefa, ok := store.tarefas[id]
// 			store.mu.RUnlock()
//
// 			if !ok {
// 				c.JSON(http.StatusNotFound, gin.H{"erro": "não encontrada"})
// 				return
// 			}
//
// 			c.JSON(http.StatusOK, tarefa)
// 		})
//
// 		api.PUT("/tarefas/:id", func(c *gin.Context) {
// 			id, _ := strconv.Atoi(c.Param("id"))
//
// 			var dto AtualizarTarefaDTO
// 			if err := c.ShouldBindJSON(&dto); err != nil {
// 				c.JSON(http.StatusBadRequest, gin.H{"erro": err.Error()})
// 				return
// 			}
//
// 			store.mu.Lock()
// 			tarefa, ok := store.tarefas[id]
// 			if !ok {
// 				store.mu.Unlock()
// 				c.JSON(http.StatusNotFound, gin.H{"erro": "não encontrada"})
// 				return
// 			}
//
// 			if dto.Titulo != nil { tarefa.Titulo = *dto.Titulo }
// 			if dto.Descricao != nil { tarefa.Descricao = *dto.Descricao }
// 			if dto.Concluida != nil { tarefa.Concluida = *dto.Concluida }
// 			store.tarefas[id] = tarefa
// 			store.mu.Unlock()
//
// 			c.JSON(http.StatusOK, tarefa)
// 		})
//
// 		api.DELETE("/tarefas/:id", func(c *gin.Context) {
// 			id, _ := strconv.Atoi(c.Param("id"))
//
// 			store.mu.Lock()
// 			if _, ok := store.tarefas[id]; !ok {
// 				store.mu.Unlock()
// 				c.JSON(http.StatusNotFound, gin.H{"erro": "não encontrada"})
// 				return
// 			}
// 			delete(store.tarefas, id)
// 			store.mu.Unlock()
//
// 			c.JSON(http.StatusOK, gin.H{"mensagem": "deletada"})
// 		})
// 	}
//
// 	fmt.Println("🚀 API Gin rodando em http://localhost:8081")
// 	r.Run(":8081")
// }

func main() {
	_ = NovaStore()
	fmt.Println("============================================================")
	fmt.Println("  API GIN — Código de referência")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Para rodar esta API com Gin, siga os passos:")
	fmt.Println()
	fmt.Println("  1. Instale o Gin:")
	fmt.Println("     go get github.com/gin-gonic/gin")
	fmt.Println()
	fmt.Println("  2. Descomente o código main() com Gin neste arquivo")
	fmt.Println("  3. Comente este main() placeholder")
	fmt.Println("  4. Rode: go run main.go")
	fmt.Println()
	fmt.Println("Vantagens do Gin vs net/http:")
	fmt.Println("  ✅ Roteamento com :params e grupos")
	fmt.Println("  ✅ Binding automático (JSON → struct)")
	fmt.Println("  ✅ Validação com tags (binding:\"required\")")
	fmt.Println("  ✅ Middleware built-in (Logger, Recovery)")
	fmt.Println("  ✅ gin.H{} para respostas rápidas")
	fmt.Println("  ✅ Context com helpers (c.JSON, c.Param, etc)")

	// Evitar unused imports
	_ = http.StatusOK
	_ = log.Println
}
