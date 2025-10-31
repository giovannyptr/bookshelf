package models

import "time"

// User represents an account.
// swagger:model User
type User struct {
	ID        uint      `json:"id"        gorm:"primaryKey"`
	Email     string    `json:"email"     gorm:"uniqueIndex"`
	Password  string    `json:"-"` // never expose
	Name      string    `json:"name"`
	Role      string    `json:"role" example:"admin"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
