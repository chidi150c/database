package server

import (
	"testing"

	"github.com/gorilla/websocket"
)

func TestCreateTradingSystem(t *testing.T) {
    // Create a mock WebSocket connection
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:35260/database-services/ws", nil)
    if err != nil {
        t.Fatalf("Failed to connect to WebSocket: %v", err)
    }
    defer conn.Close()

    // Create a message (request) to send
    request := map[string]interface{}{
        "action": "create",
        "entity": "trading-system",
        "data": map[string]interface{}{
            "Symbol":        "AAPL",
            "ClosingPrices": []float64{150.0, 151.0, 152.0},
        },
    }

    // Send the message
    err = conn.WriteJSON(request)
    if err != nil {
        t.Fatalf("Failed to send WebSocket message: %v", err)
    }

    // Receive and parse the response
    var response map[string]interface{}
    err = conn.ReadJSON(&response)
    if err != nil {
        t.Fatalf("Failed to read WebSocket response: %v", err)
    }

    // Perform assertions to verify the response
    tradeID, ok := response["trade_id"].(float64)
    if !ok {
        t.Fatalf("Invalid response format")
    }
	if tradeID != 123 {
		t.Fatal("tradrID incorrect")
	}
    // Add more assertions as needed

    // Optionally, clean up any test data created during the test
}
