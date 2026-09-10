package handler

import (
	"strconv"

	"vnfm-api/internal/event/dto"
	"vnfm-api/internal/event/service"
	"vnfm-api/internal/middleware"
	"vnfm-api/pkg/errors"
	"vnfm-api/pkg/response"
	"vnfm-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

type EventHandler struct {
	svc *service.EventService
}

func NewEventHandler(svc *service.EventService) *EventHandler {
	return &EventHandler{svc: svc}
}

func viewerID(c fiber.Ctx) int64 {
	if u := middleware.GetUser(c); u != nil {
		return u.ID
	}
	return 0
}

func (h *EventHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	rows, total, err := h.svc.List(viewerID(c), limit, offset)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Paginated(c, rows, total)
}

func (h *EventHandler) Get(c fiber.Ctx) error {
	id, perr := strconv.ParseInt(c.Params("id"), 10, 64)
	if perr != nil {
		return response.Error(c, errors.ErrBadRequest("无效的 ID"))
	}
	row, err := h.svc.Get(id, viewerID(c))
	if err != nil {
		return response.Error(c, err)
	}
	return response.OK(c, row)
}

func (h *EventHandler) Create(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	var req dto.UpsertEventRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	row, appErr := h.svc.Create(user.ID, &req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, row)
}

func (h *EventHandler) Update(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := strconv.ParseInt(c.Params("id"), 10, 64)
	if perr != nil {
		return response.Error(c, errors.ErrBadRequest("无效的 ID"))
	}
	var req dto.UpsertEventRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	row, appErr := h.svc.Update(user.ID, id, &req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, row)
}

func (h *EventHandler) Delete(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := strconv.ParseInt(c.Params("id"), 10, 64)
	if perr != nil {
		return response.Error(c, errors.ErrBadRequest("无效的 ID"))
	}
	if appErr := h.svc.Delete(user.ID, id); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已删除")
}

func (h *EventHandler) Register(c fiber.Ctx) error {
	return h.setReg(c, true)
}

func (h *EventHandler) Unregister(c fiber.Ctx) error {
	return h.setReg(c, false)
}

func (h *EventHandler) setReg(c fiber.Ctx, join bool) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := strconv.ParseInt(c.Params("id"), 10, 64)
	if perr != nil {
		return response.Error(c, errors.ErrBadRequest("无效的 ID"))
	}
	if appErr := h.svc.Register(user.ID, id, join); appErr != nil {
		return response.Error(c, appErr)
	}
	if join {
		return response.OKMessage(c, "报名成功")
	}
	return response.OKMessage(c, "已取消报名")
}

func (h *EventHandler) Mine(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	rows, appErr := h.svc.MyRegistrations(user.ID)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, rows)
}
