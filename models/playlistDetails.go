package models

import (
	"gorm.io/gorm"
)

type PlaylistDetails struct {
	gorm.Model
	UserID     uint      `json:"user_id"`
	User       *User     `gorm:"foreignKey:UserID" json:"user"`
	PlaylistID uint      `json:"playlist_id"`
	Playlist   *Playlist `gorm:"foreignKey:PlaylistID" json:"playlist"`
	Artist     string    `json:"artist"`
	Title      string    `json:"title"`
	Message    string    `jaon:"message"`
}
