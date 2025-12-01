package machineryutil

import (
    "github.com/RichardKnop/machinery/v2"
    "ThumbnailCreator/internal/service"
)

func StartWorker(server *machinery.Server) error {
    tasks := map[string]interface{}{
        "CreateThumbnail": service.CreateThumbnail, // This must match the task name sent by VideoChunker
    }
    if err := server.RegisterTasks(tasks); err != nil {
        return err
    }
    worker := server.NewWorker("thumbnail-convertor-worker", 1) 
    return worker.Launch()
}