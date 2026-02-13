package main

// ============================================================================
// EXERCÍCIO 11 — Concorrência
// ============================================================================
//
// Exercício 11.1 — Download Simulado
// Simule o download de 5 arquivos em paralelo usando goroutines.
// Cada "download" é um time.Sleep aleatório (100-500ms).
// Use WaitGroup para esperar todos terminarem.
// Imprima o tempo total (deve ser ~500ms, não ~1500ms).
//
// Exercício 11.2 — Fan-In (merge channels)
// Crie 3 goroutines, cada uma enviando números por um channel diferente.
// Crie uma função `fanIn` que recebe N channels e retorna um único channel
// que emite todos os valores. Imprima os valores conforme chegam.
//
// Exercício 11.3 — Worker Pool
// Crie um worker pool com 3 workers que processam URLs (strings).
// Cada worker "baixa" a URL (sleep aleatório) e retorna o tamanho simulado.
// Use channels para jobs e resultados.
// Processe 10 URLs e imprima os resultados.
//
// Exercício 11.4 — Rate Limiter
// Crie um rate limiter usando time.Ticker que permite no máximo
// 3 requisições por segundo. Simule 10 requisições e mostre
// que elas são espaçadas corretamente.
//
// Exercício 11.5 — Context com Timeout
// Crie uma função que "busca dados" (sleep de 1-3 segundos aleatório).
// Use context.WithTimeout para dar no máximo 2 segundos.
// Se der timeout, imprima "timeout". Se completar, imprima o resultado.
// Execute 5 vezes e veja quantas completam vs timeout.
//
// ============================================================================

func main() {
	// TODO: Exercício 11.1 — Download Simulado

	// TODO: Exercício 11.2 — Fan-In

	// TODO: Exercício 11.3 — Worker Pool

	// TODO: Exercício 11.4 — Rate Limiter

	// TODO: Exercício 11.5 — Context com Timeout
}
