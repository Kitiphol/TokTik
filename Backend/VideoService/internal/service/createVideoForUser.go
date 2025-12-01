package service

import (
	"VideoService/internal/DTO"
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"VideoService/internal/s3util"
	"fmt"
	"log"
	"time"

	"github.com/google/uuid"
)


// func CreateVideoForUser(userID uuid.UUID, req DTO.CreateVideoRequest) (string, string, error) {
//     videoID := uuid.New()
//     log.Printf("[DEBUG] Creating video for user: userID=%s, videoID=%s, filename=%s", userID.String(), videoID.String(), req.Filename)

//     // uuid_{filename}
//     // uniqueFilename := fmt.Sprintf("videos/%s_%s", videoID.String(), req.Filename)

//     uniqueFilename := fmt.Sprintf("%s/%s/%s", userID.String(), videoID.String(), req.Filename)
//     // uniqueFilename := fmt.Sprintf("%s", req.Filename)
//     video := entity.Video{
//         ID:           videoID,
//         Title:        req.Title,
//         Filename:     uniqueFilename,
//         UserID:       userID,
//     }
//     if err := database.DB.Create(&video).Error; err != nil {
//         return "", "", err
//     }
//     log.Printf("[DEBUG] Inserted video: %+v", video)
//     log.Printf("[DEBUG] Video created: ID=%s, Title=%s, Filename=%s", video.ID.String(), video.Title, video.Filename)

//     log.Printf("[DEBUG] Creating video: bucket=%s, key=%s", "toktikp2-video", uniqueFilename)
//     presignURL, err := s3util.GeneratePresignedPutURL("toktikp2-video", uniqueFilename, 300 * time.Second)
//     if err != nil {
//         return "", "", err
//     }
//     log.Printf("Generated presigned URL: %s\n", presignURL)
//     return presignURL, videoID.String(), nil
// }


func CreateVideoForUser(userID uuid.UUID, req DTO.CreateVideoRequest) (string, string, error) {

    log.Printf("[DEBUG] Creating video for user: userID=%s, title=%s, filename=%s", userID.String(), req.Title, req.Filename)
    // Do NOT generate videoID here!
    video := entity.Video{
        Title:    req.Title,
        UserID:   userID,
        Filename: req.Filename,
        // Filename will be set after we have the ID
    }

    // Insert to DB (GORM will call BeforeCreate and set the ID)
    if err := database.DB.Create(&video).Error; err != nil {
        return "", "", err
    }

    // Now video.ID is set by GORM
    uniqueFilename := fmt.Sprintf("%s/%s/%s", userID.String(), video.ID.String(), req.Filename)
    video.Filename = uniqueFilename

    playlistFilePath := fmt.Sprintf("%s/%s/playlist.m3u8", userID.String(), video.ID.String())

    // Update the filename in the database
    if err := database.DB.Model(&video).Update("filename", playlistFilePath).Error; err != nil {
        return "", "", err
    }

    log.Printf("[DEBUG] Inserted video: %+v", video)
    log.Printf("[DEBUG] Video created: ID=%s, Title=%s, Filename=%s", video.ID.String(), video.Title, video.Filename)

    log.Printf("[DEBUG] Creating video: bucket=%s, (path to file) Objectkey=%s", "toktikp2-video", video.Filename)
    log.Printf("[DEBUG] This is the unique filename that will be used in S3/R2: %s", uniqueFilename)
    log.Printf("[DEBUG] This is the filename stored in the database: %s", video.Filename)
    presignURL, err := s3util.GeneratePresignedPutURL("toktikp2-video", playlistFilePath, 1*time.Hour)
    if err != nil {
        return "", "", err
    }
    log.Printf("Generated PUT presigned URL: %s\n", presignURL)
    return presignURL, video.ID.String(), nil
}