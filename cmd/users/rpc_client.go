package main

import (
	"errors"
	"fmt"
	"github.com/google/uuid"
	"log"
	"math/rand"
	"strconv"
	"sync"
	"time"

	"github.com/streadway/amqp"
)

type Client interface {
	Close()
	RemoteCall(n int) (int, error)
}

type ClientConfig struct {
	Url         string
	ServerQueue string
	Timeout     time.Duration
}

func NewClientConfig(url, serverQueue string, timeout time.Duration) ClientConfig {
	return ClientConfig{
		Url:         url,
		ServerQueue: serverQueue,
		Timeout:     timeout,
	}
}

func Connect(cfg ClientConfig) (Client, error) {
	conn, err := amqp.Dial(cfg.Url)
	if err != nil {
		return nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	queue, err := channel.QueueDeclare(
		"reply.user.command.fib.get",    // name
		false, // durable
		true,  // delete when used
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
		true,       // auto-ack
		false,      // exclusive
		false,      // no-local
		false,      // no-wait
		nil,        // args
	)
	if err != nil {
		return nil, err
	}
	client := newClient(cfg.ServerQueue, conn, channel, &queue, cfg.Timeout)
	go client.handleDeliveries(msgs)

	return client, nil
}

type clientImpl struct {
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

func newClient(serverQueue string, conn *amqp.Connection, channel *amqp.Channel, queue *amqp.Queue, timeout time.Duration) *clientImpl {
	return &clientImpl{
		serverQueue: serverQueue,
		conn:        conn,
		channel:     channel,
		queue:       queue,
		calls:       make(map[string]*pendingCall),
		timeout:     timeout,
		done:        make(chan bool)}
}

func (client *clientImpl) RemoteCall(n int) (int, error) {

	expiration := fmt.Sprintf("%d", client.timeout)
	corrId := newCorrId()
	err := client.channel.Publish(
		"",                 // exchange
		client.serverQueue, // routing key
		false,              // mandatory
		false,              // immediate
		amqp.Publishing{
			ContentType:   "application/octet-stream",
			CorrelationId: corrId,
			ReplyTo:       client.queue.Name,
			Body:          []byte(strconv.Itoa(n)),
			Expiration:    expiration,
		})
	if err != nil {
		return 0, err
	}
	fmt.Println("Queue name: ", client.queue.Name)
	call := &pendingCall{done: make(chan bool)}

	client.guard.Lock()
	client.calls[corrId] = call
	client.guard.Unlock()

	var respData int
	var respError error = errors.New("timeout")

	select {
	case <-call.done:
		respData, respError = strconv.Atoi(string(call.data))
		if respError != nil {
			fmt.Println(err)
		}
	case <-time.After(client.timeout):
		break
	}

	client.guard.Lock()
	delete(client.calls, corrId)
	client.guard.Unlock()

	return respData, respError
}

func (client *clientImpl) Close() {
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

func newCorrId() string {
	return uuid.New().String()
}

func (client *clientImpl) handleDeliveries(msgs <-chan amqp.Delivery) {
	finish := false
	for !finish {
		select {
		case msg := <-msgs:
			fmt.Println(msg.ReplyTo)
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

func main() {
	rand.Seed(time.Now().UTC().UnixNano())

	cfg := NewClientConfig("amqp://guest:guest@localhost:5672/", "user.command.fib.get", time.Second)
	client, err := Connect(cfg)
	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	for i := 1; i < 10; i++ {
		time.Sleep(time.Second)
		res, err := client.RemoteCall(i)
		if err != nil {
			log.Println(err)
		}
		fmt.Println("Result: ", res)
	}
}
