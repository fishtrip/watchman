package main

import (
	"log"
	"net/http"
	"time"

	"github.com/streadway/amqp"
)

type NotifyResponse int

const (
	NotifySuccess = 1
	NotifyFailure = 0
)

type Message struct {
	queueConfig    QueueConfig
	amqpDelivery   *amqp.Delivery // message read from rabbitmq
	notifyResponse NotifyResponse // notify result from callback url
}

func (m Message) CurrentMessageRetries() int {
	msg := m.amqpDelivery

	xDeathArray, ok := msg.Headers["x-death"].([]interface{})
	if !ok {
		m.Printf("x-death array case fail")
		return 0
	}

	if len(xDeathArray) <= 0 {
		return 0
	}

	for _, h := range xDeathArray {
		xDeathItem := h.(amqp.Table)

		if xDeathItem["reason"] == "rejected" {
			return int(xDeathItem["count"].(int64))
		}
	}

	return 0
}

func (m *Message) Notify(client *http.Client) *Message {
	qc := m.queueConfig
	msg := m.amqpDelivery

	client.Timeout = time.Duration(qc.NotifyTimeoutWithDefault()) * time.Second
	statusCode := notifyUrl(client, qc.NotifyUrl(), msg.Body)

	m.Printf("notify url %s, result: %d", qc.NotifyUrl(), statusCode)

	if statusCode == 200 || statusCode == 201 {
		m.notifyResponse = NotifySuccess
	} else {
		m.notifyResponse = NotifyFailure
	}

	return m
}

func (m Message) IsMaxRetry() bool {
	retries := m.CurrentMessageRetries()
	maxRetries := m.queueConfig.RetryTimesWithDefault()
	return retries >= maxRetries
}

func (m Message) IsNotifySuccess() bool {
	return m.notifyResponse == NotifySuccess
}

func (m Message) Ack() error {
	m.Printf("acker: ack message")
	err := m.amqpDelivery.Ack(false)
	LogOnError(err)
	return err
}

func (m Message) Reject() error {
	m.Printf("acker: reject message")
	err := m.amqpDelivery.Reject(false)
	LogOnError(err)
	return err
}

func (m Message) Republish(out chan<- Message) error {
	m.Printf("acker: ERROR republish message")
	out <- m
	err := m.amqpDelivery.Ack(false)
	LogOnError(err)
	return err
}

func (m Message) CloneAndPublish(channel *amqp.Channel) error {
	msg := m.amqpDelivery
	qc := m.queueConfig

	errMsg := cloneToPublishMsg(msg)
	err := channel.Publish(qc.ErrorExchangeName(), msg.RoutingKey, false, false, *errMsg)
	LogOnError(err)
	return err
}

func (m Message) Printf(v ...interface{}) {
	msg := m.amqpDelivery

	vv := []interface{}{}
	vv = append(vv, msg.MessageId, msg.RoutingKey)
	vv = append(vv, v[1:]...)

	log.Printf("[%s] [%s] "+v[0].(string), vv...)
}
