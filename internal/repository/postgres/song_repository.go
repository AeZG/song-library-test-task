package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/pkg/errors"
	"song-library-test-task/internal/models"
	"strings"
)

// songRepository is a Postgres-based implementation of domain.SongRepository.
type songRepository struct {
	db *sql.DB
}

// NewSongRepository returns a new instance of a Postgres song repository.
func NewSongRepository(db *sql.DB) models.SongRepository {
	return &songRepository{db: db}
}

// Create inserts a new song into the DB and returns the newly created ID.
func (r *songRepository) Create(ctx context.Context, song *models.Song) (int64, error) {
	query := `
        INSERT INTO songs (group_name, title, release_date, link, text, created_at, updated_at)
        VALUES ($1, $2, $3, $4, $5, NOW(), NOW())
        RETURNING id
    `

	var newID int64
	err := r.db.QueryRowContext(
		ctx,
		query,
		song.GroupName,
		song.Title,
		song.ReleaseDate,
		song.Link,
		song.Text,
	).Scan(&newID)
	if err != nil {
		return 0, errors.Wrap(err, "failed to insert new song")
	}

	return newID, nil
}

// GetByID retrieves a single song by its ID.
func (r *songRepository) GetByID(ctx context.Context, id int64) (*models.Song, error) {
	query := `
        SELECT
            id,
            group_name,
            title,
            release_date,
            link,
            text,
            created_at,
            updated_at
        FROM songs
        WHERE id = $1
        LIMIT 1
    `

	row := r.db.QueryRowContext(ctx, query, id)

	var s models.Song
	err := row.Scan(
		&s.ID,
		&s.GroupName,
		&s.Title,
		&s.ReleaseDate,
		&s.Link,
		&s.Text,
		&s.CreatedAt,
		&s.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, errors.Wrap(err, "failed to get song by ID")
	}

	return &s, nil
}

// GetAll retrieves songs from the DB matching the filter (if any) and applies pagination.
func (r *songRepository) GetAll(ctx context.Context, filter models.SongFilter, limit, offset int) ([]models.Song, error) {
	baseQuery := `
        SELECT
            id,
            group_name,
            title,
            release_date,
            link,
            text,
            created_at,
            updated_at
        FROM songs
    `
	whereClauses := []string{}
	args := []interface{}{}
	argPos := 1

	if filter.GroupName != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("group_name ILIKE $%d", argPos))
		args = append(args, "%"+filter.GroupName+"%")
		argPos++
	}

	if filter.Title != "" {
		whereClauses = append(whereClauses, fmt.Sprintf("title ILIKE $%d", argPos))
		args = append(args, "%"+filter.Title+"%")
		argPos++
	}

	if len(whereClauses) > 0 {
		baseQuery += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	// Add pagination
	baseQuery += fmt.Sprintf(" ORDER BY id DESC LIMIT %d OFFSET %d", limit, offset)

	rows, err := r.db.QueryContext(ctx, baseQuery, args...)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get songs")
	}
	defer rows.Close()

	var songs []models.Song
	for rows.Next() {
		var s models.Song
		err := rows.Scan(
			&s.ID,
			&s.GroupName,
			&s.Title,
			&s.ReleaseDate,
			&s.Link,
			&s.Text,
			&s.CreatedAt,
			&s.UpdatedAt,
		)
		if err != nil {
			return nil, errors.Wrap(err, "failed to scan row into Song")
		}
		songs = append(songs, s)
	}

	if err = rows.Err(); err != nil {
		return nil, errors.Wrap(err, "error iterating over song rows")
	}

	return songs, nil
}

// Update modifies an existing song's data in the DB.
func (r *songRepository) Update(ctx context.Context, song *models.Song) error {
	query := `
        UPDATE songs
        SET
            group_name   = $1,
            title        = $2,
            release_date = $3,
            link         = $4,
            text         = $5,
            updated_at   = NOW()
        WHERE id = $6
    `

	_, err := r.db.ExecContext(
		ctx,
		query,
		song.GroupName,
		song.Title,
		song.ReleaseDate,
		song.Link,
		song.Text,
		song.ID,
	)
	if err != nil {
		return errors.Wrap(err, "failed to update song")
	}

	return nil
}

// Delete removes a song record by ID.
func (r *songRepository) Delete(ctx context.Context, id int64) error {
	query := `DELETE FROM songs WHERE id = $1`

	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return errors.Wrap(err, "failed to delete song")
	}

	return nil
}
