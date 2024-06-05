package network

import (
	"bufio"
	"final_project/server/commands"
	"final_project/server/crypto"
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
)

var (
	clients    = make(map[net.Conn]string)
	addr       = make(map[net.Conn]string)
	mutex      sync.Mutex
	bannedIPs  = make(map[string]bool)
	admins     = make(map[string]bool)
	historyLog = "history.log"
	logs       []string
)

func HandleConnection(conn net.Conn) {
	nickname := "Anonymous"
	clientIP := conn.RemoteAddr().String()

	if bannedIPs[clientIP] {
		encryptedMessage, _ := crypto.EncryptMessage("You are banned from this server.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		conn.Close()
		return
	}

	mutex.Lock()
	clients[conn] = nickname
	addr[conn] = clientIP
	mutex.Unlock()

	defer func() {
		mutex.Lock()
		delete(clients, conn)
		delete(addr, conn)
		mutex.Unlock()
		conn.Close()
	}()

	logFile, err := os.OpenFile(historyLog, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Println("Error opening history log file:", err)
		return
	}
	defer logFile.Close()

	log.Printf("Client %s (%s) connected.", addr[conn], nickname)
	logs = append(logs, fmt.Sprintf("Client %s (%s) connected.", addr[conn], nickname))

	for {
		netData, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Printf("Client %s (%s) disconnected.", addr[conn], nickname)
			logs = append(logs, fmt.Sprintf("Client %s (%s) disconnected.", addr[conn], nickname))
			BroadcastMessage(fmt.Sprintf("%s disconnected from the chat!\n", nickname), conn)
			break
		}

		decryptedData, err := crypto.DecryptMessage(strings.TrimSpace(netData))
		if err != nil {
			encryptedMessage, _ := crypto.EncryptMessage("Error decrypting message.\n")
			conn.Write([]byte(encryptedMessage + "\n"))
			continue
		}

		commands.HandleCommands(conn, &nickname, decryptedData, logFile)
	}
}

func BroadcastMessage(message string, sender net.Conn) {
	mutex.Lock()
	defer mutex.Unlock()
	for conn := range clients {
		if conn != sender {
			encryptedMessage, err := crypto.EncryptMessage(message)
			if err != nil {
				fmt.Println("Error encrypting message:", err)
				continue
			}
			conn.Write([]byte(encryptedMessage + "\n"))
		}
	}
}
