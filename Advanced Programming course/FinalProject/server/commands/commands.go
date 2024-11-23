package commands

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"time"

	"final_project/server/config"
	"final_project/server/crypto"
	"final_project/server/state"
	"final_project/server/utils"
)

func HandleCommands(conn net.Conn, nickname *string, message string, logFile *os.File) {
	if strings.HasPrefix(message, "/quit") {
		encryptedMessage, _ := crypto.EncryptMessage("Goodbye!\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		conn.Close()
	} else if strings.HasPrefix(message, "/history") {
		utils.SendHistory(conn)
	} else if strings.HasPrefix(message, "/help") {
		sendHelp(conn)
	} else if strings.HasPrefix(message, "/nickname") {
		parts := strings.SplitN(message, " ", 2)
		if len(parts) == 2 {
			changeNickname(conn, nickname, parts[1])
		}
	} else if strings.HasPrefix(message, "/statistics") {
		sendStatistics(conn)
	} else if strings.HasPrefix(message, "/users") {
		sendUsersList(conn)
	} else if strings.HasPrefix(message, "/bot task") {
		handleBotTaskCommands(conn, nickname, message)
	} else if strings.HasPrefix(message, "/admin") {
		handleAdminCommands(conn, message)
	} else if strings.HasPrefix(message, "/bot timer") {
		parts := strings.SplitN(message, " ", 3)
		if len(parts) == 3 {
			timerDuration, err := time.ParseDuration(parts[2] + "m")
			if err != nil {
				encryptedMessage, _ := crypto.EncryptMessage("Invalid timer duration.\n")
				conn.Write([]byte(encryptedMessage + "\n"))
				return
			}
			setBotTimer(conn, *nickname, timerDuration)
		}
	} else {
		logMessage(*nickname, message, logFile)
		response := fmt.Sprintf("%s: %s\n", *nickname, message)
		broadcastMessage(response, conn)
		state.Mutex.Lock()
		state.MessageCount++
		state.Mutex.Unlock()
	}

	currentTime := time.Now().Format(time.RFC1123)
	logEntry := fmt.Sprintf("%s: %s - %s\n", currentTime, state.Addr[conn], message)
	state.AddLogEntry(logEntry)
}

func sendStatistics(conn net.Conn) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	totalUsers := state.CurrentUsers
	currentUsers := len(state.Clients)
	totalMessages := state.MessageCount
	message := fmt.Sprintf("Current connected users: %d\nTotal connected users: %d\nTotal messages sent: %d\n", totalUsers, currentUsers, totalMessages)
	encryptedMessage, _ := crypto.EncryptMessage(message)
	conn.Write([]byte(encryptedMessage + "\n"))
}

func handleBotTaskCommands(conn net.Conn, nickname *string, message string) {
	parts := strings.SplitN(message, " ", 4)
	if len(parts) < 3 {
		encryptedMessage, _ := crypto.EncryptMessage("Invalid /bot task command.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		return
	}

	action := parts[2]
	switch action {
	case "add":
		if len(parts) == 4 {
			utils.AddTask(conn, *nickname, parts[3])
		} else {
			encryptedMessage, _ := crypto.EncryptMessage("Usage: /bot task add <description>\n")
			conn.Write([]byte(encryptedMessage + "\n"))
		}
	case "list":
		utils.ListTasks(conn)
	case "delete":
		if len(parts) == 4 {
			utils.DeleteTask(conn, parts[3])
		} else {
			encryptedMessage, _ := crypto.EncryptMessage("Usage: /bot task delete <task_id>\n")
			conn.Write([]byte(encryptedMessage + "\n"))
		}
	default:
		encryptedMessage, _ := crypto.EncryptMessage("Unknown /bot task action.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	}
}

func handleAdminCommands(conn net.Conn, message string) {
	parts := strings.SplitN(message, " ", 3)
	clientIP := conn.RemoteAddr().String()

	if len(parts) < 2 {
		encryptedMessage, _ := crypto.EncryptMessage("Invalid admin command.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		return
	}

	if parts[1] == config.ADMIN_PASSWORD {
		state.Admins[clientIP] = true
		encryptedMessage, _ := crypto.EncryptMessage("Admin access granted.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		return
	}

	if !state.Admins[clientIP] {
		encryptedMessage, _ := crypto.EncryptMessage("You are not an admin.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		return
	}

	switch parts[1] {
	case "quit":
		delete(state.Admins, clientIP)
		encryptedMessage, _ := crypto.EncryptMessage("Admin access revoked.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	case "ban":
		if len(parts) == 3 {
			utils.BanUser(conn, parts[2])
		} else {
			encryptedMessage, _ := crypto.EncryptMessage("Usage: /admin ban <nickname>\n")
			conn.Write([]byte(encryptedMessage + "\n"))
		}
	case "kick":
		if len(parts) == 3 {
			utils.KickUser(conn, parts[2])
		} else {
			encryptedMessage, _ := crypto.EncryptMessage("Usage: /admin kick <nickname>\n")
			conn.Write([]byte(encryptedMessage + "\n"))
		}
	case "logs":
		sendLogs(conn)
	default:
		encryptedMessage, _ := crypto.EncryptMessage("Unknown admin command.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	}

	currentTime := time.Now().Format(time.RFC1123)
	logEntry := fmt.Sprintf("%s: %s - %s\n", currentTime, clientIP, message)
	state.AddLogEntry(logEntry)
}

func sendHelp(conn net.Conn) {
	commandsList := []string{
		"/help - Display this help message",
		"/nickname <name> - Change your nickname",
		"/users - List all connected users",
		"/bot task add <description> - Add a new task with the given description",
		"/bot task list - List all current tasks",
		"/bot task delete <task_id> - Delete a task by its ID",
		"/bot timer <minutes> - Set a timer for the specified number of minutes",
		"/admin <password> - Gain admin access using the specified password",
		"/admin ban <nickname> - Ban a user by their nickname",
		"/admin kick <nickname> - Kick a user by their nickname",
		"/admin logs - Display the server logs",
		"/history - Display the chat history",
		"/quit - Disconnect from the server",
	}
	for _, command := range commandsList {
		encryptedCommand, _ := crypto.EncryptMessage(command + "\n")
		conn.Write([]byte(encryptedCommand + "\n"))
		time.Sleep(10 * time.Millisecond)
	}
}

func sendUsersList(conn net.Conn) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	var users []string
	for _, nickname := range state.Clients {
		users = append(users, nickname)
	}

	usersList := strings.Join(users, ", ")
	message := fmt.Sprintf("Connected users: %s\n", usersList)
	encryptedMessage, _ := crypto.EncryptMessage(message)
	conn.Write([]byte(encryptedMessage + "\n"))
}

func changeNickname(conn net.Conn, nickname *string, newNickname string) {
	oldNickname := *nickname
	*nickname = newNickname

	state.Mutex.Lock()
	state.Clients[conn] = newNickname
	state.Mutex.Unlock()

	encryptedMessage, _ := crypto.EncryptMessage(fmt.Sprintf("Nickname changed to %s\n", newNickname))
	conn.Write([]byte(encryptedMessage + "\n"))

	log.Printf("Client %s (%s) changed nickname to %s.", state.Addr[conn], oldNickname, newNickname)
	state.AddLogEntry(fmt.Sprintf("Client %s (%s) changed nickname to %s.", state.Addr[conn], oldNickname, newNickname))
	utils.BroadcastMessage(fmt.Sprintf("'%s' changed nickname to '%s'\n", oldNickname, newNickname), conn)
}

func sendLogs(conn net.Conn) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	for _, logEntry := range state.Logs {
		encryptedLogEntry, _ := crypto.EncryptMessage(logEntry + "\n")
		conn.Write([]byte(encryptedLogEntry + "\n"))
		time.Sleep(10 * time.Millisecond)
	}
}

func logMessage(nickname string, message string, logFile *os.File) {
	currentTime := time.Now().Format(time.RFC1123)
	logEntry := fmt.Sprintf("%s: %s - %s\n", currentTime, nickname, message)
	logFile.WriteString(logEntry)
}

func broadcastMessage(message string, sender net.Conn) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()
	for conn := range state.Clients {
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

func setBotTimer(conn net.Conn, owner string, duration time.Duration) {
	encryptedMessage, _ := crypto.EncryptMessage(fmt.Sprintf("Timer set for %v minutes.\n", duration.Minutes()))
	conn.Write([]byte(encryptedMessage + "\n"))

	go func() {
		time.Sleep(duration)
		message := fmt.Sprintf("Timer set by %s has ended.\n", owner)
		utils.BroadcastMessage(message, nil)
	}()
}
