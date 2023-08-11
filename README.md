# jarvis-trading-bot

```shell
curl --location --request POST 'https://{HOST}/trading/signal' \
--header 'Content-Type: application/json' \
--data-raw '{
    "symbol": "BTCUSD",
	"price": "20.32",
	"date": "2023-02-14",
	"operation": "Sell"
}'
```