package models

import (
	"time"
)

type Song struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time  `json:"created_at"`
	UpdatedAt   time.Time  `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty"`
	Group       string     `json:"group"`
	Title       string     `json:"title"`
	ReleaseDate string     `json:"release_date"`
	Text        string     `json:"text"`
	Link        string     `json:"link"`
}
