package broker

import (
	"errors"
	"net/url"
	"strconv"
	"time"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

type BinanceApi struct {
	ApiKey    string
	SecretKey string
	BaseUrl   string
}

const RecvWindow = "5000"

func (b *BinanceApi) Init() {
	b.ApiKey = utils.GetStringConfig(consts.BinanceAPIKey)
	b.SecretKey = utils.GetStringConfig(consts.BinanceSecretKey)
	b.BaseUrl = utils.GetStringConfig(consts.BinanceUrl)
}

func (b *BinanceApi) CreateRequest() *resty.Request {
	client := resty.New()
	return client.R().
		SetHeader("X-MBX-APIKEY", b.ApiKey).
		EnableTrace()
}

func (b *BinanceApi) createSignature(params *url.Values) {
	params.Add("recvWindow", RecvWindow)
	params.Add("timestamp", utils.ConvertToTimestamp(time.Now()))
	params.Add("signature", utils.CreateHash(b.SecretKey, params.Encode()))
}

func (b *BinanceApi) GetAccountBalance() (*structs.AccountBalance, error) {
	params := url.Values{}
	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&structs.AccountBalance{}).
		Get(b.BaseUrl + "/api/v3/account?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		log.InfoLogger.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	balance := resp.Result().(*structs.AccountBalance)
	return balance, nil
}

func (b *BinanceApi) GetAssetBalance(symbol string) *structs.Balance {
	balance, err := b.GetAccountBalance()

	if err == nil {
		for _, item := range balance.Balances {
			if item.Asset == symbol {
				return &item
			}
		}
	}
	return nil
}

func (b *BinanceApi) GetBalance() (*structs.AccountBalance, error) {
	balance, err := b.GetAccountBalance()
	if err != nil {
		return nil, err
	}
	balances := make([]structs.Balance, 0)

	for _, item := range balance.Balances {
		free := item.Free
		locked := item.Locked

		if free.Cmp(decimal.NewFromInt(0)) > 0 || locked.Cmp(decimal.NewFromInt(0)) > 0 {
			balances = append(balances, item)
		}
	}
	balance.Balances = balances
	return balance, nil
}

func (b *BinanceApi) CreateOrder(symbol string, side string, qty decimal.Decimal, price decimal.Decimal) (*structs.OrderResponse, error) {
	order := new(structs.Order)
	order.Symbol = symbol
	order.Side = side
	order.Qty = qty
	order.Price = price

	return b.doCreateOrder(order)
}

func (b *BinanceApi) doCreateOrder(o *structs.Order) (*structs.OrderResponse, error) {
	params := url.Values{}

	params.Add("symbol", o.Symbol)
	params.Add("side", o.Side)
	params.Add("type", "MARKET")
	params.Add("quantity", o.Qty.String())
	params.Add("newOrderRespType", "FULL")

	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&structs.OrderResponse{}).
		SetBody(params.Encode()).
		Post(b.BaseUrl + "/api/v3/order")

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		log.InfoLogger.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	oResp := resp.Result().(*structs.OrderResponse)
	return oResp, nil
}

func (b *BinanceApi) Buy(symbol string, qty decimal.Decimal, price decimal.Decimal) (*structs.OrderResponse, error) {
	return b.CreateOrder(symbol, "BUY", qty, price)
}

func (b *BinanceApi) Sell(symbol string, qty decimal.Decimal, price decimal.Decimal) (*structs.OrderResponse, error) {
	return b.CreateOrder(symbol, "SELL", qty, price)
}

func (b *BinanceApi) QueryOrder(symbol string, orderId int64) (*structs.OrderResponse, error) {
	params := url.Values{}

	params.Add("symbol", symbol)
	params.Add("orderId", utils.ConvertInt64ToString(orderId))

	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&structs.OrderResponse{}).
		Get(b.BaseUrl + "/api/v3/order?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		log.InfoLogger.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	oResp := resp.Result().(*structs.OrderResponse)
	return oResp, nil
}

func (b *BinanceApi) GetTradeInfo(symbol string, orderId int64) (*structs.TradeResponse, error) {
	params := url.Values{}

	params.Add("symbol", symbol)
	params.Add("orderId", utils.ConvertInt64ToString(orderId))

	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&[]structs.TradeResponse{}).
		Get(b.BaseUrl + "/api/v3/myTrades?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		log.InfoLogger.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	aResp := resp.Result().(*[]structs.TradeResponse)
	if len(*aResp) > 0 {
		return &(*aResp)[0], nil
	}
	return nil, nil
}

func (b *BinanceApi) isSuccess(r *resty.Response) bool {
	return r != nil && r.StatusCode() >= 200 && r.StatusCode() <= 299
}

func GetBrokerBalance(symbol string) {
	b := new(BinanceApi)
	b.Init()

	var balance interface{}
	var err error

	if len(symbol) > 1 {
		balance = b.GetAssetBalance(symbol)
	} else {
		balance, err = b.GetBalance()
	}

	if err != nil {
		log.ErrorLogger.Println(err)
	} else {
		utils.PrintJson(balance)
	}
}
