package websocket

type Room struct {
	ID      string             `json:"id"`
	Name    string             `json:"name"`
	Clients map[string]*Client `json:"clients"`
}
type Hub struct {
	Rooms      map[string]*Room
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan *Message
}

func NewHub() *Hub {
	return &Hub{
		Rooms:      make(map[string]*Room),
		Register:   make(chan *Client),
		Unregister: make(chan *Client),
		Broadcast:  make(chan *Message, 5),
	}
}

func (h *Hub) Run() {
	for {
		select {
		case client := <-h.Register:
			if _, ok := h.Rooms[client.RoomID]; ok {
				r := h.Rooms[client.RoomID]
				if _, ok := r.Clients[client.ID]; !ok {
					if len(r.Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  client.Username + " has joined the room",
							RoomID:   client.RoomID,
							Username: client.Username,
						}
					}
					r.Clients[client.ID] = client
				}
			}
		case client := <-h.Unregister:
			if _, ok := h.Rooms[client.RoomID]; ok {
				r := h.Rooms[client.RoomID]
				if _, ok := r.Clients[client.ID]; ok {
					if len(r.Clients) != 0 {
						h.Broadcast <- &Message{
							Content:  client.Username + " has left the room",
							RoomID:   client.RoomID,
							Username: client.Username,
						}
					}
					delete(r.Clients, client.ID)
					close(client.Message)
				}
			}
		case m := <-h.Broadcast:
			if _, ok := h.Rooms[m.RoomID]; ok {
				r := h.Rooms[m.RoomID]
				for _, c := range r.Clients {
					c.Message <- m
				}
			}

		}
	}
}
