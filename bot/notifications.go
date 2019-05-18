package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/telegramStandupBot/model"
	"github.com/olebedev/when"
	"github.com/olebedev/when/rules/en"
	"github.com/olebedev/when/rules/ru"
	log "github.com/sirupsen/logrus"
)

//StartWatchers looks for new gropus from the channel and start watching it
func (b *Bot) StartWatchers() {
	for group := range b.watchersChan {
		log.Info("New group to track: ", group)
		team := &model.Team{
			Group:    group,
			QuitChan: make(chan struct{}),
		}
		b.teams = append(b.teams, team)
		b.wg.Add(1)
		go b.trackStandupersIn(team)
	}
}

func (b *Bot) trackStandupersIn(team *model.Team) {
	ticker := time.NewTicker(time.Second * 60).C
	for {
		select {
		case <-ticker:
			b.NotifyGroup(team.Group, time.Now())
		case <-team.QuitChan:
			b.wg.Done()
			break
		}
	}
}

//NotifyGroup launches go routines that notify standupers
//about upcoming deadlines
func (b *Bot) NotifyGroup(group *model.Group, t time.Time) {
	if int(t.Weekday()) == 6 || int(t.Weekday()) == 0 {
		return
	}

	if group.StandupDeadline == "" {
		return
	}
	w := when.New(nil)
	w.Add(en.All...)
	w.Add(ru.All...)

	r, err := w.Parse(group.StandupDeadline, time.Now())
	if err != nil {
		log.Errorf("Unable to parse channel standup time [%v]: [%v]", group.StandupDeadline, err)
		return
	}

	if r == nil {
		log.Infof("Could not find matches. Channel standup time: [%v]", group.StandupDeadline)
		return
	}
	if !(t.Hour() == r.Time.Hour() && t.Minute() == r.Time.Minute()) {
		return
	}
	standupers, err := b.db.ListChatStandupers(group.ChatID)
	if err != nil {
		log.Error(err)
		return
	}

	missed := []string{}

	for _, standuper := range standupers {
		if !b.submittedStandupToday(standuper) {
			missed = append(missed, "@"+standuper.Username)
		}
	}

	msg := tgbotapi.NewMessage(group.ChatID, fmt.Sprintf("@%v, you have missed deadline, please, submit standup ASAP!", strings.Join(missed, ", ")))
	_, err = b.tgAPI.Send(msg)
	if err != nil {
		log.Error(err)
	}
}
