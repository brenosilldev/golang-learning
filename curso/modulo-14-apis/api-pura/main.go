package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"time"
)

// ============================================================================
// MÓDULO 14 — API REST com net/http puro
// ============================================================================
//
// Esta é uma API CRUD completa sem NENHUMA dependência externa.
// Tudo é feito com a stdlib do Go.
//
// Conceitos:
//   - http.HandleFunc e http.ListenAndServe
//   - Roteamento manual (por path e método)
//   - JSON encoding/decoding
//   - Middleware (logging, CORS)
//   - Armazenamento em memória (map + mutex)
//
// Rode com: go run main.go
// Teste com: curl http://localhost:8080/api/tarefas
// ============================================================================

// --- Models ---

type Tarefa struct {
	ID        int       `json:"id"`
	Titulo    string    `json:"titulo"`
	Descricao string    `json:"descricao,omitempty"`
	Concluida bool      `json:"concluida"`
	CriadaEm  time.Time `json:"criada_em"`
}

type CriarTarefaDTO struct {
	Titulo    string `json:"titulo"`
	Descricao string `json:"descricao"`
}

type AtualizarTarefaDTO struct {
	Titulo    *string `json:"titulo,omitempty"`
	Descricao *string `json:"descricao,omitempty"`
	Concluida *bool   `json:"concluida,omitempty"`
}

type RespostaErro struct {
	Erro   string `json:"erro"`
	Codigo int    `json:"codigo"`
}

// --- Store (repositório em memória) ---

type TarefaStore struct {
	mu        sync.RWMutex
	tarefas   map[int]Tarefa
	proximoID int
}

func NovaStore() *TarefaStore {
	return &TarefaStore{
		tarefas:   make(map[int]Tarefa),
		proximoID: 1,
	}
}

func (s *TarefaStore) Listar() []Tarefa {
	s.mu.RLock()
	defer s.mu.RUnlock()

	lista := make([]Tarefa, 0, len(s.tarefas))
	for _, t := range s.tarefas {
		lista = append(lista, t)
	}
	return lista
}

func (s *TarefaStore) BuscarPorID(id int) (Tarefa, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	t, ok := s.tarefas[id]
	return t, ok
}

func (s *TarefaStore) Criar(dto CriarTarefaDTO) Tarefa {
	s.mu.Lock()
	defer s.mu.Unlock()

	tarefa := Tarefa{
		ID:        s.proximoID,
		Titulo:    dto.Titulo,
		Descricao: dto.Descricao,
		Concluida: false,
		CriadaEm:  time.Now(),
	}
	s.tarefas[s.proximoID] = tarefa
	s.proximoID++
	return tarefa
}

func (s *TarefaStore) Atualizar(id int, dto AtualizarTarefaDTO) (Tarefa, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	tarefa, ok := s.tarefas[id]
	if !ok {
		return Tarefa{}, false
	}

	if dto.Titulo != nil {
		tarefa.Titulo = *dto.Titulo
	}
	if dto.Descricao != nil {
		tarefa.Descricao = *dto.Descricao
	}
	if dto.Concluida != nil {
		tarefa.Concluida = *dto.Concluida
	}

	s.tarefas[id] = tarefa
	return tarefa, true
}

func (s *TarefaStore) Deletar(id int) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tarefas[id]; !ok {
		return false
	}
	delete(s.tarefas, id)
	return true
}

// --- Middleware ---

func logMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		inicio := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(inicio))
	})
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

// --- Helpers ---

func jsonResponse(w http.ResponseWriter, status int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(data)
}

func erroResponse(w http.ResponseWriter, status int, mensagem string) {
	jsonResponse(w, status, RespostaErro{Erro: mensagem, Codigo: status})
}

func extrairID(path string) (int, error) {
	partes := strings.Split(strings.Trim(path, "/"), "/")
	if len(partes) < 3 {
		return 0, fmt.Errorf("ID não encontrado no path")
	}
	return strconv.Atoi(partes[2])
}

// --- Handlers ---

func (s *TarefaStore) handleTarefas(w http.ResponseWriter, r *http.Request) {
	// Roteamento manual por path
	path := r.URL.Path

	switch {
	// GET /api/tarefas ou POST /api/tarefas
	case path == "/api/tarefas" || path == "/api/tarefas/":
		switch r.Method {
		case http.MethodGet:
			s.handleListar(w, r)
		case http.MethodPost:
			s.handleCriar(w, r)
		default:
			erroResponse(w, http.StatusMethodNotAllowed, "método não permitido")
		}

	// GET/PUT/DELETE /api/tarefas/{id}
	case strings.HasPrefix(path, "/api/tarefas/"):
		switch r.Method {
		case http.MethodGet:
			s.handleBuscar(w, r)
		case http.MethodPut:
			s.handleAtualizar(w, r)
		case http.MethodDelete:
			s.handleDeletar(w, r)
		default:
			erroResponse(w, http.StatusMethodNotAllowed, "método não permitido")
		}

	default:
		erroResponse(w, http.StatusNotFound, "rota não encontrada")
	}
}

func (s *TarefaStore) handleListar(w http.ResponseWriter, r *http.Request) {
	tarefas := s.Listar()
	jsonResponse(w, http.StatusOK, tarefas)
}

func (s *TarefaStore) handleCriar(w http.ResponseWriter, r *http.Request) {
	var dto CriarTarefaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		erroResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	if dto.Titulo == "" {
		erroResponse(w, http.StatusBadRequest, "título é obrigatório")
		return
	}

	tarefa := s.Criar(dto)
	jsonResponse(w, http.StatusCreated, tarefa)
}

func (s *TarefaStore) handleBuscar(w http.ResponseWriter, r *http.Request) {
	id, err := extrairID(r.URL.Path)
	if err != nil {
		erroResponse(w, http.StatusBadRequest, "ID inválido")
		return
	}

	tarefa, ok := s.BuscarPorID(id)
	if !ok {
		erroResponse(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}

	jsonResponse(w, http.StatusOK, tarefa)
}

func (s *TarefaStore) handleAtualizar(w http.ResponseWriter, r *http.Request) {
	id, err := extrairID(r.URL.Path)
	if err != nil {
		erroResponse(w, http.StatusBadRequest, "ID inválido")
		return
	}

	var dto AtualizarTarefaDTO
	if err := json.NewDecoder(r.Body).Decode(&dto); err != nil {
		erroResponse(w, http.StatusBadRequest, "JSON inválido")
		return
	}

	tarefa, ok := s.Atualizar(id, dto)
	if !ok {
		erroResponse(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}

	jsonResponse(w, http.StatusOK, tarefa)
}

func (s *TarefaStore) handleDeletar(w http.ResponseWriter, r *http.Request) {
	id, err := extrairID(r.URL.Path)
	if err != nil {
		erroResponse(w, http.StatusBadRequest, "ID inválido")
		return
	}

	if !s.Deletar(id) {
		erroResponse(w, http.StatusNotFound, "tarefa não encontrada")
		return
	}

	jsonResponse(w, http.StatusOK, map[string]string{"mensagem": "tarefa deletada"})
}

func main() {
	store := NovaStore()

	// Dados iniciais
	store.Criar(CriarTarefaDTO{Titulo: "Aprender Go", Descricao: "Completar o curso"})
	store.Criar(CriarTarefaDTO{Titulo: "Construir API", Descricao: "REST com net/http"})

	// Roteamento
	mux := http.NewServeMux()
	mux.HandleFunc("/api/tarefas", store.handleTarefas)
	mux.HandleFunc("/api/tarefas/", store.handleTarefas)

	// Aplicar middleware
	handler := logMiddleware(corsMiddleware(mux))

	fmt.Println("🚀 API Pura rodando em http://localhost:8080")
	fmt.Println("📋 Endpoints:")
	fmt.Println("   GET    /api/tarefas      — listar todas")
	fmt.Println("   POST   /api/tarefas      — criar nova")
	fmt.Println("   GET    /api/tarefas/{id}  — buscar por ID")
	fmt.Println("   PUT    /api/tarefas/{id}  — atualizar")
	fmt.Println("   DELETE /api/tarefas/{id}  — deletar")

	log.Fatal(http.ListenAndServe(":8080", handler))
}
