package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os/signal"
	"syscall"

	config "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/configs"
	"github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/logger"
	domain "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/model"
	rmqclient "github.com/fAv4JhG/home-works/hw12_13_14_15_calendar/internal/rmqclient"
)

var configFile string

func init() {
	flag.StringVar(&configFile, "config", "../../configs/config.yaml", "Path to configuration file")
}

func main() {
	flag.Parse()

	config, err := config.NewConfig(configFile)
	if err != nil {
		fmt.Println("Ошибка загрузки конфигурации:", err)
		return
	}
	logg := logger.New(config.Logger.Level)

	_, cancel := signal.NotifyContext(context.Background(),
		syscall.SIGINT, syscall.SIGTERM, syscall.SIGHUP)
	defer cancel()

	rmq, err := rmqclient.NewRabbitMQClient(config.Rabbit.DSN())
	if err != nil {
		logg.Fatal("Ошибка подключения к RabbitMQ :", err)
	}
	defer rmq.Close()

	queueName := "event_notifications"
	msgs, err := rmq.ConsumeMessages(queueName)
	if err != nil {
		logg.Fatal("Ошибка при подписке на очередь: ", err)
	}

	for msg := range msgs {
		handleMessage(msg.Body, logg)
	}
}

func handleMessage(body []byte, logg *logger.Logger) {
	var notification domain.Notification
	if err := json.Unmarshal(body, &notification); err != nil {
		logg.Error("Ошибка парсинга сообщения: ", err)
		return
	}
	logg.Info("Отправка уведомления:", notification)
}
