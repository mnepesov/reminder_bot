package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reminder_bot/config"
	"reminder_bot/internal/parsing/delivery"
	"reminder_bot/pkg/rabbitmq"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("No .env file found")
	}
}

func main() {

	configPath := config.GetConfigPath(os.Getenv("IsDebug"))
	cfg, err := config.GetConfig(configPath)
	if err != nil {
		log.Fatalf("Loading config: %v", err)
	}

	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		log.Fatalf("connecting to rabbitmq: %v", err)
	}

	consumers, err := delivery.NewParsingConsumer(amqpConn)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	forever := make(chan struct{})
	go func() {
		err = consumers.StartParseConsumer(ctx, cfg.RabbitMQ.WorkerPoolSize, cfg.Queue.ParsingCommandTextParse)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	}()

	fmt.Println("Parsing Service is started")

	select {
	case <-forever:
		break
	case <-ctx.Done():
		fmt.Println("Выхожу")
		break
	}

}
