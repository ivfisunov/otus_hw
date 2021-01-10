package main

import (
	"flag"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/amqp/publisher"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/config/config_scheduler"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/scheduler"
	internalgrpc "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/server/grpc"
	_ "github.com/lib/pq"
	"google.golang.org/grpc"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_scheduler.toml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := configscheduler.NewConfig(configFile)
	if err != nil {
		log.Fatal(err)
	}

	logg, err := logger.New(config.Logger.Level, config.Logger.Path)
	if err != nil {
		log.Fatal(err)
	}

	conn, err := grpc.Dial(net.JoinHostPort(config.Grpc.Host, config.Grpc.Port), grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	client := internalgrpc.NewCalendarClient(conn)

	publisher := amqppublisher.NewPublisher(
		logg,
		config.Amqp.URI,
		config.Amqp.Qname,
		config.Amqp.Exchname,
		config.Amqp.Exchtype)

	scheduler := scheduler.NewScheduler(client, logg, publisher)
	logg.Info("scheduler started...")

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	notifyEvents := time.NewTicker(time.Second * 5)
	deleteEvents := time.NewTicker(time.Hour * 24)

	go func() {
		for range notifyEvents.C {
			err := scheduler.FindEvents()
			if err != nil {
				logg.Error(err.Error())
			}
		}
	}()

	go func() {
		for range deleteEvents.C {
			err := scheduler.DeleteEvents()
			if err != nil {
				logg.Error(err.Error())
			}
		}
	}()

	<-signals
	logg.Info("scheduler shutdown")
}
