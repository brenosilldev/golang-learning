package main

// ============================================================================
// EXERCÍCIO 14 — APIs
// ============================================================================
//
// Exercício 14.1 — API de Notas (net/http puro)
// Crie uma API REST para gerenciar notas/anotações:
//   - GET    /api/notas           — listar todas
//   - POST   /api/notas           — criar nova
//   - GET    /api/notas/{id}      — buscar por ID
//   - PUT    /api/notas/{id}      — atualizar
//   - DELETE /api/notas/{id}      — deletar
//   - GET    /api/notas/busca?q=  — buscar por texto
//
// Struct: Nota { ID, Titulo, Conteudo, Tags []string, CriadaEm, AtualizadaEm }
// Armazene em memória com map + Mutex.
//
// Exercício 14.2 — Middleware
// Adicione estes middlewares à sua API:
//   a) Logger: log de todas as requisições com método, path e duração
//   b) Auth: verificar header "X-API-Key" (use uma chave hardcoded)
//   c) RateLimit: máximo de 10 requisições por minuto por IP
//
// Exercício 14.3 — Refatorar com Gin
// Pegue a API do exercício 14.1 e refaça usando Gin.
// Compare a quantidade de código e a legibilidade.
// Use:
//   - gin.Default() com Logger e Recovery
//   - ShouldBindJSON para parsing
//   - Grupos de rotas
//   - Query params com c.Query()
//
// Exercício 14.4 — Testes de API
// Escreva testes para sua API usando net/http/httptest:
//   - Teste cada endpoint (GET, POST, PUT, DELETE)
//   - Teste cenários de erro (404, 400, 401)
//   - Use table-driven tests
//
// ============================================================================

func main() {
	// TODO: Exercício 14.1 — API de Notas (net/http puro)

	// TODO: Exercício 14.2 — Middleware

	// TODO: Exercício 14.3 — Refatorar com Gin

	// TODO: Exercício 14.4 — Testes de API
}
