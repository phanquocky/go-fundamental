package main

import (
	"fmt"
	"os"

	"github.com/gorilla/websocket"
)

var filename = "000000000.fdb"

func main() {
	// This is a client implementation

	conn, _, err := websocket.DefaultDialer.Dial("ws://localhost:8080/ws", nil)
	if err != nil {
		fmt.Println("Error create Dial:", err)
		panic(err)
	}
	defer conn.Close()

	fmt.Println("Client connected, start sync data from file " + filename)
	count := 0

	for {
		_, p, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
				fmt.Println("Close connection, err: ", err)
				return
			}
			fmt.Println("read message from connection fail, err: ", err)
			return
		}

		file, err := os.OpenFile(filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
		if err != nil {
			fmt.Printf("Failed to open file: %s", err)
		}
		defer file.Close()

		_, err = file.Write(p)
		if err != nil {
			fmt.Printf("Failed to write file: %s", err)
		}

		fmt.Printf("Successfully received %d messages\n", count)
		count++
	}

}
