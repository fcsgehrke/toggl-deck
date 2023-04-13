package services

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/fcsgehrke/toggl-deck/mocks"
	"github.com/fcsgehrke/toggl-deck/pkg/deck"
	"github.com/gofrs/uuid"
)

func TestCreate(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  response, err := service.CreateDeck(ctx, false, nil)

  if err != nil {
    log.Println("Should not return error.")
    t.FailNow()
  }

  if response.Remaining != 52 {
    log.Println("Should be remanining 52 cards.")
    t.FailNow()
  }

  if response.Shuffled == true {
    log.Println("Should not be shuffled.")
    t.FailNow()
  }
}

func TestCreateShuffled(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  response, err := service.CreateDeck(ctx, true, nil)

  if err != nil {
    log.Println("Should not return error.")
    t.FailNow()
  }

  if response.Remaining != 52 {
    log.Println("Should be remanining 52 cards.")
    t.FailNow()
  }

  if response.Shuffled == false {
    log.Println("Should be shuffled.")
    t.FailNow()
  }
}

func TestCreateWithCards(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  response, err := service.CreateDeck(ctx, false, []string {"AS", "7H"})

  if err != nil {
    log.Println("Should not return error.")
    t.FailNow()
  }

  if response.Remaining != 2 {
    log.Println("Should be remanining 2 cards.")
    t.FailNow()
  }

  if response.Shuffled == true {
    log.Println("Should not be shuffled.")
    t.FailNow()
  }
}

func TestCreateWithCardsAndShuffle(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  response, err := service.CreateDeck(ctx, true, []string {"AS", "7H"})

  if err != nil {
    log.Println("Should not return error.")
    t.FailNow()
  }

  if response.Remaining != 2 {
    log.Println("Should be remanining 2 cards.")
    t.FailNow()
  }

  if response.Shuffled == false {
    log.Println("Should be shuffled.")
    t.FailNow()
  }
}

func TestCreateWithMoreCardsThanNeeded(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  _, err := service.CreateDeck(ctx, false, []string {"AS", "7H", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS", "AS"})

  if err == nil {
    log.Println("Should return error.")
    t.FailNow()
  }

  if err != deck.ErrInvalidDeckSize {
    log.Println("Should return ErrInvalidDeckSize")
    t.FailNow()
  }
}

func TestCreateWrongCardCode(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  _, err := service.CreateDeck(ctx, false, []string {"XS"})

  if err == nil {
    log.Println("Should return error.")
    t.FailNow()
  }

  if err != deck.ErrInvalidCardCode {
    log.Println("Should return ErrInvalidCardCode")
    t.FailNow()
  }
}

func TestOpen(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  deck, _ := service.CreateDeck(ctx, false, []string {"AS", "AC", "AD", "2C"})

  response, err := service.OpenDeck(ctx, deck.ID.String())
  if err != nil {
    log.Println("Should not return error.")
    t.FailNow()
  }

  if response.ID != deck.ID {
    log.Println("Shoud return same id.")
    t.FailNow()
  }

  if response.Remaining != deck.Remaining {
    log.Println("Should return 4 remaining cards.")
    t.FailNow()
  }

  if len(response.Cards) != response.Remaining {
    log.Println("Shoud have the same number of len(cards) and remaining.")
    t.FailNow()
  }

  if response.Shuffled != deck.Shuffled {
    log.Println("Shoud have the same shuffled value as the created deck.")
    t.FailNow()
  }
}

func TestOpenWithWrongId(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  _, err := service.OpenDeck(ctx, "UUID")
  if err != deck.ErrInvalidDeckID {
    log.Println("Should return ErrInvalidDeckID")
  }
}

func TestOpenNotFound(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)
  
  uid, _ := uuid.NewV4()

  _, err := service.OpenDeck(ctx, uid.String())
  if err != deck.ErrDeckNotFound {
    log.Println("Should return ErrDeckNotFound")
  }
}

func TestDraw(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  deck, _ := service.CreateDeck(ctx, false, []string {"AS", "AC", "AD", "2C"})

  response, err := service.DrawCard(ctx, deck.ID.String(), 1)
  if err != nil {
    log.Println("Should not return error.")
    t.FailNow()
  }

  if len(response.Cards) != 1 {
    log.Println("Should draw 1 card.")
    t.FailNow()
  }

  drawnDeck, _ := service.OpenDeck(ctx, deck.ID.String())
  if drawnDeck.Remaining != 3 {
    log.Println("Shoud remain 3 cards on the deck.")
    t.FailNow()
  }
}

func TestDrawInvalidID(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)

  _, err := service.DrawCard(ctx, "UUID", 1)
  if err != deck.ErrInvalidDeckID {
    log.Println("Should return ErrInvalidDeckID")
  }
}

func TestDrawNotFound(t *testing.T) {
  ctx := context.Background()
  log := log.New(os.Stderr, "", 0)
  repo := mocks.NewDeckRepoMock()
  service, _ := NewDeckService(repo, log)
  
  uid, _ := uuid.NewV4()

  _, err := service.DrawCard(ctx, uid.String(), 1)
  if err != deck.ErrDeckNotFound {
    log.Println("Should return ErrDeckNotFound")
  }
}

