package storage

import (
	"time"

	// This line is must for working MySQL database
	_ "github.com/go-sql-driver/mysql"

	"github.com/jmoiron/sqlx"
	"github.com/maddevsio/tgstandupbot/config"
	"github.com/maddevsio/tgstandupbot/model"
)

// MySQL provides api for work with mysql database
type MySQL struct {
	conn *sqlx.DB
}

// NewMySQL creates a new instance of database API
func NewMySQL(c *config.BotConfig) (*MySQL, error) {
	m := &MySQL{}
	conn, err := sqlx.Open("mysql", c.DatabaseURL)
	if err != nil {
		return nil, err
	}
	m.conn = conn
	return m, nil
}

// CreateStandup creates standup entry in database
func (m *MySQL) CreateStandup(s *model.Standup) (*model.Standup, error) {
	res, err := m.conn.Exec(
		"INSERT INTO `standups` (created, modified, username, comment, groupid) VALUES (?, ?, ?, ?, ?)",
		time.Now().UTC(), time.Now().UTC(), s.Username, s.Comment, s.ChatID,
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
		"UPDATE `standups` SET modified=?, username=?, comment=? WHERE id=?",
		time.Now().UTC(), s.Username, s.Comment, s.ID,
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
func (m *MySQL) LastStandupFor(username string, groupID int64) (*model.Standup, error) {
	standup := &model.Standup{}
	err := m.conn.Get(standup, "SELECT * FROM `standups` WHERE username=? and groupid=? ORDER BY id DESC LIMIT 1", username, groupID)
	return standup, err
}

// CreateStanduper creates Standuper
func (m *MySQL) CreateStanduper(s *model.Standuper) (*model.Standuper, error) {
	res, err := m.conn.Exec(
		"INSERT INTO `standupers` (username, groupid) VALUES (?, ?)",
		s.Username, s.ChatID,
	)
	if err != nil {
		return nil, err
	}
	id, err := res.LastInsertId()
	if err != nil {
		return nil, err
	}
	s.ID = id
	return s, nil
}

// UpdateStanduper updates Standuper entry in database
func (m *MySQL) UpdateStanduper(s *model.Standuper) (*model.Standuper, error) {
	m.conn.Exec(
		"UPDATE `standupers` SET username=? WHERE id=?",
		s.Username, s.ID,
	)
	err := m.conn.Get(s, "SELECT * FROM `standupers` WHERE id=?", s.ID)
	return s, err
}

// SelectStanduper selects Standuper entry from database
func (m *MySQL) SelectStanduper(id int64) (*model.Standuper, error) {
	s := &model.Standuper{}
	err := m.conn.Get(s, "SELECT * FROM `standupers` WHERE id=?", id)
	return s, err
}

// FindStanduper selects Standuper entry from database
func (m *MySQL) FindStanduper(name string, groupID int64) (*model.Standuper, error) {
	s := &model.Standuper{}
	err := m.conn.Get(s, "SELECT * FROM `standupers` WHERE username=? and groupid=?", name, groupID)
	return s, err
}

// DeleteStanduper deletes Standuper entry from database
func (m *MySQL) DeleteStanduper(id int64) error {
	_, err := m.conn.Exec("DELETE FROM `standupers` WHERE id=?", id)
	return err
}

// ListStandupers returns array of Standuper entries from database
func (m *MySQL) ListStandupers() ([]*model.Standuper, error) {
	standupers := []*model.Standuper{}
	err := m.conn.Select(&standupers, "SELECT * FROM `standupers`")
	return standupers, err
}
