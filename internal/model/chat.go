package model

import "time"

type Chat struct {
	ID        int64
	Usernames []string
}

type Message struct {
	ID        int64
	From      string
	Text      string
	Timestamp time.Time
}
