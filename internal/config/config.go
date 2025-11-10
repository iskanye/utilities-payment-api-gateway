package config

import (
	"os"
	"time"
)

type Config struct {
	Host        string        `yaml:"host" env-default:"localhost"`
	Port        int           `yaml:"port" env-default:"8080"`
	AuthSecret  string        `yaml:"secret"`
	Timeout     time.Duration `yaml:"timeout"`
	Auth        HostPort      `yaml:"auth"`
	Billing     HostPort      `yaml:"billing"`
	Payment     HostPort      `yaml:"payment"`
	BillingTerm int           `yaml:"billing_term" env-default:"1"` // in Months
}

type HostPort struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-required:"true"`
}

func (c *Config) LoadSecret() {
	// Если секрет не задан в конфиге ищем его в параметрах окружения
	if c.AuthSecret == "" {
		c.AuthSecret = os.Getenv("AUTH_SECRET")
	}
	if c.AuthSecret == "" {
		panic("auth secret mustnt be empty")
	}
}
