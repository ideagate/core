package pubsub

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_pubSub_Publish(t *testing.T) {
	type args struct {
		ctx   context.Context
		topic string
		name  string
		data  any
	}
	tests := []struct {
		name     string
		args     args
		funcTest func(*testing.T, IPubSub, args)
	}{
		{
			name: "without subscriber",
			args: args{
				ctx:   context.TODO(),
				topic: "topic_1",
				data:  "mock_data",
			},
			funcTest: func(t *testing.T, pubSub IPubSub, args args) {
				pubSub.Publish(args.ctx, args.topic, args.data)
			},
		},
		{
			name: "with 1 subscriber but publisher already closed",
			args: args{
				ctx:   context.TODO(),
				topic: "topic_1",
				data:  "mock_data",
			},
			funcTest: func(t *testing.T, pubSub IPubSub, args args) {
				wg := sync.WaitGroup{}

				var result1 any

				// This subscriber will not get the data
				go func() {
					defer wg.Done()
					wg.Add(1)
					subscribe := pubSub.Subscribe(args.ctx, args.topic, "subscriber_1", SubscribeSetting{
						NumBufferChan: 1,
					})
					result1 = <-subscribe.GetData()
					subscribe.Close()
				}()

				time.Sleep(1 * time.Second)
				pubSub.Close()
				pubSub.Publish(args.ctx, args.topic, args.data)

				wg.Wait()
				assert.Equal(t, nil, result1)
			},
		},
		{
			name: "with 2 subscriber",
			args: args{
				ctx:   context.TODO(),
				topic: "topic_1",
				data:  "mock_data",
			},
			funcTest: func(t *testing.T, pubSub IPubSub, args args) {
				wg := sync.WaitGroup{}

				var (
					result1 any
					result2 any
				)

				// This subscriber will get the data
				go func() {
					defer wg.Done()
					wg.Add(1)
					subscribe := pubSub.Subscribe(args.ctx, args.topic, "subscriber_1", SubscribeSetting{
						NumBufferChan: 1,
					})
					result1 = <-subscribe.GetData()
					subscribe.Close()
				}()

				// This subscriber will not get the data because it is closed
				go func() {
					defer wg.Done()
					wg.Add(1)
					subscribe := pubSub.Subscribe(args.ctx, args.topic, "subscriber_2", SubscribeSetting{
						NumBufferChan: 1,
					})
					subscribe.Close()
					result2 = <-subscribe.GetData()
				}()

				time.Sleep(1 * time.Second)
				pubSub.Publish(args.ctx, args.topic, args.data)
				pubSub.Close()

				wg.Wait()
				assert.Equal(t, "mock_data", result1)
				assert.Equal(t, nil, result2)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := New()

			tt.funcTest(t, p, tt.args)
		})
	}
}
