package main

import (
	"log"
	"net/http"

	"main.go/conection"
)

// --- Estruturas de dados ---
// type Event struct {
// 	Event string      `json:"event"`
// 	Data  interface{} `json:"data"`
// }

// type Court struct {
// 	ID       string `json:"id"`
// 	Name     string `json:"name"`
// 	Capacity int    `json:"capacity"`
// }

// type Reservation struct {
// 	ID        string `json:"id"`
// 	CourtID   string `json:"courtId"`
// 	User      string `json:"user"`
// 	StartTime int64  `json:"startTime"`
// 	ExpiresAt int64  `json:"expiresAt"`
// 	Status    string `json:"status"`
// }

// type Client struct {
// 	conn *websocket.Conn
// 	send chan Event
// }

// var upgrader = websocket.Upgrader{
// 	CheckOrigin: func(r *http.Request) bool { return true },
// }

// type Client struct {
// 	conn *websocket.Conn
// 	send chan models.Event
// }

// // --- Estado global do servidor ---
// type ServerState struct {
// 	clients       map[*Client]bool
// 	clientsMutex  sync.Mutex
// 	courts        map[string]models.Court
// 	reservations  map[string]models.Reservation
// 	stateMutex    sync.Mutex
// 	nextCourtID   int
// 	nextReserveID int64
// 	systemStarted bool
// }

// func NewServerState() *ServerState {
// 	return &ServerState{
// 		clients:      make(map[*Client]bool),
// 		courts:       make(map[string]models.Court),
// 		reservations: make(map[string]models.Reservation),
// 	}
// }

// --- Funções auxiliares ---
// func (s *ServerState) broadcast(ev models.Event) {
// 	s.clientsMutex.Lock()
// 	defer s.clientsMutex.Unlock()
// 	for c := range s.clients {
// 		select {
// 		case c.send <- ev:
// 		default:
// 			close(c.send)
// 			delete(s.clients, c)
// 		}
// 	}
// }

// // --- Leitura e escrita WebSocket ---
// func (c *Client) readLoop() {
// 	defer func() {
// 		c.conn.Close()
// 		state.clientsMutex.Lock()
// 		delete(state.clients, c)
// 		state.clientsMutex.Unlock()
// 	}()
// 	for {
// 		_, msg, err := c.conn.ReadMessage()
// 		if err != nil {
// 			log.Println("read error:", err)
// 			return
// 		}
// 		var ev models.Event
// 		if err := json.Unmarshal(msg, &ev); err != nil {
// 			log.Println("invalid event JSON:", err)
// 			continue
// 		}
// 		handleEvent(c, ev)
// 	}
// }

// func (c *Client) writeLoop() {
// 	for ev := range c.send {
// 		b, err := json.Marshal(ev)
// 		if err != nil {
// 			continue
// 		}
// 		if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
// 			log.Println("write error:", err)
// 			return
// 		}
// 	}
// }

// // --- Manipuladores de eventos ---
// func handleEvent(c *Client, ev models.Event) {
// 	switch ev.Event {
// 	case "sistema.iniciado":
// 		handleStartSystem()
// 	case "sistema.parado":
// 		handleStopSystem()
// 	case "quadra.adicionada":
// 		handleAddCourt(ev.Data)
// 	case "horario.reservado":
// 		handleReserve(ev.Data)
// 	case "reserva.cancelada":
// 		handleCancel(ev.Data)
// 	case "jogo.finalizado":
// 		handleFinalize(ev.Data)
// 	default:
// 		log.Println("Evento desconhecido:", ev.Event)
// 	}
// }

// func handleStartSystem() {
// 	state.stateMutex.Lock()
// 	state.systemStarted = true
// 	courts := make([]models.Court, 0, len(state.courts))
// 	for _, ct := range state.courts {
// 		courts = append(courts, ct)
// 	}
// 	state.stateMutex.Unlock()
// 	state.broadcast(models.Event{Event: "sistema.iniciado", Data: map[string]interface{}{"courts": courts}})
// }

// func handleAddCourt(data interface{}) {
// 	m, ok := data.(map[string]interface{})
// 	if !ok {
// 		return
// 	}
// 	name, _ := m["name"].(string)
// 	cap, _ := m["capacity"].(float64)
// 	if name == "" {
// 		return
// 	}

// 	state.stateMutex.Lock()
// 	state.nextCourtID++
// 	id := fmt.Sprintf("q%d", state.nextCourtID)
// 	ct := models.Court{ID: id, Name: name, Capacity: int(cap)}
// 	state.courts[id] = ct
// 	state.stateMutex.Unlock()

// 	state.broadcast(models.Event{Event: "quadra.adicionada", Data: ct})
// }

// func handleReserve(data interface{}) {
// 	m, ok := data.(map[string]interface{})
// 	if !ok {
// 		return
// 	}
// 	courtId, _ := m["courtId"].(string)
// 	user, _ := m["user"].(string)
// 	startFloat, _ := m["startTime"].(float64)
// 	start := int64(startFloat)

// 	state.stateMutex.Lock()

// 	// 🔒 BLOQUEIO: sistema precisa estar iniciado
// 	if !state.systemStarted {
// 		state.stateMutex.Unlock()
// 		state.broadcast(models.Event{Event: "reserva.negada", Data: map[string]interface{}{
// 			"reason":  "sistema não iniciado",
// 			"user":    user,
// 			"courtId": courtId,
// 		}})
// 		return
// 	}

// 	// Verifica conflito de horário
// 	for _, r := range state.reservations {
// 		if r.CourtID == courtId && r.StartTime == start && (r.Status == "reserved" || r.Status == "confirmed") {
// 			state.stateMutex.Unlock()
// 			state.broadcast(models.Event{Event: "reserva.negada", Data: map[string]interface{}{
// 				"reason":    "horário já reservado",
// 				"courtId":   courtId,
// 				"startTime": start,
// 			}})
// 			return
// 		}
// 	}

// 	// Cria reserva (expira em 5 minutos)
// 	state.nextReserveID++
// 	resID := fmt.Sprintf("r-%d", state.nextReserveID)
// 	exp := time.Now().Add(5 * time.Minute)
// 	res := models.Reservation{
// 		ID:        resID,
// 		CourtID:   courtId,
// 		User:      user,
// 		StartTime: start,
// 		ExpiresAt: exp.Unix(),
// 		Status:    "reserved",
// 	}
// 	state.reservations[resID] = res
// 	state.stateMutex.Unlock()

// 	state.broadcast(models.Event{Event: "horario.reservado", Data: res})

// 	// Expira após 5 minutos
// 	go func(id string) {
// 		timer := time.NewTimer(5 * time.Minute)
// 		<-timer.C
// 		state.stateMutex.Lock()
// 		r, ok := state.reservations[id]
// 		if ok && r.Status == "reserved" {
// 			r.Status = "expired"
// 			state.reservations[id] = r
// 			state.stateMutex.Unlock()
// 			state.broadcast(models.Event{Event: "reserva.expirada", Data: r})
// 			return
// 		}
// 		state.stateMutex.Unlock()
// 	}(resID)

// 	// Confirma após 5 segundos
// 	go func(id string) {
// 		time.Sleep(5 * time.Second)
// 		state.stateMutex.Lock()
// 		r, ok := state.reservations[id]
// 		if ok && r.Status == "reserved" {
// 			r.Status = "confirmed"
// 			r.ExpiresAt = 0
// 			state.reservations[id] = r
// 			state.stateMutex.Unlock()
// 			state.broadcast(models.Event{Event: "horario.confirmado", Data: r})
// 			return
// 		}
// 		state.stateMutex.Unlock()
// 	}(resID)
// }

// func handleCancel(data interface{}) {
// 	m, _ := data.(map[string]interface{})
// 	resID, _ := m["reservationId"].(string)
// 	state.stateMutex.Lock()
// 	r, ok := state.reservations[resID]
// 	if ok && (r.Status == "reserved" || r.Status == "confirmed") {
// 		r.Status = "cancelled"
// 		state.reservations[resID] = r
// 		state.stateMutex.Unlock()
// 		state.broadcast(models.Event{Event: "reserva.cancelada", Data: r})
// 		return
// 	}
// 	state.stateMutex.Unlock()
// }

// func handleStopSystem() {
// 	state.stateMutex.Lock()
// 	state.systemStarted = false
// 	state.stateMutex.Unlock()
// 	state.broadcast(models.Event{Event: "sistema.parado", Data: map[string]interface{}{
// 		"message": "Sistema parado pelo administrador",
// 	}})
// }

// func handleFinalize(data interface{}) {
// 	m, _ := data.(map[string]interface{})
// 	resID, _ := m["reservationId"].(string)
// 	state.stateMutex.Lock()
// 	r, ok := state.reservations[resID]
// 	if ok && r.Status == "confirmed" {
// 		r.Status = "finalized"
// 		state.reservations[resID] = r
// 		state.stateMutex.Unlock()
// 		state.broadcast(models.Event{Event: "jogo.finalizado", Data: r})
// 		return
// 	}
// 	state.stateMutex.Unlock()
// }

// // --- Conexão WebSocket ---
// func WsHandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := upgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		log.Println("upgrade:", err)
// 		return
// 	}
// 	client := &Client{conn: conn, send: make(chan models.Event, 16)}
// 	state.clientsMutex.Lock()
// 	state.clients[client] = true
// 	state.clientsMutex.Unlock()

// 	// Envia snapshot inicial
// 	state.stateMutex.Lock()
// 	courts := make([]models.Court, 0, len(state.courts))
// 	for _, c := range state.courts {
// 		courts = append(courts, c)
// 	}
// 	resList := make([]models.Reservation, 0, len(state.reservations))
// 	for _, r := range state.reservations {
// 		resList = append(resList, r)
// 	}
// 	systemStarted := state.systemStarted
// 	state.stateMutex.Unlock()

// 	client.send <- models.Event{Event: "state.snapshot", Data: map[string]interface{}{
// 		"courts": courts, "reservations": resList, "systemStarted": systemStarted,
// 	}}

// 	go client.writeLoop()
// 	client.readLoop()
// }

func main() {
	http.HandleFunc("/ws", conection.WsHandler)
	addr := ":8080"
	log.Println("Servidor WebSocket rodando em", addr)
	log.Fatal(http.ListenAndServe(addr, nil))
}
