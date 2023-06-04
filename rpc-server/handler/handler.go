package handler

import (
	"context"
	"net/http"

	"github.com/jh-chee/kitewave/rpc-server/kitex_gen/rpc"
	"github.com/jh-chee/kitewave/rpc-server/service"
	"github.com/rs/zerolog/log"
)

type IMServiceInterface interface {
	Send(ctx context.Context, req *rpc.SendRequest) (r *rpc.SendResponse, err error)
	Pull(ctx context.Context, req *rpc.PullRequest) (r *rpc.PullResponse, err error)
}

// IMService implements the last service interface defined in the IDL.
type IMService struct {
	messageService service.MessageService
}

func NewIMService(mssageService service.MessageService) IMServiceInterface {
	return &IMService{
		messageService: mssageService,
	}
}

func (s *IMService) Send(ctx context.Context, req *rpc.SendRequest) (*rpc.SendResponse, error) {
	resp := rpc.NewSendResponse()
	if err := s.messageService.Send(req); err != nil {
		log.Ctx(ctx).Error().Err(err).Msg("fail to send message")
		resp.Code = http.StatusInternalServerError
		return resp, nil
	}

	log.Ctx(ctx).Info().Msg("save to db successful")
	resp.Code, resp.Msg = http.StatusOK, "OK"
	return resp, nil
}

func (s *IMService) Pull(ctx context.Context, req *rpc.PullRequest) (*rpc.PullResponse, error) {
	resp := rpc.NewPullResponse()
	resp.Code, resp.Msg = http.StatusOK, "OK"
	return resp, nil
}
