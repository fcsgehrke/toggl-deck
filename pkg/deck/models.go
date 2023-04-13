package deck

import (
	"errors"

	"github.com/gofrs/uuid"
)

var (
	ErrDeckNotFound          = errors.New("Deck not found")
	ErrCardCountNotAvailable = errors.New("Can't draw more cards than available on the deck")
	ErrInvalidCardCode       = errors.New("Some of cards has an invalid code")
	ErrInvalidDeckSize       = errors.New("The amount of cards is greater than it should be")
	ErrInvalidDeckID         = errors.New("Invalid Deck ID")
)

type CreateDeckResponse struct {
	ID        uuid.UUID `json:"deck_id"`
	Shuffled  bool      `json:"shuffled"`
	Remaining int       `json:"remaining"`
}

type CardResponse struct {
	Value string `json:"value"`
	Suit  string `json:"suit"`
	Code  string `json:"code"`
}

type OpenDeckResponse struct {
	ID        uuid.UUID       `json:"deck_id"`
	Shuffled  bool            `json:"shuffled"`
	Remaining int             `json:"remaining"`
	Cards     []*CardResponse `json:"cards,omitempty"`
}

type DrawCardResponse struct {
	Cards []*CardResponse `json:"cards"`
}
