package websocket

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

type Handler struct {
	hub *Hub
}

func NewHandler(hub *Hub) *Handler {
	return &Handler{
		hub: hub,
	}
}

type CreateRoomReq struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

func (h *Handler) CreateRoom(c *gin.Context) {

	var req CreateRoomReq
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	h.hub.Rooms[req.ID] = &Room{
		ID:      req.ID,
		Name:    req.Name,
		Clients: make(map[string]*Client),
	}

	c.JSON(200, req)
}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		//origin := r.Header.Get("Origin")
		//return origin == "http://localhost:3000"
		return true
	},
}

func (h *Handler) JoinRoom(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)

	if err != nil {
		c.JSON(500, gin.H{"error": err.Error()})
		return
	}

	roomID := c.Param("roomId")
	clientID := c.Query("userId")
	username := c.Query("username")
	client := &Client{
		ID:       clientID,
		Username: username,
		RoomID:   roomID,
		Conn:     conn,
		Message:  make(chan *Message),
	}
	h.hub.Register <- client
	go client.writeMessage()
	client.readMessage(h.hub)
	c.JSON(200, h.hub.Rooms[roomID])

}

type RoomRes struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	AmountOfUsers int    `json:"amount_of_users"`
}

func (h *Handler) GetRooms(c *gin.Context) {
	rooms := make([]RoomRes, 0)
	for _, r := range h.hub.Rooms {
		rooms = append(rooms, RoomRes{
			ID:            r.ID,
			Name:          r.Name,
			AmountOfUsers: len(r.Clients),
		})
	}
	c.JSON(200, rooms)
}

type ClientRes struct {
	ID       string `json:"id"`
	Username string `json:"username"`
}

func (h *Handler) GetClients(c *gin.Context) {
	roomId := c.Param("roomId")
	clients := make([]ClientRes, 0)

	if _, ok := h.hub.Rooms[roomId]; ok {
		for _, client := range h.hub.Rooms[roomId].Clients {
			clients = append(clients, ClientRes{
				ID:       client.ID,
				Username: client.Username,
			})
		}
	}
	c.JSON(200, clients)
}
