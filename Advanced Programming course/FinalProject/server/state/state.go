package state

import (
	"net"
	"sync"
)

var (
	Clients       = make(map[net.Conn]string)
	Addr          = make(map[net.Conn]string)
	Mutex         sync.Mutex
	BannedIPs     = make(map[string]bool)
	Admins        = make(map[string]bool)
	Logs          []string
	TaskIDCounter int
	Tasks         = make(map[string]Task)
	HistoryLog    = "history.log"
)

type Task struct {
	ID          string
	Description string
	Owner       string
}

func AddLogEntry(entry string) {
	Logs = append(Logs, entry)
}

func GetHistoryLog() string {
	return HistoryLog
}

func IncrementTaskIDCounter() int {
	TaskIDCounter++
	return TaskIDCounter
}
