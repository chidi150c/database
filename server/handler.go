package server

import (
	"encoding/json"
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
            log.Println("WebSocket read error:", err)
            return
        }

        // Assuming p is the received message (in JSON)

        // Unmarshal the WebSocket message (p) into a WebSocketMessage struct
        var message WebSocketMessage
        if err := json.Unmarshal(p, &message); err != nil {
            log.Println("Error parsing WebSocket message:", err)
            continue
        }

        action := message.Action
        entity := message.Entity
        data := message.Data

        // Handle different actions and entities here
        if action == "create" {
            if entity == "trading-system" {
                var ts model.TradingSystem				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}
				// Insert the new trading system into the database
				tradeID, err := th.DBServices.CreateTradingSystem(&ts)
				if err != nil {
					log.Println("Error creating trading system:", err)
					continue
				}
				// Send the tradeID back to the client via the conn
				response := map[string]interface{}{
					"message": "TradingSystem created successfully",
					"data_id": tradeID,
				}
				err = conn.WriteJSON(response)
				if err != nil {
					log.Println("Error sending response via WebSocket:", err)
					continue
				}	
				
			} else if entity == "app-data" {
                var appData model.AppData
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &appData); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}

				// Insert the new app data into the database
				dataID, err := th.DBServices.CreateAppData(&appData)
				if err != nil {
					log.Println("Error creating app data:", err)
					return
				}

				// Send the dataID back to the client via the conn
				response := map[string]interface{}{
					"message": "AppData Created successfully",
					"data_id": dataID,
					"error": errMessage,
				}

				err = conn.WriteJSON(response)
				if err != nil {
					log.Println("Error sending response via WebSocket:", err)
					return
				}
            }
        } else if action == "read" {
            if entity == "trading-system" {
                var ts model.TradingSystem
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}
                tradeID := ts.ID
                // Fetch the trading system from the database based on tradeID
                trade, err := th.DBServices.ReadTradingSystem(tradeID)
                if err != nil {
                    log.Println("Error retrieving trading system:", err)
                    return
                }
				// Serialize the AppData object to JSON
				appDataJSON, err := json.Marshal(trade)
				if err != nil {
					log.Println("Error marshaling TradingSystem to JSON:", err)
					return
				}
				// Send the tradeID back to the client via the conn
				response := map[string]interface{}{
					"message": "TradingSystem Read successfully",
					"data":   json.RawMessage(appDataJSON), // RawMessage to keep it as JSON
				}

                // Send the trading system data to the client via the conn
                err = conn.WriteJSON(response)
                if err != nil {
                    log.Println("Error sending trading system data via WebSocket:", err)
                    return
                }
            } else if entity == "app-data" {
                var ap model.AppData
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ap); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}


				// Fetch the app data from the database based on dataID
				appData, err := th.DBServices.ReadAppData(ap.ID)
				if err != nil {
					log.Println("Error retrieving app data:", err)
					return
				}
				// Serialize the AppData object to JSON
				appDataJSON, err := json.Marshal(appData)
				if err != nil {
					log.Println("Error marshaling AppData to JSON:", err)
					return
				}
				// Send the tradeID back to the client via the conn
				response := map[string]interface{}{
					"message": "TradingSystem Read successfully",
					"data":   json.RawMessage(appDataJSON), // RawMessage to keep it as JSON
				}
				// Send the app data to the client via the conn
				err = conn.WriteJSON(response)
				if err != nil {
					log.Println("Error sending app data via WebSocket:", err)
					return
				}
            }
        } else if action == "update" {
            if entity == "trading-system" {
                var ts model.TradingSystem				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}

                // Fetch the existing trading system from the database based on tradeID
                existingTrade, err := th.DBServices.ReadTradingSystem(ts.ID)
                if err != nil {
                    log.Println("Error retrieving trading system for update:", err)
                    return
                }
                // Update the existing trading system fields with new data
                existingTrade.Symbol = ts.Symbol
				existingTrade.ClosingPrices = ts.ClosingPrices
				existingTrade.Container1 = ts.Container1
				existingTrade.Container2 = ts.Container2
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
                    log.Println("Error updating trading system:", err)
                    return
                }

                // Send a success response back to the client via the conn
                response := map[string]interface{}{
                    "message": "Trading system updated successfully",
					"data_id": existingTrade.ID,
                }

                err = conn.WriteJSON(response)
                if err != nil {
                    log.Println("Error sending response via WebSocket:", err)
                    return
                }
            } else if entity == "app-data" {
                var ap model.AppData				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ap); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}
		
				// Fetch the existing app data from the database based on dataID
				existingAppData, err := th.DBServices.ReadAppData(ap.ID)
				if err != nil {
					log.Println("Error retrieving app data for update:", err)
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
					log.Println("Error updating app data:", err)
					return
				}
		
				// Send a success response back to the client via the conn
				response := map[string]interface{}{
					"message": "App data updated successfully",
					"data_id": existingAppData.ID,
				}
		
				err = conn.WriteJSON(response)
				if err != nil {
					log.Println("Error sending response via WebSocket:", err)
					return
				}
            }
        } else if action == "delete" {
            if entity == "trading-system" {
                var ts model.AppData				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ts); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}

                // Delete the trading system from the database based on tradeID
                err := th.DBServices.DeleteTradingSystem(ts.ID)
                if err != nil {
                    log.Println("Error deleting trading system:", err)
                    return
                }

                // Send a success response back to the client via the conn
                response := map[string]interface{}{
                    "message": "Trading system deleted successfully",
					"data_id": ts.ID,
                }

                err = conn.WriteJSON(response)
                if err != nil {
                    log.Println("Error sending response via WebSocket:", err)
                    return
                }
            } else if entity == "app-data" {
                var ap model.AppData				
				dataByte, _ := json.Marshal(data)
				// Deserialize the WebSocket message directly into the struct
				if err := json.Unmarshal(dataByte, &ap); err != nil {
					log.Println("Error parsing WebSocket message:", err)
					continue
				}
		
				// Delete the app data from the database based on dataID
				err := th.DBServices.DeleteAppData(ap.ID)
				if err != nil {
					log.Println("Error deleting app data:", err)
					return
				}
		
				// Send a success response back to the client via the conn
				response := map[string]interface{}{
					"message": "App data deleted successfully",
					"data_id": ap.ID,
				}
		
				err = conn.WriteJSON(response)
				if err != nil {
					log.Println("Error sending response via WebSocket:", err)
					return
				}
            }
        } else {
            log.Println("Invalid action in WebSocket message")
        }
    }
}
