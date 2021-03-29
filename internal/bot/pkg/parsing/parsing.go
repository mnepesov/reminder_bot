package parsing

import (
	"github.com/streadway/amqp"
	"reminder_bot/config"
	"reminder_bot/internal/bot/pkg/rpc"
	"time"
)

type Service struct {
	parse *rpc.Client
}

func NewParsingService(conn *amqp.Connection, queue config.Queue, exchange config.Exchange) (*Service, error) {
	parse, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.ParsingCommandTextParse,
		Timeout:     time.Second * 2,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		parse: parse,
	}, nil
}

func (s *Service) Parse(data []byte) ([]byte, error) {
	return s.parse.RemoteCall(data)
}
