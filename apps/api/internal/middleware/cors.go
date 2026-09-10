package middleware

import (
	"strings"

	"github.com/gofiber/fiber/v3/middleware/cors"
)

func CORS(allowOrigins string) cors.Config {
	return cors.Config{
		AllowOrigins:     strings.Split(allowOrigins, ","),
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		AllowCredentials: true,
		MaxAge:           86400,
	}
}
