package queue

import "context"

type QueueMock struct {
	PublishFn      func(ctx context.Context, topic string, msg any) error
	PublishCount   int
	SubscribeFn    func(ctx context.Context, topic string, handler HandlerFunc) error
	SubscribeCount int
}

func (mock *QueueMock) Publish(ctx context.Context, topic string, message any) error {
	mock.PublishCount++
	return mock.PublishFn(ctx, topic, message)
}

func (mock *QueueMock) Subscribe(ctx context.Context, topic string, handler HandlerFunc) error {
	mock.SubscribeCount++
	return mock.SubscribeFn(ctx, topic, handler)
}
