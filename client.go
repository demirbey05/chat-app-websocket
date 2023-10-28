package main

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"log"
)

type ServerComm struct {
	RegisterChannel  chan *Client
	BroadcastChannel chan Message
}

type ServerStorage struct {
	Clients  map[*Client]bool
	Messages []Message
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

type Client struct {
	Username string
	Conn     *websocket.Conn
}
type Response struct {
	Type    string `json:"type"`
	Content string `json:"content"`
}

func (c *Client) ReadPump(comm *ServerComm, storage *ServerStorage) {
	var objmap map[string]json.RawMessage
	var unmarshalledMessage Message
	defer func() {
		c.Conn.Close()
	}()

	for {
		_, message, err := c.Conn.ReadMessage()
		json.Unmarshal(message, &objmap)
		json.Unmarshal(objmap["content"], &unmarshalledMessage)
		if err != nil {
			log.Println("Error reading message: ", err)
			break
		}
		messageObject := Message{Username: c.Username, Message: unmarshalledMessage.Message}
		log.Println("Message received: ", messageObject)
		comm.BroadcastChannel <- messageObject
		storage.Messages = append(storage.Messages, messageObject)

	}
}

func Broadcaster(comm *ServerComm, storage *ServerStorage) {
	for message := range comm.BroadcastChannel {
		m, _ := json.Marshal(message)
		responseObject := Response{Type: "message", Content: string(m)}
		message, err := json.Marshal(responseObject)
		if err != nil {
			log.Println("Error marshalling message: ", err)
			continue
		}
		log.Println(storage.Clients)
		for client := range storage.Clients {
			err := client.Conn.WriteMessage(websocket.TextMessage, message)
			log.Println("Message sent: ", "to", client.Username)
			if err != nil {
				log.Println("Error writing message: ", err)
				client.Conn.Close()
				delete(storage.Clients, client)
			}
		}
	}
}

func Registerer(comm *ServerComm, storage *ServerStorage) {
	for client := range comm.RegisterChannel {
		storage.Clients[client] = true
		log.Println("Client registered: ", client.Username)
		client.Conn.WriteJSON(Response{Type: "user_registration", Content: client.Username})
		go client.ReadPump(comm, storage)
	}
}
