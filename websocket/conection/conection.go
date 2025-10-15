package conection

import (
	"sync"

	"github.com/gorilla/websocket"
	"main.go/models"
)

type Client struct {
	conn *websocket.Conn
	send chan models.Evento
}

type ServerState struct {
	clients       map[*Client]bool
	clientsMutex  sync.Mutex
	courts        map[string]models.Quadra
	reservations  map[string]models.Reserva
	stateMutex    sync.Mutex
	nextCourtID   int
	nextReserveID int64
	systemStarted bool
}

func NewServerState() *ServerState {
	return &ServerState{
		clients:      make(map[*Client]bool),
		courts:       make(map[string]models.Quadra),
		reservations: make(map[string]models.Reserva),
	}
}

func (s *ServerState) broadcast(ev models.Evento) {
	s.clientsMutex.Lock()
	defer s.clientsMutex.Unlock()
	for c := range s.clients {
		select {
		case c.send <- ev:
		default:
			close(c.send)
			delete(s.clients, c)
		}
	}
}
