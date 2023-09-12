package model

import (
    "database/sql/driver"
    "encoding/json"
    "errors"
	"github.com/jinzhu/gorm"
)


// Float64Slice is a custom data type for a float64 slice.
type Float64Slice []float64

// Scan scans a value into Float64Slice.
func (f *Float64Slice) Scan(value interface{}) error {
    if value == nil {
        *f = nil
        return nil
    }

    byteValue, ok := value.([]byte)
    if !ok {
        return errors.New("Invalid Scan Source")
    }

    return json.Unmarshal(byteValue, f)
}

// Value converts Float64Slice to a database value.
func (f Float64Slice) Value() (driver.Value, error) {
    if f == nil {
        return nil, nil
    }

    return json.Marshal(f)
}



type AppData struct {
    gorm.Model
    DataPoint        int `json:"data_point"`
    Strategy         string  `json:"strategy"`
    ShortPeriod      int     `json:"short_period"`
    LongPeriod       int     `json:"long_period"`
    ShortEMA         float64 `json:"short_ema"`
    LongEMA          float64 `json:"long_ema"`
    TargetProfit     float64 `json:"target_profit"`
    TargetStopLoss   float64 `json:"target_stop_loss"`
    RiskPositionPercentage float64 `json:"risk_position_percentage"`
    TotalProfitLoss  float64 `json:"total_profit_loss"`
}

type TradingSystem struct {
	gorm.Model
	Symbol                   string
	ClosingPrices            float64
	Timestamps               int64
	Signals                  string
	NextInvestBuYPrice       Float64Slice `gorm:"type:json"`
	NextProfitSeLLPrice      Float64Slice `gorm:"type:json"`
	CommissionPercentage     float64
	InitialCapital           float64
	PositionSize             float64
	EntryPrice               Float64Slice `gorm:"type:json"`
	InTrade                  bool
	QuoteBalance             float64
	BaseBalance              float64
	RiskCost                 float64
	DataPoint                int
	CurrentPrice             float64
	EntryQuantity            Float64Slice `gorm:"type:json"`
	EntryCostLoss            Float64Slice `gorm:"type:json"`
	TradeCount               int
	TradingLevel             int       
	ClosedWinTrades          int      
	EnableStoploss           bool
	StopLossTrigered         bool
	StopLossRecover          Float64Slice `gorm:"type:json"`
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
