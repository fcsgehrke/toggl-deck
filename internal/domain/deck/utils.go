package deck

import (
	"math/rand"
	"strconv"
	"time"

	"github.com/fcsgehrke/toggl-deck/pkg/deck"
)

var (
	suits  = []string{"S", "D", "C", "H"}
	values = []string{"A", "2", "3", "4", "5", "6", "7", "8", "9", "0", "J", "Q", "K"}
)

func ConvertCardCodeIntoCard(code string) (*Card, error) {
	if len(code) != 2 {
		return nil, deck.ErrInvalidCardCode
	}

	card := &Card{}
	number := code[0]
	suit := code[1]

	switch number {
	case 'K':
		card.Value = "KING"
	case 'Q':
		card.Value = "QUEEN"
	case 'J':
		card.Value = "JACK"
	case 'A':
		card.Value = "ACE"
	default:
		nbr, err := strconv.ParseInt(string(number), 10, 8)
		if err != nil {
			return nil, deck.ErrInvalidCardCode
		}

		if nbr == 1 {
			return nil, deck.ErrInvalidCardCode
		}

		card.Value = string(number)
	}

	switch suit {
	case 'S':
		card.Suit = Spades
	case 'D':
		card.Suit = Diamonds
	case 'C':
		card.Suit = Clubs
	case 'H':
		card.Suit = Hearts
	default:
		return nil, deck.ErrInvalidCardCode
	}

	card.Code = code

	return card, nil
}

func ConvertCards(cards []string) ([]*Card, error) {

	if len(cards) > 52 {
		return nil, deck.ErrInvalidDeckSize
	}

	crds := []*Card{}

	for _, code := range cards {
		card, err := ConvertCardCodeIntoCard(code)
		if err != nil {
			return nil, err
		}
		crds = append(crds, card)
	}
	return crds, nil
}

func GenerateCards() []*Card {
	cards := []*Card{}

	for _, suit := range suits {
		for _, value := range values {
			card, _ := ConvertCardCodeIntoCard(value + suit)
			cards = append(cards, card)
		}
	}

	return cards
}

func ShuffleCards(cards []*Card) []*Card {
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(cards), func(i, j int) {
		cards[i], cards[j] = cards[j], cards[i]
	})
	return cards
}
