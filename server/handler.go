package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/chidi150c/database/gorm"
	"github.com/chidi150c/database/model"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

// Define a struct that matches the expected WebSocket message format
type WebSocketMessage struct {
	Action string                 `json:"action"`
	Entity string                 `json:"entity"`
	Data   map[string]interface{} `json:"data"`
}

// WebService is a user login-aware wrapper for a html/template.
var connections = struct {
	sync.RWMutex
	m map[*websocket.Conn]struct{}
}{
	m: make(map[*websocket.Conn]struct{}),
}

type TradeHandler struct {
	mux    *chi.Mux
	DbName string
}

func NewTradeHandler(dbName string) TradeHandler {
	h := TradeHandler{
		mux:    chi.NewRouter(),
		DbName: dbName,
	}
	h.mux.Get("/database-services/ws", h.DataBaseSocketHandler)
	return h
}

func (h TradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}
func (th *TradeHandler) DataBaseSocketHandler(w http.ResponseWriter, r *http.Request) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			// 		if r.Header.Get("Origin") != h.HostSite {
			// 			return false
			// 		}
			return true
		},
	}
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	//Register the conn in connections
	connections.Lock()
	connections.m[conn] = struct{}{}
	connections.Unlock()

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			log.Print("ConnRDErr")
			break
		}

		var msg WebSocketMessage
		if err := json.Unmarshal(p, &msg); err != nil {
			log.Println(err)
			continue
		}

		go processMessage(conn, msg, th.DbName)
	}
	connections.Lock()
	delete(connections.m, conn)
	connections.Unlock()
}
func processMessage(conn *websocket.Conn, message WebSocketMessage, dbName string) {
	msg := ""
	DBServices, err := gorm.NewDBServices(dbName)
	if err != nil {
		log.Printf("Error: while processing Websocket massge: %v", err)
		return
	}
	switch message.Action {
	case "create":
		// Handle create action for entities
		if message.Entity == "trading-system" {
			// Parse and process trading system creation
			var ts model.TradingSystemData
			dataByte, _ := json.Marshal(message.Data)
			// Deserialize the WebSocket message directly into the struct
			if err := json.Unmarshal(dataByte, &ts); err != nil {
				msg = fmt.Sprintf("Error3 parsing WebSocket message: %v", err)
				writeResponseWithID(msg, ts.ID, conn)
				return
			}
			// Convert standard types to custom data types
			dbTrade := &model.TradingSystem{
				Symbol:                   ts.Symbol,
				ClosingPrices:            model.Float64Slice(ts.ClosingPrices),
				Timestamps:               model.Int64Slice(ts.Timestamps),
				Signals:                  model.StringSlice(ts.Signals),
				NextInvestBuYPrice:       model.Float64Slice(ts.NextInvestBuYPrice),
				NextProfitSeLLPrice:      model.Float64Slice(ts.NextProfitSeLLPrice),
				CommissionPercentage:     ts.CommissionPercentage,
				InitialCapital:           ts.InitialCapital,
				PositionSize:             ts.PositionSize,
				EntryPrice:               model.Float64Slice(ts.EntryPrice),
				InTrade:                  ts.InTrade,
				QuoteBalance:             ts.QuoteBalance,
				BaseBalance:              ts.BaseBalance,
				RiskCost:                 ts.RiskCost,
				DataPoint:                ts.DataPoint,
				CurrentPrice:             ts.CurrentPrice,
				EntryQuantity:            model.Float64Slice(ts.EntryQuantity),
				EntryCostLoss:            model.Float64Slice(ts.EntryCostLoss),
				TradeCount:               ts.TradeCount,
				TradingLevel:             ts.TradingLevel,
				ClosedWinTrades:          ts.ClosedWinTrades,
				EnableStoploss:           ts.EnableStoploss,
				StopLossTrigered:         ts.StopLossTrigered,
				StopLossRecover:          model.Float64Slice(ts.StopLossRecover),
				RiskFactor:               ts.RiskFactor,
				MaxDataSize:              ts.MaxDataSize,
				RiskProfitLossPercentage: ts.RiskProfitLossPercentage,
				BaseCurrency:             ts.BaseCurrency,
				QuoteCurrency:            ts.QuoteCurrency,
				MiniQty:                  ts.MiniQty,
				MaxQty:                   ts.MaxQty,
				MinNotional:              ts.MinNotional,
				StepSize:                 ts.StepSize,
				TargetStopLoss:           ts.TargetStopLoss,
				TargetProfit:             ts.TargetProfit,
				TotalProfitLoss:          ts.TotalProfitLoss,
				RiskPositionPercentage:   ts.RiskPositionPercentage,
				ShortPeriod:              ts.ShortPeriod,
				LongPeriod:               ts.LongPeriod,
			}
			// Insert the new trading system into the database
			tradeID, err := DBServices.CreateTradingSystem(dbTrade)
			if err != nil {
				msg = fmt.Sprintf("Error creating trading system: %v", err)
				writeResponseWithID(msg, tradeID, conn)
				return
			}
			writeResponseWithID("TradingSystem Created successfully", tradeID, conn)
		}
	case "read":
		if message.Entity == "trading-system" {
			var ts model.TradingSystem
			dataByte, _ := json.Marshal(message.Data)
			// Deserialize the WebSocket message directly into the struct
			if err := json.Unmarshal(dataByte, &ts); err != nil {
				msg = fmt.Sprintf("Error5 parsing WebSocket message: %v", err)
				writeResponseWithData(msg, ts, conn)
				return
			}
			tradeID := ts.ID
			// Fetch the trading system from the database based on tradeID
			dbTrade, err := DBServices.ReadTradingSystem(tradeID)
			if err != nil {
				msg = fmt.Sprintf("Error retrieving trading system: %v", err)
				writeResponseWithData(msg, &model.TradingSystemData{}, conn)
				return
			}
			// Convert custom data types to standard types
			dataTrade := &model.TradingSystemData{
				ID:                       dbTrade.ID,
				Symbol:                   dbTrade.Symbol,
				ClosingPrices:            []float64(dbTrade.ClosingPrices),
				Timestamps:               []int64(dbTrade.Timestamps),
				Signals:                  []string(dbTrade.Signals),
				NextInvestBuYPrice:       []float64(dbTrade.NextInvestBuYPrice),
				NextProfitSeLLPrice:      []float64(dbTrade.NextProfitSeLLPrice),
				CommissionPercentage:     dbTrade.CommissionPercentage,
				InitialCapital:           dbTrade.InitialCapital,
				PositionSize:             dbTrade.PositionSize,
				EntryPrice:               []float64(dbTrade.EntryPrice),
				InTrade:                  dbTrade.InTrade,
				QuoteBalance:             dbTrade.QuoteBalance,
				BaseBalance:              dbTrade.BaseBalance,
				RiskCost:                 dbTrade.RiskCost,
				DataPoint:                dbTrade.DataPoint,
				CurrentPrice:             dbTrade.CurrentPrice,
				EntryQuantity:            []float64(dbTrade.EntryQuantity),
				EntryCostLoss:            []float64(dbTrade.EntryCostLoss),
				TradeCount:               dbTrade.TradeCount,
				TradingLevel:             dbTrade.TradingLevel,
				ClosedWinTrades:          dbTrade.ClosedWinTrades,
				EnableStoploss:           dbTrade.EnableStoploss,
				StopLossTrigered:         dbTrade.StopLossTrigered,
				StopLossRecover:          []float64(dbTrade.StopLossRecover),
				RiskFactor:               dbTrade.RiskFactor,
				MaxDataSize:              dbTrade.MaxDataSize,
				RiskProfitLossPercentage: dbTrade.RiskProfitLossPercentage,
				BaseCurrency:             dbTrade.BaseCurrency,
				QuoteCurrency:            dbTrade.QuoteCurrency,
				MiniQty:                  dbTrade.MiniQty,
				MaxQty:                   dbTrade.MaxQty,
				MinNotional:              dbTrade.MinNotional,
				StepSize:                 dbTrade.StepSize,
				TargetStopLoss:           dbTrade.TargetStopLoss,
				TargetProfit:             dbTrade.TargetProfit,
				TotalProfitLoss:          dbTrade.TotalProfitLoss,
				RiskPositionPercentage:   dbTrade.RiskPositionPercentage,
				ShortPeriod:              dbTrade.ShortPeriod,
				LongPeriod:               dbTrade.LongPeriod,
			}
			writeResponseWithData("TradingSystem Read successfully", dataTrade, conn)
		}
	case "update":
		if message.Entity == "trading-system" {
			var ts model.TradingSystemData
			dataByte, _ := json.Marshal(message.Data)
			// Deserialize the WebSocket message directly into the struct
			if err := json.Unmarshal(dataByte, &ts); err != nil {
				msg = fmt.Sprintf("Error6 parsing WebSocket message: %v", err)
				writeResponseWithID(msg, ts.ID, conn)
				return
			}
			// Fetch the existing trading system from the database based on tradeID
			existingTrade, err := DBServices.ReadTradingSystem(ts.ID)
			if err != nil {
				msg = fmt.Sprintf("Error retrieving trading system for update: %v", err)
				fmt.Printf("Error retrieving trading system %d for update: %v", ts.ID, err)
				writeResponseWithID(msg, ts.ID, conn)
				return
			} else if ts.ID != existingTrade.ID {
				msg = fmt.Sprintf("Error retrieving trading system for update: ts.ID %d != existingTrade.ID %d %v", ts.ID, existingTrade.ID, err)
				fmt.Printf("Error retrieving trading system %d for update: ts.ID %d != existingTrade.ID %d %v", ts.ID, ts.ID, existingTrade.ID, err)
				writeResponseWithID(msg, ts.ID, conn)
				return
			}
			// Update the existing trading system fields with new data
			existingTrade.Symbol = ts.Symbol
			existingTrade.ClosingPrices = model.Float64Slice(ts.ClosingPrices)
			existingTrade.Timestamps = model.Int64Slice(ts.Timestamps)
			existingTrade.Signals = model.StringSlice(ts.Signals)
			existingTrade.NextInvestBuYPrice = model.Float64Slice(ts.NextInvestBuYPrice)
			existingTrade.NextProfitSeLLPrice = model.Float64Slice(ts.NextProfitSeLLPrice)
			existingTrade.CommissionPercentage = ts.CommissionPercentage
			existingTrade.InitialCapital = ts.InitialCapital
			existingTrade.PositionSize = ts.PositionSize
			existingTrade.EntryPrice = model.Float64Slice(ts.EntryPrice)
			existingTrade.InTrade = ts.InTrade
			existingTrade.QuoteBalance = ts.QuoteBalance
			existingTrade.BaseBalance = ts.BaseBalance
			existingTrade.RiskCost = ts.RiskCost
			existingTrade.DataPoint = ts.DataPoint
			existingTrade.CurrentPrice = ts.CurrentPrice
			existingTrade.EntryQuantity = model.Float64Slice(ts.EntryQuantity)
			existingTrade.EntryCostLoss = model.Float64Slice(ts.EntryCostLoss)
			existingTrade.TradeCount = ts.TradeCount
			existingTrade.TradingLevel = ts.TradingLevel
			existingTrade.ClosedWinTrades = ts.ClosedWinTrades
			existingTrade.EnableStoploss = ts.EnableStoploss
			existingTrade.StopLossTrigered = ts.StopLossTrigered
			existingTrade.StopLossRecover = model.Float64Slice(ts.StopLossRecover)
			existingTrade.RiskFactor = ts.RiskFactor
			existingTrade.MaxDataSize = ts.MaxDataSize
			existingTrade.RiskProfitLossPercentage = ts.RiskProfitLossPercentage
			existingTrade.BaseCurrency = ts.BaseCurrency
			existingTrade.QuoteCurrency = ts.QuoteCurrency
			existingTrade.MiniQty = ts.MiniQty
			existingTrade.MaxQty = ts.MaxQty
			existingTrade.MinNotional = ts.MinNotional
			existingTrade.StepSize = ts.StepSize
			existingTrade.TargetStopLoss = ts.TargetStopLoss
			existingTrade.TargetProfit = ts.TargetProfit
			existingTrade.TotalProfitLoss = ts.TotalProfitLoss
			existingTrade.RiskPositionPercentage = ts.RiskPositionPercentage
			existingTrade.ShortPeriod = ts.ShortPeriod
			existingTrade.LongPeriod = ts.LongPeriod

			// Save the updated trading system back to the database
			err = DBServices.UpdateTradingSystem(existingTrade)
			if err != nil {
				msg = fmt.Sprintf("Error updating trading system: %v", err)
				writeResponseWithID(msg, existingTrade.ID, conn)
				return
			}
			writeResponseWithID("Trading system updated successfully", existingTrade.ID, conn)
		}
	case "delete":
		if message.Entity == "trading-system" {
			var ts model.TradingSystemData
			dataByte, _ := json.Marshal(message.Data)
			// Deserialize the WebSocket message directly into the struct
			if err := json.Unmarshal(dataByte, &ts); err != nil {
				msg = fmt.Sprintf("Error8 parsing WebSocket message: %v", err)
				return
			}
			// Delete the trading system from the database based on tradeID
			err := DBServices.DeleteTradingSystem(ts.ID)
			if err != nil {
				msg = fmt.Sprintf("Error deleting trading system: %v", err)
				return
			}
			// Send a success response back to the client via the conn
			response := map[string]interface{}{
				"message": "Trading system deleted successfully",
				"data_id": ts.ID,
			}
			err = conn.WriteJSON(response)
			if err != nil {
				msg = fmt.Sprintf("Error sending response via WebSocket: %v", err)
				return
			}
		}
	default:
		log.Printf("Unknown action: %s", message.Action)
		// msg = fmt.Sprintf("Invalid action in WebSocket message")
	}
}
func writeResponseWithID(msg string, id uint, conn *websocket.Conn) {
	// Send the dataID back to the client via the conn
	response := map[string]interface{}{
		"message": msg,
		"data_id": id,
	}
	err := conn.WriteJSON(response)
	if err != nil {
		log.Println("Error sending response via WebSocket:", err)
		return
	}
}
func writeResponseWithData(msg string, data interface{}, conn *websocket.Conn) {
	// Serialize the data object to JSON
	dataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println("Error marshaling Data to JSON:", err)
		return
	}
	// Send the data back to the client via the conn
	response := map[string]interface{}{
		"message": msg,
		"data":    json.RawMessage(dataJSON),
	}
	err = conn.WriteJSON(response)
	if err != nil {
		log.Println("Error sending response via WebSocket:", err)
		return
	}
}
