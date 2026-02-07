package _1_fundamentos

import "fmt"

func Estruturas() {

	// && - E - todos os valores devem ser true
	// || - OU - pelo menos um valor deve ser true
	// ! - NOT - inverte o valor

	salario := 1100.0
	desconto := 0.1

	if salario <= 1000.0 {
		bonus := salario * 0.1

		fmt.Println("Salarios com bonus: ", salario+bonus)
	} else if salario > 1000.0 && salario <= 5000.0 {
		bonus := salario * 0.20

		fmt.Println("Salarios com bonus: ", salario+bonus)

	} else {
		salarioDescontado := salario - salario*(1-desconto)

		fmt.Println("Salarios com desconto: ", salario-salarioDescontado)
	}
}
