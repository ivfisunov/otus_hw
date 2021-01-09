//go:generate protoc -I ../../../api/ EventService.proto --go_out=. --go-grpc_out=.

package internalgrpc

import (
	"context"
	"fmt"
	"net"

	"github.com/golang/protobuf/ptypes"
	"github.com/golang/protobuf/ptypes/empty"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"google.golang.org/grpc"
)

type Server struct {
	*app.App
	server *grpc.Server
	UnimplementedCalendarServer
	addr string
}

func NewServer(app *app.App, host string, port string) *Server {
	addr := net.JoinHostPort(host, port)
	return &Server{app, nil, UnimplementedCalendarServer{}, addr}
}

func (s *Server) Start() error {
	lsn, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}

	server := grpc.NewServer()
	RegisterCalendarServer(server, s.UnimplementedCalendarServer)
	s.server = server

	go func() {
		if err := server.Serve(lsn); err != nil {
			s.Logger.Error("error running server")
		}
	}()
	s.Logger.Info(fmt.Sprintf("grpc server started on %s", s.addr))
	return nil
}

func (s *Server) Stop() {
	s.server.GracefulStop()
	s.Logger.Info("grpc server is stoped")
}

func (s *Server) CreateEvent(ctx context.Context, e *CreateEventReq) (*empty.Empty, error) {
	event, err := parseEvent(e.Event)
	if err != nil {
		s.Logger.Error("event parsing error: " + err.Error())
	}
	err = s.Storage.CreateEvent(event)
	if err != nil {
		s.Logger.Error("event creating error: " + err.Error())
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *Server) UpdateEvent(ctx context.Context, e *UpdateEventReq) (*empty.Empty, error) {
	id := int(e.Id)
	event, err := parseEvent(e.Event)
	if err != nil {
		s.Logger.Error("event parsing error: " + err.Error())
	}
	err = s.Storage.UpdateEvent(id, event)
	if err != nil {
		s.Logger.Error("event updating error: " + err.Error())
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *Server) DeleteEvent(ctx context.Context, e *DeleteEventReq) (*empty.Empty, error) {
	id := int(e.Id)
	err := s.Storage.DeleteEvent(id)
	if err != nil {
		s.Logger.Error("event deleting error: " + err.Error())
		return nil, err
	}
	return &empty.Empty{}, nil
}

func (s *Server) ListEventDay(ctx context.Context, e *ListEventReq) (*ListEventRes, error) {
	date, err := ptypes.Timestamp(e.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	events, err := s.Storage.ListEventDay(date)
	if err != nil {
		return nil, fmt.Errorf("error searching events: %w", err)
	}

	evnt := make([]*Event, 0, len(events))
	for _, event := range events {
		ev, err := parseEventProto(event)
		if err != nil {
			return nil, err
		}
		evnt = append(evnt, ev)
	}

	return &ListEventRes{
		Events: evnt,
	}, nil
}

func (s *Server) ListEventWeek(ctx context.Context, e *ListEventReq) (*ListEventRes, error) {
	date, err := ptypes.Timestamp(e.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	events, err := s.Storage.ListEventWeek(date)
	if err != nil {
		return nil, fmt.Errorf("error searching events: %w", err)
	}

	evnt := make([]*Event, 0, len(events))
	for _, event := range events {
		ev, err := parseEventProto(event)
		if err != nil {
			return nil, err
		}
		evnt = append(evnt, ev)
	}

	return &ListEventRes{
		Events: evnt,
	}, nil
}

func (s *Server) ListEventMonth(ctx context.Context, e *ListEventReq) (*ListEventRes, error) {
	date, err := ptypes.Timestamp(e.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date: %w", err)
	}

	events, err := s.Storage.ListEventMonth(date)
	if err != nil {
		return nil, fmt.Errorf("error searching events: %w", err)
	}

	evnt := make([]*Event, 0, len(events))
	for _, event := range events {
		ev, err := parseEventProto(event)
		if err != nil {
			return nil, err
		}
		evnt = append(evnt, ev)
	}

	return &ListEventRes{
		Events: evnt,
	}, nil
}

func parseEvent(e *Event) (storage.Event, error) {
	date, err := ptypes.Timestamp(e.Date)
	if err != nil {
		return storage.Event{}, fmt.Errorf("error parsing event date: %w", err)
	}
	duration, err := ptypes.Duration(e.Duration)
	if err != nil {
		return storage.Event{}, fmt.Errorf("error parsing event duration: %w", err)
	}
	notifyBefore, err := ptypes.Duration(e.NotifyBefore)
	if err != nil {
		return storage.Event{}, fmt.Errorf("error parsing event notification: %w", err)
	}
	event := storage.Event{
		ID:           int(e.Id),
		Title:        e.Title,
		Date:         date,
		Duration:     duration,
		Description:  e.Description,
		UserID:       int(e.UserId),
		NotifyBefore: notifyBefore,
	}

	return event, nil
}

func parseEventProto(e storage.Event) (*Event, error) {
	date, err := ptypes.TimestampProto(e.Date)
	if err != nil {
		return nil, fmt.Errorf("error parsing date timestamp: %w", err)
	}
	duration := ptypes.DurationProto(e.Duration)
	notifyBefore := ptypes.DurationProto(e.NotifyBefore)

	return &Event{
		Id:           int32(e.ID),
		Title:        e.Title,
		Date:         date,
		Duration:     duration,
		Description:  e.Description,
		UserId:       int32(e.UserID),
		NotifyBefore: notifyBefore,
	}, nil
}
