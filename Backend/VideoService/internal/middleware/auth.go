package middleware

import (
    "errors"
    "strings"
    "log"
    "github.com/gin-gonic/gin"
    "github.com/golang-jwt/jwt/v4"
    "github.com/google/uuid"
    "VideoService/internal/config"
)

// Context key for userID
const ContextUserID = "userID"

func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        authHeader := c.GetHeader("Authorization")
        // Log the raw Authorization header
        log.Println("[AuthMiddleware] Authorization header:", authHeader)

        if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
            log.Println("[AuthMiddleware] Missing or malformed Authorization header")
            c.AbortWithStatusJSON(401, gin.H{"error": "Missing or malformed Authorization header"})
            return
        }

        tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
        cfg := config.Load()

        token, err := jwt.ParseWithClaims(tokenStr, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
            if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
                log.Println("[AuthMiddleware] Unexpected signing method:", token.Header["alg"])
                return nil, errors.New("unexpected signing method")
            }
            return []byte(cfg.JWTSecret), nil
        })

        if err != nil {
            log.Println("[AuthMiddleware] Error parsing token:", err)
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
            return
        }

        if !token.Valid {
            log.Println("[AuthMiddleware] Token is invalid")
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid or expired token"})
            return
        }

        claims, ok := token.Claims.(*jwt.RegisteredClaims)
        if !ok || claims.Subject == "" {
            log.Println("[AuthMiddleware] Invalid token claims")
            c.AbortWithStatusJSON(401, gin.H{"error": "Invalid token claims"})
            return
        }

        userID, err := uuid.Parse(claims.Subject)
        if err != nil {
            log.Println("[AuthMiddleware] Malformed userID in token:", claims.Subject)
            c.AbortWithStatusJSON(401, gin.H{"error": "Malformed userID in token"})
            return
        }

        log.Println("[AuthMiddleware] Successfully extracted userID:", userID.String())

        c.Set(ContextUserID, userID)
        c.Next()
    }
}

func GetUserIDFromContext(c *gin.Context) (uuid.UUID, error) {
    val, exists := c.Get(ContextUserID)
    if !exists {
        return uuid.Nil, errors.New("userID not found in context")
    }
    userID, ok := val.(uuid.UUID)
    if !ok {
        return uuid.Nil, errors.New("invalid userID type in context")
    }
    return userID, nil
}
