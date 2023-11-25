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

type StringSlice []string

func (s StringSlice) Value() (driver.Value, error) {
    // Serialize the StringSlice to a JSON string
    return json.Marshal(s)
}

func (s *StringSlice) Scan(value interface{}) error {
    // Deserialize the JSON string to a StringSlice
    if value == nil {
        return nil
    }
    if str, ok := value.([]byte); ok {
        return json.Unmarshal(str, s)
    }
    return errors.New("Invalid value type for StringSlice")
}

type Int64Slice []int64

func (i Int64Slice) Value() (driver.Value, error) {
    // Serialize the Int64Slice to a JSON string
    return json.Marshal(i)
}
// 
func (i *Int64Slice) Scan(value interface{}) error {
    // Deserialize the JSON string to an Int64Slice
    if value == nil {
        return nil
    }
    if str, ok := value.([]byte); ok {
        return json.Unmarshal(str, i)
    }
    return errors.New("Invalid value type for Int64Slice")
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
	ClosingPrices            Float64Slice `gorm:"type:json"`
	Timestamps               Int64Slice `gorm:"type:json"`
	Signals                  StringSlice `gorm:"type:json"`
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

type TradingSystemData struct {
	ID                       uint
	Symbol                   string    `json:"symbol"`
	ClosingPrices            []float64   `json:"closing_prices"`
	Timestamps               []int64     `json:"timestamps"`
	Signals                  []string    `json:"signals"`
	NextInvestBuYPrice       []float64 `json:"next_invest_buy_price"`
	NextProfitSeLLPrice      []float64 `json:"next_profit_sell_price"`
	CommissionPercentage     float64   `json:"commission_percentage"`
	InitialCapital           float64   `json:"initial_capital"`
	PositionSize             float64   `json:"position_size"`
	EntryPrice               []float64 `json:"entry_price"`
	InTrade                  bool      `json:"in_trade"`
	QuoteBalance             float64   `json:"quote_balance"`
	BaseBalance              float64   `json:"base_balance"`
	RiskCost                 float64   `json:"risk_cost"`
	DataPoint                int       `json:"data_point"`
	CurrentPrice             float64   `json:"current_price"`
	EntryQuantity            []float64 `json:"entry_quantity"`
	EntryCostLoss            []float64 `json:"entry_cost_loss"`
	TradeCount               int       `json:"trade_count"`
	TradingLevel             int       `json:"trading_level"`
	ClosedWinTrades          int       `json:"closed_win_trades"`
	EnableStoploss           bool      `json:"enable_stoploss"`
	StopLossTrigered         bool      `json:"stop_loss_triggered"`
	StopLossRecover          []float64 `json:"stop_loss_recover"`
	RiskFactor               float64   `json:"risk_factor"`
	MaxDataSize              int       `json:"max_data_size"`
	RiskProfitLossPercentage float64   `json:"risk_profit_loss_percentage"`
	BaseCurrency             string    `json:"base_currency"`
	QuoteCurrency            string    `json:"quote_currency"`
	MiniQty                  float64   `json:"mini_qty"`
	MaxQty                   float64   `json:"max_qty"`
	MinNotional              float64   `json:"min_notional"`
	StepSize                 float64   `json:"step_size"`
}
