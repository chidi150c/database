package model

import "github.com/jinzhu/gorm"

type AppData struct {
    gorm.Model
    ID             int
    DataPoint      int
    Strategy       string
    ShortPeriod    int
    LongPeriod     int
    ShortEMA       float64
    LongEMA        float64
    TargetProfit   float64
    TargetStopLoss float64
}

type TradingSystem struct {
    gorm.Model
    ID                       int
    Symbol                   string
    ClosingPrices            []float64 `gorm:"type:real[]"` // Use real[] type for SQLite arrays
    Container1               []float64 `gorm:"type:real[]"`
    Container2               []float64 `gorm:"type:real[]"`
    Timestamps               []int64   `gorm:"type:bigint[]"` // Use bigint[] type for SQLite arrays
    Signals                  []string
    NextInvestBuYPrice       []float64 `gorm:"type:real[]"`
    NextProfitSeLLPrice      []float64 `gorm:"type:real[]"`
    CommissionPercentage     float64
    InitialCapital           float64
    PositionSize             float64
    EntryPrice               []float64 `gorm:"type:real[]"`
    InTrade                  bool
    QuoteBalance             float64
    BaseBalance              float64
    RiskCost                 float64
    DataPoint                int
    CurrentPrice             float64
    EntryQuantity            []float64 `gorm:"type:real[]"`
    Scalping                 string
    StrategyCombLogic        string
    EntryCostLoss            []float64 `gorm:"type:real[]"`
    TradeCount               int
    EnableStoploss           bool
    StopLossTrigered         bool
    StopLossRecover          []float64 `gorm:"type:real[]"`
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

// DBServicer defines database services to tradingsystem and appdata model struct
type DBServicer interface {
	Create(trade *TradingSystem, data *AppData)(tradeID int, dataID int, err error)
	Read(tradeID int, dataID int) (trade *TradingSystem, data *AppData, err error)
	Update(trade *TradingSystem, data *AppData) error
	Delete(radeID string, dataID string) error
}
