# Módulo 03 — Controle de Fluxo

[← Variáveis e Tipos](../modulo-02-variaveis-tipos/README.md) | [Próximo: Coleções →](../modulo-04-colecoes/README.md)

---

## 📖 Teoria

### if / else
```go
if x > 10 {
    fmt.Println("maior")
} else if x == 10 {
    fmt.Println("igual")
} else {
    fmt.Println("menor")
}
```

**Diferencial do Go** — `if` com inicializador:
```go
if valor, ok := mapa["chave"]; ok {
    fmt.Println(valor)
}
// 'valor' e 'ok' só existem dentro deste if
```

### switch
```go
switch dia {
case "segunda", "terça":
    fmt.Println("início da semana")
case "sexta":
    fmt.Println("sextou!")
default:
    fmt.Println("outro dia")
}
```
- **Não precisa de `break`** — Go para automaticamente
- `fallthrough` força continuar para o próximo case

### for (o ÚNICO loop do Go)
```go
for i := 0; i < 10; i++ {}           // clássico
for i < 10 {}                         // while-like
for {}                                 // infinito
for i, v := range slice {}            // iteração
for k, v := range mapa {}             // iteração em map
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo03_fluxo.go` | if, switch, for, range, break/continue |
| `exercicios/ex03_fluxo.go` | 🏋️ Exercícios |

---

[← Variáveis e Tipos](../modulo-02-variaveis-tipos/README.md) | [Próximo: Coleções →](../modulo-04-colecoes/README.md)
