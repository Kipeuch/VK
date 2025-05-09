package handle_message

import (
	"log/slog"

	"VK/internal/entity"
)

func (c *Consumer) Consume(msg interface{}) {
	convertedMsg, ok := msg.(entity.Message)
	if !ok {
		slog.Error("failed to cast message to entity.Message")
		return
	}

	if err := c.streamManager.Send(convertedMsg.Subject, convertedMsg); err != nil {
		slog.Error("failed to send message:", slog.Any("idpKey", convertedMsg.IdempontencyKey), slog.Any("data", convertedMsg.Data), slog.Any("err", err))
		return
	}
}
