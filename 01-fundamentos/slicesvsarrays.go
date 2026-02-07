package _1_fundamentos

import "fmt"

func SlicesVSArrays() {
	lista := []int{2, 30, 5, 4, 8, 6, 7, 9, 10, 3, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10} // slice - Dinamico
	array := [10]int{2, 30, 5, 4, 8, 6, 7, 9, 10, 3}                              // array - Fixo

	fmt.Println(lista)
	fmt.Println(array)

}
