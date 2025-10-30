package models

import "time"

type Book struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Title     string    `json:"title" gorm:"index;not null"`
	Author    string    `json:"author"`
	Category  string    `json:"category" gorm:"index"`
	Price     float64   `json:"price"`
	Stock     int       `json:"stock"`
	CoverURL  string    `json:"coverUrl"` 
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
