package service

import (
    "github.com/RichardKnop/machinery/v2"
    "github.com/RichardKnop/machinery/v2/tasks"
)

// Sends a task to the thumbnail creator with the full segment key (not just folder/filename)
func SendTaskToThumbnailCreator(server *machinery.Server, bucket, segmentKey string) error {
    signature := &tasks.Signature{
        Name: "CreateThumbnail",
        Args: []tasks.Arg{
            {Type: "string", Value: bucket},
            {Type: "string", Value: segmentKey},
        },
    }
    _, err := server.SendTask(signature)
    return err
}