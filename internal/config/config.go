package config

import "time"

type Config struct {
	Host      string        `yaml:"host" env-default:"localhost"`
	Port      int           `yaml:"port" env-default:"8080"`
	CookieTTL time.Duration `yaml:"cookie_ttl" env-default:"1h"`
	Auth      HostPort      `yaml:"auth"`
}

type HostPort struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-required:"true"`
}
