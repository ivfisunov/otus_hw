package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"syscall"

	amqpsubscriber "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/amqp/subscriber"
	configscheduler "github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/config/config_scheduler"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/logger"
	"github.com/ivfisunov/otus_hw/hw12_13_14_15_calendar/internal/sender"
	_ "github.com/lib/pq"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "configs/config_sender.toml", "Path to configuration file")
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

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGTERM, syscall.SIGINT)

	subscriber := amqpsubscriber.NewSubscriber(
		logg,
		config.AMQP.URI,
		config.AMQP.Qname,
		config.AMQP.Exchname,
		config.AMQP.Exchtype)

	sender := sender.NewSender(logg, subscriber)
	logg.Info("sender started...")
	err = sender.ProcessMessages()
	if err != nil {
		logg.Error(err.Error())
		return
	}

	<-signals
	logg.Info("sender shutdown")
}
