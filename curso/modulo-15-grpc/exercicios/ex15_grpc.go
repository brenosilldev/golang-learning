package main

// ============================================================================
// EXERCÍCIO 15 — gRPC
// ============================================================================
//
// Exercício 15.1 — Serviço de Usuários
// Crie um .proto para um serviço de usuários com:
//   - CreateUser(name, email) → User
//   - GetUser(id) → User
//   - ListUsers() → [User]
//   - UpdateUser(id, name, email) → User
//   - DeleteUser(id) → Success
//
// Gere o código, implemente servidor e cliente.
//
// Exercício 15.2 — Server Streaming (Chat)
// Crie um serviço de chat simples:
//   - SendMessage(user, content) → MessageResponse
//   - StreamMessages(room) → stream Message
//
// O servidor faz stream de todas as mensagens de um "room" para os clients.
// Use goroutines para enviar mensagens em background.
//
// Exercício 15.3 — Bidirectional Streaming
// Crie um serviço de chat com bidi streaming:
//   - Chat(stream Message) → stream Message
//
// Cada mensagem enviada pelo client aparece em todos os outros clients.
// Dica: mantenha uma lista de streams ativos no servidor.
//
// Exercício 15.4 — Interceptors (Middleware gRPC)
// Implemente interceptors para:
//   a) Logging: logar cada chamada RPC com duração
//   b) Auth: verificar metadata "authorization" com um token
//   c) Recovery: capturar panics e retornar erro gRPC
// Dica: use grpc.UnaryInterceptor() e grpc.StreamInterceptor()
//
// ============================================================================

func main() {
	// TODO: Exercício 15.1 — Serviço de Usuários

	// TODO: Exercício 15.2 — Server Streaming Chat

	// TODO: Exercício 15.3 — Bidirectional Streaming

	// TODO: Exercício 15.4 — Interceptors
}
