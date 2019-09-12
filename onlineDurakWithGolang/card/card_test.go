package deck

import (
	"fmt"
	"testing"
)

func ExampleCard() {
	fmt.Println(Card{Rank: Ace, Suit: Spade})
	fmt.Println(Card{Rank: Six, Suit: Club})
	fmt.Println(Card{Rank: Ten, Suit: Diamond})
	fmt.Println(Card{Rank: Jack, Suit: Heart})

	//Output:
	//Ace of Spades
	//Six of Clubs
	//Ten of Diamonds
	//Jack of Hearts

}

func TestDeck(t *testing.T) {
	cards := New()

	fmt.Println(cards)
	if len(cards) != 36 {
		t.Error("Wrong num of cards")
	}
}

func TestShuffle(t *testing.T) {
	cards := New()
	fmt.Println(cards)
}
