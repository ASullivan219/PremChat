package main

import (
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
	http.HandleFunc("/", home)
	http.Handle("/ws", websocket.Handler(ChatRooom.connect))
	http.ListenAndServe(":8000", nil)
}
