package s3util

import (
    "context"
    "fmt"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)


func GeneratePresignedPutURL(bucketName, objectKey string , expires time.Duration) (string, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
        config.WithRegion("auto"),
    )
    if err != nil {
        return "", fmt.Errorf("failed to load config: %w", err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
    })

    presignClient := s3.NewPresignClient(client)
    params := &s3.PutObjectInput{
        Bucket:      aws.String(bucketName),
        Key:         aws.String(objectKey),
        ContentType:   aws.String("image/webp"),
    }
    presignedReq, err := presignClient.PresignPutObject(context.TODO(), params, s3.WithPresignExpires(expires))
    if err != nil {
        return "", fmt.Errorf("failed to presign PUT URL: %w", err)
    }
    return presignedReq.URL, nil
}