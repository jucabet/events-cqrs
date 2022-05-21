package events

import (
	"context"

	"github.com/jucabet/events-cqrs/models"
)

type EventsStore interface {
	Close()
	PublishCreatedFeed(ctx context.Context, feed *models.Feed) error
	SubscribeCreateFeed(ctx context.Context) (<-chan CreateFeedMessage, error)
	OnCreateFeed(f func(CreateFeedMessage)) error
}

var eventStore EventsStore

func SetEventStore(store EventsStore) {
	eventStore = store
}

func Close() {
	eventStore.Close()
}

func PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	return eventStore.PublishCreatedFeed(ctx, feed)
}

func SubscribeCreateFeed(ctx context.Context) (<-chan CreateFeedMessage, error) {
	return eventStore.SubscribeCreateFeed(ctx)
}

func OnCreateFeed(ctx context.Context, f func(CreateFeedMessage)) error {
	return eventStore.OnCreateFeed(f)
}
