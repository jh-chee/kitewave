package service

import (
	"fmt"
	"strings"

	"github.com/jh-chee/kitewave/rpc-server/kitex_gen/rpc"
	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/jh-chee/kitewave/rpc-server/repository"
)

type MessageService interface {
	Send(req *rpc.SendRequest) error
}

type messageService struct {
	messageRepository repository.MessageRepository
}

func NewMessageRepository(messageRepository repository.MessageRepository) MessageService {
	return &messageService{
		messageRepository: messageRepository,
	}
}

func (s *messageService) Send(req *rpc.SendRequest) error {
	chatroom, err := getChatRoom(req.Message.Chat)
	if err != nil {
		return fmt.Errorf("unable to save msg: %w", err)
	}

	msg := &models.Message{
		Sender:   req.Message.Sender,
		ChatRoom: chatroom,
		Body:     req.Message.Text,
	}

	if err := s.messageRepository.Save(msg); err != nil {
		return fmt.Errorf("unable to save msg: %w", err)
	}
	return nil
}

// Sort the participants in chat room in ascending order so user1:user2 is equivalent to user2:user1
func getChatRoom(chat string) (string, error) {
	participants := strings.Split(chat, ":")
	if len(participants) != 2 {
		return "", fmt.Errorf("received invalid chat format %s, expecting user1:user2", chat)
	}

	p0, p1 := participants[0], participants[1]
	if p1 < p0 {
		return fmt.Sprintf("%s:%s", p1, p0), nil
	}
	return fmt.Sprintf("%s:%s", p0, p1), nil
}
