package createStorage

import (
	"fmt"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage"
	memorystorage "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage/memory"
	sqlstorage "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage/sql"
)

var mapping = map[string]func(dsn string) (storage.EventStorage, error){
	"memory": func(dsn string) (storage.EventStorage, error) {
		return memorystorage.New(dsn)
	},
	"sql": func(dsn string) (storage.EventStorage, error) {
		return sqlstorage.New(dsn)
	},
}

func Init(storageType string, dsn string) (storage.EventStorage, error) {
	if _, ok := mapping[storageType]; !ok {
		return nil, fmt.Errorf("invalid storage type: %v", storageType)
	}
	stor, err := mapping[storageType](dsn)
	if err != nil {
		return nil, err
	}
	return stor, nil
}
