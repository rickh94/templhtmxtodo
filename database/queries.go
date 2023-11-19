package database

import (
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
)

func GetOrCreateUser(email string) (*User, error) {
	var user User
	result := DB.First(&user, "email = ?", email)
	if result.Error == gorm.ErrRecordNotFound {
		log.Default().Printf("Creating user: %v\n", email)
		// create user
		user.Email = email
		user.CreatedAt = time.Now()
		result := DB.Create(&user)
		if result.Error != nil {
			return nil, result.Error
		} else {
			return &user, nil
		}
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func GetUserByEmail(email string) (*User, error) {
	var user User
	result := DB.Preload("Credentials").First(&user, "email = ?", email)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user %s not found", email)
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}

func GetUserByID(id string) (*User, error) {
	var user User
	result := DB.First(&user, "id = ?", id)
	if result.Error == gorm.ErrRecordNotFound {
		return nil, fmt.Errorf("user %s not found", id)
	} else if result.Error != nil {
		return nil, result.Error
	} else {
		return &user, nil
	}
}
