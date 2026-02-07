package _1_fundamentos

import "fmt"

var (
	a int = -1
	b int = 2
)

func Teste() {
	//Variaveis sao tipos dinamicos
	// O tipo nao pode ser alterado
	// Const nao pode ser alterada
	a += -2
	b += 1
	metro := 2.0

	fmt.Println(a, b, metro)
	fmt.Println(a + b)
	fmt.Println(float32(metro + 0.2*float64(a+b)))
}
