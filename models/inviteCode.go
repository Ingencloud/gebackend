package models

import (
	"time"
)

type InviteCode struct {
	ID        uint   `gorm:"primaryKey"`
	Code      string `gorm:"unique"`
	IsUsed    bool
	CreatedAt time.Time
	UpdatedAt time.Time
}
