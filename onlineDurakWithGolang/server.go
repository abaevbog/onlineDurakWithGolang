package main

import (
	"fmt"
	"log"
	"net/http"

	durak "github.com/abaevbog/onlineDurakWithGolang/durak"
	pl "github.com/abaevbog/onlineDurakWithGolang/player"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

const (
	TextMessage   = 1
	BinaryMessage = 2
	CloseMessage  = 8
	PingMessage   = 9
	PongMessage   = 10
)

func endpointFuncGenerator(ch chan<- pl.ClientChannels) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		upgrader.CheckOrigin = func(r *http.Request) bool { return true }

		ws, err := upgrader.Upgrade(w, r, nil)

		if err != nil {
			log.Println(err)
		}
		toClientCh := make(chan []byte)
		toServerCh := make(chan string)
		brokeConnCh := make(chan int, 1)
		attackerCh := make(chan pl.ChannelFields)
		defenderCh := make(chan pl.ChannelFields)
		newMessageCh := make(chan bool)
		newsChannel := make(chan pl.ChannelFields)
		ch <- pl.ClientChannels{PlayerBrokeConntection: brokeConnCh, ClientToServer: toServerCh, ServerToClient: toClientCh, AttackerChannel: attackerCh, DefenderChannel: defenderCh, ClienNeedsToReloadInfoChannel: newMessageCh, NewsChannel: newsChannel}
		go clientToServer(ws, toServerCh, brokeConnCh)
		go serverToClient(ws, toClientCh)
	}
}

//server reads messages from the client
func clientToServer(conn *websocket.Conn, channel chan string, brokeConn chan int) {
	for {

		typ, p, err := conn.ReadMessage()
		if typ == CloseMessage {
			fmt.Println("Client close message!")
			return
		}
		if err != nil {
			log.Println(err)
			fmt.Println("type: ", typ)
			fmt.Println("Connection broken!!")
			toClose := false
			select {
			case brokeConn <- 1:
				fmt.Println("sent to broken conn")
				toClose = true
			default:
				fmt.Println("not sent to broken conn")
				return
			}
			if toClose {
				close(brokeConn)
			}
			fmt.Println("client to server done")
			return
		}
		res := string(p)
		channel <- res

	}
}

// server writes messages to the client
func serverToClient(conn *websocket.Conn, channel chan []byte) {
	for {
		message := <-channel
		if string(message) == "Done" {
			fmt.Println("server to client done")

			conn.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			fmt.Println("server to client done")
			return
		}
		if err := conn.WriteMessage(TextMessage, []byte(message)); err != nil {
			fmt.Println("Error in write server to client")
			fmt.Println("server to client done")
		}
	}
}

func setup() {
	players := 0
	queueChannel := make(chan pl.ClientChannels)
	channelArr := make([]pl.ClientChannels, 2)
	wsEndpoint := endpointFuncGenerator(queueChannel)
	go func() {
		for {
			channels := <-queueChannel
			fmt.Println("new player joined")
			channelArr[players] = channels
			players = (players + 1) % 2
			if players == 0 {
				fmt.Println("new game launched!")
				go durak.LaunchGame(2, channelArr)
			}
		}
	}()
	http.HandleFunc("/api", wsEndpoint)
}

func main() {
	fmt.Println("Go!")
	setup()
	http.Handle("/", http.FileServer(http.Dir("./network/build/")))
	log.Fatal(http.ListenAndServe(":8080", nil))
}
