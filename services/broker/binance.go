package broker

import (
	"errors"
	"fmt"
	"net/url"
	"os"
	"strconv"

	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/utils"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"
)

type BinanceApi struct {
	ApiKey    string
	SecretKey string
	BaseUrl   string
}

const RecvWindow = "5000"

func NewBinanceApi() *BinanceApi {
	return &BinanceApi{
		ApiKey:    os.Getenv(consts.BinanceAPIKey),
		SecretKey: os.Getenv(consts.BinanceSecretKey),
		BaseUrl:   os.Getenv(consts.BinanceUrl),
	}
}

func (b *BinanceApi) CreateRequest() *resty.Request {
	client := resty.New()
	return client.R().
		SetHeader("X-MBX-APIKEY", b.ApiKey).
		EnableTrace()
}

func (b *BinanceApi) createSignature(params *url.Values) {
	params.Add("recvWindow", RecvWindow)
	params.Add("timestamp", utils.ConvertToTimestamp(utils.GetCurrentTime()))
	params.Add("signature", utils.CreateHash(b.SecretKey, params.Encode()))
}

func (b *BinanceApi) GetAccountBalance() (*AccountBalance, error) {
	params := url.Values{}
	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&AccountBalance{}).
		Get(b.BaseUrl + "/api/v3/account?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		fmt.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	balance := resp.Result().(*AccountBalance)
	return balance, nil
}

func (b *BinanceApi) GetAssetBalance(symbol string) *Balance {
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

func (b *BinanceApi) GetBalance() (*AccountBalance, error) {
	balance, err := b.GetAccountBalance()
	if err != nil {
		return nil, err
	}
	balances := make([]Balance, 0)

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

func (b *BinanceApi) CreateOrder(symbol string, side string, qty decimal.Decimal, price decimal.Decimal) (*OrderResponse, error) {
	order := new(Order)
	order.Symbol = symbol
	order.Side = side
	order.Qty = qty
	order.Price = price

	return b.doCreateOrder(order)
}

func (b *BinanceApi) doCreateOrder(o *Order) (*OrderResponse, error) {
	params := url.Values{}

	params.Add("symbol", o.Symbol)
	params.Add("side", o.Side)
	params.Add("type", "MARKET")
	params.Add("quantity", o.Qty.String())
	params.Add("newOrderRespType", "FULL")

	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&OrderResponse{}).
		SetBody(params.Encode()).
		Post(b.BaseUrl + "/api/v3/order")

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		fmt.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	oResp := resp.Result().(*OrderResponse)
	return oResp, nil
}

func (b *BinanceApi) Buy(symbol string, qty decimal.Decimal, price decimal.Decimal) (*OrderResponse, error) {
	return b.CreateOrder(symbol, consts.OPERATION_BUY, qty, price)
}

func (b *BinanceApi) Sell(symbol string, qty decimal.Decimal, price decimal.Decimal) (*OrderResponse, error) {
	return b.CreateOrder(symbol, consts.OPERATION_SELL, qty, price)
}

func (b *BinanceApi) QueryOrder(symbol string, orderId int64) (*OrderResponse, error) {
	params := url.Values{}

	params.Add("symbol", symbol)
	params.Add("orderId", utils.ConvertInt64ToString(orderId))

	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&OrderResponse{}).
		Get(b.BaseUrl + "/api/v3/order?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		fmt.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	oResp := resp.Result().(*OrderResponse)
	return oResp, nil
}

func (b *BinanceApi) GetTradeInfo(symbol string, orderId int64) (*TradeResponse, error) {
	params := url.Values{}

	params.Add("symbol", symbol)
	params.Add("orderId", utils.ConvertInt64ToString(orderId))

	b.createSignature(&params)
	resp, err := b.CreateRequest().
		SetResult(&[]TradeResponse{}).
		Get(b.BaseUrl + "/api/v3/myTrades?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		fmt.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	aResp := resp.Result().(*[]TradeResponse)
	if len(*aResp) > 0 {
		return &(*aResp)[0], nil
	}
	return nil, nil
}

func (b *BinanceApi) GetCandlestick(symbol string) (*Kline, error) {
	params := url.Values{}

	params.Add("symbol", symbol)
	params.Add("interval", "1m")
	params.Add("limit", "1")

	resp, err := b.CreateRequest().
		SetResult([][]interface{}{}).
		Get(b.BaseUrl + "/api/v3/klines?" + params.Encode())

	if err != nil {
		return nil, err
	}

	if !b.isSuccess(resp) {
		fmt.Println("Body:", resp)
		return nil, errors.New("Error, StatusCode: " + strconv.Itoa(resp.StatusCode()))
	}

	aResp := resp.Result().(*[][]interface{})
	if len(*aResp) > 0 {
		currTime := utils.GetCurrentTime().Format("2006-01-02T15:04:05")
		r := (*aResp)[0]
		kline := &Kline{
			KlineOpenTime:    currTime,
			OpenPrice:        r[1].(string),
			HighPrice:        r[2].(string),
			LowPrice:         r[3].(string),
			ClosePrice:       r[4].(string),
			Volume:           r[5].(string),
			KlineCloseTime:   currTime,
			QuoteAssetVolume: r[7].(string),
		}
		return kline, nil
	}
	return nil, nil
}

func (b *BinanceApi) isSuccess(r *resty.Response) bool {
	return r != nil && r.StatusCode() >= 200 && r.StatusCode() <= 299
}
