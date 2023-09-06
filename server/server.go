package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/chidi150c/database/gorm"
	"github.com/chidi150c/database/model"
	"github.com/chidi150c/database/webclient"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

type TradeHandler struct {
	mux        *chi.Mux
	RESTAPI    *webclient.WebService
	WebSocket  *webclient.SocketService // Updated to use webclient.SocketService
	HostSite   string
	DBServices *gorm.DBServices
}

func NewTradeHandler(dBServices *gorm.DBServices, HostSite string) TradeHandler {
	h := TradeHandler{
		mux:        chi.NewRouter(),
		RESTAPI:    webclient.NewWebService(),
		WebSocket:  webclient.NewSocketService(HostSite), // Updated to use webclient.SocketService
		HostSite:   os.Getenv("HOSTSITE"),
		DBServices: dBServices,
	}

	h.mux.Get("/trading-system", h.UpdateTradingSystem)
	h.mux.Get("/trading-system/ws", h.tradingSystemSocketHandler)
	h.mux.Get("/app-data", h.CreateAppData)

	return h
}

func (h TradeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.mux.ServeHTTP(w, r)
}

func (th *TradeHandler) CreateTradingSystem(w http.ResponseWriter, r *http.Request) {
	var trade model.TradingSystem
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	tradeID, _, err := th.DBServices.Create(&trade, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"trade_id": tradeID})
}

func (th *TradeHandler) UpdateTradingSystem(w http.ResponseWriter, r *http.Request) {
	var trade model.TradingSystem
	if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := th.DBServices.Update(&trade, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (th *TradeHandler) CreateAppData(w http.ResponseWriter, r *http.Request) {
	var data model.AppData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	_, dataID, err := th.DBServices.Create(nil, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]int{"data_id": dataID})
}

func (th *TradeHandler) UpdateAppData(w http.ResponseWriter, r *http.Request) {
	var data model.AppData
	if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := th.DBServices.Update(nil, &data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}
func (th *TradeHandler) tradingSystemSocketHandler(w http.ResponseWriter, r *http.Request) {
    // Upgrade the HTTP connection to a WebSocket connection
	conn, err := th.WebSocket.Upgrader.Upgrade(w, r, nil)
    if err != nil {
        http.Error(w, "Could not upgrade to WebSocket", http.StatusBadRequest)
        return
    }
    defer conn.Close()

	// Create a channel to signal when the WebSocket connection is closed
	closedChannel := make(chan struct{})

    for {
        messageType, p, err := conn.ReadMessage()
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
				// Fetch the trading system from the database based on tradeID and AppDataID
				ts := &model.TradingSystem{
					Symbol: data["field1"].(string),
					ClosingPrices: data["field2"].([]float64),
				}
				tradeID, _, err := th.DBServices.Create(ts, nil)
				if err != nil {
					log.Println("Error fetching trading system:", err)
					return
				}

			// Send the trading system data to the client via the conn
			err = conn.WriteJSON(trade)
			if err != nil {
				log.Println("Error sending trading system data via WebSocket:", err)
				return
			}
                // Perform database operations and send response back to client
            } else if entity == "app-data" {
                // Handle create app data
                // You can access data["field1"], data["field2"], etc.
                // Perform database operations and send response back to client
            }
        } else if action == "retrieve" {
            // Fetch the trading system from the database based on tradeID and AppDataID
			trade, _, err := th.DBServices.Read(tradeID, AppDataID)
			if err != nil {
				log.Println("Error fetching trading system:", err)
				return
			}

			// Send the trading system data to the client via the conn
			err = conn.WriteJSON(trade)
			if err != nil {
				log.Println("Error sending trading system data via WebSocket:", err)
				return
			}
        } else {
            log.Println("Invalid action in WebSocket message")
        }

        // Send responses back to the client as needed
    }
}

