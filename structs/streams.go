package structs

type Kline struct {
	KlineStartTime           int64  `json:"t"`
	KlineCloseTime           int64  `json:"T"`
	Symbol                   string `json:"s"`
	Interval                 string `json:"i"`
	FirstTradeID             int64  `json:"f"`
	LastTradeID              int64  `json:"L"`
	OpenPrice                string `json:"o"`
	ClosePrice               string `json:"c"`
	HighPrice                string `json:"h"`
	LowPrice                 string `json:"l"`
	BaseAssetVolume          string `json:"v"`
	NumberOfTrades           int64  `json:"n"`
	IsThisKlineClosed        bool   `json:"x"`
	QuoteAssetVolume         string `json:"q"`
	TakerBuyBaseAssetVolume  string `json:"V"`
	TakerBuyQuoteAssetVolume string `json:"Q"`
	Ignore                   string `json:"B"`
}

type Event struct {
	EventType string `json:"e"`
	EventTime int64  `json:"E"`
	Symbol    string `json:"s"`
	Kline     Kline  `json:"k"`
}

type Stream struct {
	Stream string `json:"stream"`
	Data   Event  `json:"data"`
}
