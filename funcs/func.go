package funcs

import "fmt"

func ExemploFunc() {
	fmt.Println("Exemplo de funcao")
	resultado, mensagem, valor := exemplo(10, 20)
	fmt.Println(resultado)
	fmt.Println(mensagem)
	fmt.Println(valor)
}


func exemplo(a int, b int) (int,string, float64) {
	return a + b, "Hello, World!", float64(a + b)
}