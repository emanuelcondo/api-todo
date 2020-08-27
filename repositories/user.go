package repositories

import (
	"github.com/emanuelcondo/api-todo/db"
	"github.com/emanuelcondo/api-todo/models"
	"github.com/jinzhu/gorm"
)

type UserRepository struct{}

// Create a new user
func (userRepository *UserRepository) Create(user *models.User) (*models.User, error) {
	DB := db.GetDBConnection()
	createdUser := DB.Create(&user)

	if createdUser.Error != nil {
		return nil, createdUser.Error
	} else {
		return user, nil
	}
}

// Find User by email
func (userRepository *UserRepository) FindByEmail(email string) (*models.User, error) {
	DB := db.GetDBConnection()
	user := &models.User{}
	foundUser := DB.Where("email = ?", email).First(user)

	if foundUser.Error != nil {
		if foundUser.Error == gorm.ErrRecordNotFound {
			return nil, nil
		}
		return nil, foundUser.Error
	} else {
		return user, nil
	}
}