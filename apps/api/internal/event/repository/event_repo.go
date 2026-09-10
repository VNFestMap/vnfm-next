package repository

import (
	"vnfm-api/internal/event/model"

	"gorm.io/gorm"
)

type EventRepository struct {
	db *gorm.DB
}

func NewEventRepository(db *gorm.DB) *EventRepository {
	return &EventRepository{db: db}
}

func (r *EventRepository) List(limit, offset int) ([]model.Event, int64, error) {
	q := r.db.Model(&model.Event{})
	var total int64
	if err := q.Count(&total).Error; err != nil {
		return nil, 0, err
	}
	if limit <= 0 || limit > 100 {
		limit = 50
	}
	var rows []model.Event
	err := q.Order("starts_at DESC").Limit(limit).Offset(offset).Find(&rows).Error
	return rows, total, err
}

func (r *EventRepository) FindByID(id int64) (*model.Event, error) {
	var row model.Event
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *EventRepository) Create(e *model.Event) error {
	return r.db.Create(e).Error
}

func (r *EventRepository) Save(e *model.Event) error {
	return r.db.Save(e).Error
}

func (r *EventRepository) Delete(id int64) error {
	return r.db.Delete(&model.Event{}, id).Error
}

func (r *EventRepository) FindReg(eventID, userID int64) (*model.Registration, error) {
	var row model.Registration
	if err := r.db.Where("event_id = ? AND user_id = ?", eventID, userID).First(&row).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *EventRepository) CreateReg(reg *model.Registration) error {
	return r.db.Create(reg).Error
}

func (r *EventRepository) DeleteReg(eventID, userID int64) error {
	return r.db.Where("event_id = ? AND user_id = ?", eventID, userID).Delete(&model.Registration{}).Error
}

func (r *EventRepository) ListUserRegs(userID int64) ([]model.Registration, error) {
	var rows []model.Registration
	err := r.db.Where("user_id = ?", userID).Order("created_at DESC").Find(&rows).Error
	return rows, err
}

func (r *EventRepository) ClubName(id int64) string {
	var name string
	r.db.Table("clubs").Select("name").Where("id = ?", id).Scan(&name)
	return name
}
