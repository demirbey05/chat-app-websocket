package main

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {

	fmt.Println("Server is starting ...")
	comm := ServerComm{RegisterChannel: make(chan *Client, 20), BroadcastChannel: make(chan Message, 20)}
	storage := ServerStorage{Clients: make(map[*Client]bool)}

	r := gin.Default()
	r.Static("/css", "./static/css")
	r.Static("/js", "./static/js")
	r.GET("/", HomeHandler)
	r.GET("/ws", func(c *gin.Context) {
		wsHandler(c.Writer, c.Request, &comm)
	})
	r.LoadHTMLFiles("static/index.html")
	go Registerer(&comm, &storage)
	go Broadcaster(&comm, &storage)
	r.Run("0.0.0.0:8080")
}

var wsUpgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
