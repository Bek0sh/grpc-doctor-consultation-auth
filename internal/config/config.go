package config

import (
	"sync"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	IsDebug *bool `yaml:"is_debug"`
	Listen  struct {
		Type   string `yaml:"type"`
		BindIp string `yaml:"bind_ip"`
		Port   string `yaml:"port"`
	} `yaml:"listen"`
	Postgre struct {
		DbUsername string `yaml:"db_username"`
		DbPassword string `yaml:"db_password"`
		DbPort     string `yaml:"db_port"`
		DbName     string `yaml:"db_name"`
	} `yaml:"postgre"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(
		func() {
			logrus.Info("Getting yaml configurations")
			instance = &Config{}
			if err := cleanenv.ReadConfig("config.yaml", instance); err != nil {
				logrus.Fatal(err)
			}
		},
	)
	return instance
}
