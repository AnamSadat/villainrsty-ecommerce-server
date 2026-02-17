package router

import (
	"fmt"
	"net/http"
	"runtime"
	"time"

	"villainrsty-ecommerce-server/internal/adapters/http/auth/routes"
	"villainrsty-ecommerce-server/internal/adapters/http/httpx"
	"villainrsty-ecommerce-server/internal/app"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"github.com/swaggest/swgui/v5emb"
)

var startTime = time.Now()

func New(container *app.Container) *chi.Mux {
	r := chi.NewRouter()

	r.Use(middleware.Logger, middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "DELETE", "PUT", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	}))

	r.Get("/", func(w http.ResponseWriter, _ *http.Request) {
		httpx.JSON(w, http.StatusOK, map[string]any{
			"success": true,
			"message": "Server API Villainrsty Ecommerce is running",
		})
	})

	r.Get("/health", healtHandler)
	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		httpx.ErrorWithDetails(
			w,
			http.StatusNotFound,
			"Route not found",
			"NOT_FOUND",
			map[string]any{
				"path":   r.URL.Path,
				"method": r.Method,
			})
	})

	r.Get("/openapi.yml", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "./openapi.yml")
	})

	// Pasang Swagger UI v5 (Support v3.0 & v3.1)
	// New(Judul, Path_ke_YAML, Path_di_Browser)
	r.Mount("/docs", v5emb.New(
		"Villainrsty API",
		"/openapi.yml",
		"/docs",
	))

	routes.RegisterRoute(r, container.AuthHandler)

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
