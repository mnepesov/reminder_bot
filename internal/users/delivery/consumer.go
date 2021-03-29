package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"reminder_bot/internal/models"
	"reminder_bot/internal/users"
)

const (
	queueDurable    = true
	queueAutoDelete = false
	queueExclusive  = false
	queueNoWait     = false

	consumeAutoAck   = false
	consumeExclusive = false
	consumeNoLocal   = false
	consumeNoWait    = false
)

type UserConsumer struct {
	conn    *amqp.Connection
	channel *amqp.Channel
	useCase users.UseCase
}

func NewUserConsumer(amqpConn *amqp.Connection, useCase users.UseCase) (*UserConsumer, error) {
	ch, err := amqpConn.Channel()
	if err != nil {
		return nil, errors.New("Error amqpConn.Channel")
	}
	return &UserConsumer{
		conn:    amqpConn,
		channel: ch,
		useCase: useCase,
	}, nil
}

//get user by id
func (c *UserConsumer) StartGetUserConsumer(ctx context.Context, workerPoolSize int, queueName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queue, err := c.channel.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)

	if err != nil {
		return err
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		"",
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.getUserWorker(ctx, deliveries)
	}
	select {
	case chanErr := <-c.channel.NotifyClose(make(chan *amqp.Error)):
		fmt.Printf("GetUserConsumer.Close: %v\n", chanErr)
		return err
	case <-ctx.Done():
		return nil
	}
}
func (c *UserConsumer) getUserWorker(ctx context.Context, messages <-chan amqp.Delivery) {
	for d := range messages {
		fmt.Println("New Request: ")
		u := &models.GetUser{}
		if err := json.Unmarshal(d.Body, u); err != nil {
			_ = d.Reject(false)
			fmt.Println(err)
			continue
		}
		fmt.Println("[GetUserWorker]: ", u.Id)
		var resp models.Response
		user, err := c.useCase.GetUserById(u.Id)
		if err != nil {
			resp.IsSuccess = false
			resp.ErrText = err.Error()
		} else {
			resp.IsSuccess = true
			resp.Body = user
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			panic(fmt.Sprintf("Failed marshall responce: %v", err))
		}

		err = c.channel.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   d.ContentType,
				CorrelationId: d.CorrelationId,
				Body:          respData,
			})

		if err != nil {
			panic(fmt.Sprintf("Failed to publish a message: %v", err))
		}

		_ = d.Ack(false)
	}
}

//create user
func (c *UserConsumer) StartCreateUserConsumer(workerPoolSize int, queueName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queue, err := c.channel.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)

	if err != nil {
		return err
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		"",
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.createUserWorker(ctx, deliveries)
	}
	chanErr := <-c.channel.NotifyClose(make(chan *amqp.Error))
	fmt.Printf("GetUserConsumer.Close: %v\n", chanErr)
	return chanErr
}
func (c *UserConsumer) createUserWorker(ctx context.Context, messages <-chan amqp.Delivery) {
	for d := range messages {
		u := models.User{}
		if err := json.Unmarshal(d.Body, &u); err != nil {
			_ = d.Reject(false)
			fmt.Println(err)
			continue
		}
		fmt.Println("[createUserWorker]: ", u)
		var resp models.Response
		err := c.useCase.CreateUser(u)
		if err != nil {
			resp.IsSuccess = false
			resp.ErrText = err.Error()
		} else {
			resp.IsSuccess = true
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			panic(fmt.Sprintf("Failed marshall responce: %v", err))
		}

		err = c.channel.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   d.ContentType,
				CorrelationId: d.CorrelationId,
				Body:          respData,
			})

		if err != nil {
			panic(fmt.Sprintf("Failed to publish a message: %v", err))
		}

		_ = d.Ack(false)
	}
}

//create user
func (c *UserConsumer) StartUpdateTimeZoneConsumer(workerPoolSize int, queueName string) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	queue, err := c.channel.QueueDeclare(
		queueName,
		queueDurable,
		queueAutoDelete,
		queueExclusive,
		queueNoWait,
		nil,
	)

	if err != nil {
		return err
	}

	deliveries, err := c.channel.Consume(
		queue.Name,
		"",
		consumeAutoAck,
		consumeExclusive,
		consumeNoLocal,
		consumeNoWait,
		nil,
	)
	if err != nil {
		return err
	}

	for i := 0; i < workerPoolSize; i++ {
		go c.updateTimezoneWorker(ctx, deliveries)
	}
	chanErr := <-c.channel.NotifyClose(make(chan *amqp.Error))
	fmt.Printf("updateTimezoneConsumer.Close: %v\n", chanErr)
	return chanErr
}
func (c *UserConsumer) updateTimezoneWorker(ctx context.Context, messages <-chan amqp.Delivery) {
	for d := range messages {
		u := models.UpdateTimezone{}
		if err := json.Unmarshal(d.Body, &u); err != nil {
			_ = d.Reject(false)
			fmt.Println(err)
			continue
		}
		fmt.Println("[UpdateTimezoneWorker]: ", u)
		var resp models.Response
		err := c.useCase.UpdateTimezone(u.UserId, u.Timezone)
		if err != nil {
			resp.IsSuccess = false
			resp.ErrText = err.Error()
		} else {
			resp.IsSuccess = true
		}

		respData, err := json.Marshal(resp)
		if err != nil {
			panic(fmt.Sprintf("Failed marshall responce: %v", err))
		}

		err = c.channel.Publish(
			"",
			d.ReplyTo,
			false,
			false,
			amqp.Publishing{
				ContentType:   d.ContentType,
				CorrelationId: d.CorrelationId,
				Body:          respData,
			})

		if err != nil {
			panic(fmt.Sprintf("Failed to publish a message: %v", err))
		}

		_ = d.Ack(false)
	}
}
