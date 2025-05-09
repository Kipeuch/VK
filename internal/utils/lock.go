package utils

import (
	"log/slog"
	"sync"
)

func WithLock(mu *sync.Mutex, f func()) {
	mu.Lock()
	defer func() {
		if e := recover(); e != nil {
			slog.Error("Panic:", slog.Any("error", e))
		}

		mu.Unlock()
	}()

	f()
}
