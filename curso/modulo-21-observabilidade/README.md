# 🔭 Módulo 21 — Observabilidade & SRE

> **Nível**: Avançado | **Pré-requisito**: Módulos 14 (APIs), 17 (Docker), 18 (Fundamentos Distribuídos)

---

## 🤔 Por que Observabilidade?

> "Você não pode gerenciar o que não consegue medir."

Num sistema distribuído com 50 microserviços, quando algo falha você precisa responder:
- **Onde exatamente falhou?** (rastreamento)
- **O que estava acontecendo antes de falhar?** (logs)
- **Qual era o estado do sistema?** (métricas)

### Os 3 Pilares da Observabilidade

```
📊 Métricas        📝 Logs            🔍 Traces
─────────────      ─────────────      ─────────────
"O quê?"           "Por quê?"         "Onde?"
Números no tempo   Eventos com        Caminho de uma
(latência, erros,  contexto e         requisição entre
throughput, uso    detalhes sobre     múltiplos serviços
de CPU)            o que aconteceu
```

---

## 📊 Métricas com Prometheus

Prometheus usa um modelo **pull** — ele periodicamente raspa (`scrape`) um endpoint `/metrics` da sua aplicação.

### Tipos de métricas:

| Tipo | Uso | Exemplo |
|------|-----|---------|
| **Counter** | Conta eventos (só sobe) | Total de requests, erros |
| **Gauge** | Valor atual no tempo | Goroutines ativas, memória |
| **Histogram** | Distribuição de valores | Latência de requests |
| **Summary** | Percentis calculados no cliente | p99 de latência |

```go
package main

import (
    "fmt"
    "math/rand"
    "net/http"
    "time"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promauto"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

// Declara métricas globais
var (
    requestsTotal = promauto.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total de requisições HTTP",
        },
        []string{"method", "path", "status"}, // labels
    )

    requestDuration = promauto.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "http_request_duration_seconds",
            Help:    "Duração das requisições HTTP",
            Buckets: prometheus.DefBuckets, // .005, .01, .025, .05, .1, .25, .5, 1, 2.5, 5, 10
        },
        []string{"method", "path"},
    )

    activeConnections = promauto.NewGauge(prometheus.GaugeOpts{
        Name: "http_active_connections",
        Help: "Conexões HTTP ativas no momento",
    })

    dbQueryDuration = promauto.NewHistogram(prometheus.HistogramOpts{
        Name:    "db_query_duration_seconds",
        Help:    "Duração das queries ao banco",
        Buckets: []float64{.001, .005, .01, .025, .05, .1, .25, .5, 1},
    })
)

// Middleware que instrumenta automaticamente todos os handlers
func metricsMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        start := time.Now()
        activeConnections.Inc()
        defer activeConnections.Dec()

        // Wrapper para capturar o status code
        rw := &responseWriter{ResponseWriter: w, statusCode: 200}
        next.ServeHTTP(rw, r)

        duration := time.Since(start).Seconds()
        status := fmt.Sprintf("%d", rw.statusCode)

        requestsTotal.WithLabelValues(r.Method, r.URL.Path, status).Inc()
        requestDuration.WithLabelValues(r.Method, r.URL.Path).Observe(duration)
    })
}

type responseWriter struct {
    http.ResponseWriter
    statusCode int
}

func (rw *responseWriter) WriteHeader(code int) {
    rw.statusCode = code
    rw.ResponseWriter.WriteHeader(code)
}

func main() {
    mux := http.NewServeMux()

    // Handler de negócio
    mux.HandleFunc("/api/pedidos", func(w http.ResponseWriter, r *http.Request) {
        // Simula latência variável
        time.Sleep(time.Duration(rand.Intn(200)) * time.Millisecond)

        // Instrumenta query no banco
        dbStart := time.Now()
        time.Sleep(time.Duration(rand.Intn(50)) * time.Millisecond) // simula query
        dbQueryDuration.Observe(time.Since(dbStart).Seconds())

        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"pedidos":[]}`)
    })

    // Endpoint de métricas para o Prometheus raspar
    mux.Handle("/metrics", promhttp.Handler())

    // Health checks
    mux.HandleFunc("/health/live", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
    })
    mux.HandleFunc("/health/ready", func(w http.ResponseWriter, r *http.Request) {
        // Verifica se banco e dependências estão ok
        w.WriteHeader(http.StatusOK)
    })

    server := &http.Server{
        Addr:    ":8080",
        Handler: metricsMiddleware(mux),
    }
    fmt.Println("Servidor em :8080 | Métricas em :8080/metrics")
    server.ListenAndServe()
}
```

---

## 🔍 Distributed Tracing com OpenTelemetry

OpenTelemetry é o padrão aberto para traces, métricas e logs. Funciona com Jaeger, Zipkin, Datadog, etc.

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"

    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/attribute"
    "go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
    "go.opentelemetry.io/otel/propagation"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
    "go.opentelemetry.io/otel/trace"
)

var tracer trace.Tracer

func initTracer(ctx context.Context) func() {
    exporter, err := otlptracehttp.New(ctx,
        otlptracehttp.WithEndpoint("localhost:4318"), // Jaeger OTLP endpoint
        otlptracehttp.WithInsecure(),
    )
    if err != nil {
        log.Fatal(err)
    }

    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(resource.NewWithAttributes(
            semconv.SchemaURL,
            semconv.ServiceName("api-pedidos"),
            semconv.ServiceVersion("1.0.0"),
            attribute.String("environment", "production"),
        )),
        sdktrace.WithSampler(sdktrace.AlwaysSample()), // em produção: TraceIDRatioBased(0.1)
    )

    otel.SetTracerProvider(tp)
    otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
        propagation.TraceContext{}, // W3C Trace Context
        propagation.Baggage{},
    ))

    tracer = otel.Tracer("api-pedidos")
    return func() { tp.Shutdown(ctx) }
}

// Handler com tracing completo
func handleCriarPedido(w http.ResponseWriter, r *http.Request) {
    // Extrai contexto de trace do header (se vier de outro serviço)
    ctx := otel.GetTextMapPropagator().Extract(r.Context(), propagation.HeaderCarrier(r.Header))

    // Cria span para essa operação
    ctx, span := tracer.Start(ctx, "criar-pedido",
        trace.WithAttributes(
            attribute.String("user.id", r.Header.Get("X-User-ID")),
            attribute.String("http.method", r.Method),
        ),
    )
    defer span.End()

    // Sub-operação: validação
    if err := validarPedido(ctx); err != nil {
        span.RecordError(err)
        span.SetStatus(2 /* Error */, err.Error())
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    // Sub-operação: banco
    pedidoID, err := salvarNoBanco(ctx)
    if err != nil {
        span.RecordError(err)
        http.Error(w, "erro interno", http.StatusInternalServerError)
        return
    }

    span.SetAttributes(attribute.String("pedido.id", pedidoID))
    fmt.Fprintf(w, `{"id":"%s"}`, pedidoID)
}

func validarPedido(ctx context.Context) error {
    _, span := tracer.Start(ctx, "validar-pedido")
    defer span.End()
    // lógica de validação...
    return nil
}

func salvarNoBanco(ctx context.Context) (string, error) {
    _, span := tracer.Start(ctx, "db.insert",
        trace.WithAttributes(
            attribute.String("db.system", "postgresql"),
            attribute.String("db.statement", "INSERT INTO pedidos ..."),
        ),
    )
    defer span.End()
    return "pedido-123", nil
}

func main() {
    ctx := context.Background()
    cleanup := initTracer(ctx)
    defer cleanup()

    http.HandleFunc("/pedidos", handleCriarPedido)
    log.Println("Servidor em :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))
}
```

---

## 📝 Structured Logging com Zerolog

Logs em JSON são processáveis por Elasticsearch, Loki, CloudWatch, etc.

```go
package main

import (
    "net/http"
    "os"
    "time"

    "github.com/rs/zerolog"
    "github.com/rs/zerolog/log"
)

func init() {
    // Produção: JSON puro
    zerolog.TimeFieldFormat = time.RFC3339Nano
    log.Logger = zerolog.New(os.Stdout).With().
        Timestamp().
        Str("service", "api-pedidos").
        Str("version", "1.0.0").
        Logger()

    // Desenvolvimento: output legível
    if os.Getenv("ENV") != "production" {
        log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
    }
}

func handler(w http.ResponseWriter, r *http.Request) {
    start := time.Now()
    requestID := r.Header.Get("X-Request-ID")

    // Logger com contexto da requisição
    logger := log.With().
        Str("request_id", requestID).
        Str("method", r.Method).
        Str("path", r.URL.Path).
        Str("remote_addr", r.RemoteAddr).
        Logger()

    logger.Info().Msg("requisição recebida")

    // Em produção, passe o logger no contexto
    ctx := logger.WithContext(r.Context())
    _ = ctx

    // Log de erro com stack-like context
    if r.URL.Path == "/erro" {
        logger.Error().
            Err(http.ErrHandlerTimeout).
            Str("user_id", "123").
            Int("tentativa", 3).
            Msg("falha ao processar pedido")
        http.Error(w, "erro", http.StatusInternalServerError)
        return
    }

    duration := time.Since(start)
    logger.Info().
        Dur("duration_ms", duration).
        Int("status", 200).
        Msg("requisição concluída")
}
```

---

## 📏 SLI / SLO / SLA

| Sigla | Nome | Exemplo |
|-------|------|---------|
| **SLI** | Service Level Indicator | `latência p99 = 120ms` |
| **SLO** | Service Level Objective | `p99 < 200ms em 99.9% do tempo` |
| **SLA** | Service Level Agreement | `Downtime máximo 4h/mês (99.95% uptime)` |

**Error Budget**: se o SLO é 99.9% uptime → posso ter 43min/mês de downtime. Se o budget acabou, parar de fazer deploys.

---

## 📋 Exercícios

### 🟢 1. Dashboard Prometheus + Grafana
Com Docker Compose, suba Prometheus + Grafana + sua API e crie um dashboard com:
- Requests por segundo, latência p50/p95/p99, taxa de erros

### 🟡 2. Trace Propagation
Implemente dois serviços (A → B via HTTP) onde o trace ID é propagado automaticamente via headers W3C. Visualize no Jaeger que a chamada ao serviço B aparece como span filho do A.

### 🔴 3. SLO Alert
Configure um alerta no Prometheus que dispara quando o error budget estiver sendo consumido mais rápido do que o normal (multi-burn-rate alerting).

---

> **Próximo**: [Módulo 22 — Cloud Native & Kubernetes](../modulo-22-cloud-native-k8s/README.md)
