package main

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	finnhub "github.com/Finnhub-Stock-API/finnhub-go/v2"
	"github.com/gorilla/websocket"
)

const apiKey = "cjufgepr01qlodk35ebgcjufgepr01qlodk35ec0"

type SocketMsg struct {
	Type string  `json:"type"`
	Data []Trade `json:"data"`
}

type Trade struct {
	Price  float32 `json:"p"`
	Symbol string  `json:"s"`
	TimeMs float32 `json:"t"`
	Volume float32 `json:"v"`
}

type Quote struct {
	price          float32
	change         float32
	percent_change float32
	high           float32
	low            float32
	open           float32
	previous_close float32
}

func newQuote(res finnhub.Quote) *Quote {
	q := Quote{
		price:          *res.C,
		change:         *res.D,
		percent_change: *res.Dp,
		high:           *res.H,
		low:            *res.L,
		open:           *res.O,
		previous_close: *res.Pc,
	}
	return &q
}

func getVoo() {
	cfg := finnhub.NewConfiguration()
	cfg.AddDefaultHeader("X-Finnhub-Token", apiKey)
	finnhubClient := finnhub.NewAPIClient(cfg).DefaultApi

	res, _, _ := finnhubClient.Quote(context.Background()).Symbol("VOO").Execute()
	q := newQuote(res)

	fmt.Printf("%+v\n", *q)
}

func openFinnhubWs() {
	headers := make(http.Header)
	headers.Add("X-Finnhub-Token", apiKey)
	var wsUrl = "wss://ws.finnhub.io?token=" + apiKey
	w, _, err := websocket.DefaultDialer.Dial(wsUrl, headers)
	if err != nil {
		panic(err)
	}
	defer w.Close()

	symbols := []string{"VOO", "SPY", "F", "O", "REXR", "ARE", "WBD"}
	for _, s := range symbols {
		msg, _ := json.Marshal(map[string]interface{}{"type": "subscribe", "symbol": s})
		w.WriteMessage(websocket.TextMessage, msg)
	}

	var msg SocketMsg
	for {
		err := w.ReadJSON(&msg)
		if err != nil {
			panic(err)
		}

		if msg.Type == "trade" {
			fmt.Printf("Trade made: %+v\n", msg)
		} else {
			fmt.Printf("Other: %+v\n", msg)
		}
	}
}
