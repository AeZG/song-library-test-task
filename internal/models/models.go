package models

import (
	"context"
	"time"
)

// Song represents the song info (business entity)
type Song struct {
	ID          int64
	GroupName   string
	Title       string
	ReleaseDate time.Time
	Link        string
	Text        string
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// SongFilter is used to filter the results in GetAll (list) calls.
type SongFilter struct {
	GroupName string
	Title     string
}

type SongRepository interface {
	Create(ctx context.Context, song *Song) (int64, error)
	GetByID(ctx context.Context, id int64) (*Song, error)
	GetAll(ctx context.Context, filter SongFilter, limit, offset int) ([]Song, error)
	Update(ctx context.Context, song *Song) error
	Delete(ctx context.Context, id int64) error
}
