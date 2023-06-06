package models

type Request struct {
	ChatRoom string
	Cursor   int64
	Limit    int32
	Reverse  bool
}
