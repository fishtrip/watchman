package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/streadway/amqp"
)

func LogOnError(err error) {
	if err != nil {
		fmt.Printf("ERROR - %s\n", err)
	}
}

func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

func setupChannel() (*amqp.Connection, *amqp.Channel, error) {
	url := os.Getenv("AMQP_URL")

	conn, err := amqp.Dial(url)
	if err != nil {
		LogOnError(err)
		return nil, nil, err
	}

	channel, err := conn.Channel()
	if err != nil {
		LogOnError(err)
		return nil, nil, err
	}

	err = channel.Qos(1, 0, false)
	if err != nil {
		LogOnError(err)
		return nil, nil, err
	}

	log.Printf("setup channel success!")

	return conn, channel, nil
}

func cloneToPublishMsg(msg *amqp.Delivery) *amqp.Publishing {
	newMsg := amqp.Publishing{
		Headers: msg.Headers,

		ContentType:     msg.ContentType,
		ContentEncoding: msg.ContentEncoding,
		DeliveryMode:    msg.DeliveryMode,
		Priority:        msg.Priority,
		CorrelationId:   msg.CorrelationId,
		ReplyTo:         msg.ReplyTo,
		Expiration:      msg.Expiration,
		MessageId:       msg.MessageId,
		Timestamp:       msg.Timestamp,
		Type:            msg.Type,
		UserId:          msg.UserId,
		AppId:           msg.AppId,

		Body: msg.Body,
	}

	return &newMsg
}

func newHttpClient(maxIdleConns, maxIdleConnsPerHost, idleConnTimeout int) *http.Client {
	tr := &http.Transport{
		MaxIdleConns:        maxIdleConns,
		MaxIdleConnsPerHost: maxIdleConnsPerHost,
		IdleConnTimeout:     time.Duration(idleConnTimeout) * time.Second,
	}

	client := &http.Client{
		Transport: tr,
	}

	return client
}

func notifyUrl(client *http.Client, url string, body []byte) int {
	req, err := http.NewRequest("POST", url, bytes.NewReader(body))
	if err != nil {
		log.Printf("notify url create req fail: %s", err)
		return 0
	}

	req.Header.Set("Content-Type", "application/json")
	response, err := client.Do(req)

	if err != nil {
		log.Printf("notify url %s fail: %s", url, err)
		return 0
	}
	defer response.Body.Close()

	io.Copy(ioutil.Discard, response.Body)

	return response.StatusCode
}
