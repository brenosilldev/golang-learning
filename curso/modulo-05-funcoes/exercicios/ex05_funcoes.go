package main

// ============================================================================
// EXERCÍCIO 5 — Funções
// ============================================================================
//
// Exercício 5.1 — Calculadora com Múltiplos Retornos
// Crie uma função `calculadora(a, b float64, op string) (float64, error)`
// que aceita operações: "+", "-", "*", "/".
// Retorne erro se a operação for inválida ou divisão por zero.
// Teste com várias operações.
//
// Exercício 5.2 — Map / Filter / Reduce
// Implemente 3 funções usando funções como parâmetro:
//   a) mapSlice([]int, func(int) int) []int — transforma cada elemento
//   b) filterSlice([]int, func(int) bool) []int — filtra elementos
//   c) reduceSlice([]int, func(int, int) int, int) int — reduz a um valor
// Teste: dado []int{1,2,3,4,5,6,7,8,9,10}
//   - mapSlice: multiplique por 3
//   - filterSlice: apenas pares
//   - reduceSlice: soma total
//
// Exercício 5.3 — Gerador de Fibonacci com Closure
// Crie uma closure que retorne func() int.
// Cada chamada deve retornar o próximo número de Fibonacci.
// fibonacci() → 0, 1, 1, 2, 3, 5, 8, 13, 21...
//
// Exercício 5.4 — Pipeline de Texto
// Crie um tipo `type TransformFunc func(string) string`
// Implemente um pipeline que aplique N transformações em sequência.
// Transformações: trim espaços, minúsculas, substituir espaços por "-",
//                 adicionar prefixo "slug-"
// Entrada: "  Meu Título de Blog  "
// Saída:   "slug-meu-título-de-blog"
//
// ============================================================================

func main() {
	// TODO: Exercício 5.1 — Calculadora com Múltiplos Retornos

	// TODO: Exercício 5.2 — Map / Filter / Reduce

	// TODO: Exercício 5.3 — Gerador de Fibonacci com Closure

	// TODO: Exercício 5.4 — Pipeline de Texto
}
