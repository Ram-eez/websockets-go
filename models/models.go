package models

import (
	"bytes"
	"fmt"
	"text/template"
)

type User struct {
	Username string `json:"username"`
	ID       string `json:"id"`
	Password string `json:"password"`
}

type Message struct {
	Username string `json:"username"`
	Message  string `json:"message"`
}

var RegisteredUsers []User
var SecretKey = []byte("secret-key")

func (msg *Message) GetMessageHTML() []byte {
	tmpl, err := template.ParseFiles("views/message.html")
	if err != nil {
		fmt.Println("templete parsing err: ", err)
		return nil
	}

	var renderedMessage bytes.Buffer

	if err := tmpl.Execute(&renderedMessage, msg); err != nil {
		fmt.Println("execution err could not replace : ", err)
		return nil
	}

	fmt.Println("generated HTML with replaced obj: ", renderedMessage.String())

	return renderedMessage.Bytes()
}
