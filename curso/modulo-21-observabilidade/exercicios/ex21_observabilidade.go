// Módulo 21 — Exercícios: Observabilidade & SRE
//
// 🧠 Antes de começar:
//   1. Quais são os 3 pilares da observabilidade?
//   2. O que é um SLO e como ele se diferencia de um SLA?
//   3. Por que latência p99 é mais importante que média?
//
// NOTA: Estes exercícios funcionam sem dependências externas.
//       Para Prometheus real: go get github.com/prometheus/client_golang
//       Para OpenTelemetry: go get go.opentelemetry.io/otel

package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sort"
	"sync"
	"time"
)

// =============================================================================
// EXERCÍCIO 1 — Metrics Collector (Prometheus-like) (🟢 Fácil)
// =============================================================================
// Implemente um coletor de métricas simples que suporta Counter, Gauge e Histogram.
// Não use a biblioteca Prometheus — implemente do zero para entender os tipos.

type MetricsRegistry struct {
	mu         sync.RWMutex
	counters   map[string]*Counter
	gauges     map[string]*Gauge
	histograms map[string]*Histogram
}

func NewMetricsRegistry() *MetricsRegistry {
	// TODO: inicialize os maps
	panic("implemente NewMetricsRegistry")
}

type Counter struct {
	mu    sync.Mutex
	name  string
	value float64
}

func (c *Counter) Inc() {
	// TODO: incrementa por 1
	panic("implemente Counter.Inc")
}

func (c *Counter) Add(v float64) {
	// TODO: incrementa por v (v deve ser >= 0)
	panic("implemente Counter.Add")
}

func (c *Counter) Value() float64 {
	// TODO: retorne o valor atual thread-safe
	panic("implemente Counter.Value")
}

type Gauge struct {
	mu    sync.Mutex
	name  string
	value float64
}

func (g *Gauge) Set(v float64) { panic("implemente Gauge.Set") }
func (g *Gauge) Inc()          { panic("implemente Gauge.Inc") }
func (g *Gauge) Dec()          { panic("implemente Gauge.Dec") }
func (g *Gauge) Value() float64 { panic("implemente Gauge.Value") }

type Histogram struct {
	mu      sync.Mutex
	name    string
	buckets []float64 // limites: ex [0.01, 0.05, 0.1, 0.5, 1.0]
	counts  []int     // count[i] = número de obs <= buckets[i]
	sum     float64
	count   int
}

func NewHistogram(name string, buckets []float64) *Histogram {
	// TODO: inicialize com os buckets dados (ordene-os)
	panic("implemente NewHistogram")
}

func (h *Histogram) Observe(v float64) {
	// TODO: incremente count e sum, e o bucket correspondente
	// Cada bucket conta observações <= bucket[i]
	panic("implemente Histogram.Observe")
}

// Percentile calcula o percentil p (0-100) das observações
func (h *Histogram) Percentile(p float64) float64 {
	// TODO: use interpolação linear nos buckets para estimar o percentil
	// Dica: encontre os dois buckets que contêm p% das observações
	panic("implemente Histogram.Percentile")
}

// NewCounter cria ou retorna um counter existente pelo nome
func (r *MetricsRegistry) NewCounter(name string) *Counter {
	// TODO: crie se não existe, retorne o existente se sim
	panic("implemente NewCounter")
}

func (r *MetricsRegistry) NewGauge(name string) *Gauge {
	panic("implemente NewGauge")
}

func (r *MetricsRegistry) NewHistogram(name string, buckets []float64) *Histogram {
	panic("implemente NewHistogram na registry")
}

// Expose gera saída no formato Prometheus text (simplificado)
func (r *MetricsRegistry) Expose() string {
	// TODO: para cada métrica, gere linhas no formato:
	// # TYPE nome_counter counter
	// nome_counter 42
	// # TYPE nome_gauge gauge
	// nome_gauge 7
	// # TYPE nome_histogram histogram
	// nome_histogram_bucket{le="0.1"} 5
	// nome_histogram_sum 2.5
	// nome_histogram_count 10
	panic("implemente Expose")
}

// =============================================================================
// EXERCÍCIO 2 — HTTP Middleware de Instrumentação (🟡 Médio)
// =============================================================================
// Implemente um middleware que instrumenta automaticamente handlers HTTP:
//   - requests_total{method, path, status} — Counter
//   - request_duration_seconds{method, path} — Histogram
//   - active_requests — Gauge

type InstrumentedMux struct {
	mux      *http.ServeMux
	registry *MetricsRegistry
	// TODO: adicione os campos de métricas necessários (counter, histogram, gauge)
}

func NewInstrumentedMux(registry *MetricsRegistry) *InstrumentedMux {
	// TODO: crie o mux e as 3 métricas
	panic("implemente NewInstrumentedMux")
}

func (im *InstrumentedMux) Handle(pattern string, handler http.Handler) {
	// TODO: envolva o handler com instrumentação e registre no mux
	panic("implemente Handle")
}

func (im *InstrumentedMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	im.mux.ServeHTTP(w, r)
}

// responseRecorder captura o status code do response
type responseRecorder struct {
	http.ResponseWriter
	statusCode int
}

func (rr *responseRecorder) WriteHeader(code int) {
	rr.statusCode = code
	rr.ResponseWriter.WriteHeader(code)
}

// =============================================================================
// EXERCÍCIO 3 — SLO Error Budget Calculator (🔴 Difícil)
// =============================================================================
// Implemente um calculador de Error Budget baseado em SLO.
//
// Dado um SLO de 99.9% de uptime no mês:
//   - Total de requests no período: N
//   - Requests com erro: E
//   - Error rate: E/N
//   - Availability: 1 - (E/N)
//   - Error budget total: N * (1 - SLO) = N * 0.001
//   - Error budget consumido: E
//   - Error budget restante: budget_total - E
//   - Burn rate: (E/N) / (1 - SLO)
//     Burn rate 1.0 = consumindo o budget exatamente no ritmo esperado
//     Burn rate > 1.0 = vai extrapolar o budget antes do fim do período

type SLOCalculator struct {
	targetSLO float64 // ex: 0.999 para 99.9%
	windowDays int    // ex: 30 dias
	requests   []RequestRecord
	mu         sync.RWMutex
}

type RequestRecord struct {
	Timestamp time.Time
	Success   bool
	LatencyMs float64
}

func NewSLOCalculator(targetSLO float64, windowDays int) *SLOCalculator {
	// TODO: inicialize
	panic("implemente NewSLOCalculator")
}

func (s *SLOCalculator) Record(success bool, latencyMs float64) {
	// TODO: adicione um RequestRecord com o timestamp atual
	panic("implemente Record")
}

type SLOReport struct {
	TotalRequests    int
	SuccessRequests  int
	ErrorRequests    int
	Availability     float64 // 0-1
	TargetSLO        float64
	ErrorBudgetTotal int     // em número de requests que podem falhar
	ErrorBudgetUsed  int
	ErrorBudgetLeft  int
	BurnRate         float64 // 1.0 = ritmo normal
	LatencyP50Ms     float64
	LatencyP99Ms     float64
}

func (s *SLOCalculator) Report() SLOReport {
	// TODO: calcule todas as métricas do SLOReport
	// Para latência, use apenas as requests do windowDays mais recentes
	panic("implemente Report")
}

func (s *SLOCalculator) Alert() string {
	// TODO: retorne alertas baseados no burn rate:
	// burnRate > 14.4: "CRÍTICO: error budget vai acabar em 1 hora"
	// burnRate > 6.0:  "ALTO: error budget vai acabar em 6 horas"
	// burnRate > 3.0:  "MÉDIO: error budget consumição acima do normal"
	// burnRate > 1.0:  "AVISO: burn rate acima de 1.0"
	// default:         "OK"
	panic("implemente Alert")
}

// =============================================================================
// MAIN
// =============================================================================

func testMetrics() {
	fmt.Println("\n=== METRICS COLLECTOR ===")
	registry := NewMetricsRegistry()

	requests := registry.NewCounter("http_requests_total")
	errors := registry.NewCounter("http_errors_total")
	activeConns := registry.NewGauge("http_active_connections")
	latency := registry.NewHistogram("http_duration_seconds",
		[]float64{0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1.0})

	// Simula tráfego
	for i := 0; i < 100; i++ {
		duration := rand.Float64() * 0.5
		requests.Inc()
		activeConns.Inc()
		latency.Observe(duration)
		if rand.Float32() < 0.05 { // 5% de erro
			errors.Inc()
		}
		activeConns.Dec()
	}

	fmt.Printf("Requests: %.0f | Errors: %.0f | Active: %.0f\n",
		requests.Value(), errors.Value(), activeConns.Value())
	fmt.Printf("Latência p50: %.3fs | p99: %.3fs\n",
		latency.Percentile(50), latency.Percentile(99))
	fmt.Println("\nPrometheus format:")
	fmt.Println(registry.Expose())
}

func testSLO() {
	fmt.Println("\n=== SLO ERROR BUDGET ===")
	calc := NewSLOCalculator(0.999, 30) // 99.9% SLO, janela de 30 dias

	// Simula 10.000 requests com 5% de erro (burn rate ~50x!)
	for i := 0; i < 10000; i++ {
		success := rand.Float32() > 0.05
		latency := rand.Float64() * 200
		calc.Record(success, latency)
	}

	report := calc.Report()
	fmt.Printf("Total: %d | Erros: %d | Disponibilidade: %.4f%%\n",
		report.TotalRequests, report.ErrorRequests, report.Availability*100)
	fmt.Printf("Error Budget: total=%d, usado=%d, restante=%d\n",
		report.ErrorBudgetTotal, report.ErrorBudgetUsed, report.ErrorBudgetLeft)
	fmt.Printf("Burn Rate: %.1fx\n", report.BurnRate)
	fmt.Printf("Latência: p50=%.1fms | p99=%.1fms\n", report.LatencyP50Ms, report.LatencyP99Ms)
	fmt.Printf("Alerta: %s\n", calc.Alert())
}

func main() {
	rand.Seed(time.Now().UnixNano())
	_ = math.Inf(1) // evita import não-usado
	_ = sort.Ints   // evita import não-usado

	testMetrics()
	testSLO()
	fmt.Println("\n✅ Exercícios do Módulo 21 concluídos!")
}
