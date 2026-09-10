package app

import (
	"log/slog"

	clubHandler "vnfm-api/internal/club/handler"
	clubRepo "vnfm-api/internal/club/repository"
	clubService "vnfm-api/internal/club/service"
	eventHandler "vnfm-api/internal/event/handler"
	eventRepo "vnfm-api/internal/event/repository"
	eventService "vnfm-api/internal/event/service"
	"vnfm-api/internal/infrastructure/cache"
	"vnfm-api/internal/middleware"
	notifyHandler "vnfm-api/internal/notify/handler"
	notifyRepo "vnfm-api/internal/notify/repository"
	notifyService "vnfm-api/internal/notify/service"
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

	nRepo := notifyRepo.NewNotifyRepository(a.DB)
	nSvc := notifyService.NewNotifyService(nRepo)
	nH := notifyHandler.NewNotifyHandler(nSvc)

	cRepo := clubRepo.NewClubRepository(a.DB)
	cSvc := clubService.NewClubService(cRepo, nSvc)
	cH := clubHandler.NewClubHandler(cSvc)

	eRepo := eventRepo.NewEventRepository(a.DB)
	eSvc := eventService.NewEventService(eRepo, cRepo, nSvc)
	eH := eventHandler.NewEventHandler(eSvc)

	authMw := middleware.Auth(a.Redis, oauthClient)
	optMw := middleware.OptionalAuth(a.Redis, oauthClient)

	v1.Post("/auth/oauth/callback", oauthH.Callback)
	v1.Post("/auth/logout", oauthH.Logout)
	v1.Get("/auth/me", authMw, oauthH.Me)
	v1.Patch("/auth/me", authMw, oauthH.UpdateMe)

	v1.Get("/notifications", authMw, nH.List)
	v1.Get("/notifications/unread", authMw, nH.Unread)
	v1.Post("/notifications/read-all", authMw, nH.MarkAllRead)
	v1.Post("/notifications/:id/read", authMw, nH.MarkRead)

	v1.Get("/clubs", optMw, cH.List)
	v1.Get("/clubs/regions", optMw, cH.Regions)
	v1.Get("/clubs/:id", optMw, cH.Get)
	v1.Post("/clubs", authMw, cH.Create)
	v1.Patch("/clubs/:id", authMw, cH.Update)
	v1.Post("/clubs/:id/apply", authMw, cH.Apply)
	v1.Get("/clubs/:id/members", authMw, cH.Members)
	v1.Get("/clubs/:id/codes", authMw, cH.ListCodes)
	v1.Post("/clubs/:id/codes", authMw, cH.CreateCode)
	v1.Post("/clubs/:id/leave", authMw, cH.Leave)
	v1.Post("/codes/redeem", authMw, cH.Redeem)
	v1.Post("/codes/:id/revoke", authMw, cH.RevokeCode)
	v1.Get("/me/memberships", authMw, cH.Mine)
	v1.Get("/me/applications", authMw, cH.Pending)
	v1.Post("/memberships/:id/approve", authMw, cH.Approve)
	v1.Post("/memberships/:id/reject", authMw, cH.Reject)
	v1.Post("/memberships/:id/role", authMw, cH.ChangeRole)
	v1.Post("/memberships/:id/kick", authMw, cH.Kick)
	v1.Post("/memberships/:id/transfer", authMw, cH.Transfer)

	v1.Get("/events", optMw, eH.List)
	v1.Get("/events/:id", optMw, eH.Get)
	v1.Post("/events", authMw, eH.Create)
	v1.Patch("/events/:id", authMw, eH.Update)
	v1.Delete("/events/:id", authMw, eH.Delete)
	v1.Post("/events/:id/register", authMw, eH.Register)
	v1.Post("/events/:id/unregister", authMw, eH.Unregister)
	v1.Get("/me/registrations", authMw, eH.Mine)
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
