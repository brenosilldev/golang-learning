package structs

import (
	"fmt"
	"time"
)

type Pessoa struct {
	Nome           string
	Idade          int
	Email          string
	DataNascimento time.Time
	Endereco       Endereco
}

type DadosPessoais struct {
	Pessoa // Herança de Pessoa
	Cpf    string
	RG     string
}

func (p *Pessoa) IdadePessoa() int {

	p.Idade = CalcularIdade(p.DataNascimento)
	return p.Idade
}

func CalcularIdade(dataNascimento time.Time) int {
	anoNascimento := dataNascimento.Year()
	anoAtual := time.Now().Year()
	return anoAtual - anoNascimento
}

func ExemploStruct() {

	pessoa := Pessoa{
		Nome:  "Joao",
		Idade: 20,
		Email: "joao@gmail.com",
		Endereco: Endereco{
			Rua:    "Rua das Flores",
			Numero: 123,
			Bairro: "Centro",
			Cidade: "Sao Paulo",
			Estado: "SP",
			CEP:    "1234567890",
		},
		DataNascimento: time.Date(1990, 1, 1, 0, 0, 0, 0, time.UTC),
	}
	pessoa.Idade = pessoa.IdadePessoa()

	dadosPessoais := DadosPessoais{
		Pessoa: pessoa,
		Cpf:    "1234567890",
		RG:     "1234567890",
	}

	fmt.Println(dadosPessoais)

	fmt.Println(pessoa.Nome)
	fmt.Println("Idade: ", pessoa.Idade)
}
