package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/luckydiss/ocr_tgapp/internal/config"
	"github.com/luckydiss/ocr_tgapp/internal/gemini"
	"github.com/luckydiss/ocr_tgapp/internal/handler"
)

func main() {
	cfg := config.Load()

	ctx := context.Background()
	geminiService, err := gemini.NewGeminiService(ctx, cfg.GeminiAPIKey, cfg.GeminiAPIEndpoint)
	if err != nil {
		log.Fatalf("Failed to initialize Gemini service: %v", err)
	}

	mux := http.NewServeMux()

	extractHandler := handler.NewExtractHandler(geminiService)
	mux.Handle("/api/v1/extract", extractHandler)

	fs := http.FileServer(http.Dir("web"))
	mux.Handle("/", fs)

	var h http.Handler = mux
	h = handler.CORS(cfg.AllowedOrigins)(h)
	h = handler.Logging(h)
	h = handler.Recovery(h)

	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      h,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 60 * time.Second,
		IdleTimeout:  120 * time.Second,
	}

	go func() {
		log.Printf("Server starting on port %s", cfg.Port)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exited")
}
