package main

import (
	"MessageQueue"
	"bytes"
	"encoding/json"
	"fmt"
	"html/template"
	"io"
	"sync"

	"golang.org/x/net/websocket"
)

type ChatRoom struct {
	connections	map[*websocket.Conn]string
	users		map[string]*websocket.Conn
	messages	messagequeue.MessageQueue
	mutex		sync.Mutex
}

func NewChatroom(msgBufferSize int)ChatRoom{
	return ChatRoom{
		connections: make(map[*websocket.Conn]string),
		users: make(map[string]*websocket.Conn),
		messages: messagequeue.New(msgBufferSize),
	}
}

func (cr *ChatRoom) ReceiveMessage(message messagequeue.Message) {
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	cr.messages.AddMessage(message)
	fmt.Println("Message Received: ", message)

	if (cr.users[message.Username] == nil){
		//TODO: redirect user thats sending a message even though theyve left the chat
		// Maybe correlate this with their connection as well to make double sure
		fmt.Println("UH oh this user shouldnt exist, redirect them to the join page")
	}
	cr.broadcast(message)
}

func (cr *ChatRoom) broadcast(message messagequeue.Message) {
	fmt.Println(" Broadcasting message to all connections")
	for connection := range cr.connections{
		sendMessage(message, connection)
	}

}

func (cr *ChatRoom) connect(conn *websocket.Conn) {
	cr.sendAllMessages(conn)
	cr.readLoop(conn)
}


func (cr * ChatRoom) sendAllMessages(connection *websocket.Conn){
	cr.mutex.Lock()
	defer cr.mutex.Unlock()

	iterator := messagequeue.NewIterator(&cr.messages)
	for iterator.HasNext(){
		message := iterator.Next()
		go sendMessage(message, connection)
	}
}

type InitialMessage struct {
	Username	string
}


func (cr * ChatRoom) readLoop(ws *websocket.Conn){
	buf := make([]byte, 1024)
	var initialMessage InitialMessage
	for {
		n, err:= ws.Read(buf)
		if err != nil {
			if err == io.EOF{
				fmt.Println("Connection closed")
				cr.removeUser(ws, initialMessage.Username)
				break
			}
			fmt.Println("Read err:", err)
			continue
		}
		msg := buf[:n]
		err = json.Unmarshal(msg, &initialMessage)
		if err != nil {
			fmt.Println(err)
		} else {
			cr.adduser(ws, initialMessage.Username)
		}
	}

}

func (cr *ChatRoom) checkUsernameValid(username string) bool{
	return cr.users[username] == nil
}

func (cr *ChatRoom) removeUser(ws *websocket.Conn, username string){
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	delete(cr.connections, ws)
	fmt.Printf("removing user %s, from the chat room\n", username)
	delete(cr.users, username)
}

func (cr *ChatRoom) adduser(ws *websocket.Conn, username string){
	cr.mutex.Lock()
	defer cr.mutex.Unlock()
	fmt.Printf("adding user %s, to the chat room\n", username)
	cr.connections[ws] = username
	cr.users[username] = ws
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
