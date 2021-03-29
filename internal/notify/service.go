package notify

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jasonlvhit/gocron"
	"github.com/streadway/amqp"
	"reminder_bot/internal/models"
)

type Notify struct {
	conn      *amqp.Connection
	channel   *amqp.Channel
	queueName string
	uc        UseCase
}

func NewNotifyService(conn *amqp.Connection, uc UseCase, queueName string) (*Notify, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, errors.New("error amqpConn.Channel")
	}
	return &Notify{
		conn:      conn,
		uc:        uc,
		channel:   ch,
		queueName: queueName,
	}, nil
}

func (n *Notify) Start(ctx context.Context) error {

	s := gocron.NewScheduler()
	s.Every(1).Minute().Do(n.checkNew)
	select {
	case <-s.Start():
		return nil
	case chanErr := <-n.channel.NotifyClose(make(chan *amqp.Error)):
		fmt.Printf("NotifyConsumer.Close: %v\n", chanErr)
		return nil
	case <-ctx.Done():
		return nil
	}
}

func (n *Notify) checkNew() error {

	notifies, err := n.uc.GetNotifies()
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, notify := range notifies {
		go func(item models.NotifyRequest) {
			err := n.publish(item)
			if err != nil {
				fmt.Println(err)
			}
		}(notify)
	}
	return nil
}

func (n *Notify) publish(notify models.NotifyRequest) error {

	body, err := json.Marshal(notify)
	if err != nil {
		fmt.Println(err)
		return err
	}
	err = n.channel.Publish(
		"",          // exchange
		n.queueName, // routing key
		false,       // mandatory
		false,       // immediate
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		})
	if err != nil {
		return err
	}
	go func() {
		_ = n.uc.DeactivateReminder(notify.Id)
	}()
	return nil
}
