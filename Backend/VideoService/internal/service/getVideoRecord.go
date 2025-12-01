package service

import (
    "VideoService/internal/database"
    "VideoService/internal/entity"
    "github.com/google/uuid"
)

func GetVideoRecord(videoID uuid.UUID) (*entity.Video, error) {
    var video entity.Video
    if err := database.DB.Where("id = ?", videoID).First(&video).Error; err != nil {
        return nil, err
    }
    return &video, nil
}