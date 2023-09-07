package server

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/chidi150c/database/gorm"
	"github.com/chidi150c/database/model"
	"github.com/go-chi/chi"
	"github.com/gorilla/websocket"
)

// WebService is a user login-aware wrapper for a html/template.
type WebSocketService struct {	
	Upgrader websocket.Upgrader
}

type TradeHandler struct {
    mux        *chi.Mux
    WebSocket  WebSocketService
    HostSite   string
    DBServices *gorm.DBServices
}

func NewTradeHandler(dBServices *gorm.DBServices, webSocketService WebSocketService, HostSite string) TradeHandler {
    h := TradeHandler{
        mux:        chi.NewRouter(),
        WebSocket:  webSocketService,
        HostSite:   os.Getenv("HOSTSITE"),
        DBServices: dBServices,
    }

    h.mux.Get("/trading-system", h.CreateTradingSystem)
    h.mux.Get("/trading-system/{tradeID}", h.ReadTradingSystem)
    h.mux.Put("/trading-system/{tradeID}", h.UpdateTradingSystem)
    h.mux.Delete("/trading-system/{tradeID}", h.DeleteTradingSystem)

    h.mux.Get("/app-data", h.CreateAppData)
    h.mux.Get("/app-data/{dataID}", h.ReadAppData)
    h.mux.Put("/app-data/{dataID}", h.UpdateAppData)
    h.mux.Delete("/app-data/{dataID}", h.DeleteAppData)

    h.mux.Get("/trading-system/ws", h.tradingSystemSocketHandler)

    return h
}

// CreateTradingSystem handles create trading system operation
func (th *TradeHandler) CreateTradingSystem(w http.ResponseWriter, r *http.Request) {
    var trade model.TradingSystem
    if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    tradeID, err := th.DBServices.CreateTradingSystem(&trade)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"trade_id": tradeID})
}

// ReadTradingSystem handles retrieve trading system operation
func (th *TradeHandler) ReadTradingSystem(w http.ResponseWriter, r *http.Request) {
    tradeID := chi.URLParam(r, "tradeID")
    tradeIDInt, _ := strconv.Atoi(tradeID)

    trade, err := th.DBServices.ReadTradingSystem(tradeIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(trade)
}

// UpdateTradingSystem handles update trading system operation
func (th *TradeHandler) UpdateTradingSystem(w http.ResponseWriter, r *http.Request) {
    tradeID := chi.URLParam(r, "tradeID")
    tradeIDInt, _ := strconv.Atoi(tradeID)

    var trade model.TradingSystem
    if err := json.NewDecoder(r.Body).Decode(&trade); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    trade.ID = tradeIDInt
    err := th.DBServices.UpdateTradingSystem(&trade)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// DeleteTradingSystem handles delete trading system operation
func (th *TradeHandler) DeleteTradingSystem(w http.ResponseWriter, r *http.Request) {
    tradeID := chi.URLParam(r, "tradeID")
    tradeIDInt, _ := strconv.Atoi(tradeID)

    err := th.DBServices.DeleteTradingSystem(tradeIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// CreateAppData handles create app data operation
func (th *TradeHandler) CreateAppData(w http.ResponseWriter, r *http.Request) {
    var data model.AppData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    dataID, err := th.DBServices.CreateAppData(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusCreated)
    json.NewEncoder(w).Encode(map[string]int{"data_id": dataID})
}

// ReadAppData handles retrieve app data operation
func (th *TradeHandler) ReadAppData(w http.ResponseWriter, r *http.Request) {
    dataID := chi.URLParam(r, "dataID")
    dataIDInt, _ := strconv.Atoi(dataID)

    data, err := th.DBServices.ReadAppData(dataIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(data)
}

// UpdateAppData handles update app data operation
func (th *TradeHandler) UpdateAppData(w http.ResponseWriter, r *http.Request) {
    dataID := chi.URLParam(r, "dataID")
    dataIDInt, _ := strconv.Atoi(dataID)

    var data model.AppData
    if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    data.ID = dataIDInt
    err := th.DBServices.UpdateAppData(&data)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// DeleteAppData handles delete app data operation
func (th *TradeHandler) DeleteAppData(w http.ResponseWriter, r *http.Request) {
    dataID := chi.URLParam(r, "dataID")
    dataIDInt, _ := strconv.Atoi(dataID)

    err := th.DBServices.DeleteAppData(dataIDInt)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    w.WriteHeader(http.StatusOK)
}

// tradingSystemSocketHandler handles WebSocket connections
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
        for {
            select {
            case <-ticker.C:
                // Fetch the tradeID from the WebSocket client, you can implement a custom protocol for this
                // For now, I'll use a placeholder
                tradeID := "your_trade_id"

                // Fetch the trading system from the database based on tradeID
                trade, err := th.DBServices.ReadTradingSystem(tradeID)
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
            case <-closedChannel:
                // WebSocket connection is closed
                return
            }
        }
    }()

    // Block until the WebSocket connection is closed
    <-closedChannel
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