package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"villainrsty-ecommerce-server/internal/adapters/http/router"
	"villainrsty-ecommerce-server/internal/app"
	"villainrsty-ecommerce-server/internal/config"

	"github.com/joho/godotenv"
)

func main() {
	// Load .env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found")
	}

	cfg := config.MustLoad()
	db := config.ConnectDB(cfg.DatabaseUrl)
	defer db.Close()

	container := app.New(cfg, db)

	r := router.New(container)

	srv := &http.Server{
		Addr:              cfg.Addr,
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		log.Printf("🚀 Server running on http://localhost%s", cfg.Addr)
		if err := srv.ListenAndServe(); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
	<-stop

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Printf("shutdown error: %v", err)
	}

	log.Println("shutdown compolete")
}
