package s3util

import (
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

// DeleteVideoFromR2 deletes an object from a Cloudflare R2 bucket.
func DeleteVideoFromR2(bucketName, objectKey string) error {
    log.Printf("[DEBUG] DeleteVideoFromR2: bucket=%s, key=%s", bucketName, objectKey)
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
            accessKeyId,
            accessKeySecret,
            "",
        )),
        config.WithRegion("auto"),
    )
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
    })

    _, err = client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
        Bucket: &bucketName,
        Key:    aws.String(objectKey),
    })
    if err != nil {
        return fmt.Errorf("failed to delete object: %w", err)
    }

    log.Printf("[DEBUG] Deleted object %s from bucket %s", objectKey, bucketName)
    return nil
}