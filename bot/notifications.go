package bot

import (
	"fmt"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/ru"
	log "github.com/sirupsen/logrus"
)

func (b *Bot) StartNotificationThreads() {
	go func() {
		ticker := time.NewTicker(time.Second * 60).C
		for {
			select {
			case <-ticker:
				b.NotifyGroups(time.Now())
			}
		}
	}()
}

//NotifyGroups launches go routines that notify standupers
//about upcoming deadlines
func (b *Bot) NotifyGroups(t time.Time) {
	if int(t.Weekday()) == 6 || int(t.Weekday()) == 0 {
		return
	}

	groups, err := b.db.ListGroups()
	if err != nil {
		log.Error(err)
	}

	for _, group := range groups {
		if group.StandupDeadline == "" {
			continue
		}
		w := when.New(nil)
		w.Add(en.All...)
		w.Add(ru.All...)

		r, err := w.Parse(group.StandupDeadline, time.Now())
		if err != nil {
			log.Errorf("Unable to parse channel standup time [%v]: [%v]", group.StandupDeadline, err)
			continue
		}

		if r == nil {
			log.Infof("Could not find matches. Channel standup time: [%v]", group.StandupDeadline)
			continue
		}
		if !(t.Hour() == r.Time.Hour() && t.Minute() == r.Time.Minute()) {
			continue
		}
		standupers, err := b.db.ListChatStandupers(group.ChatID)
		if err != nil {
			log.Error(err)
			continue
		}

		for _, standuper := range standupers {
			if !b.submittedStandupToday(standuper) {
				msg := tgbotapi.NewMessage(standuper.ChatID, fmt.Sprintf("@%v, you have missed deadline, please, submit standup ASAP!", standuper.Username))
				_, err = b.tgAPI.Send(msg)
				if err != nil {
					log.Error(err)
				}
			}
		}
	}
}
