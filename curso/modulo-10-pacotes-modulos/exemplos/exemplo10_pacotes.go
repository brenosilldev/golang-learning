package main

import (
	"fmt"
	"math"
	"os"
	"strings"
)

// ============================================================================
// MÓDULO 10 — Pacotes e Módulos
// ============================================================================
//
// Conceitos:
//   - Organização de pacotes
//   - Visibilidade (maiúscula = exportado)
//   - Import paths e alias
//   - go mod init / go mod tidy / go get
//   - Estrutura de projeto Go
//   - internal/ (proteção em nível de compilador)
//   - init() — função especial de inicialização
//
// Rode com: go run exemplo10_pacotes.go
// ============================================================================

// init() é chamada automaticamente ANTES de main()
// Cada pacote pode ter múltiplas init()
// Ordem: variáveis globais → init() → main()
func init() {
	fmt.Println("[init] Pacote inicializado!")
}

// Variáveis de pacote (globais ao pacote)
var versao = "1.0.0"

func main() {
	// ==========================================================
	fmt.Println("\n=== IMPORTS ===")

	// Import simples — já estamos usando
	fmt.Println("Hello!")

	// Usando vários pacotes da stdlib
	fmt.Println("Pi:", math.Pi)
	fmt.Println("Upper:", strings.ToUpper("hello"))

	// os — acesso ao sistema operacional
	hostname, _ := os.Hostname()
	fmt.Println("Hostname:", hostname)
	fmt.Println("Args:", os.Args)

	// ==========================================================
	fmt.Println("\n=== VISIBILIDADE ===")

	// EM GO:
	// - NomeComMaiuscula → EXPORTADO (público) — acessível fora do pacote
	// - nomeComMinuscula → NÃO exportado (privado) — apenas dentro do pacote

	// Exemplos de exportados na stdlib:
	// fmt.Println ✅ (P maiúsculo)
	// math.Pi     ✅ (P maiúsculo)
	// os.Exit     ✅ (E maiúsculo)

	// Isso se aplica a:
	// - Funções
	// - Tipos (structs, interfaces)
	// - Campos de structs
	// - Constantes
	// - Variáveis

	type PessoaPublica struct {
		Nome  string // campo exportado
		Email string // campo exportado
		idade int    // campo NÃO exportado (minúscula)
	}

	p := PessoaPublica{Nome: "Go", Email: "go@go.dev"}
	fmt.Printf("Pessoa: %+v\n", p)
	// p.idade não seria acessível de outro pacote

	// ==========================================================
	fmt.Println("\n=== IMPORT ALIAS ===")

	// Você pode dar alias para imports:
	// import (
	//     f "fmt"                           // alias
	//     . "math"                          // dot import (traz tudo pro escopo)
	//     _ "github.com/lib/pq"             // blank import (só roda init())
	//     meuhttp "net/http"                // rename para evitar conflito
	// )

	// Blank import é usado para:
	// - Registrar drivers de banco de dados
	// - Registrar codecs de imagem
	// - Qualquer situação onde o init() é necessário

	// ==========================================================
	fmt.Println("\n=== ESTRUTURA DE PROJETO ===")

	fmt.Println(`
Estrutura recomendada para projetos Go:

projeto/
├── cmd/                    # Cada subpasta = um executável
│   ├── api/
│   │   └── main.go         # go run ./cmd/api
│   └── cli/
│       └── main.go         # go run ./cmd/cli
│
├── internal/               # Código PRIVADO do projeto
│   ├── handler/            # Handlers HTTP
│   │   └── user.go
│   ├── service/            # Lógica de negócio
│   │   └── user.go
│   ├── repository/         # Acesso a dados
│   │   └── user.go
│   └── model/              # Structs/tipos
│       └── user.go
│
├── pkg/                    # Código público (pode ser importado)
│   └── validator/
│       └── validator.go
│
├── go.mod                  # Definição do módulo
├── go.sum                  # Lock file
├── Makefile                # Comandos úteis
└── README.md
`)

	// ==========================================================
	fmt.Println("=== GO MOD — Comandos essenciais ===")

	fmt.Println(`
# Criar novo módulo
go mod init github.com/usuario/projeto

# Adicionar dependência
go get github.com/gin-gonic/gin@latest

# Limpar dependências não usadas
go mod tidy

# Ver dependências
go list -m all

# Ver por que uma dependência é necessária
go mod why github.com/pacote

# Vendor (copiar deps para pasta local)
go mod vendor

# Atualizar todas as dependências
go get -u ./...
`)

	// ==========================================================
	fmt.Println("=== INIT() — Inicialização ===")
	fmt.Println("Versão do app:", versao)
	fmt.Println("A mensagem [init] apareceu ANTES de main()!")
	fmt.Println()
	fmt.Println("Usos comuns de init():")
	fmt.Println("  - Validar configuração obrigatória")
	fmt.Println("  - Registrar drivers/codecs")
	fmt.Println("  - Inicializar variáveis complexas")
	fmt.Println("  - ⚠️ Evite lógica pesada no init()!")
}
