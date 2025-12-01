











// package service

// import (
//     "VideoConvertor/internal/s3util"
//     "github.com/RichardKnop/machinery/v2"
//     "context"
//     "encoding/json"
//     "fmt"
//     "os"
//     "os/exec"
//     "path/filepath"
//     "sort"
//     "strconv"
//     "strings"
//     "io"
//     "log"
// )



// func ConvertChunks(ctx context.Context, server *machinery.Server, bucket string, chunkKeysJSON string) error {
//     // Parse chunk keys
//     var keys []string
//     if err := json.Unmarshal([]byte(chunkKeysJSON), &keys); err != nil {
//         return fmt.Errorf("invalid JSON: %w", err)
//     }
//     if len(keys) == 0 {
//         return fmt.Errorf("no chunks to convert")
//     }

//     // Prepare local HLS temp directory
//     parts := strings.SplitN(keys[0], "/", 3)
//     base := fmt.Sprintf("%s_%s", parts[0], parts[1])
//     outputDir := filepath.Join(os.TempDir(), base+"_hls")
//     os.RemoveAll(outputDir)
//     if err := os.MkdirAll(outputDir, 0755); err != nil {
//         return fmt.Errorf("mkdir HLS dir: %w", err)
//     }
//     defer os.RemoveAll(outputDir)

//     // 1. Download each TS chunk and build concat list
//     var concatLines []string
//     var originalChunkKeys []string
//     for i, key := range keys {
//         localTS := filepath.Join(outputDir, fmt.Sprintf("input-%03d.ts", i))
//         fmt.Printf("[INFO] Downloading chunk %d from %s to %s\n", i, key, localTS)
//         if err := s3util.DownloadFromR2(bucket, key, localTS); err != nil {
//             return fmt.Errorf("download chunk %s: %w", key, err)
//         }
//         concatLines = append(concatLines, fmt.Sprintf("file '%s'\n", localTS))
//         originalChunkKeys = append(originalChunkKeys, key)
//     }

//     // 2. Write concat.txt
//     concatPath := filepath.Join(outputDir, "concat.txt")
//     fmt.Printf("[INFO] Writing concat file %s\n", concatPath)
//     if err := os.WriteFile(concatPath, []byte(strings.Join(concatLines, "")), 0644); err != nil {
//         return fmt.Errorf("write concat.txt: %w", err)
//     }

//     // 3. Run FFmpeg to generate initial HLS playlist and segments
//     fmt.Println("[INFO] Running ffmpeg HLS generation...")
//     playlistPath := filepath.Join(outputDir, "playlist.m3u8")
//     cmd := exec.Command("ffmpeg", "-y",
//         "-f", "concat", "-safe", "0", "-i", concatPath,
//         "-vf", "scale='if(gt(iw,ih),1080,-2)':'if(gt(iw,ih),-2,1920)',setsar=1",
//         "-c:v", "libx264", "-preset", "veryslow", "-crf", "28",
//         "-profile:v", "baseline", "-level", "3.0",
//         "-maxrate", "800k", "-bufsize", "1200k",
//         "-c:a", "aac", "-b:a", "96k", "-ac", "2",
//         "-f", "hls", "-hls_time", "10", "-hls_list_size", "0",
//         "-hls_segment_filename", filepath.Join(outputDir, "%03d.ts"),
//         playlistPath,
//     )
//     cmd.Stdout = os.Stdout
//     cmd.Stderr = os.Stderr
//     if err := cmd.Run(); err != nil {
//         return fmt.Errorf("ffmpeg HLS failed: %w", err)
//     }
//     fmt.Println("[INFO] ffmpeg initial HLS generation complete")

//     // 4. Rewrite playlist.m3u8 with accurate durations
//     f, err := os.Create(playlistPath)
//     if err != nil {
//         return fmt.Errorf("failed to create playlist file: %w", err)
//     }



//     // write header
//     f.WriteString("#EXTM3U\n")
//     if err := writeLog(f, "header #EXTM3U", "#EXTM3U\n"); err != nil {
//        return fmt.Errorf("write header: %w", err)
//     }

//     f.WriteString("#EXT-X-VERSION:3\n")

//     if err := writeLog(f, "header VERSION", "#EXT-X-VERSION:3\n"); err != nil {
//        return fmt.Errorf("write header: %w", err)
//     }


//     f.WriteString("#EXT-X-TARGETDURATION:10\n")

    
//     f.WriteString("#EXT-X-MEDIA-SEQUENCE:0\n")

//     // find each generated .ts (excluding input-*.ts)
//     entries, _ := os.ReadDir(outputDir)
//     var segs []string
//     for _, e := range entries {
//         name := e.Name()
//         if !e.IsDir() && strings.HasSuffix(name, ".ts") && !strings.HasPrefix(name, "input-") {
//             segs = append(segs, name)
//         }
//     }
//     sort.Strings(segs)

//     // loop and probe durations
//     for _, seg := range segs {
//         path := filepath.Join(outputDir, seg)
//         out, err := exec.Command("ffprobe", "-v", "error",
//             "-show_entries", "format=duration",
//             "-of", "default=noprint_wrappers=1:nokey=1", path,
//         ).Output()
//         duration := 10.0
//         if err == nil {
//             if d, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64); err == nil {
//                 duration = d
//             }
//         }
//         f.WriteString(fmt.Sprintf("#EXTINF:%.3f,\n", duration))
//         f.WriteString(seg + "\n")
//     }
//     f.WriteString("#EXT-X-ENDLIST\n")
//     if err := writeLog(f, "footer ENDLIST", "#EXT-X-ENDLIST\n"); err != nil {
//       return fmt.Errorf("write footer: %w", err)
//     }

//     // ensure it's flushed and closed before uploading
//     if err := f.Sync(); err != nil {
//         return fmt.Errorf("sync playlist: %w", err)
//     }
//     if err := f.Close(); err != nil {
//         return fmt.Errorf("close playlist: %w", err)
//     }
//     fmt.Println("[INFO] Rewritten playlist.m3u8 closed, size and contents finalized")

//     // 5. Now re-scan the directory and upload everything (except input-*.ts)
//     files, err := os.ReadDir(outputDir)
//     if err != nil {
//         return fmt.Errorf("read HLS dir: %w", err)
//     }
//     prefix := strings.Join(parts[:2], "/")
//     for _, f := range files {
//         name := f.Name()
//         if f.IsDir() || strings.HasPrefix(name, "input-") {
//             continue
//         }
//         localFile := filepath.Join(outputDir, name)
//         key := fmt.Sprintf("%s/%s", prefix, name)
//         fmt.Printf("[INFO] Uploading HLS file %s to %s/%s\n", name, bucket, key)
//         if err := s3util.UploadFile(bucket, key, localFile); err != nil {
//             return fmt.Errorf("upload HLS file %s: %w", name, err)
//         }
//     }

//     // 6. (optional) enqueue thumbnail task
//     // firstSegKey := fmt.Sprintf("%s/000.ts", prefix)
//     // if server != nil {
//     //     fmt.Println("[INFO] Enqueueing thumbnail task")
//     //     log.Printf("[DEBUG] Sending CreateThumbnail task for %s/%s\n", bucket, firstSegKey)
//     //     _ = SendTaskToThumbnailCreator(server, bucket, firstSegKey)
//     // }

//     // 7. Delete original TS chunks
//     for _, chunkKey := range originalChunkKeys {
//         fmt.Printf("[INFO] Deleting original chunk %s from R2\n", chunkKey)
//         if err := s3util.DeleteVideoFromR2(bucket, chunkKey); err != nil {
//             fmt.Printf("[WARN] Failed to delete chunk %s: %v\n", chunkKey, err)
//         }
//     }

//     fmt.Println("[INFO] ConvertChunks completed successfully")
//     return nil
// }


// func writeLog(w io.WriteSeeker, desc, data string) error {
//     n, err := w.Write([]byte(data))
//     log.Printf("[DEBUG] write %q (%s): %d bytes, err=%v", data, desc, n, err)
//     return err
// }




package service

import (
    "context"
    "encoding/json"
    "fmt"
    "io"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "sort"
    "strconv"
    "strings"

    "github.com/RichardKnop/machinery/v2"
    "VideoConvertor/internal/s3util"
)

// writeLog writes data to w, logging its description, byte‑count, and any error.
func writeLog(w io.Writer, desc, data string) error {
    n, err := w.Write([]byte(data))
    log.Printf("[DEBUG] write %q (%s): %d bytes, err=%v", data, desc, n, err)
    return err
}

func ConvertChunks(ctx context.Context, server *machinery.Server, bucket string, chunkKeysJSON string) error {
    // Parse chunk keys
    var keys []string
    if err := json.Unmarshal([]byte(chunkKeysJSON), &keys); err != nil {
        return fmt.Errorf("invalid JSON: %w", err)
    }
    if len(keys) == 0 {
        return fmt.Errorf("no chunks to convert")
    }

    // Prepare local HLS temp directory
    parts := strings.SplitN(keys[0], "/", 3)
    base := fmt.Sprintf("%s_%s", parts[0], parts[1])
    outputDir := filepath.Join(os.TempDir(), base+"_hls")
    os.RemoveAll(outputDir)
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        return fmt.Errorf("mkdir HLS dir: %w", err)
    }
    defer os.RemoveAll(outputDir)

    // 1. Download each TS chunk and build concat list
    var concatLines []string
    var originalChunkKeys []string
    for i, key := range keys {
        localTS := filepath.Join(outputDir, fmt.Sprintf("input-%03d.ts", i))
        log.Printf("[INFO] Downloading chunk %d from %s to %s", i, key, localTS)
        if err := s3util.DownloadFromR2(bucket, key, localTS); err != nil {
            return fmt.Errorf("download chunk %s: %w", key, err)
        }
        concatLines = append(concatLines, fmt.Sprintf("file '%s'\n", localTS))
        originalChunkKeys = append(originalChunkKeys, key)
    }

    // 2. Write concat.txt
    concatPath := filepath.Join(outputDir, "concat.txt")
    log.Printf("[INFO] Writing concat file %s", concatPath)
    if err := os.WriteFile(concatPath, []byte(strings.Join(concatLines, "")), 0644); err != nil {
        return fmt.Errorf("write concat.txt: %w", err)
    }

    // 3. Run FFmpeg to generate initial HLS playlist and segments
    log.Println("[INFO] Running ffmpeg HLS generation...")
    playlistPath := filepath.Join(outputDir, "playlist.m3u8")
    cmd := exec.Command("ffmpeg", "-y",
        "-f", "concat", "-safe", "0", "-i", concatPath,
        "-vf", "scale='if(gt(iw,ih),1080,-2)':'if(gt(iw,ih),-2,1920)',setsar=1",
        "-c:v", "libx264", "-preset", "veryslow", "-crf", "28",
        "-profile:v", "baseline", "-level", "3.0",
        "-maxrate", "800k", "-bufsize", "1200k",
        "-c:a", "aac", "-b:a", "96k", "-ac", "2",
        "-f", "hls", "-hls_time", "10", "-hls_list_size", "0",
        "-hls_segment_filename", filepath.Join(outputDir, "%03d.ts"),
        playlistPath,
    )
    cmd.Stdout = os.Stdout
    cmd.Stderr = os.Stderr
    if err := cmd.Run(); err != nil {
        return fmt.Errorf("ffmpeg HLS failed: %w", err)
    }
    log.Println("[INFO] ffmpeg initial HLS generation complete")

    // 4. Rewrite playlist.m3u8 with accurate durations, logging every write
    f, err := os.Create(playlistPath)
    if err != nil {
        return fmt.Errorf("failed to create playlist file: %w", err)
    }

    // headers
    if err := writeLog(f, "header #EXTM3U", "#EXTM3U\n"); err != nil {
        return fmt.Errorf("write header: %w", err)
    }
    if err := writeLog(f, "header VERSION", "#EXT-X-VERSION:3\n"); err != nil {
        return fmt.Errorf("write header: %w", err)
    }
    if err := writeLog(f, "header TARGETDURATION", "#EXT-X-TARGETDURATION:10\n"); err != nil {
        return fmt.Errorf("write header: %w", err)
    }
    if err := writeLog(f, "header MEDIA-SEQUENCE", "#EXT-X-MEDIA-SEQUENCE:0\n"); err != nil {
        return fmt.Errorf("write header: %w", err)
    }

    // find segments
    entries, _ := os.ReadDir(outputDir)
    var segs []string
    for _, e := range entries {
        name := e.Name()
        if !e.IsDir() && strings.HasSuffix(name, ".ts") && !strings.HasPrefix(name, "input-") {
            segs = append(segs, name)
        }
    }
    sort.Strings(segs)

    // EXTINF + filename per segment
    for _, seg := range segs {
        path := filepath.Join(outputDir, seg)
        out, err := exec.Command("ffprobe", "-v", "error",
            "-show_entries", "format=duration",
            "-of", "default=noprint_wrappers=1:nokey=1", path,
        ).Output()
        duration := 10.0
        if err == nil {
            if d, err := strconv.ParseFloat(strings.TrimSpace(string(out)), 64); err == nil {
                duration = d
            }
        }

        infoLine := fmt.Sprintf("#EXTINF:%.3f,\n", duration)
        if err := writeLog(f, "EXTINF "+seg, infoLine); err != nil {
            return fmt.Errorf("write EXTINF: %w", err)
        }
        if err := writeLog(f, "segment name "+seg, seg+"\n"); err != nil {
            return fmt.Errorf("write segment name: %w", err)
        }
    }

    // footer
    if err := writeLog(f, "footer ENDLIST", "#EXT-X-ENDLIST\n"); err != nil {
        return fmt.Errorf("write footer: %w", err)
    }

    // flush & close
    if err := f.Sync(); err != nil {
        return fmt.Errorf("sync playlist: %w", err)
    }
    if err := f.Close(); err != nil {
        return fmt.Errorf("close playlist: %w", err)
    }
    log.Printf("[DEBUG] Finished writing playlist, path=%s", playlistPath)

    // check if writing was successful
    info, err := os.Stat(playlistPath)
    if err != nil {
        return fmt.Errorf("stat playlist: %w", err)
    }
    log.Printf("[DEBUG] local playlist.m3u8 size: %d bytes", info.Size())
    if info.Size() == 0 {
        return fmt.Errorf("local playlist is empty!")
    }



    // 5. Re-scan directory and upload everything except input-*.ts
    // files, err := os.ReadDir(outputDir)
    // if err != nil {
    //     return fmt.Errorf("read HLS dir: %w", err)
    // }
    // prefix := strings.Join(parts[:2], "/")
    // for _, entry := range files {
    //     name := entry.Name()
    //     if entry.IsDir() || strings.HasPrefix(name, "input-") {
    //         continue
    //     }
    //     localFile := filepath.Join(outputDir, name)
    //     key := fmt.Sprintf("%s/%s", prefix, name)
    //     log.Printf("[INFO] Uploading HLS file %s to %s/%s", name, bucket, key)
    //     if err := s3util.UploadFile(bucket, key, localFile); err != nil {
    //         return fmt.Errorf("upload HLS file %s: %w", name, err)
    //     }
    // }


    files, err := os.ReadDir(outputDir)
    if err != nil {
        return fmt.Errorf("read HLS dir: %w", err)
    }
    prefix := strings.Join(parts[:2], "/")
    for _, entry := range files {
        name := entry.Name()
        if entry.IsDir() {
            continue
        }

        if strings.HasPrefix(name, "input-") {
            continue
        }
        
        // skip anything not .ts or .m3u8
        if !(strings.HasSuffix(name, ".ts") || name == "playlist.m3u8") {
            continue
        }


        localFile := filepath.Join(outputDir, name)
        key := fmt.Sprintf("%s/%s", prefix, name)
        log.Printf("[INFO] Uploading HLS file %s to %s/%s", name, bucket, key)
        if err := s3util.UploadFile(bucket, key, localFile); err != nil {
            return fmt.Errorf("upload HLS file %s: %w", name, err)
        }
    }






    // 6. (optional) enqueue thumbnail task
    // ...

    // 7. Delete original TS chunks
    for _, chunkKey := range originalChunkKeys {
        log.Printf("[INFO] Deleting original chunk %s from R2", chunkKey)
        if err := s3util.DeleteVideoFromR2(bucket, chunkKey); err != nil {
            log.Printf("[WARN] Failed to delete chunk %s: %v", chunkKey, err)
        }
    }

    log.Println("[INFO] ConvertChunks completed successfully")
    return nil
}
