package mocks

import (
	"context"
	"fmt"

	entities "github.com/fcsgehrke/toggl-deck/internal/domain/deck"
	"github.com/gofrs/uuid"
)

type DeckRepoMock struct {
	data map[string]*entities.Deck
}

func NewDeckRepoMock() *DeckRepoMock {
	return &DeckRepoMock{
		data: map[string]*entities.Deck{},
	}
}

func (r *DeckRepoMock) CreateDeck(ctx context.Context, deck *entities.Deck) (*entities.Deck, error) {
  uid, err := uuid.NewV4()
  if err != nil {
    return nil, err
  }

  r.data[uid.String()] = deck
  deck.ID = uid

  return deck, nil
}

func (r *DeckRepoMock) ReadDeck(ctx context.Context, id uuid.UUID) (*entities.Deck, error) {
  if dk, ok := r.data[id.String()]; ok {
    return dk, nil
  }

  return nil, fmt.Errorf("record not found.")
}

func (r *DeckRepoMock) DrawDeckCard(ctx context.Context, id uuid.UUID, count int64) ([]*entities.Card, error) {
  if _, ok := r.data[id.String()]; !ok {
  return nil, fmt.Errorf("record not found.")
  }

  dk := r.data[id.String()]  

  max := int(count)
  if int(count) > len(dk.Cards) {
    max = len(dk.Cards)
  }

  cards := dk.Cards[0:max]
  dk.Remaining -= max
  dk.Cards = dk.Cards[max:]
 
  return cards, nil
}
