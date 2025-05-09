package stream_manager

import (
	"sync"
	"time"

	generated "VK/internal/generated/service"
)

type StreamManager struct {
	mu              sync.Mutex
	conns           map[string][]generated.PubSub_SubscribeServer
	IdempotencyKeys map[string]time.Time
}

func New() *StreamManager {
	return &StreamManager{
		conns:           make(map[string][]generated.PubSub_SubscribeServer),
		IdempotencyKeys: make(map[string]time.Time),
	}
}
