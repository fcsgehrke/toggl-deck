package postgres

import (
	"context"
	"log"

	entities "github.com/fcsgehrke/toggl-deck/internal/domain/deck"
	"github.com/gofrs/uuid"
	"github.com/lib/pq"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type DeckPostgresRepo struct {
	db  *gorm.DB
	log *log.Logger
}

func NewDeckPostgresRepo(ctx context.Context, dsn string, log *log.Logger) (*DeckPostgresRepo, error) {
  log.Println("[INF] - Connecting DB...")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Printf("[ERR] - DB Conn Open failed w/ err: %s\n", err.Error())
		return nil, err
	}

  log.Println("[INF] - DB Connected Successfully")

	return &DeckPostgresRepo{
		db:  db,
		log: log,
	}, nil
}

func (r *DeckPostgresRepo) Migrate() error {
  r.log.Println("[INF] - Running DB Migrations...")
  if err := r.db.AutoMigrate(&Deck{}); err != nil {
    r.log.Printf("[ERR] - DB Migrations failed w/ err: %s\n", err.Error())
    return err
  }
  r.log.Println("[INF] - DB Migrations Run Successfully")
  return nil
}

func (r *DeckPostgresRepo) CreateDeck(ctx context.Context, deck *entities.Deck) (*entities.Deck, error) {
  cards := []string{}
  for _, card := range deck.Cards {
    cards = append(cards, card.Code)
  }

  dbDeck := Deck{
   Shuffled: deck.Shuffled,
   Cards: pq.StringArray(cards),
   Remaining: len(cards),
  }

	err := r.db.WithContext(ctx).Create(&dbDeck).Error
	if err != nil {
		r.log.Printf("[ERR] - Create Deck failed w/ err: %s\n", err.Error())
		return nil, err
	}

  deck.ID = dbDeck.ID
  deck.Remaining = dbDeck.Remaining

	return deck, nil
}

func (r *DeckPostgresRepo) ReadDeck(ctx context.Context, id uuid.UUID) (*entities.Deck, error) {
	var dbDeck Deck
	err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbDeck).Error
	if err != nil {
		return nil, err
	}

  cards, err := entities.ConvertCards(dbDeck.Cards)
  if err != nil {
    return nil, err
  }

  return &entities.Deck{
    ID: dbDeck.ID,
    Shuffled: dbDeck.Shuffled,
    Remaining: dbDeck.Remaining,
    Cards: cards,
  }, nil
}

func (r *DeckPostgresRepo) DrawDeckCard(ctx context.Context, id uuid.UUID, count int64) ([]*entities.Card, error) {
	var dbDeck Deck

  // Getting the Deck
	if err := r.db.WithContext(ctx).Where("id = ?", id).First(&dbDeck).Error; err != nil {
		return nil, err
	}

  if dbDeck.Remaining == 0 {
    return []*entities.Card{}, nil
  }

  // Ensure Array Bounds
  max := int(count)
  if int(count) > len(dbDeck.Cards) {
    max = len(dbDeck.Cards)
  }

  // Updating variables
  cards := dbDeck.Cards[0:max]
  dbDeck.Remaining -= max
  dbDeck.Cards = dbDeck.Cards[max:]

  // Updating deck
  if err := r.db.WithContext(ctx).Save(&dbDeck).Error; err != nil {
    return nil, err
  }

  return entities.ConvertCards(cards)
}
