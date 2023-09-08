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

func NewTradeHandler(dBServices *gorm.DBServices, webSocketService WebSocketService, HostSite string) TradeHandler {
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
					"data_id": dataID,
				}

				err = conn.WriteJSON(response)
				if err != nil {
					log.Println("Error sending response via WebSocket:", err)
					return
				}
            }
        } else if action == "read" {
            if entity == "trading-system" {
                tradeID := uint(data["data_id"].(float64))

                // Fetch the trading system from the database based on tradeID
                trade, err := th.DBServices.ReadTradingSystem(tradeID)
                if err != nil {
                    log.Println("Error retrieving trading system:", err)
                    return
                }

                // Send the trading system data to the client via the conn
                err = conn.WriteJSON(trade)
                if err != nil {
                    log.Println("Error sending trading system data via WebSocket:", err)
                    return
                }
            } else if entity == "app-data" {
				dataID := uint(data["data_id"].(float64))

				// Fetch the app data from the database based on dataID
				appData, err := th.DBServices.ReadAppData(dataID)
				if err != nil {
					log.Println("Error retrieving app data:", err)
					return
				}

				// Send the app data to the client via the conn
				err = conn.WriteJSON(appData)
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
				// existingTrade.Container1 = ts.Container1 // Example
				// existingTrade.Container2 = ts.Container2 // Example
				// existingTrade.Timestamps = ts.Timestamps"].(model.Int64Slice)   // Example
				// existingTrade.Signals = ts.Signals"].(model.StringSlice)        // Example
				// existingTrade.NextInvestBuYPrice = ts.NextInvestBuYPrice // Example
				// existingTrade.NextProfitSeLLPrice = ts.NextProfitSeLLPrice // Example
				// existingTrade.CommissionPercentage = ts.CommissionPercentage"].(float64)
				// existingTrade.InitialCapital = ts.InitialCapital"].(float64)
				// existingTrade.PositionSize = ts.PositionSize"].(float64)
				// existingTrade.EntryPrice = ts.EntryPrice // Example
				// existingTrade.InTrade = ts.InTrade"].(bool)
				// existingTrade.QuoteBalance = ts.QuoteBalance"].(float64)
				// existingTrade.BaseBalance = ts.BaseBalance"].(float64)
				// existingTrade.RiskCost = ts.RiskCost"].(float64)
				// existingTrade.DataPoint = int(ts.DataPoint"].(float64))
				// existingTrade.CurrentPrice = ts.CurrentPrice"].(float64)
				// existingTrade.EntryQuantity = ts.EntryQuantity // Example
				// existingTrade.Scalping = ts.Scalping
				// existingTrade.StrategyCombLogic = ts.StrategyCombLogic
				// existingTrade.EntryCostLoss = ts.EntryCostLoss // Example
				// existingTrade.TradeCount = int(ts.TradeCount"].(float64))
				// existingTrade.EnableStoploss = ts.EnableStoploss"].(bool)
				// existingTrade.StopLossTrigered = ts.StopLossTrigered"].(bool)
				// existingTrade.StopLossRecover = ts.StopLossRecover // Example
				// existingTrade.RiskFactor = ts.RiskFactor"].(float64)
				// existingTrade.MaxDataSize = int(ts.MaxDataSize"].(float64))
				// existingTrade.RiskProfitLossPercentage = ts.RiskProfitLossPercentage"].(float64)
				// existingTrade.BaseCurrency = ts.BaseCurrency
				// existingTrade.QuoteCurrency = ts.QuoteCurrency
				// existingTrade.MiniQty = ts.MiniQty"].(float64)
				// existingTrade.MaxQty = ts.MaxQty"].(float64)
				// existingTrade.MinNotional = ts.MinNotional"].(float64)
				// existingTrade.StepSize = ts.StepSize"].(float64)



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
				existingAppData.DataPoint = int(data["DataPoint"].(float64))
				existingAppData.Strategy = data["Strategy"].(string)
				existingAppData.ShortPeriod = int(data["ShortPeriod"].(float64))
				existingAppData.LongPeriod = int(data["LongPeriod"].(float64))
				existingAppData.ShortEMA = data["ShortEMA"].(float64)
				existingAppData.LongEMA = data["LongEMA"].(float64)
				existingAppData.ProfitLoss    Float64Slice `gorm:"type:real[]"`
				existingAppData.CapitalCurve  Float64Slice `gorm:"type:real[]"`
				existingAppData.WinLossRatio  Float64Slice `gorm:"type:real[]"`
				// Update other fields as needed
		
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
