# Módulo 15 — gRPC e Protocol Buffers

[← APIs](../modulo-14-apis/README.md) | [Próximo: Banco de Dados →](../modulo-16-database/README.md)

---

## 📖 Por que gRPC?

| HTTP/JSON (REST) | gRPC/Protobuf |
|-------------------|---------------|
| Texto (JSON) — lento para parsear | Binário — 10x mais rápido |
| Schema informal (documentação) | Schema formal (.proto) — contrato |
| Unidirecional (request/response) | 4 tipos de stream |
| Qualquer linguagem lê | Gera código para 10+ linguagens |

**gRPC é o padrão de comunicação entre microserviços.** Kubernetes, Envoy, etcd, CockroachDB — todos usam gRPC internamente.

### Os 4 tipos de RPC

```
1. Unary         → Request → Response (como REST)
2. Server Stream → Request → Stream de Responses  
3. Client Stream → Stream de Requests → Response
4. Bidi Stream   → Stream ↔ Stream (real-time)
```

### Como funciona

```
1. Você define o contrato em um .proto
2. protoc gera código Go (structs + interfaces)
3. Você implementa o servidor
4. O client é gerado automaticamente
```

### Arquivo .proto
```protobuf
syntax = "proto3";
package taskapi;
option go_package = "./pb";

message Task {
  int32 id = 1;
  string titulo = 2;
  bool concluida = 3;
}

message CreateTaskRequest {
  string titulo = 1;
}

message TaskResponse {
  Task task = 1;
}

message ListRequest {}

message ListResponse {
  repeated Task tasks = 1;
}

service TaskService {
  rpc CreateTask(CreateTaskRequest) returns (TaskResponse);
  rpc ListTasks(ListRequest) returns (ListResponse);
  rpc StreamTasks(ListRequest) returns (stream Task);  // server stream
}
```

### Instalação
```bash
# Instalar protoc (compilador)
sudo apt install -y protobuf-compiler

# Plugins Go para protoc
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Gerar código
protoc --go_out=. --go-grpc_out=. proto/task.proto
```

---

## 📂 Arquivos

| Arquivo | Conteúdo |
|---------|----------|
| `exemplos/proto/task.proto` | Definição do serviço |
| `exemplos/servidor/main.go` | Servidor gRPC |
| `exemplos/cliente/main.go` | Cliente gRPC |
| `exercicios/ex15_grpc.go` | 🏋️ Exercícios |

---

[← APIs](../modulo-14-apis/README.md) | [Próximo: Banco de Dados →](../modulo-16-database/README.md)
