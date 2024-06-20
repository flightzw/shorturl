package model

import (
	"time"
)

type Shorturl struct {
	ID        int64      `json:"id" bson:"_id"`
	URL       string     `json:"url"`
	CreatedAt *time.Time `json:"created_at" bson:"created_at"`
	UpdatedAt *time.Time `json:"updated_at" bson:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at" bson:"deleted_at"`
}

func NewShorturl(url string) *Shorturl {
	now := time.Now()
	return &Shorturl{
		URL:       url,
		CreatedAt: &now,
		UpdatedAt: &now,
	}
}
