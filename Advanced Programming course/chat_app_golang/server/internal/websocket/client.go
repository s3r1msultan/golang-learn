package websocket

import (
	"github.com/gorilla/websocket"
	"server/initializers"
)

type Client struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	RoomID   string `json:"roomId"`
	Conn     *websocket.Conn
	Message  chan *Message
}

type Message struct {
	Content  string `json:"content"`
	Username string `json:"username"`
	RoomID   string `json:"roomId"`
}

func (c *Client) writeMessage() {
	defer func() {
		c.Conn.Close()
	}()

	for {
		message, ok := <-c.Message
		if !ok {
			return
		}
		c.Conn.WriteJSON(message)
	}
}

func (c *Client) readMessage(hub *Hub) {
	defer func() {
		hub.Unregister <- c
		c.Conn.Close()
	}()
	for {
		_, m, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				initializers.LogError("read message", err, nil)
			}
			break
		}
		hub.Broadcast <- &Message{
			Content:  string(m),
			Username: c.Username,
			RoomID:   c.RoomID,
		}

	}
}
