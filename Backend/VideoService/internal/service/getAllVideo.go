package service

import (
    "VideoService/internal/entity"
    "VideoService/internal/database"
    "VideoService/internal/s3util"
    "log"
    "time"
    "encoding/json"
)


func GetAllVideosWithThumbnails() ([]VideoThumbInfo, error) {
    log.Println("[DEBUG] GetAllVideosWithThumbnails called")
    var videos []entity.Video
    err := database.DB.Find(&videos).Error
    if err != nil {
        log.Printf("[ERROR] DB error in GetAllVideosWithThumbnails: %v", err)
        return nil, err
    }
    log.Printf("[DEBUG] Found %d videos", len(videos))

    var results []VideoThumbInfo
    for _, video := range videos {
        log.Printf("[DEBUG] Processing videoID=%s userID=%s", video.ID.String(), video.UserID.String())
        thumbKey := video.UserID.String() + "/" + video.ID.String() + "/thumbnail.jpg"
        url, err := s3util.GeneratePresignedGetURL("toktikp2-video", thumbKey, 2 * time.Hour)
        if err != nil {
            log.Printf("[WARN] Could not generate presigned URL for %s: %v", thumbKey, err)
            continue
        }
        log.Printf("[DEBUG] Generated thumbnail URL for videoID=%s: %s", video.ID.String(), url)
        results = append(results, VideoThumbInfo{
            VideoID:      video.ID.String(),
            ThumbnailURL: url,
            Title:        video.Title,
            TotalLikeCount: int64(video.TotalLikeCount),
            TotalViewCount: uint64(video.TotalViewCount),
        })
    }
    log.Printf("[DEBUG] Returning %d video thumbnails", len(results))

    jsonBytes, _ := json.MarshalIndent(results, "", "  ")
    log.Printf("JSON sent to frontend: %s", string(jsonBytes))
    return results, nil
}