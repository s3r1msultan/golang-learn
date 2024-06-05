package network

import (
	"bufio"
	"final_project/client/crypto"
	"fmt"
	"net"
	"os"
	"strings"
)

func ReadFromServer(conn net.Conn) {
	for {
		message, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			fmt.Println("Disconnected from the server.")
			return
		}
		decryptedMessage, err := crypto.DecryptMessage(strings.TrimSpace(message))
		if err != nil {
			fmt.Println("Error decrypting message:", err)
			continue
		}
		fmt.Print("Server: ", decryptedMessage)
	}
}

func WriteToServer(conn net.Conn) {
	consoleReader := bufio.NewReader(os.Stdin)
	for {
		input, _ := consoleReader.ReadString('\n')
		trimmedInput := strings.TrimSpace(input)
		if trimmedInput == "/quit" {
			fmt.Println("Disconnecting from server...")
			conn.Write([]byte("/quit\n"))
			conn.Close()
			os.Exit(0)
		}
		encryptedInput, err := crypto.EncryptMessage(trimmedInput)
		if err != nil {
			fmt.Println("Error encrypting message:", err)
			continue
		}
		conn.Write([]byte(encryptedInput + "\n"))
	}
}
