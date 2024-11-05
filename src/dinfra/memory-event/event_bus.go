package memoryevent

import (
	"context"
	"fmt"
	"sync"

	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
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
	logger = dinfra.LoggerWithStruct(logger, "MemoryEventBus")
	logger.Trace("create")

	return &MemoryEventBus{
		logger: logger,
		lock:   sync.Mutex{},
		items:  map[string][]*SubscribeItem{},
	}
}

func (bus *MemoryEventBus) Subscribe(topic string, handler dinfra.EventHandler) (string, error) {
	logger := bus.logger

	bus.lock.Lock()
	defer bus.lock.Unlock()

	items, exist := bus.items[topic]
	if !exist {
		items = make([]*SubscribeItem, 0)
	}

	id := uuid.NewString()
	items = append(items, &SubscribeItem{
		ID:      id,
		Handler: handler,
	})
	bus.items[topic] = items

	logger.WithFields(logrus.Fields{
		"subscribeID": id,
		"topic":       topic,
	}).Info("subscribe event")
	return id, nil
}

func (bus *MemoryEventBus) Unsubscribe(subscribeID string) error {
	logger := bus.logger.WithField("subscribeID", subscribeID)

	bus.lock.Lock()
	defer bus.lock.Unlock()

	for topic, items := range bus.items {
		for i, item := range items {
			if item.ID == subscribeID {
				bus.items[topic] = append(items[:i], items[i+1:]...)
				logger.WithField("topic", topic).Info("unsubscribe event")
				return nil
			}
		}
	}

	logger.Error("subscribe id not exist")
	return fmt.Errorf("subscribe id [%s] not exist", subscribeID)
}

func (bus *MemoryEventBus) Publish(event *dinfra.Event) (*dinfra.Event, error) {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	logger := bus.logger.WithField("eventID", event.ID)
	logger.WithField("topic", event.Topic).Info("info")

	items, exist := bus.items[event.Topic]
	if !exist {
		logger.Warn("there is no subscriber")
		return event, nil
	}

	go func(items []*SubscribeItem) {

		var wg sync.WaitGroup
		for _, item := range items {
			if item == nil {
				continue
			}
			wg.Add(1)
			go func(item *SubscribeItem) {
				itemLogger := logger.WithField("subscribeID", item.ID)
				defer func() {
					if r := recover(); r != nil {
						itemLogger.WithField("r", r).Error("subscriber handle event crashed")
						wg.Done()
						return
					}

					wg.Done()
				}()

				err := item.Handler(context.TODO(), event)
				logger.WithError(err).Info("subscriber handled")
			}(item)
		}

		wg.Wait()
	}(append([]*SubscribeItem{}, items...))

	return event, nil
}
