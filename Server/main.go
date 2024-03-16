package main

import (
	messagequeue "MessageQueue"
	"fmt"
	"html/template"
	"net/http"

	"golang.org/x/net/websocket"
)

var ChatRooom = NewChatroom(25)


func join(w http.ResponseWriter, r *http.Request){
	err := r.ParseForm()
	username := r.Form.Get("username")

	var allErrors []string

	if err != nil {
		allErrors = append(allErrors, "server error")
	}

	if (!ChatRooom.checkUsernameValid(username)){
		fmt.Println("invalid username")
		allErrors = append(allErrors, "Username taken")
	}

	if ( len(allErrors) > 0){
		tmpl, err := template.ParseFiles("./templates/joinError.tmpl")
		if ( err != nil ) {
			fmt.Println(err)
		}
		err=tmpl.Execute(w, allErrors)
		if ( err != nil){
			fmt.Println(err)
		}
		return
	}
	tmpl, err := template.ParseFiles("./templates/join.tmpl")
	tmpl.Execute(w,username)
}

func message(w http.ResponseWriter, r *http.Request)  {
	err := r.ParseForm()
	if err != nil {
		fmt.Println("error parsing form")
	}
	formValues := r.Form
	rawText := formValues.Get("rawText")
	username := formValues.Get("username")
	message := messagequeue.NewMessage(rawText, username)
	ChatRooom.ReceiveMessage(message)
}

func home(w http.ResponseWriter, r *http.Request){
	http.ServeFile(w, r, "./templates/static/index.html")
}


func main (){
	fmt.Println("Starting server, listening on port 8000")
	
	ChatRooom.messages.AddMessage(messagequeue.NewMessage("hello", "alex"))
	ChatRooom.messages.AddMessage(messagequeue.NewMessage("hello2", "alex"))
	http.HandleFunc("/", home)
	staticServer := http.FileServer(http.Dir("./static"))
	http.Handle("/static/", 
		http.StripPrefix("/static/", staticServer))
	http.Handle("/ws", websocket.Handler(ChatRooom.connect))
	http.HandleFunc("/message", message) 
	http.HandleFunc("/join", join) 
	http.ListenAndServe(":8000", nil)
}
