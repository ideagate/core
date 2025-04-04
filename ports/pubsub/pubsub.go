package pubsub

import "context"

type IPubSubAdapter interface {
	Publish(ctx context.Context, topic string, data []byte) error
	Subscribe(ctx context.Context, topic string) (ISubscriber, error)
}

type ISubscriber interface {
	Data(ctx context.Context) <-chan []byte
	Close() error
}
