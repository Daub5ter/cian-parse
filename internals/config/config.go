package config

import (
	"github.com/spf13/viper"
)

type Config struct {
	Listen  `yaml:"listen"`
	Storage `yaml:"storage"`
}

type Listen struct {
	Type   string `yaml:"type"`
	BindIP string `yaml:"bind_ip"`
	Port   string `yaml:"port"`
}

type Storage struct {
	Host       string `json:"host"`
	Port       string `json:"port"`
	Database   string `json:"database"`
	Username   string `json:"username"`
	Password   string `json:"password"`
	AuthDB     string `json:"auth_db"`
	Collection string `json:"collection"`
}

func SetupConfig() *Config { //заполнение структуры
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")

	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	return &Config{
		Listen: Listen{
			Type: viper.GetString("listen.type"),
			Port: viper.GetString("listen.port"),
		},
		Storage: Storage{
			Host:       viper.GetString("storage.host"),
			Port:       viper.GetString("storage.port"),
			Database:   viper.GetString("storage.database"),
			Username:   viper.GetString("storage.username"),
			Password:   viper.GetString("storage.password"),
			AuthDB:     viper.GetString("storage.auth_db"),
			Collection: viper.GetString("storage.collection"),
		},
	}
}
