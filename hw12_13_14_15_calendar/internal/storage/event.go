package storage

import (
	"errors"
	"time"
)

type Event struct {
	ID           int           `db:"id"`
	Title        string        `db:"title"`
	Date         time.Time     `db:"date"`
	Duration     time.Duration `db:"duration"`
	Description  string        `db:"description"`
	UserID       int           `db:"user_id"`
	NotifyBefore time.Duration `db:"notify_before"`
}

type Notification struct {
	ID     int
	Title  string
	Date   time.Time
	ToUser int
}

type EventStorage interface {
	CreateEvent(Event) error
	UpdateEvent(int, Event) error
	DeleteEvent(int) error
	ListEventDay(time.Time) ([]Event, error)
	ListEventWeek(time.Time) ([]Event, error)
	ListEventMonth(time.Time) ([]Event, error)
}

var ErrDateBusy = errors.New("this time is busy")
