package service

import (
	"fmt"
	"strings"

	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/jh-chee/kitewave/rpc-server/repository"
)

type MessageService interface {
	Send(msg *models.Message) error
}

type messageService struct {
	messageRepository repository.MessageRepository
}

func NewMessageRepository(messageRepository repository.MessageRepository) MessageService {
	return &messageService{
		messageRepository: messageRepository,
	}
}

func (s *messageService) Send(msg *models.Message) (err error) {
	msg.ChatRoom, err = sortParticipants(msg.ChatRoom)
	if err != nil {
		return fmt.Errorf("unable to save msg: %w", err)
	}

	if err := s.messageRepository.Save(msg); err != nil {
		return fmt.Errorf("unable to save msg: %w", err)
	}
	return nil
}

// Sort the participants in chat room in ascending order so user2:user1 is equivalent to user1:user2
func sortParticipants(chat string) (string, error) {
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
