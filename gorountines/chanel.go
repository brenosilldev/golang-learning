package gorountines

import (
	"fmt"
	"time"
)

/*
	Channel = Canal para enviar e receber dados entre goroutines
	chan <- = Envia o valor para o canal
	<-chan = Recebe o valor do canal
	channel := make(chan int, 10) // Cria um canal para enviar e receber dados com buffer de 10
*/

func ExemploChanel() {
	channel := make(chan int, 10) // Cria um canal para enviar e receber dados com buffer de 10
	go setList(channel)           // Cria uma goroutine para executar a funcao setList
	fmt.Println(<-channel)        // Recebe o valor do canal e imprime na tela
	for i := range channel {
		fmt.Println("Valor recebido: ", i)
		time.Sleep(time.Second)

	}
}

func setList(channel chan<- int) { // Funcao que envia os valores para o canal
	for i := 0; i < 10; i++ { // Envia os valores para o canal
		channel <- i // Envia o valor para o canal
		fmt.Println("Valor enviado: ", i)
	}
	close(channel) // Fecha o canal
}
