# Módulo 01 — Introdução ao Go

[← Voltar ao índice](../README.md) | [Próximo: Variáveis e Tipos →](../modulo-02-variaveis-tipos/README.md)

---

## 📖 O que é Go?

Go é uma linguagem **compilada**, **estáticamente tipada** e com **garbage collector** criada pelo Google. Ela foi projetada para ser:

- **Simples** — Poucas palavras-chave (25!), sintaxe limpa
- **Rápida** — Compila para código nativo em segundos
- **Concorrente** — Goroutines e channels fazem parte da linguagem
- **Prática** — Ferramental completo built-in

### Filosofia do Go

```
"Less is exponentially more" — Rob Pike
```

Go **intencionalmente** não tem:
- ❌ Classes e herança (usa composição)
- ❌ Exceções (usa retorno de `error`)
- ❌ Generics complexos (adicionados de forma simples em Go 1.18)
- ❌ Operadores sobrecarregados
- ❌ Metaprogramação mágica

Isso pode parecer limitante, mas na prática torna o código **extremamente legível** — qualquer dev Go lê qualquer projeto Go com facilidade.

---

## 🛠️ Instalação

### Linux (Ubuntu/Debian)
```bash
# Baixar a versão mais recente
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz

# Extrair para /usr/local
sudo rm -rf /usr/local/go
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz

# Adicionar ao PATH (colocar no ~/.bashrc ou ~/.zshrc)
export PATH=$PATH:/usr/local/go/bin
export GOPATH=$HOME/go
export PATH=$PATH:$GOPATH/bin

# Verificar instalação
go version
```

### macOS
```bash
brew install go
```

### Windows
Baixe o instalador em [go.dev/dl](https://go.dev/dl/) e siga o wizard.

---

## 📁 Estrutura de um Programa Go

Todo programa Go segue esta estrutura básica:

```go
package main      // 1. Declaração do pacote

import "fmt"      // 2. Importações

func main() {     // 3. Função principal (ponto de entrada)
    fmt.Println("Olá, Go!")
}
```

### Regras importantes:

| Regra | Detalhe |
|-------|---------|
| `package main` | Todo executável precisa de um pacote `main` |
| `func main()` | É o ponto de entrada — o programa começa aqui |
| Importações | Só pode importar o que realmente usa (senão erro de compilação!) |
| Chave `{` | DEVE ficar na **mesma linha** da declaração |
| Ponto e vírgula | **Não precisa** — o compilador insere automaticamente |

### O que acontece se importar algo sem usar?

```go
package main

import "fmt"   // ← ERRO: imported and not used
import "os"    // ← ERRO: imported and not used

func main() {
}
```

Go não permite "lixo" no código. Isso mantém o código limpo por design.

---

## 🔧 Ferramentas Essenciais

Go vem com um toolkit completo. Você não precisa instalar nada extra para começar:

```bash
# Rodar um programa sem compilar (modo script)
go run main.go

# Compilar para um binário
go build main.go
./main

# Compilar para outro OS/arquitetura
GOOS=linux GOARCH=amd64 go build main.go
GOOS=windows GOARCH=amd64 go build main.go
GOOS=darwin GOARCH=arm64 go build main.go

# Formatar código automaticamente
go fmt ./...

# Analisar código por problemas comuns
go vet ./...

# Inicializar um módulo (projeto)
go mod init nome-do-modulo

# Baixar dependências
go mod tidy

# Rodar testes
go test ./...

# Ver documentação de um pacote
go doc fmt
go doc fmt.Println
```

### Cross-compilation — Superpoder do Go

Isso é INCRÍVEL: com **uma linha** você compila para qualquer OS:

```bash
# Compilar no Linux um binário para Windows
GOOS=windows go build -o app.exe main.go

# Compilar no Mac um binário para Linux
GOOS=linux go build -o app main.go
```

Nenhuma VM, nenhum runtime necessário no destino. Apenas o binário.

---

## 📝 Primeiro Programa Completo

Vamos criar algo um pouco mais completo que um "Hello World":

```go
package main

import (
    "fmt"
    "runtime"
    "time"
)

func main() {
    // Informações do sistema
    fmt.Println("=== Meu Primeiro Programa Go ===")
    fmt.Println()

    // Variáveis com declaração curta
    nome := "Breno"
    linguagem := "Go"
    experiencia := 3

    // Printf com formatação
    fmt.Printf("👋 Olá! Eu sou %s\n", nome)
    fmt.Printf("🔧 Estou aprendendo %s\n", linguagem)
    fmt.Printf("📅 Tenho %d anos de experiência em programação\n", experiencia)
    fmt.Println()

    // Informações do Go e do sistema
    fmt.Printf("🐹 Versão do Go: %s\n", runtime.Version())
    fmt.Printf("💻 Sistema: %s/%s\n", runtime.GOOS, runtime.GOARCH)
    fmt.Printf("🧵 CPUs disponíveis: %d\n", runtime.NumCPU())
    fmt.Printf("⏰ Data/Hora: %s\n", time.Now().Format("02/01/2006 15:04:05"))
}
```

**Rode com:**
```bash
go run curso/modulo-01-introducao/exercicios/ex01_hello.go
```

---

## 📦 Criando um Módulo (Projeto)

Todo projeto Go real usa **módulos**. Um módulo é um conjunto de pacotes com versionamento:

```bash
# Criar um novo projeto
mkdir meu-projeto
cd meu-projeto
go mod init github.com/seu-usuario/meu-projeto

# Isso cria um arquivo go.mod:
```

```
module github.com/seu-usuario/meu-projeto

go 1.22
```

O `go.mod` é como o `package.json` do Node ou `Cargo.toml` do Rust — ele rastreia dependências.

---

## 🔍 Convenções Importantes

Desde o início, siga estas convenções que o Go impõe:

1. **Nomes exportados começam com maiúscula**
   ```go
   func Publica() {}   // ✅ Acessível fora do pacote
   func privada() {}   // ❌ Só acessível dentro do pacote
   ```

2. **Formatação é obrigatória**
   ```bash
   go fmt ./...  # Formata todo o projeto
   ```
   Não existe "war" de estilo em Go. Todo mundo usa o mesmo formato.

3. **Nomes curtos e descritivos**
   ```go
   // ✅ Bom
   func (s *Server) Start() error {}
   
   // ❌ Ruim (verbose demais)
   func (server *HTTPServer) StartServer() error {}
   ```

4. **camelCase para tudo** (não snake_case)
   ```go
   userName := "breno"    // ✅
   user_name := "breno"   // ❌ não é idiomático
   ```

---

## ✏️ Exercícios

### Exercício 1.1 — Cartão de Visita
Crie um programa que imprime um "cartão de visita" formatado com:
- Seu nome
- Sua idade
- Sua linguagem favorita
- O número de CPUs do seu computador
- A versão do Go instalada

### Exercício 1.2 — Compilação Cross-Platform
1. Compile seu programa para **3 plataformas diferentes** (Linux, Windows, macOS)
2. Liste os binários gerados e seus tamanhos

### Exercício 1.3 — Explorar `go doc`
Use `go doc` para descobrir:
1. O que faz `fmt.Sprintf`?
2. O que faz `strings.Contains`?
3. O que faz `os.Args`?

> 💡 Veja o arquivo `exercicios/ex01_hello.go` para a solução do Exercício 1.1

---

[← Voltar ao índice](../README.md) | [Próximo: Variáveis e Tipos →](../modulo-02-variaveis-tipos/README.md)
