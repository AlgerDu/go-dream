package dinfra

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
)

var (
	trackContextKey contextKey = "TrackID"
	TopicSeparators            = "/"

	ErrInvalidTopic = errors.New("invalid topic")
)

type (
	contextKey string
	EventName  string // 事件名称，方便进行过滤处理

	// 事件
	Event struct {
		ID       string // 可追踪 ID ；发布时生成
		Name     EventName
		Topic    string // 主题
		Data     any    // 携带的数据，类型可能是 JSON 或者结构体
		CreateAt int64  // 创建时间（时间戳，单位 ms）
	}

	// 事件处理
	EventHandler func(context context.Context, event *Event) error

	// 事件总线
	EventBus interface {
		Subscribe(topic string, handler EventHandler) (string, error)  // 订阅事件，返回订阅 ID ，可用于取消订阅
		Unsubscribe(subscribeID string) error                          // 通过订阅 ID 取消订阅
		Publish(context context.Context, event *Event) (*Event, error) // 发布事件
	}
)

func PublishEvent(
	eventBus EventBus,
	topic string,
	data any,
) (*Event, error) {
	return eventBus.Publish(
		context.TODO(),
		&Event{
			ID:       "",
			Topic:    topic,
			Data:     data,
			CreateAt: time.Now().UnixMilli(),
		},
	)
}

func PublishEvent2(
	eventBus EventBus,
	ctx context.Context,
	topic string,
	data any,
) (*Event, error) {

	id := uuid.NewString()

	return eventBus.Publish(
		context.WithValue(ctx, trackContextKey, id),
		&Event{
			ID:       id,
			Topic:    topic,
			Data:     data,
			CreateAt: time.Now().UnixMilli(),
		},
	)
}

// 将 event 携带的数据转换为结构体
func ConvertEventDataTo[DataType any](event *Event) (DataType, error) {
	if data, ok := event.Data.(DataType); ok {
		return data, nil
	}

	var data DataType
	if jsonStr, ok := event.Data.(string); ok {
		err := json.Unmarshal([]byte(jsonStr), &data)
		return data, err
	}

	return data, fmt.Errorf("convert event data err")
}
