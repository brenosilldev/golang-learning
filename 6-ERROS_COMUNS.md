# ⚠️ Erros Comuns em Go — FAQ do Iniciante

> Os erros que TODO mundo comete ao aprender Go. Leia antes de sofrer.

---

## 1. `declared and not used`

```go
// ❌ ERRO
func main() {
    x := 42 // declared and not used
}
```

Go **não permite** variáveis sem uso. Ou use, ou remova.

```go
// ✅ CORRETO
func main() {
    x := 42
    fmt.Println(x)
    
    // Se precisa ignorar temporariamente:
    _ = 42  // blank identifier
}
```

> **Por quê?** Go prioriza código limpo. Variável não usada = provável bug.

---

## 2. Confundir `:=` com `=`

```go
// ❌ ERRO — := cria NOVA variável (shadow!)
x := 10
if true {
    x := 20  // NOVA variável x dentro do if (shadow)
}
fmt.Println(x) // imprime 10, não 20!
```

```go
// ✅ CORRETO — = reatribui a MESMA variável
x := 10
if true {
    x = 20  // reatribui o x externo
}
fmt.Println(x) // imprime 20
```

> **Regra**: `:=` = criar nova variável. `=` = modificar existente. Dentro de blocos internos, use `=` se quer mudar a variável externa.

---

## 3. Goroutine que termina antes de executar

```go
// ❌ ERRO — programa termina antes da goroutine rodar
func main() {
    go fmt.Println("olá")  // pode nunca imprimir!
}
```

```go
// ✅ CORRETO — esperar com WaitGroup
func main() {
    var wg sync.WaitGroup
    wg.Add(1)
    go func() {
        defer wg.Done()
        fmt.Println("olá")
    }()
    wg.Wait() // espera terminar
}

// ✅ CORRETO — esperar com channel
func main() {
    done := make(chan bool)
    go func() {
        fmt.Println("olá")
        done <- true
    }()
    <-done
}
```

> **Por quê?** `main()` termina = programa todo morre, incluindo goroutines.

---

## 4. Nil pointer dereference

```go
// ❌ ERRO — acessar campo de ponteiro nil
var p *Pessoa   // p é nil
fmt.Println(p.Nome)  // PANIC: nil pointer dereference
```

```go
// ✅ CORRETO — verificar antes
var p *Pessoa
if p != nil {
    fmt.Println(p.Nome)
}

// ✅ Ou inicializar
p := &Pessoa{Nome: "Alice"}
fmt.Println(p.Nome)
```

> **Onde isso acontece mais**: retorno de funções que retornam ponteiro + error. Sempre verifique o error ANTES de usar o ponteiro.

---

## 5. Slice append gotcha

```go
// ❌ SURPRESA — append pode ou NÃO modificar o slice original
s := []int{1, 2, 3}
s2 := s[:2]      // s2 = [1, 2], compartilha memória com s
s2 = append(s2, 99)
fmt.Println(s)    // [1, 2, 99] ← s FOI modificado!!
```

```go
// ✅ SEGURO — usar copy para independência
s := []int{1, 2, 3}
s2 := make([]int, 2)
copy(s2, s[:2])        // cópia independente
s2 = append(s2, 99)
fmt.Println(s)          // [1, 2, 3] ← s não mudou
```

> **Por quê?** Slices compartilham o array subjacente. `append` só aloca novo array quando a capacidade é ultrapassada.

---

## 6. Loop com goroutine capturando variável errada

```go
// ❌ ERRO CLÁSSICO — todas as goroutines veem o ÚLTIMO valor
for _, v := range []string{"a", "b", "c"} {
    go func() {
        fmt.Println(v)  // imprime "c" "c" "c"
    }()
}
```

```go
// ✅ CORRETO (Go 1.22+) — loop var tem escopo por iteração
// Em Go 1.22+, o código acima funciona corretamente!

// ✅ CORRETO (Go < 1.22) — passar como parâmetro
for _, v := range []string{"a", "b", "c"} {
    go func(val string) {
        fmt.Println(val)  // imprime "a" "b" "c"
    }(v)
}
```

> **Nota**: Go 1.22 corrigiu isso! Mas muitos projetos ainda rodam versões anteriores.

---

## 7. Map não é thread-safe

```go
// ❌ ERRO — race condition! Crash aleatório
m := make(map[string]int)
for i := 0; i < 100; i++ {
    go func(n int) {
        m[fmt.Sprint(n)] = n  // PANIC: concurrent map writes
    }(i)
}
```

```go
// ✅ CORRETO — usar sync.Mutex
var mu sync.Mutex
m := make(map[string]int)
for i := 0; i < 100; i++ {
    go func(n int) {
        mu.Lock()
        m[fmt.Sprint(n)] = n
        mu.Unlock()
    }(i)
}

// ✅ Ou usar sync.Map (para casos específicos)
var m sync.Map
m.Store("chave", "valor")
v, _ := m.Load("chave")
```

---

## 8. Esquecer de tratar error

```go
// ❌ ERRO — ignorar erro
arquivo, _ := os.Open("dados.txt")  // se falhar, arquivo é nil
dados := arquivo.Read(buf)           // PANIC!
```

```go
// ✅ CORRETO — SEMPRE verificar
arquivo, err := os.Open("dados.txt")
if err != nil {
    log.Fatal("erro ao abrir:", err)
}
defer arquivo.Close()
```

> **Regra absoluta em Go**: NUNCA ignore um `error`. Se realmente não importa, documente por quê: `_ = arquivo.Close() // erro de close é inofensivo aqui`

---

## 9. Defer em loop

```go
// ❌ PERIGOSO — defer só executa ao SAIR da função, não do loop
for _, nome := range arquivos {
    f, _ := os.Open(nome)
    defer f.Close()  // TODOS os arquivos ficam abertos até o fim da função!
}
```

```go
// ✅ CORRETO — usar função auxiliar
for _, nome := range arquivos {
    func() {
        f, _ := os.Open(nome)
        defer f.Close()  // fecha ao sair da função anônima
        // processar f
    }()
}
```

---

## 10. String é imutável

```go
// ❌ ERRO — não pode modificar caractere individual
s := "hello"
s[0] = 'H'  // ERRO DE COMPILAÇÃO

// ✅ CORRETO — converter para []byte ou []rune
s := "hello"
b := []byte(s)
b[0] = 'H'
s = string(b)  // "Hello"
```

---

## 11. Range retorna cópia

```go
// ❌ SURPRESA — range cria CÓPIA do elemento
type Pessoa struct{ Nome string; Idade int }
pessoas := []Pessoa{{Nome: "Alice", Idade: 30}}

for _, p := range pessoas {
    p.Idade = 99  // modifica a CÓPIA, não o original!
}
fmt.Println(pessoas[0].Idade)  // ainda 30
```

```go
// ✅ CORRETO — usar índice
for i := range pessoas {
    pessoas[i].Idade = 99  // modifica o original
}
```

---

## 12. Interface nil vs valor nil

```go
// ❌ SURPRESA — interface com tipo mas valor nil NÃO é nil
var p *Pessoa = nil
var i interface{} = p

fmt.Println(p == nil)  // true
fmt.Println(i == nil)  // FALSE!! 😱
```

> **Por quê?** Uma interface tem dois campos: `(tipo, valor)`. Quando `i = p`, o tipo é `*Pessoa` e o valor é `nil`. A interface só é nil quando **ambos** são nil.

```go
// ✅ Como verificar
if i == nil || reflect.ValueOf(i).IsNil() {
    fmt.Println("é nil de verdade")
}
```

---

## 13. Imports circulares

```go
// ❌ ERRO — pacote A importa B, e B importa A
// package a → import "b"
// package b → import "a"
// ERRO: import cycle not allowed
```

```go
// ✅ SOLUÇÃO — extrair interface para pacote separado
// package models → define interfaces
// package a → importa models
// package b → importa models
```

---

## 🧠 Resumo Rápido

| # | Erro | Solução em 1 frase |
|---|------|-------------------|
| 1 | declared and not used | Use ou delete a variável |
| 2 | := vs = | := cria nova, = reatribui existente |
| 3 | Goroutine não executa | Use WaitGroup ou channel |
| 4 | Nil pointer | Verifique error antes de usar ponteiro |
| 5 | Slice append | Use copy() para independência |
| 6 | Loop + goroutine | Passe variável como parâmetro (ou use Go 1.22+) |
| 7 | Map + goroutines | Use sync.Mutex |
| 8 | Error ignorado | SEMPRE trate com if err != nil |
| 9 | Defer em loop | Encapsule em função anônima |
| 10 | String imutável | Converta para []byte |
| 11 | Range copia | Use índice para modificar |
| 12 | Interface nil | Cuidado com (tipo, nil) vs nil |
| 13 | Import circular | Extraia interfaces para pacote intermediário |
