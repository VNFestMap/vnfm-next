package model

import "time"

type User struct {
	ID                   int64      `gorm:"primaryKey"`
	Name                 string     `gorm:"not null;default:''"`
	AvatarURL            string     `gorm:"column:avatar_url;not null;default:''"`
	LanguagePreference   string     `gorm:"not null;default:zh"`
	ThemePreference      string     `gorm:"not null;default:system"`
	DisplayMembershipID  *int64     `gorm:"column:display_membership_id"`
	CreatedAt            time.Time  `gorm:"not null"`
	UpdatedAt            time.Time  `gorm:"not null"`
}

func (User) TableName() string { return "users" }
