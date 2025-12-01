package service

import (
	"ThumbnailCreator/internal/s3util"
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
    // "log"

	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// func CreateThumbnail(bucket, segmentKey string) error {
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Starting CreateThumbnail: bucket=%s, segmentKey=%s\n", bucket, segmentKey)

//     // Prepare local paths
//     segmentFilename := filepath.Base(segmentKey)
//     localSegmentPath := filepath.Join(os.TempDir(), segmentFilename)
//     localThumbnailPath := filepath.Join(os.TempDir(), "thumbnail.jpg")
//     thumbnailKey := filepath.Join(filepath.Dir(segmentKey), "thumbnail.jpg")

//     fmt.Printf("[DEBUG] [ThumbnailMaker] Downloading segment from %s/%s to %s\n", bucket, segmentKey, localSegmentPath)
//     // Download .ts segment from R2
//     if err := s3util.DownloadFromR2(bucket, segmentKey, localSegmentPath); err != nil {
//         return fmt.Errorf("failed to download segment: %w", err)
//     }
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Finished downloading segment\n")

//     // Generate thumbnail using ffmpeg
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Running ffmpeg to generate thumbnail at %s\n", localThumbnailPath)
//     cmd := exec.Command("ffmpeg", "-y", "-i", localSegmentPath, "-ss", "00:00:01.000", "-vframes", "1", localThumbnailPath)
//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr
//     if err := cmd.Run(); err != nil {
//         return fmt.Errorf("ffmpeg failed: %w", err)
//     }
//     fmt.Printf("[DEBUG] [ThumbnailMaker] ffmpeg finished, thumbnail generated\n")

//     // Upload thumbnail back to R2
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Uploading thumbnail to %s/%s\n", bucket, thumbnailKey)
//     if err := s3util.UploadFile(bucket, thumbnailKey, localThumbnailPath); err != nil {
//         return fmt.Errorf("upload failed: %w", err)
//     }
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Thumbnail uploaded\n")

//     // Clean up local files
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Cleaning up local files\n")
//     _ = os.Remove(localSegmentPath)
//     _ = os.Remove(localThumbnailPath)

//     fmt.Printf("[DEBUG] [ThumbnailMaker] Thumbnail created and uploaded to %s/%s\n", bucket, thumbnailKey)
//     return nil
// }

// func CreateThumbnail(bucket, segmentKey string) error {
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Starting CreateThumbnail: bucket=%s, segmentKey=%s", bucket, segmentKey)

//     // Prepare local paths
//     segmentFile := filepath.Base(segmentKey)
//     localSeg := filepath.Join(os.TempDir(), segmentFile)
//     localThumb := filepath.Join(os.TempDir(), "thumbnail.jpg")
//     thumbKey := filepath.Join(filepath.Dir(segmentKey), "thumbnail.jpg")

//     fmt.Printf("[DEBUG] [ThumbnailMaker] Downloading segment %s to %s", segmentKey, localSeg)

//     if err := s3util.DownloadFromR2(bucket, segmentKey, localSeg); err != nil {
//         return fmt.Errorf("download segment: %w", err)
//     }
//     fmt.Println("[DEBUG] [ThumbnailMaker] Download complete")

//     // Generate thumbnail
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Running ffmpeg for thumbnail at %s", localThumb)

//     cmd := exec.Command("ffmpeg", "-y", "-ss", localSeg, "-i", "00:00:01.000", "-vframes", "1", localThumb)
//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr
//     if err := cmd.Run(); err != nil {
//         return fmt.Errorf("ffmpeg thumbnail failed: %w", err)
//     }
//     fmt.Println("[DEBUG] [ThumbnailMaker] Thumbnail generation complete")

//     // Upload thumbnail
//     fmt.Printf("[DEBUG] [ThumbnailMaker] Uploading thumbnail to %s/%s", bucket, thumbKey)

//     err := s3util.UploadFile(bucket, thumbKey, localThumb)

//     if (err != nil) {
//         return fmt.Errorf("upload thumbnail: %w", err)
//     }
//     fmt.Println("[DEBUG] [ThumbnailMaker] Upload complete")

//     // Cleanup
//     fmt.Println("[DEBUG] [ThumbnailMaker] Cleaning up local files")
//     _ = os.Remove(localSeg)
//     _ = os.Remove(localThumb)

//     fmt.Printf("[DEBUG] [ThumbnailMaker] CreateThumbnail finished: %s/%s", bucket, thumbKey)
//     return nil
// }









func CreateThumbnail(bucket, segmentKey string) error {
    fmt.Printf("[DEBUG] [ThumbnailMaker] Starting CreateThumbnail: bucket=%s, segmentKey=%s\n", bucket, segmentKey)

    // Prepare local paths
    segmentFile := filepath.Base(segmentKey)
    localSeg := filepath.Join(os.TempDir(), segmentFile)
    localThumb := filepath.Join(os.TempDir(), "thumbnail.jpg")
    thumbKey := filepath.Join(filepath.Dir(segmentKey), "thumbnail.jpg")

    fmt.Printf("[DEBUG] [ThumbnailMaker] Downloading from bucket: %s , segment*(OBJECT KEY) %s to LOCAL: %s\n", bucket, segmentKey, localSeg)

    err := s3util.DownloadFromR2(bucket, segmentKey, localSeg)

    if (err != nil) {
        return fmt.Errorf("Download Segment: %w", err)
    }
    fmt.Println("[DEBUG] [ThumbnailMaker] Download complete")

    // Generate thumbnail
    fmt.Printf("[DEBUG] [ThumbnailMaker] Running ffmpeg for thumbnail at %s\n", localThumb)

    // ✅ Fix: Correct order of arguments
    // cmd := exec.Command("ffmpeg", "-y", "-ss", "00:00:01.5", "-i", localSeg, "-vframes", "1", localThumb)
    cmd := exec.Command(
        "ffmpeg", "-y",
        "-i", localSeg, // input first, then seek
        "-ss", "00:00:00.500", // seek *after* file opens
        "-vframes", "1",
        localThumb,
    )

    // fi, err := os.Stat(localThumb)
    // if err != nil {
    //     return fmt.Errorf("[ERROR]thumbnail file not found: %w", err)
    // }
    // if fi.Size() == 0 {
    //     return fmt.Errorf("[ERROR] thumbnail is 0 bytes; ffmpeg likely failed")
    // }

    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("ffmpeg thumbnail failed: %w", err)
    }
    fmt.Println("[DEBUG] [ThumbnailMaker] Thumbnail generation complete")

    // Upload thumbnail
    fmt.Printf("[DEBUG] [ThumbnailMaker] Uploading thumbnail to %s/%s\n", bucket, thumbKey)

    //uplaod to userID/videoID/thumbnail.jpg
    err2 := s3util.UploadFile(bucket, thumbKey, localThumb)

    if err2 != nil {
        return fmt.Errorf("[Upload Thumbnail] Failed : %w", err2)
    }

    // Cleanup
    fmt.Println("[DEBUG] [ThumbnailMaker] Cleaning up local files")
    _ = os.Remove(localSeg)
    _ = os.Remove(localThumb)

    fmt.Printf("[DEBUG] [ThumbnailMaker] CreateThumbnail finished: %s/%s\n", bucket, thumbKey)
    return nil
}







// func CreateThumbnail(bucket, objectKey string) error {
//     log.Printf("[DEBUG] Starting CreateThumbnail: bucket=%s, objectKey=%s\n", bucket, objectKey)

//     // Download
//     mp4Name  := filepath.Base(objectKey)
//     localMP4 := filepath.Join(os.TempDir(), mp4Name)
//     if err := s3util.DownloadFromR2(bucket, objectKey, localMP4); err != nil {
//         return fmt.Errorf("download MP4: %w", err)
//     }
//     fi, err := os.Stat(localMP4)
//     if err != nil {
//         return fmt.Errorf("stat MP4: %w", err)
//     }
//     log.Printf("[DEBUG] Downloaded %s → %s (%d bytes)\n", objectKey, localMP4, fi.Size())
//     if fi.Size() < 100_000 {
//         return fmt.Errorf("downloaded file too small, wrong key?")
//     }

//     // Generate thumbnail
//     localThumb := filepath.Join(os.TempDir(), "thumbnail.jpg")
//     cmd := exec.Command(
//         "ffmpeg", "-y",
//         "-ss", "00:00:01",
//         "-i", localMP4,
//         "-frames:v", "1",
//         localThumb,
//     )
//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr
//     if err := cmd.Run(); err != nil {
//         return fmt.Errorf("ffmpeg failed: %w", err)
//     }

//     // Verify
//     ti, err := os.Stat(localThumb)
//     if err != nil {
//         return fmt.Errorf("stat thumb: %w", err)
//     }
//     log.Printf("[DEBUG] Created thumbnail %s (%d bytes)\n", localThumb, ti.Size())
//     if ti.Size() == 0 {
//         return fmt.Errorf("thumbnail is still 0 B — ffmpeg produced nothing")
//     }

//     // Upload
//     thumbKey := filepath.Join(filepath.Dir(objectKey), "thumbnail.jpg")
//     if err := s3util.UploadFile(bucket, thumbKey, localThumb); err != nil {
//         return fmt.Errorf("upload thumb: %w", err)
//     }
//     log.Printf("[DEBUG] Uploaded thumbnail to %s/%s\n", bucket, thumbKey)

//     // Cleanup
//     os.Remove(localMP4)
//     os.Remove(localThumb)
//     return nil
// }














// ExtractAndUploadThumbnailToR2 streams a thumbnail from a video URL and uploads as WebP.
func ExtractAndUploadThumbnailToR2(ctx context.Context, s3Client *s3.Client, inputURL, bucketName, objectKey string) error {
    fmt.Printf("[DEBUG] [ThumbnailMaker] Starting ExtractAndUploadThumbnail: url=%s, bucket=%s, key=%s", inputURL, bucketName, objectKey)

    // Create temp file
    tmpFile, err := os.CreateTemp("", "thumb-*.webp")
    if err != nil {
        return fmt.Errorf("create temp file: %w", err)
    }
    tmpPath := tmpFile.Name()
    tmpFile.Close()
    defer func() {
        fmt.Printf("[DEBUG] [ThumbnailMaker] Removing temp file %s", tmpPath)
        os.Remove(tmpPath)
    }()

    // Run ffmpeg
    width, height := "1080", "1920"
    fmt.Printf("[DEBUG] [ThumbnailMaker] Running ffmpeg to extract WebP thumbnail")
    cmd := exec.Command(
        "ffmpeg", "-y",
        "-i", inputURL,
        "-frames:v", "1",
        "-q:v", "2",
        "-vf", fmt.Sprintf("scale=%s:%s", width, height),
        "-c:v", "libwebp",
        tmpPath,
    )
    var errBuf bytes.Buffer
    cmd.Stderr = &errBuf
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("ffmpeg WebP failed: %w, %s", err, errBuf.String())
    }
    fmt.Println("[DEBUG] [ThumbnailMaker] ffmpeg extraction complete")

    // Upload
    file, err := os.Open(tmpPath)
    if err != nil {
        return fmt.Errorf("open temp thumbnail: %w", err)
    }
    defer file.Close()
    // stat, err := file.Stat()
    // if err != nil {
    //     return fmt.Errorf("stat temp file: %w", err)
    // }
    fmt.Printf("[DEBUG] [ThumbnailMaker] Uploading WebP to %s/%s", bucketName, objectKey)
    // _, err = s3Client.PutObject(ctx, &s3.PutObjectInput{
    //     Bucket:        aws.String(bucketName),
    //     Key:           aws.String(objectKey),
    //     Body:          file,
    //     ContentType:   aws.String("image/webp"),
    //     ContentLength: aws.Int64(stat.Size()),
    // })

    err = s3util.UploadFile(bucketName, objectKey, tmpPath)
    if err != nil {
        return fmt.Errorf("upload WebP: %w", err)
    }
    fmt.Println("[DEBUG] [ThumbnailMaker] WebP upload complete")
    return nil
}
