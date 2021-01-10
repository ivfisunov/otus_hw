package amqppublisher

import (
	"encoding/json"
	"time"

	"github.com/streadway/amqp"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

type Publisher struct {
	logg     *logger.Logg
	conn     *amqp.Connection
	channel  *amqp.Channel
	uri      string
	qName    string
	exchName string
	exchType string
}

type Notification struct {
	ID     int
	Title  string
	Date   time.Time
	UserID int
}

func NewPublisher(logger *logger.Logg, uri, qname, exchname, exchtype string) *Publisher {
	return &Publisher{
		logg:     logger,
		uri:      uri,
		qName:    qname,
		exchName: exchname,
		exchType: exchtype,
	}
}

func (p *Publisher) Connect() error {
	var err error
	p.conn, err = amqp.Dial(p.uri)
	if err != nil {
		return err
	}

	p.channel, err = p.conn.Channel()
	if err != nil {
		return err
	}

	if err := p.channel.ExchangeDeclare(
		p.exchName, // name
		p.exchType, // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return err
	}

	_, err = p.channel.QueueDeclare(
		p.qName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	p.logg.Info("connected to broker")
	return nil
}

func (p *Publisher) Publish(notification Notification) error {
	body, err := json.Marshal(notification)
	if err != nil {
		return err
	}

	err = p.channel.Publish(p.exchName, p.qName, false, false, amqp.Publishing{
		ContentType: "application/json",
		Body:        body,
	})
	if err != nil {
		return err
	}
	return nil
}
