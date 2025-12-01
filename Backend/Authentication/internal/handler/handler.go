// internal/handler/handler.go
package handler

import (
    "net/http"

    "github.com/gin-gonic/gin"
    "Auth/internal/service"
    "Auth/internal/DTO"
    "Auth/internal/middleware"
)

// RegisterHandler handles registration and login
type RegisterHandler struct{}

// Handle registration
func (h *RegisterHandler) Handle(c *gin.Context) {
    var req DTO.RegisterRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    token, username, err := service.Register(req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "token":    token,
        "username": username,
    })
}

// Login handles login
func (h *RegisterHandler) Login(c *gin.Context) {
    var req DTO.LoginRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    token, username, err := service.Login(req)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{
        "token":    token,
        "username": username,
    })
}



func Profile(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"userID": userID.String()})
}