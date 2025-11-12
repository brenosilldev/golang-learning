package fundamentos

import "fmt"

func Arrays() {
	lista := []int{1, 2, 3, 4, 5}
	fmt.Println(lista)

	for i := 0; i < len(lista); i++ {
		fmt.Println(lista[i])
	}

}
