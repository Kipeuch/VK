package rpc

import (
	"context"

	"VK/internal/entity"
	generated "VK/internal/generated/service"
)

type SubPubService interface {
	// Subscribe creates an asynchronous queue subscriber on the given subject.
	Subscribe(subject string, cb entity.MessageHandler) (sub entity.Subscription, err error)

	// Publish publishes the msg argument to the given subject.
	Publish(subject string, msg interface{}) error

	// Close will shutdown sub-pub system.
	// May be blocked by data delivery until the context is canceled.
	Close(ctx context.Context) error
}

type StreamManager interface {
	AddConn(subject string, conn generated.PubSub_SubscribeServer)
	RemoveConn(subject string, conn generated.PubSub_SubscribeServer)
}
