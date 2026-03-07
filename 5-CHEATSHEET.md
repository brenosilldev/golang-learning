# 📋 Go Cheatsheet — Referência Rápida

> Cola rápida para consultar enquanto coda. Não substitui os módulos — use para relembrar.

---

## Variáveis

```go
// Declaração curta (dentro de funções)
nome := "Go"
idade := 25
ativo := true

// Declaração explícita
var nome string = "Go"
var idade int              // zero value: 0
var ativo bool             // zero value: false

// Múltiplas
var (
    x int
    y string
    z bool
)

// Constantes
const PI = 3.14159
const (
    StatusOK    = 200
    StatusError = 500
)
```

## Tipos

```go
// Básicos
int, int8, int16, int32, int64
uint, uint8 (byte), uint16, uint32, uint64
float32, float64
string, bool, rune (= int32, para Unicode)

// Zero values
int → 0    string → ""    bool → false
pointer → nil    slice → nil    map → nil

// Conversão (Go NÃO converte implicitamente!)
i := 42
f := float64(i)
s := strconv.Itoa(i)     // int → string
i, _ = strconv.Atoi("42") // string → int
```

## Controle de Fluxo

```go
// if (sem parênteses!)
if x > 10 {
    // ...
} else if x > 5 {
    // ...
} else {
    // ...
}

// if com init statement
if err := fazerAlgo(); err != nil {
    log.Fatal(err)
}

// switch (sem break necessário!)
switch dia {
case "seg", "ter", "qua", "qui", "sex":
    fmt.Println("trabalho")
case "sab", "dom":
    fmt.Println("descanso")
default:
    fmt.Println("???")
}

// for (é o ÚNICO loop em Go)
for i := 0; i < 10; i++ { }   // clássico
for i < 10 { }                  // while
for { }                          // infinito
for i, v := range slice { }     // range
for k, v := range mapa { }      // range map
```

## Coleções

```go
// Array (tamanho fixo — raramente usado)
var arr [5]int
arr := [3]int{1, 2, 3}

// Slice (tamanho dinâmico — SEMPRE use isso)
s := []int{1, 2, 3}
s = append(s, 4, 5)
sub := s[1:3]              // [2, 3] (início incluso, fim excluso)
copia := make([]int, len(s))
copy(copia, s)

// Map
m := map[string]int{
    "alice": 100,
    "bob":   200,
}
m["carol"] = 300
delete(m, "alice")
valor, existe := m["bob"]  // comma ok pattern

// Iterar
for i, v := range slice { }
for chave, valor := range mapa { }
```

## Funções

```go
// Básica
func somar(a, b int) int {
    return a + b
}

// Múltiplos retornos
func dividir(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("divisão por zero")
    }
    return a / b, nil
}

// Retorno nomeado
func coordenadas() (x, y int) {
    x = 10
    y = 20
    return // retorna x e y automaticamente
}

// Variádica
func soma(nums ...int) int {
    total := 0
    for _, n := range nums {
        total += n
    }
    return total
}

// Função como variável / closure
dobro := func(x int) int { return x * 2 }

// Defer (executa ao sair da função, LIFO)
defer arquivo.Close()
```

## Ponteiros

```go
x := 42
p := &x     // p aponta para x
fmt.Println(*p)  // 42 (dereferência)
*p = 100    // x agora é 100

// Funções com ponteiro (modificam o original)
func dobrar(n *int) { *n *= 2 }
dobrar(&x)
```

## Structs

```go
type Pessoa struct {
    Nome  string
    Idade int
}

// Criar
p := Pessoa{Nome: "Alice", Idade: 30}
p2 := Pessoa{"Bob", 25}  // ordem dos campos
p3 := &Pessoa{Nome: "Carol"}  // ponteiro

// Métodos
func (p Pessoa) Saudacao() string {
    return "Olá, " + p.Nome
}

func (p *Pessoa) Aniversario() {  // pointer receiver para modificar
    p.Idade++
}

// Composição (em vez de herança)
type Funcionario struct {
    Pessoa          // embedding
    Cargo  string
}
```

## Interfaces

```go
// Definir
type Stringer interface {
    String() string
}

// Implementar (IMPLÍCITO — não precisa declarar)
func (p Pessoa) String() string {
    return fmt.Sprintf("%s (%d)", p.Nome, p.Idade)
}
// Pessoa automaticamente implementa Stringer!

// Interface vazia (aceita qualquer tipo)
var qualquer interface{} // ou: any (Go 1.18+)

// Type assertion
s, ok := qualquer.(string)

// Type switch
switch v := qualquer.(type) {
case string:  fmt.Println("string:", v)
case int:     fmt.Println("int:", v)
default:      fmt.Println("outro")
}
```

## Erros

```go
// Retornar erro
func abrir(nome string) (*Arquivo, error) {
    if nome == "" {
        return nil, errors.New("nome vazio")
    }
    return &Arquivo{}, nil
}

// Tratar erro (SEMPRE verifique!)
arq, err := abrir("teste.txt")
if err != nil {
    log.Fatal(err)  // ou return err
}

// Erro customizado
type NotFoundError struct {
    ID int
}
func (e *NotFoundError) Error() string {
    return fmt.Sprintf("ID %d não encontrado", e.ID)
}

// Wrap / Unwrap (Go 1.13+)
return fmt.Errorf("ao buscar user: %w", err)
if errors.Is(err, sql.ErrNoRows) { }
```

## Concorrência

```go
// Goroutine
go funcao()
go func() { fmt.Println("anônima") }()

// Channel
ch := make(chan int)       // sem buffer (bloqueia)
ch := make(chan int, 10)   // com buffer

ch <- 42        // enviar
valor := <-ch   // receber

// Select (switch para channels)
select {
case msg := <-ch1:
    fmt.Println(msg)
case ch2 <- "oi":
    fmt.Println("enviou")
case <-time.After(5 * time.Second):
    fmt.Println("timeout")
default:
    fmt.Println("nenhum pronto")
}

// WaitGroup
var wg sync.WaitGroup
wg.Add(1)
go func() {
    defer wg.Done()
    // trabalho
}()
wg.Wait()

// Mutex
var mu sync.Mutex
mu.Lock()
// acesso exclusivo
mu.Unlock()

// Context
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()
```

## Generics (Go 1.18+)

```go
func Map[T any, R any](slice []T, fn func(T) R) []R {
    result := make([]R, len(slice))
    for i, v := range slice {
        result[i] = fn(v)
    }
    return result
}

// Constraint
type Number interface {
    int | int64 | float64
}

func Soma[T Number](nums []T) T {
    var total T
    for _, n := range nums {
        total += n
    }
    return total
}
```

## Testes

```go
// arquivo: math_test.go
func TestSomar(t *testing.T) {
    resultado := Somar(2, 3)
    if resultado != 5 {
        t.Errorf("esperava 5, recebeu %d", resultado)
    }
}

// Table-driven
func TestSomar(t *testing.T) {
    tests := []struct{ a, b, want int }{
        {1, 2, 3},
        {0, 0, 0},
        {-1, 1, 0},
    }
    for _, tt := range tests {
        got := Somar(tt.a, tt.b)
        if got != tt.want {
            t.Errorf("Somar(%d,%d) = %d, want %d", tt.a, tt.b, got, tt.want)
        }
    }
}

// Rodar
// go test ./...
// go test -v -run TestSomar
// go test -bench=. -benchmem
```

## Pacotes Úteis

```go
fmt.Println()              // imprimir
fmt.Sprintf("x=%d", x)    // formatar string
strings.Contains(s, "go") // buscar
strings.Split(s, ",")     // separar
strconv.Itoa(42)           // int → string
strconv.Atoi("42")         // string → int
sort.Ints(slice)           // ordenar
json.Marshal(obj)          // struct → JSON
json.Unmarshal(data, &obj) // JSON → struct
os.ReadFile("f.txt")       // ler arquivo
os.WriteFile("f.txt", d, 0644) // escrever
time.Now()                 // tempo atual
time.Sleep(2 * time.Second)
log.Fatal(err)             // log + os.Exit(1)
```

---

> 💡 **Regra de ouro do Go**: se compila, provavelmente funciona. Se não compila, o erro é claro.
