package service

import (
	"strings"
	"time"

	clubrepo "vnfm-api/internal/club/repository"
	"vnfm-api/internal/event/dto"
	"vnfm-api/internal/event/model"
	"vnfm-api/internal/event/repository"
	notify "vnfm-api/internal/notify/service"
	"vnfm-api/pkg/errors"
)

type EventService struct {
	repo   *repository.EventRepository
	clubs  *clubrepo.ClubRepository
	notify *notify.NotifyService
}

func NewEventService(repo *repository.EventRepository, clubs *clubrepo.ClubRepository, n *notify.NotifyService) *EventService {
	return &EventService{repo: repo, clubs: clubs, notify: n}
}

func (s *EventService) List(viewerID int64, limit, offset int) ([]dto.EventView, int64, *errors.AppError) {
	rows, total, err := s.repo.List(limit, offset)
	if err != nil {
		return nil, 0, errors.ErrInternal("读取活动失败")
	}
	out := make([]dto.EventView, 0, len(rows))
	for i := range rows {
		out = append(out, s.toView(&rows[i], viewerID))
	}
	return out, total, nil
}

func (s *EventService) Get(id, viewerID int64) (*dto.EventView, *errors.AppError) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrNotFound("活动不存在")
	}
	v := s.toView(e, viewerID)
	return &v, nil
}

func (s *EventService) Create(userID int64, req *dto.UpsertEventRequest) (*dto.EventView, *errors.AppError) {
	if err := s.requireManage(userID, req.ClubID); err != nil {
		return nil, err
	}
	title := strings.TrimSpace(req.Title)
	if title == "" {
		return nil, errors.ErrBadRequest("请填写活动标题")
	}
	starts, err := time.Parse(time.RFC3339, req.StartsAt)
	if err != nil {
		return nil, errors.ErrBadRequest("开始时间格式无效")
	}
	now := time.Now()
	e := &model.Event{
		ClubID:      req.ClubID,
		Title:       title,
		Description: req.Description,
		Location:    req.Location,
		StartsAt:    starts,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
	e.EndsAt = parseOptTime(req.EndsAt)
	e.RegisterUntil = parseOptTime(req.RegisterUntil)
	if err := s.repo.Create(e); err != nil {
		return nil, errors.ErrInternal("创建活动失败")
	}
	v := s.toView(e, userID)
	return &v, nil
}

func (s *EventService) Update(userID, id int64, req *dto.UpsertEventRequest) (*dto.EventView, *errors.AppError) {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrNotFound("活动不存在")
	}
	if err := s.requireManage(userID, e.ClubID); err != nil {
		return nil, err
	}
	if strings.TrimSpace(req.Title) != "" {
		e.Title = strings.TrimSpace(req.Title)
	}
	e.Description = req.Description
	e.Location = req.Location
	if req.StartsAt != "" {
		starts, perr := time.Parse(time.RFC3339, req.StartsAt)
		if perr != nil {
			return nil, errors.ErrBadRequest("开始时间格式无效")
		}
		e.StartsAt = starts
	}
	e.EndsAt = parseOptTime(req.EndsAt)
	e.RegisterUntil = parseOptTime(req.RegisterUntil)
	e.UpdatedAt = time.Now()
	if err := s.repo.Save(e); err != nil {
		return nil, errors.ErrInternal("保存活动失败")
	}
	v := s.toView(e, userID)
	return &v, nil
}

func (s *EventService) Delete(userID, id int64) *errors.AppError {
	e, err := s.repo.FindByID(id)
	if err != nil {
		return errors.ErrNotFound("活动不存在")
	}
	if err := s.requireManage(userID, e.ClubID); err != nil {
		return err
	}
	if err := s.repo.Delete(id); err != nil {
		return errors.ErrInternal("删除失败")
	}
	return nil
}

func (s *EventService) Register(userID, eventID int64, join bool) *errors.AppError {
	e, err := s.repo.FindByID(eventID)
	if err != nil {
		return errors.ErrNotFound("活动不存在")
	}
	if locked(e) {
		return errors.ErrBadRequest("报名已截止")
	}
	if join {
		if _, ferr := s.repo.FindReg(eventID, userID); ferr == nil {
			return errors.ErrBadRequest("已经报名")
		}
		if err := s.repo.CreateReg(&model.Registration{EventID: eventID, UserID: userID, CreatedAt: time.Now()}); err != nil {
			return errors.ErrInternal("报名失败")
		}
		s.notify.Push(userID, "event_registered", "活动报名成功",
			"你已报名「"+e.Title+"」", "/events/"+itoa(e.ID), "event", e.ID)
		return nil
	}
	if err := s.repo.DeleteReg(eventID, userID); err != nil {
		return errors.ErrInternal("取消报名失败")
	}
	return nil
}

func (s *EventService) MyRegistrations(userID int64) ([]dto.EventView, *errors.AppError) {
	regs, err := s.repo.ListUserRegs(userID)
	if err != nil {
		return nil, errors.ErrInternal("读取报名失败")
	}
	out := make([]dto.EventView, 0, len(regs))
	for _, reg := range regs {
		e, ferr := s.repo.FindByID(reg.EventID)
		if ferr != nil {
			continue
		}
		out = append(out, s.toView(e, userID))
	}
	return out, nil
}

func (s *EventService) requireManage(userID, clubID int64) *errors.AppError {
	m, err := s.clubs.FindMembership(userID, clubID)
	if err != nil || m.Status != "active" || (m.Role != "manager" && m.Role != "representative") {
		return errors.ErrForbidden("需要同好会管理员或负责人权限")
	}
	return nil
}

func (s *EventService) toView(e *model.Event, viewerID int64) dto.EventView {
	v := dto.EventView{
		ID:          e.ID,
		ClubID:      e.ClubID,
		ClubName:    s.repo.ClubName(e.ClubID),
		Title:       e.Title,
		Description: e.Description,
		Location:    e.Location,
		CoverKey:    e.CoverKey,
		StartsAt:    e.StartsAt.Format(time.RFC3339),
		Locked:      locked(e),
	}
	if e.EndsAt != nil {
		v.EndsAt = e.EndsAt.Format(time.RFC3339)
	}
	if e.RegisterUntil != nil {
		v.RegisterUntil = e.RegisterUntil.Format(time.RFC3339)
	}
	if viewerID > 0 {
		if _, err := s.repo.FindReg(e.ID, viewerID); err == nil {
			v.Registered = true
		}
		if s.requireManage(viewerID, e.ClubID) == nil {
			v.CanManage = true
		}
	}
	return v
}

func locked(e *model.Event) bool {
	now := time.Now()
	if e.RegisterUntil != nil {
		return now.After(*e.RegisterUntil)
	}
	return now.After(e.StartsAt)
}

func parseOptTime(raw string) *time.Time {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil
	}
	t, err := time.Parse(time.RFC3339, raw)
	if err != nil {
		return nil
	}
	return &t
}

func itoa(n int64) string {
	const digits = "0123456789"
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var b [20]byte
	i := len(b)
	for n > 0 {
		i--
		b[i] = digits[n%10]
		n /= 10
	}
	if neg {
		i--
		b[i] = '-'
	}
	return string(b[i:])
}
