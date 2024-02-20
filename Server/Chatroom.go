package main

import (
	"MessageQueue"
	"bytes"
	"fmt"
	"html/template"
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
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.messages.AddMessage(message)
	cr.broadcast(message)
}

func (cr *ChatRoom) broadcast(message messagequeue.Message) {
	fmt.Println(" Broadcasting message to all connections")
	for connection := range cr.connections{
		sendMessage(message, connection)
	}

}

func (cr *ChatRoom) connect(conn *websocket.Conn) {
	cr.mutex.Lock()
	cr.connections[conn] = true
	cr.mutex.Unlock()
	cr.sendAllMessages(conn)
	cr.readLoop(conn)
}


func (cr * ChatRoom) sendAllMessages(connection *websocket.Conn){
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	iterator := messagequeue.NewIterator(&cr.messages)
	for iterator.HasNext(){
		message := iterator.Next()
		sendMessage(message, connection)
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
				delete(cr.connections, ws)
				cr.mutex.Unlock()
				break
			}
			fmt.Println("Read err:", err)
			continue
		}
		msg := buf[:n]
		fmt.Println("Message was: ", string(msg))
	}

}

func sendMessage(message messagequeue.Message, conn *websocket.Conn){
	var response bytes.Buffer
	tmpl, err := template.ParseFiles("./templates/message.tmpl")
	if err != nil {
		fmt.Println("Error parsing template")
	}

	err = tmpl.Execute(&response, message)
	if err != nil {
		fmt.Println("error executing message")
	}

	_, err = conn.Write(response.Bytes())


	if err != nil {
		println("Couldnt send the message")
	}
}
