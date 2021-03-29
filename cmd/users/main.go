package main

import (
	"context"
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"reminder_bot/config"
	"reminder_bot/internal/users/delivery"
	"reminder_bot/internal/users/repository"
	"reminder_bot/internal/users/usecase"
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

	postgresConn,err:=postgres.NewPostgresDB(cfg.Postgres)
	if err!=nil {
		log.Fatal(err)
	}
	repo := repository.NewUsersRepository(postgresConn)

	useCase := usecase.NewUserUseCase(repo)

	consumers, err := delivery.NewUserConsumer(amqpConn, useCase)
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithCancel(context.Background())

	forever := make(chan struct{})
	go func() {
		err = consumers.StartGetUserConsumer(ctx, cfg.RabbitMQ.WorkerPoolSize, cfg.Queue.UsersRepoCommandUsersGet)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	}()
	go func() {
		err = consumers.StartCreateUserConsumer(cfg.RabbitMQ.WorkerPoolSize, cfg.Queue.UsersRepoCommandUsersCreate)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	}()

	go func() {
		err = consumers.StartUpdateTimeZoneConsumer(cfg.RabbitMQ.WorkerPoolSize, cfg.Queue.UsersRepoCommandUsersUpdateTimezone)
		if err != nil {
			cancel()
			fmt.Println(err)
		}
	}()

	fmt.Println("User Service is started")

	select {
	case <-forever:
		break
	case <-ctx.Done():
		fmt.Println("Выхожу")
		break
	}
}
