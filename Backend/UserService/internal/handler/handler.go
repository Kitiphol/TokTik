package handler

import (
	"github.com/gin-gonic/gin"
	"UserService/internal/service"
	"UserService/internal/DTO"
	"UserService/internal/middleware"
	"log"

)

type RegisterHandler struct{}

func (h *RegisterHandler) Handle(c *gin.Context) {
    var req DTO.UpdateUserRequest
    if err := c.ShouldBindJSON(&req); err != nil {
        log.Println("[UpdateUserHandler] JSON bind error:", err)
        c.JSON(400, gin.H{"error": err.Error()})
        return
    }

    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        log.Println("[UpdateUserHandler] Error getting userID from context:", err)
        c.JSON(401, gin.H{"error": "unauthorized"})
        return
    }

    log.Println("[UpdateUserHandler] Received update request for userID:", userID.String())

    msg, updatedUsername, err := service.UpdateUser(userID, req)
    if err != nil {
        log.Println("[UpdateUserHandler] Error updating user:", err)
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }

    log.Println("[UpdateUserHandler] Update successful for userID:", userID.String())
    c.JSON(200, gin.H{
        "message":  msg,
        "username": updatedUsername,
    })
}
