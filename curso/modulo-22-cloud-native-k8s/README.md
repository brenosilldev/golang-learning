# ☸️ Módulo 22 — Cloud Native & Kubernetes

> **Nível**: Avançado | **Pré-requisito**: Módulos 14 (APIs), 17 (Docker), 21 (Observabilidade)

---

## 🤔 Cloud Native em Go

Uma aplicação Cloud Native foi **projetada para rodar em Kubernetes** — ela é stateless, configurável via env vars/ConfigMaps, sabe reportar sua própria saúde e desligar graciosamente.

Checklist de uma Go app K8s-ready:
- ✅ **Graceful shutdown** — fecha conexões em andamento quando recebe SIGTERM
- ✅ **Health checks** — `/health/live` e `/health/ready`
- ✅ **Config via env vars** — nunca hardcode credenciais
- ✅ **Métricas Prometheus** — `/metrics`
- ✅ **Logging estruturado** — JSON para stdout
- ✅ **Resource limits** — define CPU/memória para o K8s não sofrer
- ✅ **Dockerfile multi-stage** — imagem mínima (scratch ou distroless)

---

## 🚀 Go App Pronta para Kubernetes

```go
package main

import (
    "context"
    "errors"
    "fmt"
    "log/slog"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Config struct {
    Port        string
    DatabaseURL string
    AppVersion  string
}

func configFromEnv() Config {
    port := os.Getenv("PORT")
    if port == "" {
        port = "8080"
    }
    return Config{
        Port:        port,
        DatabaseURL: os.Getenv("DATABASE_URL"),   // nunca hardcode!
        AppVersion:  os.Getenv("APP_VERSION"),
    }
}

type App struct {
    cfg    Config
    logger *slog.Logger
    db     *FakeDB // substitua por *sql.DB
    ready  bool
}

type FakeDB struct{}
func (db *FakeDB) Ping() error { return nil }

func NewApp(cfg Config) *App {
    logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
        Level: slog.LevelInfo,
    }))
    slog.SetDefault(logger)
    return &App{cfg: cfg, logger: logger, db: &FakeDB{}}
}

func (a *App) routes() http.Handler {
    mux := http.NewServeMux()

    // Liveness probe — a aplicação está viva?
    // K8s reinicia o pod se retornar erro
    mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprintf(w, `{"status":"ok","version":"%s"}`, a.cfg.AppVersion)
    })

    // Readiness probe — a aplicação está pronta para receber tráfego?
    // K8s remove o pod do load balancer se retornar erro
    mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
        if !a.ready {
            http.Error(w, `{"status":"starting"}`, http.StatusServiceUnavailable)
            return
        }
        if err := a.db.Ping(); err != nil {
            a.logger.Error("banco indisponível", "error", err)
            http.Error(w, `{"status":"db_unavailable"}`, http.StatusServiceUnavailable)
            return
        }
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"ready"}`)
    })

    // Handler de negócio
    mux.HandleFunc("/api/v1/pedidos", func(w http.ResponseWriter, r *http.Request) {
        a.logger.Info("requisição recebida",
            "method", r.Method,
            "path", r.URL.Path,
            "request_id", r.Header.Get("X-Request-ID"),
        )
        fmt.Fprint(w, `{"pedidos":[]}`)
    })

    return mux
}

func (a *App) Run() error {
    slog.Info("iniciando aplicação", "port", a.cfg.Port, "version", a.cfg.AppVersion)

    server := &http.Server{
        Addr:         ":" + a.cfg.Port,
        Handler:      a.routes(),
        ReadTimeout:  10 * time.Second,
        WriteTimeout: 30 * time.Second,
        IdleTimeout:  120 * time.Second,
    }

    // Inicia em goroutine separada
    errCh := make(chan error, 1)
    go func() {
        if err := server.ListenAndServe(); !errors.Is(err, http.ErrServerClosed) {
            errCh <- err
        }
    }()

    // Simula inicialização (conectar ao banco, etc.)
    time.Sleep(500 * time.Millisecond)
    a.ready = true
    slog.Info("aplicação pronta")

    // Aguarda sinal de parada
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

    select {
    case err := <-errCh:
        return fmt.Errorf("servidor falhou: %w", err)
    case sig := <-quit:
        slog.Info("recebeu sinal de parada", "signal", sig.String())
    }

    // Graceful shutdown — K8s envia SIGTERM e aguarda terminationGracePeriodSeconds
    slog.Info("desligando graciosamente...")
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()

    a.ready = false // para de aceitar novas requisições no readiness probe
    time.Sleep(5 * time.Second) // aguarda o K8s remover do load balancer (deregistration delay)

    if err := server.Shutdown(ctx); err != nil {
        return fmt.Errorf("shutdown forçado: %w", err)
    }

    slog.Info("desligado com sucesso")
    return nil
}

func main() {
    cfg := configFromEnv()
    app := NewApp(cfg)
    if err := app.Run(); err != nil {
        slog.Error("aplicação encerrou com erro", "error", err)
        os.Exit(1)
    }
}
```

---

## 🐳 Dockerfile Produção

```dockerfile
# ───── Estágio 1: Build ─────
FROM golang:1.22-alpine AS builder

WORKDIR /app

# Copia e baixa dependências primeiro (melhor cache de layer)
COPY go.mod go.sum ./
RUN go mod download

# Copia o código e compila
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build \
    -ldflags="-s -w -X main.version=$(git rev-parse --short HEAD)" \
    -o /app/server ./cmd/server

# ───── Estágio 2: Runtime ─────
# distroless tem apenas CA certificates e timezone data
FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /app/server /server

# Não rodar como root
USER nonroot:nonroot

EXPOSE 8080

ENTRYPOINT ["/server"]
```

---

## ☸️ Manifestos Kubernetes

### Deployment

```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-pedidos
  labels:
    app: api-pedidos
spec:
  replicas: 3
  selector:
    matchLabels:
      app: api-pedidos
  strategy:
    type: RollingUpdate
    rollingUpdate:
      maxSurge: 1        # cria 1 pod novo antes de matar o antigo
      maxUnavailable: 0  # nunca deixa pod sem réplica durante deploy
  template:
    metadata:
      labels:
        app: api-pedidos
      annotations:
        prometheus.io/scrape: "true"
        prometheus.io/path: "/metrics"
        prometheus.io/port: "8080"
    spec:
      terminationGracePeriodSeconds: 60  # tempo para graceful shutdown
      containers:
        - name: api-pedidos
          image: registry.company.com/api-pedidos:v1.2.3
          ports:
            - containerPort: 8080
          env:
            - name: PORT
              value: "8080"
            - name: DATABASE_URL
              valueFrom:
                secretKeyRef:
                  name: api-pedidos-secrets
                  key: database-url
            - name: APP_VERSION
              valueFrom:
                fieldRef:
                  fieldPath: metadata.labels['app']  # injeta versão via downward API
          
          # Limites de recursos (OBRIGATÓRIO em produção)
          resources:
            requests:
              memory: "64Mi"
              cpu: "50m"
            limits:
              memory: "256Mi"
              cpu: "500m"
          
          # Probes
          livenessProbe:
            httpGet:
              path: /health/live
              port: 8080
            initialDelaySeconds: 10
            periodSeconds: 15
            failureThreshold: 3
          
          readinessProbe:
            httpGet:
              path: /health/ready
              port: 8080
            initialDelaySeconds: 5
            periodSeconds: 10
            failureThreshold: 3
          
          # Lifecycle hook: aguarda remoção do load balancer antes do shutdown
          lifecycle:
            preStop:
              exec:
                command: ["/bin/sh", "-c", "sleep 5"]
```

### HPA — Horizontal Pod Autoscaler

```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: api-pedidos-hpa
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: api-pedidos
  minReplicas: 2
  maxReplicas: 20
  metrics:
    - type: Resource
      resource:
        name: cpu
        target:
          type: Utilization
          averageUtilization: 70  # escala quando CPU > 70%
    - type: Pods
      pods:
        metric:
          name: http_requests_per_second  # métrica customizada do Prometheus
        target:
          type: AverageValue
          averageValue: "1000"  # max 1000 req/s por pod
```

---

## 🎛️ Helm — Gerenciador de Pacotes K8s

```bash
# Criar chart
helm create api-pedidos

# Estrutura gerada:
# api-pedidos/
# ├── Chart.yaml        → metadados
# ├── values.yaml       → valores padrão
# ├── templates/
# │   ├── deployment.yaml
# │   ├── service.yaml
# │   ├── hpa.yaml
# │   └── _helpers.tpl

# Deploy
helm install api-pedidos ./api-pedidos \
  --set image.tag=v1.2.3 \
  --set replicaCount=3 \
  --namespace producao

# Upgrade com rollback automático em falha
helm upgrade api-pedidos ./api-pedidos \
  --set image.tag=v1.2.4 \
  --atomic \        # reverte automaticamente se falhar
  --timeout 5m
```

---

## 🤖 Kubernetes Operator em Go

Operators estendem o K8s com lógica de negócio customizada usando a API de controllers.

```go
// Esqueleto de um Operator com controller-runtime
package main

import (
    "context"
    "os"

    ctrl "sigs.k8s.io/controller-runtime"
    "sigs.k8s.io/controller-runtime/pkg/client"
    "sigs.k8s.io/controller-runtime/pkg/log/zap"
    appsv1 "k8s.io/api/apps/v1"
)

type PedidoReconciler struct {
    client.Client
}

// Reconcile é chamado sempre que o estado desejado ou real muda
func (r *PedidoReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    log := ctrl.LoggerFrom(ctx)

    // Busca o objeto real
    var deployment appsv1.Deployment
    if err := r.Get(ctx, req.NamespacedName, &deployment); err != nil {
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    log.Info("reconciliando", "deployment", req.NamespacedName)

    // Lógica: garantir que replicas >= 2 em produção
    desired := int32(2)
    if deployment.Namespace == "producao" && *deployment.Spec.Replicas < desired {
        deployment.Spec.Replicas = &desired
        if err := r.Update(ctx, &deployment); err != nil {
            return ctrl.Result{}, err
        }
        log.Info("ajustou réplicas para mínimo de produção")
    }

    return ctrl.Result{}, nil
}

func main() {
    ctrl.SetLogger(zap.New())

    mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{})
    if err != nil {
        os.Exit(1)
    }

    if err := (&PedidoReconciler{Client: mgr.GetClient()}).
        SetupWithManager(mgr); err != nil {
        os.Exit(1)
    }

    if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
        os.Exit(1)
    }
}
```

---

## 📋 Exercícios

### 🟢 1. App Completa K8s-Ready
Pegue uma API do módulo 14 e adicione: graceful shutdown, health checks, config via env vars, logging JSON. Faça o deployment usando o Dockerfile multi-stage.

### 🟡 2. Deploy no Kind (K8s local)
Usando `kind` (Kubernetes in Docker):
```bash
kind create cluster
kubectl apply -f deployment.yaml -f service.yaml -f hpa.yaml
kubectl port-forward svc/api-pedidos 8080:8080
```
Teste o rolling update e o graceful shutdown observando os logs.

### 🔴 3. Operator Simples
Crie um CRD `ManagedApp` e um operator que:
- Ao criar um `ManagedApp`, cria automaticamente Deployment + Service + HPA
- Ao deletar o `ManagedApp`, faz cleanup de todos os recursos filhos

---

> **Projeto Final**: [NexusMQ — Message Broker Distribuído](../projeto-final/README.md) — junte tudo que aprendeu!
