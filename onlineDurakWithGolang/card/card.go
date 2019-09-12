//go:generate stringer -type=Suit,Rank

package deck

import (
	"fmt"
	"math/rand"
	"time"
)

// Suit type, small int
type Suit uint8

//Suits defs
const (
	Spade   Suit = iota //0
	Diamond             //1
	Club                //2
	Heart               //3
)

//Rank type
type Rank uint8

//Rank defs
const (
	Six Rank = iota + 6
	Seven
	Eight
	Nine
	Ten
	Jack
	Queen
	King
	Ace
)

//Card type def
type Card struct {
	Suit
	Rank
}

func (c Card) String() string {
	return fmt.Sprintf("%s of %ss", c.Rank.String(), c.Suit.String())
}

//New card deck
func New() []Card {
	var cards []*Card
	suits := []Suit{Spade, Diamond, Club, Heart}
	for _, suit := range suits {
		for rank := 6; rank <= 14; rank++ {
			cards = append(cards, &Card{Suit: suit, Rank: Rank(rank)})
		}
	}
	return shuffle(cards)
}

//shuffle card deck
func shuffle(cards []*Card) []Card {
	ret := make([]Card, 36)
	rand.Seed(time.Now().UTC().UnixNano())
	for i := 0; i < 36; i++ {
		ind := rand.Intn(len(cards))
		ret[i] = *(cards[ind])
		cards = append(cards[:ind], cards[ind+1:]...)
	}
	return ret
}
