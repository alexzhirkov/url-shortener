package main

import (
	"github.com/alexzhirkov/url-shortener/internal/config"
	"github.com/alexzhirkov/url-shortener/internal/lib/logger/sl"
	"github.com/alexzhirkov/url-shortener/internal/storage/sqlite"
	"log"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
)

const (
	envLocal = "local"
	envDev   = "dev"
	envProd  = "prod"
)

func main() {
	if err := run(); err != nil {
		log.Fatal(err)
	}
	os.Exit(0)
}

func run() error {
	//init config
	cfg := config.MustLoad()

	//init logger
	logger := setupLogger(cfg.Env)

	//logger.With(slog.String("env", cfg.Env)).Info("info message")
	logger.Info("start url-shortener", slog.String("env", cfg.Env))
	//logger.Debug("debug message")

	//todo: init storage: sqlite
	storage, err := sqlite.New(cfg.StoragePath)
	if err != nil {
		logger.Error("storage initialization failed", sl.Err(err))
		os.Exit(1)
	}
	_ = storage

	//todo: init router: chi

	//todo: run server

	// listen to OS signals and gracefully shutdown HTTP server
	stopped := make(chan struct{})
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
		<-sigint
		//ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		//defer cancel()
		//if err := srv.Shutdown(ctx); err != nil {
		//	log.Printf("HTTP Server Shutdown Error: %v", err)
		//}
		close(stopped)
	}()

	//log.Printf("Starting HTTP server on %s", cfg.HTTPAddr)

	// start HTTP server
	//if err := srv.ListenAndServe(); err != http.ErrServerClosed {
	//	log.Fatalf("HTTP server ListenAndServe Error: %v", err)
	//}

	<-stopped

	log.Printf("Have a nice day!")

	return nil
}

func setupLogger(env string) *slog.Logger {
	var logger *slog.Logger

	switch env {
	case envLocal:
		logger = slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envDev:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelDebug}))
	case envProd:
		logger = slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	}

	return logger
}
