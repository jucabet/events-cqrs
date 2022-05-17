package database

import (
	"context"
	"database/sql"

	"github.com/jucabet/events-cqrs/models"

	_ "github.com/lib/pq"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}

	return &PostgresRepository{db}, nil
}

func (repo *PostgresRepository) Close() {
	repo.db.Close()
}

func (repo *PostgresRepository) InsertFeed(ctx context.Context, feed *models.Feed) error {
	_, err := repo.db.ExecContext(ctx, "INSERT INTO feeds (id, tittle, description) VALUES ($1, $2, $3)", feed.Id, feed.Title, feed.Description)
	return err
}

func (repo *PostgresRepository) ListFeeds(ctx context.Context) ([]*models.Feed, error) {
	rows, err := repo.db.QueryContext(ctx, "SELECT id, tittle, description, create_at FROM feeds")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var feeds []*models.Feed
	for rows.Next() {
		var feed models.Feed
		if err = rows.Scan(&feed.Id, &feed.Title, &feed.Description, &feed.CreateAt); err == nil {
			feeds = append(feeds, &feed)
		}
	}
	if rows.Err() != nil {
		return nil, rows.Err()
	}

	return feeds, nil
}
