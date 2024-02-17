package main

import (
	"MessageQueue"
	"fmt"
	"io"
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

func (cr *ChatRoom) ReceiveMessage(message messagequeue.Message) {
	cr.messages.AddMessage(message)
}

func (cr *ChatRoom) broadcast() {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	fmt.Printf(" Broadcasting to all connections")
}

func (cr *ChatRoom) connect(conn *websocket.Conn) {
	cr.mutex.Lock()
	fmt.Println("New connection:")
	fmt.Println(conn.Config().Origin)
	fmt.Println(conn.Config().Location)
	fmt.Println(conn.Config().Dialer)
	sendMessage([]byte("<p> Connected!</p>"), conn)
	cr.connections[conn] = true
	fmt.Println(cr.connections)
	cr.mutex.Unlock()
	
	cr.sendAllMessages(conn)
	cr.readLoop(conn)
}


func (cr * ChatRoom) sendAllMessages(ws *websocket.Conn){
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	iterator := messagequeue.NewIterator(&cr.messages)
	for iterator.HasNext(){
		message := iterator.Next()
		fmt.Println(message)
		ws.Write([]byte("<h1 id=\"ws\" hx-swap-oob=\"outerHTML\"> message </h1><br>"))
	}
}


func (cr * ChatRoom) readLoop( ws *websocket.Conn){
	buf := make([]byte, 1024)
	for {
		n, err:= ws.Read(buf)
		if err != nil {
			if err == io.EOF{
				fmt.Println("Connection closed")
				cr.mutex.Lock()
				defer cr.mutex.Unlock()
				delete(cr.connections, ws)
				break
			}
			fmt.Println("Read err:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println("Message was: ", string(msg))
		ws.Write([]byte("<h1 id=\"ws\" hx-swap-oob=\"outerHTML\"> Message received </h1><br>"))
	}

}

func sendMessage(message []byte, conn *websocket.Conn){
	_, err := conn.Write(message)
	if err != nil {
		println("Couldnt send the message for some reason or another")
	}
}
