package memorystorage

import (
	"sync"
	"time"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type Storage struct {
	mu     sync.RWMutex
	events map[int]storage.Event
}

func New(_ string) (*Storage, error) {
	return &Storage{
		mu:     sync.RWMutex{},
		events: make(map[int]storage.Event),
	}, nil
}

func checkBusyEvent(events map[int]storage.Event, event storage.Event) error {
	for _, e := range events {
		if e.Date.Equal(event.Date) {
			return storage.ErrDateBusy
		}
	}
	return nil
}

func findEvents(events map[int]storage.Event, date time.Time, days int) []storage.Event {
	evts := make([]storage.Event, 0)
	duration := time.Duration(days * 24)
	for _, e := range events {
		if e.Date.After(date) && e.Date.Before(date.Add(time.Hour*duration)) {
			evts = append(evts, e)
		}
	}
	return evts
}

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.RLock()
	err := checkBusyEvent(s.events, event)
	if err != nil {
		return err
	}
	s.mu.RUnlock()

	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) UpdateEvent(id int, event storage.Event) error {
	s.mu.RLock()
	err := checkBusyEvent(s.events, event)
	if err != nil {
		return err
	}
	s.mu.RUnlock()

	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	s.mu.Lock()
	delete(s.events, id)
	s.mu.Unlock()
	return nil
}

func (s *Storage) ListEventDay(date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	events := findEvents(s.events, date, 1)
	s.mu.RUnlock()
	return events, nil
}

func (s *Storage) ListEventWeek(date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	events := findEvents(s.events, date, 7)
	s.mu.RUnlock()
	return events, nil
}

func (s *Storage) ListEventMonth(date time.Time) ([]storage.Event, error) {
	s.mu.RLock()
	events := findEvents(s.events, date, 30)
	s.mu.RUnlock()
	return events, nil
}
