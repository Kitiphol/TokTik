package s3util

import (
    "context"
    "fmt"
    "io"
    "log"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/joho/godotenv"
)

var (
    _ = godotenv.Load() // Load environment variables from .env file
    accountId       = os.Getenv(("R2_ACCOUNT_ID"))
    accessKeyId     = os.Getenv("R2_ACCESS_KEY_ID")
    accessKeySecret = os.Getenv("R2_ACCESS_KEY_SECRET")
)

func DownloadFromR2(bucketName, objectKey, localFilePath string) error {
    // 1. Load config and create client
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
        config.WithRegion("auto"),
    )
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
    })

    // Download the object
    out, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: &bucketName,
        Key:    aws.String(objectKey),
    })
    if err != nil {
        return fmt.Errorf("failed to download object: %w", err)
    }
    defer out.Body.Close()

    // Save to local file --> /tmp/file.mp4
    f, err := os.Create(localFilePath)
    if err != nil {
        return fmt.Errorf("failed to create file: %w", err)
    }
    defer f.Close()

    _, err = io.Copy(f, out.Body)
    if err != nil {
        return fmt.Errorf("failed to save file: %w", err)
    }

    log.Printf("Downloaded %s from bucket %s to %s\n", objectKey, bucketName, localFilePath)
    return nil
}

/*
 	bucket := "your-bucket-name"
    objectKey := "path/in/bucket/video.mp4"
    localPath := "/tmp/video.mp4" // or any path you want to save to

    err := s3util.DownloadFromR2(bucket, objectKey, localPath)
    if err != nil {
        log.Fatalf("Download failed: %v", err)
    }
    log.Println("Download successful!")
*/