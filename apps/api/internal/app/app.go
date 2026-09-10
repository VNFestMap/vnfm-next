package app

import (
	"log/slog"

	"vnfm-api/internal/infrastructure/cache"
	"vnfm-api/internal/middleware"
	userHandler "vnfm-api/internal/user/handler"
	"vnfm-api/internal/user/oauth"
	userRepo "vnfm-api/internal/user/repository"
	userService "vnfm-api/internal/user/service"
	"vnfm-api/pkg/config"
	"vnfm-api/pkg/errors"
	"vnfm-api/pkg/response"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/cors"
	"github.com/gofiber/fiber/v3/middleware/recover"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

type App struct {
	Fiber  *fiber.App
	DB     *gorm.DB
	Redis  *redis.Client
	Config *config.Config
}

func New(cfg *config.Config, db *gorm.DB, rdb *redis.Client) *App {
	fiberApp := fiber.New(fiber.Config{
		ErrorHandler: globalErrorHandler,
		BodyLimit:    10 * 1024 * 1024,
	})
	fiberApp.Use(recover.New())
	fiberApp.Use(cors.New(middleware.CORS(cfg.CORS.AllowOrigins)))

	app := &App{Fiber: fiberApp, DB: db, Redis: rdb, Config: cfg}
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

	if a.DB == nil || a.Redis == nil || a.Config.OIDC.ClientID == "" {
		slog.Warn("auth routes skipped (need DATABASE_URL, redis, OIDC_CLIENT_ID)")
		return
	}

	oauthClient := oauth.NewClient(a.Config.OIDC)
	users := userRepo.NewUserRepository(a.DB)
	authSvc := userService.NewAuthService(users, a.Redis, oauthClient, a.Config.OIDC)
	oauthH := userHandler.NewOAuthHandler(authSvc, a.Config.Server.Mode == "prod")

	authMw := middleware.Auth(a.Redis, oauthClient)
	v1.Post("/auth/oauth/callback", oauthH.Callback)
	v1.Post("/auth/logout", oauthH.Logout)
	v1.Get("/auth/me", authMw, oauthH.Me)
	v1.Patch("/auth/me", authMw, oauthH.UpdateMe)
}

func globalErrorHandler(c fiber.Ctx, err error) error {
	if appErr, ok := err.(*errors.AppError); ok {
		return response.Error(c, appErr)
	}
	slog.Error("unhandled error", "error", err.Error(), "path", c.Path(), "method", c.Method())
	return response.Error(c, errors.ErrInternal("服务器内部错误"))
}

func RedisFromConfig(cfg *config.Config) *redis.Client {
	return cache.Open(cache.RedisConfig{
		Host:     cfg.Redis.Host,
		Port:     cfg.Redis.Port,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.DB,
	})
}
