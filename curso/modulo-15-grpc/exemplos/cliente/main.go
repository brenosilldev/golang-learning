package main

import (
	"fmt"
)

// ============================================================================
// MÓDULO 15 — Cliente gRPC
// ============================================================================
//
// O client é gerado pelo protoc junto com o server.
// Você só precisa criar a conexão e chamar os métodos.
//
// Para rodar:
//   1. Inicie o servidor primeiro: go run ../servidor/main.go
//   2. Descomente o código abaixo
//   3. go run main.go
// ============================================================================

// import (
// 	"context"
// 	"io"
// 	"log"
// 	"time"
//
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/credentials/insecure"
// 	pb "seu-modulo/curso/modulo-15-grpc/exemplos/pb"
// )
//
// func main() {
// 	// Conectar ao servidor
// 	conn, err := grpc.Dial("localhost:50051",
// 		grpc.WithTransportCredentials(insecure.NewCredentials()),
// 	)
// 	if err != nil {
// 		log.Fatalf("Falha ao conectar: %v", err)
// 	}
// 	defer conn.Close()
//
// 	// Criar client (gerado pelo protoc)
// 	client := pb.NewTaskServiceClient(conn)
// 	ctx := context.Background()
//
// 	// --- Criar tarefas ---
// 	fmt.Println("=== Criando tarefas ===")
// 	tarefas := []string{"Aprender gRPC", "Estudar Protobuf", "Montar API"}
// 	for _, titulo := range tarefas {
// 		resp, err := client.CreateTask(ctx, &pb.CreateTaskRequest{
// 			Titulo:    titulo,
// 			Descricao: "Criada via gRPC client",
// 		})
// 		if err != nil {
// 			log.Printf("Erro ao criar: %v", err)
// 			continue
// 		}
// 		fmt.Printf("Criada: ID=%d, Título=%s\n", resp.Task.Id, resp.Task.Titulo)
// 	}
//
// 	// --- Listar todas ---
// 	fmt.Println("\n=== Listando tarefas ===")
// 	list, err := client.ListTasks(ctx, &pb.Empty{})
// 	if err != nil {
// 		log.Fatalf("Erro ao listar: %v", err)
// 	}
// 	for _, t := range list.Tasks {
// 		fmt.Printf("  [%d] %s (concluída: %t)\n", t.Id, t.Titulo, t.Concluida)
// 	}
//
// 	// --- Buscar por ID ---
// 	fmt.Println("\n=== Buscando tarefa 1 ===")
// 	resp, err := client.GetTask(ctx, &pb.GetTaskRequest{Id: 1})
// 	if err != nil {
// 		log.Printf("Erro: %v", err)
// 	} else {
// 		fmt.Printf("  Encontrada: %s\n", resp.Task.Titulo)
// 	}
//
// 	// --- Server Streaming ---
// 	fmt.Println("\n=== Streaming de tarefas ===")
// 	stream, err := client.StreamTasks(ctx, &pb.Empty{})
// 	if err != nil {
// 		log.Fatalf("Erro no stream: %v", err)
// 	}
// 	for {
// 		task, err := stream.Recv()
// 		if err == io.EOF {
// 			break  // Stream terminou
// 		}
// 		if err != nil {
// 			log.Fatalf("Erro recebendo: %v", err)
// 		}
// 		fmt.Printf("  Stream recebido: [%d] %s\n", task.Id, task.Titulo)
// 	}
//
// 	// --- Deletar ---
// 	fmt.Println("\n=== Deletando tarefa 2 ===")
// 	delResp, err := client.DeleteTask(ctx, &pb.DeleteTaskRequest{Id: 2})
// 	if err != nil {
// 		log.Printf("Erro: %v", err)
// 	} else {
// 		fmt.Printf("  Deletada: %t — %s\n", delResp.Success, delResp.Message)
// 	}
// }

func main() {
	fmt.Println("============================================================")
	fmt.Println("  CLIENTE gRPC — Código de referência")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Para rodar:")
	fmt.Println("  1. Inicie o servidor: go run ../servidor/main.go")
	fmt.Println("  2. Descomente o código neste arquivo")
	fmt.Println("  3. go run main.go")
	fmt.Println()
	fmt.Println("O client gRPC é tipado e gerado pelo protoc.")
	fmt.Println("Você chama métodos como funções Go normais!")
	fmt.Println()
	fmt.Println("  client.CreateTask(ctx, &pb.CreateTaskRequest{...})")
	fmt.Println("  client.ListTasks(ctx, &pb.Empty{})")
	fmt.Println("  stream, _ := client.StreamTasks(ctx, &pb.Empty{})")
}
