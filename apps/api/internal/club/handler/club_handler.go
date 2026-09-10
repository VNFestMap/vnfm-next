package handler

import (
	"strconv"

	"vnfm-api/internal/club/dto"
	"vnfm-api/internal/club/service"
	"vnfm-api/internal/middleware"
	"vnfm-api/pkg/errors"
	"vnfm-api/pkg/response"
	"vnfm-api/pkg/utils"

	"github.com/gofiber/fiber/v3"
)

type ClubHandler struct {
	svc *service.ClubService
}

func NewClubHandler(svc *service.ClubService) *ClubHandler {
	return &ClubHandler{svc: svc}
}

func viewerID(c fiber.Ctx) int64 {
	if u := middleware.GetUser(c); u != nil {
		return u.ID
	}
	return 0
}

func paramID(c fiber.Ctx, name string) (int64, *errors.AppError) {
	id, err := strconv.ParseInt(c.Params(name), 10, 64)
	if err != nil || id <= 0 {
		return 0, errors.ErrBadRequest("无效的 ID")
	}
	return id, nil
}

func (h *ClubHandler) List(c fiber.Ctx) error {
	limit, _ := strconv.Atoi(c.Query("limit", "50"))
	offset, _ := strconv.Atoi(c.Query("offset", "0"))
	rows, total, err := h.svc.List(viewerID(c), c.Query("q"), c.Query("country"), c.Query("province"), c.Query("type"), limit, offset)
	if err != nil {
		return response.Error(c, err)
	}
	return response.Paginated(c, rows, total)
}

func (h *ClubHandler) Get(c fiber.Ctx) error {
	id, err := paramID(c, "id")
	if err != nil {
		return response.Error(c, err)
	}
	row, appErr := h.svc.Get(id, viewerID(c))
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, row)
}

func (h *ClubHandler) Create(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	var req dto.CreateClubRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	row, appErr := h.svc.Create(user.ID, &req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, row)
}

func (h *ClubHandler) Update(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	var req dto.UpdateClubRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	row, appErr := h.svc.Update(user.ID, id, &req)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, row)
}

func (h *ClubHandler) Apply(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	var req dto.ApplyRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	if appErr := h.svc.Apply(user.ID, id, &req); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "申请已提交")
}

func (h *ClubHandler) Approve(c fiber.Ctx) error {
	return h.review(c, true)
}

func (h *ClubHandler) Reject(c fiber.Ctx) error {
	return h.review(c, false)
}

func (h *ClubHandler) review(c fiber.Ctx, ok bool) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	if appErr := h.svc.Approve(user.ID, id, ok); appErr != nil {
		return response.Error(c, appErr)
	}
	if ok {
		return response.OKMessage(c, "已通过")
	}
	return response.OKMessage(c, "已拒绝")
}

func (h *ClubHandler) ChangeRole(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	var req dto.ChangeRoleRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	if appErr := h.svc.ChangeRole(user.ID, id, req.Role); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "角色已更新")
}

func (h *ClubHandler) Kick(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	if appErr := h.svc.Kick(user.ID, id); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已移出")
}

func (h *ClubHandler) Leave(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	if appErr := h.svc.Leave(user.ID, id); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已退出")
}

func (h *ClubHandler) Transfer(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	if appErr := h.svc.Transfer(user.ID, id); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已转让负责人")
}

func (h *ClubHandler) Members(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	rows, appErr := h.svc.Members(user.ID, id)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, rows)
}

func (h *ClubHandler) Mine(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	rows, appErr := h.svc.MyMemberships(user.ID)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, rows)
}

func (h *ClubHandler) Pending(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	rows, appErr := h.svc.PendingForManager(user.ID)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, rows)
}

func (h *ClubHandler) CreateCode(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	var req dto.CreateCodeRequest
	_ = utils.BindJSON(c, &req)
	row, appErr := h.svc.CreateCode(user.ID, id, req.MaxUses, req.Hours)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, row)
}

func (h *ClubHandler) ListCodes(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	rows, appErr := h.svc.ListCodes(user.ID, id)
	if appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OK(c, rows)
}

func (h *ClubHandler) RevokeCode(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	id, perr := paramID(c, "id")
	if perr != nil {
		return response.Error(c, perr)
	}
	if appErr := h.svc.RevokeCode(user.ID, id); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "已吊销")
}

func (h *ClubHandler) Redeem(c fiber.Ctx) error {
	user, err := middleware.MustGetUser(c)
	if err != nil {
		return response.Error(c, err)
	}
	var req dto.RedeemRequest
	if bindErr := utils.BindJSON(c, &req); bindErr != nil {
		return response.Error(c, bindErr)
	}
	if appErr := h.svc.Redeem(user.ID, req.Code); appErr != nil {
		return response.Error(c, appErr)
	}
	return response.OKMessage(c, "加入成功")
}
