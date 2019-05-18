package bot

import tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

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
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Adding...")
	_, err := b.tgAPI.Send(msg)
	return err
}

//Show standupers
func (b *Bot) Show(event tgbotapi.Update) error {
	msg := tgbotapi.NewMessage(event.Message.Chat.ID, "Showing...")
	_, err := b.tgAPI.Send(msg)
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
