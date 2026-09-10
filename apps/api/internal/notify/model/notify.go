package model

import "time"

type Notification struct {
	ID          int64      `gorm:"primaryKey"`
	UserID      int64      `gorm:"not null"`
	Type        string     `gorm:"not null"`
	Title       string     `gorm:"not null"`
	Message     string     `gorm:"not null;default:''"`
	Link        string     `gorm:"not null;default:''"`
	RelatedType string     `gorm:"not null;default:''"`
	RelatedID   int64      `gorm:"not null;default:0"`
	IsRead      bool       `gorm:"not null;default:false"`
	CreatedAt   time.Time  `gorm:"not null"`
	ReadAt      *time.Time `gorm:""`
}

func (Notification) TableName() string { return "notifications" }
