package main

import (
	"context"
	"time"

	"github.com/AlgerDu/go-dream/src/dinfra"
	extlogrus "github.com/AlgerDu/go-dream/src/dinfra/ext-logrus"
	memoryevent "github.com/AlgerDu/go-dream/src/dinfra/memory-event"
)

type (
	Book struct {
		ID      string
		Content string
	}
)

func main() {
	logger := extlogrus.New(&extlogrus.LoggerOptions{MinLevel: extlogrus.Trace})
	bus := memoryevent.NewMemoryEventBus(logger)

	bus.Subscribe(
		"/book/+/update",
		func(context context.Context, event *dinfra.Event) error {
			book, err := dinfra.ConvertEventDataTo[*Book](event)
			if err != nil {
				return err
			}

			logger.WithField("book", book).Trace("handler topic [/book/+/update]")
			return nil
		},
	)

	dinfra.PublishEvent(bus, "/book/1/update", &Book{ID: "1"})

	time.Sleep(time.Second * 3)
}
