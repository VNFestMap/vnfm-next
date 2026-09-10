package repository

import (
	"time"

	"vnfm-api/internal/notify/model"

	"gorm.io/gorm"
)

type NotifyRepository struct {
	db *gorm.DB
}

func NewNotifyRepository(db *gorm.DB) *NotifyRepository {
	return &NotifyRepository{db: db}
}

func (r *NotifyRepository) Create(n *model.Notification) error {
	n.CreatedAt = time.Now()
	return r.db.Create(n).Error
}

func (r *NotifyRepository) List(userID int64, limit, offset int) ([]model.Notification, int64, error) {
	q := r.db.Model(&model.Notification{}).Where("user_id = ?", userID)
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var rows []model.Notification
	err := q.Order("created_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

func (r *NotifyRepository) UnreadCount(userID int64) (int64, error) {
	var n int64
	err := r.db.Model(&model.Notification{}).Where("user_id = ? AND is_read = FALSE", userID).Count(&n).Error
	return n, err
}

func (r *NotifyRepository) MarkRead(userID, id int64) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Where("id = ? AND user_id = ?", id, userID).
		Updates(map[string]any{"is_read": true, "read_at": now}).Error
}

func (r *NotifyRepository) MarkAllRead(userID int64) error {
	now := time.Now()
	return r.db.Model(&model.Notification{}).
		Where("user_id = ? AND is_read = FALSE", userID).
		Updates(map[string]any{"is_read": true, "read_at": now}).Error
}
