package models

import (
	"time"
)

type URL struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	Address          string    `gorm:"not null" json:"address"`
	Status           string    `json:"status"` // queued, running, done, error
	Title            string    `json:"title"`
	HTMLVersion      string    `json:"html_version"`
	HasLoginForm     bool      `json:"has_login_form"`
	InternalLinks    int       `json:"internal_links"`
	ExternalLinks    int       `json:"external_links"`
	BrokenLinks      string    `gorm:"type:text" json:"broken_links"` // JSON string or comma-separated
	Error            string    `json:"error"`                         // Error message if failed
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
