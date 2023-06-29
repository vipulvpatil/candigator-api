package model

import "time"

type Candidate struct {
	id        string
	email     string
	createdAt time.Time
}
