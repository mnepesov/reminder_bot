package useCase

import (
	"encoding/json"
	"reminder_bot/internal/bot/pkg/parsing"
	"reminder_bot/internal/bot/pkg/reminder"
	"reminder_bot/internal/bot/pkg/user"
	"reminder_bot/internal/models"
)

type BotUseCase struct {
	userService     *user.Service
	parsingService  *parsing.Service
	reminderService *reminder.Service
}

func NewBotUseCase(userService *user.Service, parsingService *parsing.Service, reminderService *reminder.Service) *BotUseCase {
	return &BotUseCase{
		userService:     userService,
		parsingService:  parsingService,
		reminderService: reminderService,
	}
}

func (b *BotUseCase) GetUser(userId int) (models.User, error) {

	data := models.GetUser{Id: userId}
	jsonData, err := json.Marshal(data)
	if err != nil {
		return models.User{}, err
	}
	respUser, err := b.userService.GetUser(jsonData)
	if err != nil {
		return models.User{}, err
	}
	u := models.User{}
	if err := json.Unmarshal(respUser, &u); err != nil {
		return models.User{}, err
	}

	return u, err
}

func (b *BotUseCase) CreateUser(user models.User) error {
	jsonUser, err := json.Marshal(user)
	if err != nil {
		return err
	}

	err = b.userService.CreateUser(jsonUser)
	return err
}

func (b *BotUseCase) UpdateTimezone(userId int, timezone string) error {
	tz := models.UpdateTimezone{
		UserId:   userId,
		Timezone: timezone,
	}
	jsonData, err := json.Marshal(tz)
	if err != nil {
		return err
	}
	err = b.userService.UpdateTimezone(jsonData)
	if err != nil {
		return err
	}
	return nil
}

func (b *BotUseCase) Parse(req models.ParseRequest) (models.ParseResponse, error) {

	jsonData, err := json.Marshal(req)
	if err != nil {
		return models.ParseResponse{}, err
	}

	data, err := b.parsingService.Parse(jsonData)
	if err != nil {
		return models.ParseResponse{}, err
	}

	res := models.ParseResponse{}
	if err := json.Unmarshal(data, &res); err != nil {
		return models.ParseResponse{}, err
	}

	return res, nil
}

func (b *BotUseCase) AddReminder(reminder models.AddReminderRequest) error {
	data, err := json.Marshal(reminder)
	if err != nil {
		return err
	}

	err = b.reminderService.AddReminder(data)
	return err
}

func (b *BotUseCase) GetRemindersByUserId(req models.GetRemindersRequest) ([]models.Reminder, error) {
	data, err := json.Marshal(req)
	if err != nil {
		return nil, err
	}

	res, err := b.reminderService.GetRemindersByUserId(data)
	if err != nil {
		return nil, err
	}

	var reminders []models.Reminder
	err = json.Unmarshal(res, &reminders)
	if err != nil {
		return nil, err
	}
	return reminders, nil
}
