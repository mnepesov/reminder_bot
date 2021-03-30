package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reminder_bot/config"
	"reminder_bot/internal/reminder/delivery"
	"reminder_bot/internal/reminder/repository"
	"reminder_bot/internal/reminder/usecase"
	"reminder_bot/pkg/postgres"
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

	postgresConn, err := postgres.NewPostgresDB(cfg.Postgres)
	if err != nil {
		log.Fatal(err)
	}
	repo := repository.NewReminderPostgres(postgresConn)

	useCase := usecase.NewReminderUseCase(repo)

	consumers, err := delivery.NewReminderConsumer(amqpConn, useCase)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	forever := make(chan struct{})
	go func() {
		err = consumers.StartAddReminderConsumer(ctx, cfg.RabbitMQ.WorkerPoolSize, cfg.Queue.ReminderCommandReminderAdd)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	}()

	go func() {
		err = consumers.StartGetRemindersConsumer(ctx, cfg.RabbitMQ.WorkerPoolSize, cfg.Queue.ReminderCommandRemindersGet)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	}()

	fmt.Println("Reminder Service is started")

	select {
	case <-forever:
		break
	case <-ctx.Done():
		fmt.Println("Выхожу")
		break
	}
}
