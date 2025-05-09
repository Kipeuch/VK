package handle_message

import "VK/internal/entity"

type StreamManager interface {
	Send(subject string, msg entity.Message) error
}
