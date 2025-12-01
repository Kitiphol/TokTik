// package machineryutil

// import (
//     "github.com/RichardKnop/machinery/v2/tasks"
// )

// func SendChunkVideoTask(input string) error {
    
//     server, err := CreateMachineryServer()
//     if err != nil {
//         return err
//     }
//     signature := &tasks.Signature{
//         Name: "ChunkVideo",
//         Args: []tasks.Arg{
//             {
//                 Type:  "string",
//                 Value: input,
//             },
//         },
//     }
//     _, err = server.SendTask(signature)
//     return err
// }

package machineryutil

import (
    "encoding/json"
    "github.com/RichardKnop/machinery/v2/tasks"
    "fmt"
)

type ChunkVideoPayload struct {
    Bucket    string `json:"bucket"`
    ObjectKey string `json:"objectKey"`
}

func SendChunkVideoTask(bucket, objectKey string) error {
    server, err := CreateMachineryServer()
    if err != nil {
        return err
    }


    // bucket = bucket name
    // objectKey = path to the video file in the bucket
    payload := ChunkVideoPayload{
        Bucket:    bucket, //toktikp2-video
        ObjectKey: objectKey, //"userID/videoID/filename.mp4"
    }

    // Validate JSON formatting
    payloadBytes, err := json.Marshal(payload)
    if err != nil {
        return err
    }


    // Optional: Validate that it can be unmarshaled back (format check)
    var check ChunkVideoPayload
    if err := json.Unmarshal(payloadBytes, &check); err != nil {
        return fmt.Errorf("payload format invalid: %w", err)
    }


    signature := &tasks.Signature{
        Name: "ChunkVideo",
        Args: []tasks.Arg{
            {
                Type:  "string",
                Value: string(payloadBytes),
            },
        },
    }
    _, err = server.SendTask(signature)

    if err != nil {
        return fmt.Errorf("failed to send task: %w", err)
    }
    fmt.Printf("[DEBUG] Sent ChunkVideo task: %s\n", string(payloadBytes))
    return err
}


func SendCreateThumbnailTask(bucket, objectKey string) error {
    server, err := CreateMachineryServer()
    if err != nil {
        return fmt.Errorf("create server: %w", err)
    }

    signature := &tasks.Signature{
        Name: "CreateThumbnail",
        Args: []tasks.Arg{
            {
                Type:  "string",
                Value: bucket,
            },
            {
                Type:  "string",
                Value: objectKey,
            },
        },
    }

    asyncResult, err := server.SendTask(signature)
    if err != nil {
        return fmt.Errorf("failed to send task: %w", err)
    }

    fmt.Printf("[DEBUG] Sent CreateThumbnail task for bucket=%s, key=%s. Task UUID: %s\n",
        bucket, objectKey, asyncResult.Signature.UUID)

    return nil
}