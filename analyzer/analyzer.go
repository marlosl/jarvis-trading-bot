package analyzer

import (
	"fmt"
	"sort"
	"time"

	"jarvis-trading-bot/analyzer/indicators"
	"jarvis-trading-bot/broker"
	"jarvis-trading-bot/consts"
	"jarvis-trading-bot/notification"
	"jarvis-trading-bot/structs"
	"jarvis-trading-bot/structs/builders"
	"jarvis-trading-bot/utils"
	"jarvis-trading-bot/utils/log"

	"github.com/shopspring/decimal"
)

type Analyzer struct {
	BrokerApi      *broker.BinanceApi
	Simulation     bool
	OnlyCalculate  bool
	ActiveTradings []structs.TradingStatusKey
	BotsParams     map[structs.TradingStatusKey]structs.BotParams
	Indicators     map[structs.TradingStatusKey][]indicators.Indicator
}

func (a *Analyzer) Init() {
	a.Simulation = false
	a.OnlyCalculate = false
	a.initAnalyzer()
}

func (a *Analyzer) InitSimulation() {
	a.Simulation = true
	a.OnlyCalculate = false
	a.initAnalyzer()
}

func (a *Analyzer) InitOnlyCalculate(idBotParam string) {
	a.Simulation = false
	a.OnlyCalculate = true
	a.initOnlyCalculateAnalyzer(idBotParam)
}

func (a *Analyzer) initProperties() {
	a.BrokerApi = new(broker.BinanceApi)
	a.ActiveTradings = make([]structs.TradingStatusKey, 0)
	a.BotsParams = make(map[structs.TradingStatusKey]structs.BotParams)
	a.Indicators = make(map[structs.TradingStatusKey][]indicators.Indicator, 0)
}

func (a *Analyzer) initAnalyzer() {
	a.initProperties()
	a.BrokerApi.Init()

	botParams := GetCryptoBotParameters()
	log.InfoLogger.Printf("botParams: %s\n", utils.SPrintJson(botParams))

	for _, b := range botParams {
		trStatus := GetTradingStatus(b.Id, b.UserId, a.Simulation)
		bp := a.GetBotParams(b)
		a.FillAnalyzerParams(b.Id, b.UserId, bp, trStatus, 1)
	}

	sort.Sort(SortByInstance(a.ActiveTradings))
}

func (a *Analyzer) initOnlyCalculateAnalyzer(idBotParam string) {
	a.initProperties()

	botParams := GetCryptoBotParametersById(idBotParam)
	log.InfoLogger.Printf("botParams: %s\n", utils.SPrintJson(botParams))

	if botParams != nil {
		bp := a.GetBotParams(*botParams)
		a.FillAnalyzerParams(botParams.Id, botParams.UserId, bp, nil, 1)
	}
}

func (a *Analyzer) GetBotParams(b structs.BotParameters) structs.BotParams {
	bb := builders.NewBotParamsBuilder()
	bb.
		SetSymbol(b.Symbol).
		SetBuyingAsset(b.BuyingAsset).
		SetSellingAsset(b.SellingAsset).
		SetBuyingQty(b.BuyingQty).
		SetPercentageTax(b.PercentageTax).
		SetStopLossPercentage(b.StopLossPercentage).
		SetTrailingStopLoss(b.TrailingStopLoss).
		SetMinimumLimitPercentage(b.MinimumLimitPercentage).
		SetRsiPeriod(b.RsiPeriod).
		SetRsiOverbought(b.RsiOverbought).
		SetRsiOversold(b.RsiOversold).
		SetMaxNumberNegotiations(b.MaxNumberNegotiations).
		SetMinPeriodNextNegotiation(b.MinPeriodNextNegotiation).
		SetStreamSymbol(b.StreamSymbol).
		SetStreamInterval(b.StreamInterval)
	bp := bb.Build()
	return *bp
}

func (a *Analyzer) FillAnalyzerParams(botParamId uint, userId uint, b structs.BotParams, trStatus *[]structs.TradingStatus, initialInstance uint) {
	atk := new(structs.TradingStatusKey)
	atk.Symbol = b.Symbol
	atk.BotParameterId = botParamId
	atk.UserId = userId
	atk.Active = true

	if trStatus != nil {
		for _, t := range *trStatus {
			atk.InstanceId = t.InstanceId

			psi := new(indicators.Psi)
			psi.BotParams = b
			psi.TradingStatus = &t
			psi.OnlyCalculate = a.OnlyCalculate
			a.Indicators[*atk] = append(a.Indicators[*atk], psi)

			bband := new(indicators.BBands)
			bband.BotParams = b
			bband.OnlyCalculate = a.OnlyCalculate
			a.Indicators[*atk] = append(a.Indicators[*atk], bband)
		}
	} else {
		atk.InstanceId = initialInstance
		psi := new(indicators.Psi)
		psi.BotParams = b
		psi.TradingStatus = new(structs.TradingStatus)
		psi.TradingStatus.Simulation = a.Simulation
		psi.TradingStatus.UserId = userId
		psi.TradingStatus.BotParameterId = botParamId
		psi.TradingStatus.Symbol = b.Symbol
		psi.TradingStatus.InstanceId = initialInstance
		psi.OnlyCalculate = a.OnlyCalculate
		a.Indicators[*atk] = append(a.Indicators[*atk], psi)

		bband := new(indicators.BBands)
		bband.BotParams = b
		bband.OnlyCalculate = a.OnlyCalculate
		a.Indicators[*atk] = append(a.Indicators[*atk], bband)
	}

	a.BotsParams[*atk] = b
	a.ActiveTradings = append(a.ActiveTradings, *atk)
	log.InfoLogger.Printf("a.ActiveTradings: %s\n", utils.SPrintJson(a.ActiveTradings))
}

func (a *Analyzer) getIndicators(atk structs.TradingStatusKey, indicator string) indicators.Indicator {
	for _, i := range a.Indicators[atk] {
		if i.Name() == indicator {
			return i
		}
	}
	return nil
}

func (a *Analyzer) Process(candle *structs.Candlestick) {
	log.InfoLogger.Printf("a.ActiveTradings: %s\n", utils.SPrintJson(a.ActiveTradings))
	for _, k := range a.ActiveTradings {
		if candle.Symbol == k.Symbol {
			p := a.getIndicators(k, indicators.PSIIndicator).(*indicators.Psi)
			b := a.BotsParams[k]
			aReturn := p.CalcIndicator(candle)

			log.InfoLogger.Printf("Return Analysis: %s\n", utils.SPrintJson(aReturn))

			switch aReturn.Operation {
			case consts.OperationBuy:
				a.BuyOp(&k, &b, aReturn)
			case consts.OperationSell:
				a.SellOp(&k, &b, aReturn)
			}

			if b.MinPeriodNextNegotiation > 0 && int(k.InstanceId) < b.MaxNumberNegotiations {
				interval := time.Duration(b.MinPeriodNextNegotiation) * time.Minute
				log.InfoLogger.Printf("Now: %s\n", time.Now().Format(time.RFC3339))
				log.InfoLogger.Printf("Last Operation: %s\n", p.TradingStatus.LastOperationTime.Format(time.RFC3339))
				log.InfoLogger.Printf("Last Operation + Interval: %s\n", p.TradingStatus.LastOperationTime.Add(interval).Format(time.RFC3339))

				if !p.TradingStatus.LastOperationTime.IsZero() && time.Now().After(p.TradingStatus.LastOperationTime.Add(interval)) && !a.hasNextOperation(k) {
					log.InfoLogger.Println("Adding Next Operation")
					bp := a.BotsParams[k]
					a.FillAnalyzerParams(k.BotParameterId, k.UserId, bp, nil, k.InstanceId+1)
				}
			}
		}
	}
}

func (a *Analyzer) ProcessOnlyCalculate(k structs.TradingStatusKey, candle *structs.Candlestick) *structs.AnalysisReturn {
	aReturn := new(structs.AnalysisReturn)

	for _, indicator := range a.Indicators[k] {
		switch indicator.Name() {
		case indicators.BBandsIndicator:
			b := indicator.(*indicators.BBands)
			v := b.CalcIndicator(candle)
			aReturn.BBandUpper = v.BBandUpper
			aReturn.BBandMiddle = v.BBandMiddle
			aReturn.BBandLower = v.BBandLower
		case indicators.PSIIndicator:
			p := indicator.(*indicators.Psi)
			v := p.CalcIndicator(candle)
			aReturn.PSI = v.PSI
			aReturn.Operation = v.Operation
			aReturn.Price = v.Price
		}
	}
	return aReturn
}

func (a *Analyzer) hasNextOperation(k structs.TradingStatusKey) bool {
	nextKey := k
	nextKey.InstanceId = k.InstanceId + 1
	_, ok := a.Indicators[nextKey]
	return ok
}

func (a *Analyzer) BuyOp(k *structs.TradingStatusKey, b *structs.BotParams, ar *structs.AnalysisReturn) {
	log.InfoLogger.Println("BuyOperation...")

	if a.Simulation {
		log.InfoLogger.Println("It is a simulation...")
		if !a.OnlyCalculate {
			a.SaveSimulatedOperation(k, b, ar)
		}
		return
	}

	success := false
	balanceAmount := decimal.NewFromInt(0)

	balance := a.BrokerApi.GetAssetBalance(b.SellingAsset)
	if balance != nil {
		balanceAmount = balance.Free
	} else {
		log.InfoLogger.Printf("NIL Balance for %s\n", b.SellingAsset)
	}

	total := ar.Price.Mul(b.BuyingQty)

	if a.Simulation {
		log.InfoLogger.Printf("SIMULATION BUYING %s - Asset: %s, Balance: %s, Buying Asset: %s, Buying Amount: %s \n",
			b.Symbol,
			b.SellingAsset,
			balanceAmount,
			b.BuyingAsset,
			total,
		)
		return
	}

	if balance != nil && balance.Free.Cmp(total) >= 0 {
		if r, e := a.BrokerApi.Buy(b.Symbol, b.BuyingQty, ar.Price); e == nil {
			if !a.OnlyCalculate {
				a.SaveOperation(k, b, ar, r)
			}
			utils.PrintJson(r)
			success = true
		}
	}

	if !success {
		logMessage := fmt.Sprintf("NOT BUYING %s - Asset: %s, Balance: %s, Buying Asset: %s, Buying Amount: %s \n",
			b.Symbol,
			b.SellingAsset,
			balance.Free,
			b.BuyingAsset,
			total,
		)
		log.InfoLogger.Println(logMessage)
		notification.SendMessage(logMessage, true)
	}
}

func (a *Analyzer) SellOp(k *structs.TradingStatusKey, b *structs.BotParams, ar *structs.AnalysisReturn) {
	log.InfoLogger.Println("SellOperation...")

	if a.Simulation {
		log.InfoLogger.Println("It is a simulation...")
		if !a.OnlyCalculate {
			a.SaveSimulatedOperation(k, b, ar)
		}
		return
	}

	success := false
	balanceAmount := decimal.NewFromInt(0)

	balance := a.BrokerApi.GetAssetBalance(b.BuyingAsset)
	if balance != nil {
		balanceAmount = balance.Free
	} else {
		log.InfoLogger.Printf("NIL Balance for %s\n", b.BuyingAsset)
	}

	if a.Simulation {
		log.InfoLogger.Printf("SIMULATION SELLING %s - Asset: %s, Balance: %s, Price: %s \n",
			b.Symbol,
			b.BuyingAsset,
			balanceAmount,
			ar.Price,
		)
		return
	}

	if balance != nil && balance.Free.Cmp(decimal.Zero) > 0 {
		qty := balance.Free
		if r, e := a.BrokerApi.Sell(b.Symbol, qty, ar.Price); e == nil {
			if !a.OnlyCalculate {
				a.SaveOperation(k, b, ar, r)
			}
			utils.PrintJson(r)
			success = true
		}
	}

	if !success {
		logMessage := fmt.Sprintf("NOT SELLING %s - Asset: %s, Balance: %s, Price: %s \n",
			b.Symbol,
			b.BuyingAsset,
			balance.Free,
			ar.Price,
		)
		log.InfoLogger.Println(logMessage)
		notification.SendMessage(logMessage, true)
	}
}

func (a *Analyzer) SaveOperation(k *structs.TradingStatusKey, b *structs.BotParams, ar *structs.AnalysisReturn, or *structs.OrderResponse) {
	ob := builders.NewOperationBuilder()
	ob.
		SetBotParameterId(k.BotParameterId).
		SetInstanceId(k.InstanceId).
		SetSymbol(b.Symbol).
		SetOperation(ar.Operation).
		SetOrigQty(or.OrigQty).
		SetBaseAsset(b.BuyingAsset).
		SetBasePrice(or.Price).
		SetOpened(time.Now()).
		SetOrderId(or.OrderId).
		SetExecutedQty(or.ExecutedQty).
		SetCummulativeQuoteQty(or.CummulativeQuoteQty).
		SetType(or.Type).
		SetStatus(or.Status).
		SetTransactTime(utils.ConvertToTime(or.TransactTime)).
		SetFinished(false).
		SetCommissionBase(or.Price.Mul(or.OrigQty).Mul(b.PercentageTax))

	o := ob.Build()

	utils.PrintJson(o)
	notification.SendJson(utils.SPrintJson(o))
	utils.DB.Debug().Save(o)

	p := a.getIndicators(*k, indicators.PSIIndicator).(*indicators.Psi)
	p.TradingStatus.LastOperationTime = time.Now()
	p.SaveStatus()
}

func (a *Analyzer) SaveSimulatedOperation(k *structs.TradingStatusKey, b *structs.BotParams, ar *structs.AnalysisReturn) {
	ob := builders.NewOperationBuilder()
	ob.
		SetBotParameterId(k.BotParameterId).
		SetInstanceId(k.InstanceId).
		SetSymbol(b.Symbol).
		SetOperation(ar.Operation).
		SetOrigQty(b.BuyingQty).
		SetBaseAsset(b.BuyingAsset).
		SetBasePrice(ar.Price).
		SetOpened(time.Now()).
		SetOrderId(0).
		SetExecutedQty(b.BuyingQty).
		SetCummulativeQuoteQty(b.BuyingQty.Mul(ar.Price)).
		SetType("").
		SetStatus("FILLED").
		SetTransactTime(time.Now()).
		SetFinished(true).
		SetCommissionBase(b.BuyingQty.Mul(ar.Price).Mul(b.PercentageTax))

	o := ob.Build()

	utils.PrintJson(o)
	notification.SendJson(utils.SPrintJson(o))
	utils.DB.Save(o)

	p := a.getIndicators(*k, indicators.PSIIndicator).(*indicators.Psi)
	p.TradingStatus.LastOperationTime = time.Now()
	p.SaveStatus()
}

func (a *Analyzer) GetCandlestickAnalysis(k structs.TradingStatusKey, startDate string, endDate string) []structs.CandlestickAnalysis {
	candlesAnalysis := make([]structs.CandlestickAnalysis, 0)

	candles := GetCandlesticks(k.Symbol, startDate, endDate)
	lastCandleTime := time.Date(0001, 1, 1, 00, 00, 00, 00, time.UTC)
	for _, candle := range *candles {
		if candle.EventTime.Before(lastCandleTime) || candle.EventTime.Equal(lastCandleTime) {
			lastCandleTime = candle.EventTime
			continue
		}

		ca := new(structs.CandlestickAnalysis)
		ca.EventTime = candle.EventTime
		ca.Symbol = candle.Symbol
		ca.OpenPrice = candle.OpenPrice
		ca.ClosePrice = candle.ClosePrice
		ca.HighPrice = candle.HighPrice
		ca.LowPrice = candle.LowPrice
		ca.Volume = candle.BaseAssetVolume

		aReturn := a.ProcessOnlyCalculate(k, &candle)
		ca.Operation = aReturn.Operation
		ca.PSI = aReturn.PSI

		ca.BBandUpper = aReturn.BBandUpper
		ca.BBandMiddle = aReturn.BBandMiddle
		ca.BBandLower = aReturn.BBandLower

		lastCandleTime = candle.EventTime
		candlesAnalysis = append(candlesAnalysis, *ca)
	}

	return candlesAnalysis
}

func (a *Analyzer) CreateTradingStatusKeyByBotParamId(botParamId string) structs.TradingStatusKey {
	b := GetCryptoBotParametersById(botParamId)
	atk := new(structs.TradingStatusKey)
	atk.Symbol = b.Symbol
	atk.BotParameterId = b.Id
	atk.UserId = b.UserId
	atk.Active = true
	atk.InstanceId = 1
	return *atk
}

func QueryOperation(symbol string, orderId string) {
	a := new(Analyzer)
	a.Init()

	oi := utils.ConvertStringToInt64(orderId)
	r, e := a.BrokerApi.GetTradeInfo(symbol, oi)
	if e != nil {
		log.InfoLogger.Printf("Error Executing Operation: %s\n", e)
	} else {
		utils.PrintJson(r)
	}
}
