package avancado

import "fmt"

func ExemploRecover() {
	defer func() { // Defer = Executa a funcao apos a funcao principal ser executada
		if r := recover(); r != nil { // Recover = Recupera o erro e continua a execucao do programa
			fmt.Println("Recuperado:", r)
		}
	}() // Defer = Executa a funcao apos a funcao principal ser executada
	panic("Erro fatal") // Panic = Gera um erro fatal e para a execucao do programa
}
