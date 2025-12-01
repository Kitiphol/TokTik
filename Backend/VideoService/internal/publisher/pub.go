package publisher

import (
    "context"
    "encoding/json"
    "os"
    "strings"

    "github.com/redis/go-redis/v9"
	"fmt"
)


var RedisClient *redis.ClusterClient


func InitCluster() {
    addrsCSV := os.Getenv("REDIS_ADDRS")     
    password := os.Getenv("REDIS_PASSWORD")   
    addrs := strings.Split(addrsCSV, ",")

	fmt.Print("Connecting to Redis cluster at: ", addrs, "\n")

    RedisClient = redis.NewClusterClient(&redis.ClusterOptions{
        Addrs:    addrs,
        Password: password,
    })


	if err := RedisClient.Ping(context.Background()).Err(); err != nil {
        // could be connection refused, auth failure, DNS lookup, etc.
        fmt.Printf("❌ Redis cluster ping failed: %v\n", err)
        // you can choose to os.Exit(1) here, or retry, or just leave RedisClient nil
    } else {
        fmt.Println("✅ Connected to Redis cluster")
    }
}

// Publish marshals `msg` to JSON and PUBLISH-es it on `channel`.
func Publish(channel string, msg interface{}) error {
    b, err := json.Marshal(msg)
    if err != nil {
		fmt.Printf("Failed to marshal message: %v\n", err)
        return err

    }
    return RedisClient.Publish(context.Background(), channel, b).Err()
}
