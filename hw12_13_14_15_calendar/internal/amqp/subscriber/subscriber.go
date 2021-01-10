package amqpsubscriber

import (
	"encoding/json"

	amqppublisher "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/amqp/publisher"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/streadway/amqp"
)

type Subscriber struct {
	logg     *logger.Logg
	conn     *amqp.Connection
	channel  *amqp.Channel
	uri      string
	qName    string
	exchName string
	exchType string
}

func NewSubscriber(logger *logger.Logg, uri, qname, exchname, exchtype string) *Subscriber {
	return &Subscriber{
		logg:     logger,
		uri:      uri,
		qName:    qname,
		exchName: exchname,
		exchType: exchtype,
	}
}

func (s *Subscriber) Listen(processMessage func(*amqppublisher.Notification)) error {
	var err error
	s.conn, err = amqp.Dial(s.uri)
	if err != nil {
		return err
	}

	s.channel, err = s.conn.Channel()
	if err != nil {
		return err
	}

	if err := s.channel.ExchangeDeclare(
		s.exchName, // name
		s.exchType, // type
		true,       // durable
		false,      // auto-deleted
		false,      // internal
		false,      // noWait
		nil,        // arguments
	); err != nil {
		return err
	}

	_, err = s.channel.QueueDeclare(
		s.qName, // name
		false,   // durable
		false,   // delete when unused
		false,   // exclusive
		false,   // no-wait
		nil,     // arguments
	)
	if err != nil {
		return err
	}

	if err = s.channel.QueueBind(
		s.qName,
		s.qName,
		s.exchName,
		false,
		nil,
	); err != nil {
		return err
	}

	msgs, err := s.channel.Consume(
		s.qName,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for {
		for msg := range msgs {
			notification := &amqppublisher.Notification{}
			err := json.Unmarshal(msg.Body, notification)
			if err != nil {
				s.logg.Error("error reading queue")
			}

			processMessage(notification)
		}
	}
}
