package subpub

import (
	"log/slog"
	"sync"

	"VK/internal/entity"
	"VK/internal/utils"
)

type Subscriber struct {
	handler      entity.MessageHandler
	queue        utils.Queue[entity.Message]
	subPublisher *SubPubBus
	subject      string
	closeCh      chan struct{}
	processing   bool
	mu           sync.Mutex
}

func NewSubscriber(handler entity.MessageHandler, subPublisher *SubPubBus, subject string) *Subscriber {
	sub := &Subscriber{
		handler:      handler,
		queue:        utils.Queue[entity.Message]{},
		subPublisher: subPublisher,
		subject:      subject,
		closeCh:      make(chan struct{}),
		processing:   false,
		mu:           sync.Mutex{},
	}
	return sub
}

func (s *Subscriber) Unsubscribe() {
	utils.WithLock(&s.subPublisher.mu, func() {
		close(s.closeCh)

		subs, ok := s.subPublisher.bus[s.subject]
		if !ok {
			return
		}

		for i, sub := range subs {
			if sub == s {
				s.subPublisher.bus[s.subject] = append(subs[:i], subs[i+1:]...)
				break
			}
		}
	})
}

func (s *Subscriber) Add(msg any, idempKey string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.queue.Push(entity.Message{
		IdempontencyKey: idempKey,
		Data:            msg,
		Subject:         s.subject,
	})
	if !s.processing {
		s.processing = true
		go s.listenQueue()
	}
}

func (s *Subscriber) listenQueue() {
	slog.Info("Trying to handle messages")
	for {
		s.mu.Lock()
		if s.queue.Len() == 0 {
			s.processing = false
			s.mu.Unlock()
			return
		}
		msg, ok := s.queue.Pop()
		if !ok {
			s.mu.Unlock()
			return
		}
		s.mu.Unlock()
		select {
		case <-s.closeCh:
			return
		default:
			s.handler(msg)
		}
	}
}
