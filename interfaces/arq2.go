package interfaces

import (
	"errors"
	"fmt"
)

type ErroRede struct {
	rede     bool
	hardware bool
}

func (e ErroRede) Error() string {
	if e.rede {
		return "Erro de rede"
	} else if e.hardware {
		return "Erro de hardware"
	} else {
		return "Erro desconhecido"
	}
}

func ExibirErro(e error) {
	fmt.Println(e.Error())
}

func Arq2() {
	fmt.Println("----------- Inicio do programa -----------")

	erro := ErroRede{rede: true, hardware: false}
	ExibirErro(errors.New("Erro de rede 1,/"))
	ExibirErro(erro)
}
