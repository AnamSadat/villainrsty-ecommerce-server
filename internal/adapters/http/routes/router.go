package routes

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"villainrsty-ecommerce-server/internal/adapters/http/httpx"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

var startTime = time.Now()

func New() *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		w.Write([]byte("Hello World"))
	})

	r.Get("/health", healtHandler)

	return r
}

func healtHandler(w http.ResponseWriter, _ *http.Request) {
	requestStart := time.Now()

	uptime := time.Since(startTime)
	hours := int(uptime.Hours())
	minutes := int(uptime.Minutes()) % 60
	seconds := int(uptime.Seconds()) % 60

	var mem runtime.MemStats
	runtime.ReadMemStats(&mem)

	latency := time.Since(requestStart)

	httpx.JSON(w, http.StatusOK, map[string]any{
		"status":    "ok",
		"timestamp": time.Now().Format(time.RFC3339),
		"uptime":    fmt.Sprintf("%dh %dm %ds", hours, minutes, seconds),
		"latency":   fmt.Sprintf("%.2f ms", float64(latency.Microseconds())/1000),
		"memory": map[string]string{
			"heapUsed":  fmt.Sprintf("%d MB", mem.HeapAlloc/1024/1024),
			"heapTotal": fmt.Sprintf("%d MB", mem.HeapSys/1024/1024),
		},
	})
}
