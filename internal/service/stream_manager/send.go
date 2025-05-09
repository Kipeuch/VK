package stream_manager

import (
	"time"

	"VK/internal/entity"
	"VK/internal/errs"
	generated "VK/internal/generated/service"
)

func (m *StreamManager) Send(subject string, msg entity.Message) error {

	m.mu.Lock()
	_, ok := m.IdempotencyKeys[msg.IdempontencyKey]
	if ok {
		m.mu.Unlock()
		return nil
	}

	m.IdempotencyKeys[msg.IdempontencyKey] = time.Now().Add(time.Hour * 24)

	conns, ok := m.conns[subject]
	if !ok {
		m.mu.Unlock()
		return errs.ErrNoConnection
	}
	m.mu.Unlock()

	for _, v := range conns {
		if err := v.Send(&generated.Event{
			Data: msg.Data.(string),
		}); err != nil {
			return err
		}
	}

	return nil
}
