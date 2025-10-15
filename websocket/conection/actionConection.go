package conection

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
	"main.go/models"
)

var state = NewServerState()
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

// Funções responsáveis pela leitura e escrita do websocket
func (c *Client) readLoop() {
	defer func() {
		c.conn.Close()
		state.clientsMutex.Lock()
		delete(state.clients, c)
		state.clientsMutex.Unlock()
	}()
	for {
		_, msg, err := c.conn.ReadMessage()
		if err != nil {
			log.Println("read error:", err)
			return
		}
		var ev models.Evento
		if err := json.Unmarshal(msg, &ev); err != nil {
			log.Println("invalid event JSON:", err)
			continue
		}
		handleEvent(c, ev)
	}
}

func (c *Client) writeLoop() {
	for ev := range c.send {
		b, err := json.Marshal(ev)
		if err != nil {
			continue
		}
		if err := c.conn.WriteMessage(websocket.TextMessage, b); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}

// Funções relacionadas aos eventos
func handleEvent(c *Client, ev models.Evento) {
	switch ev.Event {
	case "sistema.iniciado":
		handleStartSystem()
	case "sistema.parado":
		handleStopSystem()
	case "quadra.adicionada":
		handleAddCourt(ev.Data)
	case "horario.reservado":
		handleReserve(ev.Data)
	case "reserva.cancelada":
		handleCancel(ev.Data)
	case "jogo.finalizado":
		handleFinalize(ev.Data)
	default:
		log.Println("Evento desconhecido:", ev.Event)
	}
}

// Inciar o sistema pra a reserva de quadras
func handleStartSystem() {
	state.stateMutex.Lock()
	state.systemStarted = true
	courts := make([]models.Quadra, 0, len(state.courts))
	for _, ct := range state.courts {
		courts = append(courts, ct)
	}
	state.stateMutex.Unlock()
	state.broadcast(models.Evento{Event: "sistema.iniciado", Data: map[string]interface{}{"courts": courts}})
}

// Parar o sistema e bloquear reserva de quadra
func handleStopSystem() {
	state.stateMutex.Lock()
	state.systemStarted = false
	state.stateMutex.Unlock()
	state.broadcast(models.Evento{Event: "sistema.parado", Data: map[string]interface{}{
		"message": "Sistema parado pelo administrador",
	}})
}

// Adicionar quadras
func handleAddCourt(data interface{}) {
	m, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	name, _ := m["name"].(string)
	cap, _ := m["capacity"].(float64)
	if name == "" {
		return
	}

	state.stateMutex.Lock()
	state.nextCourtID++
	id := fmt.Sprintf("q%d", state.nextCourtID)
	ct := models.Quadra{ID: id, Name: name, Capacity: int(cap)}
	state.courts[id] = ct
	state.stateMutex.Unlock()

	state.broadcast(models.Evento{Event: "quadra.adicionada", Data: ct})
}

// Reservar quadra
func handleReserve(data interface{}) {

	m, ok := data.(map[string]interface{})
	if !ok {
		return
	}
	courtId, _ := m["courtId"].(string)
	user, _ := m["user"].(string)
	startFloat, _ := m["startTime"].(float64)
	start := int64(startFloat)

	duration := 60 * 60 // 1 hora em segundos
	end := start + int64(duration)

	state.stateMutex.Lock()

	//Só libera a reserva caso o sistema tenah sido iniciado pelo administrador
	if !state.systemStarted {
		state.stateMutex.Unlock()
		state.broadcast(models.Evento{Event: "reserva.negada", Data: map[string]interface{}{
			"reason":  "sistema não iniciado",
			"user":    user,
			"courtId": courtId,
		}})
		return
	}

	for _, r := range state.reservations {
		if r.CourtID != courtId {
			continue
		}
		// só verifica conflitos em reservas ativas
		if r.Status != "reserved" && r.Status != "confirmed" {
			continue
		}

		// Verifica confluito de horário
		if start < r.EndTime && end > r.StartTime {
			state.stateMutex.Unlock()
			state.broadcast(models.Evento{
				Event: "reserva.negada",
				Data: map[string]interface{}{
					"reason":    "Conflito de horário",
					"courtId":   courtId,
					"startTime": start,
					"endTime":   end,
					"conflict":  r.ID,
				},
			})
			return
		}
	}

	// Aṕos criar a reserva ela fica em uso por 1 minuto
	state.nextReserveID++
	resID := fmt.Sprintf("r-%d", state.nextReserveID)
	exp := time.Now().Add(1 * time.Minute)
	res := models.Reserva{
		ID:        resID,
		CourtID:   courtId,
		User:      user,
		StartTime: start,
		EndTime:   end,
		ExpiresAt: exp.Unix(),
		Status:    "reserved",
	}
	state.reservations[resID] = res
	state.stateMutex.Unlock()

	state.broadcast(models.Evento{Event: "horario.reservado", Data: res})

	// Expira após 1 minutos
	go func(id string) {
		timer := time.NewTimer(1 * time.Minute)
		<-timer.C
		state.stateMutex.Lock()
		r, ok := state.reservations[id]
		if ok && r.Status == "reserved" {
			r.Status = "expired"
			state.reservations[id] = r
			state.stateMutex.Unlock()
			state.broadcast(models.Evento{Event: "reserva.expirada", Data: r})
			return
		}
		state.stateMutex.Unlock()
	}(resID)

	// Confirma após 10 segundos
	go func(id string) {
		time.Sleep(10 * time.Second)
		state.stateMutex.Lock()
		r, ok := state.reservations[id]
		if ok && r.Status == "reserved" {
			r.Status = "confirmed"
			r.ExpiresAt = 0
			state.reservations[id] = r
			state.stateMutex.Unlock()
			state.broadcast(models.Evento{Event: "horario.confirmado", Data: r})
			return
		}
		state.stateMutex.Unlock()
	}(resID)
}

// Cancelar reserva de quadra
func handleCancel(data interface{}) {
	m, _ := data.(map[string]interface{})
	resID, _ := m["reservationId"].(string)
	state.stateMutex.Lock()
	r, ok := state.reservations[resID]
	if ok && (r.Status == "reserved" || r.Status == "confirmed") {
		r.Status = "cancelled"
		state.reservations[resID] = r
		state.stateMutex.Unlock()
		state.broadcast(models.Evento{Event: "reserva.cancelada", Data: r})
		return
	}
	state.stateMutex.Unlock()
}

func handleFinalize(data interface{}) {
	m, _ := data.(map[string]interface{})
	resID, _ := m["reservationId"].(string)
	state.stateMutex.Lock()
	r, ok := state.reservations[resID]
	if ok && r.Status == "confirmed" {
		r.Status = "finalized"
		state.reservations[resID] = r
		state.stateMutex.Unlock()
		state.broadcast(models.Evento{Event: "jogo.finalizado", Data: r})
		return
	}
	state.stateMutex.Unlock()
}

// --- Conexão WebSocket ---
func WsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	client := &Client{conn: conn, send: make(chan models.Evento, 16)}
	state.clientsMutex.Lock()
	state.clients[client] = true
	state.clientsMutex.Unlock()

	// Envia snapshot inicial
	state.stateMutex.Lock()
	courts := make([]models.Quadra, 0, len(state.courts))
	for _, c := range state.courts {
		courts = append(courts, c)
	}
	resList := make([]models.Reserva, 0, len(state.reservations))
	for _, r := range state.reservations {
		resList = append(resList, r)
	}
	systemStarted := state.systemStarted
	state.stateMutex.Unlock()

	client.send <- models.Evento{Event: "state.snapshot", Data: map[string]interface{}{
		"courts": courts, "reservations": resList, "systemStarted": systemStarted,
	}}

	go client.writeLoop()
	client.readLoop()
}
