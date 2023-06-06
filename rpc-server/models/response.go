package models

type Response struct {
	Messages   []*Message
	HasMore    *bool
	NextCursor *int64
}
