package models

import (
	"time"
)

type Booking struct {
	ID         int        `json:"id"`
	UserID     int        `json:"user_id"`
	ResourceID int        `json:"resource_id"`
	Title      string     `json:"title"`
	StartAt    time.Time  `json:"start_at"`
	EndAt      time.Time  `json:"end_at"`
	IsHoliday  bool       `json:"is_holiday"`
	CreatedAt  time.Time  `json:"created_at"`
	UpdatedAt  time.Time  `json:"updated_at"`
	DeletedAt  *time.Time `json:"deleted_at"`
}
