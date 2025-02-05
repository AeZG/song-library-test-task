package http

import (
	"context"
	"encoding/json"
	"net/http"
	"strconv"

	kithttp "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"

	"song-library-test-task/internal/handler/http/endpoints"
)

// NewHTTPHandler constructs a http.Handler with all the Song routes.
func NewHTTPHandler(eps endpoints.SongEndpoints) http.Handler {
	r := mux.NewRouter()

	// --------------------------------------------------------------------------------
	// Create a new song
	// --------------------------------------------------------------------------------
	// CreateSong godoc
	// @Summary     Create a new song
	// @Description Provide a JSON body with "group" and "song" fields. This also calls an external API to enrich the data with release date, text, and link.
	// @Tags        songs
	// @Accept      json
	// @Produce     json
	// @Param       input body endpoints.CreateSongRequest true "New Song Data"
	// @Success     201 {object} endpoints.CreateSongResponse
	// @Failure     400 {object} errorResponse
	// @Failure     500 {object} errorResponse
	// @Router      /songs [post]
	r.Handle("/songs",
		kithttp.NewServer(
			eps.CreateSongEndpoint,
			decodeCreateSongRequest,
			encodeJSONResponse,
		),
	).Methods("POST")

	// --------------------------------------------------------------------------------
	// List songs with optional filtering and pagination
	// --------------------------------------------------------------------------------
	// ListSongs godoc
	// @Summary     List songs
	// @Description Returns a list of songs with optional filtering by group, title, and pagination.
	// @Tags        songs
	// @Produce     json
	// @Param       group  query   string false "Filter by group name (partial match)"
	// @Param       title  query   string false "Filter by song title (partial match)"
	// @Param       limit  query   int    false "Max records to return (default 10)"
	// @Param       offset query   int    false "Offset from first record (default 0)"
	// @Success     200 {object} endpoints.ListSongsResponse
	// @Failure     500 {object} errorResponse
	// @Router      /songs [get]
	r.Handle("/songs",
		kithttp.NewServer(
			eps.ListSongsEndpoint,
			decodeListSongsRequest,
			encodeJSONResponse,
		),
	).Methods("GET")

	// --------------------------------------------------------------------------------
	// Get a single song by ID
	// --------------------------------------------------------------------------------
	// GetSong godoc
	// @Summary     Get song by ID
	// @Description Returns a single song's data by its ID.
	// @Tags        songs
	// @Produce     json
	// @Param       id   path int true "Song ID"
	// @Success     200 {object} endpoints.GetSongResponse
	// @Failure     400 {object} errorResponse
	// @Failure     500 {object} errorResponse
	// @Router      /songs/{id} [get]
	r.Handle("/songs/{id}",
		kithttp.NewServer(
			eps.GetSongEndpoint,
			decodeGetSongRequest,
			encodeJSONResponse,
		),
	).Methods("GET")

	// --------------------------------------------------------------------------------
	// Update an existing song by ID
	// --------------------------------------------------------------------------------
	// UpdateSong godoc
	// @Summary     Update an existing song
	// @Description Updates a song's fields by ID (e.g. title, group, release date, link, text).
	// @Tags        songs
	// @Accept      json
	// @Produce     json
	// @Param       id    path   int  true "Song ID"
	// @Param       input body   endpoints.UpdateSongRequest true "Song Data"
	// @Success     200 {object} endpoints.UpdateSongResponse
	// @Failure     400 {object} errorResponse
	// @Failure     500 {object} errorResponse
	// @Router      /songs/{id} [put]
	r.Handle("/songs/{id}",
		kithttp.NewServer(
			eps.UpdateSongEndpoint,
			decodeUpdateSongRequest,
			encodeJSONResponse,
		),
	).Methods("PUT")

	// --------------------------------------------------------------------------------
	// Delete a song by ID
	// --------------------------------------------------------------------------------
	// DeleteSong godoc
	// @Summary     Delete a song
	// @Description Removes a song from the database by ID.
	// @Tags        songs
	// @Produce     json
	// @Param       id   path  int  true "Song ID"
	// @Success     204 "No Content"
	// @Failure     400 {object} errorResponse
	// @Failure     500 {object} errorResponse
	// @Router      /songs/{id} [delete]
	r.Handle("/songs/{id}",
		kithttp.NewServer(
			eps.DeleteSongEndpoint,
			decodeDeleteSongRequest,
			encodeJSONResponse,
		),
	).Methods("DELETE")

	// --------------------------------------------------------------------------------
	// Get song lyrics (verses) with pagination
	// --------------------------------------------------------------------------------
	// GetLyrics godoc
	// @Summary     Get lyrics by verse
	// @Description Returns paginated verses of the song text, by ID. For example, page=1&pageSize=1 returns the first verse.
	// @Tags        songs
	// @Produce     json
	// @Param       id        path  int true "Song ID"
	// @Param       page      query int false "Verse page (default 1)"
	// @Param       pageSize  query int false "Verses per page (default 1)"
	// @Success     200 {object} endpoints.GetLyricsResponse
	// @Failure     400 {object} errorResponse
	// @Failure     500 {object} errorResponse
	// @Router      /songs/{id}/lyrics [get]
	r.Handle("/songs/{id}/lyrics",
		kithttp.NewServer(
			eps.GetLyricsEndpoint,
			decodeGetLyricsRequest,
			encodeJSONResponse,
		),
	).Methods("GET")

	return r
}

// --------------------------------------------------------------------------------
// Decode functions
// --------------------------------------------------------------------------------

func decodeCreateSongRequest(_ context.Context, r *http.Request) (interface{}, error) {
	var req endpoints.CreateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		return nil, err
	}
	return req, nil
}

func decodeListSongsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vals := r.URL.Query()
	group := vals.Get("group")
	title := vals.Get("title")
	limit, _ := strconv.Atoi(vals.Get("limit"))
	offset, _ := strconv.Atoi(vals.Get("offset"))

	req := endpoints.ListSongsRequest{
		GroupName: group,
		Title:     title,
		Limit:     limit,
		Offset:    offset,
	}
	return req, nil
}

func decodeGetSongRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRoute // define an error
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return endpoints.GetSongRequest{ID: id}, nil
}

func decodeUpdateSongRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	var body endpoints.UpdateSongRequest
	if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
		return nil, err
	}
	body.ID = id
	return body, nil
}

func decodeDeleteSongRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}
	return endpoints.DeleteSongRequest{ID: id}, nil
}

func decodeGetLyricsRequest(_ context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)
	idStr, ok := vars["id"]
	if !ok {
		return nil, errBadRoute
	}
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		return nil, err
	}

	q := r.URL.Query()
	page, _ := strconv.Atoi(q.Get("page"))
	if page < 1 {
		page = 1
	}
	pageSize, _ := strconv.Atoi(q.Get("pageSize"))
	if pageSize < 1 {
		pageSize = 1
	}

	return endpoints.GetLyricsRequest{
		ID:       id,
		Page:     page,
		PageSize: pageSize,
	}, nil
}

// --------------------------------------------------------------------------------
// Encode (response) functions
// --------------------------------------------------------------------------------

func encodeJSONResponse(_ context.Context, w http.ResponseWriter, response interface{}) error {
	if f, ok := response.(failureer); ok && f.Failed() != nil {
		http.Error(w, f.Failed().Error(), http.StatusInternalServerError)
		return nil
	}
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	return json.NewEncoder(w).Encode(response)
}

var errBadRoute = &BadRouteError{"bad route"}

type BadRouteError struct{ msg string }

func (e *BadRouteError) Error() string { return e.msg }

type failureer interface {
	Failed() error
}
