package queue

import "context"

type Message struct {
	Headers map[string]string
	Data    []byte
}

type HandlerFunc func(msg Message)

type Publisher interface {
	Publish(ctx context.Context, topic string, msg any) error
}

type Subscriber interface {
	Subscribe(ctx context.Context, topic string, handler HandlerFunc) error
}
