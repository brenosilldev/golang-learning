package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime"
	"time"
)

// ============================================================================
// MÓDULO 17 — App Go pronto para Docker
// ============================================================================
//
// Esta é uma API simples otimizada para rodar em container:
//   - Lê configuração de variáveis de ambiente
//   - Health check endpoint
//   - Graceful shutdown
//   - Versão injetada no build
//
// Build local:
//   go build -ldflags="-X main.version=1.0.0" -o api .
//
// Docker:
//   docker build -t minha-api .
//   docker run -p 8080:8080 -e PORT=8080 minha-api
//
// Docker Compose:
//   docker compose up --build
// ============================================================================

// Injetada no build com -ldflags
var version = "dev"

type HealthResponse struct {
	Status  string `json:"status"`
	Version string `json:"version"`
	Uptime  string `json:"uptime"`
	GoVer   string `json:"go_version"`
}

type InfoResponse struct {
	App          string `json:"app"`
	Version      string `json:"version"`
	Environment  string `json:"environment"`
	Port         string `json:"port"`
	GoVersion    string `json:"go_version"`
	NumCPU       int    `json:"num_cpu"`
	NumGoroutine int    `json:"num_goroutine"`
}

var startTime = time.Now()

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(HealthResponse{
		Status:  "ok",
		Version: version,
		Uptime:  time.Since(startTime).Round(time.Second).String(),
		GoVer:   runtime.Version(),
	})
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(InfoResponse{
		App:          "NexusMQ API",
		Version:      version,
		Environment:  getEnv("ENV", "development"),
		Port:         getEnv("PORT", "8080"),
		GoVersion:    runtime.Version(),
		NumCPU:       runtime.NumCPU(),
		NumGoroutine: runtime.NumGoroutine(),
	})
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"message": "🐳 API Go rodando em Docker!",
		"version": version,
	})
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}

func loggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %v", r.Method, r.URL.Path, time.Since(start))
	})
}

func main() {
	port := getEnv("PORT", "8080")

	mux := http.NewServeMux()
	mux.HandleFunc("/", rootHandler)
	mux.HandleFunc("/health", healthHandler)
	mux.HandleFunc("/info", infoHandler)

	handler := loggingMiddleware(mux)

	fmt.Printf("🐳 API v%s rodando em :%s\n", version, port)
	fmt.Printf("   ENV: %s\n", getEnv("ENV", "development"))
	fmt.Printf("   Go:  %s\n", runtime.Version())
	fmt.Printf("   CPU: %d cores\n", runtime.NumCPU())
	fmt.Println()
	fmt.Println("Endpoints:")
	fmt.Println("  GET /        — mensagem de boas-vindas")
	fmt.Println("  GET /health  — health check")
	fmt.Println("  GET /info    — informações do sistema")

	log.Fatal(http.ListenAndServe(":"+port, handler))
}
