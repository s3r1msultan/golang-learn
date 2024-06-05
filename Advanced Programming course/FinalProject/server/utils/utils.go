package utils

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
	"time"

	"final_project/server/crypto"
	"final_project/server/state"
)

func SendHistory(conn net.Conn) {
	file, err := os.Open(state.GetHistoryLog())
	if err != nil {
		encryptedMessage, _ := crypto.EncryptMessage("Error reading history.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		message := scanner.Text() + "\n"
		encryptedMessage, err := crypto.EncryptMessage(message)
		if err != nil {
			log.Printf("Error encrypting history message: %s", err)
			continue
		}
		conn.Write([]byte(encryptedMessage + "\n"))
		time.Sleep(10 * time.Millisecond)
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading from history file: %s", err)
		encryptedMessage, _ := crypto.EncryptMessage("Error occurred while reading history.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	}
}

func AddTask(conn net.Conn, owner, description string) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	taskID := fmt.Sprintf("%d", state.IncrementTaskIDCounter())
	state.Tasks[taskID] = state.Task{ID: taskID, Description: description, Owner: owner}

	encryptedMessage, _ := crypto.EncryptMessage(fmt.Sprintf("Task added with ID %s\n", taskID))
	conn.Write([]byte(encryptedMessage + "\n"))
}

func ListTasks(conn net.Conn) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	var taskDescriptions []string
	for _, task := range state.Tasks {
		taskDescriptions = append(taskDescriptions, fmt.Sprintf("ID: %s, Owner: %s, Description: %s", task.ID, task.Owner, task.Description))
	}

	if len(taskDescriptions) == 0 {
		encryptedMessage, _ := crypto.EncryptMessage("No tasks found.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	} else {
		for _, description := range taskDescriptions {
			encryptedDescription, _ := crypto.EncryptMessage(description + "\n")
			conn.Write([]byte(encryptedDescription + "\n"))
			time.Sleep(10 * time.Millisecond)
		}
	}
}

func DeleteTask(conn net.Conn, taskID string) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	if _, ok := state.Tasks[taskID]; ok {
		delete(state.Tasks, taskID)
		encryptedMessage, _ := crypto.EncryptMessage("Task deleted successfully.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	} else {
		encryptedMessage, _ := crypto.EncryptMessage("Task not found.\n")
		conn.Write([]byte(encryptedMessage + "\n"))
	}
}

func BanUser(conn net.Conn, nickname string) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	for clientConn, clientNickname := range state.Clients {
		if clientNickname == nickname {
			clientIP := clientConn.RemoteAddr().String()
			state.BannedIPs[clientIP] = true
			encryptedMessage, _ := crypto.EncryptMessage("You have been banned from the server.\n")
			clientConn.Write([]byte(encryptedMessage + "\n"))
			clientConn.Close()
			encryptedMessage, _ = crypto.EncryptMessage(fmt.Sprintf("User %s banned successfully.\n", nickname))
			conn.Write([]byte(encryptedMessage + "\n"))
			return
		}
	}

	encryptedMessage, _ := crypto.EncryptMessage(fmt.Sprintf("User %s not found.\n", nickname))
	conn.Write([]byte(encryptedMessage + "\n"))
}

func KickUser(conn net.Conn, nickname string) {
	state.Mutex.Lock()
	defer state.Mutex.Unlock()

	for clientConn, clientNickname := range state.Clients {
		if clientNickname == nickname {
			encryptedMessage, _ := crypto.EncryptMessage("You have been kicked from the server.\n")
			clientConn.Write([]byte(encryptedMessage + "\n"))
			clientConn.Close()
			encryptedMessage, _ = crypto.EncryptMessage(fmt.Sprintf("User %s kicked successfully.\n", nickname))
			conn.Write([]byte(encryptedMessage + "\n"))
			return
		}
	}

	encryptedMessage, _ := crypto.EncryptMessage(fmt.Sprintf("User %s not found.\n", nickname))
	conn.Write([]byte(encryptedMessage + "\n"))
}

func BroadcastMessage(message string, sender net.Conn) {
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
