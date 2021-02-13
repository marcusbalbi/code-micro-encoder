package domain

import "time"

type job struct {
	ID               string
	OutputBucketPath string
	Status           string
	Video            *Video
	Error            string
	CreatedAt        time.Time
	UpdatedAt        time.Time
}
