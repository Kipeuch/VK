package rpc

import (
	"log/slog"

	generated "VK/internal/generated/service"
)

func (h *Handler) Subscribe(req *generated.SubscribeRequest, stream generated.PubSub_SubscribeServer) (err error) {
	defer func() {
		if err != nil {
			slog.Error("Subscribe:", slog.Any("error", err))
		}
	}()

	if err != nil {
		err = h.getResErr(err)
		return
	}

	slog.Info("Trying to add connection:")
	h.streamManager.AddConn(req.Key, stream)
	slog.Info("Connection added to stream manager:")

	slog.Info("Trying subscribe", slog.Any("key", req.Key))
	subscription, err := h.subPubService.Subscribe(req.Key, h.cb)
	if err != nil {
		err = h.getResErr(err)
		return
	}
	slog.Info("Subscribed successfully", slog.Any("key", req.Key))
	select {
	case <-stream.Context().Done():
		h.streamManager.RemoveConn(req.Key, stream)
		subscription.Unsubscribe()
	}

	return
}
