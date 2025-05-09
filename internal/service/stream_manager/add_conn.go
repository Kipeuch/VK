package stream_manager

import (
	generated "VK/internal/generated/service"
	"VK/internal/utils"
)

func (m *StreamManager) AddConn(subject string, conn generated.PubSub_SubscribeServer) {
	utils.WithLock(&m.mu, func() {
		conns, ok := m.conns[subject]
		if !ok {
			m.conns[subject] = []generated.PubSub_SubscribeServer{}
			m.conns[subject] = append(m.conns[subject], conn)
			return
		}
		for _, v := range conns {
			if v == conn {
				return
			}
		}
		m.conns[subject] = append(m.conns[subject], conn)
	})
}
