package services

import (
	"context"
	"log"

	entities "github.com/fcsgehrke/toggl-deck/internal/domain/deck"
	"github.com/fcsgehrke/toggl-deck/pkg/deck"
	"github.com/gofrs/uuid"
)

type DeckRepo interface {
	CreateDeck(ctx context.Context, deck *entities.Deck) (*entities.Deck, error)
	ReadDeck(ctx context.Context, id uuid.UUID) (*entities.Deck, error)
	DrawDeckCard(ctx context.Context, id uuid.UUID, count int64) ([]*entities.Card, error)
}

type DeckService struct {
	repo DeckRepo
	log  *log.Logger
}

func NewDeckService(repo DeckRepo, log *log.Logger) (*DeckService, error) {
	return &DeckService{
		repo: repo,
		log:  log,
	}, nil
}

func (s *DeckService) CreateDeck(ctx context.Context, shuffled bool, cards []string) (*deck.CreateDeckResponse, error) {
	newDeck := &entities.Deck{Shuffled: shuffled}

	if len(cards) > 0 {
		crds, err := entities.ConvertCards(cards)
		if err != nil {
			s.log.Printf("[ERR] - Converting cards failed w/ err: %s", err.Error())
			return nil, err
		}
		newDeck.Cards = crds
	} else {
		newDeck.Cards = entities.GenerateCards()
	}

	if shuffled {
		newDeck.Cards = entities.ShuffleCards(newDeck.Cards)
	}

  newDeck.Remaining = len(newDeck.Cards)

	newDeck, err := s.repo.CreateDeck(ctx, newDeck)
	if err != nil {
		s.log.Printf("[ERR] - Could not create deck w/ err: %s\n", err.Error())
		return nil, err
	}

	return &deck.CreateDeckResponse{
		ID:        newDeck.ID,
		Shuffled:  newDeck.Shuffled,
		Remaining: newDeck.Remaining,
	}, nil
}

func (s *DeckService) OpenDeck(ctx context.Context, id string) (*deck.OpenDeckResponse, error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		s.log.Printf("[ERR] - UUID Convertion failed w/ err: %s\n", err.Error())
		return nil, deck.ErrInvalidDeckID
	}

	dk, err := s.repo.ReadDeck(ctx, uid)
	if err != nil {
		s.log.Printf("[ERR] - ReadDeck failed w/ err: %s\n", err.Error())
		return nil, deck.ErrDeckNotFound
	}

	return &deck.OpenDeckResponse{
		ID:        dk.ID,
		Shuffled:  dk.Shuffled,
		Remaining: dk.Remaining,
		Cards: func() []*deck.CardResponse {
			cards := []*deck.CardResponse{}
			for _, card := range dk.Cards {
				cards = append(cards, &deck.CardResponse{
					Value: card.Value,
					Suit:  string(card.Suit),
					Code:  card.Code,
				})
			}
			return cards
		}(),
	}, nil
}

func (s *DeckService) DrawCard(ctx context.Context, id string, count int64) (*deck.DrawCardResponse, error) {
	uid, err := uuid.FromString(id)
	if err != nil {
		s.log.Printf("[ERR] - UUID Convertion failed w/ err: %s\n", err.Error())
		return nil, deck.ErrInvalidDeckID
	}

	cards, err := s.repo.DrawDeckCard(ctx, uid, count)
	if err != nil {
		s.log.Printf("[ERR] - DrawDeckCard failed w/ err: %s\n", err.Error())
		return nil, deck.ErrDeckNotFound
	}

	return &deck.DrawCardResponse{
		Cards: func() []*deck.CardResponse {
			crds := []*deck.CardResponse{}
			for _, card := range cards {
				crds = append(crds, &deck.CardResponse{
					Value: card.Value,
					Suit:  string(card.Suit),
					Code:  card.Code,
				})
			}
			return crds
		}(),
	}, nil
}
