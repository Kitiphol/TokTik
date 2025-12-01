package service

import (

	"errors"
    "golang.org/x/crypto/bcrypt"

    "Auth/internal/database"
    "Auth/internal/entity"
	"Auth/internal/DTO"
)


func Register(req DTO.RegisterRequest) (string, string, error) {
    if req.Password == "" {
        return "", "", errors.New("password cannot be empty")
    }

    if req.Password != req.ConfirmPassword {
        return "", "", errors.New("passwords do not match")
    }

    hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
    if err != nil {
        return "", "", errors.New("failed to hash password")
    }

    if database.DB.Where("username = ?", req.Username).First(&entity.User{}).RowsAffected > 0 {
        return "", "", errors.New("username already exists")
    }

    if database.DB.Where("email = ?", req.Email).First(&entity.User{}).RowsAffected > 0 {
        return "", "", errors.New("email already exists")
    }

    user := entity.User{Username: req.Username, PasswordHash: string(hash), Email: req.Email}
    if err := database.DB.Create(&user).Error; err != nil {
        return "", "", errors.New("failed to create user")
    }

    token, err := generateJWT(user.ID)
    return token, user.Username, err
}