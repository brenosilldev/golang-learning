package main

// ============================================================================
// EXERCÍCIO 7 — Structs e Métodos
// ============================================================================
//
// Exercício 7.1 — Sistema de Biblioteca
// Crie as structs:
//   - Livro: Titulo, Autor, AnoPublicacao, Disponivel (bool)
//   - Biblioteca: Nome, Livros ([]Livro)
//
// Implemente métodos para Biblioteca:
//   - AdicionarLivro(livro Livro)
//   - Emprestar(titulo string) error  — marca como indisponível
//   - Devolver(titulo string) error   — marca como disponível
//   - ListarDisponiveis() []Livro
//   - BuscarPorAutor(autor string) []Livro
//
// Exercício 7.2 — Composição
// Crie as structs usando embedding:
//   - Veiculo: Marca, Modelo, Ano
//   - Carro: Veiculo + NumPortas
//   - Moto: Veiculo + Cilindradas
//
// Cada tipo deve ter um método Descricao() string.
// Crie um slice com carros e motos e imprima todas as descrições.
//
// Exercício 7.3 — JSON
// Crie uma struct Usuario com: ID, Nome, Email, Senha, CriadoEm
// - Use tags JSON para que os campos fiquem em snake_case no JSON
// - Senha nunca deve aparecer no JSON (use tag "-")
// - CriadoEm deve ser omitido se vazio (omitempty)
// Serialize para JSON e deserialize de volta.
//
// ============================================================================

func main() {
	// TODO: Exercício 7.1 — Sistema de Biblioteca

	// TODO: Exercício 7.2 — Composição

	// TODO: Exercício 7.3 — JSON
}
