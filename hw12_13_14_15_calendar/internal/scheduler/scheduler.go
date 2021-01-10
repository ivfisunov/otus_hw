package scheduler

import (
	"context"
	"fmt"
	"time"

	"github.com/golang/protobuf/ptypes"
	amqppublisher "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/amqp/publisher"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
)

type Scheduler struct {
	client    internalgrpc.CalendarClient
	logger    *logger.Logg
	publisher *amqppublisher.Publisher
}

func NewScheduler(
	client internalgrpc.CalendarClient,
	logger *logger.Logg,
	publisher *amqppublisher.Publisher) *Scheduler {
	return &Scheduler{client, logger, publisher}
}

func (s *Scheduler) FindEvents() error {
	ctx := context.Background()
	eventReqProto := &internalgrpc.ListEventReq{Date: ptypes.TimestampNow()}
	eventsProto, err := s.client.ListEventDay(ctx, eventReqProto)
	if err != nil {
		return fmt.Errorf("error finding events: %w", err)
	}

	if eventsProto == nil {
		s.logger.Info("there are no events to notify")
		return nil
	}

	nowTime := time.Now()
	for _, event := range eventsProto.Events {
		notifyBefore, err := ptypes.Duration(event.NotifyBefore)
		if err != nil {
			return err
		}
		date, err := ptypes.Timestamp(event.Date)
		if err != nil {
			return err
		}
		ev := amqppublisher.Notification{
			ID:     int(event.Id),
			Title:  event.Title,
			Date:   date,
			UserID: int(event.UserId),
		}

		notify := date.Add(-notifyBefore)
		if nowTime.Equal(notify) || nowTime.After(notify) {
			err := s.publisher.Publish(ev)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Scheduler) DeleteEvents() error {
	ctx := context.Background()
	yearPast := time.Now().AddDate(-1, 0, -1)
	yearPastProto, err := ptypes.TimestampProto(yearPast)
	if err != nil {
		return err
	}

	eventReqProto := &internalgrpc.ListEventReq{Date: yearPastProto}
	events, err := s.client.ListEventDay(ctx, eventReqProto)
	if err != nil {
		return fmt.Errorf("error finding events: %w", err)
	}

	for _, event := range events.Events {
		_, err := s.client.DeleteEvent(ctx, &internalgrpc.DeleteEventReq{
			Id: event.Id,
		})
		if err != nil {
			s.logger.Error("error deleting event: " + err.Error())
		}
	}
	return nil
}
