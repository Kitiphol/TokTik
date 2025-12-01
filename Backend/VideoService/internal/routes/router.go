package routes

import (
	"VideoService/internal/handler"
	"VideoService/internal/middleware"
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



	videoHandler := handler.VideoHandler{}

	// Public: list all videos
	r.GET("/api/video/videos", videoHandler.GetAllVideos)

	// Protected routes
	protected := r.Group("/api/video", middleware.AuthMiddleware())
	{
		protected.GET("/videos/user", videoHandler.GetAllVideosForUser)
		protected.GET("/videos/:videoID", videoHandler.GetVideoForUser)
		protected.POST("/videos", videoHandler.CreateVideoForUser)
		protected.GET("/videos/check", videoHandler.CheckUpload)
		protected.PUT("/videos/:videoID", videoHandler.UpdateVideoForUser)
		protected.DELETE("/videos/:videoID", videoHandler.DeleteVideoForUser)
	}

	return r
}
