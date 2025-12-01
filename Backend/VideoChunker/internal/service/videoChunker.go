// package service

// import (
//     "fmt"
//     "os"
//     "os/exec"
//     "path/filepath"
//     "strings"
//     "encoding/json"

//     "VideoChuncker/internal/s3util"
//     "VideoChuncker/internal/machineryUtil"
//     "github.com/RichardKnop/machinery/v2"
//     // "github.com/RichardKnop/machinery/v2/tasks"
// )

// func ChunkAndUpload(server *machinery.Server, bucket, objectKey string) error {
//     fmt.Printf("[DEBUG] ChunkAndUpload: bucket=%s, objectKey=%s\n", bucket, objectKey)

//     // Download the video from R2 to a local temp file
//     localInput := filepath.Join(os.TempDir(), filepath.Base(objectKey))
//     fmt.Printf("[DEBUG] Downloading video to %s\n", localInput)
//     if err := s3util.DownloadFromR2(bucket, objectKey, localInput); err != nil {
//         return fmt.Errorf("failed to download video: %w", err)
//     }

//     // Use the objectKey prefix (userID/videoID) for chunk keys
//     prefix := filepath.Dir(objectKey)
//     fmt.Printf("[DEBUG] Calculated prefix: %s\n", prefix)

//     // userID/videoID/chunk_XXX.mp4
//     name := strings.TrimSuffix(filepath.Base(objectKey), filepath.Ext(objectKey))
//     outputDir := filepath.Join(os.TempDir(), name+"_output")
//     fmt.Printf("[DEBUG] Output directory for chunks: %s\n", outputDir)

//     // Chunk the video into segments
//     fmt.Printf("[DEBUG] Chunking video to %s\n", outputDir)
//     if err := ChunkVideo(localInput, outputDir); err != nil {
//         _ = DeleteLocalFile(localInput)
//         return fmt.Errorf("failed to chunk video: %w", err)
//     }

//     // Upload each chunked file to R2 as {userID}/{videoID}/chunk_XXX.mp4
//     files, err := os.ReadDir(outputDir)
//     if err != nil {
//         _ = DeleteLocalFile(localInput)
//         _ = DeleteLocalFile(outputDir)
//         return fmt.Errorf("failed to read output dir: %w", err)
//     }

//     var chunkKeys []string
//     for _, f := range files {
//         if f.IsDir() {
//             continue
//         }
//         chunkPath := filepath.Join(outputDir, f.Name())
//         chunkKey := filepath.Join(prefix, f.Name())
//         fmt.Printf("[DEBUG] For chunk file: %s\n", f.Name())
//         fmt.Printf("[DEBUG]   chunkPath: %s\n", chunkPath)
//         fmt.Printf("[DEBUG]   chunkKey: %s\n", chunkKey)
//         fmt.Printf("[DEBUG]   Uploading chunk: %s -> %s/%s\n", chunkPath, bucket, chunkKey)
//         if err := s3util.UploadFile(bucket, chunkKey, chunkPath); err != nil {
//             fmt.Printf("Failed to upload chunk %s: %v\n", chunkPath, err)
//         }
//         chunkKeys = append(chunkKeys, chunkKey)
//     }
//     fmt.Printf("[DEBUG] All chunk keys: %v\n", chunkKeys)

//     // Serialize chunkKeys to JSON for the next service
//     chunkKeysJSON, err := json.Marshal(chunkKeys)
//     if err != nil {
//         return fmt.Errorf("failed to marshal chunk keys: %w", err)
//     }
//     fmt.Printf("[DEBUG] chunkKeysJSON: %s\n", string(chunkKeysJSON))

//     // Notify the next service (e.g., Convertor) with all chunk keys
//     // if server != nil {
//     //     fmt.Printf("[DEBUG] Sending task to convertor with bucket: %s and chunkKeysJSON: %s\n", bucket, string(chunkKeysJSON))
//     //     err = machineryUtil.SendTaskToConvertor(
//     //         server,
//     //         "ConvertChunks",
//     //         tasks.Arg{Type: "string", Value: bucket},
//     //         tasks.Arg{Type: "string", Value: string(chunkKeysJSON)},
//     //     )
//     //     if err != nil {
//     //         fmt.Printf("Failed to send convert task: %v\n", err)
//     //     }
//     // }
//     if server != nil {

//     fmt.Printf("[DEBUG: Sending message to Convertor] Trying to send task to convertor with bucket: %s and chunkKeysJSON: %s\n", bucket, string(chunkKeysJSON))
//     err = machineryUtil.SendTaskToConvertor(
//         server,
//         bucket,
//         string(chunkKeysJSON),
//     )

//     fmt.Printf("[DEBUG: Sending message to Convertor] Task SUCCESSFULLY SENT to convertor with bucket: %s and chunkKeysJSON: %s\n", bucket, string(chunkKeysJSON))
//     if err != nil {
//         fmt.Printf("Failed to send convert task: %v\n", err)
//     }
// }

//     // Delete the original video from R2 after chunking
//     fmt.Printf("[DEBUG] Deleting original video from R2: %s/%s\n", bucket, objectKey)
//     if err := s3util.DeleteVideoFromR2(bucket, objectKey); err != nil {
//         fmt.Printf("Failed to delete original video from R2: %v\n", err)
//     }

//     // Clean up local files
//     _ = DeleteLocalFile(localInput)
//     _ = DeleteLocalFile(outputDir)

//     fmt.Println("[DEBUG] ChunkAndUpload completed successfully")
//     return nil
// }

// // ChunkVideo splits the input video into 10-second segments in the outputDir.
// func ChunkVideo(inputPath, outputDir string) error {
//     segmentDuration := 10 // seconds

//     if err := os.MkdirAll(outputDir, 0755); err != nil {
//         return fmt.Errorf("failed to create output directory: %w", err)
//     }

//     outputPattern := filepath.Join(outputDir, "chunk_%03d.mp4")

//     cmd := exec.Command(
//         "ffmpeg",
//         "-i", inputPath,
//         "-c", "copy",
//         "-map", "0",
//         "-f", "segment",
//         "-segment_time", fmt.Sprintf("%d", segmentDuration),
//         "-reset_timestamps", "1",
//         outputPattern,
//     )

//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr

//     if err := cmd.Run(); err != nil {
//         return fmt.Errorf("ffmpeg failed: %w", err)
//     }

//     return nil
// }

// // DeleteLocalFile deletes a file or directory at the given path.
// func DeleteLocalFile(path string) error {
//     return os.RemoveAll(path)
// }

package service

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"time"
    "log"

	"VideoChuncker/internal/machineryUtil"
	"VideoChuncker/internal/s3util"

	"github.com/RichardKnop/machinery/v2"
)

// ChunkAndUpload streams MPEG-TS segments directly from the source URL into R2
// and then sends the list of TS keys to ConvertChunks.
func ChunkAndUpload(server *machinery.Server, bucket, objectKey string) error {
    fmt.Printf("[INFO] Starting ChunkAndUpload for %s/%s\n", bucket, objectKey)

    // 1) Generate a presigned GET URL for the source MP4
    fmt.Println("[INFO] Presigning GET URL...")
    srcURL, err := s3util.GeneratePresignedGetURL(bucket, objectKey, 1*time.Hour)
    if err != nil {
        return fmt.Errorf("presign GET URL: %w", err)
    }
    fmt.Printf("[INFO] Got source URL: %s\n", srcURL)
    


    // 2) Probe total duration
    fmt.Println("[INFO] Probing duration with ffprobe...")
    durationSec, err := ProbeDurationSeconds(srcURL)
    if err != nil {
        return fmt.Errorf("probe duration: %w", err)
    }
    fmt.Printf("[INFO] Total duration: %.2f seconds\n", durationSec)

    // 3) Compute chunk count
    const segmentSeconds = 10
    numChunks := int((durationSec + float64(segmentSeconds) - 1) / float64(segmentSeconds))
    fmt.Printf("[INFO] Will create %d chunks of %d seconds each\n", numChunks, segmentSeconds)

    base := strings.TrimSuffix(objectKey, filepath.Ext(objectKey))
    var tsKeys []string

    // 4) Loop: for each segment, stream TS to R2 via PUT
    for i := 0; i < numChunks; i++ {
        start := i * segmentSeconds
        key := fmt.Sprintf("%s/chunk-%03d.ts", base, i)
        fmt.Printf("[INFO] Presigning PUT URL for chunk %d (%s)...\n", i, key)
        putURL, err := s3util.GeneratePresignedPutURL(bucket, key, 1*time.Hour)
        if err != nil {
            return fmt.Errorf("presign PUT URL for chunk %d: %w", i, err)
        }
        fmt.Printf("[INFO] Uploading chunk #%d: start=%d → PUT %s\n", i, start, putURL)

        if err := FFmpegUploadChunkAsTS(srcURL, start, segmentSeconds, putURL); err != nil {
            fmt.Printf("[ERROR] Chunk #%d failed: %v\n", i, err)
            return fmt.Errorf("chunk %d upload failed: %w", i, err)
        }
        fmt.Printf("[INFO] Chunk #%d uploaded successfully as %s\n", i, key)
        tsKeys = append(tsKeys, key)
    }


    // send the first chunk to Thumbnail Creator
    // fmt.Println("[INFO] Preparing to send first chunk to Thumbnail Creator...")
    // if server != nil && len(tsKeys) > 0 {
    //     fmt.Printf("[INFO] Sending first chunk to Thumbnail Creator: %s\n", tsKeys[0])
    //     if err := machineryUtil.SendTaskToThumbnailCreator(server, bucket, tsKeys[0]); err != nil {
    //         fmt.Printf("[WARN] Failed to send thumbnail task: %v\n", err)
    //     }
    //     fmt.Println("[INFO] Thumbnail task enqueued successfully")
    // }


    // 5) Send list of TS keys to ConvertChunks
    fmt.Println("[INFO] Sending TS keys to ConvertChunks task...")
    payload, _ := json.Marshal(tsKeys)
    if server != nil {
        if err := machineryUtil.SendTaskToConvertor(server, bucket, string(payload)); err != nil {
            return fmt.Errorf("send convert task: %w", err)
        }
        fmt.Println("[INFO] ConvertChunks task enqueued successfully")
    }

    fmt.Println("[INFO] ChunkAndUpload completed successfully")




    err2 := s3util.DeleteVideoFromR2(bucket, objectKey)
    if err != nil {
        log.Printf("[WARN] Failed to delete original video %s: %v", objectKey, err2)
    } else {
        log.Printf("[INFO] Deleted original video: %s", objectKey)
    }



    
    return nil
}

// ProbeDurationSeconds returns total duration of the media at inputURL.
func ProbeDurationSeconds(inputURL string) (float64, error) {
    fmt.Printf("[INFO] Start Probing for %s\n", inputURL)
    out, err := exec.Command(
        "ffprobe", "-v", "error",
        "-show_entries", "format=duration",
        "-of", "csv=p=0", inputURL,
    ).Output()
    if err != nil {
        return 0, fmt.Errorf("ffprobe failed: %w", err)
    }
    s := strings.TrimSpace(string(out))
    f, err := strconv.ParseFloat(s, 64)
    if err != nil {
        return 0, fmt.Errorf("parse duration %q: %w", s, err)
    }
    return f, nil
}


// FFmpegUploadChunkAsTS streams a single TS segment to putURL via HTTP PUT.
func FFmpegUploadChunkAsTS(inputURL string, startSec, durationSec int, putURL string) error {
    fmt.Printf("[INFO] Starting ffmpeg for segment start=%d, duration=%d\n", startSec, durationSec)
    args := []string{
        "-ss", strconv.Itoa(startSec),
        "-i", inputURL,
        "-t", strconv.Itoa(durationSec),
        "-reset_timestamps", "1",
        "-c", "copy",
        "-map", "0",
        "-f", "mpegts",
        "pipe:1",
    }
    ff := exec.Command("ffmpeg", args...)
    ff.Stderr = os.Stderr
    pipe, err := ff.StdoutPipe()
    if err != nil {
        return fmt.Errorf("ffmpeg stdout pipe: %w", err)
    }
    if err := ff.Start(); err != nil {
        return fmt.Errorf("start ffmpeg: %w", err)
    }

    // curl PUT
    fmt.Println("[INFO] Streaming TS via curl PUT...")
    curl := exec.Command("curl", "-s", "-X", "PUT",
        "-H", "Content-Type: video/mp2t",
        "--data-binary", "@-", putURL,
    )
    // curl := exec.Command("curl", "-s", "-X", "PUT",
    //     "--data-binary", "@-", putURL,
    // )
    curl.Stdin = pipe
    curl.Stderr = os.Stderr
    if err := curl.Run(); err != nil {
        ff.Process.Kill()
        return fmt.Errorf("curl PUT failed: %w", err)
    }

    if err := ff.Wait(); err != nil {
        return fmt.Errorf("ffmpeg exit error: %w", err)
    }
    fmt.Println("[INFO] ffmpeg segment pipeline finished")
    return nil
}








// func FFmpegUploadChunkAsTS(inputURL string, startSec, durationSec int, putURL string) error {
//     fmt.Printf("[INFO] Starting ffmpeg for segment start=%d, duration=%d\n", startSec, durationSec)

//     args := []string{
//         "-re",
//         "-ss", strconv.Itoa(startSec),
//         "-i", inputURL,
//         "-t", strconv.Itoa(durationSec),
//         "-reset_timestamps", "1",
//         "-c", "copy",
//         "-map", "0",
//         "-f", "mpegts",
//         "-flush_packets", "1",
//         "pipe:1",
//     }

//     ff := exec.Command("ffmpeg", args...)
//     ff.Stderr = os.Stderr

//     // Get ffmpeg stdout pipe
//     pipe, err := ff.StdoutPipe()
//     if err != nil {
//         return fmt.Errorf("ffmpeg stdout pipe: %w", err)
//     }

//     // Start ffmpeg process
//     if err := ff.Start(); err != nil {
//         return fmt.Errorf("start ffmpeg: %w", err)
//     }

//     // Prepare curl PUT command
//     fmt.Println("[INFO] Streaming TS via curl PUT...")
//     curl := exec.Command("curl", "-sS", "-X", "PUT",
//         "-H", "Content-Type: video/mp2t",
//         "--data-binary", "@-", putURL,
//     )
//     curl.Stdin = pipe

//     // Capture curl stderr and output for debugging
//     var curlOut bytes.Buffer
//     curl.Stdout = &curlOut
//     curl.Stderr = &curlOut

//     if err := curl.Start(); err != nil {
//         ff.Process.Kill()
//         return fmt.Errorf("start curl PUT: %w", err)
//     }

//     // Wait for curl to finish
//     if err := curl.Wait(); err != nil {
//         ff.Process.Kill()
//         return fmt.Errorf("curl PUT failed: %w\n[DEBUG] curl output: %s", err, curlOut.String())
//     }

//     // Wait for ffmpeg to finish
//     if err := ff.Wait(); err != nil {
//         return fmt.Errorf("ffmpeg exit error: %w", err)
//     }

//     fmt.Println("[INFO] ffmpeg segment pipeline finished successfully")
//     return nil
// }
