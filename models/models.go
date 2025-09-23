package models

import (
	"bytes"
	"fmt"
	"os"
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
	RoomID   string `json:"roomid"`
}

var RegisteredUsers []User

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

func GetJWTSecret() []byte {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		secret = "secret-key"
	}
	return []byte(secret)
}
