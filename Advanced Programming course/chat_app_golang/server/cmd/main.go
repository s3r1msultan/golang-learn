package main

import (
	"server/db"
	"server/initializers"
	"server/internal/user"
	"server/internal/websocket"
	"server/router"
)

func main() {
	err := db.Connect()
	if err != nil {
		initializers.LogError("connection to db", err, nil)
	}

	database := db.MongoClient.Database(db.GetDB())

	userRep := user.NewRepository(database)
	userService := user.NewService(userRep)
	userHandler := user.NewHandler(userService)

	hub := websocket.NewHub()
	websocketHandler := websocket.NewHandler(hub)
	go hub.Run()

	router.InitRouter(userHandler, websocketHandler)
	router.Start(initializers.GetPort())
}

func init() {
	initializers.InitLogger()
	initializers.InitDotEnv()
}
