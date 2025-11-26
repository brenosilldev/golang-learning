package avancado

import (
	"fmt"
	"os"
)

/*
*
Defer = Executa a funcao apos a funcao principal ser executada
defer file.Close() // Fecha o arquivo apos a funcao principal ser executada
defer file.Write([]byte("Hello, World!")) // Escreve no arquivo apos a funcao principal ser executada
defer file.WriteString("Hello, World!") // Escreve no arquivo apos a funcao principal ser executada
*/
func ExemploDefer() {
	file, err := os.Create("arquivo.txt")
	if err != nil {
		panic(err)
	}
	file.Write([]byte("Hello, World!"))
	defer file.Close() // Defer = Executa a funcao apos a funcao principal ser executada

	fmt.Println("Arquivo criado com sucesso!")
}
