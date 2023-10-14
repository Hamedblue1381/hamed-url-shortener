package model

import (
	"errors"
	"fmt"

	"gorm.io/gorm"
)

func CreateUser(user *User) error {
	result := db.Create(user)
	if result.Error != nil {
		return result.Error
	}
	return nil
}

func GetUserByID(userID uint64) (*User, error) {
	var user User
	result := db.First(&user, userID)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found")
		}
		return nil, result.Error
	}
	return &user, nil
}

// GetUserByUsername fetches a user from the database by their username.
func GetUser(username string) (*User, error) {
	var user User
	result := db.Where("username = ?", username).First(&user)
	if result.Error != nil {
		if errors.Is(result.Error, gorm.ErrRecordNotFound) {
			return nil, fmt.Errorf("User not found")
		}
		return nil, result.Error
	}
	return &user, nil
}
