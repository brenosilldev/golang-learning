# 🔐 Módulo 23 — Segurança em Go

> **Nível**: Avançado | **Pré-requisito**: Módulos 14 (APIs), 17 (Docker)

---

> **Antes de ler — tente responder:**
> 1. O que é a diferença entre autenticação e autorização?
> 2. Por que não devemos armazenar senhas em texto plano?
> 3. O que é um JWT e por que ele é "stateless"?

---

## 🤔 Por que Segurança é Obrigatório?

Segurança não é opcional em 2026. **90% das vagas Go** pedem conhecimento de:
- Autenticação com JWT ou OAuth2
- TLS mútuo (mTLS) para comunicação entre microserviços
- Middleware de autorização (RBAC)
- Proteção contra vulnerabilidades comuns (OWASP Top 10)

---

## 1. Hashing de Senhas com bcrypt

**Nunca** armazene senhas em texto plano ou com MD5/SHA.

```go
package main

import (
    "errors"
    "fmt"
    "log"

    "golang.org/x/crypto/bcrypt"
)

var (
    ErrInvalidPassword = errors.New("senha inválida")
    ErrWeakPassword    = errors.New("senha muito fraca (mínimo 8 caracteres)")
)

func HashPassword(password string) (string, error) {
    if len(password) < 8 {
        return "", ErrWeakPassword
    }
    // bcrypt inclui salt automaticamente — nunca use o mesmo hash para a mesma senha
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", fmt.Errorf("erro ao gerar hash: %w", err)
    }
    return string(hash), nil
}

func CheckPassword(password, hash string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
        return ErrInvalidPassword
    }
    return err
}

func main() {
    hash, err := HashPassword("minhasenha123")
    if err != nil {
        log.Fatal(err)
    }
    fmt.Println("Hash:", hash)

    // Verifica senha correta
    if err := CheckPassword("minhasenha123", hash); err != nil {
        fmt.Println("Erro:", err)
    } else {
        fmt.Println("Senha correta ✅")
    }

    // Verifica senha errada
    if err := CheckPassword("senhaerrada", hash); err != nil {
        fmt.Println("Senha errada:", err) // ErrInvalidPassword
    }
}
```

---

## 2. JWT — JSON Web Tokens

JWT é um token **stateless**: o servidor não precisa de banco para validar — toda a informação está no token.

```
Header.Payload.Signature

eyJhbGciOiJIUzI1NiJ9.eyJ1c2VySWQiOiIxMjMiLCJleHAiOjE3MDAwMDAwMDB9.abc123
```

```go
package main

import (
    "errors"
    "fmt"
    "time"

    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("seu-segredo-aqui-mínimo-32-bytes!")

type Claims struct {
    UserID string `json:"user_id"`
    Email  string `json:"email"`
    Role   string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(userID, email, role string) (string, error) {
    claims := &Claims{
        UserID: userID,
        Email:  email,
        Role:   role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
            Issuer:    "api-pedidos",
            Subject:   userID,
        },
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

var (
    ErrTokenExpired  = errors.New("token expirado")
    ErrTokenInvalid  = errors.New("token inválido")
)

func ValidateToken(tokenString string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenString, &Claims{},
        func(token *jwt.Token) (interface{}, error) {
            // Valida o algoritmo — essencial para evitar "alg:none" attack
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                return nil, fmt.Errorf("algoritmo inesperado: %v", token.Header["alg"])
            }
            return jwtSecret, nil
        },
    )
    if err != nil {
        if errors.Is(err, jwt.ErrTokenExpired) {
            return nil, ErrTokenExpired
        }
        return nil, ErrTokenInvalid
    }

    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid {
        return nil, ErrTokenInvalid
    }
    return claims, nil
}

func main() {
    // Gera token
    token, err := GenerateToken("user-123", "joao@example.com", "admin")
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    fmt.Println("Token:", token[:50]+"...")

    // Valida token
    claims, err := ValidateToken(token)
    if err != nil {
        fmt.Println("Erro:", err)
        return
    }
    fmt.Printf("UserID: %s | Role: %s | Expira: %v\n",
        claims.UserID, claims.Role, claims.ExpiresAt.Time.Format(time.RFC3339))
}
```

---

## 3. Middleware de Autenticação & Autorização (RBAC)

```go
package main

import (
    "context"
    "net/http"
    "strings"
)

type contextKey string
const claimsKey contextKey = "claims"

// AuthMiddleware valida o JWT e adiciona claims ao contexto
func AuthMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        authHeader := r.Header.Get("Authorization")
        if authHeader == "" {
            http.Error(w, `{"error":"não autenticado"}`, http.StatusUnauthorized)
            return
        }

        // Formato: "Bearer <token>"
        parts := strings.SplitN(authHeader, " ", 2)
        if len(parts) != 2 || parts[0] != "Bearer" {
            http.Error(w, `{"error":"formato inválido"}`, http.StatusUnauthorized)
            return
        }

        claims, err := ValidateToken(parts[1])
        if err != nil {
            status := http.StatusUnauthorized
            if errors.Is(err, ErrTokenExpired) {
                http.Error(w, `{"error":"token expirado"}`, status)
            } else {
                http.Error(w, `{"error":"token inválido"}`, status)
            }
            return
        }

        // Adiciona claims ao contexto para uso nos handlers
        ctx := context.WithValue(r.Context(), claimsKey, claims)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// RequireRole retorna 403 se o usuário não tiver o role necessário
func RequireRole(role string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            claims, ok := r.Context().Value(claimsKey).(*Claims)
            if !ok {
                http.Error(w, `{"error":"claims não encontradas"}`, http.StatusInternalServerError)
                return
            }

            if claims.Role != role && claims.Role != "admin" {
                http.Error(w, `{"error":"acesso negado"}`, http.StatusForbidden)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// GetClaims extrai claims do contexto (helper para handlers)
func GetClaims(r *http.Request) *Claims {
    claims, _ := r.Context().Value(claimsKey).(*Claims)
    return claims
}

// Exemplo de uso com mux
func setupRouter() http.Handler {
    mux := http.NewServeMux()

    // Rota pública
    mux.HandleFunc("/api/login", handleLogin)

    // Rota autenticada (qualquer usuário logado)
    mux.Handle("/api/pedidos", AuthMiddleware(
        http.HandlerFunc(handlePedidos),
    ))

    // Rota com role específico
    mux.Handle("/api/admin/users", AuthMiddleware(
        RequireRole("admin")(
            http.HandlerFunc(handleAdminUsers),
        ),
    ))

    return mux
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
    // Validaria usuário/senha no banco, retornaria JWT
    token, _ := GenerateToken("user-1", "user@example.com", "user")
    w.Header().Set("Content-Type", "application/json")
    fmt.Fprintf(w, `{"token":"%s"}`, token)
}

func handlePedidos(w http.ResponseWriter, r *http.Request) {
    claims := GetClaims(r)
    fmt.Fprintf(w, `{"pedidos":[],"user":"%s"}`, claims.UserID)
}

func handleAdminUsers(w http.ResponseWriter, r *http.Request) {
    fmt.Fprint(w, `{"users":[]}`)
}
```

---

## 4. TLS e mTLS entre microserviços

```go
package main

import (
    "crypto/tls"
    "crypto/x509"
    "fmt"
    "net/http"
    "os"
)

// Servidor com TLS (HTTPS)
func serverWithTLS() *http.Server {
    return &http.Server{
        Addr: ":8443",
        TLSConfig: &tls.Config{
            MinVersion: tls.VersionTLS13, // nunca TLS 1.0 ou 1.1
            CipherSuites: []uint16{
                tls.TLS_AES_128_GCM_SHA256,
                tls.TLS_AES_256_GCM_SHA384,
            },
        },
    }
    // Inicia com: server.ListenAndServeTLS("cert.pem", "key.pem")
}

// Cliente com mTLS (verifica certificado do servidor E apresenta o próprio)
func clientWithMTLS() *http.Client {
    // Carrega CA que assinou o certificado do servidor
    caCert, err := os.ReadFile("ca.crt")
    if err != nil {
        panic(err)
    }
    caCertPool := x509.NewCertPool()
    caCertPool.AppendCertsFromPEM(caCert)

    // Carrega certificado e chave do cliente
    clientCert, err := tls.LoadX509KeyPair("client.crt", "client.key")
    if err != nil {
        panic(err)
    }

    tlsConfig := &tls.Config{
        RootCAs:      caCertPool,   // verifica servidor
        Certificates: []tls.Certificate{clientCert}, // apresenta ao servidor
        MinVersion:   tls.VersionTLS13,
    }

    return &http.Client{
        Transport: &http.Transport{TLSClientConfig: tlsConfig},
    }
}

// Gera certificados para desenvolvimento (não use em produção!)
// Use: openssl genrsa -out key.pem 4096
//      openssl req -new -x509 -sha256 -key key.pem -out cert.pem -days 365
```

---

## 5. Proteções OWASP Essenciais

```go
// Rate Limiting — previne DDoS e brute force
import "golang.org/x/time/rate"

type RateLimiter struct {
    mu       sync.Mutex
    limiters map[string]*rate.Limiter
    rate     rate.Limit
    burst    int
}

func NewRateLimiter(r rate.Limit, burst int) *RateLimiter {
    return &RateLimiter{
        limiters: make(map[string]*rate.Limiter),
        rate:     r,
        burst:    burst,
    }
}

func (rl *RateLimiter) Get(ip string) *rate.Limiter {
    rl.mu.Lock()
    defer rl.mu.Unlock()
    if l, ok := rl.limiters[ip]; ok {
        return l
    }
    l := rate.NewLimiter(rl.rate, rl.burst)
    rl.limiters[ip] = l
    return l
}

func RateLimitMiddleware(rl *RateLimiter) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            ip := r.RemoteAddr
            limiter := rl.Get(ip)
            if !limiter.Allow() {
                w.Header().Set("Retry-After", "60")
                http.Error(w, `{"error":"rate limit exceeded"}`, http.StatusTooManyRequests)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}

// Security Headers — previne XSS, clickjacking, etc.
func SecurityHeadersMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "1; mode=block")
        w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubDomains")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        next.ServeHTTP(w, r)
    })
}
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exercicios/ex23_seguranca.go` | 🏋️ Exercícios: JWT, bcrypt, RBAC, rate limiting |

---

## 📋 Exercícios

### 🟢 1. Sistema de Login/Registro
Implemente endpoints `/register` e `/login` com:
- Hash de senha com bcrypt (custo 12)
- JWT com claims customizadas (user_id, role, email)
- Refresh token com validade maior (7 dias) e access token curto (15 min)

### 🟡 2. RBAC com múltiplos roles
Estenda o middleware de autorização para suportar hierarquia de roles:
- `admin` > `moderator` > `user`
- Permissões específicas por endpoint (ex: `DELETE /resources` requer `admin`)

### 🟡 3. Rate Limiter por usuário
Implemente rate limiting por `user_id` (não por IP) usando o JWT decodificado:
- 100 requests/minuto para users
- 1000 requests/minuto para admins
- Retorne headers `X-RateLimit-Remaining` e `X-RateLimit-Reset`

### 🔴 4. JWT Blacklist (logout)
O problema do JWT stateless: um token válido não pode ser revogado.
Implemente uma blacklist em memória (+ Redis simulado) para tokens revogados:
- `POST /logout` → adiciona o JTI (JWT ID) à blacklist
- Middleware verifica se o JTI está na blacklist
- Expirar entradas da blacklist automaticamente quando o token expiraria de qualquer jeito

---

> **Próximo projeto**: [NexusMQ — Message Broker Distribuído](../projeto-final/README.md) com autenticação JWT!
