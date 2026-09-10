package handler

import (
	"strconv"

	"vnfm-api/internal/middleware"
	"vnfm-api/internal/notify/service"
	"vnfm-api/pkg/errors"
	"vnfm-api/pkg/response"

	"github.com/gofiber/fiber/v3"
)

type NotifyHandler struct {
	svc *service.NotifyService
}

func NewNotifyHandler(svc *service.NotifyService) *NotifyHandler {
	return &NotifyHandler{svc: svc}
}

func (h *NotifyHandler) List(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	rows, total, appErr := h.svc.List(user.ID, limit, offset)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	unread, _ := h.svc.UnreadCount(user.ID)
	return response.OK(c, fiber.Map{"items": rows, "total": total, "unread": unread})
}

func (h *NotifyHandler) Unread(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	n, appErr := h.svc.UnreadCount(user.ID)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, fiber.Map{"unread": n})
}

func (h *NotifyHandler) MarkRead(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := strconv.ParseInt(c.Params("id"), 10, 64)
	if perr != nil {
		return response.Error(c, errors.ErrBadRequest("无效的 ID"))
	}
	if appErr := h.svc.MarkRead(user.ID, id); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已读")
}

func (h *NotifyHandler) MarkAllRead(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	if appErr := h.svc.MarkAllRead(user.ID); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已全部标为已读")
}
