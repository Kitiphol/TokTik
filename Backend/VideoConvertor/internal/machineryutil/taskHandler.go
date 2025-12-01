package machineryutil

import (
    "context"
    "github.com/RichardKnop/machinery/v2"
    "VideoConvertor/internal/service"
    "fmt"
)

func StartWorker(server *machinery.Server) error {
    tasks := map[string]interface{}{
        "ConvertChunks": func(bucket string, chunkKeysJSON string) error {
            fmt.Printf("[DEBUG] [Convertor Worker] Received ConvertChunks task: bucket=%s, chunkKeysJSON=%s\n", bucket, chunkKeysJSON)
            err := service.ConvertChunks(context.Background(), server, bucket, chunkKeysJSON)
            if err != nil {
                fmt.Printf("[DEBUG] [Convertor Worker] ConvertChunks task failed: %v\n", err)
            } else {
                fmt.Println("[DEBUG] [Convertor Worker] ConvertChunks task completed successfully")
            }
            return err
        },
    }
    if err := server.RegisterTasks(tasks); err != nil {
        return err
    }
    worker := server.NewWorker("video-convertor-worker", 5)
    return worker.Launch()
}