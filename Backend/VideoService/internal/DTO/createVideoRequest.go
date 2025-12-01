package DTO

type CreateVideoRequest struct {
    Title    string `json:"title" binding:"required"`
    Filename string `json:"filename" binding:"required"`
    // S3URL    string `json:"s3_url" binding:"required"`
    // ThumbnailURL string `json:"thumbnail_url" binding:"required"`
}