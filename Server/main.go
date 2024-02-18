package main

import (
	messagequeue "MessageQueue"
	"fmt"
	"net/http"
	"golang.org/x/net/websocket"
)

var ChatRooom = NewChatroom(25)

func connect(w http.ResponseWriter, r *http.Request)  {
	fmt.Fprintln(w, "setting up the server")
}

func home(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./templates/static/index.html")
}


func main (){
	fmt.Println("Starting server, listening on port 8000")
	ChatRooom.ReceiveMessage(messagequeue.Message{RawText: "M1"})
	ChatRooom.ReceiveMessage(messagequeue.Message{RawText: "M2"})
	http.HandleFunc("/", home)
	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", 
		http.StripPrefix("/static/", staticServer))
	http.Handle("/ws", websocket.Handler(ChatRooom.connect))
	http.ListenAndServe(":8000", nil)
}
