package user

import (
	"github.com/streadway/amqp"
	"reminder_bot/config"
	"reminder_bot/internal/bot/pkg/rpc"
	"time"
)

type Service struct {
	getUser        *rpc.Client
	createUser     *rpc.Client
	updateTimezone *rpc.Client
}

func NewUserService(conn *amqp.Connection, queue config.Queue, exchange config.Exchange) (*Service, error) {
	get, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.UsersRepoCommandUsersGet,
		Timeout:     time.Second,
	})
	if err != nil {
		return nil, err
	}

	create, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.UsersRepoCommandUsersCreate,
		Timeout:     time.Second,
	})
	if err != nil {
		return nil, err
	}

	updateTimezone, err := rpc.Connect(conn, rpc.ClientConfig{
		ServerQueue: queue.UsersRepoCommandUsersUpdateTimezone,
		Timeout:     time.Second * 2,
	})
	if err != nil {
		return nil, err
	}

	return &Service{
		getUser:        get,
		createUser:     create,
		updateTimezone: updateTimezone,
	}, nil
}

func (u *Service) GetUser(data []byte) ([]byte, error) {
	return u.getUser.RemoteCall(data)
}

func (u *Service) CreateUser(data []byte) error {
	_, err := u.createUser.RemoteCall(data)
	return err
}

func (u *Service) UpdateTimezone(data []byte) error {
	_, err := u.updateTimezone.RemoteCall(data)
	return err
}
