package events

import (
	"bytes"
	"context"
	"encoding/gob"

	"github.com/jucabet/events-cqrs/models"
	"github.com/nats-io/nats.go"
)

type NatsEventStore struct {
	conn            *nats.Conn
	feedCreatedSub  *nats.Subscription
	feedCreatedChan chan CreateFeedMessage
}

func NewNats(url string) (*NatsEventStore, error) {
	conn, err := nats.Connect(url)
	if err != nil {
		return nil, err
	}

	return &NatsEventStore{
		conn: conn,
	}, nil
}

func (n *NatsEventStore) Close() {
	if n.conn != nil {
		n.conn.Close()
	}

	if n.feedCreatedSub != nil {
		n.feedCreatedSub.Unsubscribe()
	}

	close(n.feedCreatedChan)
}

func (n *NatsEventStore) encodeMessage(m Message) ([]byte, error) {
	b := bytes.Buffer{}
	err := gob.NewEncoder(&b).Encode(m)
	if err != nil {
		return nil, err
	}

	return b.Bytes(), nil
}

func (n *NatsEventStore) PublishCreatedFeed(ctx context.Context, feed *models.Feed) error {
	msg := CreateFeedMessage{
		Id:          feed.Id,
		Title:       feed.Title,
		Description: feed.Description,
		CreateAt:    feed.CreateAt,
	}

	data, err := n.encodeMessage(msg)
	if err != nil {
		return err
	}

	return n.conn.Publish(msg.Type(), data)
}

func (n *NatsEventStore) decodeMessage(data []byte, m interface{}) error {
	b := bytes.Buffer{}
	b.Write(data)

	return gob.NewDecoder(&b).Decode(m)
}

func (n *NatsEventStore) OnCreateFeed(f func(CreateFeedMessage)) (err error) {
	msg := CreateFeedMessage{}
	n.feedCreatedSub, err = n.conn.Subscribe(msg.Type(), func(m *nats.Msg) {
		n.decodeMessage(m.Data, &msg)
		f(msg)
	})

	return
}

func (n *NatsEventStore) SubscribeCreateFeed(ctx context.Context) (<-chan CreateFeedMessage, error) {
	m := CreateFeedMessage{}
	n.feedCreatedChan = make(chan CreateFeedMessage, 64)
	ch := make(chan *nats.Msg, 64)
	var err error
	n.feedCreatedSub, err = n.conn.ChanSubscribe(m.Type(), ch)
	if err != nil {
		return nil, err
	}

	go func() {
		for {
			select {
			case msg := <-ch:
				n.decodeMessage(msg.Data, m)
				n.feedCreatedChan <- m
			}
		}
	}()

	return (<-chan CreateFeedMessage)(n.feedCreatedChan), nil
}
