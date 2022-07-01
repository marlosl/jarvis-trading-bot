package indicators

import (
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/talib"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"

	"github.com/shopspring/decimal"
)

const (
	None       = 0
	Overbought = 1
	Oversold   = 2
)

var zeroDecimal = decimal.NewFromInt(0)

type Psi struct {
	BotParams     structs.BotParams
	TradingStatus *structs.TradingStatus
	Closes        []decimal.Decimal
	OnlyCalculate bool
}

func (p *Psi) CalcRSI(candle *structs.Candlestick) *structs.AnalysisReturn {
	return p.doCalcRSI(candle)
}

func (p *Psi) doCalcRSI(candle *structs.Candlestick) *structs.AnalysisReturn {
	period := p.BotParams.RsiPeriod
	overboughtLimit := p.BotParams.RsiOverbought
	oversoldLimit := p.BotParams.RsiOversold
	minMargPerc := p.BotParams.MinimumLimitPercentage

	aReturn := new(structs.AnalysisReturn)
	aReturn.Operation = consts.OperationNone
	aReturn.Price = zeroDecimal

	if candle.EventTime.Before(p.TradingStatus.LastEvent) {
		return aReturn
	}

	p.TradingStatus.LastEvent = candle.EventTime
	p.Closes = append(p.Closes, candle.ClosePrice)

	if len(p.Closes) < period {
		p.SaveStatus()
		return aReturn
	}

	relForce := p.getRSIRelativeForce(p.Closes)
	relForce = decimal.NewFromInt(1).Add(relForce)
	relForce = decimal.NewFromInt(100).Div(relForce)
	rsi := decimal.NewFromInt(100).Sub(relForce)

	aReturn.PSI = rsi

	log.InfoLogger.Printf("doCalcRSI - Symbol: %s, RSI: %s\n", p.TradingStatus.Symbol, rsi)

	if rsi.Cmp(overboughtLimit) >= 0 && p.TradingStatus.LastStatus != Overbought {
		p.TradingStatus.LastStatus = Overbought

		if p.TradingStatus.BuyAmount.Cmp(zeroDecimal) > 0 {
			sellAmount := candle.ClosePrice.Mul(p.BotParams.BuyingQty)
			p.TradingStatus.SellTax = sellAmount.Mul(p.BotParams.PercentageTax)

			p.TradingStatus.ProfitAmount = sellAmount.Sub(p.TradingStatus.BuyAmount)
			netProfit := sellAmount.Sub(p.TradingStatus.SellTax).Sub(p.TradingStatus.BuyAmount).Sub(p.TradingStatus.BuyTax)

			if p.isTakeProfit(minMargPerc, p.TradingStatus.BuyAmount, netProfit) {
				if p.BotParams.TrailingStopLoss {
					if p.TradingStatus.InitialBuyAmount.IsZero() {
						p.TradingStatus.InitialBuyAmount = p.TradingStatus.BuyAmount
					}
					p.TradingStatus.BuyAmount = sellAmount
					log.InfoLogger.Printf("TRAILING STOP LOSS - sellAmount: %s\n", sellAmount)
					log.InfoLogger.Printf("InitialBuyAmount: %s, profitAmount: %s, sellTax: %s, newBuyAmount: %s\n", p.TradingStatus.InitialBuyAmount, p.TradingStatus.ProfitAmount, p.TradingStatus.SellTax, p.TradingStatus.BuyAmount)

					p.SaveStatus()
				} else {
					p.TradingStatus.TotalProfit = p.TradingStatus.TotalProfit.Add(p.TradingStatus.ProfitAmount)
					p.TradingStatus.TotalTaxes = p.TradingStatus.TotalTaxes.Add(p.TradingStatus.SellTax)

					log.InfoLogger.Printf("SELLING - sellAmount: %s\n", sellAmount)
					log.InfoLogger.Printf("buyAmount: %s, profitAmount: %s, sellTax: %s\n", p.TradingStatus.BuyAmount, p.TradingStatus.ProfitAmount, p.TradingStatus.SellTax)

					p.TradingStatus.BuyAmount = decimal.NewFromInt(0)

					log.InfoLogger.Printf("RSI: %s, Overbought, ClosePrice: %s\n", rsi, candle.ClosePrice)

					p.SaveStatus()
					aReturn.Operation = consts.OperationSell
					aReturn.Price = candle.ClosePrice
					aReturn.PSI = rsi
					return aReturn
				}
			}
		}

		log.InfoLogger.Printf("RSI: %s, Overbought, ClosePrice: %s\n", rsi, candle.ClosePrice)
		log.InfoLogger.Println("")
	} else if rsi.Cmp(oversoldLimit) <= 0 && p.TradingStatus.LastStatus != Oversold {
		p.TradingStatus.LastStatus = Oversold

		if p.TradingStatus.BuyAmount.Equal(zeroDecimal) {
			p.TradingStatus.BuyAmount = candle.ClosePrice.Mul(p.BotParams.BuyingQty)

			p.TradingStatus.BuyTax = p.TradingStatus.BuyAmount.Mul(p.BotParams.PercentageTax)
			p.TradingStatus.TotalTaxes = p.TradingStatus.TotalTaxes.Add(p.TradingStatus.BuyTax)

			log.InfoLogger.Printf("BUYING - buyAmount: %s, buyTax: %s, closePrice: %s, buyingQty:%s\n",
				p.TradingStatus.BuyAmount,
				p.TradingStatus.BuyTax,
				candle.ClosePrice,
				p.BotParams.BuyingQty,
			)
			log.InfoLogger.Printf("RSI: %s, Oversold, ClosePrice: %s\n", rsi, candle.ClosePrice)

			p.SaveStatus()
			aReturn.Operation = consts.OperationBuy
			aReturn.Price = candle.ClosePrice
			aReturn.PSI = rsi
			return aReturn
		}

		log.InfoLogger.Printf("RSI: %s, Oversold, ClosePrice: %s\n", rsi, candle.ClosePrice)
	}

	if p.isStopLoss(p.TradingStatus.BuyAmount, p.TradingStatus.BuyTax, candle.ClosePrice) {
		sellAmount := candle.ClosePrice.Mul(p.BotParams.BuyingQty)
		p.TradingStatus.SellTax = sellAmount.Mul(p.BotParams.PercentageTax)
		p.TradingStatus.ProfitAmount = sellAmount.Sub(p.TradingStatus.BuyAmount)

		p.TradingStatus.TotalProfit = p.TradingStatus.TotalProfit.Add(p.TradingStatus.ProfitAmount)
		p.TradingStatus.TotalTaxes = p.TradingStatus.TotalTaxes.Add(p.TradingStatus.SellTax)

		log.InfoLogger.Printf("SELLING - STOP LOSS - sellAmount: %s\n", sellAmount)
		log.InfoLogger.Printf("buyAmount: %s, profitAmount: %s, sellTax: %s\n", p.TradingStatus.BuyAmount, p.TradingStatus.ProfitAmount, p.TradingStatus.SellTax)

		p.TradingStatus.BuyAmount = decimal.NewFromInt(0)

		p.SaveStatus()
		aReturn.Operation = consts.OperationSell
		aReturn.Price = candle.ClosePrice
		aReturn.PSI = rsi
		return aReturn
	}

	p.Closes = p.Closes[1:]

	p.SaveStatus()
	return aReturn
}

func (p *Psi) SaveStatus() {
	if !p.OnlyCalculate {
		utils.DB.Debug().Save(p.TradingStatus)
	}
}

func (p *Psi) isTakeProfit(minMargPerc decimal.Decimal, buyAmount decimal.Decimal, netProfit decimal.Decimal) bool {
	if minMargPerc.Cmp(zeroDecimal) > 0 {
		if netProfit.Div(buyAmount).Cmp(minMargPerc) >= 0 {
			log.InfoLogger.Printf("buyAmount: %s, netProfit: %s\n", buyAmount, netProfit)
			log.InfoLogger.Printf("TAKEPROFIT: %s\n", netProfit.Div(buyAmount))
			log.InfoLogger.Println("")
		}
		return netProfit.Div(buyAmount).Cmp(minMargPerc) >= 0
	}
	return netProfit.Cmp(zeroDecimal) > 0
}

func (p *Psi) isStopLoss(buyAmount decimal.Decimal, buyTax decimal.Decimal, closePrice decimal.Decimal) bool {
	stopLossPerc := utils.GetDecimalConfig(consts.StopLossPerc)
	if stopLossPerc.Cmp(zeroDecimal) > 0 && buyAmount.Cmp(zeroDecimal) > 0 {
		buyValue := buyAmount.Add(buyTax)
		sellValue := closePrice.Mul(p.BotParams.BuyingQty)
		sellValue = sellValue.Add(sellValue.Mul(p.BotParams.PercentageTax))
		diff := buyValue.Sub(sellValue)
		if diff.Cmp(zeroDecimal) > 0 {
			if diff.Div(buyAmount).Cmp(stopLossPerc) >= 0 {
				log.InfoLogger.Printf("buyValue: %s, sellValue: %s, diff: %s\n", buyValue, sellValue, diff)
				log.InfoLogger.Printf("STOPLOSS: %s\n", diff.Div(buyAmount))
				log.InfoLogger.Println("")
			}
			return diff.Div(buyAmount).Cmp(stopLossPerc) >= 0
		}
	}
	return false
}

func (p *Psi) getRSIRelativeForce(closes []decimal.Decimal) decimal.Decimal {
	relForce := decimal.NewFromFloat(0)
	lows := decimal.NewFromFloat(0)
	highs := decimal.NewFromFloat(0)
	for i, close := range closes {
		if i == 0 {
			continue
		}

		diff := close.Sub(closes[i-1])
		if diff.Cmp(zeroDecimal) >= 0 {
			highs = highs.Add(diff)
		} else {
			lows = lows.Sub(diff)
		}
	}

	if lows.Cmp(zeroDecimal) > 0 {
		relForce = highs.Div(lows)
	}

	return relForce
}

func (p *Psi) GetTaLibRSIRelativeForce(closes []decimal.Decimal, period int) decimal.Decimal {
	inReal := make([]float64, 0)
	for _, close := range closes {
		f, _ := close.Float64()
		inReal = append(inReal, f)
	}
	outReal := talib.Rsi(inReal, period)
	rsi := outReal[len(outReal)-1]
	return decimal.NewFromFloat(rsi)
}
