package reminder

import (
	"github.com/streadway/amqp"
	"reminder_bot/config"
	"reminder_bot/internal/bot/pkg/rpc"
	"time"
)

type Service struct {
	add *rpc.Client
}

func NewReminderService(conn *amqp.Connection, queue config.Queue, exchange config.Exchange) (*Service, error) {
	addReminder, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.ReminderCommandReminderAdd,
		Timeout:     time.Second,
	})
	if err != nil {
		return nil, err
	}
	
	return &Service{
		add: addReminder,
	}, nil
}

func (s *Service) AddReminder(data []byte) error {
	
	_, err := s.add.RemoteCall(data)
	return err
}
