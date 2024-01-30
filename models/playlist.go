package models

import (
	"gorm.io/gorm"
)

type Playlist struct {
	gorm.Model
	PlaylistDetails []PlaylistDetails `json:"playlistDetails"`
}
