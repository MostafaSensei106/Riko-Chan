package db

import (
	"gorm.io/gorm"

	"github.com/MostafaSensei106/Riko-Chan/internal/models"
)

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *models.User) error {
	return r.db.Create(user).Error
}

func (r *UserRepository) GetByID(id int64) (*models.User, error) {
	var user models.User
	result := r.db.First(&user, "id = ?", id)
	if result.Error != nil {
		return nil, result.Error
	}
	return &user, nil
}

func (r *UserRepository) Update(user *models.User) error {
	return r.db.Save(user).Error
}

// func (r *UserRepository) Upsert(user *models.User) error {
// 	return r.db.Clauses(
// 		gorm.Clauses{gorm.OnConflict{
// 			UpdateAll: true,
// 		}},
// 	).Create(user).Error
// }
