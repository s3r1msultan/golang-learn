package main

import (
	"final_project/server/config"
	"final_project/server/network"
	"log"
	"net"
)

func main() {
	listener, err := net.Listen(config.CONN_TYPE, config.CONN_PORT)
	if err != nil {
		log.Fatal("Error starting TCP server:", err)
	}
	defer listener.Close()
	log.Println("Server listening on", config.CONN_PORT)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("Error accepting connection:", err)
			continue
		}
		go network.HandleConnection(conn)
	}
}
