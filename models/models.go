package models

import (
	"time"
)

// Song represents the structure of a song in the database
type Song struct {
	ID          uint       `json:"id" gorm:"primaryKey" example:"1"`
	CreatedAt   time.Time  `json:"created_at" example:"2024-10-02T15:04:05Z07:00"`
	UpdatedAt   time.Time  `json:"updated_at" example:"2024-10-02T15:04:05Z07:00"`
	DeletedAt   *time.Time `json:"deleted_at,omitempty" example:"2024-10-02T15:04:05Z07:00"`
	Group       string     `json:"group" example:"Muse"`
	Title       string     `json:"title" example:"Supermassive Black Hole"`
	ReleaseDate string     `json:"release_date" example:"16.07.2006"`
	Text        string     `json:"text" example:"Ooh baby, don't you know I suffer?"`
	Link        string     `json:"link" example:"https://www.youtube.com/watch?v=Xsp3_a-PMTw"`
}
