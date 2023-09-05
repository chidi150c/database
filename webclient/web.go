package webclient

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

// WebService is a user login-aware wrapper for a html/template.
type WebService struct {
}

// parseTemplate applies a given file to the body of the base template.
func NewWebService() *WebService {
	return &WebService{}
}

// Execute writes the template using the provided data, adding login and user
// information to the base template.
func (resp *WebService) Execute(w http.ResponseWriter, r *http.Request, data interface{}, usr interface{}, msg interface{}) error {
	d := struct {
		Data      interface{} `json:"data"`
		LogoutURL string      `json:"logoutUrl"`
		User      interface{} `json:"user"`
		Msg       interface{} `json:"msg"`
	}{
		Data:      data,
		LogoutURL: "/logout?redirect=" + r.URL.RequestURI(),
		User:      usr,
		Msg:       msg,
	}
	
	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(d)
	if err != nil {
		log.Printf("mobileWarning %v\n", err)
	}
	return err
}

// WebService is a user login-aware wrapper for a html/template.
type SocketService struct {	
	Upgrader websocket.Upgrader
}

// parseTemplate applies a given file to the body of the base template.
func NewSocketService(HostSite string) *SocketService {
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
	return &SocketService{
		Upgrader: upgrader,
	}
}