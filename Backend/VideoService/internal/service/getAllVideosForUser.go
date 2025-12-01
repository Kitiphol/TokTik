package service

import (
	"VideoService/internal/database"
	"VideoService/internal/entity"
	"VideoService/internal/s3util"
	"log"
	"time"

	"github.com/google/uuid"
)

type CommentInfo struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	UserID  string `json:"userID"`
}

type VideoThumbInfoWithComments struct {
	VideoID        string        `json:"videoID"`
	Title          string        `json:"title"`
	TotalLikeCount int64         `json:"totalLikeCount"`
	TotalViewCount uint64        `json:"totalViewCount"`
	ThumbnailURL   string        `json:"thumbnailURL"`
	Comments       []CommentInfo `json:"comments"`
}


type VideoThumbInfo struct {
	VideoID        string `json:"videoID"`
	Title          string `json:"title"`
	TotalLikeCount int64  `json:"totalLikeCount"`
	TotalViewCount uint64 `json:"totalViewCount"`
	ThumbnailURL   string `json:"thumbnailURL"`
}

func GetAllVideosForUser(userID uuid.UUID) ([]VideoThumbInfoWithComments, error) {
	log.Printf("[DEBUG] GetAllVideosForUser called for userID=%s", userID.String())

	var videos []entity.Video
	if err := database.DB.Where("user_id = ?", userID).Find(&videos).Error; err != nil {
		log.Printf("[ERROR] DB error in GetAllVideosForUser: %v", err)
		return nil, err
	}

	var results []VideoThumbInfoWithComments

	for _, video := range videos {
		// Get thumbnail
		thumbKey := userID.String() + "/" + video.ID.String() + "/thumbnail.jpg"
		url, err := s3util.GeneratePresignedGetURL("toktikp2-video", thumbKey, 20*time.Minute)
		if err != nil {
			log.Printf("[WARN] Could not generate presigned URL for %s: %v", thumbKey, err)
			continue
		}

		// 🔥 Manually load comments for this video
		var comments []entity.Comment
		if err := database.DB.Where("video_id = ?", video.ID).Find(&comments).Error; err != nil {
			log.Printf("[WARN] Failed to load comments for video %s: %v", video.ID.String(), err)
			continue
		}

		var commentInfos []CommentInfo
		for _, comment := range comments {
			commentInfos = append(commentInfos, CommentInfo{
				ID:      comment.ID.String(),
				Content: comment.Content,
				UserID:  comment.UserID.String(),
			})
		}

		results = append(results, VideoThumbInfoWithComments{
			VideoID:        video.ID.String(),
			Title:          video.Title,
			TotalLikeCount: int64(video.TotalLikeCount),
			TotalViewCount: uint64(video.TotalViewCount),
			ThumbnailURL:   url,
			Comments:       commentInfos,
		})
	}

	return results, nil
}
