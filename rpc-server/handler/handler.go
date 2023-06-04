package handler

import (
	"context"
	"net/http"

	"github.com/jh-chee/kitewave/rpc-server/kitex_gen/rpc"
	"github.com/jh-chee/kitewave/rpc-server/models"
	"github.com/jh-chee/kitewave/rpc-server/service"
	"github.com/rs/zerolog/log"
)

type MessageHandler interface {
	Send(ctx context.Context, req *rpc.SendRequest) (r *rpc.SendResponse, err error)
	Pull(ctx context.Context, req *rpc.PullRequest) (r *rpc.PullResponse, err error)
}

// IMService implements the last service interface defined in the IDL.
type messageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(mssageService service.MessageService) MessageHandler {
	return &messageHandler{
		messageService: mssageService,
	}
}

func (s *messageHandler) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()
	msg := &models.Message{
		Sender:   req.Message.Sender,
		ChatRoom: req.Message.Chat,
		Body:     req.Message.Text,
	}

	if err := s.messageService.Send(msg); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("fail to send message")
		resp.Code = http.StatusInternalServerError
		return resp, err
	}

	log.Ctx(ctx).Info().Msg("save to db successful")
	resp.Code, resp.Msg = http.StatusOK, "OK"
	return resp, nil
}

func (s *messageHandler) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = http.StatusOK, "OK"
	return resp, nil
}
