package service

import (
    "VideoService/internal/database"
    "VideoService/internal/entity"
    "VideoService/internal/s3util"
    "bytes"
    "fmt"
    "log"
    "strings"
    "time"
    "gorm.io/gorm"
    "VideoService/internal/publisher"
)

type CommentDTO struct {
    ID       string `json:"id"`
    Content  string `json:"content"`
    Username string `json:"username"`
}

type VideoDetailResponse struct {
    Playlist        string       `json:"playlist"`
    Title           string       `json:"title"`
    TotalViewCount  uint         `json:"totalViewCount"`
    TotalLikeCount  int          `json:"totalLikeCount"`
    Comments        []CommentDTO `json:"comments"`
}

// // Utility: Save playlist to Desktop
// func SavePlaylistToDesktop(filename, content string) error {
//     desktopDir := "/Users/kitipholk/Desktop"
//     if _, err := os.Stat(desktopDir); os.IsNotExist(err) {
//         log.Printf("[DEBUG] Desktop directory does not exist, creating: %s", desktopDir)
//         if err := os.MkdirAll(desktopDir, 0755); err != nil {
//             log.Printf("[ERROR] Failed to create Desktop directory: %v", err)
//             return err
//         }
//     }
//     desktopPath := filepath.Join(desktopDir, filename)
//     log.Printf("[DEBUG] Saving playlist to %s", desktopPath)
//     f, err := os.Create(desktopPath)
//     if err != nil {
//         log.Printf("[ERROR] Failed to create file: %v", err)
//         return err
//     }
//     defer f.Close()
//     _, err = f.WriteString(content)
//     if err != nil {
//         log.Printf("[ERROR] Failed to write to file: %v", err)
//         return err
//     }
//     log.Printf("[DEBUG] Successfully saved playlist to %s", desktopPath)
//     return nil
// }

// Returns the playlist.m3u8 content with all .ts segment lines replaced by presigned URLs.
func GetVideoForUser(videoID string) (*VideoDetailResponse, error) {
    log.Printf("[DEBUG] GetVideoForUser called with videoID=%s", videoID)

    var video entity.Video
    log.Printf("[DEBUG] Querying video with id=%s", videoID)
    if err := database.DB.Where("id = ?", videoID).First(&video).Error; err != nil {
        log.Printf("[ERROR] Video not found: %v", err)
        return nil, fmt.Errorf("video not found: %w", err)
    }
    log.Printf("[DEBUG] Found video: %+v", video)

    if err := database.DB.Model(&entity.Video{}).
        Where("id = ?", videoID).
        UpdateColumn("total_view_count", gorm.Expr("total_view_count + ?", 1)).Error; err != nil {
        log.Printf("[ERROR] Failed to increment view count: %v", err)
    } else {
        video.TotalViewCount++ 
    }

     var updated entity.Video
    database.DB.First(&updated, "id = ?", videoID)


    publisher.Publish("notifications", map[string]interface{}{
        "type": "video:view",
        "data": map[string]interface{}{
            "videoID":        video.ID.String(),
            "videoTitle":     video.Title,
            "totalViewCount": video.TotalViewCount,
        },
    })

    
    var likeCount int64
    log.Printf("[DEBUG] Counting likes for video_id=%s", videoID)
    if err := database.DB.Model(&entity.Like{}).Where("video_id = ?", videoID).Count(&likeCount).Error; err != nil {
        log.Printf("[ERROR] Failed to count likes: %v", err)
    }


    var comments []entity.Comment
    log.Printf("[DEBUG] Querying comments for video_id=%s", videoID)
    if err := database.DB.Where("video_id = ?", videoID).Find(&comments).Error; err != nil {
        log.Printf("[ERROR] Failed to get comments: %v", err)
    }
    log.Printf("[DEBUG] Found %d comments", len(comments))

    var commentDTOs []CommentDTO
    for _, c := range comments {
        var user entity.User
        log.Printf("[DEBUG] Querying user for comment user_id=%s", c.UserID.String())
        if err := database.DB.Select("username").Where("id = ?", c.UserID).First(&user).Error; err != nil {
            log.Printf("[ERROR] Failed to get user for comment: %v", err)
            user.Username = "Unknown"
        }
        commentDTOs = append(commentDTOs, CommentDTO{
            ID:       c.ID.String(),
            Content:  c.Content,
            Username: user.Username,
        })
    }
    log.Printf("[DEBUG] Built %d commentDTOs", len(commentDTOs))

    // Download and rewrite playlist
    userID := video.UserID.String()
    playlistKey := video.Filename
    log.Printf("[DEBUG] Downloading playlist from S3: bucket=toktikp2-video, key=%s", playlistKey)
    var buf bytes.Buffer
    err := s3util.DownloadToWriter("toktikp2-video", playlistKey, &buf)
    if err != nil {
        log.Printf("[ERROR] Failed to download playlist: %v", err)
        return nil, fmt.Errorf("failed to download playlist: %w", err)
    }
    playlistContent := buf.String()
    log.Printf("[DEBUG] Downloaded playlist content length: %d", len(playlistContent))

    lines := strings.Split(playlistContent, "\n")
    for i, line := range lines {
        line = strings.TrimSpace(line)
        if strings.HasSuffix(line, ".ts") {
            segKey := fmt.Sprintf("%s/%s/%s", userID, videoID, line)
            log.Printf("[DEBUG] Generating presigned URL for segment: %s", segKey)
            url, err := s3util.GeneratePresignedGetURL("toktikp2-video", segKey, 10*time.Minute)
            if err != nil {
                log.Printf("[ERROR] Failed to presign segment %s: %v", segKey, err)
                return nil, fmt.Errorf("failed to presign segment %s: %w", segKey, err)
            }
            lines[i] = url
        }
    }
       newPlaylist := strings.Join(lines, "\n")
    log.Printf("[DEBUG] Rewritten playlist content length: %d", len(newPlaylist))

    // Upload the rewritten playlist as playlist2.m3u8 to R2
    playlist2Key := fmt.Sprintf("%s/%s/playlist2.m3u8", userID, videoID)
    err = s3util.UploadToR2("toktikp2-video", playlist2Key, []byte(newPlaylist))
    if err != nil {
        log.Printf("[ERROR] Failed to upload playlist2.m3u8: %v", err)
        return nil, fmt.Errorf("failed to upload playlist2.m3u8: %w", err)
    }
    log.Printf("[DEBUG] Uploaded playlist2.m3u8 to R2 at key: %s", playlist2Key)

    // Generate a presigned URL for playlist2.m3u8
    playlist2Url, err := s3util.GeneratePresignedGetURL("toktikp2-video", playlist2Key, 10*time.Minute)
    if err != nil {
        log.Printf("[ERROR] Failed to generate presigned URL for playlist2.m3u8: %v", err)
        return nil, fmt.Errorf("failed to generate presigned URL for playlist2.m3u8: %w", err)
    }
    log.Printf("[DEBUG] Presigned URL for playlist2.m3u8: %s", playlist2Url)

    log.Printf("[DEBUG] Returning VideoDetailResponse for videoID=%s", videoID)
    return &VideoDetailResponse{
        Playlist:       playlist2Url, 
        Title:          video.Title,
        TotalViewCount: video.TotalViewCount,
        TotalLikeCount: int(likeCount),
        Comments:       commentDTOs,
    }, nil
}