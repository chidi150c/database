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

	// Handle WebSocket close event
	go func() {
		for {
			_, _, err := conn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway) {
					// WebSocket connection is closed
					closedChannel <- struct{}{}
					return
				}
				log.Println("WebSocket read error:", err)
				return
			}
		}
	}()

	// Periodically fetch and send TradingSystem records to the client
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	go func() {
		// You can fetch the tradeID and AppDataID from the user or HTTP request
		tradeID := "your_trade_id"
		AppDataID := "your_appdata_id"

		for {
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
			// Sleep for a while (you can adjust the duration)
			time.Sleep(5 * time.Second)
		}
	}()

	// Block until the WebSocket connection is closed
	<-closedChannel
}
