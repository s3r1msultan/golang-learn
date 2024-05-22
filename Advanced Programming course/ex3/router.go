package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func setupRouter(openAICall OpenAICaller) *httprouter.Router {
	router := httprouter.New()
	router.ServeFiles("/public/*filepath", http.Dir("public"))
	router.POST("/chat", handleChat(openAICall))
	router.GET("/", homeHandler)
	router.GET("/history", historyHandler)
	return router
}
