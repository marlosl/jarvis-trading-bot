package server

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"jarvis-trading-bot/analyzer"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/notification"
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/segmentio/kafka-go"
)

type ListenerServer struct {
	Analyzer     *analyzer.Analyzer
	KafkaCluster []string
	TopicName    string
	Reader       *kafka.Reader
	Writer       *kafka.Writer
	Ctx          context.Context
	Simulation   bool
}

func (l *ListenerServer) Init() {
	l.KafkaCluster = utils.GetStringSliceConfig(consts.KafkaCluster)
	l.TopicName = utils.GetStringConfig(consts.TopicName)
	l.Ctx = context.Background()
}

func (l *ListenerServer) InitReader() {
	l.Init()

	l.Analyzer = new(analyzer.Analyzer)

	if l.Simulation {
		l.Analyzer.InitSimulation()
	} else {
		l.Analyzer.Init()
	}

	l.Reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers: l.KafkaCluster,
		Topic:   l.TopicName,
		GroupID: "group-1",
	})
}

func (l *ListenerServer) InitWriter() {
	l.Init()

	l.Writer = kafka.NewWriter(kafka.WriterConfig{
		Brokers: l.KafkaCluster,
		Topic:   l.TopicName,
	})
}

func (l *ListenerServer) StartListenerServer() {
	l.InitWriter()

	interrupt := make(chan os.Signal, 1)
	done := make(chan struct{})

	signal.Notify(interrupt, os.Interrupt)

	streamHost := utils.GetStringConfig(consts.WssBinanceUrl)
	streams := l.getStreams()

	signal.Notify(interrupt, os.Interrupt)
	u := url.URL{Scheme: "wss", Host: streamHost, Path: "/stream", RawQuery: fmt.Sprintf("streams=%s", streams)}

	log.InfoLogger.Printf("connecting to %s", u.String())
	c, _, err := websocket.DefaultDialer.Dial(u.String(), nil)

	if err != nil {
		notification.SendMessage(fmt.Sprintf("dial: %s", err), true)
		log.ErrorLogger.Fatal("dial:", err)
	}

	defer c.Close()
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				notification.SendMessage(fmt.Sprintf("read: %s", err), true)
				log.InfoLogger.Println("read:", err)
				os.Exit(1)
				return
			}

			stream := new(structs.Stream)
			if err := json.Unmarshal(message, &stream); err != nil {
				notification.SendMessage(fmt.Sprintf("read: %s", err), true)
				log.InfoLogger.Println("read:", err)
				os.Exit(2)
				return
			}

			if stream.Data.Kline.IsThisKlineClosed {
				log.InfoLogger.Printf("recv: %s", message)
				candle := utils.ConvertStreamToCandle(stream)

				id := uuid.New()
				err := l.Writer.WriteMessages(l.Ctx, kafka.Message{
					Key:   []byte(id.String()),
					Value: []byte(utils.SPrintJson(candle)),
				})
				if err != nil {
					log.ErrorLogger.Fatal("could not write message " + err.Error())
				}
			}
		}
	}()

	for {
		select {
		case <-done:
			return

		case <-interrupt:
			log.InfoLogger.Println("interrupt")
			err := c.WriteMessage(websocket.CloseMessage,
				websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			)
			if err != nil {
				notification.SendMessage(fmt.Sprintf("write close: %s", err), true)
				log.InfoLogger.Println("write close:", err)
				os.Exit(3)
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

func (l *ListenerServer) getStreams() string {
	symbols := utils.GetStringSliceConfig(consts.StreamsSymbols)
	intervals := utils.GetStringSliceConfig(consts.StreamsIntervals)

	streams := ""

	for i, s := range symbols {
		interval := ""
		if i < len(intervals) {
			interval = intervals[i]
		}

		if i > 0 {
			streams += "/"
		}

		streams += fmt.Sprintf("%s@kline_%s", s, interval)
	}
	return streams
}

func (l *ListenerServer) StartAnalyzerServer() {
	l.InitReader()

	for {
		msg, err := l.Reader.ReadMessage(l.Ctx)
		if err != nil {
			log.ErrorLogger.Fatal("could not read message " + err.Error())
		}
		candle := new(structs.Candlestick)
		if err := json.Unmarshal(msg.Value, &candle); err != nil {
			notification.SendMessage(fmt.Sprintf("read: %s", err), true)
			log.InfoLogger.Println("read:", err)
			return
		}

		if !l.Simulation {
			utils.DB.Create(candle)
		}
		l.Analyzer.Process(candle)

		log.InfoLogger.Println("received: ", string(msg.Value))
	}
}

func ServerStartAndListen() {
	s := new(ListenerServer)
	s.StartListenerServer()
}

func AnalyzerServerStart() {
	s := new(ListenerServer)
	s.StartAnalyzerServer()
}

func AnalyzerServerSimulationStart() {
	s := new(ListenerServer)
	s.Simulation = true
	s.StartAnalyzerServer()
}
