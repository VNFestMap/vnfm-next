package service

import (
	"crypto/rand"
	"encoding/hex"
	"strconv"
	"strings"
	"time"

	"vnfm-api/internal/club/dto"
	"vnfm-api/internal/club/model"
	"vnfm-api/internal/club/repository"
	notify "vnfm-api/internal/notify/service"
	"vnfm-api/pkg/errors"

	"gorm.io/gorm"
)

var manageRoles = map[string]bool{"manager": true, "representative": true}
var formalRoles = map[string]bool{"member": true, "manager": true, "representative": true}

type ClubService struct {
	repo   *repository.ClubRepository
	notify *notify.NotifyService
}

func NewClubService(repo *repository.ClubRepository, n *notify.NotifyService) *ClubService {
	return &ClubService{repo: repo, notify: n}
}

func (s *ClubService) List(viewerID int64, q, country, province, typ string, limit, offset int) ([]dto.ClubSummary, int64, *errors.AppError) {
	rows, total, err := s.repo.List(q, country, province, typ, limit, offset)
	if err != nil {
		return nil, 0, errors.ErrInternal("读取同好会失败")
	}
	out := make([]dto.ClubSummary, 0, len(rows))
	for _, c := range rows {
		out = append(out, s.toSummary(&c, viewerID, false))
	}
	return out, total, nil
}

func (s *ClubService) Regions(country string) ([]dto.RegionCount, *errors.AppError) {
	if country != "japan" {
		country = "china"
	}
	rows, err := s.repo.RegionCounts(country)
	if err != nil {
		return nil, errors.ErrInternal("读取地区统计失败")
	}
	out := make([]dto.RegionCount, 0, len(rows))
	for _, r := range rows {
		out = append(out, dto.RegionCount{Key: r.Key, Count: r.Count})
	}
	return out, nil
}

func (s *ClubService) Get(id, viewerID int64) (*dto.ClubSummary, *errors.AppError) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		return nil, errors.ErrNotFound("同好会不存在")
	}
	sum := s.toSummary(c, viewerID, true)
	return &sum, nil
}

func (s *ClubService) Create(userID int64, req *dto.CreateClubRequest) (*dto.ClubSummary, *errors.AppError) {
	name := strings.TrimSpace(req.Name)
	if name == "" {
		return nil, errors.ErrBadRequest("请填写同好会名称")
	}
	country := req.Country
	if country != "japan" {
		country = "china"
	}
	typ := req.Type
	if typ == "" {
		typ = "school"
	}
	hidden := true
	if req.ContactHidden != nil {
		hidden = *req.ContactHidden
	}
	now := time.Now()
	club := &model.Club{
		Country:       country,
		Name:          name,
		School:        strings.TrimSpace(req.School),
		Province:      strings.TrimSpace(req.Province),
		Prefecture:    strings.TrimSpace(req.Prefecture),
		City:          strings.TrimSpace(req.City),
		Type:          typ,
		Info:          req.Info,
		Remark:        req.Remark,
		ContactHidden: hidden,
		CreatedBy:     userID,
		CreatedAt:     now,
		UpdatedAt:     now,
	}
	if err := s.repo.Create(club); err != nil {
		return nil, errors.ErrInternal("创建同好会失败")
	}
	m := &model.Membership{
		UserID:   userID,
		ClubID:   club.ID,
		Country:  country,
		Role:     "representative",
		Status:   "active",
		JoinedAt: now,
	}
	if err := s.repo.CreateMembership(m); err != nil {
		return nil, errors.ErrInternal("写入会籍失败")
	}
	sum := s.toSummary(club, userID, true)
	return &sum, nil
}

func (s *ClubService) Update(userID, clubID int64, req *dto.UpdateClubRequest) (*dto.ClubSummary, *errors.AppError) {
	if _, err := s.requireManage(userID, clubID); err != nil {
		return nil, err
	}
	club, ferr := s.repo.FindByID(clubID)
	if ferr != nil {
		return nil, errors.ErrNotFound("同好会不存在")
	}
	if req.Name != nil && strings.TrimSpace(*req.Name) != "" {
		club.Name = strings.TrimSpace(*req.Name)
	}
	if req.School != nil {
		club.School = strings.TrimSpace(*req.School)
	}
	if req.Province != nil {
		club.Province = strings.TrimSpace(*req.Province)
	}
	if req.Prefecture != nil {
		club.Prefecture = strings.TrimSpace(*req.Prefecture)
	}
	if req.City != nil {
		club.City = strings.TrimSpace(*req.City)
	}
	if req.Type != nil && *req.Type != "" {
		club.Type = *req.Type
	}
	if req.Info != nil {
		club.Info = *req.Info
	}
	if req.Remark != nil {
		club.Remark = *req.Remark
	}
	if req.ContactHidden != nil {
		club.ContactHidden = *req.ContactHidden
	}
	club.UpdatedAt = time.Now()
	if err := s.repo.Update(club); err != nil {
		return nil, errors.ErrInternal("保存失败")
	}
	sum := s.toSummary(club, userID, true)
	return &sum, nil
}

func (s *ClubService) Apply(userID, clubID int64, req *dto.ApplyRequest) *errors.AppError {
	club, err := s.repo.FindByID(clubID)
	if err != nil {
		return errors.ErrNotFound("同好会不存在")
	}
	existing, _ := s.repo.FindMembership(userID, clubID)
	if existing != nil && existing.Status == "active" {
		return errors.ErrBadRequest("你已经是该同好会成员")
	}
	if existing != nil && existing.Status == "pending" {
		return errors.ErrBadRequest("申请正在审核中")
	}
	role := req.ApplyRole
	if role != "external" {
		role = "member"
	}
	now := time.Now()
	if existing != nil {
		existing.Status = "pending"
		existing.ApplyRole = role
		existing.ApplyReason = req.ApplyReason
		existing.ContactAccount = req.ContactAccount
		existing.LeftAt = nil
		existing.JoinedAt = now
		if err := s.repo.SaveMembership(existing); err != nil {
			return errors.ErrInternal("提交申请失败")
		}
	} else {
		m := &model.Membership{
			UserID:         userID,
			ClubID:         clubID,
			Country:        club.Country,
			Role:           "member",
			Status:         "pending",
			ContactAccount: req.ContactAccount,
			ApplyReason:    req.ApplyReason,
			ApplyRole:      role,
			JoinedAt:       now,
		}
		if err := s.repo.CreateMembership(m); err != nil {
			return errors.ErrInternal("提交申请失败")
		}
	}
	managers, _ := s.repo.Managers(clubID)
	for _, mgr := range managers {
		s.notify.Push(mgr.UserID, "join_request", "新的加入申请",
			"有人申请加入 "+club.Name, "/user", "club", clubID)
	}
	return nil
}

func (s *ClubService) Approve(actorID, membershipID int64, approve bool) *errors.AppError {
	m, err := s.repo.FindMembershipByID(membershipID)
	if err != nil {
		return errors.ErrNotFound("申请不存在")
	}
	if _, err := s.requireManage(actorID, m.ClubID); err != nil {
		return err
	}
	now := time.Now()
	m.ReviewedBy = &actorID
	m.ReviewedAt = &now
	club, _ := s.repo.FindByID(m.ClubID)
	name := ""
	if club != nil {
		name = club.Name
	}
	if approve {
		m.Status = "active"
		m.Role = m.ApplyRole
		if m.Role == "" {
			m.Role = "member"
		}
		s.notify.Push(m.UserID, "join_approved", "加入申请已通过",
			"你已加入 "+name, "/clubs/"+itoa(m.ClubID), "club_membership", m.ID)
	} else {
		m.Status = "rejected"
		s.notify.Push(m.UserID, "join_rejected", "加入申请未通过",
			name+" 未通过你的加入申请", "/clubs/"+itoa(m.ClubID), "club_membership", m.ID)
	}
	if err := s.repo.SaveMembership(m); err != nil {
		return errors.ErrInternal("审核失败")
	}
	return nil
}

func (s *ClubService) ChangeRole(actorID, membershipID int64, role string) *errors.AppError {
	if !formalRoles[role] && role != "external" {
		return errors.ErrBadRequest("无效角色")
	}
	m, err := s.repo.FindMembershipByID(membershipID)
	if err != nil {
		return errors.ErrNotFound("会籍不存在")
	}
	actor, appErr := s.requireManage(actorID, m.ClubID)
	if appErr != nil {
		return appErr
	}
	if actor.Role != "representative" && (role == "representative" || m.Role == "representative") {
		return errors.ErrForbidden("只有负责人可以变更负责人角色")
	}
	m.Role = role
	if err := s.repo.SaveMembership(m); err != nil {
		return errors.ErrInternal("改角色失败")
	}
	club, _ := s.repo.FindByID(m.ClubID)
	title := "角色已更新"
	msg := "你在同好会的角色已变更为 " + role
	if club != nil {
		msg = "你在 " + club.Name + " 的角色已变更为 " + role
	}
	s.notify.Push(m.UserID, "role_changed", title, msg, "/user", "club_membership", m.ID)
	return nil
}

func (s *ClubService) Kick(actorID, membershipID int64) *errors.AppError {
	m, err := s.repo.FindMembershipByID(membershipID)
	if err != nil {
		return errors.ErrNotFound("会籍不存在")
	}
	if _, appErr := s.requireManage(actorID, m.ClubID); appErr != nil {
		return appErr
	}
	if m.Role == "representative" {
		return errors.ErrForbidden("不能移出负责人，请先转让")
	}
	now := time.Now()
	m.Status = "left"
	m.LeftAt = &now
	if err := s.repo.SaveMembership(m); err != nil {
		return errors.ErrInternal("移出失败")
	}
	club, _ := s.repo.FindByID(m.ClubID)
	name := ""
	if club != nil {
		name = club.Name
	}
	s.notify.Push(m.UserID, "member_kicked", "你已被移出同好会",
		"你已被移出 "+name, "/", "club_membership", m.ID)
	return nil
}

func (s *ClubService) Leave(userID, clubID int64) *errors.AppError {
	m, err := s.repo.FindMembership(userID, clubID)
	if err != nil {
		return errors.ErrNotFound("会籍不存在")
	}
	if m.Role == "representative" && m.Status == "active" {
		return errors.ErrForbidden("负责人请先转让后再退出")
	}
	now := time.Now()
	m.Status = "left"
	m.LeftAt = &now
	if err := s.repo.SaveMembership(m); err != nil {
		return errors.ErrInternal("退出失败")
	}
	return nil
}

func (s *ClubService) Transfer(actorID, membershipID int64) *errors.AppError {
	target, err := s.repo.FindMembershipByID(membershipID)
	if err != nil || target.Status != "active" {
		return errors.ErrNotFound("目标会籍不存在")
	}
	actor, appErr := s.requireManage(actorID, target.ClubID)
	if appErr != nil {
		return appErr
	}
	if actor.Role != "representative" {
		return errors.ErrForbidden("只有负责人可以转让")
	}
	actor.Role = "manager"
	target.Role = "representative"
	if err := s.repo.SaveMembership(actor); err != nil {
		return errors.ErrInternal("转让失败")
	}
	if err := s.repo.SaveMembership(target); err != nil {
		return errors.ErrInternal("转让失败")
	}
	s.notify.Push(target.UserID, "role_changed", "你已成为负责人",
		"同好会负责人已转让给你", "/clubs/"+itoa(target.ClubID), "club_membership", target.ID)
	return nil
}

func (s *ClubService) CreateCode(userID, clubID int64, maxUses, hours int) (*dto.CodeView, *errors.AppError) {
	if _, err := s.requireManage(userID, clubID); err != nil {
		return nil, err
	}
	if maxUses <= 0 {
		maxUses = 1
	}
	raw := make([]byte, 6)
	if _, err := rand.Read(raw); err != nil {
		return nil, errors.ErrInternal("生成绑定码失败")
	}
	code := strings.ToUpper(hex.EncodeToString(raw))
	row := &model.VerificationCode{
		ClubID:    clubID,
		Code:      code,
		MaxUses:   maxUses,
		IsActive:  true,
		CreatedBy: userID,
		CreatedAt: time.Now(),
	}
	if hours > 0 {
		t := time.Now().Add(time.Duration(hours) * time.Hour)
		row.ExpiresAt = &t
	}
	if err := s.repo.CreateCode(row); err != nil {
		return nil, errors.ErrInternal("保存绑定码失败")
	}
	return toCodeView(row), nil
}

func (s *ClubService) ListCodes(userID, clubID int64) ([]dto.CodeView, *errors.AppError) {
	if _, err := s.requireManage(userID, clubID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListCodes(clubID)
	if err != nil {
		return nil, errors.ErrInternal("读取绑定码失败")
	}
	out := make([]dto.CodeView, 0, len(rows))
	for i := range rows {
		out = append(out, *toCodeView(&rows[i]))
	}
	return out, nil
}

func (s *ClubService) RevokeCode(userID, codeID int64) *errors.AppError {
	row, err := s.repo.FindCodeByID(codeID)
	if err != nil {
		return errors.ErrNotFound("绑定码不存在")
	}
	if _, appErr := s.requireManage(userID, row.ClubID); appErr != nil {
		return appErr
	}
	row.IsActive = false
	if err := s.repo.SaveCode(row); err != nil {
		return errors.ErrInternal("吊销失败")
	}
	return nil
}

func (s *ClubService) Redeem(userID int64, code string) *errors.AppError {
	code = strings.TrimSpace(strings.ToUpper(code))
	if code == "" {
		return errors.ErrBadRequest("请输入绑定码")
	}
	row, err := s.repo.FindCode(code)
	if err != nil {
		return errors.ErrNotFound("绑定码无效")
	}
	if !row.IsActive {
		return errors.ErrBadRequest("绑定码已失效")
	}
	if row.ExpiresAt != nil && row.ExpiresAt.Before(time.Now()) {
		return errors.ErrBadRequest("绑定码已过期")
	}
	if row.UseCount >= row.MaxUses {
		return errors.ErrBadRequest("绑定码次数已用尽")
	}
	club, err := s.repo.FindByID(row.ClubID)
	if err != nil {
		return errors.ErrNotFound("同好会不存在")
	}
	existing, _ := s.repo.FindMembership(userID, row.ClubID)
	now := time.Now()
	if existing != nil && existing.Status == "active" {
		return errors.ErrBadRequest("你已经是该同好会成员")
	}
	if existing != nil {
		existing.Status = "active"
		existing.Role = "member"
		existing.LeftAt = nil
		existing.JoinedAt = now
		if err := s.repo.SaveMembership(existing); err != nil {
			return errors.ErrInternal("加入失败")
		}
	} else {
		m := &model.Membership{
			UserID:   userID,
			ClubID:   row.ClubID,
			Country:  club.Country,
			Role:     "member",
			Status:   "active",
			JoinedAt: now,
		}
		if err := s.repo.CreateMembership(m); err != nil {
			return errors.ErrInternal("加入失败")
		}
	}
	row.UseCount++
	if err := s.repo.SaveCode(row); err != nil {
		return errors.ErrInternal("更新绑定码失败")
	}
	return nil
}

func (s *ClubService) Members(userID, clubID int64) ([]dto.MembershipView, *errors.AppError) {
	if _, err := s.requireManage(userID, clubID); err != nil {
		return nil, err
	}
	rows, err := s.repo.ListMembers(clubID, "")
	if err != nil {
		return nil, errors.ErrInternal("读取成员失败")
	}
	return s.toMembershipViews(rows), nil
}

func (s *ClubService) MyMemberships(userID int64) ([]dto.MembershipView, *errors.AppError) {
	rows, err := s.repo.ListUserMemberships(userID)
	if err != nil {
		return nil, errors.ErrInternal("读取会籍失败")
	}
	return s.toMembershipViews(rows), nil
}

func (s *ClubService) PendingForManager(userID int64) ([]dto.MembershipView, *errors.AppError) {
	rows, err := s.repo.ListPendingForManager(userID)
	if err != nil {
		return nil, errors.ErrInternal("读取申请失败")
	}
	return s.toMembershipViews(rows), nil
}

func (s *ClubService) IsFormalActive(userID, membershipID int64) bool {
	m, err := s.repo.FindMembershipByID(membershipID)
	if err != nil || m.UserID != userID || m.Status != "active" {
		return false
	}
	return formalRoles[m.Role]
}

func (s *ClubService) requireManage(userID, clubID int64) (*model.Membership, *errors.AppError) {
	m, err := s.repo.FindMembership(userID, clubID)
	if err != nil || m.Status != "active" || !manageRoles[m.Role] {
		return nil, errors.ErrForbidden("需要同好会管理员或负责人权限")
	}
	return m, nil
}

func (s *ClubService) toSummary(c *model.Club, viewerID int64, detail bool) dto.ClubSummary {
	sum := dto.ClubSummary{
		ID:            c.ID,
		Country:       c.Country,
		Name:          c.Name,
		School:        c.School,
		Province:      c.Province,
		Prefecture:    c.Prefecture,
		City:          c.City,
		Type:          c.Type,
		LogoKey:       c.LogoKey,
		ContactHidden: c.ContactHidden,
		CreatedAt:     c.CreatedAt.UTC().Format(time.RFC3339),
	}
	canSee := !c.ContactHidden
	if viewerID > 0 {
		if m, err := s.repo.FindMembership(viewerID, c.ID); err == nil {
			sum.MyRole = m.Role
			sum.MyStatus = m.Status
			if m.Status == "active" && (formalRoles[m.Role] || manageRoles[m.Role]) {
				canSee = true
			}
			if m.Status == "active" {
				sum.CanApply = false
			} else {
				sum.CanApply = true
			}
		} else if err == gorm.ErrRecordNotFound {
			sum.CanApply = true
		}
	}
	if detail && canSee {
		sum.Info = c.Info
		sum.Remark = c.Remark
		sum.ExternalLinks = c.ExternalLinks
	}
	return sum
}

func (s *ClubService) toMembershipViews(rows []model.Membership) []dto.MembershipView {
	userIDs := make([]int64, 0, len(rows))
	clubIDs := make([]int64, 0, len(rows))
	for _, m := range rows {
		userIDs = append(userIDs, m.UserID)
		clubIDs = append(clubIDs, m.ClubID)
	}
	users := s.repo.UserNames(userIDs)
	clubs := s.repo.ClubNames(clubIDs)
	out := make([]dto.MembershipView, 0, len(rows))
	for _, m := range rows {
		out = append(out, dto.MembershipView{
			ID:             m.ID,
			UserID:         m.UserID,
			UserName:       users[m.UserID],
			ClubID:         m.ClubID,
			ClubName:       clubs[m.ClubID],
			Country:        m.Country,
			Role:           m.Role,
			Status:         m.Status,
			ContactAccount: m.ContactAccount,
			ApplyReason:    m.ApplyReason,
			ApplyRole:      m.ApplyRole,
		})
	}
	return out
}

func toCodeView(c *model.VerificationCode) *dto.CodeView {
	v := &dto.CodeView{
		ID:       c.ID,
		Code:     c.Code,
		MaxUses:  c.MaxUses,
		UseCount: c.UseCount,
		IsActive: c.IsActive,
	}
	if c.ExpiresAt != nil {
		v.ExpiresAt = c.ExpiresAt.Format(time.RFC3339)
	}
	return v
}

func itoa(n int64) string {
	return strconv.FormatInt(n, 10)
}
