package env_config

import (
	"fmt"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	RabbitBrokerURI  string `envconfig:"BROKER_URI" required:"true"`
	RabbitBackendURI string `envconfig:"BACKEND_URI" required:"true"`

	R2AccountID       string `envconfig:"R2_ACCOUNT_ID" required:"true"`
	R2AccessKeyID     string `envconfig:"R2_ACCESS_KEY_ID" required:"true"`
	R2SecretAccessKey string `envconfig:"R2_ACCESS_KEY_SECRET" required:"true"`

	// UploadBucket  string `envconfig:"UPLOAD_BUCKET" default:"video-uploads"`
	// ChunkBucket   string `envconfig:"CHUNK_BUCKET" default:"video-chunks"`
	// ConvertBucket string `envconfig:"CONVERT_BUCKET" default:"video-converted"`

	WorkerConcurrency int `envconfig:"WORKER_CONCURRENCY" default:"5"`
}

func LoadConfig() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, fmt.Errorf("failed to load env config: %w", err)
	}
	return &cfg, nil
}