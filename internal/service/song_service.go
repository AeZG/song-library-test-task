package service

import (
	"context"
	"errors"
	"fmt"
	"log"

	"song-library-test-task/internal/models"
)

// ExternalClient is an interface that will define the methods to call the external Swagger-based API
type ExternalClient interface {
	FetchSongInfo(ctx context.Context, groupName, songTitle string) (*SongInfo, error)
}

// SongInfo is a simple struct that represents the data from the external service
type SongInfo struct {
	ReleaseDate string
	Text        string
	Link        string
}

// SongService is the business logic layer for songs.
type SongService struct {
	repo   models.SongRepository
	client ExternalClient
}

// NewSongService constructs a new service object with the required dependencies.
func NewSongService(repo models.SongRepository, client ExternalClient) *SongService {
	return &SongService{
		repo:   repo,
		client: client,
	}
}

// CreateSong orchestrates adding a new song to the library.
// 1. Calls external API to get enrichment (releaseDate, text, link).
// 2. Inserts the record into Postgres via the repository.
// 3. Returns the new ID or an error.
func (uc *SongService) CreateSong(ctx context.Context, groupName, songTitle string) (int64, error) {
	log.Printf("[INFO] createSong: group=%s, title=%s", groupName, songTitle)

	// 1. Get external info (assuming it's required to store a complete record)
	songInfo, err := uc.client.FetchSongInfo(ctx, groupName, songTitle)
	if err != nil {
		// This could be a partial failure if you want to still create the record
		// but let's assume we want to fail if we cannot fetch enrichment
		return 0, fmt.Errorf("failed to fetch external data: %w", err)
	}

	// 2. Create models Song object
	song := &models.Song{
		GroupName:   groupName,
		Title:       songTitle,
		ReleaseDate: songInfo.ReleaseDate, // or parse to time.Time if you prefer
		Link:        songInfo.Link,
		Text:        songInfo.Text,
	}

	// 3. Insert into DB
	newID, err := uc.repo.Create(ctx, song)
	if err != nil {
		return 0, fmt.Errorf("failed to create new song: %w", err)
	}

	log.Printf("[INFO] Created song with ID=%d", newID)
	return newID, nil
}

// GetSong retrieves a song by ID from the repository.
func (uc *SongService) GetSong(ctx context.Context, songID int64) (*models.Song, error) {
	log.Printf("[DEBUG] getSong: id=%d", songID)

	s, err := uc.repo.GetByID(ctx, songID)
	if err != nil {
		return nil, fmt.Errorf("failed to retrieve song with ID=%d: %w", songID, err)
	}
	if s == nil {
		return nil, errors.New("song not found")
	}

	return s, nil
}

// ListSongs retrieves a paginated list of songs matching an optional filter.
func (uc *SongService) ListSongs(ctx context.Context, filter models.SongFilter, limit, offset int) ([]models.Song, error) {
	log.Printf("[DEBUG] listSongs: filter=%+v, limit=%d, offset=%d", filter, limit, offset)

	songs, err := uc.repo.GetAll(ctx, filter, limit, offset)
	if err != nil {
		return nil, fmt.Errorf("failed to list songs: %w", err)
	}
	return songs, nil
}

// GetSongLyrics is an example method that returns a slice of verses based on page/pageSize
func (s *SongService) GetSongLyrics(ctx context.Context, id int64, page, pageSize int) ([]string, int, error) {
	song, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, 0, err
	}
	if song == nil {
		return nil, 0, fmt.Errorf("song not found")
	}

	verses := splitByVerse(song.Text)
	total := len(verses)

	start := (page - 1) * pageSize
	if start < 0 {
		start = 0
	}
	end := start + pageSize
	if end > total {
		end = total
	}
	if start >= total {
		return []string{}, total, nil
	}

	return verses[start:end], total, nil
}

// UpdateSong updates the specified fields of an existing song.
func (uc *SongService) UpdateSong(ctx context.Context, song models.Song) error {
	log.Printf("[INFO] updateSong: id=%d", song.ID)

	existing, err := uc.repo.GetByID(ctx, song.ID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing song: %w", err)
	}
	if existing == nil {
		return errors.New("song not found")
	}

	if song.GroupName == "" {
		song.GroupName = existing.GroupName
	}
	if song.Title == "" {
		song.Title = existing.Title
	}
	if song.ReleaseDate == "" {
		song.ReleaseDate = existing.ReleaseDate
	}
	if song.Link == "" {
		song.Link = existing.Link
	}
	if song.Text == "" {
		song.Text = existing.Text
	}

	// Update in DB
	if err := uc.repo.Update(ctx, &song); err != nil {
		return fmt.Errorf("failed to update song: %w", err)
	}
	return nil
}

// DeleteSong removes the specified song from the DB.
func (uc *SongService) DeleteSong(ctx context.Context, songID int64) error {
	log.Printf("[INFO] deleteSong: id=%d", songID)

	existing, err := uc.repo.GetByID(ctx, songID)
	if err != nil {
		return fmt.Errorf("failed to fetch existing song: %w", err)
	}
	if existing == nil {
		return errors.New("song not found")
	}

	if err := uc.repo.Delete(ctx, songID); err != nil {
		return fmt.Errorf("failed to delete song: %w", err)
	}

	log.Printf("[INFO] Song with ID=%d deleted", songID)
	return nil
}

func splitByVerse(text string) []string {
	// For instance, split by double newlines
	// or do something more advanced
	return SplitByDoubleNewline(text)
}

func SplitByDoubleNewline(text string) []string {
	// E.g.: strings.Split(text, "\n\n")
	return []string{}
}
