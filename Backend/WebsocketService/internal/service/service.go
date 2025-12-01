package service

import (
    "WebsocketService/internal/entity"
    "WebsocketService/internal/middleware"
    "WebsocketService/internal/database"
    "github.com/gin-gonic/gin"
    "net/http"
    "github.com/google/uuid"
	"gorm.io/gorm"
    "WebsocketService/internal/publisher"
    "fmt"
    "gorm.io/gorm/clause"
    "time"
)




type CommentResponse struct {
    ID       uuid.UUID `json:"id"`
    Content  string    `json:"content"`
    Username string    `json:"username"`
    IsUser   bool      `json:"isUser"`
}



func AddOrUpdateComment(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	videoIDStr := c.Param("videoID")
	videoID, err := uuid.Parse(videoIDStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid comment"})
		return
	}

	// Create the comment
	comment := entity.Comment{
		Content: req.Content,
		UserID:  userID,
		VideoID: videoID,
	}

	if err := database.DB.Create(&comment).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add comment"})
		return
	}

	// Fetch user for response
	var user entity.User
	database.DB.First(&user, "id = ?", userID)

	// Prepare response
	response := CommentResponse{
        ID:       comment.ID,
        Content:  comment.Content,
        Username: user.Username,
        IsUser:   true,
    }

	// ✅ Notify relevant users (likes + comments on the video)
	go NotifyRelevantUsers(
        videoID, userID, 
        fmt.Sprintf("%s commented on a video you interacted with", user.Username))

    
	var video entity.Video
	go BroadcastNewComment(videoID, response, video.Title)

	// Send back JSON
	c.JSON(http.StatusOK, gin.H{
		"message": "Comment added/updated",
		"comment": response,
	})
}


func BroadcastNewComment(videoID uuid.UUID, commentDTO CommentResponse, title string) {
    publisher.Publish("notifications", map[string]interface{}{
    "type": "video:comment",
    "data": map[string]interface{}{
        "videoID":    videoID.String(),
        "videoTitle": title,
        "comment":    commentDTO,
    },
})
}

func NotifyRelevantUsers(videoID uuid.UUID, actorID uuid.UUID, message string) {
	var userIDs []uuid.UUID

	database.DB.Raw(`
		SELECT DISTINCT user_id FROM likes WHERE video_id = ? AND user_id != ?
		UNION
		SELECT DISTINCT user_id FROM comments WHERE video_id = ? AND user_id != ?
	`, videoID, actorID, videoID, actorID).Scan(&userIDs)

	for _, uid := range userIDs {
		notif := entity.Notification{
			UserID:  uid,
			VideoID: videoID,
			Message: message,
		}
		database.DB.Create(&notif)

		// Broadcast to WebSocket users (optional)
		publisher.Publish("notifications", map[string]interface{}{
        "to":   uid.String(),
        "type": "notification",
        "data": map[string]interface{}{
            "id":      notif.ID.String(),
            "message": notif.Message,
            "read":      notif.Read,
            "createdAt": notif.CreatedAt.Format(time.RFC3339),
        },
        })
	}
}







// AddOrUpdateLike toggles a like and broadcasts the updated count.
// It wraps everything in a transaction + row lock to avoid race conditions.
func AddOrUpdateLike(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }

    videoID, err := uuid.Parse(c.Param("videoID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }

    fmt.Printf("[📝] ToggleLike start: userID=%s videoID=%s\n", userID, videoID)

    tx := database.DB.Begin()
    defer func() {
        if r := recover(); r != nil {
            tx.Rollback()
            panic(r)
        }
    }()

    // Ensure video exists
    var video entity.Video
    if err := tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        First(&video, "id = ?", videoID).Error; err != nil {
        tx.Rollback()
        fmt.Printf("[❌] Video lookup failed: %v\n", err)
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }

    // Try to delete existing like
    var like entity.Like
    err = tx.Clauses(clause.Locking{Strength: "UPDATE"}).
        Where("user_id = ? AND video_id = ?", userID, videoID).
        First(&like).Error

    var hasLiked bool
    if err == nil {
        // Like exists → unlike
        fmt.Printf("[🔄] Unlike: likeID=%s userID=%s videoID=%s\n", like.ID, userID, videoID)
        if err := tx.Delete(&like).Error; err != nil {
            tx.Rollback()
            fmt.Printf("[❌] Failed to delete like: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to unlike"})
            return
        }
        hasLiked = false

    } else if err == gorm.ErrRecordNotFound {
        // No existing like → create one
        fmt.Printf("[➕] Creating new like for userID=%s videoID=%s\n", userID, videoID)
        newLike := entity.Like{UserID: userID, VideoID: videoID}
        if err := tx.Clauses(clause.OnConflict{DoNothing: true}).
            Create(&newLike).Error; err != nil {
            tx.Rollback()
            fmt.Printf("[❌] Failed to create like: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to like"})
            return
        }
        hasLiked = true

    } else {
        // Unexpected error
        tx.Rollback()
        fmt.Printf("[❌] Unexpected error checking like: %v\n", err)
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check like"})
        return
    }


        // Re-count total likes to ensure accuracy
        var totalLikes int64
        if err := tx.Model(&entity.Like{}).
            Where("video_id = ?", videoID).
            Count(&totalLikes).Error; err != nil {
            tx.Rollback()
            fmt.Printf("[❌] Failed to recount likes: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get like count"})
            return
        }

        // Update Video.total_like_count
        if err := tx.Model(&entity.Video{}).
            Where("id = ?", videoID).
            Update("total_like_count", totalLikes).Error; err != nil {
            tx.Rollback()
            fmt.Printf("[❌] Failed to update video like count: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update like count"})
            return
        }

        // Commit the transaction
        if err := tx.Commit().Error; err != nil {
            fmt.Printf("[❌] Transaction commit failed: %v\n", err)
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save like change"})
            return
        }

            var liker entity.User
            database.DB.Select("username").First(&liker, "id = ?", userID)

            if hasLiked {
                // Only notify on like, not on unlike
                go NotifyRelevantUsers(
                    videoID,
                    userID,
                    fmt.Sprintf("%s liked video %q", liker.Username, video.Title),
                )

                // Also only publish notification on like
                publisher.Publish("notifications", map[string]interface{}{
                    "type": "video:like",
                    "data": map[string]interface{}{
                        "videoID":        video.ID.String(),
                        "videoTitle":     video.Title,
                        "totalLikeCount": totalLikes,
                        "hasLiked":       hasLiked,
                        "username":       liker.Username,
                    },
                })
            } else {
                // Optional: log that unlike happened, no notification sent
                fmt.Printf("[ℹ️] User %s unliked video %q - no notification sent\n", liker.Username, video.Title)
            }


        

        // Broadcast via Redis/WebSocket
        publisher.Publish("notifications", map[string]interface{}{
            "type": "video:like",
            "data": map[string]interface{}{
                "videoID":        video.ID.String(),
                "videoTitle":     video.Title,
                "totalLikeCount": totalLikes,
                "hasLiked":       hasLiked,
                "username":       liker.Username,
            },
        })

        fmt.Printf("[✅] ToggleLike done: hasLiked=%v totalLikes=%d\n", hasLiked, totalLikes)
        c.JSON(http.StatusOK, gin.H{
            "likes":    totalLikes,
            "hasLiked   ": hasLiked,
        })
    }







// Add or update a view on a video
func AddOrUpdateView(c *gin.Context) {
    _, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }

    // Increment view count in Video
    if err := database.DB.Model(&entity.Video{}).Where("id = ?", videoID).UpdateColumn("total_view_count", gorm.Expr("total_view_count + ?", 1)).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update view count"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "View added/updated"})
}

// Get total likes for a video
func GetTotalLikes(c *gin.Context) {
    videoID, err := uuid.Parse(c.Param("videoID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    // fetch video
    var video entity.Video
    if err := database.DB.First(&video, "id = ?", videoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }
    // check user like
    userID, _ := middleware.GetUserIDFromContext(c)
    var count int64
    if userID != uuid.Nil {
        database.DB.Model(&entity.Like{}).
            Where("user_id = ? AND video_id = ?", userID, videoID).
            Count(&count)
    }
    c.JSON(http.StatusOK, gin.H{"likes": video.TotalLikeCount, "hasLiked": count > 0})
}


// Get total views for a video
func GetTotalViews(c *gin.Context) {
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    var video entity.Video
    if err := database.DB.First(&video, "id = ?", videoID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Video not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"views": video.TotalViewCount})
}

// List all comments for a video
func ListComments(c *gin.Context) {
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    var comments []entity.Comment
    if err := database.DB.Preload("User").Where("video_id = ?", videoID).Find(&comments).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch comments"})
        return
    }

    userID, _ := middleware.GetUserIDFromContext(c)

    type CommentResponse struct {
        ID       uuid.UUID `json:"id"`
        Content  string    `json:"content"`
        Username string    `json:"username"`
        IsUser   bool      `json:"isUser"`
    }

    var res []CommentResponse
    for _, comment := range comments {
        res = append(res, CommentResponse{
            ID: comment.ID,
            Content: comment.Content,
            Username: comment.User.Username,
            IsUser: comment.UserID == userID,
        })
    }

    c.JSON(http.StatusOK, gin.H{"comments": res})
}



// Delete a like on a video (unlike)
// DeleteLike removes a user's like and decrements the count
func DeleteLike(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoID, err := uuid.Parse(c.Param("videoID"))
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }

    // Delete record
    res := database.DB.Where("user_id = ? AND video_id = ?", userID, videoID).Delete(&entity.Like{})
    if res.Error != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete like"})
        return
    }
    if res.RowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Like not found"})
        return
    }

    // Decrement count
    if err := database.DB.Model(&entity.Video{}).
        Where("id = ?", videoID).
        UpdateColumn("total_like_count", gorm.Expr("GREATEST(total_like_count - 1, 0)")).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decrement like count"})
        return
    }

    var video entity.Video
    database.DB.First(&video, "id = ?", videoID)
    


    publisher.Publish("notifications", map[string]interface{}{
        "type": "video:like",
        "data": map[string]interface{}{
            "videoID":        video.ID.String(),
            "totalLikeCount": video.TotalLikeCount,
        },
    })

    c.JSON(http.StatusOK, gin.H{"likes": video.TotalLikeCount, "hasLiked": false})
}

// Delete a comment for a video (user specified)
func DeleteComment(c *gin.Context) {
    userID, err := middleware.GetUserIDFromContext(c)
    if err != nil {
        c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
        return
    }
    videoIDStr := c.Param("videoID")
    videoID, err := uuid.Parse(videoIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid videoID"})
        return
    }
    commentIDStr := c.Param("commentID")
    commentID, err := uuid.Parse(commentIDStr)
    if err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid commentID"})
        return
    }

    // Only allow the user who created the comment to delete it
    var comment entity.Comment
    if err := database.DB.First(&comment, "id = ? AND video_id = ? AND user_id = ?", commentID, videoID, userID).Error; err != nil {
        c.JSON(http.StatusNotFound, gin.H{"error": "Comment not found or not authorized"})
        return
    }
    if err := database.DB.Delete(&comment).Error; err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete comment"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"message": "Comment deleted"})
}



// GET /notifications/unread
func GetNotifications(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	var notifs []entity.Notification
	if err := database.DB.
		Where("user_id = ? AND read = false", userID).
		Order("created_at DESC").
		Find(&notifs).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"notifications": notifs})
}

// POST /notifications/:id/read
func GetReadNotifications(c *gin.Context) {
	userID, err := middleware.GetUserIDFromContext(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	notifID, err := uuid.Parse(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid notification ID"})
		return
	}

	// Mark as read
	res := database.DB.Model(&entity.Notification{}).
		Where("id = ? AND user_id = ?", notifID, userID).
		Update("read", true)
	if res.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to mark as read"})
		return
	}
	if res.RowsAffected == 0 {
		c.JSON(http.StatusNotFound, gin.H{"error": "Notification not found or not yours"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification marked as read"})
}
