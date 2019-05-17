package bot

import (
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/tgstandupbot/config"
	"github.com/maddevsio/tgstandupbot/storage"
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

var yesterdayWorkKeywords = []string{"вчера"}
var todayPlansKeywords = []string{"сегодня"}
var issuesKeywords = []string{"мешает"}

// New creates a new bot instance
func New(c *config.BotConfig) (*Bot, error) {
	newBot, err := tgbotapi.NewBotAPI(c.TelegramToken)
	if err != nil {
		return nil, err
	}

	newBot.Debug = false

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
	log.Info("Starting tg bot\n")
	for update := range b.updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {

	if update.Message.Text != "" {
		log.Info("The event contains text!\n")
	}

	if update.Message.LeftChatMember != nil {
		log.Info("Chat member left channel!\n")
	}

	if update.Message.NewChatMembers != nil {
		log.Info("Chat member joined!\n")
	}

	//? need to handle user change username
	//? need to handle slash commands
	/*
		? assign/view/unassign users to standup
		? set/view/remove standup deadline
	*/
}

func isStandup(message string) bool {
	message = strings.ToLower(message)

	var mentionsYesterdayWork, mentionsTodayPlans, mentionsProblem bool

	for _, work := range yesterdayWorkKeywords {
		if strings.Contains(message, work) {
			mentionsYesterdayWork = true
		}
	}

	for _, plan := range todayPlansKeywords {
		if strings.Contains(message, plan) {
			mentionsTodayPlans = true
		}
	}

	for _, problem := range issuesKeywords {
		if strings.Contains(message, problem) {
			mentionsProblem = true
		}
	}

	return mentionsProblem && mentionsYesterdayWork && mentionsTodayPlans
}
