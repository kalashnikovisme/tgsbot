package main

import (
	"log"

	"github.com/maddevsio/tgstandupbot/bot"
	"github.com/maddevsio/tgstandupbot/config"
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
