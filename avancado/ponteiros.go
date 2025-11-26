package avancado

import "fmt"

func Ponteiros() {
 
	x := 5 // valor da variavel
	y := 10 // valor da variavel

	p := &x // ponteiro para o endereço de x
	q := &y // ponteiro para o endereço de y

	fmt.Println("Ponteiros: ", p, q) // imprime o endereço de x e y // imprime o endereço dos ponteiros	
	fmt.Println("Valores: ", *p, *q) // imprime o valor de x e y // imprime o valor dos ponteiros

	*p = 10 // altera o valor de x para 10 // alterar o valor do ponteiro
	*q = 20 // altera o valor de y para 20 // alterar o valor do ponteiro

	fmt.Println("Ponteiros: ", &x, &y) // imprime o endereço de x e y // imprime o endereço dos ponteiros
	fmt.Println("Valores: ", x, y) // imprime o valor de x e y // imprime o valor dos ponteiros
}
