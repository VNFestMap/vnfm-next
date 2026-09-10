package model

import "time"

type Event struct {
	ID            int64      `gorm:"primaryKey"`
	ClubID        int64      `gorm:"not null"`
	Title         string     `gorm:"not null"`
	Description   string     `gorm:"not null;default:''"`
	Location      string     `gorm:"not null;default:''"`
	CoverKey      string     `gorm:"not null;default:''"`
	StartsAt      time.Time  `gorm:"not null"`
	EndsAt        *time.Time `gorm:""`
	RegisterUntil *time.Time `gorm:""`
	CreatedBy     int64      `gorm:"not null"`
	CreatedAt     time.Time  `gorm:"not null"`
	UpdatedAt     time.Time  `gorm:"not null"`
}

func (Event) TableName() string { return "events" }

type Registration struct {
	ID        int64     `gorm:"primaryKey"`
	EventID   int64     `gorm:"not null"`
	UserID    int64     `gorm:"not null"`
	CreatedAt time.Time `gorm:"not null"`
}

func (Registration) TableName() string { return "event_registrations" }
