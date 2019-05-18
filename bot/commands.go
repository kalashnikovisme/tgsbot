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
	case "add":
		return b.Add(event)
	case "show":
		return b.Show(event)
	case "remove":
		return b.Remove(event)
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

//Add assign user a standuper role
func (b *Bot) Add(event tgbotapi.Update) error {
	added := []string{}
	exist := []string{}
	failed := []string{}

	toAdd := event.Message.CommandArguments()
	users := strings.Split(toAdd, " ")
	for _, user := range users {
		_, err := b.db.FindStanduper(user[1:], event.Message.Chat.ID) // user[1:] to remove leading @
		if err == nil {
			exist = append(exist, user)
			continue
		}
		chatMember, err := b.tgAPI.GetChatMember(tgbotapi.ChatConfigWithUser{
			ChatID:             event.Message.Chat.ID,
			SuperGroupUsername: user[1:],
		})
		if err != nil {
			failed = append(failed, user)
			continue
		}
		_, err = b.db.CreateStanduper(&model.Standuper{
			Username:     chatMember.User.UserName,
			ChatID:       event.Message.Chat.ID,
			LanguageCode: chatMember.User.LanguageCode,
		})
		if err != nil {
			failed = append(failed, user)
			continue
		}
		added = append(added, user)
	}

	var message string

	if len(added) > 0 {
		message += fmt.Sprintf("Users assigned to standup: %v. ", strings.Join(added, ", "))
	}
	if len(failed) > 0 {
		message += fmt.Sprintf("Failed to assign: %v. ", strings.Join(failed, ", "))
	}
	if len(exist) > 0 {
		message += fmt.Sprintf("Already standupers: %v. ", strings.Join(exist, ", "))
	}

	msg := tgbotapi.NewMessage(event.Message.Chat.ID, message)
	_, err := b.tgAPI.Send(msg)
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

//Remove standupers
func (b *Bot) Remove(event tgbotapi.Update) error {
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
