package interfaces

import (
	"fmt"
	"math"
)

// Geometria define o “contrato” que qualquer forma geométrica precisa cumprir.
// Basta implementar o método Area para ser considerada uma Geometria.
type Geometria interface {
	Area() float64
}

type Retangulo struct {
	largura, altura float64
}

type Circulo struct {
	raio float64
}

// Circulo implementa a interface fornecendo Area,
// logo pode ser tratado como uma Geometria.
func (c Circulo) Area() float64 {
	return math.Pi * c.raio * c.raio
}

// O mesmo vale para Retangulo: ao ter Area, ele satisfaz Geometria.
func (r Retangulo) Area() float64 {
	return r.largura * r.altura
}

// ExibirGeometria recebe qualquer tipo que atenda ao contrato Geometria.
// Assim conseguimos trabalhar com Circulo, Retangulo ou outros tipos que
// implementem Area de forma polimórfica.
func ExibirGeometria(g Geometria) {
	fmt.Println(g.Area())
}

func Arq1() {
	fmt.Println("----------- Inicio do programa -----------")

	retangulo := Retangulo{largura: 10, altura: 20} // Retangulo implementa Geometria
	circulo := Circulo{raio: 10}

	ExibirGeometria(retangulo)
	ExibirGeometria(circulo)

}
