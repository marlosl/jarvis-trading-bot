package structs

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/gorm"
)

type Candlestick struct {
	gorm.Model
	EventTime       time.Time       `gorm:"column:event_time"`
	Symbol          string          `gorm:"column:symbol"`
	OpenPrice       decimal.Decimal `gorm:"column:open_price" sql:"type:decimal(20,10);"`
	ClosePrice      decimal.Decimal `gorm:"column:close_price" sql:"type:decimal(20,10);"`
	HighPrice       decimal.Decimal `gorm:"column:high_price" sql:"type:decimal(20,10);"`
	LowPrice        decimal.Decimal `gorm:"column:low_price" sql:"type:decimal(20,10);"`
	BaseAssetVolume decimal.Decimal `gorm:"column:volume" sql:"type:decimal(20,10);"`
	NumberOfTrades  int64           `gorm:"column:trades_number"`
}

type TradingStatus struct {
	UserId            uint            `gorm:"primaryKey;column:user_id"`
	BotParameterId    uint            `gorm:"primaryKey;column:bot_param_id"`
	InstanceId        uint            `gorm:"primaryKey;column:instance_id"`
	Symbol            string          `gorm:"primaryKey;column:symbol"`
	LastStatus        int             `gorm:"column:last_status"`
	BuyAmount         decimal.Decimal `gorm:"column:buy_amount" sql:"type:decimal(20,10);"`
	InitialBuyAmount  decimal.Decimal `gorm:"column:initial_buy_amount" sql:"type:decimal(20,10);"`
	ProfitAmount      decimal.Decimal `gorm:"column:profit_amount" sql:"type:decimal(20,10);"`
	BuyTax            decimal.Decimal `gorm:"column:buy_tax" sql:"type:decimal(20,10);"`
	SellTax           decimal.Decimal `gorm:"column:sell_tax" sql:"type:decimal(20,10);"`
	TotalProfit       decimal.Decimal `gorm:"column:total_profit" sql:"type:decimal(20,10);"`
	TotalTaxes        decimal.Decimal `gorm:"column:total_taxes" sql:"type:decimal(20,10);"`
	LastEvent         time.Time       `gorm:"column:last_event"`
	LastOperationTime time.Time       `gorm:"column:last_operation_time"`
	Simulation        bool            `gorm:"column:simulation"`
	CreatedAt         time.Time
	UpdatedAt         time.Time
	DeletedAt         gorm.DeletedAt `gorm:"index"`
}

type Operation struct {
	Id                  uint            `gorm:"primaryKey;autoIncrement;column:id"`
	BotParameterId      uint            `gorm:"column:bot_param_id"`
	InstanceId          uint            `gorm:"column:instance_id"`
	Symbol              string          `gorm:"column:symbol"`
	Operation           string          `gorm:"column:operation"`
	BaseAsset           string          `gorm:"column:base_asset"`
	BasePrice           decimal.Decimal `gorm:"column:base_price" sql:"type:decimal(20,10);"`
	OrderId             int64           `gorm:"column:order_id"`
	OrigQty             decimal.Decimal `gorm:"column:orig_qty" sql:"type:decimal(20,10);"`
	ExecutedQty         decimal.Decimal `gorm:"column:executed_qty" sql:"type:decimal(20,10);"`
	CummulativeQuoteQty decimal.Decimal `gorm:"column:cummul_quote_qty" sql:"type:decimal(20,10);"`
	CommissionBase      decimal.Decimal `gorm:"column:commission_base" sql:"type:decimal(20,10);"`
	Commission          decimal.Decimal `gorm:"column:commission" sql:"type:decimal(20,10);"`
	CommissionAsset     string          `gorm:"column:commission_asset"`
	Type                string          `gorm:"column:type"`
	Status              string          `gorm:"column:status"`
	Opened              time.Time       `gorm:"column:opened"`
	TransactTime        time.Time       `gorm:"column:transact_time"`
	Finished            bool            `gorm:"column:finished"`
	Simulation          bool            `gorm:"column:simulation"`
}

type User struct {
	Id       uint   `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	Username string `gorm:"column:username" json:"username"`
	Password string `gorm:"column:password" json:"password"`
	Email    string `gorm:"column:email" json:"email"`
	Name     string `gorm:"column:name" json:"name"`
	IsAdmin  bool   `gorm:"column:is_admin" json:"isAdmin"`
}

type Parameters struct {
	UserId                        uint            `gorm:"primaryKey;column:user_id" json:"userId"`
	DefaultPercentageTax          decimal.Decimal `gorm:"column:default_percentage_tax" sql:"type:decimal(20,10);" json:"defaultPercentageTax"`
	DefaultStopLossPercentage     decimal.Decimal `gorm:"column:default_stoploss_percentage" sql:"type:decimal(20,10);" json:"defaultStopLossPercentage"`
	DefaultMinimumLimitPercentage decimal.Decimal `gorm:"column:default_minimumLimit_percentage" sql:"type:decimal(20,10);" json:"defaultMinimumLimitPercentage"`
	DefaultRsiPeriod              int             `gorm:"column:default_rsi_period" json:"defaultRsiPeriod"`
	DefaultRsiOverbought          decimal.Decimal `gorm:"column:default_rsi_overbought" sql:"type:decimal(20,10);" json:"defaultRsiOverbought"`
	DefaultRsiOversold            decimal.Decimal `gorm:"column:default_rsi_oversold" sql:"type:decimal(20,10);" json:"defaultRsiOversold"`
	TelegramBotId                 string          `gorm:"column:telegram_bot_id" json:"telegramBotId"`
	TelegramApiId                 string          `gorm:"column:telegram_api_id" json:"telegramApiId"`
	TelegramChatId                string          `gorm:"column:telegram_chat_id" json:"telegramChatId"`
	BinanceApiKey                 string          `gorm:"column:binance_api_key" json:"binanceApiKey"`
	BinanceSecretKey              string          `gorm:"column:binance_secret_key" json:"binanceSecretKey"`
}

type BotParameters struct {
	Id                       uint            `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UserId                   uint            `gorm:"column:user_id" json:"userId"`
	Broker                   string          `gorm:"column:broker" json:"broker"`
	PercentageTax            decimal.Decimal `gorm:"column:percentage_tax" sql:"type:decimal(20,10);" json:"percentageTax"`
	StopLossPercentage       decimal.Decimal `gorm:"column:stoploss_percentage" sql:"type:decimal(20,10);" json:"stopLossPercentage"`
	TrailingStopLoss         bool            `gorm:"column:trailing_stop_loss"`
	MinimumLimitPercentage   decimal.Decimal `gorm:"column:minimumLimit_percentage" sql:"type:decimal(20,10);" json:"minimumLimitPercentage"`
	RsiPeriod                int             `gorm:"column:rsi_period" json:"rsiPeriod"`
	RsiOverbought            decimal.Decimal `gorm:"column:rsi_overbought" sql:"type:decimal(20,10);" json:"rsiOverbought"`
	RsiOversold              decimal.Decimal `gorm:"column:rsi_oversold" sql:"type:decimal(20,10);" json:"rsiOversold"`
	MaxNumberNegotiations    int             `gorm:"column:max_number_negotiations" json:"maxNumberNegotiations"`
	MinPeriodNextNegotiation int             `gorm:"column:min_period_next_negotiation" json:"minPeriodNextNegotiation"`
	Symbol                   string          `gorm:"column:symbol" json:"symbol"`
	BuyingQty                decimal.Decimal `gorm:"column:buying_qty" sql:"type:decimal(20,10);" json:"buyingQty"`
	BuyingAsset              string          `gorm:"column:buying_asset" json:"buyingAsset"`
	SellingAsset             string          `gorm:"column:selling_asset" json:"sellingAsset"`
	StreamSymbol             string          `gorm:"column:stream_symbol" json:"streamSymbol"`
	StreamInterval           int             `gorm:"column:stream_interval" json:"streamInterval"`
	Created                  *Timestamp      `gorm:"column:created" sql:"type:timestamp;" json:"created,omitempty"`
	Closed                   *Timestamp      `gorm:"column:closed" sql:"type:timestamp;" json:"closed,omitempty"`
}

type CalculatedOperations struct {
	Time           time.Time       `gorm:"column:time" json:"time"`
	Operation      string          `gorm:"column:operation" json:"operation"`
	Price          decimal.Decimal `gorm:"column:price" sql:"type:decimal(20,10);" json:"price"`
	IndicatorValue decimal.Decimal `gorm:"column:indicator_value" sql:"type:decimal(20,10);" json:"indicatorValue"`
	Symbol         string          `gorm:"column:symbol" json:"symbol"`
	SellingAsset   string          `gorm:"column:selling_asset" json:"sellingAsset"`
	Balance        decimal.Decimal `gorm:"column:balance" sql:"type:decimal(20,10);" json:"balance"`
	BuyingAsset    string          `gorm:"column:buying_asset" json:"buyingAsset"`
}

func (User) TableName() string {
	return "users"
}

func (Parameters) TableName() string {
	return "parameters"
}

func (BotParameters) TableName() string {
	return "bot_parameters"
}

func (Candlestick) TableName() string {
	return "candlesticks"
}

func (TradingStatus) TableName() string {
	return "trading_status"
}

func (Operation) TableName() string {
	return "operations"
}

func (CalculatedOperations) TableName() string {
	return "calculated_operations"
}

type Timestamp time.Time

func (m *Timestamp) UnmarshalJSON(data []byte) error {
	if string(data) == "null" || string(data) == `""` {
		return nil
	}
	return json.Unmarshal(data, (*time.Time)(m))
}

func (m Timestamp) MarshalJSON() ([]byte, error) {
	if time.Time(m).IsZero() {
		return []byte("\"\""), nil
	}

	stamp := fmt.Sprintf("\"%s\"", time.Time(m).Format(time.RFC3339))
	return []byte(stamp), nil
}
