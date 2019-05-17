package bot

import (
	"fmt"
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
		err := b.HandleMessageEvent(update)
		if err != nil {
			log.Error("Failed to Handle Message Event! ", err)
		}
	}

	if update.Message.LeftChatMember != nil {
		err := b.HandleChannelLeftEvent(update)
		if err != nil {
			log.Error("Failed to Handle Channel left Event! ", err)
		}
	}

	if update.Message.NewChatMembers != nil {
		err := b.HandleChannelJoinEvent(update)
		if err != nil {
			log.Error("Failed to Handle Channel Join Event! ", err)
		}
	}

	//? need to handle user change username
	//? need to handle slash commands
	/*
		? assign/view/unassign users to standup
		? set/view/remove standup deadline
	*/
}

//HandleMessageEvent function to analyze and save standups
func (b *Bot) HandleMessageEvent(event tgbotapi.Update) error {

	if !strings.Contains(event.Message.Text, b.tgAPI.Self.UserName) {
		return nil
	}

	if !isStandup(event.Message.Text) {
		return fmt.Errorf("Message is not a standup")
	}

	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Спасибо, стендап принят!")
	msg.ReplyToMessageID = event.Message.MessageID
	_, err := b.tgAPI.Send(msg)
	return err
}

//HandleChannelLeftEvent function to remove bot and standupers from channels
func (b *Bot) HandleChannelLeftEvent(event tgbotapi.Update) error {
	log.Info("Chat member left!\n")
	return nil
}

//HandleChannelJoinEvent function to add bot and standupers t0 channels
func (b *Bot) HandleChannelJoinEvent(event tgbotapi.Update) error {
	log.Info("Chat member joined!\n")
	text := "Hello! Nice to meet you all! I am here to help you with standups :}"
	_, err := b.tgAPI.Send(tgbotapi.NewMessage(event.Message.Chat.ID, text))
	return err
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
