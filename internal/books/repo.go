package books

import (
	"fmt"

	"github.com/giovannyptr/bookshelf/models"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }

func (r *Repository) Migrate() error { return r.db.AutoMigrate(&models.Book{}) }

func (r *Repository) List(q, category string, page, limit int, sort, order string) (items []models.Book, total int64, err error) {
	tx := r.db.Model(&models.Book{})
	if q != "" {
		like := "%" + q + "%"
		tx = tx.Where("title ILIKE ? OR author ILIKE ?", like, like)
	}
	if category != "" {
		tx = tx.Where("category = ?", category)
	}
	_ = tx.Count(&total).Error

	allowed := map[string]bool{"title": true, "category": true, "price": true, "created_at": true}
	if !allowed[sort] {
		sort = "created_at"
	}
	if order != "ASC" {
		order = "DESC"
	}

	err = tx.Order(fmt.Sprintf("%s %s", sort, order)).
		Offset((page - 1) * limit).Limit(limit).
		Find(&items).Error
	return
}

func (r *Repository) ByID(id string) (models.Book, error) {
	var b models.Book
	err := r.db.First(&b, id).Error
	return b, err
}

func (r *Repository) Create(b *models.Book) error { return r.db.Create(b).Error }
func (r *Repository) Save(b *models.Book) error   { return r.db.Save(b).Error }
func (r *Repository) Delete(b *models.Book) error { return r.db.Delete(b).Error }
