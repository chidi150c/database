package model

import "github.com/jinzhu/gorm"

type AppData struct {
	gorm.Model
	DataPoint              int
	Strategy               string
	ShortPeriod            int
	LongPeriod             int
	ShortEMA               float64
	LongEMA                float64
	TargetProfit           float64
	TargetStopLoss         float64
	RiskPositionPercentage float64
	TotalProfitLoss        float64
}

type TradingSystem struct {
	gorm.Model
	Symbol                   string
	ClosingPrices            float64
	Container1               float64
	Container2               float64
	Timestamps               int64
	Signals                  string
	NextInvestBuYPrice       float64
	NextProfitSeLLPrice      float64
	CommissionPercentage     float64
	InitialCapital           float64
	PositionSize             float64
	EntryPrice               float64
	InTrade                  bool
	QuoteBalance             float64
	BaseBalance              float64
	RiskCost                 float64
	DataPoint                int
	CurrentPrice             float64
	EntryQuantity            float64
	Scalping                 string
	StrategyCombLogic        string
	EntryCostLoss            float64
	TradeCount               int
	EnableStoploss           bool
	StopLossTrigered         bool
	StopLossRecover          float64
	RiskFactor               float64
	MaxDataSize              int
	RiskProfitLossPercentage float64
	BaseCurrency             string
	QuoteCurrency            string
	MiniQty                  float64
	MaxQty                   float64
	MinNotional              float64
	StepSize                 float64
}

type DBServicer interface {
	CreateTradingSystem(trade *TradingSystem) (tradeID uint, err error)
	ReadTradingSystem(tradeID uint) (*TradingSystem, error)
	UpdateTradingSystem(trade *TradingSystem) error
	DeleteTradingSystem(tradeID uint) error
	CreateAppData(data *AppData) (dataID uint, err error)
	ReadAppData(dataID uint) (*AppData, error)
	UpdateAppData(data *AppData) error
	DeleteAppData(dataID uint) error
}
