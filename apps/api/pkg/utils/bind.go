package utils

import (
	"vnfm-api/pkg/errors"

	"github.com/gofiber/fiber/v3"
)

func BindJSON(c fiber.Ctx, dst any) *errors.AppError {
	if err := c.Bind().Body(dst); err != nil {
		return errors.ErrBadRequest("请求格式错误")
	}
	return nil
}
