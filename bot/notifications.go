package bot

import (
	"fmt"
	"strings"
	"time"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/maddevsio/tgsbot/model"
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
		b.wg.Done()
	}
}

func (b *Bot) trackStandupersIn(team *model.Team) {
	ticker := time.NewTicker(time.Second * 60).C
	for {
		select {
		case <-ticker:
			loc, err := time.LoadLocation(team.Group.TZ)
			if err != nil {
				log.Error("LoadLocation failed! ", team)
				b.NotifyGroup(team.Group, time.Now().UTC())
				continue
			}
			b.NotifyGroup(team.Group, time.Now().In(loc))
		case <-team.QuitChan:
			log.Info("Finish working with the group: ", team.QuitChan)
			return
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

	if t.Hour() == r.Time.Hour() && t.Minute() == r.Time.Minute() {
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

		if len(missed) == 0 {
			msg := tgbotapi.NewMessage(group.ChatID, "Nice job, all standups submitted!")
			_, err = b.tgAPI.Send(msg)
			if err != nil {
				log.Error(err)
			}
			return
		}

		msg := tgbotapi.NewMessage(group.ChatID, fmt.Sprintf("Attention! Missed deadline: %v. Please, submit standups ASAP!", strings.Join(missed, ", ")))
		_, err = b.tgAPI.Send(msg)
		if err != nil {
			log.Error(err)
		}
	}
}
