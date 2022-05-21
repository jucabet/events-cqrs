package main

import (
	"fmt"
	"net/http"

	"github.com/jucabet/events-cqrs/events"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	NatsAddress string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	hub := NewHub()

	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		panic(err)
	}

	err = n.OnCreateFeed(func(m events.CreateFeedMessage) {
		hub.Broadcast(NewCreateFeedMessage(m.Id, m.Title, m.Description, m.CreateAt), nil)
	})
	if err != nil {
		panic(err)
	}

	events.SetEventStore(n)
	defer events.Close()

	go hub.Run()

	http.HandleFunc("/ws", hub.HandleWebSocket)
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}
