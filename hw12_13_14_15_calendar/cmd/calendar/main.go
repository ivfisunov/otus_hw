package main

import (
	"context"
	"errors"
	"flag"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/app"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/config"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	internalgrpc "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	internalhttp "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/server/http"
	createstorage "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/storage/create_storage"
	_ "github.com/lib/pq"
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

	config, err := config.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal(err)
	}

	storage, err := createstorage.Init(config.Storage.Type, config.Storage.Dsn)
	if err != nil {
		log.Fatal(err)
	}
	calendar := app.New(logg, storage)

	restServer := internalhttp.NewServer(calendar, config.Http.Host, config.Http.Port)
	grpcServer := internalgrpc.NewServer(calendar, config.Grpc.Host, config.Grpc.Port)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		signals := make(chan os.Signal, 1)
		signal.Notify(signals)

		<-signals
		logg.Info("Servers are stopping...")
		signal.Stop(signals)
		// cancel()

		ctx, cancel := context.WithTimeout(context.Background(), time.Second*3)
		defer cancel()

		grpcServer.Stop()
		if err := restServer.Stop(ctx); err != nil {
			logg.Error("failed to stop http server: " + err.Error())
		}
	}()

	logg.Info("calendar is running...")

	if err := grpcServer.Start(); err != nil {
		logg.Error("failed to start grpc server: " + err.Error())
	}
	if err := restServer.Start(ctx); err != nil && !errors.Is(err, http.ErrServerClosed) {
		logg.Error("failed to start http server: " + err.Error())
	}
}
