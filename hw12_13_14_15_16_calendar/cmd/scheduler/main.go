package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"sync"
	"syscall"
	"time"

	config "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/configs"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/app"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
	rmqclient "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/rmqclient"
	sqlstorage "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/storage/sql"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

// go run ./cmd/scheduler/ -config ./configs/config.yaml.
func main() {
	flag.Parse()

	config, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Println(configFile)
		fmt.Println(err)
		return
	}
	logg := logger.New(config.Logger.Level)

	ctx, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rmq, err := rmqclient.NewRabbitMQClient(config.Rabbit.DSN())
	if err != nil {
		logg.Error("RabbitMQ connection error: ", err)
	}
	defer rmq.Close()

	storage := sqlstorage.New(config.Database.DSN())
	err = storage.Connect(ctx)
	if err != nil {
		logg.Error(err)
	}

	queueName := "event_notifications"
	_, err = rmq.DeclareQueue(queueName)
	if err != nil {
		logg.Error("RMQ queue creation error: ", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		for {
			select {
			case <-ctx.Done():
				logg.Info("scheduler stopped")
				return
			default:
				processEvents(ctx, rmq, queueName, logg, storage)
				cleanupOldEvents(ctx, storage, config.Scheduler.Life)
				time.Sleep(config.Scheduler.Interval)
			}
		}
	}()

	wg.Wait()
}

func processEvents(
	ctx context.Context,
	rmq *rmqclient.RabbitMQClient,
	queueName string,
	logg *logger.Logger,
	storage app.Storage,
) {
	events, err := storage.EventList(ctx)
	if err != nil {
		logg.Error("Get events error: ", err)
	}

	for _, event := range events {
		notification := domain.Notification{
			EventID: event.ID,
			Title:   event.Title,
			Time:    event.Date.Format(time.RFC3339),
		}

		body, err := json.Marshal(notification)
		if err != nil {
			logg.Error("JSON marshal error: ", err)
		}

		err = rmq.PublishMessage(queueName, body)
		if err != nil {
			logg.Error("Rabbit send notify error: ", err)
		} else {
			logg.Info("!!! Sent: ", notification)
		}
	}
}

func cleanupOldEvents(ctx context.Context, storage app.Storage, life string) {
	storage.ClearEvents(ctx, life)
}
