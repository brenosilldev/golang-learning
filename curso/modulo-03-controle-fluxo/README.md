# Módulo 03 — Controle de Fluxo

[← Variáveis e Tipos](../modulo-02-variaveis-tipos/README.md) | [Próximo: Coleções →](../modulo-04-colecoes/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que Go tem apenas um tipo de loop (`for`)?
> 2. O que faz o `fallthrough` no `switch` e quando usá-lo?
> 3. O que é o "initializer statement" do `if` e por que ele existe?

---

## 1. if / else — Mais Poderoso que Parece

### Forma básica

```go
if x > 10 {
    fmt.Println("maior que 10")
} else if x == 10 {
    fmt.Println("igual a 10")
} else {
    fmt.Println("menor que 10")
}
```

> **Regra Go**: a chave `{` DEVE estar na mesma linha. Isso não é estilo — é sintaxe obrigatória.

### Initializer statement — O diferencial do Go

O `if` do Go aceita um statement antes da condição. O valor declarado existe **apenas dentro do bloco** `if/else`:

```go
// Forma idiomática — MUITO usada em Go
if val, ok := mapa["chave"]; ok {
    fmt.Println("encontrei:", val)
} else {
    fmt.Println("não existe")
}
// val e ok NÃO existem aqui — escopo limitado ao if

// Excelente para tratar erros localmente
if err := arquivo.Close(); err != nil {
    log.Printf("erro ao fechar arquivo: %v", err)
}

// Sem initializer (menos idiomático)
val, ok := mapa["chave"]
if ok {
    // val vaza para o escopo externo — às vezes indesejado
}
```

**Por que isso é valioso?** Reduz o escopo de variáveis de erro/resultado — código mais limpo e seguro.

---

## 2. switch — O switch do Go não tem fall-through por padrão

```go
// Não precisa de break — Go para automaticamente em cada case
switch dia {
case "segunda", "terça", "quarta", "quinta", "sexta":
    fmt.Println("dia útil")
case "sábado", "domingo":
    fmt.Println("fim de semana")
default:
    fmt.Println("dia inválido")
}
```

### switch com initializer (igual ao if)

```go
switch os := runtime.GOOS; os {
case "linux":
    fmt.Println("rodando no Linux")
case "darwin":
    fmt.Println("rodando no macOS")
default:
    fmt.Printf("sistema não suportado: %s\n", os)
}
```

### switch sem expressão — substitui cadeia de if/else if

```go
// Muito mais legível que vários else if
switch {
case temperatura < 0:
    fmt.Println("congelando")
case temperatura < 15:
    fmt.Println("frio")
case temperatura < 25:
    fmt.Println("agradável")
default:
    fmt.Println("quente")
}
```

### fallthrough — forçar execução do próximo case

```go
// fallthrough é explícito — diferente de C/Java onde é o padrão
switch n {
case 1:
    fmt.Println("um")
    fallthrough  // executa o case 2 também
case 2:
    fmt.Println("dois ou um")
case 3:
    fmt.Println("três")
}
// Para n=1: imprime "um" e "dois ou um"
// Para n=2: imprime apenas "dois ou um"
```

> **Na prática**: `fallthrough` é raramente usado. Se você precisa de `fallthrough` com frequência, provavelmente a estrutura do código pode ser melhorada.

---

## 3. for — O Único Loop do Go (e por que isso é bom)

Go tem **apenas** `for`. Mas o `for` pode se comportar como `while`, `do-while` e o `for` clássico. Menos palavras-chave = menos para aprender.

### As 4 formas do for

```go
// Forma 1: clássico (C-style)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// Forma 2: while-like (só a condição)
n := 1
for n < 100 {
    n *= 2
}

// Forma 3: infinito (use com break ou return)
for {
    input, _ := bufio.NewReader(os.Stdin).ReadString('\n')
    if input == "quit\n" {
        break
    }
    fmt.Println("você digitou:", input)
}

// Forma 4: range — iteração sobre coleções
numeros := []int{10, 20, 30, 40, 50}
for i, v := range numeros {
    fmt.Printf("índice=%d, valor=%d\n", i, v)
}
```

### range em diferentes tipos

```go
// Slice/Array: índice e valor
for i, v := range slice {}
for _, v := range slice {}   // ignora índice
for i := range slice {}      // só índice

// Map: chave e valor (ordem NÃO é garantida)
for k, v := range mapa {}
for k := range mapa {}       // só chaves

// String: índice de byte e rune (caractere Unicode)
for i, r := range "Olá" {
    fmt.Printf("byte=%d, char=%c\n", i, r)
}

// Channel: recebe valores até o channel fechar
for v := range ch {}
```

### break e continue com labels — para loops aninhados

```go
// Label: sai do loop externo de dentro do loop interno
outer:
    for i := 0; i < 5; i++ {
        for j := 0; j < 5; j++ {
            if i+j == 6 {
                break outer  // sai dos DOIS loops
            }
            fmt.Printf("(%d,%d) ", i, j)
        }
    }
```

---

## 4. Armadilhas Comuns

### Armadilha 1: Modificar slice durante range

```go
// ❌ Problemático — comportamento indefinido modificar o slice sendo iterado
for i, v := range numeros {
    if v == 3 {
        numeros = append(numeros, 99) // pode funcionar, mas é confuso
    }
}

// ✅ Melhor — itere por índice ou faça uma cópia
for i := 0; i < len(numeros); i++ {
    if numeros[i] == 3 {
        numeros = append(numeros, 99)
    }
}
```

### Armadilha 2: Captura de variável de loop em goroutine (Go < 1.22)

```go
// ❌ Bug clássico antes do Go 1.22
for _, v := range items {
    go func() {
        fmt.Println(v) // captura a VARIÁVEL v, não o VALOR
        // quando a goroutine roda, v provavelmente já é o último item
    }()
}

// ✅ Fix (necessário no Go < 1.22)
for _, v := range items {
    v := v  // nova variável v no escopo do loop
    go func() {
        fmt.Println(v) // agora captura o valor correto
    }()
}

// ✅ Go 1.22+: o comportamento foi corrigido automaticamente
// Cada iteração do loop cria uma nova variável
```

### Armadilha 3: range em map não tem ordem

```go
m := map[string]int{"a": 1, "b": 2, "c": 3}

// Cada execução pode imprimir em ordem diferente!
for k, v := range m {
    fmt.Printf("%s: %d\n", k, v) // ordem não garantida
}

// Para ordem consistente: ordene as chaves
keys := make([]string, 0, len(m))
for k := range m {
    keys = append(keys, k)
}
sort.Strings(keys)
for _, k := range keys {
    fmt.Printf("%s: %d\n", k, m[k])
}
```

---

## 5. Patterns Idiomáticos

### Early return — prefira a este em vez de if aninhados

```go
// ❌ Pirâmide da morte — difícil de ler
func processar(req *Request) error {
    if req != nil {
        if req.Valido() {
            if usuario, err := buscarUsuario(req.ID); err == nil {
                // lógica principal afundada em aninhamentos
            }
        }
    }
    return nil
}

// ✅ Guard clauses — retorne cedo, mantenha o caminho feliz sem aninhamento
func processar(req *Request) error {
    if req == nil {
        return errors.New("request nulo")
    }
    if !req.Valido() {
        return errors.New("request inválido")
    }
    usuario, err := buscarUsuario(req.ID)
    if err != nil {
        return fmt.Errorf("buscar usuário: %w", err)
    }
    // lógica principal aqui, sem aninhamento
    return nil
}
```

---

## ✅ Checklist de Controle de Fluxo

- [ ] `if` com initializer para variáveis de escopo local (especialmente erros)
- [ ] `switch` sem expressão como alternativa mais legível a `else if` encadeados
- [ ] `range` com `_` quando o índice ou valor não é necessário
- [ ] Goroutines em loops capturam o valor correto da variável (Go 1.22+ ou `v := v`)
- [ ] Early return (guard clauses) ao invés de aninhamentos profundos

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo03_fluxo.go` | if, switch, for, range, break/continue, labels |
| `exercicios/ex03_fluxo.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. FizzBuzz com switch
Implemente FizzBuzz (1 a 100) usando `switch` sem expressão (não use `if/else`). Adicione: múltiplos de 7 imprimem "Bazz", múltiplos de 3, 5 e 7 imprimem "FizzBuzzBazz".

### 🟢 2. Buscador de Palavras
Dada uma `[]string`, use `for range` para:
- Contar quantas strings têm mais de 5 caracteres
- Encontrar a primeira string que começa com "Go"
- Verificar se TODAS as strings são não-vazias (retorne `bool`)

### 🟡 3. Calculadora de Primos
Use um loop `for` para encontrar todos os primos até N usando o Crivo de Eratóstenes. Imprima os primeiros 20 primos. Use `continue` para pular não-primos.

### 🟡 4. Analisador de Strings UTF-8
Use `for range` sobre uma string com emojis e caracteres especiais para:
- Contar bytes vs runes (caracteres)
- Listar cada caractere com seu índice de byte e code point Unicode
- Identificar caracteres não-ASCII

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Variáveis e Tipos](../modulo-02-variaveis-tipos/README.md) | [Próximo: Coleções →](../modulo-04-colecoes/README.md)
