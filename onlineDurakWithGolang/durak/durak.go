package durak

import (
	"context"
	"encoding/json"
	"fmt"
	"runtime"
	"sync"
	"time"

	deck "github.com/abaevbog/onlineDurakWithGolang/card"
	pl "github.com/abaevbog/onlineDurakWithGolang/player"
)

//Game struct with list of players
type Game struct {
	cardsLeft      []deck.Card
	players        []pl.Player
	kozir          deck.Suit
	cardsBeaten    []deck.Card
	cardsNotBeaten []deck.Card
	attack         bool
	firstAttack    bool
	attackerID     int
}

func (g Game) shouldContinue(ctx context.Context) bool {
	select {
	case <-ctx.Done():

		return false
	default:
		onePLayersHandNonEmpty := false
		for _, v := range g.players {
			if len(v.HandOfCards) > 0 {
				if !onePLayersHandNonEmpty {
					onePLayersHandNonEmpty = true
				} else {
					return true
				}
			}
		}
		return false
	}

}

//NotifyAllPlayers sends a message to all players
func (g Game) NotifyAllPlayers(message string) {
	for _, p := range g.players {
		options := make(map[string]string)
		options["message"] = message
		mesEnc, err := json.Marshal(options)
		if err != nil {
			fmt.Println(err)
		}
		p.Channels.ServerToClient <- []byte(mesEnc)
	}
}

func (g *Game) refillCards() {
	for i := range g.players {
		needed := 6 - len(g.players[i].HandOfCards)
		for len((*g).cardsLeft) > 0 && needed > 0 {
			(*g).players[i].HandOfCards = append((*g).players[i].HandOfCards, g.cardsLeft[0])
			(*g).cardsLeft = (*g).cardsLeft[1:]
			needed--
		}
	}
}

func (g *Game) controlPanel(ctx context.Context) {
	for ind := range g.players {
		fmt.Println("player ", ind, "in control")
		//g.players[ind].Channels.ServerToClient <- []byte(message)
		go pl.PlayGame(ctx, &(g.players[ind]), g.kozir)
	}
	var mainUpdatesChan = make(chan pl.ChannelFields)

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Recovered in durak ", r)
		}
	}()

	go func(ctx context.Context, mainChan chan pl.ChannelFields) {
		for {
			select {
			case news := <-mainChan:
				for i := range g.players {
					if i != news.From {
						g.players[i].Channels.NewsChannel <- news
					}
				}
			case <-ctx.Done():
				close(mainUpdatesChan)
				return
			}

		}
	}(ctx, mainUpdatesChan)

	for g.shouldContinue(ctx) {
		var waitgroup sync.WaitGroup
		var mutex sync.Mutex
		waitgroup.Add(len(g.players))
		idOfDefender := (g.attackerID + 1) % len(g.players)
		for i := range g.players {
			go func(i int, waitgroup *sync.WaitGroup, mtex *sync.Mutex, newsChannel chan pl.ChannelFields) {
				if i == g.attackerID {
					for {
						doWeBreak := false
						clientNeedsToReload := false
						(*g).players[i].Channels.AttackerChannel <- pl.ChannelFields{CardsBeaten: g.cardsBeaten, CardsNotBeaten: g.cardsNotBeaten, From: i}
						select {
						case response := <-g.players[i].Channels.AttackerChannel:

							// attacked successfully
							if !response.DoneWithAttack {
								mtex.Lock()
								(*g).cardsNotBeaten = response.CardsNotBeaten
								mtex.Unlock()
								newsChannel <- response
								//attack complete
							} else {
								if len((*g).cardsNotBeaten) > 0 {
									panic("Attack was finished but there are unbeaten cards")
								}
								newsChannel <- response
								mtex.Lock()
								(*g).attackerID = (g.attackerID + 1) % len(g.players)
								(*g).cardsBeaten = []deck.Card{}
								mtex.Unlock()
								doWeBreak = true
								clientNeedsToReload = false
							}
						case news := <-g.players[i].Channels.NewsChannel:
							clientNeedsToReload = true
							if news.TookCards || news.DoneWithAttack {
								doWeBreak = true
							}
						case <-ctx.Done():

							waitgroup.Done()
							return
						}
						if clientNeedsToReload {
							go func() {
								select {
								case g.players[i].Channels.ClienNeedsToReloadInfoChannel <- true:

								case <-ctx.Done():

									waitgroup.Done()
									return
								}
							}()
						}
						if doWeBreak {
							break
						}
					}
				} else if i == idOfDefender {
					for {
						doWeBreak := false
						clientNeedsToReload := false
						g.players[i].Channels.DefenderChannel <- pl.ChannelFields{CardsBeaten: g.cardsBeaten, CardsNotBeaten: g.cardsNotBeaten, From: i}
						select {
						case response := <-g.players[i].Channels.DefenderChannel:
							if !response.TookCards {
								mtex.Lock()
								(*g).cardsBeaten = response.CardsBeaten
								(*g).cardsNotBeaten = response.CardsNotBeaten
								mtex.Unlock()
								newsChannel <- response
							} else {
								mtex.Lock()
								(*g).cardsNotBeaten = []deck.Card{}
								(*g).cardsBeaten = []deck.Card{}
								(*g).attackerID = ((*g).attackerID + 2) % len(g.players)
								mtex.Unlock()
								newsChannel <- response
								doWeBreak = true
								clientNeedsToReload = false
							}

						case news := <-g.players[i].Channels.NewsChannel:
							clientNeedsToReload = true
							if news.TookCards || news.DoneWithAttack {
								doWeBreak = true
							}
						case <-ctx.Done():

							waitgroup.Done()
							return
						}
						if clientNeedsToReload {
							go func() {
								select {
								case g.players[i].Channels.ClienNeedsToReloadInfoChannel <- true:

								case <-ctx.Done():

									waitgroup.Done()
									return
								}
							}()
						}
						if doWeBreak {
							break
						}
					}
				}
				waitgroup.Done()
			}(i, &waitgroup, &mutex, mainUpdatesChan)
		}
		waitgroup.Wait()
		(*g).refillCards()
	}
}

//LaunchGame launches one game
func LaunchGame(playersNum int, clientChannels []pl.ClientChannels) {
	cards := deck.New()
	ctx := context.Background()
	ctxWithCancel, cancelFunction := context.WithCancel(ctx)
	game := Game{
		players:        make([]pl.Player, playersNum, playersNum),
		cardsLeft:      cards,
		kozir:          cards[len(cards)-1].Suit,
		cardsBeaten:    []deck.Card{},
		cardsNotBeaten: []deck.Card{},
		attack:         true,
		firstAttack:    true,
		attackerID:     0,
	}
	for i := range game.players {
		game.players[i].HandOfCards = game.cardsLeft[:6]
		game.players[i].AvailableDefenceOptions = map[int][]int{0: {}, 1: {}, 2: {}, 3: {}, 4: {}, 5: {}}
		game.players[i].Channels = clientChannels[i]
		game.players[i].ID = i
		game.cardsLeft = game.cardsLeft[6:]
	}

	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in player ", r)
			}
		}()
		select {
		case <-game.players[0].Channels.PlayerBrokeConntection:
			fmt.Println("Player 0 broke connection")
			game.NotifyAllPlayers("Other player broke connection")
			cancelFunction()
			return
		case <-game.players[1].Channels.PlayerBrokeConntection:
			fmt.Println("Player 1 broke connection")
			game.NotifyAllPlayers("Other player broke connection")
			cancelFunction()
			return
		}
	}()
	defer func() {
		fmt.Println("running defered")
		game.players[0].CloseAllChannels()
		game.players[1].CloseAllChannels()
		time.Sleep(3 * time.Second)
		fmt.Println("Almost done! Number of goroutines: ", runtime.NumGoroutine())
		if r := recover(); r != nil {
			fmt.Println("Recovered ", r)
		}
	}()

	(&game).controlPanel(ctxWithCancel)

	(&game).NotifyAllPlayers("Done!")

	return
}
