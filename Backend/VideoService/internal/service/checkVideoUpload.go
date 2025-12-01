package service

import (
    "VideoService/internal/database"
    "github.com/google/uuid"
    "VideoService/internal/s3util"
	"VideoService/internal/entity"
)

func CheckVideoUploaded(userID, videoID uuid.UUID) (bool, error) {
    // Find the video record to get the filename/key
    var video entity.Video
    if err := database.DB.Where("user_id = ? AND id = ?", userID, videoID).First(&video).Error; err != nil {
        return false, err
    }

    // Check if the file exists in S3/R2, exists = boolean
    exists, err := s3util.ObjectExists("toktikp2-video", video.Filename)

    if(!exists) {
        // If the object exists, we can return true
        database.DB.Delete(&video) // Delete the video record if it doesn't exist in S3
    }


	

    return exists, err
}