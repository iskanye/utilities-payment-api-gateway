package config

type Config struct {
	Host string   `yaml:"host" env-default:"localhost"`
	Port int      `yaml:"port" env-default:"8080"`
	Auth HostPort `yaml:"auth"`
}

type HostPort struct {
	Host string `yaml:"host" env-default:"localhost"`
	Port int    `yaml:"port" env-required:"true"`
}
