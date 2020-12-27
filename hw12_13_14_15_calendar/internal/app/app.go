package app

import (
	"context"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage"
)

type App struct {
	Logger  Logger
	Storage Storage
}

type Logger interface {
	logger.Logger
}

type Storage interface {
	storage.EventStorage
}

func New(logger Logger, storage Storage) *App {
	return &App{Logger: logger, Storage: storage}
}

func (a *App) CreateEvent(ctx context.Context, id string, title string) error {
	// TODO

	return nil
	// return a.storage.CreateEvent(storage.Event{ID: id, Title: title})
}

// TODO
