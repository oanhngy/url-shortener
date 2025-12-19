package model

import "time"

type Link struct {
	ID        string
	LongURL   string
	ShortCode string
	CreatedAt time.Time
}
