package models

import "time"

type Message struct {
	Id       string    `json:"id"`
	ChatRoom string    `json:"chatroom"`
	Sender   string    `json:"sender"`
	Body     string    `json:"body"`
	Created  time.Time `json:"created"`
}
