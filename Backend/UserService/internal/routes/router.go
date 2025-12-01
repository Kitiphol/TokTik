package routes

import (
	"UserService/internal/handler"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-contrib/cors"
	"UserService/internal/middleware"
)

type RegisterHandler = handler.RegisterHandler

func Setup() *gin.Engine {
	r := gin.Default()

	// r.Use(cors.New(cors.Config{
	// 	AllowOrigins:     []string{"http://localhost:3000"},
	// 	AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
	// 	AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
	// 	ExposeHeaders:    []string{"Content-Length"},
	// 	AllowCredentials: true,
	// 	MaxAge:           12 * time.Hour,
	// }))


    r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization"},
		ExposeHeaders:    []string{"Content-Length", "Authorization"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	auth := r.Group("/api/user")
	auth.Use(middleware.AuthMiddleware())

	reg := RegisterHandler{}
	auth.POST("/updateUser", reg.Handle)

	return r
}
