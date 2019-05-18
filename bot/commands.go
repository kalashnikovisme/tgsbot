package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/telegramStandupBot/model"
)

//HandleCommand handles imcomming commands
func (b *Bot) HandleCommand(event tgbotapi.Update) error {
	switch event.Message.Command() {
	case "help":
		return b.Help(event)
	case "join_standupers":
		return b.JoinStandupers(event)
	case "show":
		return b.Show(event)
	case "leave_standupers":
		return b.LeaveStandupers(event)
	case "edit_deadline":
		return b.EditDeadline(event)
	case "show_deadline":
		return b.ShowDeadline(event)
	case "remove_deadline":
		return b.RemoveDeadline(event)
	default:
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "I do not know this command...")
		_, err := b.tgAPI.Send(msg)
		return err
	}
}

//Help displays help message
func (b *Bot) Help(event tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "I am here to help")
	_, err := b.tgAPI.Send(msg)
	return err
}

//JoinStandupers assign user a standuper role
func (b *Bot) JoinStandupers(event tgbotapi.Update) error {
	_, err := b.db.FindStanduper(event.Message.From.UserName, event.Message.Chat.ID) // user[1:] to remove leading @
	if err == nil {
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "You are already in the standup team!")
		msg.ReplyToMessageID = event.Message.MessageID
		_, err := b.tgAPI.Send(msg)
		return err
	}

	_, err = b.db.CreateStanduper(&model.Standuper{
		Username:     event.Message.From.UserName,
		ChatID:       event.Message.Chat.ID,
		LanguageCode: event.Message.From.LanguageCode,
	})
	if err != nil {
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Could not add you to standup team")
		msg.ReplyToMessageID = event.Message.MessageID
		_, err := b.tgAPI.Send(msg)
		return err
	}

	group, err := b.db.FindGroup(event.Message.Chat.ID)
	if err != nil {
		//! Need to create the group here... this might happen
		//! if the bot was added to the group when being inactive
	}

	var msg tgbotapi.MessageConfig

	if group.StandupDeadline == "" {
		msg = tgbotapi.NewMessage(event.Message.Chat.ID, "Welcome to standup team! No deadlines for standup submittions in the team yet!")
	} else {
		msg = tgbotapi.NewMessage(event.Message.Chat.ID, fmt.Sprintf("Welcome to standup team! Please submit your standups till %s each day exept weekends!", group.StandupDeadline))
	}

	msg.ReplyToMessageID = event.Message.MessageID
	_, err = b.tgAPI.Send(msg)
	return err
}

//Show standupers
func (b *Bot) Show(event tgbotapi.Update) error {
	standupers, err := b.db.ListChatStandupers(event.Message.Chat.ID)
	if err != nil {
		return err
	}

	if len(standupers) == 0 {
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Currently no standupser in the group. To add, please use `/add` command")
		_, err := b.tgAPI.Send(msg)
		return err
	}

	list := []string{}
	for _, standuper := range standupers {
		list = append(list, "@"+standuper.Username)
	}

	msg := tgbotapi.NewMessage(event.Message.Chat.ID, fmt.Sprintf("Standupers in the group: %v", strings.Join(list, ", ")))
	_, err = b.tgAPI.Send(msg)
	return err
}

//LeaveStandupers standupers
func (b *Bot) LeaveStandupers(event tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Removing...")
	_, err := b.tgAPI.Send(msg)
	return err
}

//EditDeadline modifies standup time
func (b *Bot) EditDeadline(event tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Changing stanup time...")
	_, err := b.tgAPI.Send(msg)
	return err
}

//ShowDeadline shows current standup time
func (b *Bot) ShowDeadline(event tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Showing standup time...")
	_, err := b.tgAPI.Send(msg)
	return err
}

//RemoveDeadline sets standup deadline to empty string
func (b *Bot) RemoveDeadline(event tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Removing standup time...")
	_, err := b.tgAPI.Send(msg)
	return err
}
