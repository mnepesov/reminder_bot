package reminder

import (
	"github.com/streadway/amqp"
	"reminder_bot/config"
	"reminder_bot/internal/bot/pkg/rpc"
	"time"
)

type Service struct {
	add *rpc.Client
	get *rpc.Client
}

func NewReminderService(conn *amqp.Connection, queue config.Queue, exchange config.Exchange) (*Service, error) {
	addReminder, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.ReminderCommandReminderAdd,
		Timeout:     time.Second,
	})
	get, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.ReminderCommandRemindersGet,
		Timeout:     time.Second,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		add: addReminder,
		get: get,
	}, nil
}

func (s *Service) AddReminder(data []byte) error {

	_, err := s.add.RemoteCall(data)
	return err
}

func (s *Service) GetRemindersByUserId(data []byte) ([]byte, error) {
	data, err := s.get.RemoteCall(data)
	return data, err
}
