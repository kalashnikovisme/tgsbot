package bot

import (
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/telegramStandupBot/config"
	"github.com/maddevsio/telegramStandupBot/storage"
	log "github.com/sirupsen/logrus"
)

const (
	telegramAPIUpdateInterval = 60
)

// Bot structure
type Bot struct {
	c       *config.BotConfig
	tgAPI   *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
	db      *storage.MySQL
}

var yesterdayWorkKeywords = []string{"yesterday"}
var todayPlansKeywords = []string{"today"}
var issuesKeywords = []string{"block"}

// New creates a new bot instance
func New(c *config.BotConfig) (*Bot, error) {
	newBot, err := tgbotapi.NewBotAPI(c.TelegramToken)
	if err != nil {
		return nil, err
	}

	newBot.Debug = c.Debug

	u := tgbotapi.NewUpdate(0)

	u.Timeout = telegramAPIUpdateInterval

	updates, err := newBot.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}

	conn, err := storage.NewMySQL(c)
	if err != nil {
		return nil, err
	}

	b := &Bot{
		c:       c,
		tgAPI:   newBot,
		updates: updates,
		db:      conn,
	}

	return b, nil
}

// Start bot
func (b *Bot) Start() {
	b.StartNotificationThreads()
	log.Info("Listening for updates... \n")
	for update := range b.updates {
		b.handleUpdate(update)
	}
}
