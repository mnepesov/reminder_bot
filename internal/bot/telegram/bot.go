package telegram

import (
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"googlemaps.github.io/maps"
	"log"
	"reminder_bot/internal/bot/pkg/cache"
	"reminder_bot/internal/bot/useCase"
	"reminder_bot/internal/models"
)

type Bot struct {
	bot      *tgbotapi.BotAPI
	cache    cache.Cache
	useCase  useCase.UseCase
	gmClient *maps.Client
}

func NewBot(bot *tgbotapi.BotAPI, cache cache.Cache, useCase useCase.UseCase, mapClient *maps.Client) *Bot {
	return &Bot{
		bot:      bot,
		cache:    cache,
		useCase:  useCase,
		gmClient: mapClient,
	}
}

func (b *Bot) Start() error {
	b.bot.Debug = false
	
	log.Printf("Authorized on account %s", b.bot.Self.UserName)
	
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60
	updates, err := b.bot.GetUpdatesChan(u)
	if err != nil {
		return err
	}
	for update := range updates {
		
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}
		message := update.Message
		user, err := b.cache.GetUserInfo(message.From.ID)
		if err != nil || user == nil {
			u, err := b.useCase.GetUser(message.From.ID)
			if err != nil {
				err := b.useCase.CreateUser(models.User{
					Id:     message.From.ID,
					ChatId:   message.Chat.ID,
					Username: message.From.UserName,
					FullName: message.From.FirstName + " " + message.From.LastName,
				})
				if err != nil {
					fmt.Println(err)
					msg := tgbotapi.NewMessage(message.Chat.ID, errorText)
					msg.ReplyMarkup = nil
					_, err = b.bot.Send(msg)
					continue
				} else {
					_ = b.cache.SetNavigation(message.From.ID, "start")
					
					msg := tgbotapi.NewMessage(message.Chat.ID, helloText)
					msg.ReplyMarkup = mainMenu
					_, err = b.bot.Send(msg)
					continue
				}
			} else {
				_ = b.cache.SetUserInfo(message.From.ID, &u)
				user = &u
			}
		}
		if update.Message.IsCommand() {
			err := b.handleCommand(update.Message)
			if err != nil {
				fmt.Println(err)
			}
			continue
		}
		err = b.handleMessage(update.Message)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}
