package s3util

import (
    "bytes"
    "context"
    "fmt"
    "log"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func UploadToR2(bucketName, objectKey string, content []byte) error {
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

    _, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(objectKey),
        Body:   bytes.NewReader(content),
        ContentType: aws.String("application/vnd.apple.mpegurl"),
    })
    if err != nil {
        return fmt.Errorf("failed to upload object: %w", err)
    }
    log.Printf("[DEBUG] Uploaded %s to bucket %s", objectKey, bucketName)
    return nil
}