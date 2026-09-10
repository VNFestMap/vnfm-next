package service

import (
	"vnfm-api/internal/notify/model"
	"vnfm-api/internal/notify/repository"
	"vnfm-api/pkg/errors"
)

type NotifyService struct {
	repo *repository.NotifyRepository
}

func NewNotifyService(repo *repository.NotifyRepository) *NotifyService {
	return &NotifyService{repo: repo}
}

func (s *NotifyService) Push(userID int64, typ, title, message, link, relatedType string, relatedID int64) {
	_ = s.repo.Create(&model.Notification{
		UserID:      userID,
		Type:        typ,
		Title:       title,
		Message:     message,
		Link:        link,
		RelatedType: relatedType,
		RelatedID:   relatedID,
	})
}

func (s *NotifyService) List(userID int64, limit, offset int) ([]model.Notification, int64, *errors.AppError) {
	rows, total, err := s.repo.List(userID, limit, offset)
	if err != nil {
		return nil, 0, errors.ErrInternal("读取通知失败")
	}
	return rows, total, nil
}

func (s *NotifyService) UnreadCount(userID int64) (int64, *errors.AppError) {
	n, err := s.repo.UnreadCount(userID)
	if err != nil {
		return 0, errors.ErrInternal("读取未读数失败")
	}
	return n, nil
}

func (s *NotifyService) MarkRead(userID, id int64) *errors.AppError {
	if err := s.repo.MarkRead(userID, id); err != nil {
		return errors.ErrInternal("标记已读失败")
	}
	return nil
}

func (s *NotifyService) MarkAllRead(userID int64) *errors.AppError {
	if err := s.repo.MarkAllRead(userID); err != nil {
		return errors.ErrInternal("标记已读失败")
	}
	return nil
}
