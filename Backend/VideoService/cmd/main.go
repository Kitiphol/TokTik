package main

import (
    "VideoService/internal/database"
	"VideoService/internal/config"
	"VideoService/internal/routes"
    "VideoService/internal/machineryUtil"
    "VideoService/internal/publisher"
    

    "fmt"
    "os"


)
func main() {
    database.InitDB()

    publisher.InitCluster()

 
    _, err := machineryutil.CreateMachineryServer()
    if err != nil {
        fmt.Printf("Failed to create machinery server: %v\n", err)
        os.Exit(1)
    }

 
    serverPort := config.Load().Port
    r := routes.Setup()
    r.Run(":" + serverPort)
}
