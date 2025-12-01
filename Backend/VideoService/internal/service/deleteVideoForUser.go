package service

import (
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"github.com/google/uuid"
)


func DeleteVideoForUser(userID, videoID uuid.UUID) error {
    tx := database.DB.Begin()

    // Step 1: Delete likes
    if err := tx.Where("video_id = ?", videoID).Delete(&entity.Like{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Step 2: Delete comments
    if err := tx.Where("video_id = ?", videoID).Delete(&entity.Comment{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    // Step 3: Delete the video (with user ID check for ownership)
    if err := tx.Where("user_id = ? AND id = ?", userID, videoID).Delete(&entity.Video{}).Error; err != nil {
        tx.Rollback()
        return err
    }

    return tx.Commit().Error
}
