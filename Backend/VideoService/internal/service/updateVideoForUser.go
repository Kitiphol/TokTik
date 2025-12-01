package service

import (
    "errors"
    "VideoService/internal/entity"
    "VideoService/internal/DTO"
    "VideoService/internal/database"
    "github.com/google/uuid"
)

func UpdateVideoForUser(userID, videoID uuid.UUID, req DTO.UpdateVideoRequest) error {
    var video entity.Video
    if err := database.DB.Where("user_id = ? AND id = ?", userID, videoID).First(&video).Error; err != nil {
        return errors.New("video not found")
    }
    if req.Title != "" {
        video.Title = req.Title
    }
    if err := database.DB.Save(&video).Error; err != nil {
        return err
    }
    return nil
}