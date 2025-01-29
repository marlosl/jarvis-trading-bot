package broker

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/shopspring/decimal"

	"jarvis-trading-bot/consts"
)

// Example environment variable constants (you can rename as needed).
// Make sure they're set in your environment or .env file, etc.
// const (
//     OKXApiKey       = "OKX_API_KEY"
//     OKXApiSecret    = "OKX_API_SECRET"
//     OKXApiPassphrase= "OKX_API_PASSPHRASE"
//     OKXBaseUrl      = "OKX_BASE_URL"
// )

// OKXApi holds the credentials and base URL for OKX v5
type OKXApi struct {
	ApiKey     string
	ApiSecret  string
	Passphrase string
	BaseUrl    string
}

// AccountBalance is an example structure that captures part of OKX's
// /api/v5/account/balance response.
type AccountBalance struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		TotalEq string    `json:"totalEq"`
		Details []Balance `json:"details"`
	} `json:"data"`
}

// Balance captures the "details" array items in /api/v5/account/balance.
type Balance struct {
	Ccy      string `json:"ccy"`
	CashBal  string `json:"cashBal"`
	AvailBal string `json:"availBal"`
	// add other fields you need ...
}

// Order is a simple struct for local usage (placing new orders).
type Order struct {
	InstId  string          // e.g. "BTC-USDT"
	Side    string          // "buy" or "sell"
	OrdType string          // e.g. "limit", "market"
	Sz      decimal.Decimal // Size of the order
	Px      decimal.Decimal // Price (for limit orders)
}

// OrderResponse for the OKX trade/order API
type OrderResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		OrdId   string `json:"ordId"`
		ClOrdId string `json:"clOrdId"`
		SMsg    string `json:"sMsg"`
	} `json:"data"`
}

// CancelOrderResponse for OKX trade/cancel-order API
type CancelOrderResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []struct {
		OrdId string `json:"ordId"`
		SMsg  string `json:"sMsg"`
	} `json:"data"`
}

// NewOKXApi initializes the client from environment variables.
// Adjust as needed for your project’s constants or config.
func NewOKXApi() *OKXApi {
	return &OKXApi{
		ApiKey:     os.Getenv(consts.OKXApiKey),
		ApiSecret:  os.Getenv(consts.OKXApiSecret),
		Passphrase: os.Getenv(consts.OKXApiPassphrase),
		BaseUrl:    os.Getenv(consts.OKXBaseUrl),
	}
}

// CreateRequest returns a resty.Request with any default config you want.
// We’ll set content-type JSON on each call, plus the OKX required headers
// in a separate method so we can sign each request with a signature.
func (o *OKXApi) CreateRequest() *resty.Request {
	client := resty.New().
		SetHeader("Content-Type", "application/json").
		SetTimeout(10 * time.Second)

	return client.R()
}

// signRequest signs the request according to OKX docs:
// signature = base64(hmac_sha256(timestamp + method + path + body))
//
// - timestamp must be in ISO8601 format, e.g. "2006-01-02T15:04:05.000Z"
func (o *OKXApi) signRequest(r *resty.Request, method, path string, body string) *resty.Request {
	// 1) Generate timestamp
	timestamp := time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	// 2) The string to sign
	signStr := timestamp + method + path + body

	// 3) HMAC-SHA256 sign, then base64
	h := hmac.New(sha256.New, []byte(o.ApiSecret))
	h.Write([]byte(signStr))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 4) Set required OKX headers
	r.Header.Set("OK-ACCESS-KEY", o.ApiKey)
	r.Header.Set("OK-ACCESS-SIGN", signature)
	r.Header.Set("OK-ACCESS-TIMESTAMP", timestamp)
	r.Header.Set("OK-ACCESS-PASSPHRASE", o.Passphrase)

	return r
}

// isHTTPSuccess is a helper to check whether the HTTP status is 2xx
func (o *OKXApi) isHTTPSuccess(resp *resty.Response) bool {
	return resp != nil && resp.StatusCode() >= 200 && resp.StatusCode() < 300
}

// -----------------------------------------------------------------------------
// BALANCES
// -----------------------------------------------------------------------------

// GetAccountBalance calls GET /api/v5/account/balance to retrieve balances.
// OKX uses "ccy" for currency code. The response is inside a "data" array.
func (o *OKXApi) GetAccountBalance() (*AccountBalance, error) {
	path := "/api/v5/account/balance"
	req := o.CreateRequest()

	// Sign with no request body (GET method). The body is "" for GET.
	req = o.signRequest(req, http.MethodGet, path, "")

	url := o.BaseUrl + path
	resp, err := req.
		SetResult(&AccountBalance{}).
		Get(url)

	if err != nil {
		return nil, err
	}
	if !o.isHTTPSuccess(resp) {
		return nil, fmt.Errorf("HTTP Error: status code %d, body: %s",
			resp.StatusCode(), resp.String())
	}

	balance := resp.Result().(*AccountBalance)
	// Check if code != "0"
	if balance.Code != "0" {
		return nil, fmt.Errorf("OKX error: code=%s, msg=%s", balance.Code, balance.Msg)
	}
	return balance, nil
}

// GetBalance consolidates the above data and returns only nonzero balances,
// as an example of post-filtering.
func (o *OKXApi) GetBalance() (*AccountBalance, error) {
	balance, err := o.GetAccountBalance()
	if err != nil {
		return nil, err
	}
	if len(balance.Data) == 0 {
		return balance, nil // no data, empty
	}

	filteredDetails := make([]Balance, 0)
	for _, detail := range balance.Data[0].Details {
		// Check if AvailBal > 0 (as decimal) to filter.
		avail, errDec := decimal.NewFromString(detail.AvailBal)
		if errDec == nil && avail.GreaterThan(decimal.NewFromInt(0)) {
			filteredDetails = append(filteredDetails, detail)
		}
	}
	balance.Data[0].Details = filteredDetails
	return balance, nil
}

// GetAssetBalance returns the balance info for a single currency, e.g. "USDT".
func (o *OKXApi) GetAssetBalance(ccy string) *Balance {
	acctBal, err := o.GetAccountBalance()
	if err != nil {
		return nil
	}
	if len(acctBal.Data) == 0 {
		return nil
	}
	for _, item := range acctBal.Data[0].Details {
		if strings.EqualFold(item.Ccy, ccy) {
			return &item
		}
	}
	return nil
}

// -----------------------------------------------------------------------------
// ORDERS
// -----------------------------------------------------------------------------

// CreateOrder places a new order using the POST /api/v5/trade/order endpoint.
// Minimal fields: instId, tdMode, side, ordType, sz. Additional optional fields
// can be sent as needed (e.g. px, posSide, clOrdId, etc.).
func (o *OKXApi) CreateOrder(order *Order) (*OrderResponse, error) {
	// For OKX, the request body is JSON, e.g.:
	// {
	//   "instId": "BTC-USDT",
	//   "tdMode": "cash" | "cross" | "isolated" etc,
	//   "side":   "buy" | "sell",
	//   "ordType":"market" | "limit" ...
	//   "sz":     "0.001",
	//   "px":     "20000"  <-- only if limit
	// }
	// This is an example approach. Adjust as needed:
	reqBody := map[string]string{
		"instId":  order.InstId,
		"tdMode":  "cash",        // e.g. "cash", "cross", "isolated"
		"side":    order.Side,    // "buy" or "sell"
		"ordType": order.OrdType, // "market", "limit", etc.
		"sz":      order.Sz.String(),
	}

	// If limit order, add the price:
	if strings.EqualFold(order.OrdType, "limit") && !order.Px.IsZero() {
		reqBody["px"] = order.Px.String()
	}

	path := "/api/v5/trade/order"
	bodyBytes, _ := json.Marshal(reqBody)
	bodyString := string(bodyBytes)

	req := o.CreateRequest()
	req = o.signRequest(req, http.MethodPost, path, bodyString)

	url := o.BaseUrl + path
	resp, err := req.
		SetBody(bodyString).
		SetResult(&OrderResponse{}).
		Post(url)

	if err != nil {
		return nil, err
	}
	if !o.isHTTPSuccess(resp) {
		return nil, fmt.Errorf("HTTP Error: status %d, body: %s",
			resp.StatusCode(), resp.String())
	}

	oResp := resp.Result().(*OrderResponse)
	if oResp.Code != "0" {
		return nil, fmt.Errorf("OKX create order error: code=%s, msg=%s", oResp.Code, oResp.Msg)
	}
	return oResp, nil
}

// CancelOrder cancels an open order. OKX uses POST /api/v5/trade/cancel-order
// Requires: instId + ordId (or clOrdId).
func (o *OKXApi) CancelOrder(instId, ordId string) (*CancelOrderResponse, error) {
	reqBody := map[string]string{
		"instId": instId,
		"ordId":  ordId,
	}
	path := "/api/v5/trade/cancel-order"
	bodyBytes, _ := json.Marshal(reqBody)
	bodyString := string(bodyBytes)

	req := o.CreateRequest()
	req = o.signRequest(req, http.MethodPost, path, bodyString)

	url := o.BaseUrl + path
	resp, err := req.
		SetBody(bodyString).
		SetResult(&CancelOrderResponse{}).
		Post(url)

	if err != nil {
		return nil, err
	}
	if !o.isHTTPSuccess(resp) {
		return nil, fmt.Errorf("HTTP Error: status %d, body: %s",
			resp.StatusCode(), resp.String())
	}

	cResp := resp.Result().(*CancelOrderResponse)
	if cResp.Code != "0" {
		return nil, fmt.Errorf("OKX cancel order error: code=%s, msg=%s", cResp.Code, cResp.Msg)
	}
	return cResp, nil
}

// Convenience buy/sell methods that create a market order.
// You might adapt them for your own logic.
func (o *OKXApi) Buy(instId string, size decimal.Decimal) (*OrderResponse, error) {
	ord := &Order{
		InstId:  instId,
		Side:    "buy",
		OrdType: "market",
		Sz:      size,
	}
	return o.CreateOrder(ord)
}

func (o *OKXApi) Sell(instId string, size decimal.Decimal) (*OrderResponse, error) {
	ord := &Order{
		InstId:  instId,
		Side:    "sell",
		OrdType: "market",
		Sz:      size,
	}
	return o.CreateOrder(ord)
}

// -----------------------------------------------------------------------------
// EXAMPLE: GET Candlestick (Optional Example)
// -----------------------------------------------------------------------------

// Kline is a local struct to hold OHLC data, if you want to fetch it from OKX.
// /api/v5/market/candles?instId=BTC-USDT&bar=1m&limit=1
type Kline struct {
	Ts     string `json:"ts"`     // timestamp
	O      string `json:"open"`   // open
	H      string `json:"high"`   // high
	L      string `json:"low"`    // low
	C      string `json:"close"`  // close
	Vol    string `json:"vol"`    // trading volume
	VolCcy string `json:"volCcy"` // quote volume
}

// GetCandlestick retrieves 1m candle data for the given instrument ID.
func (o *OKXApi) GetCandlestick(instId string) (*Kline, error) {
	path := "/api/v5/market/candles"
	query := fmt.Sprintf("?instId=%s&bar=1m&limit=1", instId)
	finalPath := path + query

	req := o.CreateRequest()
	// sign not required for public endpoints, but let's show how (OKX typically
	// does not require it for /market endpoints). We skip sign in this example:
	// req = o.signRequest(req, http.MethodGet, finalPath, "")

	url := o.BaseUrl + finalPath
	resp, err := req.
		Get(url)

	if err != nil {
		return nil, err
	}
	if !o.isHTTPSuccess(resp) {
		return nil, fmt.Errorf("HTTP Error: status %d, body: %s",
			resp.StatusCode(), resp.String())
	}

	// The response for /candles is typically an array of arrays:
	// {
	//   "code":"0",
	//   "msg":"",
	//   "data":[
	//     ["1634493600000","62617.1","62617.1","62617.1","62617.1","0.008","0.5081233"]
	//   ]
	// }
	type CandleResp struct {
		Code string     `json:"code"`
		Msg  string     `json:"msg"`
		Data [][]string `json:"data"`
	}

	var candleResp CandleResp
	if err := json.Unmarshal(resp.Body(), &candleResp); err != nil {
		return nil, err
	}
	if candleResp.Code != "0" {
		return nil, fmt.Errorf("OKX candle error: code=%s, msg=%s", candleResp.Code, candleResp.Msg)
	}

	if len(candleResp.Data) == 0 {
		return nil, errors.New("no candle data returned")
	}

	// Each entry is [ts, open, high, low, close, volume, volumeInQuote, ...]
	row := candleResp.Data[0]
	k := &Kline{
		Ts:     row[0],
		O:      row[1],
		H:      row[2],
		L:      row[3],
		C:      row[4],
		Vol:    row[5],
		VolCcy: row[6],
	}
	return k, nil
}
