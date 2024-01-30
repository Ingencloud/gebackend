package models

import (
	"gorm.io/gorm"
)

type User struct {
	gorm.Model
	FullName        string            `json:"name"`
	Number          string            `json:"number" gorm:"unique"`
	PlaylistDetails []PlaylistDetails `json:"playlistDetails" gorm:"foreignKey:UserID"`
}
