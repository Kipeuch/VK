package rpc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"VK/internal/entity"
	"VK/internal/errs"
	generated "VK/internal/generated/service"
)

type Handler struct {
	subPubService SubPubService
	streamManager StreamManager
	cb            entity.MessageHandler

	generated.UnimplementedPubSubServer
}

func New(subPubService SubPubService, streamManager StreamManager, cb entity.MessageHandler) *Handler {
	return &Handler{
		subPubService: subPubService,
		streamManager: streamManager,
		cb:            cb,
	}
}

func (h *Handler) getResErr(err error) error {
	switch {
	case errors.Is(err, errs.ErrSubjectNotFound):
		return status.Errorf(codes.NotFound, err.Error())
	case errors.Is(err, errs.ErrSubPubClosed):
		return status.Errorf(codes.Unavailable, err.Error())
	case errors.Is(err, errs.ErrNoConnection):
		return status.Errorf(codes.Unavailable, err.Error())
	default:
		return err
	}
}
