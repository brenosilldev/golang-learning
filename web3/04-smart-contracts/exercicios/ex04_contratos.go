package main

// ============================================================================
// EXERCÍCIOS — Smart Contracts
// ============================================================================
//
// Exercício 4.1 — Deploy e Interação
// Usando o Cofre.sol de exemplo:
//   a) Compile com solc e gere bindings com abigen
//   b) Faça deploy em blockchain local (Ganache)
//   c) Deposite 2 ETH
//   d) Consulte o saldo
//   e) Saque 1 ETH
//   f) Verifique que o saldo atualizou
//
// Exercício 4.2 — Event Listener
// Implemente um listener que:
//   - Escuta eventos de Deposito e Saque em tempo real
//   - Para cada evento, imprime: timestamp, endereço, valor em ETH
//   - Mantém um log em arquivo (JSON lines)
//   - Use goroutine para escutar sem bloquear
//
// Exercício 4.3 — Seu Próprio Contrato (Votação)
// Escreva um contrato de votação em Solidity:
//   - Dono adiciona candidatos
//   - Qualquer endereço vota UMA VEZ
//   - Função para ver resultado
//   - Eventos: NovoCandidato, Voto
// Gere bindings Go e faça: deploy, adicionar candidatos, votar, ver resultado.
//
// Exercício 4.4 — Contract Reader (sem abigen)
// Interaja com um contrato SEM usar abigen:
//   - Use abi.JSON() para carregar o ABI
//   - Use abi.Pack() para encodar chamadas
//   - Use client.CallContract() para executar
//   - Decodifique o resultado manualmente
// Isso é importante para entender o que abigen faz por baixo!
//
// ============================================================================

func main() {
	// TODO: Exercício 4.1 — Deploy e Interação

	// TODO: Exercício 4.2 — Event Listener

	// TODO: Exercício 4.3 — Contrato de Votação

	// TODO: Exercício 4.4 — Contract Reader manual
}
