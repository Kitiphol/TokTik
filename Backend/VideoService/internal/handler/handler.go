package handler

import (
	"VideoService/internal/DTO"
	"VideoService/internal/machineryUtil"
	"VideoService/internal/middleware"
	"VideoService/internal/service"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type VideoHandler struct{}

// GetAllVideos returns all videos in the system (public)
func (h *VideoHandler) GetAllVideos(c *gin.Context) {
    videos, err := service.GetAllVideosWithThumbnails()
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

// GetAllVideosForUser returns videos for the authenticated user
func (h *VideoHandler) GetAllVideosForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videos, err := service.GetAllVideosForUser(userID)
    if err != nil {
         log.Printf("[ERROR] GetAllVideosForUser: %v", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, videos)
}

// GetVideoForUser returns a presigned GET URL for a user's video
func (h *VideoHandler) GetVideoForUser(c *gin.Context) {
    videoID := c.Param("videoID")
    resp, err := service.GetVideoForUser(videoID)
    if err != nil {
        c.JSON(404, gin.H{"error": err.Error()})
        return
    }
    c.JSON(200, resp)
}   

// CreateVideoForUser issues a presigned PUT URL for a new user video
func (h *VideoHandler) CreateVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    var req DTO.CreateVideoRequest
    if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
        return
    }

    presignURL, videoID, err := service.CreateVideoForUser(userID, req)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusCreated, gin.H{"message": "video record created", "presignedURL": presignURL, "videoID": videoID})
}

// UpdateVideoForUser updates metadata for an existing video
func (h *VideoHandler) UpdateVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videoIDStr := c.Param("videoID")
    videoID, parseErr := uuid.Parse(videoIDStr)
    if parseErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    var req DTO.UpdateVideoRequest
    if bindErr := c.ShouldBindJSON(&req); bindErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": bindErr.Error()})
        return
    }

    if err := service.UpdateVideoForUser(userID, videoID, req); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "video updated"})
}

// DeleteVideoForUser deletes a user's video
func (h *VideoHandler) DeleteVideoForUser(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videoIDStr := c.Param("videoID")
    videoID, parseErr := uuid.Parse(videoIDStr)
    if parseErr != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid video ID"})
        return
    }

    if err := service.DeleteVideoForUser(userID, videoID); err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "video deleted"})
}



func (h *VideoHandler) CheckUpload(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "unauthorized"})
        return
    }

    videoIDStr := c.Query("videoID")
    if videoIDStr == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "videoID is required"})
        return
    }
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid videoID"})
        return
    }

    log.Printf("Checking upload status for user: %s, videoID: %s", userID.String(), videoID.String());

    video, err := service.GetVideoRecord(videoID)
    if err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "video not found"})
        return
    }


    if video.UserID != userID {
        c.JSON(http.StatusForbidden, gin.H{"error": "you do not own this video"})
        return
    }

    uploaded, err := service.CheckVideoUploaded(userID, videoID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    //send the message to video Chunker

    if uploaded {
        // Import your machineryutil package at the top if not already:
        log.Printf("The uploaded video name (Before Chunking and Converting) is : %s", video.Filename)



        err := machineryutil.SendCreateThumbnailTask("toktikp2-video", video.Filename) 

        if err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue thumbnail creation task", "details": err.Error()})
            return
        }


        err2 := machineryutil.SendChunkVideoTask("toktikp2-video", video.Filename)
        if err2 != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to queue chunking task", "details": err.Error()})
            return
        }
    }

    log.Printf("Sending chunking task for video:", videoID.String(), "uploaded:", uploaded)
    log.Printf("The Video Name is: %s   (Should be the path to file in R2)", video.Filename)

    c.JSON(http.StatusOK, gin.H{
        "videoID":  videoID.String(),
        "uploadStatus": uploaded,
    })
}