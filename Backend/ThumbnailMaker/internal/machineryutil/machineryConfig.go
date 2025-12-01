package machineryutil

import (
	"fmt"
    "ThumbnailCreator/internal/env_config"
	"net/url"
	_ "strings"

	// Machinery core + config
	"github.com/RichardKnop/machinery/v2"
	"github.com/RichardKnop/machinery/v2/config"

	// Interfaces
	backendsiface "github.com/RichardKnop/machinery/v2/backends/iface"
	brokersiface "github.com/RichardKnop/machinery/v2/brokers/iface"
	locksiface "github.com/RichardKnop/machinery/v2/locks/iface"

	// RabbitMQ broker
	"github.com/RichardKnop/machinery/v2/brokers/amqp"

	// Redis backend
	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
)

// GetMachineryConfig returns a v2 Config that uses RabbitMQ for broker
// and Redis for the result backend.
func GetMachineryConfig() *config.Config {
	var err error
	cfg, err := env_config.LoadConfig()
	if err != nil {
		panic(fmt.Sprintf("Failed to load environment config: %v", err))
	}

	return &config.Config{
		// ——— BROKER (RabbitMQ) ———————————————————————————————————————————————
		Broker: cfg.RabbitBrokerURI,

		// ——— RESULT BACKEND (Redis) ——————————————————————————————————————————
		// Store task state & return values in Redis db 0
		ResultBackend:   cfg.RabbitBackendURI,
		ResultsExpireIn: 3600, // expire stored results after 1h

		// ——— DEFAULT QUEUE ————————————————————————————————————————————————
		// Name of the RabbitMQ queue that workers will consume from
		DefaultQueue: "machinery_tasks_queue",

		// ——— REDIS CONFIG ————————————————————————————————————————————————
		// Necessary for both lock (if used) and result‐backend connection pooling.
		Redis: &config.RedisConfig{
			MaxIdle:                10,
			MaxActive:              50,
			IdleTimeout:            300,
			Wait:                   true,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  500,
			DelayedTasksPollPeriod: 20,
			DelayedTasksKey:        "machinery:delayed_tasks",
			MasterName:             "",
		},

		// ——— AMQP (RabbitMQ) SETTINGS ————————————————————————————————————————
		AMQP: &config.AMQPConfig{
			Exchange:      "machinery_exchange",
			ExchangeType:  "direct",
			BindingKey:    "task.#",
			PrefetchCount: 1,
			AutoDelete:    false,
			DelayedQueue:  "",
		},

		// ——— LOCK BACKEND ————————————————————————————————————————————————
		// We’re not using a separate lock, so leave empty.
		Lock: "",

		// ——— EVERYTHING ELSE ——————————————————————————————————————————————
		MultipleBrokerSeparator: "|",
	}
}

// buildBroker returns a RabbitMQ (AMQP) broker implementation.
func buildBroker(cfg *config.Config) (brokersiface.Broker, error) {
	return amqp.New(cfg), nil
}

// buildBackend returns a Redis‐based result backend.
func buildBackend(cfg *config.Config) (backendsiface.Backend, error) {
	u, err := url.Parse(cfg.ResultBackend)
	if err != nil {
		return nil, fmt.Errorf("invalid ResultBackend URI %q: %w", cfg.ResultBackend, err)
	}
	if u.Scheme != "redis" {
		return nil, fmt.Errorf("ResultBackend URI must use redis:// scheme, got %q", u.Scheme)
	}
	host := u.Host
	var password string
	if u.User != nil {
		if pw, ok := u.User.Password(); ok {
			password = pw
		}
	}
	var db = 0
	socketPath := ""

	return redisbackend.New(cfg, host, password, socketPath, db), nil
}

// buildLock returns nil because we’re not using a separate distributed lock.
func buildLock() (locksiface.Lock, error) {
	return nil, nil
}

// CreateMachineryServer wires everything together and returns *machinery.Server.
func CreateMachineryServer() (*machinery.Server, error) {
	// 1) Load config
	cfg := GetMachineryConfig()

	// 2) Instantiate the broker (RabbitMQ)
	broker, err := buildBroker(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create AMQP broker: %w", err)
	}

	// 3) Instantiate the Redis backend
	backend, err := buildBackend(cfg)
	if err != nil {
		return nil, fmt.Errorf("cannot create Redis backend: %w", err)
	}

	// 4) No lock backend in this example
	lock, err := buildLock()
	if err != nil {
		return nil, fmt.Errorf("cannot create lock: %w", err)
	}

	server := machinery.NewServer(cfg, broker, backend, lock)
	return server, nil
}