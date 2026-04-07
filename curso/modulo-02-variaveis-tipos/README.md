# Módulo 02 — Variáveis e Tipos

[← Introdução](../modulo-01-introducao/README.md) | [Próximo: Controle de Fluxo →](../modulo-03-controle-fluxo/README.md)

---

> **Antes de ler — tente responder:**
> 1. O que é o "zero value" e por que Go tem esse conceito?
> 2. Qual a diferença entre `var x int` e `x := 0`?
> 3. Por que Go não faz conversão implícita entre `int` e `float64`?

---

## 1. Declaração de Variáveis — As Duas Formas

Go tem exatamente **duas formas** de declarar variáveis. Entender quando usar cada uma é fundamental:

```go
// Forma 1: var — pode usar em qualquer escopo (incluindo package-level)
var nome string = "Go"
var idade int = 10

// Forma 1b: var com inferência de tipo (Go deduz o tipo)
var nome = "Go"       // Go infere: string
var ativo = true      // Go infere: bool

// Forma 2: := (short declaration) — APENAS dentro de funções
nome := "Go"          // string inferida
idade := 10           // int inferido
pi := 3.14            // float64 inferido

// Declaração múltipla — elegante para relacionados
x, y := 10, 20
min, max := encontrarMinMax(numeros)
```

### Quando usar cada forma?

| Situação | Use | Exemplo |
|----------|-----|---------|
| Dentro de função | `:=` | `resultado := calcular()` |
| Nível do pacote (global) | `var` | `var db *sql.DB` |
| Tipo precisa ser explícito | `var` | `var buf bytes.Buffer` |
| Zero value intencional | `var` | `var count int` (mais claro que `count := 0`) |
| Múltiplas variáveis relacionadas | `var (...)` | veja abaixo |

```go
// Bloco var — agrupa variáveis relacionadas no topo da função ou pacote
var (
    host     = "localhost"
    port     = 5432
    maxConns = 25
    timeout  = 30 * time.Second
)
```

---

## 2. Zero Values — O Sistema de Segurança do Go

Em Go, **toda variável tem um valor padrão**. Não existe "undefined", "null reference" ou valor lixo de memória (como em C).

```go
var i int       // 0
var f float64   // 0.0
var s string    // ""  (string vazia, não nil!)
var b bool      // false
var p *int      // nil (ponteiro nil)
var sl []int    // nil (slice nil)
var m map[string]int // nil (map nil)
```

### Por que isso importa na prática?

```go
// ✅ Você pode usar um contador sem inicializar
var contagem int
contagem++ // funciona! começa do zero value 0

// ✅ Você pode appendar em um slice nil — funciona!
var nomes []string
nomes = append(nomes, "Alice") // slice nil é equivalente a []string{}

// ❌ CUIDADO: escrever em map nil causa PANIC
var m map[string]int
m["chave"] = 1 // PANIC: assignment to entry in nil map

// ✅ Map precisa ser inicializado com make
m := make(map[string]int)
m["chave"] = 1 // seguro
```

---

## 3. Tipos Primitivos — Guia Prático

### Inteiros

```go
// Regra prática: use int para quase tudo
// Go garante que int tem o tamanho do ponteiro nativo (32 ou 64 bits)
var n int = 42

// Use tamanhos específicos quando:
// - protocolo de rede/arquivo especifica o tamanho
// - precisa de uint8 para byte individual
var b byte = 255     // byte = uint8, para dados binários
var r rune = '世'     // rune = int32, para caracteres Unicode
```

| Tipo | Tamanho | Range | Use quando |
|------|---------|-------|------------|
| `int` | 32 ou 64 bits | depende da plataforma | **Padrão** — use este |
| `int8` | 8 bits | -128 a 127 | Protocolo de rede específico |
| `int16` | 16 bits | -32768 a 32767 | Protocolo de rede específico |
| `int32` / `rune` | 32 bits | ±2 bilhões | Caractere Unicode |
| `int64` | 64 bits | ±9 quintilhões | IDs grandes, timestamps Unix |
| `uint8` / `byte` | 8 bits | 0 a 255 | Dados binários, bytes individuais |

### Float

```go
// float64 é o padrão — use sempre a não ser que tenha motivo específico
var f float64 = 3.14159

// ⚠️ Armadilha clássica: floats não são exatos!
fmt.Println(0.1 + 0.2)          // 0.30000000000000004
fmt.Println(0.1 + 0.2 == 0.3)   // false!

// Para dinheiro: use int (centavos) ou github.com/shopspring/decimal
preco := 1999  // R$ 19,99 em centavos — nunca float para dinheiro!
```

### String

```go
// String em Go é imutável e UTF-8 por padrão
s := "Olá, 世界"

// len() retorna BYTES, não caracteres!
fmt.Println(len("hello"))  // 5 (5 bytes)
fmt.Println(len("世界"))    // 6 (cada ideograma = 3 bytes em UTF-8)

// Para iterar por caracteres Unicode, use range
for i, r := range "Olá" {
    fmt.Printf("índice=%d, rune=%c (U+%04X)\n", i, r, r)
}
// índice=0, rune=O (U+004F)
// índice=1, rune=l (U+006C)
// índice=2, rune=á (U+00E1)  ← ocupa 2 bytes, por isso próximo índice é 4

// Strings são concatenadas com + mas para loops use strings.Builder
var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString("a") // eficiente — aloca uma vez
}
resultado := sb.String()
```

---

## 4. Constantes e iota

```go
// Constantes são avaliadas em tempo de compilação
const Pi = 3.14159
const MaxRetries = 3
const Prefixo = "api/v1"

// Bloco de constantes
const (
    KB = 1024
    MB = KB * 1024
    GB = MB * 1024
)

// iota — enumeração automática, começa em 0, incrementa por constante
type DiaSemana int
const (
    Domingo DiaSemana = iota  // 0
    Segunda                   // 1
    Terca                     // 2
    Quarta                    // 3
    Quinta                    // 4
    Sexta                     // 5
    Sabado                    // 6
)

// iota com expressão — bitmask de permissões
type Permissao int
const (
    Leitura   Permissao = 1 << iota // 1 (0001)
    Escrita                         // 2 (0010)
    Execucao                        // 4 (0100)
)

// Combinando com bitwise OR
minhaPermissao := Leitura | Escrita     // 3 (0011)
temLeitura := minhaPermissao & Leitura  // 1 (verdadeiro)
```

---

## 5. Conversão de Tipos — Sempre Explícita

Go **nunca** converte tipos automaticamente. Isso parece chato mas previne uma classe inteira de bugs.

```go
var i int = 42
var f float64

// ❌ ERRO DE COMPILAÇÃO — tipo incompatível
f = i

// ✅ Conversão explícita
f = float64(i)

// Conversões comuns
n := 42
f64 := float64(n)       // int → float64
i32 := int32(n)         // int → int32
s := strconv.Itoa(n)    // int → string  (NÃO use string(n) — dá rune!)
n2, err := strconv.Atoi("42") // string → int

// ⚠️ Armadilha clássica!
fmt.Println(string(65))   // "A"  ← converte int para rune, depois para string
fmt.Println(strconv.Itoa(65)) // "65" ← isso que você provavelmente quer
```

### Conversões com perda de dados

```go
// ⚠️ Go não avisa sobre overflow — você precisa verificar
big := int64(300)
small := int8(big) // silenciosamente torna-se 44 (300 - 256 = 44)

// Para conversões seguras, verifique os limites:
if big > math.MaxInt8 || big < math.MinInt8 {
    return fmt.Errorf("valor %d não cabe em int8", big)
}
```

---

## 6. Armadilhas Comuns (que derrubam desenvolvedores experientes)

### Armadilha 1: Shadowing com :=

```go
x := 10
if condicao {
    x := 20        // ← NOVA variável x, não modifica a de cima!
    fmt.Println(x) // 20
}
fmt.Println(x) // ainda 10 — shadowing silencioso
```

### Armadilha 2: Variável declarada e não usada

```go
func main() {
    x := 10  // ← ERRO DE COMPILAÇÃO se não usar x
    // Go não permite variáveis declaradas e não usadas
    // Use _ para valores que você quer descartar explicitamente
    _, err := strconv.Atoi("abc")
    if err != nil {
        fmt.Println("erro")
    }
}
```

### Armadilha 3: Múltiplos retornos com :=

```go
conn, err := net.Dial("tcp", "localhost:8080")
// ... mais tarde ...
resp, err := http.Get("https://go.dev") // ← ok! err é reatribuída, resp é nova
```

---

## ✅ Checklist de Variáveis e Tipos

- [ ] Uso `int` por padrão (não `int64` sem motivo)
- [ ] Nunca uso `float64` para valores monetários — uso `int` em centavos
- [ ] Inicializo mapas com `make` antes de escrever neles
- [ ] Uso `strings.Builder` para construção de strings em loop
- [ ] Conversões de tipo são sempre explícitas
- [ ] Uso `strconv.Itoa` (não `string(n)`) para converter int para string

---

## 📂 Arquivos deste módulo

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/exemplo02_variaveis.go` | Declaração, tipos, zero values, conversão, iota |
| `exercicios/ex02_tipos.go` | 🏋️ Exercícios para praticar |

---

## 📋 Exercícios

### 🟢 1. Calculadora de IMC
Declare variáveis `peso` (float64) e `altura` (float64). Calcule o IMC e imprima a classificação (abaixo do peso, normal, sobrepeso, obeso). Use constantes para os limites de cada categoria.

### 🟢 2. Temperatura
Crie um conversor de temperatura que usa constantes `AbsoluteZeroC = -273.15` e funções para converter entre Celsius, Fahrenheit e Kelvin. Mostre os zero values de cada variável antes de inicializar.

### 🟡 3. Sistema de Permissões com iota
Use iota para criar um sistema de permissões bitmask (Leitura, Escrita, Execução, Admin). Crie funções `TemPermissao(user, perm Permissao) bool` e `AdicionarPermissao(user, perm Permissao) Permissao`.

### 🟡 4. Contador de Bytes vs Runes
Escreva uma função que recebe uma string e retorna: número de bytes, número de runes (caracteres Unicode) e um slice com cada rune. Teste com strings ASCII e com strings contendo emojis/caracteres especiais.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← Introdução](../modulo-01-introducao/README.md) | [Próximo: Controle de Fluxo →](../modulo-03-controle-fluxo/README.md)
