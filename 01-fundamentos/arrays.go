package _1_fundamentos

import "fmt"

func Arrays() {
	//lista := []int{1, 2, 3, 4, 5}
	//fmt.Println(lista)
	//
	//for i := 0; i < len(lista); i++ {
	//	fmt.Println(lista[i])
	//}
	//
	//fmt.Println(len(lista))

	lista := []int{2, 30, 5, 4, 8, 6, 7, 9, 10, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	listanova := make([]int, 0) // Make = Cria um array com 5 posições e 10 capacidade
	for i := 0; i < len(lista); i++ {
		if lista[i] < 10 {
			listanova = append(listanova, lista[i])

		}
	}

	novalista1 := lista[:5]          // 0 a 5 da lista -- [:] -> pega tudo / [5:] -> pega a partir da 5 / [2:5] -> pega a partir da 2 até a 5
	fmt.Println(lista[len(lista)-1]) // Pega o ultimo elemento da lista
	fmt.Println(listanova)
	fmt.Println(novalista1)
}
