package fundamentos

import "fmt"

func ExercicioMaps() {
	numeros := []int{1, 20} // array de numeros

	soma := 0 // soma dos numeros
	for i := 0; i < len(numeros); i++ {
		soma += numeros[i]
	}
	fmt.Println(soma)

	lista := []int{2,8,3,10,5,4,9,9}
	numeroate5 := 0
	numeroacima5 := 0
	for i := 0; i < len(lista); i++ {
		if lista[i] <= 5 {
			numeroate5 += lista[i] // numeroate5 = numeroate5 + lista[i]
		} else {
			numeroacima5 += lista[i]
		}
	}
	fmt.Println(numeroate5)
	fmt.Println(numeroacima5)
}
