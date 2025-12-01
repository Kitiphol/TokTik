// package machineryUtil

// import (
//     "github.com/RichardKnop/machinery/v2"
//     "github.com/RichardKnop/machinery/v2/tasks"
// )

// func SendTaskToConvertor(server *machinery.Server, taskName string, args ...tasks.Arg) error {
//     signature := &tasks.Signature{
//         Name: taskName,
//         Args: args,
//     }
//     _, err := server.SendTask(signature)
//     return err
// }

package machineryUtil

import (
    "github.com/RichardKnop/machinery/v2"
    "github.com/RichardKnop/machinery/v2/tasks"
)

// Sends a ConvertChunks task to the VideoConvertor with bucket and chunkKeysJSON
func SendTaskToConvertor(server *machinery.Server, bucket string, chunkKeysJSON string) error {
    signature := &tasks.Signature{
        Name: "ConvertChunks",
        Args: []tasks.Arg{
            {Type: "string", Value: bucket},
            {Type: "string", Value: chunkKeysJSON},
        },
    }
    _, err := server.SendTask(signature)
    return err
}


func SendTaskToThumbnailCreator(server *machinery.Server, bucket string, segmentKey string) error {
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
