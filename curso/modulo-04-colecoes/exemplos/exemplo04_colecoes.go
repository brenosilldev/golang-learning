package main

import "fmt"

// ============================================================================
// MÓDULO 04 — Coleções: Arrays, Slices e Maps
// ============================================================================
//
// Conceitos:
//   - Arrays (tamanho fixo, raramente usados)
//   - Slices (dinâmicos — o que você vai usar 99% do tempo)
//   - len vs cap, append, copy, reslicing
//   - Maps (chave/valor)
//   - Sets (map[T]struct{})
//
// Rode com: go run exemplo04_colecoes.go
// ============================================================================

func main() {
	// ==========================================================
	fmt.Println("=== ARRAYS (tamanho fixo) ===")

	// Array: tamanho faz parte do tipo! [3]int ≠ [5]int
	var notas [5]float64
	notas[0] = 9.5
	notas[1] = 8.0
	notas[2] = 7.5
	notas[3] = 9.0
	notas[4] = 8.5
	fmt.Println("Notas:", notas)
	fmt.Printf("Tipo: %T, Len: %d\n", notas, len(notas))

	// Array literal
	cores := [3]string{"vermelho", "verde", "azul"}
	fmt.Println("Cores:", cores)

	// [...] deixa o compilador contar
	nums := [...]int{10, 20, 30, 40}
	fmt.Printf("Nums: %v (tamanho: %d)\n", nums, len(nums))

	// ⚠️ Arrays são COPIADOS por valor!
	original := [3]int{1, 2, 3}
	copia := original
	copia[0] = 999
	fmt.Printf("Original: %v, Cópia: %v\n", original, copia) // original NÃO muda

	// ==========================================================
	fmt.Println("\n=== SLICES (dinâmicos) ===")

	// Slice literal (sem tamanho nos colchetes)
	frutas := []string{"🍎", "🍌", "🍊", "🍇", "🍉"}
	fmt.Println("Frutas:", frutas)
	fmt.Printf("Len: %d, Cap: %d\n", len(frutas), cap(frutas))

	// Criando com make(tipo, length, capacity)
	buffer := make([]int, 3, 10)
	fmt.Printf("Buffer: %v (len=%d, cap=%d)\n", buffer, len(buffer), cap(buffer))

	// Append — adicionar elementos
	buffer = append(buffer, 1, 2, 3)
	fmt.Printf("Após append: %v (len=%d, cap=%d)\n", buffer, len(buffer), cap(buffer))

	// Append de outro slice (spread com ...)
	extras := []int{7, 8, 9}
	buffer = append(buffer, extras...)
	fmt.Printf("Após append slice: %v\n", buffer)

	// ==========================================================
	fmt.Println("\n=== SLICING (fatiamento) ===")

	numeros := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

	fmt.Println("Original:", numeros)
	fmt.Println("[2:5]  =", numeros[2:5]) // [2, 3, 4]
	fmt.Println("[:3]   =", numeros[:3])  // [0, 1, 2]
	fmt.Println("[7:]   =", numeros[7:])  // [7, 8, 9]
	fmt.Println("[:]    =", numeros[:])   // tudo (cópia rasa!)

	// ⚠️ CUIDADO: slices compartilham o array subjacente!
	a := []int{1, 2, 3, 4, 5}
	b := a[1:4] // [2, 3, 4]
	b[0] = 999
	fmt.Printf("a: %v, b: %v\n", a, b) // a[1] também virou 999!

	// Para fazer uma cópia independente, use copy()
	original2 := []int{1, 2, 3}
	copia2 := make([]int, len(original2))
	copy(copia2, original2)
	copia2[0] = 999
	fmt.Printf("Original: %v, Cópia independente: %v\n", original2, copia2)

	// ==========================================================
	fmt.Println("\n=== FILTRAGEM DE SLICE ===")

	valores := []int{15, 3, 42, 7, 88, 1, 23, 56}
	maioresQue10 := make([]int, 0)

	for _, v := range valores {
		if v > 10 {
			maioresQue10 = append(maioresQue10, v)
		}
	}
	fmt.Println("Valores > 10:", maioresQue10)

	// ==========================================================
	fmt.Println("\n=== REMOVENDO ELEMENTO DE SLICE ===")

	// Remover elemento no índice 2
	lista := []string{"a", "b", "c", "d", "e"}
	indice := 2
	lista = append(lista[:indice], lista[indice+1:]...)
	fmt.Println("Após remover índice 2:", lista) // [a b d e]

	// ==========================================================
	fmt.Println("\n=== MAPS ===")

	// Map literal
	capitais := map[string]string{
		"Brasil":    "Brasília",
		"Argentina": "Buenos Aires",
		"Japão":     "Tóquio",
	}
	fmt.Println("Capitais:", capitais)

	// Criando com make
	estoque := make(map[string]int)
	estoque["camiseta"] = 50
	estoque["calça"] = 30
	estoque["tênis"] = 20
	fmt.Println("Estoque:", estoque)

	// Leitura SEGURA com ok idiom
	valor, existe := capitais["Brasil"]
	if existe {
		fmt.Printf("Capital do Brasil: %s\n", valor)
	}

	valor2, existe2 := capitais["Alemanha"]
	if !existe2 {
		fmt.Printf("Alemanha não encontrada (valor zero: '%s')\n", valor2)
	}

	// Iteração
	fmt.Println("\nTodas as capitais:")
	for pais, capital := range capitais {
		fmt.Printf("  %s → %s\n", pais, capital)
	}

	// Delete
	delete(capitais, "Argentina")
	fmt.Println("Após delete:", capitais)

	// Tamanho
	fmt.Printf("Total de capitais: %d\n", len(capitais))

	// ==========================================================
	fmt.Println("\n=== MAP COMO CONTADOR ===")

	frase := "banana abacaxi banana manga banana manga"
	contador := make(map[string]int)

	// Usando um mini-parser manual
	palavra := ""
	for _, ch := range frase + " " {
		if ch == ' ' {
			if palavra != "" {
				contador[palavra]++
				palavra = ""
			}
		} else {
			palavra += string(ch)
		}
	}

	fmt.Println("Contagem:")
	for palavra, count := range contador {
		fmt.Printf("  %s: %d\n", palavra, count)
	}

	// ==========================================================
	fmt.Println("\n=== SET (simulado com map[T]struct{}) ===")

	// Go não tem Set nativo. A forma idiomática é usar map[T]struct{}
	// struct{} ocupa 0 bytes de memória!
	set := make(map[string]struct{})

	// Adicionar
	set["go"] = struct{}{}
	set["rust"] = struct{}{}
	set["python"] = struct{}{}
	set["go"] = struct{}{} // duplicata — nada acontece

	// Verificar existência
	if _, ok := set["go"]; ok {
		fmt.Println("'go' está no set")
	}

	if _, ok := set["java"]; !ok {
		fmt.Println("'java' NÃO está no set")
	}

	// Remover
	delete(set, "python")
	fmt.Printf("Set (tamanho %d): ", len(set))
	for item := range set {
		fmt.Printf("%s ", item)
	}
	fmt.Println()
}
