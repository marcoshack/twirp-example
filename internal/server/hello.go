package server

import (
	"context"
	service "twirp-example/helloworld"

	"github.com/rs/zerolog/log"
)

func (s *HelloWorldServer) Hello(ctx context.Context, req *service.HelloReq) (*service.HelloResp, error) {
	log.Info().Interface("request", req).Msg("Processing request")
	return &service.HelloResp{Text: "Hello " + req.Subject}, nil
}
