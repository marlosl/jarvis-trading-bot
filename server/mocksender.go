package server

import (
	"bufio"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"

	"github.com/gorilla/websocket"
)

func StartMockSenderClient() {
	messageOut := make(chan string)
	interrupt := make(chan os.Signal, 1)

	wssHost := utils.GetStringConfig(consts.WssMockHost)
	wssPort := utils.GetStringConfig(consts.WssMockPort)
	wssServer := fmt.Sprintf("%s:%s", wssHost, wssPort)

	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "ws", Host: wssServer, Path: "/stream", RawQuery: "streams=ethbusd@kline_1m/btcbusd@kline_1m"}

	log.InfoLogger.Printf("connecting to %s", u.String())
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		log.ErrorLogger.Fatal("dial:", err)
	}

	if resp != nil {
		log.ErrorLogger.Printf("handshake failed with status %d", resp.StatusCode)
	}

	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			importFile := os.Getenv("IMPORT_FILE")
			log.InfoLogger.Printf("opening file %s\n", importFile)

			f, err := os.Open(importFile)
			if err != nil {
				panic(fmt.Sprintf("error opening file: %v\n", err))
			}

			r := bufio.NewReader(f)
			s, e := utils.Readln(r)

			for e == nil {
				messageOut <- s
				s, e = utils.Readln(r)
			}
		}
	}()

	for {
		select {
		case <-done:
			return
		case m := <-messageOut:
			log.InfoLogger.Printf("Send Message M %s", m)
			err := c.WriteMessage(websocket.TextMessage, []byte(m))
			if err != nil {
				log.InfoLogger.Println("write:", err)
				return
			}
		case <-interrupt:
			log.InfoLogger.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				log.InfoLogger.Println("write close:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}
