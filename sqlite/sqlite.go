package sqlite

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/chidi150c/database/model"
    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

type DBServices struct{
	DB *sql.DB
	Result sql.Result
}

func NewDBServices(dbName string) *DBServices{
    // Open or create the SQLite database
    db, err := sql.Open("sqlite3", dbName)
    if err != nil {
        log.Fatal(err)
    }
	a := &DBServices{
		DB: db,
	}
	return a
}
var _ model.DBServicer = &DBServices{}

func (a *DBServices)Create(trade *model.TradingSystem, data *model.AppData)error{
	 // Create the table if it doesn't exist
	 createTableSQL := `
	 CREATE TABLE IF NOT EXISTS app_data (
		 id INTEGER PRIMARY KEY AUTOINCREMENT,
		 data_point INTEGER,
		 strategy TEXT,
		 short_period INTEGER,
		 long_period INTEGER,
		 short_ema REAL,
		 long_ema REAL,
		 target_profit REAL,
		 target_stop_loss REAL
	 );`
 
	 _, err := a.DB.Exec(createTableSQL)
	 if err != nil {
		 log.Fatal(err)
	 }

	 return err
}

func (a *DBServices)Update(trade*model.TradingSystem, data *model.AppData)(err error){
	insertSQL := `
	INSERT INTO app_data (data_point, strategy, short_period, long_period, short_ema, long_ema, target_profit, target_stop_loss)
	VALUES (?, ?, ?, ?, ?, ?, ?, ?);`
 
	a.Result, err = a.DB.Exec(insertSQL, data.DataPoint, data.Strategy, data.ShortPeriod, data.LongPeriod,
	 data.ShortEMA, data.LongEMA, data.TargetProfit, data.TargetStopLoss)
	if err != nil {
	 log.Fatal(err)
	}
	return err
}

func (a *DBServices)Read(tradeId string, dataID string)(trade *model.TradingSystem, retrievedData *model.AppData, err error){
	 // Retrieve data from the database
	 resultIns , err := a.Result.LastInsertId()
	 if err != nil{
		 log.Fatal(err)
	 }
	 row := a.DB.QueryRow("SELECT * FROM app_data WHERE id = ?", resultIns)
	 err = row.Scan(&retrievedData.ID, &retrievedData.DataPoint, &retrievedData.Strategy,
		 &retrievedData.ShortPeriod, &retrievedData.LongPeriod, &retrievedData.ShortEMA,
		 &retrievedData.LongEMA, &retrievedData.TargetProfit, &retrievedData.TargetStopLoss)
	 if err != nil {
		 log.Fatal(err)
	 }
 
	 // Print the retrieved data
	 fmt.Printf("Retrieved Data: %+v\n", retrievedData)
	 return trade, retrievedData, nil
}

func (a *DBServices)Delete(tradeID string, dataID string)(err error){
	return err
}