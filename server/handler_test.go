package server

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/chidi150c/database/model"
	"github.com/gorilla/websocket"
)

func TestCreateAppData(t *testing.T) {
    // Create a mock WebSocket connection
    conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:35261/database-services/ws", nil)
    if err != nil {
        t.Fatalf("Failed to connect to WebSocket: %v", err)
    }
    defer conn.Close()

    // Create a AppData object
    ap := model.AppData{
        DataPoint: 234,
		Strategy: "EMA",
		ShortPeriod: 65,
    }
	ap.ID = uint(2)

	// Serialize the AppData object to JSON
	appDataJSON, err := json.Marshal(ap)
	if err != nil {
		t.Fatalf("Error marshaling AppData to JSON: %v", err)
	}

	// Create a message (request) to send
	request := map[string]interface{}{
		"action": "create",
		"entity": "app-data",
		"data":   json.RawMessage(appDataJSON), // RawMessage to keep it as JSON
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
    ms, ok := response["message"].(string)
    if !ok {
        t.Fatalf("Invalid response format %v and resp: %v", ms, response)
    }
	fmt.Printf("tradeID in uint correct: %v\n", ms)

    // Add more assertions as needed
	tradeID := uint(2)
	// Create a message (request) to send
	request = map[string]interface{}{
		"action": "read",
		"entity": "app-data",
		"data": map[string]interface{}{
			"data_id":       tradeID,
		},
	}

	// Send the message
	err = conn.WriteJSON(request)
	if err != nil {
		t.Fatalf("Failed to send WebSocket message: %v", err)
	}

	err = conn.ReadJSON(&response)
	if err != nil {
		t.Fatalf("Failed to read WebSocket response: %v", err)
	}
	// Marshal the map into a JSON byte slice
	tradingSystemBytes, err := json.Marshal(response)
	if err != nil {
		t.Fatalf("Error marshaling AppData: %v", err)
	}

	// Unmarshal the JSON byte slice into a AppData struct
	var responseAppData model.AppData
	err = json.Unmarshal(tradingSystemBytes, &responseAppData)
	if err != nil {
		t.Fatalf("Failed to unmarshal WebSocket response: %v", err)
	}

	fmt.Printf("AppData correct: %+v", responseAppData)

}
