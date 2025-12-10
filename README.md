## Guia Rápido – Golang Learning

Repositório com exemplos curtos para revisar Go de forma rápida.  
Cada seção abaixo tem **o que é** e **um exemplo mínimo**.

---

## 1. Fundamentos

### 1.1 Variáveis, tipos e operações (`fundamentos/aula1.go`)
- **Ideia**: declaração de variáveis, tipos numéricos e operações básicas.

```go
package main

import "fmt"

func main() {
	a := -1
	b := 2
	metro := 2.0

	fmt.Println(a, b, metro)
	fmt.Println(a + b)
	fmt.Println(float32(metro + 0.2*float64(a+b)))
}
```

### 1.2 Arrays e slices (`fundamentos/arrays.go`)
- **Ideia**: uso de slices, `append`, filtragem e “fatiamento” (`[:]`).

```go
lista := []int{2, 30, 5, 4, 8}
filtrada := make([]int, 0)

for _, v := range lista {
	if v < 10 {
		filtrada = append(filtrada, v)
	}
}

primeiros := lista[:3] // 0 até 2
fmt.Println(filtrada, primeiros)
```

### 1.3 Conversão de tipos (`fundamentos/conversao.go`)
- **Ideia**: converter entre `float64`, `int` e imprimir tipo.

```go
x := 3.2
y := 4.4
z := float64(x) + float64(y)
total := int(z)

fmt.Println(z, total)
fmt.Printf("%T\n", total) // int
```

### 1.4 Estruturas condicionais (`fundamentos/estrututuras.go`)
- **Ideia**: uso de `if/else if/else` com operadores lógicos.

```go
salario := 1100.0

if salario <= 1000 {
	fmt.Println("Faixa baixa")
} else if salario <= 5000 {
	fmt.Println("Faixa média")
} else {
	fmt.Println("Faixa alta")
}
```

### 1.5 Laços `for` (`fundamentos/for.go`)
- **Ideia**: `for` simples e aninhado (tabuada).

```go
for i := 0; i < 5; i++ {
	if i%2 == 0 {
		fmt.Println(i, "par")
	} else {
		fmt.Println(i, "ímpar")
	}
}

for base := 1; base <= 3; base++ {
	for i := 1; i <= 3; i++ {
		fmt.Printf("%d x %d = %d\n", base, i, base*i)
	}
}
```

### 1.6 Maps (`fundamentos/maps.go`)
- **Ideia**: criação de `map`, leitura segura, `range` e `delete`.

```go
m := map[string]int{"a": 1, "b": 2}

v, ok := m["c"]
if ok {
	fmt.Println("existe:", v)
} else {
	fmt.Println("não existe")
}

for chave, valor := range m {
	fmt.Println(chave, valor)
}

delete(m, "a")
fmt.Println(m)
```

### 1.7 Diferença entre arrays e slices (`fundamentos/slicesvsarrays.go`)
- **Ideia**: array tamanho fixo, slice dinâmico.

```go
array := [3]int{1, 2, 3}   // tamanho fixo
slice := []int{1, 2, 3, 4} // tamanho variável

slice = append(slice, 5)
fmt.Println(array, slice)
```

---

## 2. Recursos Avançados

### 2.1 Funções com múltiplos retornos (`avancado/func.go`)
- **Ideia**: função que retorna mais de um valor.

```go
func somaMensagem(a, b int) (int, string) {
	return a + b, "resultado da soma"
}

total, msg := somaMensagem(10, 20)
fmt.Println(total, msg)
```

### 2.2 `defer` e arquivos (`avancado/defer.go`)
- **Ideia**: garantir fechamento de recurso ao fim da função.

```go
file, err := os.Create("arquivo.txt")
if err != nil {
	panic(err)
}
defer file.Close()

file.Write([]byte("Hello, World!"))
fmt.Println("Arquivo criado")
```

### 2.3 `panic` e `recover` (`avancado/panic.go`, `avancado/recover.go`)
- **Ideia**: `panic` dispara erro fatal; `recover` captura para continuar.

```go
func comRecover() {
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recuperado:", r)
		}
	}()

	panic("erro fatal")
}
```

### 2.4 Ponteiros (`avancado/ponteiros.go`)
- **Ideia**: acessar e alterar valor pelo endereço de memória.

```go
x := 5
p := &x      // ponteiro para x
fmt.Println(*p) // 5

*p = 10      // altera x via ponteiro
fmt.Println(x)  // 10
```

---

## 3. Concorrência: Goroutines e Channels

### 3.1 Goroutines básicas (`gorountines/arq1.go`)
- **Ideia**: rodar funções em paralelo com `go`.

```go
func tarefa1() { fmt.Println("goroutine 1") }
func tarefa2() { fmt.Println("goroutine 2") }

func main() {
	go tarefa1()
	go tarefa2()
	time.Sleep(time.Second) // dá tempo de rodar
}
```

### 3.2 `sync.WaitGroup` (`gorountines/arq2.go`)
- **Ideia**: esperar todas as goroutines terminarem.

```go
var wg sync.WaitGroup

func tarefa(&wg *sync.WaitGroup) {
	defer wg.Done()
	time.Sleep(time.Second)
	fmt.Println("fim tarefa")
}

func main() {
	wg.Add(2)
	go tarefa(&wg)
	go tarefa(&wg)
	wg.Wait()
}
```

### 3.3 `sync.Mutex` (`gorountines/arq3.go`)
- **Ideia**: proteger variável compartilhada.

```go
var (
	i     = 0
	mutex sync.Mutex
)

for x := 0; x < 100; x++ {
	go func() {
		mutex.Lock()
		i++
		mutex.Unlock()
	}()
}

time.Sleep(time.Second)
fmt.Println("i =", i)
```

### 3.4 Channels (`gorountines/chanel.go`)
- **Ideia**: comunicação entre goroutines via canal com buffer.

```go
ch := make(chan int, 3)

go func() {
	for i := 0; i < 3; i++ {
		ch <- i
	}
	close(ch)
}()

for v := range ch {
	fmt.Println("recebido:", v)
}
```

---

## 4. Interfaces

### 4.1 Interfaces e polimorfismo (`interfaces/arq1.go`)
- **Ideia**: tipo interface com método `Area`, implementado por vários tipos.

```go
type Geometria interface {
	Area() float64
}

type Retangulo struct{ L, A float64 }
func (r Retangulo) Area() float64 { return r.L * r.A }

type Circulo struct{ Raio float64 }
func (c Circulo) Area() float64 { return math.Pi * c.Raio * c.Raio }

func Exibir(g Geometria) {
	fmt.Println(g.Area())
}
```

### 4.2 Implementando `error` (`interfaces/arq2.go`)
- **Ideia**: struct que implementa `Error()` para virar um `error` personalizado.

```go
type ErroRede struct {
	rede, hardware bool
}

func (e ErroRede) Error() string {
	if e.rede {
		return "Erro de rede"
	}
	if e.hardware {
		return "Erro de hardware"
	}
	return "Erro desconhecido"
}

var err error = ErroRede{rede: true}
fmt.Println(err.Error())
```

### 4.3 Interface vazia (`interfaces/arq3.go`)
- **Ideia**: `interface{}` aceita qualquer tipo (uso genérico).

```go
var lista []interface{}
lista = append(lista, 10, "Hello", true, 10.5)

for i, v := range lista {
	fmt.Println(i, v)
}
```

---

## 5. Structs e composição

### 5.1 Struct simples (`structs/endereco.go`)
- **Ideia**: agrupar campos relacionados em um tipo.

```go
type Endereco struct {
	Rua, Bairro, Cidade, Estado, CEP string
	Numero                           int
}
```

### 5.2 Struct com método e “herança” (`structs/pessoa.go`)
- **Ideia**: método com receptor, struct dentro de struct (composição).

```go
type Pessoa struct {
	Nome     string
	Idade    int
	Endereco Endereco
}

func (p *Pessoa) Aniversario() {
	p.Idade++
}

type DadosPessoais struct {
	Pessoa
	Cpf string
}
```

---

## 6. Pacotes externos (`pacotes/arq1.go`)

- **Ideia**: importar pacote de terceiros e usar funções.

```go
import "github.com/fatih/color"

func main() {
	color.Green("Exemplo de pacote")
	color.Red("Erro em vermelho")
}
```

---

## 7. Generics (`genericsts/arq1.go`)

- **Ideia**: função genérica que funciona para `int` e `string`.

```go
type Tipos interface {
	int | string
}

func reverse[T Tipos](slice []T) []T {
	n := len(slice)
	out := make([]T, n)
	for i := 0; i < n; i++ {
		out[i] = slice[n-1-i]
	}
	return out
}
```

---

## 8. Exercícios (`exercicios/ex1.go`)

- **Ideia**: exercício de modelar uma compra de mercado com structs e slice.

```go
type Item struct {
	Produto    string
	Quantidade int
	Preco      float64
}

type Compra struct {
	NomeMercado string
	Items       []Item
}
```

Use este README como índice rápido: encontre o tema, abra o arquivo correspondente e rode/edite os exemplos conforme precisar.


