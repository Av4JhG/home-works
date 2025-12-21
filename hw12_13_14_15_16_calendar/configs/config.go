package config

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Logger    LoggerConf
	Database  DatabaseConf
	Server    ServerConf
	Storage   StorageConf
	Rabbit    RabbitConf
	Scheduler SchedulerConf
}

type LoggerConf struct {
	Level string
	File  string
}

type DatabaseConf struct {
	User     string
	Password string
	Name     string
	Host     string
	Port     string
	SSL      string
}

type ServerConf struct {
	Host string
	Port string
}

type StorageConf struct {
	Type string
}

type RabbitConf struct {
	Host     string
	Port     string
	User     string
	Password string
}
type SchedulerConf struct {
	Interval time.Duration
	Life     string
}

func NewConfig(configFile string) (*Config, error) {
	viper.SetConfigFile(configFile)
	if err := viper.ReadInConfig(); err != nil {
		return nil, fmt.Errorf("не удалось прочитать конфиг: %w", err)
	}
	var config Config

	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("не удалось распарсить конфиг: %w", err)
	}

	return &config, nil
}

func (d *DatabaseConf) DSN() string {
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		d.User, d.Password, d.Host, d.Port, d.Name, d.SSL)
}

func (s *ServerConf) DSN() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

func (r *RabbitConf) DSN() string {
	return fmt.Sprintf("amqp://%s:%s@%s:%s",
		r.User, r.Password, r.Host, r.Port)
}
