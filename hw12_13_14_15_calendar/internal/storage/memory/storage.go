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

func (s *Storage) CreateEvent(event storage.Event) error {
	s.mu.RLock()
	for _, e := range s.events {
		if e.Date.Equal(event.Date) {
			return storage.ErrDateBusy
		}
	}
	s.mu.RUnlock()

	s.mu.Lock()
	s.events[event.ID] = event
	s.mu.Unlock()
	return nil
}

func (s *Storage) UpdateEvent(id int, event storage.Event) error {
	s.mu.RLock()
	for _, e := range s.events {
		if e.Date.Equal(event.Date) {
			return storage.ErrDateBusy
		}
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

func (s *Storage) ListEventDay(date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event
	s.mu.RLock()
	for _, e := range s.events {
		if e.Date.After(date) && e.Date.Before(date.AddDate(0, 0, 1)) {
			events = append(events, &e)
		}
	}
	s.mu.RUnlock()
	return events, nil
}

func (s *Storage) ListEventWeek(date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event
	s.mu.RLock()
	for _, e := range s.events {
		if e.Date.After(date) && e.Date.Before(date.AddDate(0, 0, 7)) {
			events = append(events, &e)
		}
	}
	s.mu.RUnlock()
	return events, nil

}

func (s *Storage) ListEventMonth(date time.Time) ([]*storage.Event, error) {
	var events []*storage.Event
	s.mu.RLock()
	for _, e := range s.events {
		if e.Date.After(date) && e.Date.Before(date.AddDate(0, 1, 0)) {
			events = append(events, &e)
		}
	}
	s.mu.RUnlock()
	return events, nil

}
