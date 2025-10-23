package config

type Config struct {
	Auth HostPort `yaml:"auth"`
}

type HostPort struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-required:"true"`
}
