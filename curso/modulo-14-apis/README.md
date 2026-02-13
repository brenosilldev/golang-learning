# Módulo 14 — Construindo APIs

[← Testes](../modulo-13-testes/README.md) | [Próximo: Projeto Final →](../projeto-final/README.md)

---

## 📖 Teoria

### 3 abordagens para APIs em Go

| Abordagem | Prós | Contras |
|-----------|------|---------|
| **`net/http` puro** | Zero dependências, controle total | Mais verboso, sem roteamento avançado |
| **Gin** | Mais popular, middleware rico, rápido | Dependência externa |
| **Fiber** | API similar ao Express.js, ultra-rápido | Baseado em fasthttp (não net/http) |

### Arquitetura de API

```
Request → Router → Middleware → Handler → Service → Repository → DB
                                  ↓
Response ← JSON ← Handler ← Service ← Repository ← DB
```

### Patterns importantes
- **Handler**: Recebe HTTP request, chama service, retorna response
- **Service**: Lógica de negócio (não sabe de HTTP)
- **Repository**: Acesso a dados (não sabe de negócio)
- **DTO**: Data Transfer Object (struct para request/response)
- **Middleware**: Logging, auth, CORS, rate limiting

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `api-pura/main.go` | API CRUD completa com `net/http` |
| `api-gin/main.go` | API CRUD completa com Gin |
| `api-fiber/main.go` | API CRUD completa com Fiber |
| `exercicios/ex14_api.go` | 🏋️ Exercícios |

### Como rodar cada API:
```bash
# API pura (porta 8080)
go run curso/modulo-14-apis/api-pura/main.go

# API Gin (porta 8081) — precisa: go get github.com/gin-gonic/gin
go run curso/modulo-14-apis/api-gin/main.go

# API Fiber (porta 8082) — precisa: go get github.com/gofiber/fiber/v2
go run curso/modulo-14-apis/api-fiber/main.go
```

### Testando com curl:
```bash
# Listar todos
curl http://localhost:8080/api/tarefas

# Criar
curl -X POST http://localhost:8080/api/tarefas \
  -H "Content-Type: application/json" \
  -d '{"titulo": "Aprender Go", "descricao": "Terminar o curso"}'

# Buscar por ID
curl http://localhost:8080/api/tarefas/1

# Atualizar
curl -X PUT http://localhost:8080/api/tarefas/1 \
  -H "Content-Type: application/json" \
  -d '{"titulo": "Go Avançado", "concluida": true}'

# Deletar
curl -X DELETE http://localhost:8080/api/tarefas/1
```

---

[← Testes](../modulo-13-testes/README.md) | [Próximo: Projeto Final →](../projeto-final/README.md)
