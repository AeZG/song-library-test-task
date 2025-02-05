package endpoints

import (
	"context"
	"song-library-test-task/internal/models"

	"github.com/go-kit/kit/endpoint"
	"song-library-test-task/internal/service"
)

// SongEndpoints bundles all endpoints for the SongService
type SongEndpoints struct {
	CreateSongEndpoint endpoint.Endpoint
	GetSongEndpoint    endpoint.Endpoint
	ListSongsEndpoint  endpoint.Endpoint
	UpdateSongEndpoint endpoint.Endpoint
	DeleteSongEndpoint endpoint.Endpoint
	GetLyricsEndpoint  endpoint.Endpoint
}

// MakeSongEndpoints constructs a SongEndpoints struct with all endpoints
func MakeSongEndpoints(s service.SongService) SongEndpoints {
	return SongEndpoints{
		CreateSongEndpoint: makeCreateSongEndpoint(s),
		GetSongEndpoint:    makeGetSongEndpoint(s),
		ListSongsEndpoint:  makeListSongsEndpoint(s),
		UpdateSongEndpoint: makeUpdateSongEndpoint(s),
		DeleteSongEndpoint: makeDeleteSongEndpoint(s),
		GetLyricsEndpoint:  makeGetLyricsEndpoint(s),
	}
}

// Request/Response for each operation:

type CreateSongRequest struct {
	GroupName string `json:"group"`
	Title     string `json:"song"`
}
type CreateSongResponse struct {
	ID  int64  `json:"id"`
	Err string `json:"error,omitempty"`
}

func makeCreateSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(CreateSongRequest)
		id, err := s.CreateSong(ctx, req.GroupName, req.Title)
		if err != nil {
			return CreateSongResponse{Err: err.Error()}, nil
		}
		return CreateSongResponse{ID: id}, nil
	}
}

// Get Song
type GetSongRequest struct {
	ID int64
}
type GetSongResponse struct {
	Song *models.Song `json:"song,omitempty"`
	Err  string       `json:"error,omitempty"`
}

func makeGetSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetSongRequest)
		song, err := s.GetSong(ctx, req.ID)
		if err != nil {
			return GetSongResponse{Err: err.Error()}, nil
		}
		return GetSongResponse{Song: song}, nil
	}
}

// List Songs
type ListSongsRequest struct {
	GroupName string
	Title     string
	Limit     int
	Offset    int
}
type ListSongsResponse struct {
	Songs []models.Song `json:"songs"`
	Err   string        `json:"error,omitempty"`
}

func makeListSongsEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(ListSongsRequest)
		filter := models.SongFilter{
			GroupName: req.GroupName,
			Title:     req.Title,
		}
		songs, err := s.ListSongs(ctx, filter, req.Limit, req.Offset)
		if err != nil {
			return ListSongsResponse{Err: err.Error()}, nil
		}
		return ListSongsResponse{Songs: songs}, nil
	}
}

// Update Song
type UpdateSongRequest struct {
	ID          int64  `json:"-"`
	GroupName   string `json:"group"`
	Title       string `json:"song"`
	ReleaseDate string `json:"releaseDate"`
	Link        string `json:"link"`
	Text        string `json:"text"`
}
type UpdateSongResponse struct {
	Err string `json:"error,omitempty"`
}

func makeUpdateSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(UpdateSongRequest)
		err := s.UpdateSong(ctx, models.Song{
			ID:          req.ID,
			GroupName:   req.GroupName,
			Title:       req.Title,
			ReleaseDate: req.ReleaseDate,
			Link:        req.Link,
			Text:        req.Text,
		})
		if err != nil {
			return UpdateSongResponse{Err: err.Error()}, nil
		}
		return UpdateSongResponse{}, nil
	}
}

// Delete Song
type DeleteSongRequest struct {
	ID int64
}
type DeleteSongResponse struct {
	Err string `json:"error,omitempty"`
}

func makeDeleteSongEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(DeleteSongRequest)
		err := s.DeleteSong(ctx, req.ID)
		if err != nil {
			return DeleteSongResponse{Err: err.Error()}, nil
		}
		return DeleteSongResponse{}, nil
	}
}

// GetLyrics
type GetLyricsRequest struct {
	ID       int64
	Page     int
	PageSize int
}
type GetLyricsResponse struct {
	Lyrics []string `json:"lyrics"`
	Total  int      `json:"total"`
	Err    string   `json:"error,omitempty"`
}

func makeGetLyricsEndpoint(s service.SongService) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		req := request.(GetLyricsRequest)
		verses, total, err := s.GetSongLyrics(ctx, req.ID, req.Page, req.PageSize)
		if err != nil {
			return GetLyricsResponse{Err: err.Error()}, nil
		}
		return GetLyricsResponse{Lyrics: verses, Total: total}, nil
	}
}
