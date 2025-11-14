package fundamentos

import "fmt"

func ExemploMaps() {
	maps := map[string]int{ // map[chave]valor
		"a": 1,
		"b": 2,
		"c": 3,
	}

	// maps := make(map[string]int) // map[chave]valor - Cria um map vazio

	fmt.Println(maps)
	fmt.Print("----")

	valor, existe := maps["d"] // valor,existe := maps["chave"] - Pega o valor da chave e verifica se existe
	if existe {
		fmt.Println("Chave encontrada: ", valor)
	} else {
		fmt.Println("Chave nao encontrada")
	}

	for chave, valor := range maps { //range = percorre o map
		if chave == "a" {
			fmt.Println("Chave encontrada: ", chave, valor)
		} else {
			fmt.Println("Chave nao encontrada", chave, valor)
		}
	}

	delete(maps, "a") // delete(maps, "chave") - Deleta a chave do map
	fmt.Println(maps)
}
