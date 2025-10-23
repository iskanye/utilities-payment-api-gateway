package config

import (
	"flag"
	"os"

	"github.com/ilyakaznacheev/cleanenv"
)

func MustLoad[T any](modifyCfg func(*T)) *T {
	configPath := fetchConfigPath()
	if configPath == "" {
		panic("config path is empty")
	}

	return MustLoadPath(configPath, modifyCfg)
}

func MustLoadPath[T any](configPath string, modifyCfg func(*T)) *T {
	// check if file exists
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		panic("config file does not exist: " + configPath)
	}

	var cfg T

	if err := cleanenv.ReadConfig(configPath, &cfg); err != nil {
		panic("cannot read config: " + err.Error())
	}

	modifyCfg(&cfg)

	return &cfg
}

func fetchConfigPath() string {
	var res string

	flag.StringVar(&res, "config", "", "path to config file")
	flag.Parse()

	if res == "" {
		res = os.Getenv("CONFIG_PATH")
	}

	return res
}
