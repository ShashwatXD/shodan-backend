package models

import "time"

// Result is the GORM model for storing analysis results.
type Result struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	Text       string    `gorm:"type:TEXT" json:"text"`
	Output     string    `gorm:"type:TEXT" json:"output"`
	Confidence float64   `json:"confidence"`
	Model      string    `json:"model"`
	CreatedAt  time.Time `json:"created_at"`
}
