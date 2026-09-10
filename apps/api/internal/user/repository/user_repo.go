package repository

import (
	"time"

	"vnfm-api/internal/user/model"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Ensure(id int64, name, avatar string) error {
	now := time.Now()
	row := model.User{
		ID:        id,
		Name:      name,
		AvatarURL: avatar,
		CreatedAt: now,
		UpdatedAt: now,
	}
	return r.db.Clauses(clause.OnConflict{
		Columns:   []clause.Column{{Name: "id"}},
		DoUpdates: clause.AssignmentColumns([]string{"name", "avatar_url", "updated_at"}),
	}).Create(&row).Error
}

func (r *UserRepository) FindByID(id int64) (*model.User, error) {
	var row model.User
	if err := r.db.First(&row, id).Error; err != nil {
		return nil, err
	}
	return &row, nil
}

func (r *UserRepository) UpdatePrefs(id int64, language, theme *string, displayID *int64, clearDisplay bool) error {
	updates := map[string]any{"updated_at": time.Now()}
	if language != nil {
		updates["language_preference"] = *language
	}
	if theme != nil {
		updates["theme_preference"] = *theme
	}
	if clearDisplay {
		updates["display_membership_id"] = nil
	} else if displayID != nil {
		updates["display_membership_id"] = *displayID
	}
	return r.db.Model(&model.User{}).Where("id = ?", id).Updates(updates).Error
}
