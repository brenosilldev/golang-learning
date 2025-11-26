package gorountines

import (
	"fmt"
	"sync"
	"time"
)

/*
	var wg sync.Mutex // Mutex = Mutex = Garante que apenas uma goroutine possa acessar o recurso compartilhado
	wg.Lock() // Lock = Bloqueia o recurso compartilhado
	wg.Unlock() // Unlock = Desbloqueia o recurso compartilhado
	wg.TryLock() // TryLock = Tenta bloquear o recurso compartilhado
	wg.TryUnlock() // TryUnlock = Tenta desbloqueiar o recurso compartilhado
	wg.RLock() // RLock = Bloqueia o recurso compartilhado para leitura
	wg.RUnlock() // RUnlock = Desbloqueia o recurso compartilhado para leitura
	wg.RLock() // RLock = Bloqueia o recurso compartilhado para leitura


*/

func ExemploGoroutine3() {
	var mutex sync.Mutex // Mutex = Mutex = Garante que apenas uma goroutine possa acessar o recurso compartilhado
	i := 0

	for x := 0; x < 100; x++ {
		go func() { // Cria uma goroutine para executar a funcao
			mutex.Lock()   // Bloqueia o recurso compartilhado
			i++            // Incrementa o valor de i
			mutex.Unlock() // Desbloqueia o recurso compartilhado
		}() // Cria uma goroutine para executar a funcao
	}

	time.Sleep(time.Second * 1) // Espera 1 segundo
	fmt.Println("Valor final de i: ", i)

}
