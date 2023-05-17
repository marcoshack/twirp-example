package server

import (
	"context"
	"time"

	"github.com/marcoshack/twirp-example/internal/storage"
	service "github.com/marcoshack/twirp-example/rpc/helloworld"
)

func (s *HelloServer) GetAll(ctx context.Context, req *service.GetAllReq) (*service.GetAllResp, error) {
	output, err := s.dao.GetAll(ctx, &storage.GetAllInput{
		Limit: req.Limit,
	})
	if err != nil {
		return nil, err
	}

	items := make([]*service.HelloItem, 0, len(output.Entries))
	for _, entry := range output.Entries {
		items = append(items, &service.HelloItem{
			Id:        entry.ID,
			CreatedAt: entry.CreatedAt.Format(time.RFC3339),
			Message:   entry.Message,
		})
	}

	return &service.GetAllResp{
		Items: items,
		Size:  int32(len(output.Entries)),
	}, nil
}
