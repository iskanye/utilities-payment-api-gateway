package config

import (
	"os"
	"time"
)

type Config struct {
	Host        string `yaml:"host" env-default:"localhost"`
	Port        int    `yaml:"port" env-default:"8080"`
	AuthSecret  string
	Timeout     time.Duration `yaml:"timeout"`
	Auth        HostPort      `yaml:"auth"`
	Billing     HostPort      `yaml:"billing"`
	Payment     HostPort      `yaml:"payment"`
	Memcached   HostPort      `yaml:"memcached"`
	BillingTerm int           `yaml:"billing_term" env-default:"1"` // in Months
}

type HostPort struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-required:"true"`
}

func (c *Config) MustLoadSecret() {
	c.AuthSecret = os.Getenv("AUTH_SECRET")
	if c.AuthSecret == "" {
		panic("auth secret mustnt be empty")
	}
}
