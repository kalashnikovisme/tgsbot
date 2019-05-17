package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/tgstandupbot/config"
	"github.com/maddevsio/tgstandupbot/model"
	"github.com/maddevsio/tgstandupbot/storage"
	"github.com/sirupsen/logrus"
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

// New creates a new bot instance
func New(c *config.BotConfig) (*Bot, error) {
	newBot, err := tgbotapi.NewBotAPI(c.TelegramToken)
	if err != nil {
		return nil, err
	}
	b := &Bot{
		c:     c,
		tgAPI: newBot,
	}
	u := tgbotapi.NewUpdate(0)
	u.Timeout = telegramAPIUpdateInterval
	updates, err := b.tgAPI.GetUpdatesChan(u)
	if err != nil {
		return nil, err
	}
	conn, err := storage.NewMySQL(c)
	if err != nil {
		return nil, err
	}
	b.updates = updates
	b.db = conn
	return b, nil
}

// Start ...
func (b *Bot) Start() {
	logrus.Info("Starting tg bot\n")
	for update := range b.updates {
		b.handleUpdate(update)
	}
}

func (b *Bot) handleUpdate(update tgbotapi.Update) {

	text := update.Message.Text
	if !strings.Contains(text, "@"+b.tgAPI.Self.UserName) {
		return
	}

	chatID := update.Message.Chat.ID

	if isStandup(update.Message.Text) {
		standup := &model.Standup{
			Comment:  update.Message.Text,
			Username: update.Message.From.UserName,
		}
		_, err := b.db.CreateStandup(standup)
		if err != nil {
			logrus.Errorf("CreateStandup failed: %v\n", err)
			b.tgAPI.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("@%s у тебя кажется норм стендап, но сохранять его не буду.", update.Message.From.UserName)))
			return
		}
		b.tgAPI.Send(tgbotapi.NewMessage(chatID, fmt.Sprintf("@%s спасибо. Я принял твой стендап", update.Message.From.UserName)))
	}
}

func (b *Bot) senderIsAdminInChannel(sendername string, chatID int64) (bool, error) {
	chat := tgbotapi.ChatConfig{
		ChatID:             chatID,
		SuperGroupUsername: "",
	}
	admins, err := b.tgAPI.GetChatAdministrators(chat)
	if err != nil {
		return false, err
	}
	for _, admin := range admins {
		if admin.User.UserName == sendername {
			return true, nil
		}
	}
	return false, nil
}

func isStandup(message string) bool {
	logrus.Info("checking message...\n")
	var mentionsProblem, mentionsYesterdayWork, mentionsTodayPlans bool

	problemKeys := []string{"роблем", "рудност", "атруднен", "блок"}
	for _, problem := range problemKeys {
		if strings.Contains(message, problem) {
			mentionsProblem = true
		}
	}

	yesterdayWorkKeys := []string{"чера", "ятницу", "делал", "делано"}
	for _, work := range yesterdayWorkKeys {
		if strings.Contains(message, work) {
			mentionsYesterdayWork = true
		}
	}

	todayPlansKeys := []string{"егодн", "обираюс", "ланир"}
	for _, plan := range todayPlansKeys {
		if strings.Contains(message, plan) {
			mentionsTodayPlans = true
		}
	}
	return mentionsProblem && mentionsYesterdayWork && mentionsTodayPlans
}
