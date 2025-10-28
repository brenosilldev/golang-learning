package fundamentos

import (
	"fmt"
)

func ConversaoDados() {
	// Conversao de dados
	x := 3.2
	y := 4.4
	var z float64 = float64(x) + float64(y)
	var total = int(z)
	fmt.Println(float32(x))
	fmt.Println(float64(y))
	fmt.Println(z)
	fmt.Println(total)
	fmt.Printf("%T", total) //Descobre  o tipo
}
