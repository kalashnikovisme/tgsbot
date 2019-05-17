package main

import (
	"log"

	"github.com/maddevsio/telegramStandupBot/bot"
	"github.com/maddevsio/telegramStandupBot/config"
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
