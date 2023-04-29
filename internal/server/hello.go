package server

import (
	"context"
	"fmt"

	service "github.com/marcoshack/twirp-example/rpc/helloworld"

	"github.com/rs/zerolog/log"
)

func (s *HelloWorldServer) Hello(ctx context.Context, req *service.HelloReq) (*service.HelloResp, error) {
	log.Info().Interface("request", req).Msg("Processing request")

	return &service.HelloResp{Text: fmt.Sprintf("Hello, %s!", req.Subject)}, nil
}
