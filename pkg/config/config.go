package config

import (
	"fmt"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken     string
	PocketConsumerKey string
	AuthServerUrl     string
	TelegramBotUrl    string `mapstructure:"bot_url"`
	DbPath            string `mapstructure:"db_file"`
	Messages          Messages
}

type Messages struct {
	Errors
	Responses
}

type Errors struct {
	Default      string `mapstructure:"default"`
	Unauthorised string `mapstructure:"unauthorised"`
	InvalidUrl   string `mapstructure:"invalid_url"`
	UnableToSave string `mapstructure:"unable_to_save"`
}

type Responses struct {
	Start             string `mapstructure:"start"`
	AlreadyAuthorized string `mapstructure:"already_authorized"`
	SavedSuccessfully string `mapstructure:"saved_successfully"`
	UnknownCommand    string `mapstructure:"unknown_command"`
}

func Init() (*Config, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	if err := viper.ReadInConfig(); err != nil {
		return nil, err
	}

	var cfg Config
	if err := viper.Unmarshal(&cfg); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.responses", &cfg.Messages.Responses); err != nil {
		return nil, err
	}
	if err := viper.UnmarshalKey("messages.errors", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	if err := fromEnv(&cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func fromEnv(cfg *Config) error {
	viper.SetConfigName(".env")
	viper.AddConfigPath(".")
	viper.SetConfigType("env")

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("failed to read configuration: %s", err))
	}

	if err := viper.BindEnv("token"); err != nil {
		return err
	}
	cfg.TelegramToken = viper.GetString("token")

	if err := viper.BindEnv("consumer_key"); err != nil {
		return err
	}
	cfg.PocketConsumerKey = viper.GetString("consumer_key")

	if err := viper.BindEnv("auth_server_url"); err != nil {
		return err
	}
	cfg.AuthServerUrl = viper.GetString("auth_server_url")

	return nil
}
