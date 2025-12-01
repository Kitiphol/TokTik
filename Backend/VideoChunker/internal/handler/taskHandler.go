package handler

import (
    "context"
    "encoding/json"
    "fmt"
    "VideoChuncker/internal/service"
    "github.com/RichardKnop/machinery/v2"
)

// Payload struct
type ChunkVideoPayload struct {
    Bucket    string `json:"bucket"`
    ObjectKey string `json:"objectKey"`
}

// Task function to receive and process the message
// func ChunkVideo(ctx context.Context, payloadStr string) error {
//     fmt.Println("[DEBUG] Received payload:", payloadStr)
//     var payload ChunkVideoPayload
//     if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
//         return fmt.Errorf("failed to parse payload: %w", err)
//     }
//     fmt.Printf("[DEBUG] Processing: bucket=%s, objectKey=%s\n", payload.Bucket, payload.ObjectKey)
//     // Call your chunking logic
//     return service.ChunkAndUpload(nil, payload.Bucket, payload.ObjectKey)
// }



func ChunkVideoHandler(server *machinery.Server) func(ctx context.Context, payloadStr string) error {
    return func(ctx context.Context, payloadStr string) error {
        fmt.Printf("[DEBUG] [Chunker Worker] Received ChunkVideo task: payload=%s\n", payloadStr)
        var payload struct {
            Bucket    string `json:"bucket"`
            ObjectKey string `json:"objectKey"`
        }
        if err := json.Unmarshal([]byte(payloadStr), &payload); err != nil {
            fmt.Printf("[DEBUG] [Chunker Worker] Failed to parse payload: %v\n", err)
            return fmt.Errorf("failed to parse payload: %w", err)
        }
        fmt.Printf("[DEBUG] [Chunker Worker] Starting ChunkAndUpload: bucket=%s, objectKey=%s\n", payload.Bucket, payload.ObjectKey)
        err := service.ChunkAndUpload(server, payload.Bucket, payload.ObjectKey)
        if err != nil {
            fmt.Printf("[DEBUG] [Chunker Worker] ChunkAndUpload failed: %v\n", err)
        } else {
            fmt.Printf("[DEBUG] [Chunker Worker] ChunkAndUpload completed successfully\n")
        }
        return err
    }
}

// StartWorker registers tasks and starts the worker
func StartWorker(server *machinery.Server) error {
    tasks := map[string]interface{}{
        "ChunkVideo": ChunkVideoHandler(server),
    }
    if err := server.RegisterTasks(tasks); err != nil {
        return err
    }
    worker := server.NewWorker("video-chuncker-worker", 5)
    return worker.Launch()
}