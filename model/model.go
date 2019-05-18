package model

import "time"

//Group represents separate chat that bot was added to to handle standups
type Group struct {
	ID              int64
	ChatID          int64  `db:"chat_id" json:"chat_id,omitempty"`
	Title           string `db:"title" json:"title"`
	Description     string `db:"description" json:"description,omitempty"`
	StandupDeadline string `db:"standup_deadline" json:"standup_deadline,omitempty"`
}

// Standuper rerpesents standuper
type Standuper struct {
	ID           int64  `db:"id" json:"id"`
	Username     string `db:"username" json:"username"`
	ChatID       int64  `db:"chat_id" json:"chat_id"`
	LanguageCode string `db:"language_code" json:"language_code"`
}

// Standup model used for serialization/deserialization stored standups
type Standup struct {
	ID        int64     `db:"id" json:"id"`
	MessageID int       `db:"message_id" json:"message_id"`
	Created   time.Time `db:"created" json:"created"`
	Modified  time.Time `db:"modified" json:"modified"`
	Username  string    `db:"username" json:"userName"`
	Text      string    `db:"text" json:"text"`
	ChatID    int64     `db:"chat_id" json:"chat_id"`
}
