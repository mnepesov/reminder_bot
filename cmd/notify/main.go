package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reminder_bot/config"
	"reminder_bot/internal/notify"
	"reminder_bot/internal/notify/repository"
	"reminder_bot/internal/notify/usecase"
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
	
	repo := repository.NewNotifyPostgres(postgresConn)
	
	uc := usecase.NewNotifyUseCase(repo)
	
	srvc, err := notify.NewNotifyService(amqpConn, uc, cfg.Queue.BotEventNotifySend)
	if err != nil {
		log.Fatal(err)
	}
	
	ctx, cancel := context.WithCancel(context.Background())
	
	go func() {
		if err := srvc.Start(ctx); err != nil {
			cancel()
			fmt.Println(err)
		}
	}()
	
	forever := make(chan struct{})
	
	fmt.Println("Notify Service is started")
	
	select {
	case <-forever:
		break
	case <-ctx.Done():
		fmt.Println("Выхожу")
		break
	}
}
