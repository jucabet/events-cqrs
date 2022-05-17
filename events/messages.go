package events

import "time"

type Message interface {
	Type() string
}

type CreateFeedMessage struct {
	Id          string    `json:"id"`
	Title       string    `json:"tittle"`
	Description string    `json:"description"`
	CreateAt    time.Time `json:"create_at"`
}

func (m CreateFeedMessage) Type() string {
	return "create_feed"
}
