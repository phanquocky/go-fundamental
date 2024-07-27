package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/websocket"
)

// upgrader is used to upgrade the HTTP connection to a WebSocket connection
var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var filename = "000000000.fdb"

func handlerWS(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil) // upgrade the HTTP connection to a WebSocket connection
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close()

	fmt.Println("Client connected, start sync data from file " + filename)

	for {
		messageType := websocket.TextMessage // messagetype use to specify the type of message when sending data to the client

		file, err := os.Open(filename) // open the file
		if err != nil {
			log.Fatalf("Failed to open file: %s", err)
		}
		defer file.Close()

		// Create a new buffered reader to read the file
		reader := bufio.NewReader(file)
		count := 0

		for {
			data, err := reader.ReadBytes('\n') // read the file line by line
			if err != nil {
				if err.Error() == "EOF" {
					log.Println("End of file, close connection")
					closeMsg := websocket.FormatCloseMessage(websocket.CloseNormalClosure, "Server is closing the connection")
					err = conn.WriteMessage(websocket.CloseMessage, closeMsg) // send a close message to the client
					return
				}
				log.Fatalf("Failed to read file: %s", err)
			}

			err = conn.WriteMessage(messageType, data) // send the data to the client, with messageType as TextMessage
			if err != nil {
				fmt.Println("len: ", len(data))
				log.Println("Error write message:", err)
				return
			}

			fmt.Printf("Successfully sent %d messages\n", count)
			count++
		}
	}
}

func main() {
	http.HandleFunc("/ws", handlerWS) // handle the WebSocket connection, and upgrade the HTTP connection to a WebSocket connection
	fmt.Println("Server started on localhost:8080")
	if err := http.ListenAndServe(":8080", nil); err != nil { // start the server
		log.Fatal(err)
	}
}
