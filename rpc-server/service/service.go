package service

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/jh-chee/kitewave/rpc-server/repository"
)

type MessageService interface {
	Send(msg *models.Message) error
	Pull(req *models.Request) (resp *models.Response, err error)
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
		return err
	}

	if err := s.messageRepository.Save(msg); err != nil {
		return fmt.Errorf("unable to save msg: %w", err)
	}
	return nil
}

func (s *messageService) Pull(req *models.Request) (resp *models.Response, err error) {
	req.ChatRoom, err = sortParticipants(req.ChatRoom)
	if err != nil {
		return nil, err
	}

	if req.Limit < 0 {
		return nil, fmt.Errorf("invalid limit value, got %v", req.Limit)
	}

	// Cursor should exists in db
	if err := s.messageRepository.CheckCursorExistence(req.Cursor); err != nil {
		return nil, err
	}

	msgs, nextCursor, err := s.messageRepository.Pull(req)
	if err != nil {
		return nil, fmt.Errorf("unable to save msg: %w", err)
	}

	hasMore := nextCursor != 0
	return &models.Response{
		Messages:   msgs,
		HasMore:    &hasMore,
		NextCursor: &nextCursor,
	}, nil
}

// Sort the participants in chat room in ascending order so user2:user1 is equivalent to user1:user2
func sortParticipants(chat string) (string, error) {
	// Chatroom pattern, e.g. John123:Mary555
	regex := regexp.MustCompile(`^[a-zA-Z0-9]+:[a-zA-Z0-9]+$`)
	if ok := regex.MatchString(chat); !ok {
		return "", fmt.Errorf("received invalid chat format %s, expecting user1:user2", chat)
	}

	participants := strings.Split(chat, ":")
	p0, p1 := participants[0], participants[1]
	if p1 < p0 {
		return fmt.Sprintf("%s:%s", p1, p0), nil
	}
	return fmt.Sprintf("%s:%s", p0, p1), nil
}
