package bot

import (
	"fmt"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/telegramStandupBot/model"
	log "github.com/sirupsen/logrus"
)

//HandleCommand handles imcomming commands
func (b *Bot) HandleCommand(event tgbotapi.Update) error {
	switch event.Message.Command() {
	case "help":
		return b.Help(event)
	case "join":
		return b.JoinStandupers(event)
	case "show":
		return b.Show(event)
	case "leave":
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
	text := ` Here is the list of available commands:
	/help - Display list of available commands
	/join - Adds you to standup team of the group
	/show - Shows who submit standups
	/leave - Removes you from standup team of the group
	/edit_deadline - Sets new standup deadline (you can use 10am format or 15:30 format)
	/show_deadline - Shows current standup deadline 
	/remove_deadline - Removes standup deadline at all

	Looking forward for your standups!
	`
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, text)
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
		log.Error("CreateStanduper failed: ", err)
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
	standuper, err := b.db.FindStanduper(event.Message.From.UserName, event.Message.Chat.ID) // user[1:] to remove leading @
	if err != nil {
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "You are not in the standup team yet!")
		msg.ReplyToMessageID = event.Message.MessageID
		_, err := b.tgAPI.Send(msg)
		return err
	}

	err = b.db.DeleteStanduper(standuper.ID)
	if err != nil {
		log.Error("DeleteStanduper failed: ", err)
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Could not remove you from standup team")
		msg.ReplyToMessageID = event.Message.MessageID
		_, err := b.tgAPI.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "You are no longer in stanup team, thank you for all your standups!")
	msg.ReplyToMessageID = event.Message.MessageID
	_, err = b.tgAPI.Send(msg)
	return err
}

//EditDeadline modifies standup time
func (b *Bot) EditDeadline(event tgbotapi.Update) error {
	deadline := event.Message.CommandArguments()

	group, err := b.db.FindGroup(event.Message.Chat.ID)
	if err != nil {
		//! Need to create the group here... this might happen
		//! if the bot was added to the group when being inactive
	}

	group.StandupDeadline = deadline

	_, err = b.db.UpdateGroup(group)
	if err != nil {
		log.Error("UpdateGroup in EditDeadline failed: ", err)
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Could not update deadline")
		msg.ReplyToMessageID = event.Message.MessageID
		_, err = b.tgAPI.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(event.Message.Chat.ID, fmt.Sprintf("Deadline updated! New deadline is %s", deadline))
	msg.ReplyToMessageID = event.Message.MessageID
	_, err = b.tgAPI.Send(msg)
	return err
}

//ShowDeadline shows current standup time
func (b *Bot) ShowDeadline(event tgbotapi.Update) error {
	group, err := b.db.FindGroup(event.Message.Chat.ID)
	if err != nil {
		//! Need to create the group here... this might happen
		//! if the bot was added to the group when being inactive
	}

	var msg tgbotapi.MessageConfig

	if group.StandupDeadline == "" {
		msg = tgbotapi.NewMessage(event.Message.Chat.ID, "No deadlines for standup submittions in the team yet!")
	} else {
		msg = tgbotapi.NewMessage(event.Message.Chat.ID, fmt.Sprintf("Deadline is %s each day exept weekends!", group.StandupDeadline))
	}

	msg.ReplyToMessageID = event.Message.MessageID
	_, err = b.tgAPI.Send(msg)
	return err
}

//RemoveDeadline sets standup deadline to empty string
func (b *Bot) RemoveDeadline(event tgbotapi.Update) error {
	group, err := b.db.FindGroup(event.Message.Chat.ID)
	if err != nil {
		//! Need to create the group here... this might happen
		//! if the bot was added to the group when being inactive
	}

	group.StandupDeadline = ""

	_, err = b.db.UpdateGroup(group)
	if err != nil {
		log.Error("UpdateGroup in RemoveDeadline failed: ", err)
		msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Could not remove deadline")
		msg.ReplyToMessageID = event.Message.MessageID
		_, err = b.tgAPI.Send(msg)
		return err
	}

	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Deadline removed!")
	msg.ReplyToMessageID = event.Message.MessageID
	_, err = b.tgAPI.Send(msg)
	return err
}
