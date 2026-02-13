# Módulo 10 — Pacotes e Módulos

[← Tratamento de Erros](../modulo-09-tratamento-erros/README.md) | [Próximo: Concorrência →](../modulo-11-concorrencia/README.md)

---

## 📖 Teoria

### Módulos Go
```bash
go mod init github.com/usuario/projeto  # criar módulo
go mod tidy                              # limpar/baixar deps
go get github.com/pacote@v1.2.3         # adicionar dependência
```

### Visibilidade em Go
```go
func Publica() {}   // começa com maiúscula = exportado ✅
func privada() {}   // começa com minúscula = interno ❌
type Pessoa struct { // tipo exportado
    Nome  string     // campo exportado
    idade int        // campo privado
}
```

### Estrutura recomendada de projeto
```
projeto/
├── cmd/             # pontos de entrada (main packages)
│   └── api/
│       └── main.go
├── internal/        # código que NÃO pode ser importado externamente
│   ├── handler/
│   ├── service/
│   └── repository/
├── pkg/             # código que PODE ser importado externamente
├── go.mod
└── go.sum
```

### internal/ — A pasta mágica
O Go **proíbe** que código de fora do módulo importe `internal/`. É proteção em nível de compilador!

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo10_pacotes.go` | Organização, imports, visibilidade |
| `exercicios/ex10_pacotes.go` | 🏋️ Exercícios |

---

[← Tratamento de Erros](../modulo-09-tratamento-erros/README.md) | [Próximo: Concorrência →](../modulo-11-concorrencia/README.md)
