package users

import (
	"github.com/giovannyptr/bookshelf/models"
	"gorm.io/gorm"
)

type Repository struct{ db *gorm.DB }

func NewRepository(db *gorm.DB) *Repository { return &Repository{db: db} }
func (r *Repository) Migrate() error        { return r.db.AutoMigrate(&models.User{}) }

func (r *Repository) ByEmail(email string) (*models.User, error) {
	var u models.User
	if err := r.db.Where("email = ?", email).First(&u).Error; err != nil {
		return nil, err
	}
	return &u, nil
}

func (r *Repository) Create(u *models.User) error { return r.db.Create(u).Error }
func (r *Repository) ByID(id uint) (*models.User, error) {
	var u models.User
	if err := r.db.First(&u, id).Error; err != nil {
		return nil, err
	}
	return &u, nil
}
