package auth

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserRepository interface {
	RegisterUser(user *User) error
	FindByEmail(email string) (*User, error)
	FindByID(id uuid.UUID) (*User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) RegisterUser(user *User) error {
	result := r.db.Create(user)
	return result.Error
}

func (r *userRepository) FindByEmail(email string) (*User, error) {
	var user User

	result := r.db.Where("email = ?", email).First(&user)
	return &user, result.Error
}

func (r *userRepository) FindByID(id uuid.UUID) (*User, error) {
	var user User

	result := r.db.Where("id = ?", id).First(&user)

	return &user, result.Error
}
