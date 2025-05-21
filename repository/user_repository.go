package repository

import (
	"github.com/thanavatC/auth-service-go/model"
	"gorm.io/gorm"
)

type IUserRepository interface {
	FindByEmail(email string) (model.User, error)
	Create(user model.User) error
	FindByID(id string) (model.User, error)
}

type UserRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) IUserRepository {
	return &UserRepository{
		db: db,
	}
}

func (repo *UserRepository) FindByEmail(email string) (model.User, error) {
	var user model.User
	if err := repo.db.Where("email = ?", email).First(&user).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}

func (repo *UserRepository) Create(user model.User) error {
	return repo.db.Create(&user).Error
}

func (repo *UserRepository) FindByID(id string) (model.User, error) {
	var user model.User
	if err := repo.db.First(&user, "id = ?", id).Error; err != nil {
		return model.User{}, err
	}
	return user, nil
}
