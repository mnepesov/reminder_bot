package main

import (
	"encoding/json"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/joho/godotenv"
	"github.com/streadway/amqp"
	"googlemaps.github.io/maps"
	"log"
	"os"
	"reminder_bot/config"
	cache2 "reminder_bot/internal/bot/pkg/cache"
	"reminder_bot/internal/bot/pkg/parsing"
	"reminder_bot/internal/bot/pkg/reminder"
	"reminder_bot/internal/bot/pkg/user"
	"reminder_bot/internal/bot/telegram"
	useCase2 "reminder_bot/internal/bot/useCase"
	"reminder_bot/internal/models"
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
	
	botApi, err := tgbotapi.NewBotAPI(cfg.Telegram.Token)
	if err != nil {
		log.Fatal(err)
	}
	
	amqpConn, err := rabbitmq.NewRabbitMQConn(cfg)
	if err != nil {
		log.Fatal(err)
	}
	
	userService, err := user.NewUserService(amqpConn, cfg.Queue, cfg.Exchange)
	if err != nil {
		log.Fatal(err)
	}
	parsingService, err := parsing.NewParsingService(amqpConn, cfg.Queue, cfg.Exchange)
	if err != nil {
		log.Fatal(err)
	}
	
	reminderService, err := reminder.NewReminderService(amqpConn, cfg.Queue, cfg.Exchange)
	if err != nil {
		log.Fatal(err)
	}
	
	useCase := useCase2.NewBotUseCase(userService, parsingService, reminderService)
	
	cache := cache2.NewMemoryCache()
	
	mapClient, err := maps.NewClient(maps.WithAPIKey(cfg.GoogleMapApiKey))
	if err != nil {
		log.Fatal(err)
	}
	
	go NotifyConsumer(amqpConn, botApi, cfg.Queue.BotEventNotifySend)
	
	bot := telegram.NewBot(botApi, cache, useCase, mapClient)
	if err := bot.Start(); err != nil {
		log.Fatal(err)
	}
}

func NotifyConsumer(conn *amqp.Connection, bot *tgbotapi.BotAPI, queueName string) {
	
	ch, err := conn.Channel()
	if err != nil {
		return
	}
	queue, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	
	if err != nil {
		return
	}
	
	deliveries, err := ch.Consume(
		queue.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return
	}
	
	for d := range deliveries {
		req := models.NotifyRequest{}
		if err := json.Unmarshal(d.Body, &req); err != nil {
			fmt.Println(err)
			continue
		}
		
		msg := tgbotapi.NewMessage(req.ChatId, "Напоминаем: \n" + req.Text)
		_, err := bot.Send(msg)
		if err != nil {
			fmt.Println(err)
			continue
		}
	}
}
