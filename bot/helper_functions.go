package bot

import (
	"strings"
	"time"

	"github.com/maddevsio/telegramStandupBot/model"
)

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

func (b *Bot) submittedStandupToday(standuper *model.Standuper) bool {
	standup, err := b.db.LastStandupFor(standuper.Username, standuper.ChatID)
	if err != nil {
		return false
	}
	if standup.Created.Day() == time.Now().Day() {
		return true
	}
	return false
}
