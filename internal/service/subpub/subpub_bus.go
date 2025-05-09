package subpub

import (
	"context"
	"log/slog"
	"sync"

	"github.com/google/uuid"

	"VK/internal/entity"
	"VK/internal/errs"
	"VK/internal/utils"
)

type SubPubBus struct {
	mu     sync.Mutex
	bus    map[string][]*Subscriber
	closed bool
	done   <-chan struct{}
}

func NewSubPub() *SubPubBus {
	return &SubPubBus{
		mu:     sync.Mutex{},
		bus:    make(map[string][]*Subscriber),
		closed: false,
		done:   make(<-chan struct{}),
	}
}

func (sp *SubPubBus) Subscribe(subject string, cb entity.MessageHandler) (sub entity.Subscription, err error) {
	utils.WithLock(&sp.mu, func() {
		if sp.closed {
			err = errs.ErrSubPubClosed
			return
		}

		s := NewSubscriber(cb, sp, subject)
		sp.bus[subject] = append(sp.bus[subject], s)

		sub = s
	})

	return
}

func (sp *SubPubBus) Publish(subject string, msg interface{}) (err error) {
	idempKey, err := uuid.NewUUID()
	if err != nil {
		return
	}

	var subs []*Subscriber
	utils.WithLock(&sp.mu, func() {
		if sp.closed {
			err = errs.ErrSubPubClosed
			return
		}

		s, ok := sp.bus[subject]
		if !ok {
			err = errs.ErrSubjectNotFound
			return
		}

		subs = s
	})
	for _, sub := range subs {
		sub.Add(msg, idempKey.String())
	}

	return
}

func (sp *SubPubBus) Close(ctx context.Context) (err error) {
	var allSubs []*Subscriber
	slog.Info("Trying to close all subscribers")
	utils.WithLock(&sp.mu, func() {
		if sp.closed {
			err = errs.ErrSubPubClosed

			return
		}

		sp.closed = true

		for key, subs := range sp.bus {
			allSubs = append(allSubs, subs...)
			delete(sp.bus, key)
		}
	})

	var wg sync.WaitGroup
	for _, sub := range allSubs {
		wg.Add(1)

		go func() {
			defer wg.Done()
			close(sub.closeCh)
		}()
	}

	done := make(chan struct{}, 1)

	go func() {
		wg.Wait()

		close(done)
	}()

	select {
	case <-done:
		slog.Info("All subscribers have been closed")
		return
	case <-ctx.Done():
		slog.Info("Close:Context cancelled")
		return ctx.Err()
	}
}
