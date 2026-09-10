package model

import "time"

type Club struct {
	ID            int64     `gorm:"primaryKey"`
	Country       string    `gorm:"not null;default:china"`
	Name          string    `gorm:"not null"`
	School        string    `gorm:"not null;default:''"`
	Province      string    `gorm:"not null;default:''"`
	Prefecture    string    `gorm:"not null;default:''"`
	City          string    `gorm:"not null;default:''"`
	Type          string    `gorm:"not null;default:school"`
	Info          string    `gorm:"not null;default:''"`
	Remark        string    `gorm:"not null;default:''"`
	LogoKey       string    `gorm:"not null;default:''"`
	ExternalLinks string    `gorm:"not null;default:''"`
	ContactHidden bool      `gorm:"not null;default:true"`
	CreatedBy     int64     `gorm:"not null"`
	CreatedAt     time.Time `gorm:"not null"`
	UpdatedAt     time.Time `gorm:"not null"`
}

func (Club) TableName() string { return "clubs" }

type Membership struct {
	ID             int64      `gorm:"primaryKey"`
	UserID         int64      `gorm:"not null"`
	ClubID         int64      `gorm:"not null"`
	Country        string     `gorm:"not null;default:china"`
	Role           string     `gorm:"not null;default:member"`
	Status         string     `gorm:"not null;default:pending"`
	ContactAccount string     `gorm:"not null;default:''"`
	ApplyReason    string     `gorm:"not null;default:''"`
	ApplyRole      string     `gorm:"not null;default:member"`
	ReviewedBy     *int64     `gorm:""`
	ReviewedAt     *time.Time `gorm:""`
	JoinedAt       time.Time  `gorm:"not null"`
	LeftAt         *time.Time `gorm:""`
}

func (Membership) TableName() string { return "club_memberships" }

type VerificationCode struct {
	ID        int64      `gorm:"primaryKey"`
	ClubID    int64      `gorm:"not null"`
	Code      string     `gorm:"not null"`
	MaxUses   int        `gorm:"not null;default:1"`
	UseCount  int        `gorm:"not null;default:0"`
	ExpiresAt *time.Time `gorm:""`
	IsActive  bool       `gorm:"not null;default:true"`
	CreatedBy int64      `gorm:"not null"`
	CreatedAt time.Time  `gorm:"not null"`
}

func (VerificationCode) TableName() string { return "club_verification_codes" }
