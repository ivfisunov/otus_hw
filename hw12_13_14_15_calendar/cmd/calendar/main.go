package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalhttp "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	createStorage "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage/create-storage"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	if flag.Arg(0) == "version" {
		printVersion()
		return
	}

	config, err := NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := createStorage.Init(config.Storage.Type, config.Storage.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	calendar := app.New(logg, storage)

	server := internalhttp.NewServer(calendar, config.Server.Host, config.Server.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		fmt.Println("\nServer is stopping...")
		signal.Stop(signals)
		cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		if err := server.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := server.Start(ctx); err != nil && err != http.ErrServerClosed {
		logg.Error("failed to start http server: " + err.Error())
		os.Exit(1)
	}
}
