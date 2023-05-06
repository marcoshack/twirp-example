package server

import (
	"context"
	"fmt"

	"github.com/marcoshack/twirp-example/internal/storage"
	service "github.com/marcoshack/twirp-example/rpc/helloworld"
	"github.com/twitchtv/twirp"

	"github.com/rs/zerolog/log"
)

func (s *HelloServer) Hello(ctx context.Context, req *service.HelloReq) (*service.HelloResp, error) {
	log.Info().Interface("request", req).Msg("Processing request")
	entry, err := s.dao.AddHelloWorld(ctx, &storage.HelloInput{
		Message: req.Subject,
	})
	if err != nil {
		log.Error().Err(err).Msg("failed to add hello world")
		return nil, twirp.InternalError("failed to add hello world")
	}
	return &service.HelloResp{Text: fmt.Sprintf("Hello, %s! Your id is %s", req.Subject, entry.ID)}, nil
}
