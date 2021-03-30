package delivery

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	"github.com/streadway/amqp"
	"reminder_bot/internal/models"
	"reminder_bot/internal/reminder"
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
	useCase reminder.UseCase
}

func NewReminderConsumer(amqpConn *amqp.Connection, useCase reminder.UseCase) (*UserConsumer, error) {
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

//add reminder
func (c *UserConsumer) StartAddReminderConsumer(ctx context.Context, workerPoolSize int, queueName string) error {
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
		go c.addReminderWorker(ctx, deliveries)
	}
	select {
	case chanErr := <-c.channel.NotifyClose(make(chan *amqp.Error)):
		fmt.Printf("GetUserConsumer.Close: %v\n", chanErr)
		return err
	case <-ctx.Done():
		return nil
	}
}

func (c *UserConsumer) addReminderWorker(ctx context.Context, messages <-chan amqp.Delivery) {
	for d := range messages {
		fmt.Println("New Request: ")
		reminder := models.AddReminderRequest{}
		if err := json.Unmarshal(d.Body, &reminder); err != nil {
			_ = d.Reject(false)
			fmt.Println(err)
			continue
		}
		fmt.Println("[AddReminderWorker]: ", reminder.UserId)
		var resp models.Response
		err := c.useCase.AddReminders(reminder)
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

//get user by id
func (c *UserConsumer) StartGetRemindersConsumer(ctx context.Context, workerPoolSize int, queueName string) error {
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
		go c.getRemindersWorker(ctx, deliveries)
	}
	select {
	case chanErr := <-c.channel.NotifyClose(make(chan *amqp.Error)):
		fmt.Printf("GetUserConsumer.Close: %v\n", chanErr)
		return err
	case <-ctx.Done():
		return nil
	}
}

func (c *UserConsumer) getRemindersWorker(ctx context.Context, messages <-chan amqp.Delivery) {
	for d := range messages {
		fmt.Println("New Request: ")
		req := models.GetRemindersRequest{}
		if err := json.Unmarshal(d.Body, &req); err != nil {
			_ = d.Reject(false)
			fmt.Println(err)
			continue
		}
		fmt.Println("[GetRemindersWorker]: ", req.UserId)
		var resp models.Response
		remindersList, err := c.useCase.GetRemindersByUserId(req.UserId)
		if err != nil {
			resp.IsSuccess = false
			resp.ErrText = err.Error()
		} else {
			resp.IsSuccess = true
			data, _ := json.Marshal(remindersList)
			resp.Body = data
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
