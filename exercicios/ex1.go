package exercicios

import (
	"fmt"
	"time"
)

type ItemsCarrinho struct {
	Produto    string
	Quantidade int
	Preco      float64
}

type CompraMercado struct {
	DataCompra  time.Time
	NomeMercado string
	Items       []ItemsCarrinho
}

func Exercicio1() {

	items := []ItemsCarrinho{}
	items = append(items, ItemsCarrinho{Produto: "Arroz", Quantidade: 2, Preco: 10.0})
	items = append(items, ItemsCarrinho{Produto: "Feijao", Quantidade: 1, Preco: 5.0})

	compra := CompraMercado{
		DataCompra:  time.Now(),
		Items:       items,
		NomeMercado: "Mercado do Joao",
	}

	fmt.Println(items)
	fmt.Println(compra)
}

func Exercicio2() {

}
