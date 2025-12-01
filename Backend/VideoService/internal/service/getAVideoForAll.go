package service

import (
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"VideoService/internal/s3util"
	"errors"
	"log"
	"time"

	"github.com/google/uuid"
)

func GetVideoForAllUser(videoID uuid.UUID) (string, error) {
    log.Printf("[DEBUG] GetVideoForAllUser called with videoID=%s", videoID.String())
    var video entity.Video
    if err := database.DB.Where("id = ?", videoID).Preload("Comments").First(&video).Error; err != nil {
        log.Printf("[ERROR] Video not found: %v", err)
        return "", errors.New("video not found")
    }
    if err := database.DB.Save(&video).Error; err != nil {
        log.Printf("[ERROR] Failed to update view count: %v", err)
        return "", errors.New("failed to update view count")
    }
    database.DB.Model(&video).UpdateColumn("total_view_count", video.TotalViewCount+1)
    log.Printf("[DEBUG] Returning presigned URL for videoID=%s, key=%s", videoID.String(), video.Filename)
    return s3util.GeneratePresignedGetURL("toktikp2-video", video.Filename, 1*time.Hour)
}