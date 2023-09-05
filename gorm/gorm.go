package gorm

import (
	"fmt"
	"log"

	"github.com/chidi150c/database/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBServices struct {
	DB *gorm.DB
}

//NewDBServices has an initializeDatabase function that checks if the required tables (TradingSystem and AppData) exist in the database.
// It creates these tables if they don't exist. This function is called during the creation of a new DBServices instance.
func NewDBServices(dbName string) (*DBServices, error) {
	db, err := gorm.Open("sqlite3", dbName)
	if err != nil {
		return &DBServices{}, fmt.Errorf("NewDBServices error: %v", err)
	}
	// Initialize the database, create tables, and perform migrations
	initializeDatabase(db)
	a := &DBServices{
		DB: db,
	}

	return a, nil
}

var _ model.DBServicer = &DBServices{}

func initializeDatabase(db *gorm.DB){
	// Check if the TradingSystem table exists
	tradingSystemTableExists := tableExists(db, "trading_systems")

	// Check if the AppData table exists
	appDataTableExists := tableExists(db, "app_data")

	// Start a new transaction
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create tables and perform migrations conditionally
	if !tradingSystemTableExists {
		if err := tx.AutoMigrate(&model.TradingSystem{}).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Error migrating TradingSystem table: %v", err)
		}
	}

	if !appDataTableExists {
		if err := tx.AutoMigrate(&model.AppData{}).Error; err != nil {
			tx.Rollback()
			log.Fatalf("Error migrating AppData table: %v", err)
		}
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		log.Fatalf("Error committing transaction: %v", err)
	}
}



func tableExists(db *gorm.DB, tableName string) bool {
	// Check if the table exists in the database
	return db.HasTable(tableName)
}

func (a *DBServices) Create(trade *model.TradingSystem, data *model.AppData) (tradeID int, dataID int, err error) {
	// Start a new transaction
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Create TradingSystem and AppData records within the transaction
	if err := tx.Create(trade).Error; err != nil {
		tx.Rollback()
		return 0,0,err
	}
	if err := tx.Create(data).Error; err != nil {
		tx.Rollback()
		return 0,0,err
	}
	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return 0,0,err
	}
	return trade.ID, data.ID, nil
}


func (a *DBServices) Update(trade *model.TradingSystem, data *model.AppData) (err error) {
	// Start a new transaction
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	// Update multiple fields of the AppData model
	if err := tx.Model(data).Where("id = ?", data.ID).Updates(map[string]interface{}{
		"DataPoint":      data.DataPoint,
		"Strategy":       data.Strategy,
		"ShortPeriod":    data.ShortPeriod,
		"LongPeriod":     data.LongPeriod,
		"ShortEMA":       data.ShortEMA,
		"LongEMA":        data.LongEMA,
		"TargetProfit":   data.TargetProfit,
		"TargetStopLoss": data.TargetStopLoss,
	}).Error; err != nil {
		return err
	}
	// Update multiple fields of the TradingSystem (ts) model
	if err := tx.Model(trade).Where("id = ?", trade.ID).Updates(map[string]interface{}{
		"Symbol":                   trade.Symbol,
		"ClosingPrices":            trade.ClosingPrices,
		"Container1":               trade.Container1,
		"Container2":               trade.Container2,
		"Timestamps":               trade.Timestamps,
		"Signals":                  trade.Signals,
		"NextInvestBuYPrice":       trade.NextInvestBuYPrice,
		"NextProfitSeLLPrice":      trade.NextProfitSeLLPrice,
		"CommissionPercentage":     trade.CommissionPercentage,
		"InitialCapital":           trade.InitialCapital,
		"PositionSize":             trade.PositionSize,
		"EntryPrice":               trade.EntryPrice,
		"InTrade":                  trade.InTrade,
		"QuoteBalance":             trade.QuoteBalance,
		"BaseBalance":              trade.BaseBalance,
		"RiskCost":                 trade.RiskCost,
		"DataPoint":                trade.DataPoint,
		"CurrentPrice":             trade.CurrentPrice,
		"EntryQuantity":            trade.EntryQuantity,
		"Scalping":                 trade.Scalping,
		"StrategyCombLogic":        trade.StrategyCombLogic,
		"EntryCostLoss":            trade.EntryCostLoss,
		"TradeCount":               trade.TradeCount,
		"EnableStoploss":           trade.EnableStoploss,
		"StopLossTrigered":         trade.StopLossTrigered,
		"StopLossRecover":          trade.StopLossRecover,
		"RiskFactor":               trade.RiskFactor,
		"MaxDataSize":              trade.MaxDataSize,
		"RiskProfitLossPercentage": trade.RiskProfitLossPercentage,
		"BaseCurrency":             trade.BaseCurrency,
		"QuoteCurrency":            trade.QuoteCurrency,
		"MiniQty":                  trade.MiniQty,
		"MaxQty":                   trade.MaxQty,
		"MinNotional":              trade.MinNotional,
		"StepSize":                 trade.StepSize,
	}).Error; err != nil {
		return err
	}
	return err
}

func (a *DBServices) Read(tradeID int, dataID int) (trade *model.TradingSystem, data *model.AppData, err error) {
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
    trade = &model.TradingSystem{}
    data = &model.AppData{}
    if err := tx.First(trade, "id = ?", tradeID).Error; err != nil {
        return nil, nil, err
    }
    if err := tx.First(data, "id = ?", dataID).Error; err != nil {
        return nil, nil, err
    }
    return trade, data, nil
}

func (a *DBServices) Delete(tradeID string, dataID string) (err error) {
    // Delete a TradingSystem record
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
    if err := tx.Where("id = ?", tradeID).Delete(&model.TradingSystem{}).Error; err != nil {
        return err
    }
    // Delete an AppData record
    if err := tx.Where("id = ?", dataID).Delete(&model.AppData{}).Error; err != nil {
        return err
    }
    return nil
}

