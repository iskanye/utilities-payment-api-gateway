package config

import "time"

type Config struct {
	Host       string        `yaml:"host" env-default:"localhost"`
	Port       int           `yaml:"port" env-default:"8080"`
	AuthSecret string        `yaml:"secret"`
	Timeout    time.Duration `yaml:"timeout"`
	Auth       HostPort      `yaml:"auth"`
	Billing    HostPort      `yaml:"billing"`
}

type HostPort struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-required:"true"`
}
