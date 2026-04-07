# Módulo 15 — gRPC e Protocol Buffers

[← APIs](../modulo-14-apis/README.md) | [Próximo: Banco de Dados →](../modulo-16-database/README.md)

---

> **Antes de ler — tente responder:**
> 1. Por que gRPC é mais eficiente que REST/JSON para comunicação entre serviços?
> 2. O que é um `.proto` e qual sua vantagem sobre documentação informal de API?
> 3. Quando você escolheria gRPC sobre REST?

---

## 1. Por que gRPC?

| | HTTP/JSON (REST) | gRPC/Protobuf |
|---|---|---|
| **Serialização** | Texto (JSON) — lento para parsear | Binário — 5-10x menor e mais rápido |
| **Contrato** | Informal (OpenAPI/docs) | `.proto` — tipado, versionado, validado |
| **Geração de código** | Manual ou Swagger | Automático para 10+ linguagens |
| **Streaming** | Sem suporte nativo | 4 tipos (unary, server, client, bidi) |
| **Casos de uso** | APIs públicas, browsers | Comunicação interna entre microserviços |

**gRPC é o padrão para comunicação entre microserviços.** Kubernetes, etcd, CockroachDB, Envoy — todos usam gRPC internamente.

---

## 2. Protocol Buffers — Definindo o Contrato

O `.proto` é o coração do gRPC. Define mensagens e serviços de forma independente de linguagem:

```protobuf
syntax = "proto3";
package tarefas;
option go_package = "./pb;pb";  // caminho de geração

// Mensagem = struct equivalente
message Tarefa {
  int32  id        = 1;  // field number — NÃO mude em produção (quebra compatibilidade)
  string titulo    = 2;
  string descricao = 3;
  bool   concluida = 4;
}

// Request/Response separados — boa prática (permite evoluir sem breaking changes)
message CriarTarefaRequest {
  string titulo    = 1;
  string descricao = 2;
}

message CriarTarefaResponse {
  Tarefa tarefa = 1;
}

message BuscarTarefaRequest {
  int32 id = 1;
}

message ListarTarefasRequest {
  int32  pagina   = 1;
  int32  por_page = 2;
}

message ListarTarefasResponse {
  repeated Tarefa tarefas = 1;  // repeated = slice em Go
  int32           total   = 2;
}

// Serviço = interface com métodos RPC
service TarefaService {
  rpc CriarTarefa(CriarTarefaRequest) returns (CriarTarefaResponse);
  rpc BuscarTarefa(BuscarTarefaRequest) returns (Tarefa);
  rpc ListarTarefas(ListarTarefasRequest) returns (ListarTarefasResponse);
  rpc StreamTarefas(ListarTarefasRequest) returns (stream Tarefa); // server streaming
}
```

### Gerando código Go

```bash
# Instalar ferramentas (uma vez)
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
sudo apt install -y protobuf-compiler  # protoc

# Gerar
protoc \
  --go_out=. \
  --go_opt=paths=source_relative \
  --go-grpc_out=. \
  --go-grpc_opt=paths=source_relative \
  proto/tarefa.proto

# Resultado:
# proto/tarefa.pb.go       ← structs das mensagens
# proto/tarefa_grpc.pb.go  ← interface do servidor + client stub
```

---

## 3. Implementando o Servidor

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net"
    "sync"

    "google.golang.org/grpc"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pb "github.com/alice/myapp/pb"
)

// Implementa a interface gerada pelo protoc
type server struct {
    pb.UnimplementedTarefaServiceServer  // embed para forward compatibility
    mu      sync.RWMutex
    tarefas map[int32]*pb.Tarefa
    nextID  int32
}

func (s *server) CriarTarefa(ctx context.Context, req *pb.CriarTarefaRequest) (*pb.CriarTarefaResponse, error) {
    // Validação de input
    if req.Titulo == "" {
        return nil, status.Error(codes.InvalidArgument, "título é obrigatório")
    }

    s.mu.Lock()
    defer s.mu.Unlock()

    s.nextID++
    tarefa := &pb.Tarefa{
        Id:       s.nextID,
        Titulo:   req.Titulo,
        Descricao: req.Descricao,
    }
    s.tarefas[s.nextID] = tarefa

    return &pb.CriarTarefaResponse{Tarefa: tarefa}, nil
}

func (s *server) BuscarTarefa(ctx context.Context, req *pb.BuscarTarefaRequest) (*pb.Tarefa, error) {
    s.mu.RLock()
    defer s.mu.RUnlock()

    tarefa, ok := s.tarefas[req.Id]
    if !ok {
        // Use status.Error para erros gRPC — o cliente recebe o código correto
        return nil, status.Errorf(codes.NotFound, "tarefa %d não encontrada", req.Id)
    }
    return tarefa, nil
}

// Server streaming — envia múltiplas respostas para um request
func (s *server) StreamTarefas(req *pb.ListarTarefasRequest, stream pb.TarefaService_StreamTarefasServer) error {
    s.mu.RLock()
    defer s.mu.RUnlock()

    for _, tarefa := range s.tarefas {
        // Verificar contexto em cada iteração — respeitar cancelamento do cliente
        if err := stream.Context().Err(); err != nil {
            return status.Error(codes.Canceled, "stream cancelado pelo cliente")
        }

        if err := stream.Send(tarefa); err != nil {
            return fmt.Errorf("erro ao enviar tarefa: %w", err)
        }
    }
    return nil
}

func main() {
    lis, err := net.Listen("tcp", ":50051")
    if err != nil {
        log.Fatalf("falhou ao escutar: %v", err)
    }

    // Interceptors = middleware do gRPC
    s := grpc.NewServer(
        grpc.UnaryInterceptor(loggingInterceptor),
    )

    pb.RegisterTarefaServiceServer(s, &server{
        tarefas: make(map[int32]*pb.Tarefa),
    })

    log.Println("servidor gRPC em :50051")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("falhou ao servir: %v", err)
    }
}

// Interceptor = middleware para RPCs unários
func loggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
    log.Printf("RPC: %s", info.FullMethod)
    resp, err := handler(ctx, req)
    if err != nil {
        log.Printf("RPC %s erro: %v", info.FullMethod, err)
    }
    return resp, err
}
```

---

## 4. Implementando o Cliente

```go
package main

import (
    "context"
    "io"
    "log"
    "time"

    "google.golang.org/grpc"
    "google.golang.org/grpc/credentials/insecure"
    "google.golang.org/grpc/codes"
    "google.golang.org/grpc/status"

    pb "github.com/alice/myapp/pb"
)

func main() {
    // Conectar (insecure para desenvolvimento — use TLS em produção)
    conn, err := grpc.Dial("localhost:50051",
        grpc.WithTransportCredentials(insecure.NewCredentials()),
        grpc.WithBlock(),  // bloqueia até conectar
    )
    if err != nil {
        log.Fatalf("não conectou: %v", err)
    }
    defer conn.Close()

    client := pb.NewTarefaServiceClient(conn)
    ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
    defer cancel()

    // Criar tarefa
    resp, err := client.CriarTarefa(ctx, &pb.CriarTarefaRequest{
        Titulo:   "Aprender gRPC",
        Descricao: "Implementar servidor e cliente",
    })
    if err != nil {
        // Extrair código de status gRPC do erro
        if st, ok := status.FromError(err); ok {
            log.Printf("código: %v, mensagem: %s", st.Code(), st.Message())
        }
        log.Fatalf("erro ao criar: %v", err)
    }
    log.Printf("criado: ID=%d, Título=%s", resp.Tarefa.Id, resp.Tarefa.Titulo)

    // Buscar tarefa
    tarefa, err := client.BuscarTarefa(ctx, &pb.BuscarTarefaRequest{Id: resp.Tarefa.Id})
    if err != nil {
        if status.Code(err) == codes.NotFound {
            log.Println("tarefa não encontrada")
        }
        log.Fatalf("erro ao buscar: %v", err)
    }
    log.Printf("encontrado: %s", tarefa.Titulo)

    // Consumir server stream
    stream, err := client.StreamTarefas(ctx, &pb.ListarTarefasRequest{})
    if err != nil {
        log.Fatalf("erro ao abrir stream: %v", err)
    }
    for {
        tarefa, err := stream.Recv()
        if err == io.EOF {
            break  // stream encerrado pelo servidor
        }
        if err != nil {
            log.Fatalf("erro no stream: %v", err)
        }
        log.Printf("stream recebeu: %s", tarefa.Titulo)
    }
}
```

---

## 5. Códigos de Status gRPC — Equivalentes HTTP

| Código gRPC | HTTP | Quando usar |
|------------|------|-------------|
| `codes.OK` | 200 | Sucesso |
| `codes.InvalidArgument` | 400 | Input inválido |
| `codes.NotFound` | 404 | Recurso não existe |
| `codes.AlreadyExists` | 409 | Conflito |
| `codes.PermissionDenied` | 403 | Sem autorização |
| `codes.Unauthenticated` | 401 | Não autenticado |
| `codes.ResourceExhausted` | 429 | Rate limit |
| `codes.Internal` | 500 | Erro interno |
| `codes.Unavailable` | 503 | Serviço fora |
| `codes.DeadlineExceeded` | 504 | Timeout |

```go
// No servidor: retorne status.Error
return nil, status.Errorf(codes.NotFound, "usuário %d não encontrado", id)

// No cliente: verifique o código
if status.Code(err) == codes.NotFound {
    // trate especificamente
}
```

---

## ✅ Checklist de gRPC para Produção

- [ ] `.proto` com **field numbers estáveis** (nunca mude ou reutilize field numbers)
- [ ] Servidor embute `UnimplementedXxxServer` para forward compatibility
- [ ] Todos os RPCs verificam `ctx.Err()` antes de operações longas
- [ ] Erros retornam `status.Error` com código correto (não `fmt.Errorf`)
- [ ] Interceptors configurados para logging e recovery de panics
- [ ] TLS configurado em produção (não `insecure`)
- [ ] Reflection habilitada em desenvolvimento para ferramentas como `grpcurl`

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/proto/task.proto` | Definição do serviço |
| `exemplos/servidor/main.go` | Servidor gRPC completo |
| `exemplos/cliente/main.go` | Cliente gRPC com streaming |
| `exercicios/ex15_grpc.go` | 🏋️ Exercícios |

---

## 📋 Exercícios

### 🟢 1. Serviço de Calculadora
Crie um serviço gRPC `Calculator` com métodos `Add`, `Subtract`, `Multiply`, `Divide`. O `.proto` deve ter request/response separados. O servidor deve retornar `codes.InvalidArgument` para divisão por zero.

### 🟡 2. Server Streaming — Log de Eventos
Implemente um método que faz server streaming: `WatchEvents(filter) returns (stream Event)`. O servidor deve enviar eventos periódicos enquanto o stream está aberto e parar quando o contexto for cancelado.

### 🟡 3. Interceptor de Autenticação
Crie um `UnaryInterceptor` que verifica um token JWT no metadata da request gRPC. Retorne `codes.Unauthenticated` se o token for inválido. Demonstre como o cliente envia o token via metadata.

### 🔴 4. gRPC Gateway — REST + gRPC ao mesmo tempo
Use `grpc-gateway` para expor o mesmo serviço gRPC como API REST. Adicione as anotações no `.proto` e configure o gateway. Demonstre que um `curl` para `/v1/tarefas` e um cliente gRPC funcionam ao mesmo tempo.

---

> **Confirme seu aprendizado**: releia as 3 perguntas do início. Consegue responder agora?

[← APIs](../modulo-14-apis/README.md) | [Próximo: Banco de Dados →](../modulo-16-database/README.md)
