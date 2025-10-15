package models

type Evento struct {
	Event string      `json:"event"`
	Data  interface{} `json:"data"`
}

type Quadra struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Capacity int    `json:"capacity"`
}

type Reserva struct {
	ID        string `json:"id"`
	CourtID   string `json:"courtId"`
	User      string `json:"user"`
	StartTime int64  `json:"startTime"`
	EndTime   int64  `json:"endTime"`
	ExpiresAt int64  `json:"expiresAt"`
	Status    string `json:"status"`
}
