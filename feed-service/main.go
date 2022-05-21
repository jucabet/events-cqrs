package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jucabet/events-cqrs/database"
	"github.com/jucabet/events-cqrs/events"
	"github.com/jucabet/events-cqrs/repository"
	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	PostgresDB       string `envconfig:"POSTGRES_DB"`
	PostgresUser     string `envconfig:"POSTGRES_USER"`
	PostgresPassword string `envconfig:"POSTGRES_PASSWORD"`
	NatsAddress      string `envconfig:"NATS_ADDRESS"`
}

func main() {
	var cfg Config
	err := envconfig.Process("", &cfg)
	if err != nil {
		panic(err)
	}

	addr := fmt.Sprintf("postgres://%s:%s@postgres/%s?sslmode=disable", cfg.PostgresUser, cfg.PostgresPassword, cfg.PostgresDB)
	repo, err := database.NewPostgresRepository(addr)
	if err != nil {
		panic(err)
	}

	repository.SetRepository(repo)

	n, err := events.NewNats(fmt.Sprintf("nats://%s", cfg.NatsAddress))
	if err != nil {
		panic(err)
	}

	events.SetEventStore(n)

	defer events.Close()

	router := newRouter()
	if err := http.ListenAndServe(":8080", router); err != nil {
		panic(err)
	}
}

func newRouter() *mux.Router {
	var router = mux.NewRouter()
	router.HandleFunc("/feeds", createFeedHandler).Methods(http.MethodPost)

	return router
}
