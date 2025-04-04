package pubsub

import (
	"context"
	"sync"
	"sync/atomic"
)

type IPubSub interface {
	Publish(ctx context.Context, topic string, data any)
	Subscribe(ctx context.Context, topic, name string, setting SubscribeSetting) ISubscribe
	Close()
}

type ISubscribe interface {
	GetData() <-chan any
	Close()
}

type SubscribeSetting struct {
	NumBufferChan int
}

func New() IPubSub {
	return &pubSub{}
}

type pubSub struct {
	subscribers sync.Map // map[topic]map[subscriberName]*subscriber
	isClosed    atomic.Bool
}

func (p *pubSub) Publish(_ context.Context, topic string, data any) {
	if p.isClosed.Load() {
		return
	}

	subscribersPerTopic, ok := p.subscribers.Load(topic)
	if !ok {
		// if not found, return
		return
	}

	subscribersPerTopicMap := subscribersPerTopic.(*sync.Map)
	subscribersPerTopicMap.Range(func(subscribeName, subscribe interface{}) bool {
		subs := subscribe.(*subscriber)
		if subs.isClosed.Load() {
			return true
		}
		subs.dataChan <- data
		return true
	})
}

func (p *pubSub) Subscribe(_ context.Context, topic, name string, setting SubscribeSetting) ISubscribe {
	if p.isClosed.Load() {
		return nil
	}

	newSubscriber := &subscriber{
		dataChan:  make(chan any, setting.NumBufferChan),
		closeFunc: p.removeSubscriber(topic, name),
	}

	subscribersPerTopic, _ := p.subscribers.LoadOrStore(topic, &sync.Map{})
	subscribersPerTopicMap := subscribersPerTopic.(*sync.Map)
	subscribersPerTopicMap.Store(name, newSubscriber)

	return newSubscriber
}

func (p *pubSub) Close() {
	p.isClosed.Store(true)

	p.subscribers.Range(func(topic, subscribers interface{}) bool {
		subscribersPerTopicMap := subscribers.(*sync.Map)
		subscribersPerTopicMap.Range(func(_, subscribe interface{}) bool {
			subscribe.(*subscriber).Close()
			return true
		})
		return true
	})
}

func (p *pubSub) removeSubscriber(topic, name string) func() {
	return func() {
		subscribersPerTopic, ok := p.subscribers.Load(topic)
		if !ok {
			return
		}

		subscribersPerTopicMap := subscribersPerTopic.(*sync.Map)
		subscribersPerTopicMap.Delete(name)
	}
}

type subscriber struct {
	dataChan  chan any
	isClosed  atomic.Bool
	closeFunc func()
}

func (s *subscriber) GetData() <-chan any {
	return s.dataChan
}

func (s *subscriber) Close() {
	if s.isClosed.Load() {
		return
	}

	s.isClosed.Store(true)
	close(s.dataChan)
	s.closeFunc()
}
