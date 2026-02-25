package models

import "time"

type Resource struct {
	ID        int       `json:"id" db:"id"`
	Name      string    `json:"name" db:"name"`
	Capacity  int       `json:"capacity" db:"capacity"` // Вместимость (сколько человек влезет)
	CreatedAt time.Time `json:"created_at" db:"created_at"`
}
