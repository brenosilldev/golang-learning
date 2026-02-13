package main

import (
	"fmt"
)

// ============================================================================
// MÓDULO 15 — Servidor gRPC
// ============================================================================
//
// Para rodar este servidor:
//   1. Instale as dependências:
//      go get google.golang.org/grpc
//      go get google.golang.org/protobuf
//
//   2. Gere o código do .proto:
//      protoc --go_out=. --go-grpc_out=. proto/task.proto
//
//   3. Descomente o código abaixo e rode:
//      go run main.go
//
// O servidor roda na porta 50051 (padrão gRPC)
// ============================================================================

// import (
// 	"context"
// 	"log"
// 	"net"
// 	"sync"
// 	"time"
//
// 	"google.golang.org/grpc"
// 	pb "seu-modulo/curso/modulo-15-grpc/exemplos/pb"
// )

// --- Implementação do servidor ---
//
// type taskServer struct {
// 	pb.UnimplementedTaskServiceServer  // Embedding obrigatório
// 	mu        sync.RWMutex
// 	tasks     map[int32]*pb.Task
// 	nextID    int32
// }
//
// func newServer() *taskServer {
// 	return &taskServer{
// 		tasks:  make(map[int32]*pb.Task),
// 		nextID: 1,
// 	}
// }
//
// // Unary RPC — Criar tarefa
// func (s *taskServer) CreateTask(ctx context.Context, req *pb.CreateTaskRequest) (*pb.TaskResponse, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
//
// 	task := &pb.Task{
// 		Id:        s.nextID,
// 		Titulo:    req.Titulo,
// 		Descricao: req.Descricao,
// 		Concluida: false,
// 		CriadoEm:  time.Now().Unix(),
// 	}
// 	s.tasks[s.nextID] = task
// 	s.nextID++
//
// 	log.Printf("Tarefa criada: %s (ID: %d)", task.Titulo, task.Id)
// 	return &pb.TaskResponse{Task: task}, nil
// }
//
// // Unary RPC — Buscar tarefa
// func (s *taskServer) GetTask(ctx context.Context, req *pb.GetTaskRequest) (*pb.TaskResponse, error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
//
// 	task, ok := s.tasks[req.Id]
// 	if !ok {
// 		return nil, fmt.Errorf("tarefa %d não encontrada", req.Id)
// 	}
// 	return &pb.TaskResponse{Task: task}, nil
// }
//
// // Unary RPC — Listar todas
// func (s *taskServer) ListTasks(ctx context.Context, _ *pb.Empty) (*pb.TaskList, error) {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
//
// 	list := make([]*pb.Task, 0, len(s.tasks))
// 	for _, t := range s.tasks {
// 		list = append(list, t)
// 	}
// 	return &pb.TaskList{Tasks: list}, nil
// }
//
// // Server Streaming RPC — Stream de tarefas (uma por uma)
// func (s *taskServer) StreamTasks(_ *pb.Empty, stream pb.TaskService_StreamTasksServer) error {
// 	s.mu.RLock()
// 	defer s.mu.RUnlock()
//
// 	for _, task := range s.tasks {
// 		if err := stream.Send(task); err != nil {
// 			return err
// 		}
// 		time.Sleep(500 * time.Millisecond)  // Simular delay
// 	}
// 	return nil
// }
//
// // Delete RPC
// func (s *taskServer) DeleteTask(ctx context.Context, req *pb.DeleteTaskRequest) (*pb.DeleteResponse, error) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()
//
// 	if _, ok := s.tasks[req.Id]; !ok {
// 		return &pb.DeleteResponse{Success: false, Message: "não encontrada"}, nil
// 	}
// 	delete(s.tasks, req.Id)
// 	return &pb.DeleteResponse{Success: true, Message: "deletada"}, nil
// }
//
// func main() {
// 	// Criar listener TCP
// 	lis, err := net.Listen("tcp", ":50051")
// 	if err != nil {
// 		log.Fatalf("Falha ao ouvir: %v", err)
// 	}
//
// 	// Criar servidor gRPC
// 	grpcServer := grpc.NewServer()
//
// 	// Registrar nosso serviço
// 	pb.RegisterTaskServiceServer(grpcServer, newServer())
//
// 	log.Println("🚀 Servidor gRPC rodando em :50051")
// 	if err := grpcServer.Serve(lis); err != nil {
// 		log.Fatalf("Falha ao servir: %v", err)
// 	}
// }

func main() {
	fmt.Println("============================================================")
	fmt.Println("  SERVIDOR gRPC — Código de referência")
	fmt.Println("============================================================")
	fmt.Println()
	fmt.Println("Para rodar:")
	fmt.Println("  1. go get google.golang.org/grpc")
	fmt.Println("  2. go get google.golang.org/protobuf")
	fmt.Println("  3. Instale protoc: sudo apt install protobuf-compiler")
	fmt.Println("  4. go install google.golang.org/protobuf/cmd/protoc-gen-go@latest")
	fmt.Println("  5. go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest")
	fmt.Println("  6. protoc --go_out=. --go-grpc_out=. proto/task.proto")
	fmt.Println("  7. Descomente o código neste arquivo")
	fmt.Println("  8. go run main.go")
}
