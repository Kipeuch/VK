package stream_manager

import "time"

func (m *StreamManager) ClearIdpMap() {
	m.mu.Lock()
	for i, v := range m.IdempotencyKeys {
		if time.Now().After(v) {
			delete(m.IdempotencyKeys, i)
		}
	}
	m.mu.Unlock()
}
