package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/jucabet/events-cqrs/events"
	"github.com/jucabet/events-cqrs/models"
	"github.com/jucabet/events-cqrs/repository"
	"github.com/jucabet/events-cqrs/search"
)

func onCreatedFeed(m events.CreateFeedMessage) {
	feed := models.Feed{
		Id:          m.Id,
		Title:       m.Title,
		Description: m.Description,
		CreateAt:    m.CreateAt,
	}

	if err := search.IndexFeed(context.Background(), feed); err != nil {
		fmt.Printf("Error on indexing feed: %v", err)
	}
}

func listFeedsHandler(w http.ResponseWriter, r *http.Request) {
	feeds, err := repository.ListFeeds(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}

func searchHandler(w http.ResponseWriter, r *http.Request) {
	var query string = r.URL.Query().Get("q")
	if len(query) == 0 {
		http.Error(w, "Empty Query", http.StatusBadRequest)
		return
	}

	feeds, err := search.SearchFeed(r.Context(), query)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(feeds)
}
