package entity

// MessageHandler is a callback function that processes messages delivered to subscribers.
type MessageHandler func(msg interface{})
