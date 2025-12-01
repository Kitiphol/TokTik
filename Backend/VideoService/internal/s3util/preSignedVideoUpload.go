// package s3util

// import (
//     "context"
//     "fmt"
//     "time"

//     "github.com/aws/aws-sdk-go-v2/aws"
//     "github.com/aws/aws-sdk-go-v2/config"
//     "github.com/aws/aws-sdk-go-v2/credentials"
//     "github.com/aws/aws-sdk-go-v2/service/s3"
// )

// func GeneratePresignedPutURL(bucketName, objectKey string , expires time.Duration) (string, error) {
//     cfg, err := config.LoadDefaultConfig(context.TODO(),
//         config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
//         config.WithRegion("auto"),
//     )
//     if err != nil {
//         return "", fmt.Errorf("failed to load config: %w", err)
//     }

//     client := s3.NewFromConfig(cfg, func(o *s3.Options) {
//         o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
//     })

//     presignClient := s3.NewPresignClient(client)
//     params := &s3.PutObjectInput{
//         Bucket:      aws.String(bucketName),
//         Key:         aws.String(objectKey),
//     }
//     presignedReq, err := presignClient.PresignPutObject(context.TODO(), params, s3.WithPresignExpires(expires))
//     if err != nil {
//         return "", fmt.Errorf("failed to presign PUT URL: %w", err)
//     }
//     return presignedReq.URL, nil
// }

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

func GeneratePresignedPutURL(bucketName, objectKey string, expires time.Duration) (string, error) {

	// Load AWS config with static credentials and "auto" region (Cloudflare uses "auto")
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return "", fmt.Errorf("failed to load config: %w", err)
	}

	// fmt.Print("Using Cloudflare R2 with account ID: ", accountId, "\n")
	// fmt.Print("Using access key ID: ", accessKeyId, "\n")
	// fmt.Print("Using access key secret: ", accessKeySecret, "\n")
	// fmt.Print("Using bucket name: ", bucketName, "\n")
	// fmt.Print("Using object key: ", objectKey, "\n")

	// Override the endpoint to point to Cloudflare R2
	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
		o.UsePathStyle = false // important for R2 URL style
	})

	fmt.Print("Client: ", client, "\n")
	presignClient := s3.NewPresignClient(client)

	params := &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(objectKey),
		// ContentType: aws.String("video/mp4"),
		// ContentLength: aws.Int64(5636752),
	}

	presignedReq, err := presignClient.PresignPutObject(context.TODO(), params, s3.WithPresignExpires(expires))
	if err != nil {
		return "", fmt.Errorf("failed to presign PUT URL: %w", err)
	}

	return presignedReq.URL, nil
}
