package domain

import "time"

type Video struct {
	ID         string
	ResourceID string
	FilePath   string
	CreatedAt  time.Time
}
