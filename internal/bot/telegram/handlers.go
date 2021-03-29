package telegram

import (
	"context"
	"fmt"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"googlemaps.github.io/maps"
	"reminder_bot/internal/models"
	"time"
)

const (
	errorText = "Упс! Что-то пошло не так\n Попробуйте еще раз\n /start"
	helloText = "Добро пожаловать в чат-бот Reminder!\nВаш часовой пояс: Europe/Minsk\n\nО чем вам напомнить?"
	
	commandStart     = "start"
	commandStartText = "О чем вам напомнить?"
	
	addText      = "Пожалуйста, введите текст напоминания и время. ☝️Также обратите внимание, что вы можете создавать напоминания не нажимая никаких дополнительных кнопок и не отправляя команд, просто отправьте текст напоминания и, опционально, время."
	settingsText = "Ваши текущие настройки:\n\n - часовой пояс: %s"
	timezoneText = "Пожалуйста, отправьте мне название своего города, чтобы я мог определить ваш часовой пояс. \nНапример: Minsk"
	webText      = "Сайт временно не доступен"
	listText     = "Пусто"
)

//commands
func (b *Bot) handleCommand(message *tgbotapi.Message) error {
	switch message.Command() {
	case commandStart:
		return b.handleStartCommand(message)
	default:
		return b.handleUnknownCommand(message)
	}
}

func (b *Bot) handleStartCommand(message *tgbotapi.Message) error {
	
	_ = b.cache.SetNavigation(message.From.ID, "start")
	
	msg := tgbotapi.NewMessage(message.Chat.ID, commandStartText)
	msg.ReplyMarkup = mainMenu
	_, err := b.bot.Send(msg)
	return err
}

func (b *Bot) handleUnknownCommand(message *tgbotapi.Message) error {
	msg := tgbotapi.NewMessage(message.Chat.ID, "я не знаю эту команду")
	_, err := b.bot.Send(msg)
	return err
}

//handle messages
func (b *Bot) handleMessage(message *tgbotapi.Message) error {
	switch message.Text {
	case "➕Добавить":
		return b.handleAddMessage(message)
	case "⚠Настройки":
		return b.handleSettingsMessage(message)
	case "Часовой пояс":
		return b.handleTimezoneMessage(message)
	case "🌏Web":
		return b.handleWebMessage(message)
	case "📝Список":
		return b.handleListMessage(message)
	case "Отмена":
		return b.handleExitMessage(message)
	default:
		return b.handleUnknownMessage(message)
	}
}

//add message
func (b *Bot) handleAddMessage(message *tgbotapi.Message) error {
	_ = b.cache.SetNavigation(message.From.ID, "add")
	msg := tgbotapi.NewMessage(message.Chat.ID, addText)
	msg.ReplyMarkup = exitMenu
	_, err := b.bot.Send(msg)
	return err
}

//settings message
func (b *Bot) handleSettingsMessage(message *tgbotapi.Message) error {
	user, _ := b.cache.GetUserInfo(message.From.ID)
	_ = b.cache.SetNavigation(message.From.ID, "settings")
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf(settingsText, user.Timezone))
	msg.ReplyMarkup = settingMenu
	_, err := b.bot.Send(msg)
	return err
}

//timezone message
func (b *Bot) handleTimezoneMessage(message *tgbotapi.Message) error {
	_ = b.cache.SetNavigation(message.From.ID, "timezone")
	msg := tgbotapi.NewMessage(message.Chat.ID, timezoneText)
	msg.ReplyMarkup = exitMenu
	_, err := b.bot.Send(msg)
	return err
}

//web message
func (b *Bot) handleWebMessage(message *tgbotapi.Message) error {
	_ = b.cache.SetNavigation(message.From.ID, "web")
	msg := tgbotapi.NewMessage(message.Chat.ID, webText)
	msg.ReplyMarkup = mainMenu
	_, err := b.bot.Send(msg)
	return err
}

//list
func (b *Bot) handleListMessage(message *tgbotapi.Message) error {
	_ = b.cache.SetNavigation(message.From.ID, "list")
	msg := tgbotapi.NewMessage(message.Chat.ID, listText)
	msg.ReplyMarkup = mainMenu
	_, err := b.bot.Send(msg)
	return err
}

//exit message
func (b *Bot) handleExitMessage(message *tgbotapi.Message) error {
	navigation, _ := b.cache.GetNavigation(message.From.ID)
	switch navigation {
	case "add":
		return b.handleStartCommand(message)
	case "settings":
		return b.handleStartCommand(message)
	case "timezone":
		return b.handleSettingsMessage(message)
	default:
		return b.handleStartCommand(message)
	}
}

//unknown message
func (b *Bot) handleUnknownMessage(message *tgbotapi.Message) error {
	navigation, _ := b.cache.GetNavigation(message.From.ID)
	switch navigation {
	case "start":
		return b.addReminders(message)
	case "add":
		return b.addReminders(message)
	case "timezone":
		{
			city := message.Text
			
			gcReq := &maps.GeocodingRequest{
				Address: city,
			}
			gcRes, err := b.gmClient.Geocode(context.Background(), gcReq)
			if err != nil {
				_ = b.cache.SetNavigation(message.From.ID, "timezone")
				msg := tgbotapi.NewMessage(message.Chat.ID, "😐Не могу определить ваше местоположение, пожалуйста, попробуйте изменить ваш запрос")
				msg.ReplyMarkup = exitMenu
				_, err = b.bot.Send(msg)
			}
			tzReq := &maps.TimezoneRequest{
				Location: &maps.LatLng{
					Lat: gcRes[0].Geometry.Location.Lat,
					Lng: gcRes[0].Geometry.Location.Lng,
				},
				Timestamp: time.Now(),
				Language:  "en",
			}
			tzRes, err := b.gmClient.Timezone(context.Background(), tzReq)
			if err != nil {
				_ = b.cache.SetNavigation(message.From.ID, "timezone")
				msg := tgbotapi.NewMessage(message.Chat.ID, "😐Не могу определить ваше местоположение, пожалуйста, попробуйте изменить ваш запрос")
				msg.ReplyMarkup = exitMenu
				_, err = b.bot.Send(msg)
			}
			err = b.useCase.UpdateTimezone(message.From.ID, tzRes.TimeZoneID)
			if err != nil {
				fmt.Println(err)
				_ = b.cache.SetNavigation(message.From.ID, "timezone")
				msg := tgbotapi.NewMessage(message.Chat.ID, "😐Не могу определить ваше местоположение, пожалуйста, попробуйте изменить ваш запрос")
				msg.ReplyMarkup = exitMenu
				_, err = b.bot.Send(msg)
			}
			_ = b.cache.SetNavigation(message.From.ID, "settings")
			msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Ваш часовой пояс определен\nВаши текущие настройки:\n\n - язык: Русский\n - часовой пояс: %s", tzRes.TimeZoneID))
			msg.ReplyMarkup = settingMenu
			_, err = b.bot.Send(msg)
			return err
		}
	default:
		_ = b.cache.SetNavigation(message.From.ID, "start")
		msg := tgbotapi.NewMessage(message.Chat.ID, "Не могу обработать ваш текст")
		msg.ReplyMarkup = mainMenu
		_, err := b.bot.Send(msg)
		return err
	}
}

func (b *Bot) addReminders(message *tgbotapi.Message) error {
	res, err := b.useCase.Parse(models.ParseRequest{Text: message.Text})
	if err != nil {
		_ = b.cache.SetNavigation(message.From.ID, "start")
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("✅Напоминание не добавлено\n\n%s", err.Error()))
		msg.ReplyMarkup = mainMenu
		_, err = b.bot.Send(msg)
		return err
	}
	
	err = b.useCase.AddReminder(models.AddReminderRequest{
		UserId: message.From.ID,
		Text:   res.Text,
		Date:   res.Time,
	})
	
	if err != nil {
		_ = b.cache.SetNavigation(message.From.ID, "start")
		msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминание не добавлено\n\n%s", err.Error()))
		msg.ReplyMarkup = mainMenu
		_, err = b.bot.Send(msg)
		return err
	}
	
	_ = b.cache.SetNavigation(message.From.ID, "start")
	msg := tgbotapi.NewMessage(message.Chat.ID, fmt.Sprintf("Напоминание добавлено\n\n%s\n%s", res.Text, res.Time))
	msg.ReplyMarkup = mainMenu
	_, err = b.bot.Send(msg)
	return err
}
