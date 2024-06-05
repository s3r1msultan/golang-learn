package main

import (
	"final_project/client/config"
	"final_project/client/crypto"
	"final_project/client/network"
	"final_project/client/utils"
	"fmt"
	"net"
	"os"
	"strings"
	"sync"
)

func main() {
	conn, err := net.Dial(config.CONN_TYPE, config.CONN_PORT)
	if err != nil {
		fmt.Println("Error connecting to server:", err)
		os.Exit(1)
	}

	fmt.Print("Enter your nickname: ")
	nickname := utils.ReadInput()
	nickname = strings.TrimSpace(nickname)
	encryptedNickname, err := crypto.EncryptMessage("/nickname " + nickname)
	if err != nil {
		fmt.Println("Error encrypting nickname:", err)
		os.Exit(1)
	}
	conn.Write([]byte(encryptedNickname + "\n"))

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		network.ReadFromServer(conn)
	}()

	network.WriteToServer(conn)

	wg.Wait()
}
