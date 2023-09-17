package config

import (
	"errors"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

type APIConfig struct {
	Env string `yaml:"env"`
}

type SchedulerConfig struct {
	Env string `yaml:"env"`
}

type EmailSenderConfig struct {
	Env      string   `yaml:"env"`
	RabbitMQ rabbitmq `yaml:"rabbitmq"`
	SMTP     smtp     `yaml:"smtp"`
}

type rabbitmq struct {
	URL       string `yaml:"url" env:"RABBITMQ_URL"`
	QueueName string `yaml:"queue_name"`
}

type smtp struct {
	Host     string `yaml:"host" env:"SMTP_HOST"`
	Port     string `yaml:"port" env:"SMTP_PORT"`
	Email    string `env:"SMTP_EMAIL"`
	Password string `env:"SMTP_PASS"`
}

func LoadEmailSenderConfig() (*EmailSenderConfig, error) {
	path := os.Getenv("EMAIL_SENDER_CONFIG_PATH")
	if path == "" {
		return nil, errors.New("EMAIL_SENDER_CONFIG_PATH is not set")
	}

	if _, err := os.Stat(path); err != nil {
		return nil, err
	}

	var cfg EmailSenderConfig
	if err := cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}
