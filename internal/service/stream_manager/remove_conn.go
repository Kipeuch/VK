package stream_manager

import (
	generated "VK/internal/generated/service"
	"VK/internal/utils"
)

func (m *StreamManager) RemoveConn(subject string, conn generated.PubSub_SubscribeServer) {
	utils.WithLock(&m.mu, func() {
		for i, v := range m.conns[subject] {
			if v == conn {
				m.conns[subject] = append(m.conns[subject][:i], m.conns[subject][i+1:]...)
			}
		}
	})
}
