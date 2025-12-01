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


func GeneratePresignedGetURL(bucketName, objectKey string, expires time.Duration) (string, error) {
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
    params := &s3.GetObjectInput{
        Bucket: aws.String(bucketName),
        Key:    aws.String(objectKey),
    }
    presignedReq, err := presignClient.PresignGetObject(context.TODO(), params, s3.WithPresignExpires(expires))
    if err != nil {
        return "", fmt.Errorf("failed to presign GET URL: %w", err)
    }
    return presignedReq.URL, nil
}