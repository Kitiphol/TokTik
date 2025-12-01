package s3util

import (
    "errors"
    "fmt"
    "context"
    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"
    "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/credentials"
	"strings"
)

func ObjectExists(bucket, key string) (bool, error) {
    cfg, err := config.LoadDefaultConfig(context.TODO(),
        config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyId, accessKeySecret, "")),
        config.WithRegion("auto"),
    )
    if err != nil {
        return false, fmt.Errorf("failed to load config: %w", err)
    }

    client := s3.NewFromConfig(cfg, func(o *s3.Options) {
        o.BaseEndpoint = aws.String(fmt.Sprintf("https://%s.r2.cloudflarestorage.com", accountId))
    })

    _, err = client.HeadObject(context.TODO(), &s3.HeadObjectInput{
        Bucket: aws.String(bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        var notFound *types.NotFound
        if errors.As(err, &notFound) {
            return false, nil
        }
        // If error is 404, treat as not found (for some SDKs)
        if strings.Contains(err.Error(), "NotFound") {
            return false, nil
        }
        return false, err
    }

	
	


    return true, nil
}