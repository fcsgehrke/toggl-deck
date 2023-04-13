package api

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	"github.com/fcsgehrke/toggl-deck/pkg/deck"
)

type DeckService interface {
	CreateDeck(ctx context.Context, shuffled bool, cards []string) (*deck.CreateDeckResponse, error)
	OpenDeck(ctx context.Context, id string) (*deck.OpenDeckResponse, error)
	DrawCard(ctx context.Context, id string, count int64) (*deck.DrawCardResponse, error)
}

type DeckHandler struct {
	service DeckService
	log     *log.Logger
}

func NewDeckHandler(service DeckService, log *log.Logger) (*DeckHandler, error) {
	return &DeckHandler{
		service: service,
		log:     log,
	}, nil
}

// CreateDeck godoc
//	@Summary		Creates a new deck
//	@Description	Creates new decks with all cards or specific ones. Also, it's possible to shuffle the deck.
//	@Accept			json
//	@Produce		json
//	@Param			shuffled	query		bool	false	"Shuffle cards on the deck."
//	@Param			cards		query		array	false	"The card codes to be added to the deck."
//	@Success		200			{object}	deck.CreateDeckResponse
//	@Failure		400			{object}	HttpError
//	@Failure		404			{object}	HttpError
//	@Failure		422			{object}	HttpError
//	@Failure		500			{object}	HttpError
//	@Router			/ [post]
func (h *DeckHandler) CreateDeck(w http.ResponseWriter, r *http.Request) {
	var shuffled bool = false
	var cards []string

	query := r.URL.Query()

	// Reading Shuffled param from URL
	if data, ok := query["shuffled"]; ok {
		s, err := strconv.ParseBool(data[0])
		if err != nil {
    httpError(w, err, http.StatusUnprocessableEntity)
			return
		}

		shuffled = s
	}

	// Reading Cards param from the URL
	if data, ok := query["cards"]; ok {
		if strings.Contains(data[0], ",") {
      cards = strings.Split(data[0], ",")
    }
	}

	response, err := h.service.CreateDeck(r.Context(), shuffled, cards)
	if err != nil {
    httpError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
    httpError(w, err, http.StatusInternalServerError)
		h.log.Printf("[ERR] - JSON Encode failed w/ err: %s\n", err.Error())
	}
}

// OpenDeck godoc
//	@Summary		Opens an existing deck
//	@Description	Returns an existing opened deck with all cards.
//	@Accept			json
//	@Produce		json
//	@Param			id	path		string	true	"The ID (uuid) of the deck."
//	@Success		200	{object}	deck.OpenDeckResponse
//	@Failure		400	{object}	HttpError
//	@Failure		404	{object}	HttpError
//	@Failure		422	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/{id} [get]
func (h *DeckHandler) OpenDeck(w http.ResponseWriter, r *http.Request) {
	id, err := extractId(r)
	if err != nil {
    httpError(w, err, http.StatusBadRequest)
	}

	response, err := h.service.OpenDeck(r.Context(), id)
	if err != nil {
		if err == deck.ErrDeckNotFound {
    httpError(w, err, http.StatusNotFound)
			h.log.Printf("[ERR] - Deck id: %s, not found\n", id)
			return
		}

    httpError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
    httpError(w, err, http.StatusInternalServerError)
		h.log.Printf("[ERR] - JSON Encode failed w/ err: %s\n", err.Error())
	}
}

// DrawCard godoc
//	@Summary		Draws a card from the deck
//	@Description	Draws one or more cards from the deck
//	@Param			id		path	string	true	"The ID (uuid) of the deck."
//	@Param			count	query	int		true	"The number of cards to draw from the deck."
//	@Accept			json
//	@Produce		json
//	@Success		200	{object}	deck.DrawCardResponse
//	@Failure		400	{object}	HttpError
//	@Failure		404	{object}	HttpError
//	@Failure		422	{object}	HttpError
//	@Failure		500	{object}	HttpError
//	@Router			/{id}/draw [get]
func (h *DeckHandler) DrawCard(w http.ResponseWriter, r *http.Request) {
	id, err := extractId(r)
	if err != nil {
    httpError(w, err, http.StatusBadRequest)
	}

	query := r.URL.Query()

	var count int64

	// Reading Count param from URL
	if data, ok := query["count"]; ok {
		cnt, err := strconv.ParseInt(data[0], 10, 8)
		if err != nil {
      httpError(w, err, http.StatusUnprocessableEntity)
			return
		}

		count = cnt
	}

	response, err := h.service.DrawCard(r.Context(), id, count)
	if err != nil {
		if err == deck.ErrDeckNotFound {
			httpError(w, err, http.StatusNotFound)
			h.log.Printf("[ERR] - Deck id: %s, not found\n", id)
			return
		}

    httpError(w, err, http.StatusBadRequest)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
    httpError(w, err, http.StatusInternalServerError)
		h.log.Printf("[ERR] - JSON Encode failed w/ err: %s\n", err.Error())
	}
}

// You might be wondering why this ? 
// Well, I decided to do it to make it agnostic in terms of web framework.
// It's easier this way to put this endpoints to run on a GCP cloud function
// or AWS lamdba for example. You're not gonna see any 
// chi/echo/fiber/gin import in this file.
func extractId(r *http.Request) (string, error) {
	var re = regexp.MustCompile(`(?m)[a-fA-F0-9]{8}-[a-fA-F0-9]{4}-4[a-fA-F0-9]{3}-[8|9|aA|bB][a-fA-F0-9]{3}-[a-fA-F0-9]{12}`)
	for _, match := range re.FindAllString(r.URL.Path, -1) {
    return match, nil
	}
  return "", deck.ErrInvalidDeckID
}

func httpError(w http.ResponseWriter, err error, status int) {
  w.WriteHeader(status)
	if err := json.NewEncoder(w).Encode(HttpError{Message: err.Error()}); err != nil {
    httpError(w, err, http.StatusInternalServerError)
	}
}
