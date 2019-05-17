package storage

import (
	"time"

	// This line is must for working MySQL database
	_ "github.com/go-sql-driver/mysql"

	"github.com/maddevsio/tgstandupbot/model"
)

// CreateStandup creates standup entry in database
func (m *MySQL) CreateStandup(s *model.Standup) (*model.Standup, error) {
	res, err := m.conn.Exec(
		"INSERT INTO `standups` (message_id, created, modified, username, text, chat_id) VALUES (?, ?, ?, ?, ?, ?)",
		s.MessageID, time.Now().UTC(), time.Now().UTC(), s.Username, s.Text, s.ChatID,
	)
	if err != nil {
		return s, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return s, err
	}
	s.ID = id
	return s, nil
}

// UpdateStandup updates standup entry in database
func (m *MySQL) UpdateStandup(s *model.Standup) (*model.Standup, error) {
	standup := &model.Standup{}
	m.conn.Exec(
		"UPDATE `standups` SET message_id=?, modified=?, username=?, text=? WHERE id=?",
		s.MessageID, time.Now().UTC(), s.Username, s.Text, s.ID,
	)
	err := m.conn.Get(&standup, "SELECT * FROM `standups` WHERE id=?", s.ID)
	return standup, err
}

// SelectStandup selects standup entry from database
func (m *MySQL) SelectStandup(id int64) (*model.Standup, error) {
	s := &model.Standup{}
	err := m.conn.Get(s, "SELECT * FROM `standups` WHERE id=?", id)
	return s, err
}

// DeleteStandup deletes standup entry from database
func (m *MySQL) DeleteStandup(id int64) error {
	_, err := m.conn.Exec("DELETE FROM `standups` WHERE id=?", id)
	return err
}

// ListStandups returns array of standup entries from database
func (m *MySQL) ListStandups() ([]*model.Standup, error) {
	items := []*model.Standup{}
	err := m.conn.Select(&items, "SELECT * FROM `standups`")
	return items, err
}

//LastStandupFor returns last standup for Standuper
func (m *MySQL) LastStandupFor(username string, chatID int64) (*model.Standup, error) {
	standup := &model.Standup{}
	err := m.conn.Get(standup, "SELECT * FROM `standups` WHERE username=? and chat_id=? ORDER BY id DESC LIMIT 1", username, chatID)
	return standup, err
}
