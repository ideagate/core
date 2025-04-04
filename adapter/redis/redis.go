package redis

import (
	"context"
	"fmt"
	"time"

	"github.com/bayu-aditya/ideagate/backend/core/ports/distributionlock"
	"github.com/bayu-aditya/ideagate/backend/core/ports/pubsub"
	"github.com/redis/go-redis/v9"
)

type IRedisAdapter interface {
	distributionlock.IDistributionLock
	pubsub.IPubSubAdapter
}

func NewRedisAdapter(conn redis.UniversalClient) IRedisAdapter {
	return &redisAdapter{
		conn: conn,
	}
}

type redisAdapter struct {
	conn redis.UniversalClient
}

func (r *redisAdapter) Lock(ctx context.Context, key string) (isAllow bool, err error) {
	key = fmt.Sprintf("lock:%s", key)
	isAllow, err = r.conn.SetNX(ctx, key, true, 60*time.Second).Result()
	if err != nil {
		return false, err
	}
	return isAllow, nil
}
func (r *redisAdapter) Unlock(ctx context.Context, key string) error {
	key = fmt.Sprintf("lock:%s", key)
	if _, err := r.conn.Del(ctx, key).Result(); err != nil {
		return err
	}
	return nil
}
func (r *redisAdapter) Publish(ctx context.Context, topic string, data []byte) error {
	if _, err := r.conn.Publish(ctx, topic, data).Result(); err != nil {
		return err
	}
	return nil
}
func (r *redisAdapter) Subscribe(ctx context.Context, topic string) (pubsub.ISubscriber, error) {
	subscriber := r.conn.Subscribe(ctx, topic)
	if subscriber == nil {
		return nil, fmt.Errorf("failed to subscribe to topic %s", topic)
	}
	return &subscribe{
		subscriber: subscriber,
		topic:      topic,
		dataChan:   make(chan []byte),
	}, nil
}

type subscribe struct {
	subscriber *redis.PubSub
	topic      string
	dataChan   chan []byte
}

func (s *subscribe) Data(_ context.Context) <-chan []byte {
	go func() {
		for message := range s.subscriber.Channel() {
			s.dataChan <- []byte(message.Payload)
		}
	}()
	return s.dataChan
}
func (s *subscribe) Close() error {
	if err := s.subscriber.Unsubscribe(context.Background(), s.topic); err != nil {
		return err
	}
	if err := s.subscriber.Close(); err != nil {
		return err
	}
	close(s.dataChan)
	return nil
}
