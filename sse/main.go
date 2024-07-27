package main

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
)

var clients = make(map[chan string]struct{})

// broadcast sends an event to all connected clients
func broadcast(data string) {
	for client := range clients {
		client <- data
	}
}

func main() {
	router := gin.Default()
	// Clients is a list of channels to send events to connected clients

	router.GET("/sse", func(c *gin.Context) {
		c.Header("Content-Type", "text/event-stream")
		c.Header("Cache-Control", "no-cache")
		c.Header("Connection", "keep-alive")

		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		println("Client connected")
		eventChan := make(chan string)

		clients[eventChan] = struct{}{} // Add the client to the clients map
		defer func() {
			delete(clients, eventChan) // Remove the client when they disconnect
			close(eventChan)
		}()

		notify := c.Writer.CloseNotify()
		go func() {
			<-notify
			fmt.Println("Client disconnected")
		}()

		// Continuously send data to the client
		for {
			data := <-eventChan
			println("Sending data to client", data)
			fmt.Fprintf(c.Writer, "data: %s\n\n", data)
			c.Writer.Flush()
		}

	})

	router.POST("/send-data", func(c *gin.Context) {
		data := c.PostForm("data")
		// print data to console
		println("Data received from client :", data)
		broadcast(data)
		c.JSON(http.StatusOK, gin.H{"message": "Data sent to clients"})
	})

	// Start the server
	err := router.Run(":3000")
	if err != nil {
		fmt.Println(err)
	}

}
