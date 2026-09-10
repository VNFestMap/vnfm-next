package app

import (
	"log/slog"

	"vnfm-api/internal/middleware"
	"vnfm-api/pkg/config"
	"vnfm-api/pkg/errors"
	"vnfm-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"gorm.io/gorm"
)

type App struct {
	Fiber  *fiber.App
	DB     *gorm.DB
	Config *config.Config
}

func New(cfg *config.Config, db *gorm.DB) *App {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: globalErrorHandler,
		BodyLimit:    10 * 1024 * 1024,
	})
	fiberApp.Use(recover.New())
	fiberApp.Use(cors.New(middleware.CORS(cfg.CORS.AllowOrigins)))

	app := &App{Fiber: fiberApp, DB: db, Config: cfg}
	app.setupRoutes()
	return app
}

func (a *App) setupRoutes() {
	a.Fiber.Get("/healthz", func(c fiber.Ctx) error {
		return c.SendString("ok")
	})

	v1 := a.Fiber.Group("/api/v1")
	v1.Get("/health", func(c fiber.Ctx) error {
		return response.OK(c, fiber.Map{"ok": true})
	})
}

func globalErrorHandler(c fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return response.Error(c, appErr)
	}
	slog.Error("unhandled error", "error", err.Error(), "path", c.Path(), "method", c.Method())
	return response.Error(c, errors.ErrInternal("服务器内部错误"))
}
