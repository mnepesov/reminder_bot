package rpc

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
	"reminder_bot/internal/models"
	"reminder_bot/pkg/service_errors"
	"sync"
	"time"
)

type ClientConfig struct {
	ServerQueue string
	Timeout     time.Duration
}

func Connect(amqpConn *amqp.Connection, cfg ClientConfig) (*Client, error) {

	channel, err := amqpConn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"",    // name
		false, // durable
		true,  // delete when usused
		false, // exclusive
		false, // noWait
		nil,   // arguments
	)
	if err != nil {
		return nil, err
	}

	msgs, err := channel.Consume(
		queue.Name, // queue
		"",         // consumer
		false,      // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)

	if err != nil {
		return nil, err
	}

	client := newClient(cfg.ServerQueue, amqpConn, channel, &queue, cfg.Timeout)
	go client.handleDeliveries(msgs)

	return client, nil
}

type Client struct {
	conn        *amqp.Connection
	channel     *amqp.Channel
	queue       *amqp.Queue
	serverQueue string
	guard       sync.Mutex
	calls       map[string]*pendingCall
	timeout     time.Duration
	done        chan bool
}

type pendingCall struct {
	done chan bool
	data []byte
}

func newClient(serverQueue string, conn *amqp.Connection, channel *amqp.Channel, queue *amqp.Queue, timeout time.Duration) *Client {
	return &Client{
		serverQueue: serverQueue,
		conn:        conn,
		channel:     channel,
		queue:       queue,
		calls:       make(map[string]*pendingCall),
		timeout:     timeout,
		done:        make(chan bool)}
}

func (client *Client) Close() {
	if client == nil {
		return
	}

	client.done <- true

	if client.channel != nil {
		client.channel.Close()
	}

	if client.conn != nil {
		client.conn.Close()
	}
}

func (client *Client) RemoteCall(body []byte) ([]byte, error) {

	expiration := fmt.Sprintf("%d", client.timeout)
	corrId := newCorrId()
	err := client.channel.Publish(
		"",                 // exchange
		client.serverQueue, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "application/json",
			CorrelationId: corrId,
			ReplyTo:       client.queue.Name,
			Body:          body,
			Expiration:    expiration,
		})
	if err != nil {
		return nil, err
	}

	call := &pendingCall{done: make(chan bool)}

	client.guard.Lock()
	client.calls[corrId] = call
	client.guard.Unlock()

	var respData []byte
	var respError error

	select {
	case <-call.done:
		var resp models.Response
		respError = resp.Unmarshal(call.data)
		if respError == nil {
			if resp.IsSuccess {
				respData = resp.Body
			} else {
				respError = errors.New(resp.ErrText)
			}
		}

	case <-time.After(client.timeout):
		respError = service_errors.Timeout
		break
	}

	client.guard.Lock()
	delete(client.calls, corrId)
	client.guard.Unlock()

	return respData, respError
}

func (client *Client) handleDeliveries(msgs <-chan amqp.Delivery) {
	finish := false
	for !finish {
		select {
		case msg := <-msgs:
			call, ok := client.calls[msg.CorrelationId]
			if ok {
				call.data = msg.Body
				call.done <- true
			}
		case <-client.done:
			finish = true
		}
	}
}

func newCorrId() string {
	return uuid.New().String()
}

