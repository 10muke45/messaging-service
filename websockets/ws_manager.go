package websockets

import (
	"sync"

	"github.com/gorilla/websocket"
)

type WSManager struct {
	Clients map[string]*websocket.Conn
	mu      sync.Mutex
}

var WSM *WSManager

func InitWSManager() {
	WSM = &WSManager{
		Clients: make(map[string]*websocket.Conn),
	}
}

func (wsm *WSManager) AddClient(username string, conn *websocket.Conn) {
	wsm.mu.Lock()
	wsm.Clients[username] = conn
	wsm.mu.Unlock()
}

func (wsm *WSManager) RemoveClient(username string) {
	wsm.mu.Lock()
	delete(wsm.Clients, username)
	wsm.mu.Unlock()
}
