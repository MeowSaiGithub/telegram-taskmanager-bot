package config

import (
	"github.com/go-playground/validator/v10"
	"github.com/spf13/viper"
	"log"
)

type Config struct {
	BotAPIKey string `json:"bot_api_key" validate:"required"`

	Host     string `json:"database_host" validate:"required"`
	User     string `json:"database_user" validate:"required"`
	Password string `json:"database_password" validate:"required"`
	Name     string `json:"database_name" validate:"required"`
	Port     int    `json:"database_port" validate:"required"`
	Sslmode  string `json:"database_sslmode" validate:"required"`
}

func GetConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("json")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			log.Println("there is no configuration file.")
		} else {
			log.Fatalf("failed to read configuration file %v\n", err)
		}
	}

	var conf Config

	conf.BotAPIKey = viper.GetString("bot_api_key")

	conf.Host = viper.GetString("database_host")
	conf.User = viper.GetString("database_user")
	conf.Password = viper.GetString("database_password")
	conf.Name = viper.GetString("database_name")
	conf.Port = viper.GetInt("database_port")
	conf.Sslmode = viper.GetString("database_sslmode")

	validate := validator.New()
	if err := validate.Struct(conf); err != nil {
		log.Fatalf("configuration validation error: %s\n", err)
	}
	return &conf
}
