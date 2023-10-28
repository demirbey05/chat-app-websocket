package main

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
)

func HomeHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.html", nil)
}

func wsHandler(w http.ResponseWriter, r *http.Request, comm *ServerComm) {
	conn, err := wsUpgrader.Upgrade(w, r, nil)
	var objmap map[string]json.RawMessage
	fmt.Println("Connected to server")
	if err != nil {
		return
	}
	_, message, err := conn.ReadMessage()
	if err != nil {
		log.Println("Error reading message: ", err)
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection: ", err)
		}
	}
	err = json.Unmarshal(message, &objmap)
	if err != nil {
		log.Println("Error unmarshalling message: ", err)
	}
	reqType := strings.Trim(string(objmap["type"]), "\"")
	if reqType != "username" {
		log.Println("Error: First message must be username")
		err := conn.Close()
		if err != nil {
			log.Println("Error closing connection: ", err)
		}
	}

	content := objmap["content"]
	comm.RegisterChannel <- &Client{Username: string(content), Conn: conn}
}
