package player

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"

	deck "github.com/abaevbog/onlineDurakWithGolang/card"
)

//ChannelFields has fields needed for the defender and attacker
type ChannelFields struct {
	CardsBeaten    []deck.Card
	CardsNotBeaten []deck.Card
	TookCards      bool
	DoneWithAttack bool
	From           int
}

//ClientChannels struct has channels to server from client and vice versa
type ClientChannels struct {
	ClientToServer                chan string
	ServerToClient                chan []byte
	AttackerChannel               chan ChannelFields
	DefenderChannel               chan ChannelFields
	ClienNeedsToReloadInfoChannel chan bool //reprint info to the client
	NewsChannel                   chan ChannelFields
	PlayerBrokeConntection        chan int
}

//Player struct with hand of cards
type Player struct {
	HandOfCards             []deck.Card
	AvailableDefenceOptions map[int][]int //list of maps: key = id of the available option,value = list of ids of non beaten cards on the table
	ID                      int
	Channels                ClientChannels
}

//PrintStatus prints  the status of the game + deck of the player
func (p Player) PrintStatus(notBeaten []deck.Card, beaten []deck.Card, attack bool, kozir deck.Suit) {

	var status = make(map[string]map[string][]deck.Card)
	var statusData = make(map[string][]deck.Card)
	var options = make(map[string]map[int][]int)
	options["options"] = p.AvailableDefenceOptions

	statusData["kozir"] = []deck.Card{deck.Card{Rank: 0, Suit: kozir}}
	statusData["nonBeaten"] = notBeaten
	statusData["beaten"] = beaten
	statusData["deck"] = p.HandOfCards
	if attack {
		status["attack"] = statusData
	} else {
		status["defence"] = statusData
	}

	statusEnc, err := json.Marshal(status)
	optionsEnc, err2 := json.Marshal(options)
	if err != nil || err2 != nil {
		fmt.Println(err)
		fmt.Println(err2)
	}
	p.Channels.ServerToClient <- statusEnc
	p.Channels.ServerToClient <- optionsEnc
}

//PrintDefenceOptions prints options for defender
func (p Player) PrintDefenceOptions(notBeaten []deck.Card, selectedOption int) {
	var options = make(map[string][]int)
	options["defenceChoices"] = p.AvailableDefenceOptions[selectedOption]
	fmt.Println("printing defence choices")
	opt, err := json.Marshal(options)
	if err != nil {
		fmt.Println(err)
	}
	p.Channels.ServerToClient <- opt
}

func contains(arr map[int][]int, val int) bool {
	_, ok := arr[val]
	return ok
}

func containsArr(arr []int, val int) bool {
	for _, v := range arr {
		if v == val {
			return true
		}
	}
	return false
}

//PutCard puts card by a player
func (p *Player) PutCard(cardID int) deck.Card {
	result := p.HandOfCards[cardID]
	(*p).HandOfCards = append(p.HandOfCards[:cardID], p.HandOfCards[cardID+1:]...)
	return result
}

//CloseAllChannels closes all players' channels
func (p *Player) CloseAllChannels() {
	fmt.Println("CLOSE ALL CHANNELS")
	p.Channels.ServerToClient <- []byte("Done")
	close(p.Channels.AttackerChannel)
	close(p.Channels.ClienNeedsToReloadInfoChannel)
	close(p.Channels.ClientToServer)
	close(p.Channels.DefenderChannel)
	close(p.Channels.NewsChannel)
	close(p.Channels.ServerToClient)
}

//SetAvailableOptions for the players move
func (p *Player) SetAvailableOptions(attack bool, cards []deck.Card, kozir deck.Suit) {
	if attack {
		inArr := func(cards []deck.Card, card deck.Card) bool {
			for _, v := range cards {
				if v.Rank == card.Rank {
					return true
				}
			}
			return false
		}
		options := map[int][]int{}
		for i, v := range p.HandOfCards {
			if len(cards) == 0 || inArr(cards, v) {
				options[i] = []int{}
			}
		}
		(*p).AvailableDefenceOptions = options
	} else {
		cardsBeatable := func(notBeaten []deck.Card, card deck.Card) []int {
			result := []int{}
			for i, v := range notBeaten {
				if (v.Rank < card.Rank && v.Suit == card.Suit) || (card.Suit == kozir && v.Suit != kozir) {
					result = append(result, i)
				}
			}
			return result
		}
		options := map[int][]int{}
		for i, v := range p.HandOfCards {
			cardMatches := cardsBeatable(cards, v)
			if len(cardMatches) > 0 {
				options[i] = cardMatches
			}
		}
		(*p).AvailableDefenceOptions = options
	}
}

//PlayGame defines the player's behavior during the game in response to
//messages from the control panel's channels
func PlayGame(ctx context.Context, p *Player, kozir deck.Suit) {
	var mutex sync.Mutex
	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in player ", r)
		}
	}()

	for {
		select {
		case cardsOnTheTableStruct := <-p.Channels.AttackerChannel:
			cardsNotBeaten, cardsBeaten := cardsOnTheTableStruct.CardsNotBeaten, cardsOnTheTableStruct.CardsBeaten
			cardsOnTheTable := append(cardsBeaten, cardsNotBeaten...)
			mutex.Lock()
			p.SetAvailableOptions(true, cardsOnTheTable, kozir)
			mutex.Unlock()
			p.PrintStatus(cardsNotBeaten, cardsBeaten, true, kozir)
			willBreak := false
			for {
				select {
				case cardIdy := <-p.Channels.ClientToServer:
					cardID, _ := strconv.Atoi(cardIdy)
					//cardID = cardID - 48
					fmt.Println("REceived attack: ", cardID)
					if contains(p.AvailableDefenceOptions, cardID) {
						fmt.Println("card placed")
						card := (*p).PutCard(int(cardID))
						p.Channels.AttackerChannel <- ChannelFields{CardsBeaten: cardsBeaten, CardsNotBeaten: append(cardsNotBeaten, card), DoneWithAttack: false, From: cardsOnTheTableStruct.From}
						//p.Channels.ServerToClient <- []byte("You placed " + card.String())
						willBreak = true
					} else if cardID == 1000 && len(cardsOnTheTable) != 0 {
						fmt.Println("finished")
						p.Channels.AttackerChannel <- ChannelFields{CardsBeaten: cardsBeaten, CardsNotBeaten: cardsNotBeaten, DoneWithAttack: true, From: cardsOnTheTableStruct.From}
						//p.Channels.ServerToClient <- []byte("Your attack is complete")
						willBreak = true
					} else {
						//p.Channels.ServerToClient <- []byte("Please select the correct option")
						fmt.Println("not correct")
						willBreak = false
					}
				case <-p.Channels.ClienNeedsToReloadInfoChannel:
					willBreak = true

				case <-ctx.Done():
					fmt.Println("player: inside attack done!")
					return
				}
				if willBreak {
					break
				}
			}

		case cardsOnTheTableStruct := <-p.Channels.DefenderChannel:
			cardsNotBeaten, cardsBeaten := cardsOnTheTableStruct.CardsNotBeaten, cardsOnTheTableStruct.CardsBeaten
			mutex.Lock()
			p.SetAvailableOptions(false, cardsNotBeaten, kozir)
			mutex.Unlock()
			p.PrintStatus(cardsNotBeaten, cardsBeaten, false, kozir)
			willBreak := false
			for {
				select {
				case cardIdy := <-p.Channels.ClientToServer:
					cardID, _ := strconv.Atoi(cardIdy)
					if contains(p.AvailableDefenceOptions, cardID) {
						fmt.Println("defence: card chosen", cardID)
						card := (*p).PutCard(cardID)
						p.PrintDefenceOptions(cardsNotBeaten, cardID)
						cardToBeatStr := <-p.Channels.ClientToServer
						cardToBeat, _ := strconv.Atoi(cardToBeatStr)
						fmt.Println("AAaaaa")
						fmt.Println(cardToBeat)
						fmt.Println("AAaaaa")
						if containsArr(p.AvailableDefenceOptions[cardID], cardToBeat) {
							fmt.Println("defence: card to beat chosen")
							fmt.Println(cardID)
							fmt.Println(cardToBeat)
							fmt.Println(p.AvailableDefenceOptions)
							fmt.Println("+++")
							ind := p.AvailableDefenceOptions[cardID][cardToBeat]
							p.Channels.DefenderChannel <- ChannelFields{CardsBeaten: append(cardsBeaten, card, cardsNotBeaten[ind]), CardsNotBeaten: append(cardsNotBeaten[:ind], cardsNotBeaten[ind+1:]...), TookCards: false, From: cardsOnTheTableStruct.From}
							willBreak = true
						} else {
							fmt.Println("defence: not correct. Card to beat", cardToBeat)
							willBreak = false
						}
					} else if cardID == 1000 {
						p.Channels.ServerToClient <- []byte("taken")
						onTheTable := append(cardsBeaten, cardsNotBeaten...)
						totalDeck := append(onTheTable, (*p).HandOfCards...)
						(*p).HandOfCards = totalDeck
						p.Channels.DefenderChannel <- ChannelFields{CardsBeaten: cardsBeaten, CardsNotBeaten: cardsNotBeaten, TookCards: true, From: cardsOnTheTableStruct.From}
						willBreak = true
					} else {
						fmt.Println("defence: not correct2", cardID)
						willBreak = false
					}

				case <-p.Channels.ClienNeedsToReloadInfoChannel:
					willBreak = true

				case <-ctx.Done():
					fmt.Println("player: inside defence done!")
					return
				}
				if willBreak {
					break
				}
			}
		case <-ctx.Done():
			fmt.Println("player: outside done!")
			return
		}
	}
}
