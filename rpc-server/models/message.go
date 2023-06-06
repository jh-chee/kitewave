package models

type Message struct {
	Id       int    `json:"id"`
	ChatRoom string `json:"chatroom"`
	Sender   string `json:"sender"`
	Body     string `json:"body"`
	Created  int64  `json:"created"`
}
