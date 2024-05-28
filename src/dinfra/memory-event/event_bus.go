package memoryevent

import (
	"context"
	"fmt"
	"sync"

	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/google/uuid"
)

type (
	SubscribeItem struct {
		ID      string
		Handler dinfra.EventHandler
	}

	MemoryEventBus struct {
		logger dinfra.Logger
		lock   sync.Mutex
		items  map[string][]*SubscribeItem
	}
)

func NewMemoryEventBus(
	logger dinfra.Logger,
) *MemoryEventBus {
	return &MemoryEventBus{
		logger: logger,
		lock:   sync.Mutex{},
		items:  map[string][]*SubscribeItem{},
	}
}

func (bus *MemoryEventBus) Subscribe(topic string, handler dinfra.EventHandler) (string, error) {

	bus.lock.Lock()
	defer bus.lock.Unlock()

	items, exist := bus.items[topic]
	if !exist {
		items = make([]*SubscribeItem, 1)
	}

	id := uuid.NewString()

	items = append(items, &SubscribeItem{
		ID:      id,
		Handler: handler,
	})

	bus.items[topic] = items
	return id, nil
}

func (bus *MemoryEventBus) Unsubscribe(subscribeID string) error {
	bus.lock.Lock()
	defer bus.lock.Unlock()

	for topic, items := range bus.items {
		for i, item := range items {
			if item.ID == subscribeID {
				bus.items[topic] = append(items[:i], items[i+1:]...)
				return nil
			}
		}
	}

	return fmt.Errorf("subscribe not exist")
}

func (bus *MemoryEventBus) Publish(event *dinfra.Event) (*dinfra.Event, error) {

	items, exist := bus.items[event.Topic]
	if !exist {
		return event, nil
	}

	go func(items []*SubscribeItem) {

		var wg sync.WaitGroup
		for _, item := range items {
			wg.Add(1)
			go func(item *SubscribeItem) {
				defer func() {
					if r := recover(); r != nil {
						wg.Done()
						return
					}

					wg.Done()
				}()

				item.Handler(context.TODO(), event)
			}(item)
		}

		wg.Wait()
	}(append([]*SubscribeItem{}, items...))

	return event, nil
}
