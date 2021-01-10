package sender

import (
	"fmt"

	amqppublisher "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/amqp/publisher"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/amqp/subscriber"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
)

type Sender struct {
	logg      *logger.Logg
	subsciber *amqpsubscriber.Subscriber
}

func NewSender(logger *logger.Logg, sub *amqpsubscriber.Subscriber) *Sender {
	return &Sender{logger, sub}
}

func (s *Sender) ProcessMessages() error {
	err := s.subsciber.Listen(s.handleMessage)
	if err != nil {
		return err
	}
	return nil
}

func (s *Sender) handleMessage(msg *amqppublisher.Notification) {
	s.logg.Info(fmt.Sprintf("Received and sent message: %q", *msg))
}
