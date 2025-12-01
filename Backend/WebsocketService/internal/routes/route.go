package routes

import (
	"WebsocketService/internal/service"
	"WebsocketService/internal/middleware"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"time"
)

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

	protected := r.Group("/api/ws", middleware.AuthMiddleware())
	{
		protected.POST("/comments/videos/:videoID", service.AddOrUpdateComment)
		protected.POST("/likes/videos/:videoID", service.AddOrUpdateLike)
		protected.POST("/view/videos/:videoID", service.AddOrUpdateView)

		protected.GET("/likes/videos/:videoID", service.GetTotalLikes)
		protected.GET("/views/videos/:videoID", service.GetTotalViews)
		protected.GET("/comments/videos/:videoID", service.ListComments)

		protected.DELETE("/likes/videos/:videoID", service.DeleteLike)
		protected.DELETE("/comments/videos/:videoID/:commentID", service.DeleteComment)

		protected.GET("/notifications/unread", service.GetNotifications)
		protected.POST("/notifications/:id/read", service.GetReadNotifications)


	}

	return r
}
