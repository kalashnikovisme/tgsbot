package storage

import (
	"testing"

	"github.com/maddevsio/tgsbot/config"
	"github.com/maddevsio/tgsbot/model"
	"github.com/stretchr/testify/require"
)

func CreateGroupTest(t *testing.T) {
	c, err := config.Get()
	require.NoError(t, err)
	db, err := NewMySQL(c)
	require.NoError(t, err)

	group := &model.Group{
		ChatID:          int64(12345),
		Title:           "test",
		Description:     "Foo bar",
		StandupDeadline: "10am",
		TZ:              "Asia/Bishkek",
	}

	newGroup, err := db.CreateGroup(group)
	require.NoError(t, err)
	require.NotEqual(t, 0, newGroup.ID)

}
