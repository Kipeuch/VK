package handle_message

type Consumer struct {
	streamManager StreamManager
}

func New(streamManager StreamManager) *Consumer {
	return &Consumer{
		streamManager: streamManager,
	}
}
