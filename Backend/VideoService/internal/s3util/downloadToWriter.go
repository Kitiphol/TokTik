package s3util

import (
    "context"
    "fmt"
    "io"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
    "github.com/aws/aws-sdk-go-v2/service/s3"
)

func DownloadToWriter(bucketName, objectKey string, w io.Writer) error {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
        config.WithRegion("auto"),
    )
    if err != nil {
        return fmt.Errorf("failed to load config: %w", err)
    }


	fmt.Print("Using Cloudflare R2 with account ID: ", accountId, "\n")
	fmt.Print("Using access key ID: ", accessKeyId, "\n")
	fmt.Print("Using access key secret: ", accessKeySecret, "\n")
	fmt.Print("Using bucket name: ", bucketName, "\n")
	fmt.Print("Using object key: ", objectKey, "\n")



    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
    })

    out, err := client.GetObject(context.TODO(), &s3.GetObjectInput{
        Bucket: &bucketName,
        Key:    aws.String(objectKey),
    })
    if err != nil {
        return fmt.Errorf("failed to download object: %w", err)
    }
    defer out.Body.Close()

    _, err = io.Copy(w, out.Body)
    if err != nil {
        return fmt.Errorf("failed to copy object to writer: %w", err)
    }
    return nil
}