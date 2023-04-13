package deck

import (
	"github.com/gofrs/uuid"
)

type Suit string

const (
	Clubs    Suit = "CLUBS"
	Diamonds Suit = "DIAMONDS"
	Hearts   Suit = "HEARTS"
	Spades   Suit = "SPADES"
)

type Card struct {
	Value string
	Suit  Suit
	Code  string
}

type Deck struct {
	ID        uuid.UUID
	Shuffled  bool
	Remaining int
	Cards     []*Card
}
