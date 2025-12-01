package main

import (
    "fmt"
    "os"
    "VideoConvertor/internal/machineryutil"
)

func main() {
    server, err := machineryutil.CreateMachineryServer()
    if err != nil {
        fmt.Printf("Failed to create machinery server: %v\n", err)
        os.Exit(1)
    }

    // Use the StartWorker helper to register and launch the worker
    if err := machineryutil.StartWorker(server); err != nil {
        fmt.Printf("Worker failed: %v\n", err)
        os.Exit(1)
    }
}