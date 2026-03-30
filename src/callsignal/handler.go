package callsignal

import (
	"log"
	"ws/src/chat"

	"github.com/gin-gonic/gin"
)

func ServeSignalingWS(c *gin.Context) {

	userID := c.Query("user")

	conn, err := chat.UPGRADER.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket Upgrade error:", err)
		return
	}

	client := &chat.Client{
		Conn:   conn,
		UserID: userID,
		Send:   make(chan []byte),
	}

	//Register client to signallingclients
	signalingClients[userID] = client

	go client.WritePump()

	//read message from client
	go SignalMessages(client)
}
