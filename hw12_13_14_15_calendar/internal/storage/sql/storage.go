package sqlstorage

import (
	"fmt"
	"time"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	"github.com/jmoiron/sqlx"
)

type Storage struct {
	db *sqlx.DB
}

func New(dsn string) (*Storage, error) {
	// setup postgres connection
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("error connecting to Postgres: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping db error: %w", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) CreateEvent(event storage.Event) error {
	_, err := s.db.NamedExec("INSER INTO events (id, title, date, duration, description, user_id, notify_before) VALUES (:id, :title, :date, :duration, :description, :user_id, :notify_before)", event)
	if err != nil {
		return fmt.Errorf("db error creating new event, %w", err)
	}
	return nil
}

func (s *Storage) UpdateEvent(id int, event storage.Event) error {
	_, err := s.db.NamedExec("UPDATE events SET title=:title, date=:date, duration=:duration, description=:description, user_pd=:user_id, notify_before=:notify_before WHERE id = :id",
		&storage.Event{
			ID:           event.ID,
			Title:        event.Title,
			Date:         event.Date,
			Duration:     event.Duration,
			Description:  event.Description,
			UserID:       event.UserID,
			NotifyBefore: event.NotifyBefore})
	if err != nil {
		return fmt.Errorf("db error updating event, %w", err)
	}
	return nil
}

func (s *Storage) DeleteEvent(id int) error {
	_, err := s.db.NamedExec("DELETE FROM events WHERE id = $id", id)
	if err != nil {
		return fmt.Errorf("db errro deleting event %w", err)
	}
	return nil
}

func (s *Storage) ListEventDay(date time.Time) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.Select(&events, "SELECT * FROM events WHERE date BETWEEN $1 AND $1 + interval '1 day'", date)
	if err != nil {
		return nil, fmt.Errorf("db error selecting events, %w", err)
	}
	return events, nil
}

func (s *Storage) ListEventWeek(date time.Time) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.Select(&events, "SELECT * FROM events WHERE date BETWEEN $1 AND $1 + interval '7 days'", date)
	if err != nil {
		return nil, fmt.Errorf("db error selecting events, %w", err)
	}
	return events, nil
}

func (s *Storage) ListEventMonth(date time.Time) ([]storage.Event, error) {
	var events []storage.Event
	err := s.db.Select(&events, "SELECT * FROM events WHERE date BETWEEN $1 AND $1 + interval '1 month'", date)
	if err != nil {
		return nil, fmt.Errorf("db error selecting events, %w", err)
	}
	return events, nil
}
