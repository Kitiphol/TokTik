package service

import (
    "errors"
    "time"

    "github.com/golang-jwt/jwt/v4"
    "golang.org/x/crypto/bcrypt"
	"github.com/google/uuid"
	
    "Auth/internal/config"
    "Auth/internal/database"
    "Auth/internal/entity"
	"Auth/internal/DTO"
)

func Login(req DTO.LoginRequest) (string, string, error) {
    var user entity.User
    if err := database.DB.Where("username = ?", req.Username).First(&user).Error; err != nil {
        return "", "", errors.New("invalid credentials")
    }

    if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
        return "", "", errors.New("invalid credentials")
    }

    token, err := generateJWT(user.ID)
    return token, user.Username, err
}


// generateJWT signs a new token for a given user ID
func generateJWT(userID uuid.UUID) (string, error) {
    cfg := config.Load()

    claims := jwt.RegisteredClaims{
        Subject:   userID.String(),                           // "sub"
        
        ExpiresAt: jwt.NewNumericDate(time.Now().Add(72 * time.Hour)), // "exp"
        IssuedAt:  jwt.NewNumericDate(time.Now()),            // "iat" (optional but good)
    }

    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString([]byte(cfg.JWTSecret))
}

