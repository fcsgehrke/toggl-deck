package deck

import (
	"log"
	"testing"

	"github.com/fcsgehrke/toggl-deck/pkg/deck"
)

func TestConvertCardCodeIntoCard(t *testing.T) {
  card, _ := ConvertCardCodeIntoCard("AS")

  if card.Value != "ACE" {
    log.Println("The value of the card shoud be ACE.")
    t.FailNow()
  }

  if card.Suit != Spades {
    log.Println("The suit of the should be Spades.")
    t.FailNow()
  }

  if card.Code != "AS" {
    log.Println("The code of the card should be AS.")
    t.FailNow()
  }
}

func TestConvertCardCodeIntoCardWrongCode(t *testing.T) {
  _, err := ConvertCardCodeIntoCard("AX")

  if err != deck.ErrInvalidCardCode {
    log.Println("The code of the card should be invalid.")
    t.FailNow()
  }
}

func TestConvertCards(t *testing.T) {
  cards, err := ConvertCards([]string{"AS", "2C"})
  if err != nil {
    log.Println("Should not return any error.")
    t.FailNow()
  }

  if len(cards) != 2 {
    log.Println("Should return exacltly 2 cards.")
    t.FailNow()
  }
}

func TestGenerateCards(t *testing.T) {
  cards := GenerateCards()

  if len(cards) != 52 {
    log.Println("Should return exacltly 52 cards.")
    t.FailNow()
  }
}

func TestShuffleCards(t *testing.T) {
  var shuffled []*Card

  cards := GenerateCards()
  shuffled = append(shuffled, cards...)
  
  shuffled = ShuffleCards(shuffled)

  equal := true
  for idx := range cards {
    if cards[idx].Code != shuffled[idx].Code {
      equal = false
      break
    }
  }

  if equal {
    log.Println("Should be shuffled.")
    t.FailNow()
  }
}
