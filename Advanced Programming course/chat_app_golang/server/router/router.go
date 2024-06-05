package router

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"server/internal/user"
	"server/internal/websocket"
	"time"
)

var r *gin.Engine

func InitRouter(userHandler *user.Handler, webSocketHandler *websocket.Handler) *gin.Engine {
	r = gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000", "http://localhost:3001"},
		AllowMethods:     []string{"GET", "POST"},
		AllowHeaders:     []string{"Content-Type"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "http://localhost:3000"
		},
		MaxAge: 12 * time.Hour,
	}))

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
