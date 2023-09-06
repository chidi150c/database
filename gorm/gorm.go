package gorm

import (
	"fmt"
	"log"

    "github.com/chidi150c/database/model"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

// DBServices is an implementation of the DBServicer interface
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
	a := &DBServices{
		DB: db,
	}

	return a, nil
}

func (a *DBServices) CheckAndCreateTables() error {
	// Check if the TradingSystem table exists
	tradingSystemTableExists := tableExists(a.DB, "trading_systems")

	// Check if the AppData table exists
	appDataTableExists := tableExists(a.DB, "app_data")

	// Start a new transaction
	tx := a.DB.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Create tables conditionally
	if !tradingSystemTableExists {
		if err := tx.AutoMigrate(&model.TradingSystem{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Error migrating TradingSystem table: %v", err)
		}
	}

	if !appDataTableExists {
		if err := tx.AutoMigrate(&model.AppData{}).Error; err != nil {
			tx.Rollback()
			return fmt.Errorf("Error migrating AppData table: %v", err)
		}
	}

	// Commit the transaction if everything is successful
	if err := tx.Commit().Error; err != nil {
		return fmt.Errorf("Error committing transaction: %v", err)
	}

	return nil
}

func tableExists(db *gorm.DB, tableName string) bool {
	// Check if the table exists in the database
	return db.HasTable(tableName)
}

func (s *DBServices) CreateTradingSystem(trade *model.TradingSystem) (int, error) {
    if err := s.DB.Create(trade).Error; err != nil {
        return 0, err
    }
    return trade.ID, nil
}

func (s *DBServices) ReadTradingSystem(tradeID int) (*model.TradingSystem, error) {
    var trade model.TradingSystem
    if err := s.DB.First(&trade, tradeID).Error; err != nil {
        return nil, err
    }
    return &trade, nil
}

func (s *DBServices) UpdateTradingSystem(trade *model.TradingSystem) error {
    if err := s.DB.Save(trade).Error; err != nil {
        return err
    }
    return nil
}

func (s *DBServices) DeleteTradingSystem(tradeID int) error {
    if err := s.DB.Delete(&model.TradingSystem{}, tradeID).Error; err != nil {
        return err
    }
    return nil
}

func (s *DBServices) CreateAppData(data *model.AppData) (int, error) {
    if err := s.DB.Create(data).Error; err != nil {
        return 0, err
    }
    return data.ID, nil
}

func (s *DBServices) ReadAppData(dataID int) (*model.AppData, error) {
    var data model.AppData
    if err := s.DB.First(&data, dataID).Error; err != nil {
        return nil, err
    }
    return &data, nil
}

func (s *DBServices) UpdateAppData(data *model.AppData) error {
    if err := s.DB.Save(data).Error; err != nil {
        return err
    }
    return nil
}

func (s *DBServices) DeleteAppData(dataID int) error {
    if err := s.DB.Delete(&model.AppData{}, dataID).Error; err != nil {
        return err
    }
    return nil
}
