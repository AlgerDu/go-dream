package memoryevent

import (
	"context"
	"fmt"
	"strings"
	"sync"

	"github.com/AlgerDu/go-dream/src/dinfra"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
)

type (
	MemoryEventBus struct {
		logger dinfra.Logger
		lock   sync.Mutex

		root *Node
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
		root:   BuildeNode(""),
	}
}

func (bus *MemoryEventBus) Subscribe(
	topic string,
	handler dinfra.EventHandler,
) (string, error) {
	logger := bus.logger

	keys := strings.Split(topic, dinfra.TopicSeparators)
	if len(keys) == 0 {
		return "", dinfra.ErrInvalidTopic
	}

	bus.lock.Lock()
	defer bus.lock.Unlock()

	node, keys := FindChild(bus.root, keys)

	if len(keys) > 0 {
		for i, key := range keys {
			if key == "#" && i != len(keys)-1 {
				return "", fmt.Errorf("# must at last")
			}

			child := BuildeNode(key)
			child.Topic = fmt.Sprintf("%s/%s", node.Topic, key)

			node.Children[key] = child
			node = child
		}
	}

	id := uuid.NewString()
	record := &HandlerRecord{
		ID:      id,
		Handler: handler,
	}

	node.Handlers[id] = record

	logger.WithFields(logrus.Fields{
		"subscribeID": id,
		"topic":       topic,
	}).Info("subscribe event")
	return id, nil
}

func (bus *MemoryEventBus) Unsubscribe(
	subscribeID string,
) error {
	logger := bus.logger.WithField("subscribeID", subscribeID)

	bus.lock.Lock()
	defer bus.lock.Unlock()

	record := RemoveRecord(bus.root, subscribeID)
	if record == nil {
		logger.Error("subscribe id not exist")
		return fmt.Errorf("subscribe id [%s] not exist", subscribeID)
	}

	return nil
}

func (bus *MemoryEventBus) Publish(
	context context.Context,
	event *dinfra.Event,
) (*dinfra.Event, error) {
	if event.ID == "" {
		event.ID = uuid.NewString()
	}

	logger := bus.logger.WithField("eventID", event.ID)
	logger.WithField("topic", event.Topic).Info("info")

	records := MacthHandlers(bus.root, strings.Split(event.Topic, dinfra.TopicSeparators))

	if len(records) == 0 {
		logger.Warn("there is no subscriber")
		return event, nil
	}

	go func(records []*HandlerRecord) {

		var wg sync.WaitGroup
		for _, record := range records {
			if record == nil {
				continue
			}
			wg.Add(1)

			go func(record *HandlerRecord) {
				itemLogger := logger.WithField("subscribeID", record.ID)
				defer func() {
					if r := recover(); r != nil {
						itemLogger.WithField("r", r).Error("subscriber handle event crashed")
						wg.Done()
						return
					}

					wg.Done()
				}()

				err := record.Handler(context, event)
				logger.WithError(err).Info("subscriber handled")
			}(record)
		}

		wg.Wait()
		logger.Info("event handle end")

	}(records)

	return event, nil
}
