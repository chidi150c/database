package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

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
type WebSocketService struct {	
	Upgrader websocket.Upgrader
}

// parseTemplate applies a given file to the body of the base template.
func NewWebSocketService(HostSite string) WebSocketService {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {			
	// 		if r.Header.Get("Origin") != h.HostSite {
	// 			return false
	// 		}
			return true 
		},
	}
	return WebSocketService{
		Upgrader: upgrader,
	}
}

type TradeHandler struct {
    mux        *chi.Mux
    WebSocket  WebSocketService
    DBServices *gorm.DBServices
}

func NewTradeHandler(dBServices *gorm.DBServices, webSocketService WebSocketService) TradeHandler {
    h := TradeHandler{
        mux:        chi.NewRouter(),
		WebSocket:  webSocketService,
        DBServices: dBServices,
    }
	h.mux.Get("/database-services/ws", h.DataBaseSocketHandler)
    return h
}

func (h TradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    h.mux.ServeHTTP(w, r)
}

func (th *TradeHandler) DataBaseSocketHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade the HTTP connection to a WebSocket connection
    conn, err := th.WebSocket.Upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
        return
    }
    defer conn.Close()

    for {
        _, p, err := conn.ReadMessage()
        if err != nil {
            log.Println( "WebSocket read error:", err)
            return
        }

        // Assuming p is the received message (in JSON)

        // Unmarshal the WebSocket message (p) into a WebSocketMessage struct
        var message WebSocketMessage
        if err := json.Unmarshal(p, &message); err != nil {
            log.Println( "Error parsing WebSocket message:", err)
            continue
        }

        action := message.Action
        entity := message.Entity
        data := message.Data
		msg := ""
        // Handle different actions and entities here
        if action == "create" {
            if entity == "trading-system" {
                var ts model.TradingSystem				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					writeResponseWithID(msg, ts.ID, conn)
					continue
				}
				// Insert the new trading system into the database
				tradeID, err := th.DBServices.CreateTradingSystem(&ts)
				if err != nil {
					msg = fmt.Sprintf("Error creating trading system: %v", err)
					writeResponseWithID(msg, tradeID, conn)
					continue
				}
				writeResponseWithID("TradingSystem Created successfully", tradeID, conn)
			} else if entity == "app-data" {
                var appData model.AppData
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &appData); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					writeResponseWithID(msg, appData.ID, conn)
					continue
				}

				// Insert the new app data into the database
				dataID, err := th.DBServices.CreateAppData(&appData)
				if err != nil {
					msg = fmt.Sprintf("Error creating app data: %v", err)
					writeResponseWithID(msg, dataID, conn)
					return
				}
				writeResponseWithID("AppData Created successfully", dataID, conn)
            }
        } else if action == "read" {
            if entity == "trading-system" {
                var ts model.TradingSystem
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					writeResponseWithData(msg, ts, conn)
					continue
				}
                tradeID := ts.ID
                // Fetch the trading system from the database based on tradeID
                trade, err := th.DBServices.ReadTradingSystem(tradeID)
                if err != nil {
                    msg = fmt.Sprintf("Error retrieving trading system: %v", err)
					writeResponseWithData(msg, trade, conn)
                    return
                }
				writeResponseWithData("TradingSystem Read successfully", trade, conn)
            } else if entity == "app-data" {
                var ap model.AppData
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ap); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					writeResponseWithData(msg, ap, conn)
					continue
				}
				// Fetch the app data from the database based on dataID
				appData, err := th.DBServices.ReadAppData(ap.ID)
				if err != nil {
					msg = fmt.Sprintf("Error retrieving app data: %v", err)
					writeResponseWithData(msg, appData, conn)
					return
				}
				writeResponseWithData("AppData Read successfully", appData, conn)				
            }
        } else if action == "update" {
            if entity == "trading-system" {
                var ts model.TradingSystem				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					writeResponseWithID(msg, ts.ID, conn)
					continue
				}

                // Fetch the existing trading system from the database based on tradeID
                existingTrade, err := th.DBServices.ReadTradingSystem(ts.ID)
                if err != nil {
                    msg = fmt.Sprintf("Error retrieving trading system for update: %v", err)
					writeResponseWithID(msg, existingTrade.ID, conn)
                    return
                }
                // Update the existing trading system fields with new data
                existingTrade.Symbol = ts.Symbol
				existingTrade.ClosingPrices = ts.ClosingPrices
				existingTrade.Timestamps = ts.Timestamps   
				existingTrade.Signals = ts.Signals       
				existingTrade.NextInvestBuYPrice = ts.NextInvestBuYPrice
				existingTrade.NextProfitSeLLPrice = ts.NextProfitSeLLPrice
				existingTrade.CommissionPercentage = ts.CommissionPercentage
				existingTrade.InitialCapital = ts.InitialCapital
				existingTrade.PositionSize = ts.PositionSize
				existingTrade.EntryPrice = ts.EntryPrice
				existingTrade.InTrade = ts.InTrade
				existingTrade.QuoteBalance = ts.QuoteBalance
				existingTrade.BaseBalance = ts.BaseBalance
				existingTrade.RiskCost = ts.RiskCost
				existingTrade.DataPoint = ts.DataPoint
				existingTrade.CurrentPrice = ts.CurrentPrice
				existingTrade.EntryQuantity = ts.EntryQuantity
				existingTrade.Scalping = ts.Scalping
				existingTrade.StrategyCombLogic = ts.StrategyCombLogic
				existingTrade.EntryCostLoss = ts.EntryCostLoss
				existingTrade.TradeCount = ts.TradeCount
				existingTrade.EnableStoploss = ts.EnableStoploss
				existingTrade.StopLossTrigered = ts.StopLossTrigered
				existingTrade.StopLossRecover = ts.StopLossRecover
				existingTrade.RiskFactor = ts.RiskFactor
				existingTrade.MaxDataSize = ts.MaxDataSize
				existingTrade.RiskProfitLossPercentage = ts.RiskProfitLossPercentage
				existingTrade.BaseCurrency = ts.BaseCurrency
				existingTrade.QuoteCurrency = ts.QuoteCurrency
				existingTrade.MiniQty = ts.MiniQty
				existingTrade.MaxQty = ts.MaxQty
				existingTrade.MinNotional = ts.MinNotional
				existingTrade.StepSize = ts.StepSize

                // Save the updated trading system back to the database
                err = th.DBServices.UpdateTradingSystem(existingTrade)
                if err != nil {
                    msg = fmt.Sprintf("Error updating trading system: %v", err)
					writeResponseWithID(msg, existingTrade.ID, conn)
                    return
                }
				writeResponseWithID("Trading system updated successfully", existingTrade.ID, conn)
            } else if entity == "app-data" {
                var ap model.AppData				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ap); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					writeResponseWithID(msg, ap.ID, conn)
					continue
				}
		
				// Fetch the existing app data from the database based on dataID
				existingAppData, err := th.DBServices.ReadAppData(ap.ID)
				if err != nil {
					msg = fmt.Sprintf("Error retrieving app data for update: %v", err)
					writeResponseWithID(msg, existingAppData.ID, conn)
					return
				}
		
				// Update the existing app data fields with new data
				existingAppData.DataPoint = ap.DataPoint
				existingAppData.Strategy = ap.Strategy
				existingAppData.ShortPeriod = ap.ShortPeriod
				existingAppData.LongPeriod = ap.LongPeriod
				existingAppData.ShortEMA = ap.ShortEMA
				existingAppData.LongEMA = ap.LongEMA
				existingAppData.TargetProfit = ap.TargetProfit
				existingAppData.TargetStopLoss = ap.TargetStopLoss
				existingAppData.RiskPositionPercentage = ap.RiskPositionPercentage
				existingAppData.TotalProfitLoss = ap.TotalProfitLoss
		
				// Save the updated app data back to the database
				err = th.DBServices.UpdateAppData(existingAppData)
				if err != nil {
					msg = fmt.Sprintf("Error updating app data: %v", err)
					writeResponseWithID(msg, existingAppData.ID, conn)
					return
				}
				writeResponseWithID("App data updated successfully", existingAppData.ID, conn)
            }
        } else if action == "delete" {
            if entity == "trading-system" {
                var ts model.AppData				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					continue
				}
                // Delete the trading system from the database based on tradeID
                err := th.DBServices.DeleteTradingSystem(ts.ID)
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
            } else if entity == "app-data" {
                var ap model.AppData				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ap); err != nil {
					msg = fmt.Sprintf("Error parsing WebSocket message: %v", err)
					continue
				}		
				// Delete the app data from the database based on dataID
				err := th.DBServices.DeleteAppData(ap.ID)
				if err != nil {
					msg = fmt.Sprintf("Error deleting app data: %v", err)
					return
				}		
				// Send a success response back to the client via the conn
				response := map[string]interface{}{
					"message": "App data deleted successfully",
					"data_id": ap.ID,
				}		
				err = conn.WriteJSON(response)
				if err != nil {
					msg = fmt.Sprintf("Error sending response via WebSocket: %v", err)
					return
				}
            }
        } else {
            msg = fmt.Sprintf("Invalid action in WebSocket message")
        }
    }
}
func writeResponseWithID(msg string, id uint, conn *websocket.Conn){
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
func writeResponseWithData(msg string, data interface{}, conn *websocket.Conn){
	// Serialize the AppData object to JSON
	appDataJSON, err := json.Marshal(data)
	if err != nil {
		log.Println( "Error marshaling Data to JSON:", err)
		return
	}
	// Send the dataID back to the client via the conn
	response := map[string]interface{}{
		"message": msg,
		"data": json.RawMessage(appDataJSON),
	}
	err = conn.WriteJSON(response)
	if err != nil {
		log.Println( "Error sending response via WebSocket:", err)
		return
	}
}