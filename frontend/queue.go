package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/streadway/amqp"
	"log"
	"time"
)

const rabbitHost = "queue"
const rabbitPort = "5672"

const queueName = "tasks"

type Queue interface {
	Close()
	PublishTask(interface{}) error
}

type rabbitMQ struct {
	conn *amqp.Connection
	ch 	 *amqp.Channel
}

func NewRabbitMQQueue() (r *rabbitMQ, err error) {
	conn, err := retryDial()
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	_, err = ch.QueueDeclare(queueName, true, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	log.Printf("Initialized connection to rabbit")

	return &rabbitMQ{
		conn: conn,
		ch: ch,
	}, nil
}

func (r *rabbitMQ) Close() {
	r.ch.Close()
	r.conn.Close()
}

func (r *rabbitMQ) PublishTask(task interface{}) error {
	body, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return r.ch.Publish(
		"",
		queueName,
		false,
		false,
		amqp.Publishing{
			DeliveryMode: amqp.Persistent,
			ContentType:  "text/json",
			Body: 		  body,
		},
	)
}

func retryDial() (*amqp.Connection, error) {
	var retries = 10
	var err error
	for i := 0; i < retries; i++ {
		var conn *amqp.Connection
		conn, err = amqp.Dial("amqp://" + rabbitHost + ":" + rabbitPort)
		if err == nil {
			return conn, nil
		}
		time.Sleep(3*time.Second)
	}

	return nil, errors.New(fmt.Sprintf("Didn't succeed after %d retries: %s", retries, err))
}
