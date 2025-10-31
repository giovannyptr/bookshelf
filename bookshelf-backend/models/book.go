package models

import "time"

// Book represents a book entity.
// swagger:model Book
type Book struct {
	ID        uint      `json:"id"        gorm:"primaryKey"`
	Title     string    `json:"title"     gorm:"index;not null" example:"1984"`
	Author    string    `json:"author"    example:"George Orwell"`
	Category  string    `json:"category"  example:"Fiction"`
	Price     float64   `json:"price"     example:"60000"`
	Stock     int       `json:"stock"     example:"9"`
	CoverURL  string    `json:"coverUrl"  example:"/uploads/uuid.jpg"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}
