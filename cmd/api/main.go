package main

import (
	"context"
	"errors"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"example/internal/config"
	"example/internal/controller"
	"example/pkg/pgxpool"

	userCtrl "example/internal/controller/user"
	txmanager "example/internal/infra/postgres/tx-manager"
	userRepo "example/internal/infra/postgres/user"
	userSvc "example/internal/service/user"
)

func main() {
	log := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	cfg, err := config.GetConfig()
	if err != nil {
		log.Error("failed to load config", "error", err)
		os.Exit(1)
	}

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	pool, err := connectDB(ctx, cfg)
	if err != nil {
		log.Error("failed to connect to postgres", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	userRepository := userRepo.New(pool)
	txm := txmanager.NewRepository(pool)

	userService := userSvc.New(userRepository, txm)

	userHandler := userCtrl.New(userService, log)

	ctrl := controller.New(userHandler, log)
	router := ctrl.Setup()
	server := &http.Server{
		Addr:         ":" + cfg.HttpPort,
		Handler:      router,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
		IdleTimeout:  60 * time.Second,
	}
	go func() {
		log.Info("starting HTTP server", "port", cfg.HttpPort)
		if err := server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Error("HTTP server failed", "error", err)
			os.Exit(1)
		}
	}()
	<-ctx.Done()
	log.Info("shutdown signal received")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	if err := server.Shutdown(shutdownCtx); err != nil {
		log.Error("failed to shutdown server", "error", err)
		os.Exit(1)
	}

	log.Info("server stopped")
}

func connectDB(ctx context.Context, cfg *config.Config) (*pgxpool.Pool, error) {
	if cfg.DatabaseURL != "" {
		return pgxpool.NewPoolFromDSN(ctx, cfg.DatabaseURL)
	}
	return pgxpool.NewPool(ctx, pgxpool.Config{
		Host:     cfg.Postgres.Host,
		Port:     cfg.Postgres.Port,
		UserName: cfg.Postgres.UserName,
		Password: cfg.Postgres.Password,
		DBName:   cfg.Postgres.DBName,
	})
}
