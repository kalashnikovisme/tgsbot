package main

import (
	"github.com/maddevsio/telegramStandupBot/bot"
	"github.com/maddevsio/telegramStandupBot/config"
	log "github.com/sirupsen/logrus"
)

func main() {
	c, err := config.Get()
	if err != nil {
		log.Fatal(err)
	}
	b, err := bot.New(c)
	if err != nil {
		log.Fatal(err)
	}

	b.Start()
}
