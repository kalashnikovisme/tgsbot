package config

import "github.com/kelseyhightower/envconfig"

// BotConfig ...
type BotConfig struct {
	TelegramToken string `envconfig:"TELEGRAM_TOKEN" required:"false"`
	DatabaseURL   string `envconfig:"DATABASE_URL" default:"tgbot:tgbot@tcp(db:3306)/tgbot?parseTime=true" required:"true"`
	Debug         bool   `envconfig:"DEBUG" default:"false"`
}

// Get config data from environment
func Get() (*BotConfig, error) {
	var c BotConfig
	err := envconfig.Process("", &c)
	return &c, err
}
