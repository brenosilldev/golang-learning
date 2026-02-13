package main

// ============================================================================
// EXERCÍCIOS — Criptografia
// ============================================================================
//
// Exercício 1.1 — Hash Comparator
// Crie uma função que recebe dois arquivos (strings simulando conteúdo)
// e verifica se são idênticos usando SHA-256.
// Teste com strings iguais e diferentes.
// Dica: se os hashes são iguais, os conteúdos são iguais.
//
// Exercício 1.2 — Password Manager Simples
// Crie um gerenciador de senhas em memória:
//   - Registrar(usuario, senha) → armazena hash da senha (NÃO a senha!)
//   - Verificar(usuario, senha) → true se o hash bate
//   - Listar() → mostra usuários (sem senhas!)
// Adicione "salt" (string aleatória concatenada antes do hash) para segurança extra.
//
// Exercício 1.3 — Proof of Work Variável
// Implemente um miner que aceita dificuldade variável (1 a 6 zeros).
// Para cada dificuldade, meça:
//   - Tempo até encontrar o nonce
//   - Número de tentativas
// Monte uma tabela mostrando como a dificuldade cresce exponencialmente.
//
// Exercício 1.4 — Multi-Signature Wallet
// Crie uma wallet que requer 2 de 3 assinaturas para validar uma transação:
//   - Crie 3 wallets (Alice, Bob, Carol)
//   - Uma transação precisa de pelo menos 2 assinaturas para ser válida
//   - Implemente: Assinar(tx, wallet) e Verificar(tx, assinaturas[], chaves[])
//   - Teste: 2 de 3 = válido, 1 de 3 = inválido
//
// Exercício 1.5 — Merkle Tree com Prova
// Use a MerkleTree do exemplo e:
//   a) Construa uma árvore com 8 transações
//   b) Gere prova de inclusão para a transação 5
//   c) Verifique que a prova é válida
//   d) Tente verificar com uma transação adulterada
//   e) Adicione uma 9ª transação e mostre que o Root muda
//
// ============================================================================

func main() {
	// TODO: Exercício 1.1 — Hash Comparator

	// TODO: Exercício 1.2 — Password Manager

	// TODO: Exercício 1.3 — Proof of Work Variável

	// TODO: Exercício 1.4 — Multi-Signature Wallet

	// TODO: Exercício 1.5 — Merkle Tree com Prova
}
