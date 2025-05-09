package rpc

import (
	"context"
	"log/slog"

	"github.com/golang/protobuf/ptypes/empty"

	generated "VK/internal/generated/service"
)

func (h *Handler) Publish(ctx context.Context, req *generated.PublishRequest) (e *empty.Empty, err error) {
	defer func() {
		if err != nil {
			slog.Error("Publish:", slog.Any("error", err))
		}
	}()
	slog.Info("Trying publish:", slog.Any("Key", req.Key), slog.Any("Data", req.Data))
	if err = h.subPubService.Publish(req.Key, req.Data); err != nil {
		err = h.getResErr(err)
		return
	}
	slog.Info("Publish successfully:", slog.Any("Key", req.Key), slog.Any("Data", req.Data))
	return
}
