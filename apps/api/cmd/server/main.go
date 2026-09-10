package main

import (
	"log/slog"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"vnfm-api/internal/app"
	"vnfm-api/internal/infrastructure/database"
	"vnfm-api/pkg/config"
	"vnfm-api/pkg/health"
	"vnfm-api/pkg/logger"

	"github.com/joho/godotenv"
	"gorm.io/gorm"
)

func main() {
	_ = godotenv.Load()

	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "3711"
	}
	port, _ := strconv.Atoi(serverPort)
	health.MaybeProbe(port, "/healthz")

	cfg, err := config.Load()
	if err != nil {
		slog.Error("load config", "error", err)
		os.Exit(1)
	}
	logger.Init(cfg.Server.Mode)

	rdb := app.RedisFromConfig(cfg)

	var db *gorm.DB
	if cfg.Database.URL != "" {
		db, err = database.Open(cfg.Database, cfg.Server.Mode)
		if err != nil {
			slog.Error("open database", "error", err)
			os.Exit(1)
		}
	}

	application := app.New(cfg, db, rdb)

	go func() {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit
		slog.Info("shutting down")
		_ = application.Fiber.Shutdown()
	}()

	addr := ":" + cfg.Server.Port
	slog.Info("listening", "addr", addr)
	if err := application.Fiber.Listen(addr); err != nil {
		slog.Error("listen", "error", err)
		os.Exit(1)
	}
}
