package service

import (
	"errors"
	"fmt"
	"testing"

	mocks "github.com/jh-chee/kitewave/rpc-server/mocks/repository"
	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageRepository(t *testing.T) {
	repository := mocks.NewMessageRepository(t)
	service := NewMessageRepository(repository)
	assert.Implements(t, (*MessageService)(nil), service)
	assert.Equal(t, repository, service.(*messageService).messageRepository)
}

func TestMessageService_Send(t *testing.T) {
	mockMsg := &models.Message{
		ChatRoom: "John:Mary",
		Sender:   "John",
		Body:     "Hello",
	}

	tests := []struct {
		name    string
		mockFn  func(messageRepo *mocks.MessageRepository, mockMsg *models.Message)
		wantErr error
	}{
		{
			name: "success",
			mockFn: func(messageRepo *mocks.MessageRepository, mockMsg *models.Message) {
				messageRepo.On("Save", mockMsg).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "save error",
			mockFn: func(messageRepo *mocks.MessageRepository, mockMsg *models.Message) {
				messageRepo.On("Save", mockMsg).Return(errors.New("save error"))
			},
			wantErr: fmt.Errorf("unable to save msg: %w", errors.New("save error")),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageRepo := mocks.NewMessageRepository(t)
			tt.mockFn(messageRepo, mockMsg)

			messageSvc := &messageService{
				messageRepository: messageRepo,
			}

			err := messageSvc.Send(mockMsg)

			assert.Equal(t, tt.wantErr, err)
			messageRepo.AssertExpectations(t)
		})
	}
}

func TestGetChatRoom(t *testing.T) {
	tests := []struct {
		chat           string
		expectedResult string
		expectedError  error
	}{
		{
			chat:           "user1:user2",
			expectedResult: "user1:user2",
			expectedError:  nil,
		},
		{
			chat:           "user2:user1",
			expectedResult: "user1:user2",
			expectedError:  nil,
		},
		{
			chat:           "User2:User1",
			expectedResult: "User1:User2",
			expectedError:  nil,
		},
		{
			chat:           "uSEr1:UsEr2",
			expectedResult: "UsEr2:uSEr1",
			expectedError:  nil,
		},
		{
			chat:           "user1",
			expectedResult: "",
			expectedError:  fmt.Errorf("received invalid chat format user1, expecting user1:user2"),
		},
		{
			chat:           "user1:user2:user3",
			expectedResult: "",
			expectedError:  fmt.Errorf("received invalid chat format user1:user2:user3, expecting user1:user2"),
		},
	}

	for _, test := range tests {
		result, err := sortParticipants(test.chat)

		if result != test.expectedResult {
			t.Errorf("Expected result: %s, but got: %s", test.expectedResult, result)
		}

		if (err == nil && test.expectedError != nil) || (err != nil && err.Error() != test.expectedError.Error()) {
			t.Errorf("Expected error: %v, but got: %v", test.expectedError, err)
		}
	}
}
