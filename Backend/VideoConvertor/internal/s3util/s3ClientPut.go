package s3util

import (
    "context"
    "fmt"
    "log"
    "os"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)


// UploadFile uploads the local file to the specified bucket and object key in R2.
func UploadFile(bucketName, objectKey, localFilePath string) error {
    log.Printf("[DEBUG] UploadFile: bucket=%s, key=%s, localFile=%s", bucketName, objectKey, localFilePath)
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

    // Open local file
    file, err := os.Open(localFilePath)
    if err != nil {
        return fmt.Errorf("failed to open local file: %w", err)
    }
    defer file.Close()

    // Upload to R2
    _, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
       Bucket:      aws.String(bucketName),
        Key:         aws.String(objectKey),
        Body:       file,
    })
    if err != nil {
        return fmt.Errorf("failed to upload file to R2: %w", err)
    }

    log.Printf("[DEBUG] Uploaded file %s to bucket %s with key %s", localFilePath, bucketName, objectKey)
    return nil
}
