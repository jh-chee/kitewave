package handler

import (
	"context"
	"fmt"
	"net/http"
	"testing"

	"github.com/jh-chee/kitewave/rpc-server/kitex_gen/rpc"
	mocks "github.com/jh-chee/kitewave/rpc-server/mocks/service"
	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/stretchr/testify/assert"
)

func TestNewMessageHandler(t *testing.T) {
	messageService := mocks.NewMessageService(t)
	handler := NewMessageHandler(messageService)
	assert.Implements(t, (*MessageHandler)(nil), handler)
	assert.Equal(t, messageService, handler.(*messageHandler).messageService)
}

func TestSend(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.SendRequest
	}
	tests := []struct {
		name    string
		args    args
		mockFn  func(messageService *mocks.MessageService, mockMsg *models.Message)
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &rpc.SendRequest{
					Message: &rpc.Message{
						Chat:     "John:Mary",
						Text:     "Hello",
						Sender:   "John",
						SendTime: 0,
					},
				},
			},
			mockFn: func(messageService *mocks.MessageService, mockMsg *models.Message) {
				messageService.On("Send", mockMsg).Return(nil)
			},
			wantErr: nil,
		},
		{
			name: "internal server error",
			args: args{
				ctx: context.Background(),
				req: &rpc.SendRequest{
					Message: &rpc.Message{
						Chat:     "John:Mary",
						Text:     "Hello",
						Sender:   "John",
						SendTime: 0,
					},
				},
			},
			mockFn: func(messageService *mocks.MessageService, mockMsg *models.Message) {
				messageService.On("Send", mockMsg).Return(fmt.Errorf("internal server error"))
			},
			wantErr: fmt.Errorf("internal server error"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageService := mocks.NewMessageService(t)
			mockMsg := &models.Message{
				ChatRoom: tt.args.req.Message.Chat,
				Sender:   tt.args.req.Message.Sender,
				Body:     tt.args.req.Message.Text,
			}
			tt.mockFn(messageService, mockMsg)

			handler := &messageHandler{
				messageService: messageService,
			}

			resp, err := handler.Send(tt.args.ctx, tt.args.req)
			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(http.StatusOK), resp.Code)
				assert.Equal(t, "OK", resp.Msg)
				assert.NotNil(t, resp)
			}

			messageService.AssertExpectations(t)
		})
	}
}

func TestPull(t *testing.T) {
	type args struct {
		ctx context.Context
		req *rpc.PullRequest
	}
	reverse := false
	tests := []struct {
		name    string
		args    args
		mockFn  func(messageService *mocks.MessageService, mockMsg *models.Request)
		wantErr error
	}{
		{
			name: "success",
			args: args{
				ctx: context.Background(),
				req: &rpc.PullRequest{
					Chat:    "John:Mary",
					Cursor:  0,
					Limit:   10,
					Reverse: &reverse,
				},
			},
			mockFn: func(messageService *mocks.MessageService, mockReq *models.Request) {
				dbResp := &models.Response{
					Messages:   []*models.Message{{}},
					HasMore:    nil,
					NextCursor: nil,
				}
				messageService.On("Pull", mockReq).Return(dbResp, nil)
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			messageService := mocks.NewMessageService(t)
			mockReq := &models.Request{
				ChatRoom: tt.args.req.Chat,
				Cursor:   tt.args.req.Cursor,
				Limit:    tt.args.req.Limit,
				Reverse:  *tt.args.req.Reverse,
			}
			tt.mockFn(messageService, mockReq)

			handler := &messageHandler{
				messageService: messageService,
			}

			resp, err := handler.Pull(tt.args.ctx, tt.args.req)
			if err != nil {
				assert.Error(t, err)
				assert.Equal(t, tt.wantErr.Error(), err.Error())
			} else {
				assert.NoError(t, err)
				assert.Equal(t, int32(http.StatusOK), resp.Code)
				assert.Equal(t, "OK", resp.Msg)
				assert.NotNil(t, resp)
			}

			messageService.AssertExpectations(t)
		})
	}
}
