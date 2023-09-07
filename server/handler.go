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

        // Assuming p is the received message (a JSON string)

        var request map[string]interface{}
        if err := json.Unmarshal(p, &request); err != nil {
            log.Println("Error parsing WebSocket message:", err)
            continue
        }

        action, ok := request["action"].(string)
        if !ok {
            log.Println("Invalid action in WebSocket message")
            continue
        }

        entity, ok := request["entity"].(string)
        if !ok {
            log.Println("Invalid entity in WebSocket message")
            continue
        }

        data, ok := request["data"].(map[string]interface{})
        if !ok {
            log.Println("Invalid data in WebSocket message")
            continue
        }

        // Handle different actions and entities here
        if action == "create" {
            if entity == "trading-system" {
                // Handle create trading system
                // You can access data["field1"], data["field2"], etc.

                // Create a new TradingSystem object with the provided data
                ts := &model.TradingSystem{
                    Symbol:        data["Symbol"].(string),
                    ClosingPrices: data["ClosingPrices"].([]float64),
                    // Add other fields as needed
                }

                // Insert the new trading system into the database
                tradeID, err := th.DBServices.CreateTradingSystem(ts)
                if err != nil {
                    log.Println("Error creating trading system:", err)
                    return
                }

                // Send the tradeID back to the client via the conn
                response := map[string]interface{}{
                    "trade_id": tradeID,
                }

                err = conn.WriteJSON(response)
                if err != nil {
                    log.Println("Error sending response via WebSocket:", err)
                    return
                }
            } else if entity == "app-data" {
                // Handle create app data
                // You can access data["field1"], data["field2"], etc.

                // Create a new AppData object with the provided data
                appData := &model.AppData{
                    DataPoint:      data["DataPoint"].(int),
                    Strategy:       data["Strategy"].(string),
                    ShortPeriod:    data["ShortPeriod"].(int),
                    LongPeriod:     data["LongPeriod"].(int),
                    ShortEMA:       data["ShortEMA"].(float64),
                    LongEMA:        data["LongEMA"].(float64),
                    // Add other fields as needed
                }

                // Insert the new app data into the database
                dataID, err := th.DBServices.CreateAppData(appData)
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
        } else if action == "retrieve" {
            // Handle retrieve operation
            // You can access data["trade_id"] or data["data_id"] to retrieve data from the database

            // Fetch the trading system or app data based on tradeID or dataID
            if entity == "trading-system" {
                tradeID := data["trade_id"].(int)

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
                dataID := data["data_id"].(int)

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
        } else {
            log.Println("Invalid action in WebSocket message")
        }
    }
}


