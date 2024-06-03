package service

import (
	"crypto/sha256"
	"fmt"

	"github.com/omarelweshy/EcomMaster-user-service/internal/model"
	"gorm.io/gorm"
)

type UserService struct {
	DB *gorm.DB
}

func (s *UserService) RegisterUser(username, email, password string) error {
	hashedPassword := fmt.Sprintf("%x", sha256.Sum256([]byte(password)))
	user := model.User{
		Username:     username,
		Email:        email,
		PasswordHash: hashedPassword,
	}
	return s.DB.Create(&user).Error
}
