package fundamentos

import "fmt"

func For() {
	j := []int{1, 2, 3, 4, 5}
	for i := 0; i < len(j); i++ {
		if i%2 == 0 {
			fmt.Println("Par")
		} else {
			fmt.Println("Impar")
		}
		fmt.Println(i)
	}

	for base := 1; base < 11; base++ {
		for i := 1; i <= 10; i++ {
			fmt.Printf("%d x %d = %d\n", base, i, base*i)
		}
	}

}
