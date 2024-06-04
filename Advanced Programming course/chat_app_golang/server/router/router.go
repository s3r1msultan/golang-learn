package router

import (
	"github.com/gin-gonic/gin"
	"server/internal/user"
	"server/internal/websocket"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, webSocketHandler *websocket.Handler) *gin.Engine {
	r = gin.Default()
	r.POST("/signup", userHandler.CreateUser)
	r.POST("/login", userHandler.Login)
	r.GET("/logout", userHandler.Logout)

	r.POST("/ws/createRoom", webSocketHandler.CreateRoom)
	r.GET("/ws/joinRoom/:roomId", webSocketHandler.JoinRoom)
	r.GET("/ws/getRooms", webSocketHandler.GetRooms)
	r.GET("/ws/getClients/:roomId", webSocketHandler.GetClients)

	return r
}

func Start(addr string) error {
	return r.Run(addr)
}
