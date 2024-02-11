package main

import (
	"MessageQueue"
	"fmt"
	"sync"

	"golang.org/x/net/websocket"
)

type ChatRoom struct {
	connections	map[*websocket.Conn]bool
	messages	messagequeue.MessageQueue
	mutex		sync.Mutex
}

func NewChatroom(msgBufferSize int)ChatRoom{
	return ChatRoom{
		connections: make(map[*websocket.Conn]bool),
		messages: messagequeue.New(msgBufferSize),
	}
}

func (cr *ChatRoom) broadcast() {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	fmt.Printf(" Broadcasting to all connections")
}

func (cr *ChatRoom) connect(conn *websocket.Conn) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	sendMessage([]byte("Connected!"), conn)
	fmt.Printf(" received a new connection adding it to list")
	cr.connections[conn] = true
}

func sendMessage(message []byte, conn *websocket.Conn){
	_, err := conn.Write(message)
	if err != nil {
		println("Couldnt send the message for some reason or another")
	}
}
