package model

import "time"

//Group represents separate chat that bot was added to to handle standups
type Group struct {
	ID              int64
	ChatID          int64  `json:"chat_id,omitempty"`
	Title           string `json:"title"`
	Description     string `json:"description,omitempty"`
	StandupDeadline string `json:"standup_deadline,omitempty"`
}

// Standuper rerpesents standuper
type Standuper struct {
	ID       int64  `db:"id"`
	Username string `db:"username"`
	ChatID   int64  `db:"chat_id" json:"chat_id"`
}

// Standup model used for serialization/deserialization stored standups
type Standup struct {
	ID       int64     `db:"id" json:"id"`
	Created  time.Time `db:"created" json:"created"`
	Modified time.Time `db:"modified" json:"modified"`
	Username string    `db:"username" json:"userName"`
	Comment  string    `db:"comment" json:"comment"`
	ChatID   int64     `db:"chat_id" json:"chat_id"`
}
