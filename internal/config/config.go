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
	Env      string   `yaml:"env"`
	Schedule string   `yaml:"schedule"`
	RabbitMQ rabbitmq `yaml:"rabbitmq"`
	Postgres postgres `yaml:"postgres"`
}

type EmailSenderConfig struct {
	Env      string   `yaml:"env"`
	RabbitMQ rabbitmq `yaml:"rabbitmq"`
	SMTP     smtp     `yaml:"smtp"`
}

type postgres struct {
	URL string `yaml:"url" env:"PG_URL"`
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
	path, err := configPath("EMAIL_SENDER_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	var cfg EmailSenderConfig
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func LoadSchedulerConfig() (*SchedulerConfig, error) {
	path, err := configPath("STATISTIC_CONFIG_PATH")
	if err != nil {
		return nil, err
	}

	var cfg SchedulerConfig
	if err = cleanenv.ReadConfig(path, &cfg); err != nil {
		return nil, err
	}

	return &cfg, err
}

func configPath(env string) (string, error) {
	path := os.Getenv(env)
	if path == "" {
		return "", errors.New("path to config file not set")
	}

	if _, err := os.Stat(path); err != nil {
		return "", err
	}

	return path, nil
}
